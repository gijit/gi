package main

import "fmt"

func main() {
	a := 0
	b := 0
	f := func() (ret0, ret1 int) {

		defer func(a int) {
			fmt.Printf("first defer running, a=%v, b=%v, ret0=%v, ret1=%v\n", a, b, ret0, ret1)
			b = b + 3
			ret0 = (ret0 + 1) * 3
			ret1 = ret1 + 1
			recov := recover()
			fmt.Printf("defer 1 recovered '%v'\n", recov)
			if recov != nil {
				ret1 = ret1 + 9
				ret0 = ret0 + 19
			}
		}(a)

		defer func() {
			fmt.Printf("second defer running, a=%v, b=%v, ret1=%v\n", a, b, ret1)
			b = b * 7
			ret0 = ret0 + 100
			ret1 = ret1 + 100
			fmt.Printf("second defer just updated ret1 to %v\n", ret1)
			//recov := recover()
			//fmt.Printf("second defer, recov is %v\n", recov)
			//panic("panic-in-defer-2")
		}()

		a = 1
		b = 1

		panic("ouch")

		return b, 58
	}
	f1, f2 := f()
	fmt.Printf("f1 = %v, f2 = %v\n", f1, f2)
	/*
	go run flow6.go
	second defer running, a=1, b=1, ret1=0
	second defer just updated ret1 to 100
	first defer running, a=0, b=7, ret0=100, ret1=100
	defer 1 recovered 'ouch'
	f1 = 322, f2 = 110
	*/
}
