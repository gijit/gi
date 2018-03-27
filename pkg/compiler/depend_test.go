package compiler

import (
	//"fmt"
	"testing"

	"github.com/glycerine/gi/pkg/types"

	cv "github.com/glycerine/goconvey/convey"
)

func Test1200DepthFirstSearchOfTypeDependencies(t *testing.T) {

	cv.Convey("dfs on type tree should work", t, func() {

		testDFS := func() {
			s := NewDFSState()

			// verify that reset()
			// works by doing two passes.

			for i := 1; i < 2; i++ {
				s.reset()

				aPayload := types.Typ[types.Int]
				a := s.newDfsNode("a", aPayload)

				adup := s.newDfsNode("a", aPayload)
				if adup != a {
					panic("dedup failed.")
				}

				var b = s.newDfsNode("b", nil)
				var c = s.newDfsNode("c", nil)
				var d = s.newDfsNode("d", nil)
				var e = s.newDfsNode("e", nil)
				var f = s.newDfsNode("f", nil)

				// separate island:
				var g = s.newDfsNode("g", nil)

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
			}

		}
		testDFS()
		testDFS()

		cv.So(true, cv.ShouldBeTrue)
	})
}

func expectEq(a, b *dfsNode) {
	if a != b {
		panic("ouch")
	}
}
