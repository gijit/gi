package main

// what happens to the function if the first defer
// halts the panic, what does it return?  the zero value
// for its type, apparently.

import "fmt"

func main() {
	a := 0
	b := 0
	f := func() int {

		defer func() {
			fmt.Printf("first defer running, a=%v, b=%v\n", a, b)
		}()

		defer func() {
			fmt.Printf("second defer running, a=%v, b=%v\n", a, b)
			recover()
		}()

		panic("ouch")
		return 7
	}
	f1 := f()
	fmt.Printf("f1 = %v\n", f1)
	/*

		go run flow11.go
		second defer running, a=0, b=0
		first defer running, a=0, b=0
		f1 = 0

	*/
}
