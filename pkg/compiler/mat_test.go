package compiler

import (
	//"fmt"
	//"strings"
	"testing"

	//"github.com/gijit/gi/pkg/verb"
	cv "github.com/glycerine/goconvey/convey"
)

func Test500MatrixDeclOfDoubleSlice(t *testing.T) {

	cv.Convey(`[][]float inside matrix struct`, t, func() {

		src := `
type Matrix struct {
	A    [][]float64
}
m := &Matrix{A:[][]float64{[]float64{1,2},[]float64{3,4}}}
e := m.A[0][1]
f := m.A[1][1]
// g is even fast repro than slc
g:=[][]int{[]int{1,2}}
slc := m.A[1]
`
		// e == 2
		// f == 4
		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation := inc.Tr([]byte(src))
		pp("go:'%s'  -->  '%s' in lua\n", src, string(translation))
		//fmt.Printf("go:'%#v'  -->  '%#v' in lua\n", src, translation)

		//cv.So(string(translation), matchesLuaSrc, ``)

		LoadAndRunTestHelper(t, vm, translation)

		LuaMustFloat64(vm, "e", 2)
		LuaMustFloat64(vm, "f", 4)

		// bad: [][]int{[0]= [][]int{[0]= 1LL, [1]= 2LL, }, }
		rawsrc := []byte(`gs = tostring(g)`)
		LoadAndRunTestHelper(t, vm, rawsrc)
		LuaMustString(vm, "gs", "[][]int{[0]= []int{[0]= 1LL, [1]= 2LL, }, }")

		// slc was getting an extra [], e.g.
		// bad: [][]float64{[0]= 3, [1]= 4, }
		// when we want
		// good: []float64{[0]= 3, [1]= 4, }
		rawsrc = []byte(`s = tostring(slc)`)
		LoadAndRunTestHelper(t, vm, rawsrc)
		LuaMustString(vm, "s", "[]float64{[0]= 3, [1]= 4, }")
	})
}

/*
func Test501MatrixMultiply(t *testing.T) {

	cv.Convey(`full matrix multiply program.`, t, func() {

		// see _bench/matmul.go.txt
		src := `
package main

// matrix multiplication benchmark

import (
	"fmt"
	"math/rand"
	"time"
)

type Matrix struct {
	A    [][]float64
	Nrow int
	Ncol int
}

//
// vector of vectors  matrix: not necessarily the
//  fastest, but we want to compare the same
//  approach in Go as was done in thed matrix.ss chez
//  implementation.
//
func NewMatrix(nrow, ncol int, fill bool) *Matrix {
	m := &Matrix{
		A:    make([][]float64, nrow),
		Nrow: nrow,
		Ncol: ncol,
	}
	for i := range m.A {
		m.A[i] = make([]float64, ncol)
	}
	if fill {
		for i := range m.A {
			for j := range m.A[i] {
				m.A[i][j] =
					float64(rand.Intn(100)) / float64(2.0+rand.Intn(100))
			}
		}
	}
	return m
}

// m1 x m2 matrix multiplication
func mult(m1, m2 *Matrix) (r *Matrix) {
	if m1.Ncol != m2.Nrow {
		panic(fmt.Sprintf(
			"incompatible: m1.Ncol=%v, m2.Nrow=%v", m1.Ncol, m2.Nrow))
	}
	r = NewMatrix(m1.Nrow, m2.Ncol, false)
	nr1 := m1.Nrow
	nr2 := m2.Nrow
	nc2 := m2.Ncol
	for i := 0; i < nr1; i++ {
		for k := 0; k < nr2; k++ {
			for j := 0; j < nc2; j++ {
				a := r.Get(i, j)
				a += m1.Get(i, k) * m2.Get(k, j)
				r.Set(i, j, a)
			}

		}
	}
	return
}

func (m *Matrix) Set(i, j int, val float64) {
	m.A[i][j] = val
}

func (m *Matrix) Get(i, j int) float64 {
	return m.A[i][j]
}

// MatScaMul multiplies a matrix by a scalar.
func MatScaMul(m *Matrix, x float64) (r *Matrix) {
	r = NewMatrix(m.Nrow, m.Ncol, false)
	for i := 0; i < m.Nrow; i++ {
		for j := 0; j < m.Ncol; j++ {
			r.Set(i, j, x*m.Get(i, j))
		}
	}
	return
}

var done bool
func main() {
	sz := 500
	for i := 0; i < 10; i++ {
		a := NewMatrix(sz, sz, true)
		b := NewMatrix(sz, sz, true)
		t0 := time.Now()
		mult(a, b)
		elap := time.Since(t0)
		fmt.Printf("%v x %v matrix multiply in Go took %v msec\n",
			sz, sz, int(elap/time.Millisecond))
	}
    done = true
}
`
		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		//verb.VerboseVerbose = true

		importPath := ""
		translation, err := inc.FullPackage([]byte(src), importPath)
		panicOn(err)
		pp("go:'%s'  -->  '%s' in lua\n", src, string(translation))
		fmt.Printf("go origin source:")
		fmt.Printf(strings.Replace(src, "%", "%%", -1))
		fmt.Printf("\n\n  --> translation to lua -->\n\n")
		fmt.Printf(string(translation))
		fmt.Printf("\n\n")

		//cv.So(string(translation), matchesLuaSrc, ``)

		LoadAndRunTestHelper(t, vm, translation)

		LoadAndRunTestHelper(t, vm, []byte("main()"))

		LuaMustBool(vm, "done", true)
		cv.So(true, cv.ShouldBeTrue)
	})
}
*/
