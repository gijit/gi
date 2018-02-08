package assertions

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	//	success               = ""
	//	needExactValues       = "This assertion requires exactly %d comparison values (you provided %d)."
	shouldMatchModulo     = "-------------------------\nExpected expected string:\n-------------------------\n'%s'\n=========================\n  and actual string:\n=========================\n'%s'\n---------------------\n to match (ignoring %s)\n (but they did not!; first diff at '%s', pos %d); and \nFull diff -b:\n%s\n"
	shouldStartWithModulo = "Expected expected PREFIX string '%s'\n       to be found at the start of actual string '%s'\n to  (ignoring %s)\n (but they did not!; first diff at '%s', pos %d); and \nFull diff -b:\n%s\n"
	shouldContainModuloWS = "Expected expected string '%s'\n       to contain string '%s'\n (ignoring whitespace)\n (but it did not!)"
	//shouldBothBeStrings   = "Both arguments to this assertion must be strings (you provided %v and %v)."
)

// ShouldMatchModulo receives exactly two string parameters and an ignore map. It ensures that the order
// of not-ignored characters in the two strings is identical. Runes specified in the ignore map
// are ignored for the purposes of this string comparison, and each should map to true.
// ShouldMatchModulo thus allows you to do whitespace insensitive comparison, which is very useful
// in lexer/parser work.
//
func ShouldMatchModulo(ignoring map[rune]bool, actual interface{}, expected ...interface{}) string {
	if fail := need(1, expected); fail != success {
		return fail
	}

	value, valueIsString := actual.(string)
	expec, expecIsString := expected[0].(string)

	if !valueIsString || !expecIsString {
		return fmt.Sprintf(shouldBothBeStrings, reflect.TypeOf(actual), reflect.TypeOf(expected[0]))
	}

	equal, vpos, _ := stringsEqualIgnoring(value, expec, ignoring)
	if equal {
		return success
	} else {
		// extract the string fragment at the differnce point to make it easier to diagnose
		diffpoint := ""
		const diffMax = 20
		vrune := []rune(value)
		n := len(vrune) - vpos
		if n > diffMax {
			n = diffMax
		}
		if vpos == 0 {
			vpos = 1
		}
		diffpoint = string(vrune[vpos-1 : (vpos - 1 + n)])

		diff := Diffb(value, expec)

		ignored := "{"
		switch len(ignoring) {
		case 0:
			return fmt.Sprintf(shouldMatchModulo, expec, value, "nothing", diffpoint, vpos-1, diff)
		case 1:
			for k := range ignoring {
				ignored = ignored + fmt.Sprintf("'%c'", k)
			}
			ignored = ignored + "}"
			return fmt.Sprintf(shouldMatchModulo, expec, value, ignored, diffpoint, vpos-1, diff)

		default:
			for k := range ignoring {
				ignored = ignored + fmt.Sprintf("'%c', ", k)
			}
			ignored = ignored + "}"
			return fmt.Sprintf(shouldMatchModulo, expec, value, ignored, diffpoint, vpos-1, diff)
		}
	}
}

// ShouldMatchModuloSpaces compares two strings but ignores ' ' spaces.
// Serves as an example of use of ShouldMatchModulo.
//
func ShouldMatchModuloSpaces(actual interface{}, expected ...interface{}) string {
	if fail := need(1, expected); fail != success {
		return fail
	}
	return ShouldMatchModulo(map[rune]bool{' ': true}, actual, expected[0])
}

func ShouldMatchModuloWhiteSpace(actual interface{}, expected ...interface{}) string {
	if fail := need(1, expected); fail != success {
		return fail
	}
	return ShouldMatchModulo(map[rune]bool{' ': true, '\n': true, '\t': true}, actual, expected[0])
}

func ShouldStartWithModuloWhiteSpace(actual interface{}, expectedPrefix ...interface{}) string {
	if fail := need(1, expectedPrefix); fail != success {
		return fail
	}

	ignoring := map[rune]bool{' ': true, '\n': true, '\t': true}

	value, valueIsString := actual.(string)
	expecPrefix, expecIsString := expectedPrefix[0].(string)

	if !valueIsString || !expecIsString {
		return fmt.Sprintf(shouldBothBeStrings, reflect.TypeOf(actual), reflect.TypeOf(expectedPrefix[0]))
	}

	equal, vpos, _ := hasPrefixEqualIgnoring(value, expecPrefix, ignoring)
	if equal {
		return success
	} else {
		diffpoint := ""
		const diffMax = 20
		vrune := []rune(value)
		n := len(vrune) - vpos + 1
		if n > diffMax {
			n = diffMax
		}
		beg := vpos - 1
		if beg < 0 {
			beg = 0
		}
		diffpoint = string(vrune[beg:(vpos - 1 + n)])

		diff := Diffb(value, expecPrefix)

		ignored := "{"
		switch len(ignoring) {
		case 0:
			return fmt.Sprintf(shouldStartWithModulo, expecPrefix, value, "nothing", diffpoint, vpos-1, diff)
		case 1:
			for k := range ignoring {
				ignored = ignored + fmt.Sprintf("'%c'", k)
			}
			ignored = ignored + "}"
			return fmt.Sprintf(shouldStartWithModulo, expecPrefix, value, ignored, diffpoint, vpos-1, diff)

		default:
			for k := range ignoring {
				ignored = ignored + fmt.Sprintf("'%c', ", k)
			}
			ignored = ignored + "}"
			return fmt.Sprintf(shouldStartWithModulo, expecPrefix, value, ignored, diffpoint, vpos-1, diff)
		}
	}
}

