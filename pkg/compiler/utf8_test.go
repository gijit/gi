package compiler

import (
	"fmt"
	"testing"

	cv "github.com/glycerine/goconvey/convey"
)

func Test041RangeOverUtf8BytesInString(t *testing.T) {

	cv.Convey(`From the https://blog.golang.org/strings blog example`+
		` of a for-range loop over utf8 in strings: `+
		`Given the string nihongo := "日本語" containing `+
		`three multibyte utf8 characters, `+
		`then the loop:`+
		``+
		`      for i, runeValue := range s { println(r); }`+
		``+
		`should loop over whole utf8 runes and not `+
		`over individual bytes. So runeValue should `+
		`take on each of the three full runes in turn.`, t, func() {

		// Resources:
		//  run this test with: `go test -v -run 041`
		//
		// Rob Pike's blog: https://github.com/gijit/gi/blob/master/pkg/utf8/utf8.go
		//
		// see https://github.com/gijit/gi/issues/13 for discussion, and to note
		//  progress and completion.
		//
		// See the lua code in https://github.com/gijit/gi/blob/master/pkg/utf8/utf8.lua
		// See the Go code in https://github.com/gijit/gi/blob/master/pkg/utf8/utf8.go

		// From Rob Pike's blog https://blog.golang.org/strings
		//
		// Besides the axiomatic detail that Go
		// source code is UTF-8, there's really
		// only one way that Go treats UTF-8 specially,
		// and that is when using a for range loop on a string.
		//
		// We've seen what happens with a regular for loop.
		// A for range loop, by contrast, decodes one
		// UTF-8-encoded rune on each iteration. Each time
		// around the loop, the index of the loop is the
		// starting position of the current rune, measured
		// in bytes, and the code point is its value.
		//
		// ...
		//
		// const nihongo = "日本語"
		// for index, runeValue := range nihongo {
		//    fmt.Printf("%#U starts at byte position %d\n", runeValue, index)
		// }

		// should print:
		// U+65E5 '日' starts at byte position 0
		// U+672C '本' starts at byte position 3
		// U+8A9E '語' starts at byte position 6

		// since we don't have fmt/imports online yet, we'll settle
		// for just taking the string apart at the utf8 separation
		// points. No printing f

		code := `
    runes := []rune{0,0,0}
    const nihongo = "日本語"  // translated, means "Japanese"
    for i, runeValue := range nihongo {
        runes[i] = runeValue
    }
    r0 := runes[0]
    r1 := runes[1]
    r2 := runes[2]
`

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm)

		translation := inc.Tr([]byte(code))
		fmt.Printf("\n\n translation = '%s'\n\n", string(translation))

		LuaRunAndReport(vm, string(translation))
		fmt.Printf("\n past LuaRunAndReport \n")
		LuaMustString(vm, "r0", "日")
		LuaMustString(vm, "r1", "本")
		LuaMustString(vm, "r2", "語")
	})
}
