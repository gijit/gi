package main

import "fmt"

// defers access and modify named return values.

func main() {
	a := 7
	b := 0
	f := func() (ret0, ret1 int) {

		defer func(a int) {
			fmt.Printf("first defer, a = %v, ret0=%v, ret1 = %v\n", a, ret0, ret1)
			ret0 = (ret0+1)*3 + a
			ret1 = ret1 + 1 + a
		}(a)

		defer func() {
			ret0 = ret0 + 100
			ret1 = ret1 + 100
		}()

		a = 1
		b = 1

		return b, a + 58
	}
	f1, f2 := f()
	fmt.Printf("f1 = %v, f2 = %v\n", f1, f2)
	/*
	first defer, a = 7, ret0=101, ret1 = 159
	f1 = 313, f2 = 167
	*/

}
