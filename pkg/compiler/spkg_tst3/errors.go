package spkg_tst3

type ErrorW2 interface{}

type S struct{}

// fails- typechecks ErrorW2 to a nil type, because these
// two share the same name.
func (*S) ErrorW2() {}
