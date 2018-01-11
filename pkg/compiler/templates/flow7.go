package main

import "fmt"

// nested calls: can panic transfer
// from lower in the stack to higher up?

var global int = 0

func deeper(a int) {
	global += (a + 3)
	panic("panic-in-deeper")
}

func main() {
	a := 0
	b := 0

	f := func() (ret0, ret1 int) {

		defer func(a int) {
			fmt.Printf("first defer running, a=%v, b=%v, ret0=%v, ret1=%v\n", a, b, ret0, ret1)
			b = b + 3
			ret0 = (ret0+1)*3 + global
			ret1 = ret1 + 1
			recov := recover()
			fmt.Printf("defer 1 recovered '%v'\n", recov)
			if recov != nil {
				ret1 = ret1 + 9 + global
				ret0 = ret0 + 19 + global
			}
		}(a)

		deeper(a)
		return
	}
	f1, f2 := f()
	fmt.Printf("f1 = %v, f2 = %v\n", f1, f2)
	/*

	go run flow7.go
	first defer running, a=0, b=0, ret0=0, ret1=0
	defer 1 recovered 'panic-in-deeper'
	f1 = 28, f2 = 13

	*/
}
