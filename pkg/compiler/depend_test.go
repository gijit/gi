package compiler

import (
	"fmt"
	"os"
	"testing"

	"github.com/glycerine/gi/pkg/types"

	cv "github.com/glycerine/goconvey/convey"
)

type devNull int

func (devNull) Write(p []byte) (int, error) {
	return len(p), nil
}

func Test1200DepthFirstSearchOfTypeDependencies(t *testing.T) {

	cv.Convey("dfs on type tree should work", t, func() {

		dependTestMode = true
		testDFS := func() {
			s := NewDFSState()

			// verify that reset()
			// works by doing two passes.

			for i := 1; i < 2; i++ {
				s.reset()
				//anInt := types.Typ[types.Int]

				aTn := types.NewTypeName(0, nil, "A", types.Typ[types.Int])
				aTy := types.NewNamed(aTn, types.Typ[types.Int], nil)

				bTn := types.NewTypeName(0, nil, "B", types.Typ[types.String])
				bTy := types.NewNamed(bTn, types.Typ[types.String], nil)

				cv.So(aTy, cv.ShouldNotEqual, bTy)

				cTn := types.NewTypeName(0, nil, "C", types.Typ[types.Int8])
				cTy := types.NewNamed(cTn, types.Typ[types.Int8], nil)

				dTn := types.NewTypeName(0, nil, "D", types.Typ[types.Int16])
				dTy := types.NewNamed(dTn, types.Typ[types.Int16], nil)

				eTn := types.NewTypeName(0, nil, "E", types.Typ[types.Int32])
				eTy := types.NewNamed(eTn, types.Typ[types.Int32], nil)

				fTn := types.NewTypeName(0, nil, "F", types.Typ[types.Int64])
				fTy := types.NewNamed(fTn, types.Typ[types.Int64], nil)

				gTn := types.NewTypeName(0, nil, "G", types.Typ[types.Uint])
				gTy := types.NewNamed(gTn, types.Typ[types.Uint], nil)

				a := s.newDfsNode("a", aTy, []byte("//test code for a\n"))

				adup := s.newDfsNode("a", aTy, []byte("//test code for adup"))
				if adup != a {
					panic("dedup failed.")
				}

				var b = s.newDfsNode("b", bTy, []byte("//test code for b\n"))
				var c = s.newDfsNode("c", cTy, []byte("//test code for c\n"))
				var d = s.newDfsNode("d", dTy, []byte("//test code for d\n"))
				var e = s.newDfsNode("e", eTy, []byte("//test code for e\n"))
				var f = s.newDfsNode("f", fTy, []byte("//test code for f\n"))

				// separate island:
				var g = s.newDfsNode("g", gTy, []byte("//test code for g\n"))

				s.addChild(a, b)

				// check dedup of child add
				var startCount = len(a.children)
				s.addChild(a, b)
				if len(a.children) != startCount {
					panic("child dedup failed.")
				}

				s.addChild(b, c)
				s.addChild(b, d)
				s.addChild(d, e)
				s.addChild(d, f)

				// cycle rejected
				cv.So(func() {
					s.addChild(c, a)
				}, cv.ShouldPanic)

				s.doDFS()

				s.showDFSOrder()

				expectEq(s.dfsOrder[0], c)
				expectEq(s.dfsOrder[1], e)
				expectEq(s.dfsOrder[2], f)
				expectEq(s.dfsOrder[3], d)
				expectEq(s.dfsOrder[4], b)
				expectEq(s.dfsOrder[5], a)
				expectEq(s.dfsOrder[6], g)

				s.genCode(os.Stdout)
				//s.genCode(devNull(0))
			}

		}
		testDFS()
		testDFS()

		cv.So(true, cv.ShouldBeTrue)
	})
}

func expectEq(a, b *dfsNode) {
	if a != b {
		panic(fmt.Sprintf("ouch: expected equal: %#v and %#v", a.name, b.name))
	}
}
