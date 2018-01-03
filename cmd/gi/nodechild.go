package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/glycerine/idem"
	"github.com/go-interpreter/gi/pkg/compiler"
)

func NodeChildMain() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	chld := exec.CommandContext(ctx, "node", "-i")

	writerHalt := idem.NewHalter()
	readerHalt := idem.NewHalter()

	chldOut, err := chld.StdoutPipe()
	panicOn(err)
	chldIn, err := chld.StdinPipe()
	panicOn(err)

	chldOutFd, ok := chldOut.(*os.File)
	if !ok {
		panic("could not get *os.File from chldOut")
	}

	//chldPty, err := pty.Start(chld)
	err = chld.Start()
	panicOn(err)
	fmt.Printf("spawned node as process %v\n", chld.Process.Pid)

	readAvail := make(chan string, 10)

	// reader: the PRINT part of the REPL.
	go func() {
		defer func() {
			readerHalt.MarkDone()
			//fmt.Printf("reader done\n")
		}()

		var buf [4096]byte
		alldone := false

		for !alldone {
			totr := 0

		oneRead:
			// try reading until we get a newline
			for totr <= 0 {
				//fmt.Printf("<<< reader about to read from chldOut\n")
				// go1.10 only:
				// err := chldOutFd.SetReadDeadline(time.Now().Add(50 * time.Millisecond))
				panicOn(err)
				nr, err := chldOutFd.Read(buf[:])
				totr += nr
				if totr > 0 && (buf[totr-1] == '\n' || bytes.Contains(buf[:totr], []byte{'\n'})) {
					break oneRead
				}
				if err == io.EOF {
					//fmt.Printf("<<< reader sees EOF. shutting down\n")
					// wrap things up
					alldone = true
					break oneRead
				}
				// timeout errors are common, don't freakout.
				select {
				case <-time.After(time.Millisecond * 10):
				case <-writerHalt.ReqStop.Chan:
					//fmt.Printf("<<< reader sees writeHalt.ReqStop, shutting down\n")
					alldone = true
					break oneRead
				case <-readerHalt.ReqStop.Chan:
					//fmt.Printf("<<< reader sees readerHalt.ReqStop, shutting down\n")
					alldone = true
					break oneRead
				}
				//fmt.Printf("<<< reader sleep 500 msec\n")
			} // end inner oneRead loop

			if totr > 0 {
				s := string(buf[:totr])
				select {
				case readAvail <- s:
				case <-writerHalt.ReqStop.Chan:
					//fmt.Printf("<<< reader sees writeHalt.ReqStop, shutting down\n")
					return
				case <-readerHalt.ReqStop.Chan:
					//fmt.Printf("<<< reader sees readerHalt.ReqStop, shutting down\n")
					return
				}
			}
		}
		//fmt.Printf("<<< reader sees alldone. returning.\n")
	}()

	StartRepl(writerHalt, readerHalt, readAvail, chldIn)

	readerHalt.ReqStop.Close()
	cancelFunc()
	chld.Wait()
	//<-readerHalt.Done.Chan
}

func StartRepl(writerHalt *idem.Halter, readerHalt *idem.Halter, readAvail chan string, chldIn io.WriteCloser) {
	defer func() {
		writerHalt.RequestStop()
		writerHalt.MarkDone()
	}()
	reader := bufio.NewReader(os.Stdin)

	// the translator under test!
	inc := compiler.NewIncrState()

	// writer: the READ and EVAL parts of the REPL
	for {

		// sync up by reading lines until we get node's prompt
	syncLoop:
		for {
			select {
			case s := <-readAvail:
				fmt.Printf("%s", s)
				if strings.HasSuffix(s, "gi> ") {
					//fmt.Printf("we see the prompt, writer exiting wait-for-prompt top loop\n")
					break syncLoop
				}
			case <-readerHalt.ReqStop.Chan:
				return
			case <-writerHalt.ReqStop.Chan:
				return
			}
		}
		//fmt.Printf("done with top loop, asking for line from human\n")

		// read a line from our user
		src, isPrefix, err := reader.ReadLine()
		if err == io.EOF {
			fmt.Printf("[EOF]\n")
			return
		}
		panicOn(err)
		if isPrefix {
			panic("line too long")
		}

		//fmt.Printf("got line from line from human : '%s'\n", string(src))

		translation, err := translateAndCatchPanic(inc, src)
		if err != nil {
			fmt.Printf("oops: '%v' on input '%s'\n", err, string(src))
			translation = "\n"
			// still write, so we get another prompt
		}

		//fmt.Printf("got translation of line from Go into js: '%s'\n", string(translation))

		// write
		totw := len(translation)
		if totw > 0 {
			if translation[totw-1] != '\n' {
				translation += "\n"
				totw++
			}
		} else {
			translation = "\n"
			totw = 1
		}

		nw := 0
		for nw < totw {
			n, err := chldIn.Write([]byte(translation))
			nw += n
			panicOn(err)
			if nw < totw {
				time.Sleep(10 * time.Millisecond)
			}
		}
		//fmt.Printf("wrote to node a %v character string '%s'\n", totw, translation)
	}
}
