package spkg_tst4

type R struct {
	A [2][]byte
}

func NewR() *R {
	r := &R{}
	r.A[0] = []byte("hello")
	r.A[1] = []byte("world")
	return r
}

func (r *R) Get1() string {
	return string(r.A[1])
}
