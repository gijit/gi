/*
package muse provides for type
punning (converting) from types.Type
to reflect.Type.
*/
package muse

import (
	"fmt"
	"reflect"
	"unsafe"

	//"github.com/gijit/gi/pkg/token"
	"github.com/gijit/gi/pkg/types"
	"github.com/gijit/gi/pkg/verb"
)

var p = verb.P
var pp = verb.PP

// convert from types.Type to reflect.Type, so
// that we can wrap Go slices/arrays with Lua
// proxies from the very start of their creation.

type Muse struct{}

func NewMuse() *Muse { return &Muse{} }

func (m *Muse) Pun(tt types.Type) (rt reflect.Type, err error) {

	switch x := tt.(type) {
	case *types.Basic:
		return m.punBasic(x)
	case *types.Pointer:
		et, err := m.Pun(x.Elem())
		if err != nil {
			return nil, err
		}
		return reflect.PtrTo(et), nil
	case *types.Array:
		n := int(x.Len())
		et, err := m.Pun(x.Elem())
		if err != nil {
			return nil, err
		}
		return reflect.ArrayOf(n, et), nil
	case *types.Slice:
		et, err := m.Pun(x.Elem())
		if err != nil {
			return nil, err
		}
		return reflect.SliceOf(et), nil
	case *types.Map:
		kt, err := m.Pun(x.Key())
		if err != nil {
			return nil, err
		}
		et, err := m.Pun(x.Elem())
		if err != nil {
			return nil, err
		}
		return reflect.MapOf(kt, et), nil
	case *types.Chan:
		var dir reflect.ChanDir
		switch x.Dir() {
		case types.SendRecv:
			dir = reflect.BothDir
		case types.SendOnly:
			dir = reflect.SendDir
		case types.RecvOnly:
			dir = reflect.RecvDir
		default:
			panic(fmt.Errorf("unimplemented channel direction: '%#v'/'%T'",
				x.Dir(), x.Dir()))
		}
		et, err := m.Pun(x.Elem())
		if err != nil {
			return nil, err
		}
		return reflect.ChanOf(dir, et), nil
	case *types.Struct:
		nf := x.NumFields()
		fields := make([]reflect.StructField, nf)
		for i := 0; i < nf; i++ {
			f := x.Field(i) // *types.Var
			anon := f.Anonymous()
			isField := f.IsField()
			if !isField {
				panic(fmt.Errorf("huh? why isn't this a field?: '%T'/'%#v'", f, f))
			}
			pkg := f.Pkg()   // *types.Package
			name := f.Name() // string
			ftyp := f.Type() // types.Type
			rftyp, err := m.Pun(ftyp)
			if err != nil {
				return nil, err
			}
			// exported := f.Exported() // bool
			// id := f.Id() //string; Id(obj.pkg, obj.name)

			tag := x.Tag(i) // string
			fields[i] = reflect.StructField{
				Name:    name,
				PkgPath: pkg.Path(),
				Type:    rftyp,
				Tag:     reflect.StructTag(tag),

				// jea: no idea what Offset should be set to.
				Offset: 0, //    uintptr   // offset within struct, in bytes

				// jea: not sure if this is correct, in particular
				// when embedded structs are present.
				Index: []int{i}, //     []int     // index sequence for Type.FieldByIndex

				Anonymous: anon,
			}
		}
		return reflect.StructOf(fields), nil
	case *types.Tuple:
		/*
			// Tuples also represent the types of
			// the parameter list and the result list of a function
			n := x.Len()
			rt = []reflect.Type{}
			for i := 0; i < n; i++ {
				v := x.At(i)
				vrt, err := m.punVar(v)
				if err != nil {
					return nil, err
				}
				rt=append(rt, vrt)
			}
		*/

	case *types.Signature:
		// reflect.FuncOf(in, out []Type, variadic bool) Type
	case *types.Named:
		pp("in *types.Named: x ='%#v'", x)
		// x = '&types.Named{obj:(*types.TypeName)(0xc42007d450), underlying:(*types.Struct)(0xc42000ec90), methods:[]*types.Func(nil)}'

		obj := x.Obj() // returns *type.TypeName
		pkg := obj.Pkg()
		nm := obj.Name()
		ttyp := obj.Type()
		under := x.Underlying()

		numMeth := x.NumMethods()
		//x.Method(i)

		pp("obj = '%#v'", obj)
		pp("pkg = '%#v'", pkg)
		pp("nm = '%#v'", nm)
		pp("ttyp = '%#v'", ttyp)
		pp("under = '%#v'", under)
		pp("numMeth = '%#v'", numMeth)

		/*
		   util.go:178 2018-02-01 13:30:44.021 +0700 ICT we're in the token.TYPE!

		   util.go:182 2018-02-01 13:30:44.021 +0700 ICT o='&types.TypeName{object:types.object{parent:(*types.Scope)(0xc420056540), pos:7, pkg:(*types.Package)(0xc42007d400), name:"Tree", typ:(*types.Named)(0xc42000ebd0), order_:0x1, scopePos_:0}}'

		   muse_test.go:94 2018-02-01 13:30:44.021 +0700 ICT checked = '&types.Named{obj:(*types.TypeName)(0xc42007d450), underlying:(*types.Struct)(0xc42000ec90), methods:[]*types.Func(nil)}'

		   muse.go:134 2018-02-01 13:30:44.046 +0700 ICT in *types.Named: x ='&types.Named{obj:(*types.TypeName)(0xc42007d450), underlying:(*types.Struct)(0xc42000ec90), methods:[]*types.Func(nil)}'

		   muse.go:146 2018-02-01 13:30:44.046 +0700 ICT obj = '&types.TypeName{object:types.object{parent:(*types.Scope)(0xc420056540), pos:7, pkg:(*types.Package)(0xc42007d400), name:"Tree", typ:(*types.Named)(0xc42000ebd0), order_:0x1, scopePos_:0}}'

		   muse.go:147 2018-02-01 13:30:44.046 +0700 ICT pkg = '&types.Package{path:"", name:"", scope:(*types.Scope)(0xc420056540), complete:true, imports:[]*types.Package(nil), fake:false}'

		   muse.go:148 2018-02-01 13:30:44.046 +0700 ICT nm = '"Tree"'

		   muse.go:149 2018-02-01 13:30:44.046 +0700 ICT ttyp = '&types.Named{obj:(*types.TypeName)(0xc42007d450), underlying:(*types.Struct)(0xc42000ec90), methods:[]*types.Func(nil)}'

		   muse.go:150 2018-02-01 13:30:44.046 +0700 ICT under = '&types.Struct{fields:[]*types.Var{(*types.Var)(0xc42007d680)}, tags:[]string(nil)}'

		   muse.go:151 2018-02-01 13:30:44.06 +0700 ICT numMeth = '0'

		*/

		panic("not done: *types.Named")
	case *types.Interface:
	default:
		panic(fmt.Sprintf("unknown types.Type '%T'", tt))
	}
	panic(fmt.Errorf("unimplemented muse.Pun handling for type '%T'", tt))
}

