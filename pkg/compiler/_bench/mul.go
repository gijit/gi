package main

// matrix multiplication benchmark

import (
	"fmt"
	"time"
	//"math/rand"
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
	next := 2.0
	if fill {
		for i := range m.A {
			for j := range m.A[i] {
				m.A[i][j] = next
				next++
				//float64(rand.Intn(100)) / float64(2.0+rand.Intn(100))
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

func runMultiply(sz, i, j int) float64 {
	var mu *Matrix
	for k := 0; k < 1; k++ {
		a := NewMatrix(sz, sz, true)
		b := NewMatrix(sz, sz, true)
		//t0 := time.Now()
		mu = mult(a, b)
		//elap := time.Since(t0)
		//fmt.Printf("%v x %v matrix multiply in Go took %v msec\n",
		//	sz, sz, int(elap/time.Millisecond))
		//fmt.Printf("%v x %v matrix multiply mu.A[2][2] = %v\n", sz, sz, mu.A[2][2])
	}
	done = true
	return mu.A[i][j]
}

func main() {
	r := runMultiply(10, 9, 9)
	fmt.Printf("r='%v'\n", r)
	t0 := time.Now()
	fmt.Printf("runMultiply(100,9,9) -> %v\n", int(runMultiply(100, 9, 9)))
	elap := time.Since(t0)
	fmt.Printf("compiled Go elap = %v\n", elap)
}

// 3 x 3 matrix multiply mu.A[2,2] = 195
// runMultiply(10,9,9) -> 54865
// runMultiply(100,9,9) -> 480371650
// compiled Go elap = 3.085894ms
