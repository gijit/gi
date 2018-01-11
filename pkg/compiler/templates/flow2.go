package main

import "fmt"

// a defer can recover a panic value, and re-panic with
// a new value that the next defer in line receives.

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
			switch r := recov.(type) {
			case int:
				ret1 = ret1 + 9 + r*199
				ret0 = ret0 + 19 + r*3
			}
		}(a)

		defer func() {
			fmt.Printf("second defer running, a=%v, b=%v, ret1=%v\n", a, b, ret1)
			b = b * 7
			ret0 = ret0 + 100

			recov := recover()
			fmt.Printf("defer 2 recovered '%v'\n", recov)
			switch r := recov.(type) {
			case int:
				panic(r + 17)
			}
			ret1 = ret1 + 100
			fmt.Printf("second defer just updated ret1 to %v\n", ret1)

		}()

		a = 1
		b = 1

		panic(a + b)

		return b, 58
	}
	f1, f2 := f()
	fmt.Printf("f1 = %v, f2 = %v\n", f1, f2)
	/*
		go run flow2.go
		second defer running, a=1, b=1, ret1=0
		defer 2 recovered '2'
		first defer running, a=0, b=7, ret0=100, ret1=0
		defer 1 recovered '19'
		f1 = 379, f2 = 3791
	*/
}
