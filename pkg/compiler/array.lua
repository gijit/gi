-- arrays

-- create private index
_giPrivateArrayRaw = {}

_giPrivateArrayMt = {

   __newindex = function(t, k, v)
      print("newindex called for key", k)
      if t[_giPrivateArrayRaw][k] ~= nil then
         -- replace or delete
         if v == nil then 
            t.len = t.len - 1 -- delete
         else
            -- replace, no count change              
         end
      else 
         t.len = t.len + 1
      end
      t[_giPrivateArrayRaw][k] = v
   end,

   -- __index allows us to have fields to access the count.
   --
   __index = function(t, k)
      --print("index called for key", k)
      -- I don't think we need raw any more, now that
      -- __pairs works. It would be a problem for hash
      -- tables that want to store the key 'raw'.
      -- if k == 'raw' then return t[_giPrivateArrayRaw] end
      return t[_giPrivateArrayRaw][k]
   end,

   __tostring = function(t)
      local s = "array of length " .. tostring(t.len) .. " is _gi_Array{"
      local r = t[_giPrivateArrayRaw]
      -- we want to skip both the _giPrivateArrayRaw and the len
      -- when iterating, which happens automatically if we
      -- iterate on r, the inside private data, and not on the proxy.
      for i, _ in pairs(r) do s = s .. "["..tostring(i).."]" .. "= " .. tostring(r[i]) .. ", " end
      return s .. "}"
   end,

   __len = function(t)
      -- this does get called by the # operation(!)
      -- print("len called")
      return t.len
   end,

   __pairs = function(t)
      -- print("__pairs called!")
      -- this makes a _giArray work in a for k,v in pairs() do loop.

      -- Iterator function takes the table and an index and returns the next index and associated value
      -- or nil to end iteration

      local function stateless_iter(t, k)
         local v
         --  Implement your own key,value selection logic in place of next
         k, v = next(t[_giPrivateArrayRaw], k)
         if v then return k,v end
      end

      -- Return an iterator function, the table, starting point
      return stateless_iter, t, nil
   end,

   __call = function(t, ...)
      print("__call() invoked, with ... = ", ...)
      local oper, key = ...
      print("oper is", oper)
      print("key is ", key)
      if oper == "delete" then
         -- the hash table delete operation
         if key == nil then
            error("delete error: key to delete cannot be nil")
         end
         -- forward the actual delete
         t[key]=nil
      elseif oper == "slice" then

      end
   end
}

function _gi_NewArray(len)
   x={}

   if len == nil then
      error("must provide size to _gi_NewArray")
   end
   
   local proxy = {len=length, typeKind=kind}
   proxy[_giPrivateArrayRaw] = x
   setmetatable(proxy, _giPrivateArrayMt)
   return proxy
end;

arrayType = {
   zero = function()
      return _gi_NewArray(0)
   end
}
