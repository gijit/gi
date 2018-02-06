package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/gijit/gi/pkg/compiler"
	"github.com/gijit/gi/pkg/verb"
)

func translateAndCatchPanic(inc *compiler.IncrState, src []byte) (translation string, err error) {
	defer func() {
		recov := recover()
		if recov != nil {
			msg := fmt.Sprintf("problem detected during Go static type checking: '%v'", recov)
			if verb.Verbose {
				msg += fmt.Sprintf("\n%s\n", string(debug.Stack()))
			}
			err = fmt.Errorf(msg)
		}
	}()
	ssrc := string(src)
	pp("about to translate Go source '%s'", ssrc)

	translation = string(inc.Tr([]byte(src)))

	t2 := strings.TrimSpace(translation)
	nt2 := len(t2)
	if nt2 > 0 {
		if t2[nt2-1] == '\n' {
			t2 = t2[:nt2-1]
		}
	}
	p("go:'%s'  -->  '%s'\n", src, t2)
	return
}

func readHistory(histFn string) (history []string, err error) {
	if !FileExists(histFn) {
		return nil, nil
	}
	by, err := ioutil.ReadFile(histFn)
	if err != nil {
		return nil, err
	}
	splt := strings.Split(string(by), string(byte(0)))
	n := len(splt)

	// avoid returning an extra blank history line
	// at the end of the history file.
	if n > 0 && strings.TrimSpace(splt[n-1]) == "" {
		return splt[:n-1], nil
	}
	return splt, nil
}

var zeroByte = string(byte(0))

func removeCommands(history []string, histFn string, histFile *os.File, rms string) (history2 []string, histFile2 *os.File, beg int, end int, err error) {

	beg = -1
	end = -1
	history2 = history
	histFile2 = histFile
	var num []int

	num, err = getHistoryRange(rms, history)
	if err != nil {
		return
	}

	switch len(num) {
	case 1:
		k := num[0]
		beg = k - 1
		end = k - 1
		fmt.Printf("remove history %03d.\n", k)
		history2 = append(history[:k-1], history[k:]...)
	case 2:
		if num[1] < num[0] {
			err = fmt.Errorf("bad remove history request, end before beginning.")
			return
		}
		beg = num[0] - 1
		end = num[1] - 1
		fmt.Printf("remove history %03d - %03d.\n", num[0], num[1])
		history2 = append(history[:num[0]-1], history[num[1]:]...)
	}

	histFile.Close()
	os.Remove(histFn)
	histFile2, err = os.OpenFile(histFn,
		os.O_WRONLY|os.O_CREATE|os.O_APPEND|os.O_SYNC,
		0600)
	panicOn(err)
	// print new history to file
	for i := range history2 {
		fmt.Fprintf(histFile2, "%s%s", history2[i], zeroByte)
	}
	return
}

func getHistoryRange(lows string, history []string) (slc []int, err error) {
	parts := strings.Split(lows, "-")
	if len(parts) > 2 {
		return nil, fmt.Errorf("bad history range request, more than one '-' found.")
	}
	num := make([]int, len(parts))
	for i := range parts {
		s := strings.TrimSpace(parts[i])
		if s == "" {
			// allow ":rm -4" to indicate "from the beginning through 4"
			// and ":rm 4-" to mean "from 4 until the end".

			if i == 0 {
				num[i] = 1
			} else {
				num[i] = len(history)
			}
		} else {
			num[i], err = strconv.Atoi(s)
			if err != nil {
				return nil, fmt.Errorf("bad history request, could "+
					"not convert '%v' to integer.\n", s)
			}
		}
		if num[i] < 1 || num[i] > len(history) {
			return nil, fmt.Errorf("bad history request, out of range.\n")
		}
	}
	return num, nil
}

func sourceGoFiles(files []string) ([]byte, error) {
	var buf bytes.Buffer
	for _, f := range files {
		fd, err := os.Open(f)
		if err != nil {
			return nil, err
		}
		defer fd.Close()
		by, err := ioutil.ReadAll(fd)
		if err != nil {
			return nil, err
		}
		_, err = io.Copy(&buf, bytes.NewBuffer(by))
		if err != nil {
			return nil, err
		}
		// zero
		fmt.Fprintf(&buf, zeroByte)
	}
	bb := buf.Bytes()
	fmt.Printf("sourceGoFiles() returning '%s'\n", string(bb))
	return bb, nil
}
