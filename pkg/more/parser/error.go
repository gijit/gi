package parser

import (
	"fmt"
)

var ErrSyntax = fmt.Errorf("syntax error")
var ErrMoreInput = fmt.Errorf("more input required")
