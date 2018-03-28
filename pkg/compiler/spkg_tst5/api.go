package spkg_tst5

import "time"

func Astm(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	panicOn(err)
	return t
}

/*
type B
b := &ar{imin: 1, input:  &stream{sym: "AAPL"}}

sb := &streambase{dir: "teststreams"}

sm := &srcmap{base:sb, date:"2017/07/05", root:b, symsrc:hash{"AAPL":"testdata/AAPL.2017.07.05.1minute.bars.gz"}}

togo(sm) // one should suffice for all recursively.

arr := b.Val(tm)

e := arr[0]

assert(e.Cl ==  143.91)
*/
