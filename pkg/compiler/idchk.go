package compiler

import (
	"github.com/gijit/gi/pkg/ast"
)

/*
Restriction to a subset of legal `go` programs:

In `gijit`, you can't have a variable named `int`, or `float64`. These
are names of two of the pre-declared numeric types in Go.

`gijit` won't let you declare a variable name that
reuses any of the basic, pre-declared type names.

Although in Go this is technically allowed, it is highly confusing,
and poor practice.

So while:
~~~
func main() {
	var int int
        _ = int
}
~~~
is a legal `Go` program, it won't run on `gijit`.

The reason for this restriction is that otherwise
the Go type checker can get corrupted by simple syntax errors. That's not an
issue for a full-recompile from the start each time, but for
a continuously online REPL, it messes with the type checking
that follows the syntax error.
*/

type identVisitor struct {
	av *assignVisitor
}

type assignVisitor struct {
	iv    *identVisitor
	bad   bool
	which string
}

func newAssignVisitor() *assignVisitor {
	av := &assignVisitor{
		iv: &identVisitor{},
	}
	av.iv.av = av
	return av
}

func (av *assignVisitor) Visit(node ast.Node) (w ast.Visitor) {
	switch n := node.(type) {

	case *ast.ValueSpec:
		pp("Visit assignVisitor, ValueSpec, n = '%#v'", n)
		for _, id := range n.Names {
			if predeclared[id.Name] {
				pp("Visit found bad '%s'", id.Name)
				av.bad = true
				av.which = id.Name
				return nil
			}
		}

	case *ast.AssignStmt:
		pp("Visit assignVisitor, AssignStmt, n = '%#v'", n)
		walkExprList(av.iv, n.Lhs)
		if av.bad {
			return nil
		}
	default:
		pp("Visit assignVisitor ignoring node '%#v'/'%T'", node, node)
	}
	return av
}

func walkExprList(v ast.Visitor, list []ast.Expr) {
	for _, x := range list {
		ast.Walk(v, x)
	}
}

func (iv *identVisitor) Visit(node ast.Node) (w ast.Visitor) {
	switch id := node.(type) {
	case *ast.Ident:
		pp("Visit identCheckVisitor, nm = '%v'", id.Name)
		if predeclared[id.Name] {
			iv.av.bad = true
			iv.av.which = id.Name
			return nil
		}
	}
	return iv
}

func checkAllowedIdents(file *ast.File) (hasBadId bool, whichBad string) {
	pp("top of checkAllowedIdents")
	v := newAssignVisitor()
	for _, n := range file.Nodes {
		ast.Walk(v, n)
	}
	return v.bad, v.which
}
