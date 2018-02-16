package compiler

import (
	"fmt"
	"math/cmplx"
	"testing"

	cv "github.com/glycerine/goconvey/convey"
)

func Test500ComplexNumbers(t *testing.T) {

	cv.Convey(`run through Go stdlib math/cmplx to get some known good values to test complex.lua against`, t, func() {

		input := []complex128{
			(1.0017679804707456328694569 - 2.9138232718554953784519807i),
			(0.03606427612041407369636057 + 2.7358584434576260925091256i),
			(1.6249365462333796703711823 + 2.3159537454335901187730929i),
			(2.0485650849650740120660391 - 3.0795576791204117911123886i),
			(0.29621132089073067282488147 - 3.0007392508200622519398814i),
			(1.0664555914934156601503632 - 2.4872865024796011364747111i),
			(0.48681307452231387690013905 - 2.463655912283054555225301i),
			(0.6116977071277574248407752 - 1.8734458851737055262693056i),
			(1.3649311280370181331184214 + 2.8793528632328795424123832i),
			(2.6189310485682988308904501 - 2.9956543302898767795858704i),
		}

		a := input[0]

		r, theta := cmplx.Polar(a)
		fmt.Printf("\n a=%v; cmplx.Polar(a)  --->  r = %v, theta=%v\n", a, r, theta)
		fmt.Printf("\n a=%v -> cmplx.Rect(r, theta) = %v\n", a, cmplx.Rect(r, theta))

		fmt.Printf("check(cmath.Exp(%v), %v)\n", a, cmplx.Exp(a))
		fmt.Printf("check(cmath.Conj(%v), %v)\n", a, cmplx.Conj(a))
		fmt.Printf("check(cmath.Abs(%v), %v)\n", a, cmplx.Abs(a))
		fmt.Printf("check(cmath.Phase(%v), %v)\n", a, cmplx.Phase(a))
		fmt.Printf("check(cmath.Log(%v), %v)\n", a, cmplx.Log(a))
		fmt.Printf("check(cmath.Sqrt(%v), %v)\n", a, cmplx.Sqrt(a))
		fmt.Printf("check(cmath.Sin(%v), %v)\n", a, cmplx.Sin(a))
		fmt.Printf("check(cmath.Cos(%v), %v)\n", a, cmplx.Cos(a))
		fmt.Printf("check(cmath.Tan(%v), %v)\n", a, cmplx.Tan(a))
		fmt.Printf("check(cmath.Cot(%v), %v)\n", a, cmplx.Cot(a))
		fmt.Printf("check(cmath.Sinh(%v), %v)\n", a, cmplx.Sinh(a))
		fmt.Printf("check(cmath.Cosh(%v), %v)\n", a, cmplx.Cosh(a))
		fmt.Printf("check(cmath.Tanh(%v), %v)\n", a, cmplx.Tanh(a))
		fmt.Printf("check(cmath.Asin(%v), %v)\n", a, cmplx.Asin(a))
		fmt.Printf("check(cmath.Acos(%v), %v)\n", a, cmplx.Acos(a))
		fmt.Printf("check(cmath.Atan(%v), %v)\n", a, cmplx.Atan(a))
		fmt.Printf("check(cmath.Asinh(%v), %v)\n", a, cmplx.Asinh(a))
		fmt.Printf("check(cmath.Acosh(%v), %v)\n", a, cmplx.Acosh(a))
		fmt.Printf("check(cmath.Atanh(%v), %v)\n", a, cmplx.Atanh(a))

		//cmplx.Atan2(c2, c1)
		//cmplx.ComplexLog(b, z)

		fmt.Printf("done\n")
	})
}
