package compiler

import (
	"fmt"
	"os"
	"testing"

	"github.com/glycerine/gi/pkg/types"

	cv "github.com/glycerine/goconvey/convey"
)

func Test1200DepthFirstSearchOfTypeDependencies(t *testing.T) {

	cv.Convey("dfs on type tree should work", t, func() {

		dependTestMode = true
		testDFS := func() {
			s := NewDFSState()

			// verify that reset()
			// works by doing two passes.

			for i := 1; i < 2; i++ {
				s.reset()
				anInt := types.Typ[types.Int]

				aTn := types.NewTypeName(0, nil, "A", anInt)
				aTy := types.NewNamed(aTn, anInt, nil)

				bTn := types.NewTypeName(0, nil, "B", anInt)
				bTy := types.NewNamed(bTn, anInt, nil)

				cv.So(aTy, cv.ShouldNotEqual, bTy)

				cTn := types.NewTypeName(0, nil, "C", anInt)
				cTy := types.NewNamed(cTn, anInt, nil)

				dTn := types.NewTypeName(0, nil, "D", anInt)
				dTy := types.NewNamed(dTn, anInt, nil)

				eTn := types.NewTypeName(0, nil, "E", anInt)
				eTy := types.NewNamed(eTn, anInt, nil)

				fTn := types.NewTypeName(0, nil, "F", anInt)
				fTy := types.NewNamed(fTn, anInt, nil)

				gTn := types.NewTypeName(0, nil, "G", anInt)
				gTy := types.NewNamed(gTn, anInt, nil)

				a := s.newDfsNode("a", aTy, []byte("//test code for a"))

				adup := s.newDfsNode("a", aTy, []byte("//test code for adup"))
				if adup != a {
					panic("dedup failed.")
				}

				var b = s.newDfsNode("b", bTy, []byte("//test code for b"))
				var c = s.newDfsNode("c", cTy, []byte("//test code for c"))
				var d = s.newDfsNode("d", dTy, []byte("//test code for d"))
				var e = s.newDfsNode("e", eTy, []byte("//test code for e"))
				var f = s.newDfsNode("f", fTy, []byte("//test code for f"))

				// separate island:
				var g = s.newDfsNode("g", gTy, []byte("//test code for g"))

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
