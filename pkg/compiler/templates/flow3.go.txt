package main

import "fmt"

// one defer throws a different panic, a later
//  defer recovers completely.

func main() {
	a := 0
	b := 0
	f := func() (ret0, ret1 int) {

		defer func(a int) {
			fmt.Printf("first defer running, a=%v, b=%v\n", a, b)
			b = b + 3
			ret0 = (ret0 + 1) * 3
			ret1 = ret1 + 1
			recov := recover()
			fmt.Printf("defer 1 recovered '%v'\n", recov)
		}(a)

		defer func() {
			fmt.Printf("second defer running, a=%v, b=%v\n", a, b)
			b = b * 7
			ret0 = ret0 + 100
			ret1 = ret1 + 100
			recov := recover()
			fmt.Printf("second defer, recov is %v\n", recov)
			panic("panic-in-defer-2")
		}()

		a = 1
		b = 1

		panic("ouch")

		return b, 58
	}
	f1, f2 := f()
	fmt.Printf("f1 = %v, f2 = %v\n", f1, f2)
	/*
		go run flow3.go
		second defer running, a=1, b=1
		second defer, recov is ouch
		first defer running, a=0, b=7
		defer 1 recovered 'panic-in-defer-2'
		f1 = 303, f2 = 101

	*/
}
