-- arrays

-- create private index
_giPrivateRaw = _giPrivateRaw or {}
_giPrivateArrayProps = {}

_giPrivateArrayMt = {

    __newindex = function(t, k, v)
      local len = t[_giPrivateArrayProps]["len"]
      --print("newindex called for key", k, " len at start is ", len)
      if t[_giPrivateRaw][k] == nil then
         if  v ~= nil then
         -- new value
            len = len +1
         end
      else
         -- key already present, are we replacing or deleting?
          if v == nil then 
              len = len - 1 -- delete
          else
              -- replace, no count change              
          end
      end
      t[_giPrivateRaw][k] = v
      t[_giPrivateArrayProps]["len"] = len
      --print("len at end of newidnex is ", len)
    end,

   -- __index allows us to have fields to access the count.
   --
   __index = function(t, k)
      --print("index called for key", k)
      -- I don't think we need raw any more, now that
      -- __pairs works. It would be a problem for hash
      -- tables that want to store the key 'raw'.
      -- if k == 'raw' then return t[_giPrivateRaw] end
      return t[_giPrivateRaw][k]
   end,

   __tostring = function(t)
      local len =  t[_giPrivateArrayProps]["len"]
      local s = "array of length " .. tostring(len) .. " is _gi_Array{"
      --print("t.len is", len)
      local r = t[_giPrivateRaw]
      -- we want to skip both the _giPrivateRaw and the len
      -- when iterating, which happens automatically if we
      -- iterate on r, the inside private data, and not on the proxy.
      for i, _ in pairs(r) do s = s .. "["..tostring(i).."]" .. "= " .. tostring(r[i]) .. ", " end
      return s .. "}"
   end,

   __len = function(t)
      -- this does get called by the # operation, for slices.
      --print("len called")
      return t[_giPrivateArrayProps]["len"]
   end,

   __pairs = function(t)
      -- print("__pairs called!")
      -- this makes a _giArray work in a for k,v in pairs() do loop.

      -- Iterator function takes the table and an index and returns the next index and associated value
      -- or nil to end iteration

      local function stateless_iter(t, k)
         local v
         --  Implement your own key,value selection logic in place of next
         k, v = next(t[_giPrivateRaw], k)
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
   end
}

function _gi_NewArray(x, typeKind, len, zeroVal)
   --print("_gi_NewArray constructor called with x=",x," and typeKind=", typeKind," len=", len)
   if typeKind == nil or typeKind == "" then
      error("must provide typeKind to _gi_NewArray")
   end
   
   if len == nil then
      error("must provide len to _gi_NewArray")
   end

   local proxy = {}
   proxy[_giPrivateRaw] = x

   -- zero any tail that is not set
   if zeroVal ~= nil then
      for i =0,len-1 do
         if x[i] == nil then
            x[i] = zeroVal
         end
      end
   end
   
   local props = {len=len, typeKind=typeKind}
   proxy[_giPrivateArrayProps] = props

   --print("upon init, len is ",proxy[_giPrivateArrayProps]["len"])
   
   setmetatable(proxy, _giPrivateArrayMt)
   return proxy
end;

function _gi_clone(t, typ)
   print("_gi_clone called with typ ", typ)
   if type(t) ~= 'table' then
      error "_gi_clone called on non-table"
   end

   if typ == "arrayType" then
      local props = t[_giPrivateArrayProps]
      if props == nil then
         error "_gi_clone for arrayType could not get props" 
      end
      -- make a copy of the data
      local src = t[_giPrivateRaw]
      local dest = {}
      for i,v in pairs(src) do
         dest[i] = v
      end
      -- unpack ignores the [0] value, so less useful.
      
      local b = _gi_NewArray(dest, props.typeKind, props.len)
      return b
   end
   print("unimplemented typ in _gi_clone:", typ)
   error "unimplemented typ in _gi_clone"
end


-- _gi_UnpackArrayRaw is a helper, used in
-- generated Lua code.
--
-- This helper unpacks the raw _giArray
-- arguments. It returns non tables unchanged,
-- and non _giArray tables unpacked.
--
function _gi_UnpackArrayRaw(t)
   if type(t) ~= 'table' then
      return t
   end
   
   raw = t[_giPrivateRaw]
   
   if raw == nil then
      -- unpack of empty table is ok. returns nil.
      return unpack(t) 
   end

   if #raw == 0 then
      return nil
   end
   return raw[0], unpack(raw)
end

function _gi_ArrayRaw(t)
   if type(t) ~= 'table' then
      return t
   end
   
   return t[_giPrivateRaw]
end
