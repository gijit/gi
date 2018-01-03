package main

import (
	"strings"
)

func panicOn(err error) {
	if err != nil {
		panic(err)
	}
}

func panicOnNotNA(err error) {
	if err != nil {
		if strings.HasPrefix(err.Error(), "ErrNA: ") {
			return
		}
		panic(err)
	}
}
