package spkg_tst3

type ErrorW2 interface{}

type S struct{}

// fixed now, but was failing to typechecks:
// ErrorW2 to a nil type, because these
// two share the same name. The deletion
// and replacement code for re-defining
// types at the REPL needed refinement.
func (*S) ErrorW2() {}
