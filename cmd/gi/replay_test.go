package main

import (
	"flag"
	"fmt"
	//"github.com/shurcooL/go-goon"
	"strings"
	"testing"

	"github.com/gijit/gi/pkg/types"
	"github.com/gijit/gi/pkg/verb"
	cv "github.com/glycerine/goconvey/convey"
)

var _ = fmt.Printf

func Test001ReplayOfStructdDef(t *testing.T) {

	cv.Convey(`if we replay struct defn and method call, the method call should succeed the 2nd time (was failing to replay from history, stumbling on the type checker)`, t, func() {
		src := `
 type S struct{}
 var s S
 func (s *S) Hi() {println("S.Hi() called")}
 s.Hi()
`
		fmt.Printf("replay 2x, src='%s'\n", src)

		myflags := flag.NewFlagSet("gi", flag.ExitOnError)
		cfg := &GIConfig{}
		cfg.DefineFlags(myflags)

		err := myflags.Parse([]string{"-q", "-no-liner"}) // , "-vv"})
		err = cfg.ValidateConfig()
		panicOn(err)
		r := NewRepl(cfg)

		verb.VerboseVerbose = true
		verb.Verbose = true

		/*
			// oddly, when sent as a 4 line chunk, not
			// broken up into separate lines, we don't see
			// the issue.
			err = r.Eval(src)
			cv.So(err, cv.ShouldBeNil)
			err = r.Eval(src)
			cv.So(err, cv.ShouldBeNil)
		*/

		// now split into lines: then we get the oops.
		lines := strings.Split(src, "\n")
		// and do the lines set 2x
		var om_S types.Object
		for j := 0; j < 2; j++ {

			if j == 1 {
				// display the ObjMap:
				for k, v := range r.inc.CurPkg.Arch.Check.ObjMap {
					pp("Arch.Check.ObjMap: key k='%#v' v ='%#c'", k, v)
					pp("k.Name()='%#v'", k.Name())
					switch k.Name() {
					case "S":
						om_S = k
						_ = om_S
					}
				}
				// works, so try to put this into types/check.go
				pp("on j=%v pass, deleting 'om_S'", j)
				delete(r.inc.CurPkg.Arch.Check.ObjMap, om_S) // green!!!
			}
			for i := range lines {
				err = r.Eval(lines[i])
				panicOn(err)

				// CurPkg.Arch is nil until i > 0 the first time.
				if i > 0 || j > 0 {
					fmt.Printf("TypesInfo.Types='%#v'\n", r.inc.CurPkg.Arch.TypesInfo) // .Types, .Defs, .Name2node

					//for nm := range r.inc.CurPkg.Arch.TypesInfo.Name2node {
					//	fmt.Printf("j=%v, i=%v, Name2Node[nm='%s']\n", j, i, nm) // S.Hi
					//}
					//fmt.Printf("\n node = '%#v'\n", node)
					k := 0
					for def, obj := range r.inc.CurPkg.Arch.TypesInfo.Defs {
						if def != nil {
							defnm := def.Name
							_ = obj
							fmt.Printf("\n on j=%v, after Eval of line i=%v:  TypesInfo.Defs[k=%v]='%v' = '%p'\n", j, i, k, defnm, obj) // .Types, .Defs, .Name2node
							// see output below [1]
						}
						k++
					}
				}

			}
			fmt.Printf("\n pass j=%v complete.\n", j)
		}
	})
}

/*
focus on the first main.S decl: 0xc4200818b0

decl.go:245 2018-02-08 15:28:42.51 +0700 ICT jea debug, decl.go: Named.setUnderlying for typ='&types.Named{obj:(*types.TypeName)(0xc4200818b0), underlying:types.Type(nil), methods:[]*types.Func(nil)}'/'main.S'

utils.go:570 2018-02-08 15:28:42.608 +0700 ICT typeKind called on ty='&types.Named{obj:(*types.TypeName)(0xc4200818b0), underlying:(*types.Struct)(0xc420011e60), methods:[]*types.Func(nil)}', returning res='__gi_kind_Struct'


translate.go:171 2018-02-08 15:28:42.634 +0700 ICT writing tr.CurPkg.Arch.NewCode[i=0].Code = '	__type__S = __gi_NewType(0, __gi_kind_Struct, "main", "S", "main.S", true, "main", true, nil);
	__type__S.__init("", {});

	 __type__S.__constructor = function(self)
		 return self; end;

'

elapsed: '116.164µs'

replay_test.go:57 2018-02-08 15:28:42.661 +0700 ICT TypesInfo.Types='&types.Info{Types:map[ast.Expr]types.TypeAndValue{(*ast.StructType)(0xc42000cd40):types.TypeAndValue{mode:0x3, Type:(*types.Struct)(0xc420011e60), Value:constant.Value(nil)}}, Defs:map[*ast.Ident]types.Object{(*ast.Ident)(0xc42000cde0):types.Object(nil), (*ast.Ident)(0xc42000cd20):(*types.TypeName)(0xc4200818b0)}, Uses:map[*ast.Ident]types.Object{}, Implicits:map[ast.Node]types.Object{}, Selections:map[*ast.SelectorExpr]*types.Selection{}, Scopes:map[ast.Node]*types.Scope{(*ast.File)(0xc4200b8500):(*types.Scope)(nil)}, Name2node:map[string]*types.FtypeAndScope(nil), InitOrder:[]*types.Initializer(nil), NewCode:[]*types.NewStuff{(*types.NewStuff)(0xc420011d40)}}'

 on j=0, after Eval of line i=1:  TypesInfo.Defs[k=0]='S' = '0xc4200818b0'

 on j=0, after Eval of line i=1:  TypesInfo.Defs[k=1]='' = '%!p(<nil>)'

parser.go:352 2018-02-08 15:28:42.685 +0700 ICT fileOrNil past the switch on p.tok

*/

