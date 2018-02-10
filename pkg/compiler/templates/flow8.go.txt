package main

import "fmt"

// nested calls: can panic transfer
// from lower in the stack to higher up?

var global int = 8

func deeper(a int) {
	global += (a + 3)
	panic("panic-in-deeper")
}

func intermed(a int) {
	deeper(a)
}

func main() {
	a := 0
	b := 0

	f := func() (ret0, ret1 int) {

		defer func(a int) {
			fmt.Printf("first defer running, a=%v, b=%v, ret0=%v, ret1=%v, global=%v\n",
				a, b, ret0, ret1, global)
			b = b + 3
			ret0 = (ret0+1)*3 + global
			fmt.Printf("debug, after ret0 = (ret0+1) * 3 + global: ret0=%v", ret0)
			ret1 = ret1 + 1
			recov := recover()
			fmt.Printf("defer 1 recovered '%v'\n", recov)
			if recov != nil {
				fmt.Printf("recov was not nil, ret0=%v  and ret1=%v", ret0, ret1)
				ret1 = ret1 + 9 + global
				ret0 = ret0 + 19 + global
			}
			fmt.Printf("at end of 1st defer, ret0='%v', ret1='%v'\n", ret0, ret1)
		}(a)

		intermed(a)
		return
	}
	f1, f2 := f()
	fmt.Printf("f1 = %v, f2 = %v\n", f1, f2)
	/*
	go run flow8.go
	first defer running, a=0, b=0, ret0=0, ret1=0, global=11
	debug, after ret0 = (ret0+1) * 3 + global: ret0=14defer 1 recovered 'panic-in-deeper'
	recov was not nil, ret0=14  and ret1=1at end of 1st defer, ret0='44', ret1='21'
	f1 = 44, f2 = 21
	*/
}
