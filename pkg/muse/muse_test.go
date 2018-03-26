package muse

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/glycerine/gi/pkg/ast"
	//"github.com/glycerine/gi/pkg/gostd/build"
	"github.com/glycerine/gi/pkg/parser"
	"github.com/glycerine/gi/pkg/token"
	"github.com/glycerine/gi/pkg/types"
	"github.com/glycerine/gi/pkg/verb"

	cv "github.com/glycerine/goconvey/convey"
)

var pp = verb.PP

func init() {
	verb.Verbose = true
	verb.VerboseVerbose = true
}

func Test001BasicTypeConversion(t *testing.T) {

	cv.Convey(`muse.Pun() should convert basic type.Types to basic reflect.Types`, t, func() {
		m := NewMuse()

		b := types.Typ[types.Bool]
		rt, err := m.Pun(b)
		cv.So(err, cv.ShouldBeNil)
		cv.So(rt.String(), cv.ShouldResemble, reflect.Bool.String())

		s := types.Typ[types.String]
		rt, err = m.Pun(s)
		cv.So(err, cv.ShouldBeNil)
		cv.So(rt.String(), cv.ShouldResemble, reflect.String.String())

		v := types.Typ[types.Int]
		rt, err = m.Pun(v)
		cv.So(err, cv.ShouldBeNil)
		cv.So(rt.String(), cv.ShouldResemble, reflect.Int.String())

		v = types.Typ[types.Int64]
		rt, err = m.Pun(v)
		cv.So(err, cv.ShouldBeNil)
		cv.So(rt.String(), cv.ShouldResemble, reflect.Int64.String())

		v = types.Typ[types.Uint64]
		rt, err = m.Pun(v)
		cv.So(err, cv.ShouldBeNil)
		cv.So(rt.String(), cv.ShouldResemble, reflect.Uint64.String())

	})
}

func Test002StructTypeConversion(t *testing.T) {

	cv.Convey(`muse.Pun() should convert a struct type.Types `+
		`to a struct reflect.Types`, t, func() {

		src := `
type Tree struct {
   Name       string
}
`
		m := NewMuse()

		fileSet := token.NewFileSet()
		file, err := parser.ParseFile(fileSet, "", src, 0)
		panicOn(err)

		pp("len file.Nodes = '%#v'", len(file.Nodes)) // 1
		pp("file.Nodes[0] = '%#v'", file.Nodes[0])    // '&ast.GenDecl{Doc:(*ast.CommentGroup)(nil), TokPos:2, Tok:84, Lparen:0, Specs:[]ast.Spec{(*ast.TypeSpec)(0xc42000e870)}, Rparen:0}'

		switch x := file.Nodes[0].(type) {
		case *ast.GenDecl:

			pp("len x.Specs = '%#v'", len(x.Specs)) // 1
			pp("x.Specs[0] = '%#v'", x.Specs[0])    // '&ast.TypeSpec{Doc:(*ast.CommentGroup)(nil), Name:(*ast.Ident)(0xc42000a280), Assign:0, Type:(*ast.StructType)(0xc42000a2e0), Comment:(*ast.CommentGroup)(nil)}'
			switch y := x.Specs[0].(type) {
			case *ast.TypeSpec:

				nm := y.Name  // nm = '&ast.Ident{NamePos:7, Name:"Tree", Obj:(*ast.Object)(0xc42007d360)}'
				typ := y.Type // '&ast.StructType{Struct:12, Fields:(*ast.FieldList)(0xc42000e8a0), Incomplete:false}'

				pp("nm = '%#v', typ= '%#v'", nm, typ)
				// nm = '&ast.Ident{NamePos:7, Name:"Tree", Obj:(*ast.Object)(0xc42007d360)}', typ= '&ast.StructType{Struct:12, Fields:(*ast.FieldList)(0xc42000e8a0), Incomplete:false}'

				pp("file = '%#v'", file)

				checked := typeCheck(typ, fileSet, file)
				pp("checked = '%#v'", checked)
				if checked == nil {
					panic(fmt.Sprintf("got nil checked back for typ='%#v'", typ))
				}

				rt, err := m.Pun(checked)
				cv.So(err, cv.ShouldBeNil)
				cv.So(rt.String(), cv.ShouldResemble, ``)

			}
		}

	})
}

func Test003InterfaceTypeConversion(t *testing.T) {

	cv.Convey(`muse.Pun() should convert an interface with its method set from type.Types `+
		`to an interface in reflect.Types`, t, func() {

		m := NewMuse()

		b := types.Typ[types.Bool]
		rt, err := m.Pun(b)
		cv.So(err, cv.ShouldBeNil)
		cv.So(rt.String(), cv.ShouldResemble, reflect.Bool.String())

	})
}
