--  depend.lua:
--
--  Implement Depth-First-Search (DFS)
--  on the graph of depedencies
--  between types. A pre-order
--  traversal will print
--  leaf types before the compound
--  types that need them defined.


function __newDfsNode(self, name, payload)
   if payload == nil then
      error "payload cannot be nil in __newDfsNode"
   end
   local nd = self.dfsDedup[payload]
   if nd ~= nil then
      return nd
   end
   local node= {
      visited=false,
      children={},
      dedupChildren={},
      id = self.dfsNextID,
      name=name,
      payload=payload,
   }
   self.dfsNextID=self.dfsNextID+1
   self.dfsDedup[payload] = node
   table.insert(self.dfsNodes, node)
   
   return node
end

-- par should be a node; e.g. typ.__dfsNode
function __addChild(self, par, ch)

   -- we can skip all basic types,
   -- as they are already defined.
   if not ch.named then
      if ch.kind <= 16 or -- __kindComplex128
         ch.kind == 24 or -- __kindString
      ch.kind == 26 then  -- __kindUnsafePointer
         return
      end
   end
   
   if par.dedupChildren[ch] ~= nil then
      -- avoid adding same child twice.
      return
   end
   par.dedupChildren[ch]= #par.children+1
   table.insert(par.children, ch)
end

function __markGraphUnVisited(self)
   self.dfsOrder = {}
   for _,n in ipairs(self.dfsNodes) do
      n.visited = false
   end
end

function __emptyOutGraph(self)
   self.dfsOrder = {}
   self.dfsNodes = {} -- node stored in value.
   self.dfsDedup = {} -- payload key -> node value.
end

function __dfsHelper(self, node)
   if node == nil then
      return
   end
   if node.visited then
      return
   end
   node.visited = true
   for _, ch in ipairs(node.children) do
      self:dfsHelper(ch)
   end
   print("post-order visit sees node "..tostring(node.id).." : "..node.name)
   table.insert(self.dfsOrder, node)
end

function __doDFS(self)
   __markGraphUnVisited(self)
   for _, n in ipairs(self.dfsNodes) do
      self:dfsHelper(n)
   end
end

function __NewDFSState()
   return {
      dfsNodes = {},
      dfsOrder = {},
      dfsDedup = {},
      dfsNextID = 0,

      doDFS = __doDFS,
      dfsHelper = __dfsHelper,
      emptyOutGraph = __emptyOutGraph,
      newDfsNode = __newDfsNode,
      addChild = __addChild,
      markGraphUnVisited = __markGraphUnVisited,
   }
end

--[[
-- test. To test, change the --[[ above to ---[[
--       and issue dofile('dfs.lua')
dofile 'tutil.lua'

function __testDFS()
   local s = __NewDFSState()

   -- verify that __emptyOutGraph
   -- works by doing two passes.
   
   for i =1,2 do
      s:emptyOutGraph()
      
      local aPayload = {}
      local a = s:newDfsNode("a", aPayload)
   
      local adup = s:newDfsNode("a", aPayload)
      if adup ~= a then
          error "dedup failed."
      end

      local b = s:newDfsNode("b", {})
      local c = s:newDfsNode("c", {})
      local d = s:newDfsNode("d", {})
      local e = s:newDfsNode("e", {})
      local f = s:newDfsNode("f", {})

      -- separate island:
      local g = s:newDfsNode("g", {})
      
      s:addChild(a, b)

      -- check dedup of child add
      local startCount = #a.children
      s:addChild(a, b)
      if #a.children ~= startCount then
          error("child dedup failed.")
      end

      s:addChild(b, c)
      s:addChild(b, d)
      s:addChild(d, e)
      s:addChild(d, f)

      s:doDFS()

      for i, n in ipairs(s.dfsOrder) do
         print("dfs order "..i.." is "..tostring(n.id).." : "..n.name)
      end
      
      expectEq(s.dfsOrder[1], c)
      expectEq(s.dfsOrder[2], e)
      expectEq(s.dfsOrder[3], f)
      expectEq(s.dfsOrder[4], d)
      expectEq(s.dfsOrder[5], b)
      expectEq(s.dfsOrder[6], a)
      expectEq(s.dfsOrder[7], g)
   end
   
end
__testDFS()
__testDFS()
--]]
