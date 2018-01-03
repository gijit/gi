package compiler

import (
	"fmt"
	"testing"

	cv "github.com/glycerine/goconvey/convey"
	"github.com/robertkrimen/otto"
)

func Test001JavascriptTranslation(t *testing.T) {
	inc := NewIncrState()

	cv.Convey("assignment", t, func() {
		cv.So(string(inc.Tr([]byte("a := 10;"))), cv.ShouldMatchModuloWhiteSpace, "a = 10;")
		pp("GOOD: past 1st")

		cv.So(string(inc.Tr([]byte("func adder(a, b int) int { return a + b};  sum1 := adder(5,5)"))), cv.ShouldMatchModuloWhiteSpace,
			`adder = function(a, b) {
				         var a, b;
				         return a + b >> 0;
			         }; sum1 = adder(5,5);`)

		pp("GOOD: past 2nd")

		cv.So(string(inc.Tr([]byte("sum2 := adder(a,a)"))), cv.ShouldMatchModuloWhiteSpace,
			`sum2 = adder(a, a);`)
		pp("GOOD: past 3rd")
	})
}

func Test002OttoEvalIncremental(t *testing.T) {

	// and then eval!
	vm := otto.New()
	inc := NewIncrState()

	srcs := []string{"a := 10;", "func adder(a, b int) int { return a + b}; ", "sum := adder(a,a);"}
	for _, src := range srcs {
		translation := inc.Tr([]byte(src))
		fmt.Printf("go:'%s'  -->  '%s' in js\n", src, translation)
		//fmt.Printf("go:'%#v'  -->  '%#v' in js\n", src, translation)

		v, err := vm.Eval(string(translation))
		if err != nil {
			panic(err)
		}
		fmt.Printf("v back = '%#v'\n", v)
	}
	value, err := vm.Get("sum")
	if err != nil {
		panic(err)
	}

	value_int, err := value.ToInteger()
	if err != nil {
		panic(err)
	}

	fmt.Printf("value_int=%v", value_int)
	if value_int != 20 {
		panic(fmt.Sprintf("expected 20, got %v", value_int))
	}
}

// func Test003ImportsAtRepl(t *testing.T) {
// 	inc := NewIncrState()

// 	cv.Convey("imports", t, func() {
// 		cv.So(string(inc.Tr([]byte(`import "fmt"; fmt.Printf("hello world!")`))), cv.ShouldMatchModuloWhiteSpace, "")
// 		pp("GOOD: past 1st import")
// 	})
// }

func Test004ExpressionsAtRepl(t *testing.T) {
	inc := NewIncrState()

	cv.Convey("expressions alone at top level", t, func() {
		cv.So(string(inc.Tr([]byte(`a:=10; a`))), cv.ShouldMatchModuloWhiteSpace, "a=10; a;")
	})
}
