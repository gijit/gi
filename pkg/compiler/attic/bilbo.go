package main

import (
	"fmt"
)

type Baggins interface {
	WearRing() bool
}
type Gollum interface {
	Scowl() int
}
type hobbit struct {
	hasRing bool
}

func (h *hobbit) WearRing() bool {
	h.hasRing = !h.hasRing
	return h.hasRing
}

type Wolf struct {
	Claw    int
	HasRing bool
}

func (w *Wolf) Scowl() int {
	w.Claw++
	return w.Claw
}
func battle(g Gollum, b Baggins) (int, bool) {
	return g.Scowl(), b.WearRing()
}
func tryTheTypeSwitch(i interface{}) int {
	switch x := i.(type) {
	case Gollum:
		return x.Scowl()
	case Baggins:
		if x.WearRing() {
			return 1
		}
	}
	return 0
}
func main() {
	w := &Wolf{}
	bilbo := &hobbit{}
	i0, b0 := battle(w, bilbo)
	i1, b1 := battle(w, bilbo)
	fmt.Printf("i0=%v, b0=%v\n", i0, b0)
	fmt.Printf("i1=%v, b1=%v\n", i1, b1)
	fmt.Printf("tried wolf=%v\n", tryTheTypeSwitch(w))
	fmt.Printf("tried bilbo=%v\n", tryTheTypeSwitch(bilbo))
}

/*
i0=1, b0=true
i1=2, b1=false
tried wolf=3
tried bilbo=1
*/
