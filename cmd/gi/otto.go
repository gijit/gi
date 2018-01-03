package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/go-interpreter/gi/pkg/compiler"
	"github.com/robertkrimen/otto"
)

func OttoReplMain() {

	// and then eval!
	vm := otto.New()
	inc := compiler.NewIncrState()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("gi> ")
		src, isPrefix, err := reader.ReadLine()
		if err == io.EOF {
			fmt.Printf("[EOF]\n")
			return
		}
		panicOn(err)
		if isPrefix {
			panic("line too long")
		}

		translation := inc.Tr([]byte(src))
		fmt.Printf("go:'%s'  -->  '%s' in js\n", src, translation)

		v, err := vm.Eval(string(translation))
		panicOn(err)
		fmt.Printf("v back = '%#v'\n", v)
	}
}