// [1] Defs contains 'S' twice, at different addresses:
/*

 on j=1, after Eval of line i=3:  TypesInfo.Defs[k=0]='' = '%!p(<nil>)'

 on j=1, after Eval of line i=3:  TypesInfo.Defs[k=1]='S' = '0xc4200818b0' << - orig.

 on j=1, after Eval of line i=3:  TypesInfo.Defs[k=2]='s' = '0xc420081db0'

 on j=1, after Eval of line i=3:  TypesInfo.Defs[k=3]='' = '%!p(<nil>)'

 on j=1, after Eval of line i=3:  TypesInfo.Defs[k=4]='s' = '0xc420130500'

 on j=1, after Eval of line i=3:  TypesInfo.Defs[k=5]='S' = '0xc420130fa0'

 on j=1, after Eval of line i=3:  TypesInfo.Defs[k=6]='' = '%!p(<nil>)'

 on j=1, after Eval of line i=3:  TypesInfo.Defs[k=7]='' = '%!p(<nil>)'

 on j=1, after Eval of line i=3:  TypesInfo.Defs[k=8]='Hi' = '0xc420130370'

 on j=1, after Eval of line i=3:  TypesInfo.Defs[k=9]='' = '%!p(<nil>)'

 on j=1, after Eval of line i=3:  TypesInfo.Defs[k=10]='Hi' = '0xc420131a90'

 on j=1, after Eval of line i=3:  TypesInfo.Defs[k=11]='' = '%!p(<nil>)'

 on j=1, after Eval of line i=3:  TypesInfo.Defs[k=12]='s' = '0xc4201314a0'

 on j=1, after Eval of line i=3:  TypesInfo.Defs[k=13]='' = '%!p(<nil>)'

 on j=1, after Eval of line i=3:  TypesInfo.Defs[k=14]='s' = '0xc420131c20'

parser.go:744 2018-02-08 15:28:46.46 +0700 ICT jea debug: about to call p.pexpr() at the end of unaryExpr

...

then later, with the 2nd def of 'S':


 on j=1, after Eval of line i=3:  TypesInfo.Defs[k=0]='' = '%!p(<nil>)'

 on j=1, after Eval of line i=3:  TypesInfo.Defs[k=1]='S' = '0xc4200818b0'

 on j=1, after Eval of line i=3:  TypesInfo.Defs[k=2]='s' = '0xc420081db0'

 on j=1, after Eval of line i=3:  TypesInfo.Defs[k=3]='' = '%!p(<nil>)'

 on j=1, after Eval of line i=3:  TypesInfo.Defs[k=4]='s' = '0xc420130500'

 on j=1, after Eval of line i=3:  TypesInfo.Defs[k=5]='S' = '0xc420130fa0'

 on j=1, after Eval of line i=3:  TypesInfo.Defs[k=6]='' = '%!p(<nil>)'

 on j=1, after Eval of line i=3:  TypesInfo.Defs[k=7]='' = '%!p(<nil>)'

 on j=1, after Eval of line i=3:  TypesInfo.Defs[k=8]='Hi' = '0xc420130370'

 on j=1, after Eval of line i=3:  TypesInfo.Defs[k=9]='' = '%!p(<nil>)'

 on j=1, after Eval of line i=3:  TypesInfo.Defs[k=10]='Hi' = '0xc420131a90'

 on j=1, after Eval of line i=3:  TypesInfo.Defs[k=11]='' = '%!p(<nil>)'

 on j=1, after Eval of line i=3:  TypesInfo.Defs[k=12]='s' = '0xc4201314a0'

 on j=1, after Eval of line i=3:  TypesInfo.Defs[k=13]='' = '%!p(<nil>)'

 on j=1, after Eval of line i=3:  TypesInfo.Defs[k=14]='s' = '0xc420131c20'


*/