func (m *Muse) punBasic(tt *types.Basic) (rt reflect.Type, err error) {
	k := tt.Kind()
	var x interface{}
	switch k {
	case types.Invalid:
		panic("invalid types.Basic")
	case types.Bool:
		x = false
	case types.Int:
		x = int(0)
	case types.Int8:
		x = int8(0)
	case types.Int16:
		x = int16(0)
	case types.Int32: // types.Rune is an alias
		x = int32(0)
	case types.Int64:
		x = int64(0)
	case types.Uint:
		x = uint(0)
	case types.Uint8: // types.Byte is an alias.
		x = uint8(0)
	case types.Uint16:
		x = uint16(0)
	case types.Uint32:
		x = uint32(0)
	case types.Uint64:
		x = uint64(0)
	case types.Uintptr:
		x = uintptr(0)
	case types.Float32:
		x = float32(0)
	case types.Float64:
		x = float64(0)
	case types.Complex64:
		x = complex64(0)
	case types.Complex128:
		x = complex128(0)
	case types.String:
		x = ""
	case types.UnsafePointer:
		x = unsafe.Pointer(m)
	case types.UntypedBool,
		types.UntypedInt,
		types.UntypedRune,
		types.UntypedFloat,
		types.UntypedComplex,
		types.UntypedString,
		types.UntypedNil:
		panic("can't take TypeOf an unfinished type.")
	}
	return reflect.TypeOf(x), nil
}

func (m *Muse) punVar(f *types.Var) (rt reflect.Type, err error) {
	ftyp := f.Type() // types.Type
	rftyp, err := m.Pun(ftyp)
	if err != nil {
		return nil, err
	}
	return rftyp, nil
}
