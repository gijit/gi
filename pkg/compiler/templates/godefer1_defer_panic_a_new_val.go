package main

import "fmt"

func main() {
	a := 0
	b := 0
	f := func(ret0, ret1 int) {
		defer func(a int) {
			fmt.Printf("defer 1 running.\n")
			if r := recover(); r != nil {
				fmt.Printf("defer 1 recovered '%v'\n", r)
			}
		}(a)

		defer func() {
			fmt.Printf("defer 2 running.\n")
			//	if r := recover(); r != nil {
			//		fmt.Printf("defer 2 recovered '%v', throwing something new\n", r)
			panic("panic-in-defer-2")
			//	}
		}()

		panic("ouch")
	}
	f()
}
