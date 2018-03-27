package compiler

import (
	"fmt"
	"io"

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
//  See dfsState.genCode() for that.

var dependTestMode = false

func isBasicTyp(n *dfsNode) bool {
	_, ok := n.typ.(*types.Basic)
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

	createCode []byte
}

func (me *dfsNode) bloom(w io.Writer, s *dfsState) {
	/*
		if len(me.createCode) > 0 {
			_, err := w.Write(me.createCode)
			panicOn(err)
		}
	*/
	if len(s.codeMap) > 0 {
		code, ok := s.codeMap[me.typ]
		if ok {
			_, err := w.Write([]byte(code))
			panicOn(err)
		}
	}
}

// a func on nodes to force instantiation of
// any types this node depends on, i.e. those
// types (not values) that were described but
// lazily instantated. Calls me.typ.bloom
// on our subtree in depth-first order.
//
func (me *dfsNode) makeRequiredTypes(w io.Writer, s *dfsState) {
	if me.made {
		return
	}
	me.made = true
	for _, ch := range me.children {
		ch.makeRequiredTypes(w, s)
	}
	me.bloom(w, s)
}

func (s *dfsState) newDfsNode(name string, typ types.Type, createCode []byte) *dfsNode {
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
		createCode:    createCode,
	}
	s.dfsNextID++
	s.dfsDedup[typ] = node
	s.dfsNodes = append(s.dfsNodes, node)
	s.codeMap[typ] = string(createCode)

	return node
}

// par should be a node; e.g. typ.dfsNode
func (s *dfsState) addChild(par, ch *dfsNode) {

	if par == nil {
		panic("par cannot be nil in addChild")
	}
	if ch == nil {
		panic("ch cannot be nil in addChild")
	}

	// we can skip all basic types,
	// as they are already defined.
	if !dependTestMode {
		if isBasicTyp(ch) {
			return
		}
		if isBasicTyp(par) {
			panic(fmt.Sprintf("addChild error: parent was basic type. "+
				"cannot add child to basic typ %v", par))
		}
	}

	if par.dedupChildren[ch] {
		pp("avoid adding same child twice to a parent.")
		return
	}

	n := len(par.children)
	par.children = append(par.children, ch)
	par.dedupChildren[ch] = true
	s.stale = true

	if s.hasCycle() {
		// back it out, in case we recover, as the 1200 test does.
		par.children = par.children[:n]
		delete(par.dedupChildren, ch)
		panic("cycles not allowed")
	}
}

func (s *dfsState) hasCycle() bool {
	s.markGraphUnVisited()
	s.stale = true
	for _, n := range s.dfsNodes {
		ancestors := make(map[*dfsNode]bool)
		if s.hasCycleHelper(n, ancestors) {
			return true
		}
	}
	return false
}

func (s *dfsState) hasCycleHelper(node *dfsNode, ancestors map[*dfsNode]bool) bool {
	if node == nil {
		return false
	}
	if node.visited {
		return false
	}
	node.visited = true
	ancestors[node] = true

	for _, ch := range node.children {
		if ancestors[ch] {
			vv("found cycle: backedge from '%v' to '%v'", node.name, ch.name)
			return true
		}
		if s.hasCycleHelper(ch, ancestors) {
			return true
		}
	}

	// key to the cycle detection while allowing diamonds:
	delete(ancestors, node)

	return false
}

func (s *dfsState) markGraphUnVisited() {
	s.dfsOrder = []*dfsNode{}
	for _, n := range s.dfsNodes {
		n.visited = false
	}
	s.stale = false
}

func (s *dfsState) reset() {
	// empty the graph
	s.dfsOrder = []*dfsNode{}
	s.dfsNodes = []*dfsNode{}              // node stored in value.
	s.dfsDedup = map[types.Type]*dfsNode{} // payloadTyp key -> node value.
	s.dfsNextID = 0
	s.stale = false
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

	pp("post-order visit sees node %v : %v", node.id, node.name)
	s.dfsOrder = append(s.dfsOrder, node)

}

func (s *dfsState) showDFSOrder() {
	if s.stale {
		s.doDFS()
	}
	for i, n := range s.dfsOrder {
		pp("dfs order %v is %v : %v", i, n.id, n.name)
	}
}

// in depth-first order, so dependencies are
// defined before things that depend on them.
func (s *dfsState) genCode(w io.Writer) {
	if s.stale {
		s.doDFS()
	}
	for i, n := range s.dfsOrder {
		pp("dfs order %v is %v : %v", i, n.id, n.name)
		n.makeRequiredTypes(w, s)
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
	codeMap   map[types.Type]string
}

func NewDFSState() *dfsState {
	return &dfsState{
		dfsNodes: []*dfsNode{},
		dfsOrder: []*dfsNode{},
		dfsDedup: make(map[types.Type]*dfsNode),
		codeMap:  make(map[types.Type]string),
	}
}

func (s *dfsState) setCodeMap(c map[types.Type]string) {
	s.codeMap = c
}
