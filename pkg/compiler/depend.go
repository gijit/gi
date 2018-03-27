package compiler

import (
	"fmt"

	"github.com/glycerine/gi/pkg/types"
)

//  depend.go
//
//  Implement Depth-First-Search (DFS)
//  on the graph of depedencies
//  between types. A pre-order
//  traversal will print
//  leaf types before the compound
//  types that need them defined.

var dfsTestMode = false

func isBasicTyp(typ types.Type) bool {
	_, ok := typ.(*types.Basic)
	return ok
}

type dfsNode struct {
	id            int
	name          string
	typ           types.Type
	stale         bool
	made          bool
	children      []*dfsNode
	dedupChildren map[*dfsNode]bool
	visited       bool
}

func (me *dfsNode) bloom() {
	panic("TODO")
}

// a func on nodes to force instantiation of
// any types this node depends on, i.e. those
// types (not values) that were described but
// lazily instantated. Calls me.typ.bloom
// on our subtree in depth-first order.
//
func (me *dfsNode) makeRequiredTypes() {
	if me.made {
		return
	}
	me.made = true
	if isBasicTyp(me.typ) {
		return // basic types are always leaf nodes, no children.
	}

	for _, ch := range me.children {
		ch.makeRequiredTypes()
	}
	me.bloom()
}

func (s *dfsState) newDfsNode(name string, typ types.Type) *dfsNode {
	if typ == nil {
		panic("typ cannot be nil in newDfsNode")
	}

	nd, ok := s.dfsDedup[typ]
	if ok {
		return nd
	}

	node := &dfsNode{
		id:            s.dfsNextID,
		name:          name,
		typ:           typ,
		stale:         true,
		dedupChildren: make(map[*dfsNode]bool),
	}
	s.dfsNextID++
	s.dfsDedup[typ] = node
	s.dfsNodes = append(s.dfsNodes, node)

	return node
}

// par should be a node; e.g. typ.dfsNode
func (s *dfsState) addChild(me *dfsNode, parTyp, chTyp types.Type) {

	if parTyp == nil {
		panic("parTyp cannot be nil in addChild")
	}
	if chTyp == nil {
		panic("chTyp cannot be nil in addChild")
	}

	// we can skip all basic types,
	// as they are already defined.
	if isBasicTyp(chTyp) {
		return
	}
	if isBasicTyp(parTyp) {
		panic(fmt.Sprintf("addChild error: parent was basic type. "+
			"cannot add child to basic typ %v", parTyp))
	}

	_, present := s.dfsDedup[chTyp]
	if present {
		// child was previously generated, so
		// we don't need to worry about this
		// dependency
		return
	}

	parNode := s.dfsDedup[parTyp]
	if parNode == nil {
		parNode = s.newDfsNode("TODO-par", parTyp)
	}

	chNode := s.newDfsNode("TODO-ch", chTyp)
	if parNode.dedupChildren[chNode] {
		// avoid adding same child twice.
		return
	}

	//pnc := len(parNode.children)

	// jea: huh?
	//	if pnc > 0 {
	// we lazily instantiate children
	// for better diagnostics.
	//		parNode.children = nil
	//	}

	parNode.dedupChildren[chNode] = true
	parNode.children = append(parNode.children, chNode)
	s.stale = true
}

func (s *dfsState) markGraphUnVisited() {
	s.dfsOrder = []*dfsNode{}
	for _, n := range s.dfsNodes {
		n.visited = false
	}
	s.stale = false
}

func (me *dfsState) emptyOutGraph() {
	me.dfsOrder = []*dfsNode{}
	me.dfsNodes = []*dfsNode{}              // node stored in value.
	me.dfsDedup = map[types.Type]*dfsNode{} // payloadTyp key -> node value.
	me.dfsNextID = 0
	me.stale = false
}

func (s *dfsState) dfsHelper(node *dfsNode) {
	if node == nil {
		return
	}
	if node.visited {
		return
	}
	node.visited = true

	for _, ch := range node.children {
		s.dfsHelper(ch)
	}

	vv("post-order visit sees node %v : %v", node.id, node.name)
	s.dfsOrder = append(s.dfsOrder, node)

}

func (s *dfsState) showDFSOrder() {
	if s.stale {
		s.doDFS()
	}
	for i, n := range s.dfsOrder {
		vv("dfs order %v is %v : %v", i, n.id, n.name)
	}
}

func (s *dfsState) doDFS() {
	s.markGraphUnVisited()
	for _, n := range s.dfsNodes {
		s.dfsHelper(n)
	}
	s.stale = false
}

func (s *dfsState) hasTypes() bool {
	return s.dfsNextID != 0
}

type dfsState struct {
	dfsNodes  []*dfsNode
	dfsOrder  []*dfsNode
	dfsDedup  map[types.Type]*dfsNode
	dfsNextID int
	stale     bool
}

func NewDFSState() *dfsState {
	return &dfsState{
		dfsNodes: []*dfsNode{},
		dfsOrder: []*dfsNode{},
		dfsDedup: make(map[types.Type]*dfsNode),
	}
}

/*
// test. To test, change the //[[ above to //-[[
//       and issue dofile('dfs.lua')
dofile 'tutil.lua' // must be in prelude dir to test.

func __testDFS() {
   __dfsTestMode = true
    s := __NewDFSState()

   // verify that reset()
   // works by doing two passes.

   for i =1,2 do
      s.reset()

      local aPayload = nil
      local a = s.newDfsNode("a", aPayload)

      local adup = s.newDfsNode("a", aPayload)
      if adup != a {
          panic( "dedup failed.")
      }

      local b = s.newDfsNode("b", nil)
      local c = s.newDfsNode("c", nil)
      local d = s.newDfsNode("d", nil)
      local e = s.newDfsNode("e", nil)
      local f = s.newDfsNode("f", nil)

      // separate island:
      local g = s.newDfsNode("g", nil)

      s.addChild(a, b)

      // check dedup of child add
      local startCount = #a.children
      s.addChild(a, b)
      if #a.children != startCount {
          panic("child dedup failed.")
      }

      s.addChild(b, c)
      s.addChild(b, d)
      s.addChild(d, e)
      s.addChild(d, f)

      s.doDFS()

      s.showDFSOrder()

      __expectEq(s.dfsOrder[1], c)
      __expectEq(s.dfsOrder[2], e)
      __expectEq(s.dfsOrder[3], f)
      __expectEq(s.dfsOrder[4], d)
      __expectEq(s.dfsOrder[5], b)
      __expectEq(s.dfsOrder[6], a)
      __expectEq(s.dfsOrder[7], g)
   }

}
__testDFS()
__testDFS()
*/
