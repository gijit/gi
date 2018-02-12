package main

import "fmt"

func main() {
	defer func() {
		fmt.Printf("defer 1 running.\n")
		if r := recover(); r != nil {
			fmt.Printf("defer 1 recovered '%v'\n", r)
		}
	}()

	defer func() {
		fmt.Printf("defer 2 running.\n")
		//	if r := recover(); r != nil {
		//		fmt.Printf("defer 2 recovered '%v', throwing something new\n", r)
		panic("panic-in-defer-2")
		//	}
	}()

	panic("ouch")
}
