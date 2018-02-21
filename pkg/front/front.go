// accept multiline input, only feeding
// it to the interpreter when we have
// a full statement/declaration/expression.
package front

import (
	"bytes"
	"fmt"
)

// eofSeen is returned true if our the input is incomplete.
func TopLevelParseGoSource(sourceCode []byte) (eofSeen, errorSeen, empty bool, err error) {

	// we always start with the full source to parse, it's
	// never a partial continuation with prior unseen stuff.
	// While TrimLeft is safe, TrimRight and TrimSpace may
	// take away import space in strings/comments that
	// just not be finished yet.
	sourceCode = bytes.TrimLeft(sourceCode, " \t\n\r")
	if len(sourceCode) == 0 {
		return false, false, true, nil
	}

	var errh ErrorHandler = nil // stop at first error.
	var fileh FilenameHandler = nil
	var base *PosBase
	stream := bytes.NewBuffer(sourceCode)

	defer func() {
		// jea debug:
		//return

		// catch panics that are communicating
		// parse results quickly to us. Re throw
		// any others.
		recov := recover()
		switch recov {
		case nil:
		case ErrSyntax:
			errorSeen = true
		case ErrMoreInput:
			eofSeen = true
		case CompleteNoError:
			// effectively:	errorSeen = false; and eofSeen = false;
		case EmptyInput:
			empty = true
		default:
			errorSeen = true
			err = fmt.Errorf("%v", recov)
		}
	}()

	_, err = Parse(base, stream, errh, nil, fileh)
	pp("normal return from Parse(), err = '%v'", err)

	if err == ErrSyntax {
		errorSeen = true
		err = nil
	}
	if err == ErrMoreInput {
		eofSeen = true
		err = nil
	}
	panicOn(err)
	pp("we got past the Parse !")

	// set eofSeen, errorSeen !!
	return
}
