--  depend.lua:
--
--  Implement Depth-First-Search (DFS)
--  on the graph of depedencies
--  between types. A pre-order
--  traversal will print
--  leaf types before the compound
--  types that need them defined.

__nodes = {}
__dfsOrder = {}
__nextID = 0

function __newNode(name)
   local node= {
      visited=false,
      children={},
      id = __nextID,
      name=name,
   }
   __nextID=__nextID+1
   table.insert(__nodes, node)
   return node
end

function __addChild(par, ch)
   table.insert(par.children, ch)
end

function __markGraphUnVisited()
   __dfsOrder = {}
   for _,n in ipairs(__nodes) do
      n.visited = false
   end
end

function __emptyOutGraph()
   __dfsOrder = {}
   __nodes = {}
end

function __dfsHelper(node)
   if node == nil then
      return
   end
   if node.visited then
      return
   end
   node.visited = true
   for _, ch in ipairs(node.children) do
      __dfsHelper(ch)
   end
   print("post-order visit sees node "..tostring(node.id).." : "..node.name)
   table.insert(__dfsOrder, node)
end

function __doDFS()
   __markGraphUnVisited()
   for _, n in ipairs(__nodes) do
      __dfsHelper(n)
   end
end

-- test
dofile 'tutil.lua'

function __testDFS()
   __emptyOutGraph()
   
   local a = __newNode("a")
   local b = __newNode("b")
   local c = __newNode("c")
   local d = __newNode("d")
   local e = __newNode("e")
   local f = __newNode("f")

   -- separate island:
   local g = __newNode("g")
   
   __addChild(a, b)
   __addChild(b, c)
   __addChild(b, d)
   __addChild(d, e)
   __addChild(d, f)

   __doDFS()

   for i, n in ipairs(__dfsOrder) do
      print("dfs order "..i.." is "..tostring(n.id).." : "..n.name)
   end
   
   expectEq(__dfsOrder[1], c)
   expectEq(__dfsOrder[2], e)
   expectEq(__dfsOrder[3], f)
   expectEq(__dfsOrder[4], d)
   expectEq(__dfsOrder[5], b)
   expectEq(__dfsOrder[6], a)
   expectEq(__dfsOrder[7], g)

end
__testDFS()
__testDFS()
