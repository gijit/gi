package front

import (
	"fmt"
)

// Each set of lines processed by front can be classified
// into one member of the following partition. The set
// is either a syntax error, requires more input, or is
// complete all by itself. When complete by itself, it
// may be a complete expression, a complete declaration,
// or a statement. Complete sets won't end with dangling
// commas ',' or open parenthesis '(' or braces '{' or
// open brackets '['.
//
// We also distinguish the empty string "" input, so
// as to not confuse repl users by showing a continuation
// prompt after no input.

var ErrSyntax = fmt.Errorf("syntax error")
var ErrMoreInput = fmt.Errorf("more input required")
var CompleteNoError = fmt.Errorf("complete by itself") // not syntax error, no more input needed.
var EmptyInput = fmt.Errorf("empty input, no tokens.")