// returns if equal, and if not then rpos and spos hold the position of first mismatch
func stringsEqualIgnoring(a, b string, ignoring map[rune]bool) (equal bool, rpos int, spos int) {
	r := []rune(a)
	s := []rune(b)

	nextr := 0
	nexts := 0

	for {
		// skip past spaces in both r and s
		for nextr < len(r) {
			if ignoring[r[nextr]] {
				nextr++
			} else {
				break
			}
		}

		for nexts < len(s) {
			if ignoring[s[nexts]] {
				nexts++
			} else {
				break
			}
		}

		if nextr >= len(r) && nexts >= len(s) {
			return true, -1, -1 // full match
		}

		if nextr >= len(r) {
			return false, nextr, nexts
		}
		if nexts >= len(s) {
			return false, nextr, nexts
		}

		if r[nextr] != s[nexts] {
			return false, nextr, nexts
		}
		nextr++
		nexts++
	}

	return false, nextr, nexts
}

// returns if equal, and if not then rpos and spos hold the position of first mismatch
func hasPrefixEqualIgnoring(str, prefix string, ignoring map[rune]bool) (equal bool, spos int, rpos int) {
	s := []rune(str)
	r := []rune(prefix)

	nextr := 0
	nexts := 0

	for {
		// skip past spaces in both r and s
		for nextr < len(r) {
			if ignoring[r[nextr]] {
				nextr++
			} else {
				break
			}
		}

		for nexts < len(s) {
			if ignoring[s[nexts]] {
				nexts++
			} else {
				break
			}
		}

		if nextr >= len(r) && nexts >= len(s) {
			return true, -1, -1 // full match
		}

		if nextr >= len(r) {
			return true, nexts, nextr // for prefix testing
		}
		if nexts >= len(s) {
			return false, nexts, nextr
		}

		if r[nextr] != s[nexts] {
			return false, nexts, nextr
		}
		nextr++
		nexts++
	}

	return false, nexts, nextr
}

/*func need(needed int, expected []interface{}) string {
	if len(expected) != needed {
		return fmt.Sprintf(needExactValues, needed, len(expected))
	}
	return success
}
*/

func ShouldContainModuloWhiteSpace(haystack interface{}, expectedNeedle ...interface{}) string {
	if fail := need(1, expectedNeedle); fail != success {
		return fail
	}

	value, valueIsString := haystack.(string)
	expecNeedle, expecIsString := expectedNeedle[0].(string)

	if !valueIsString || !expecIsString {
		return fmt.Sprintf(shouldBothBeStrings, reflect.TypeOf(haystack), reflect.TypeOf(expectedNeedle[0]))
	}

	elimWs := func(r rune) rune {
		if r == ' ' || r == '\t' || r == '\n' {
			return -1 // drop the rune
		}
		return r
	}

	h := strings.Map(elimWs, value)
	n := strings.Map(elimWs, expecNeedle)

	if strings.Contains(h, n) {
		return success
	}

	return fmt.Sprintf(shouldContainModuloWS, value, expecNeedle)
}

// older method

// ShouldBeEqualIgnoringSpaces receives exactly two string parameters and ensures that the order
// of non-space characters in the two strings is identical. Only one character ' '
// is considered a space. This is not a general white-space ignoring routine! Differences in tabs or
// newlines *will* be noticed, and will cause the two strings to look different.
func ShouldBeEqualIgnoringSpaces(actual interface{}, expected ...interface{}) string {
	if fail := need(1, expected); fail != success {
		return fail
	}

	value, valueIsString := actual.(string)
	expec, expecIsString := expected[0].(string)

	// handle []byte as either or both inputs.
	if !valueIsString {
		valueBy, valueIsSliceOfByte := actual.([]byte)
		if valueIsSliceOfByte {
			value = string(valueBy)
			valueIsString = true
		}
	}

	if !expecIsString {
		expecBy, expecIsSliceOfByte := expected[0].([]byte)

		if expecIsSliceOfByte {
			expec = string(expecBy)
			expecIsString = true
		}
	}

	if !valueIsString || !expecIsString {
		return fmt.Sprintf(shouldBothBeStrings, reflect.TypeOf(actual), reflect.TypeOf(expected[0]))
	}

	if equalIgnoringSpaces(value, expec) {
		return success
	} else {
		return fmt.Sprintf(shouldHaveBeenEqualIgnoringSpaces, value, expec)
	}
}

func equalIgnoringSpaces(r, s string) bool {
	nextr := 0
	nexts := 0

	for {
		// skip past spaces in both r and s
		for nextr < len(r) {
			if r[nextr] == ' ' {
				nextr++
			} else {
				break
			}
		}

		for nexts < len(s) {
			if s[nexts] == ' ' {
				nexts++
			} else {
				break
			}
		}

		if nextr >= len(r) && nexts >= len(s) {
			return true
		}

		if nextr >= len(r) {
			return false
		}
		if nexts >= len(s) {
			return false
		}

		if r[nextr] != s[nexts] {
			return false
		}
		nextr++
		nexts++
	}

	return false
}
