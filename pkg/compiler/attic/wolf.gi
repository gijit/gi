package main

import (
	"fmt"
)

type Bowser interface {
	Hi()
}

type Possum interface {
	Hi()
    Pebbles()
}

type Unsat interface {
	Hi()
    Pebbles()
    MissMe()
}

type B struct{}

func (b *B) Hi() {
	fmt.Printf("B.Hi called\n")
}
func (b *B) Pebbles() {}

    chk := 0
	var v Bowser = &B{}
	switch v.(type) {
    case Possum:
		fmt.Printf("ooh! it types as a Possum!\n")
        chk = 2
	case Bowser:
		fmt.Printf("yabadadoo! it types as a Bowser!\n")
        chk = 1
	}
    fmt.Printf("chk = '%v'\n", chk)

    // and verify that v implements Bowser too:
    asBowser, isBowser := v.(Bowser)
    asIsNil := (asBowser == nil)

    // negative check, should not convert:
    asUn, isUn := v.(Unsat)
    asUnNil := (asUn == nil)
