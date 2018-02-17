-- arrays

-- important keys
__giPrivateRaw = "__giPrivateRaw"
__giPrivateArrayProps = "__giPrivateArrayProps"
__giPrivateSliceProps = "__giPrivateSliceProps"
__giGo = "__giGo"

__giPrivateArrayMt = {

   __newindex = function(t, k, v)
      local props = rawget(t, __giPrivateArrayProps)
      local len = props.__length
      print("newindex called for key", k, " len at start is ", len)
      local raw = rawget(t, __giPrivateRaw)
      if raw[k] == nil then
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
      raw[k] = v
      props.__length = len
      --print("len at end of newidnex is ", len)
   end,

   -- __index allows us to have fields to access the count.
   --
   __index = function(t, k)
      print("index called for key", k)
      -- I don't think we need raw any more, now that
      -- __pairs works. It would be a problem for hash
      -- tables that want to store the key 'raw'.
      -- if k == 'raw' then return t[__giPrivateRaw] end
      return rawget(t, __giPrivateRaw)[k]
   end,

   __tostring = function(t)
      local props = rawget(t, __giPrivateArrayProps)
      local len =  props.__length
      local s = "array <len= " .. tostring(len) .. "> is __gi_Array{"
      --print("t.__length is", len)
      local raw = rawget(t, __giPrivateRaw)
      -- we want to skip both the __giPrivateRaw and the len
      -- when iterating, which happens automatically if we
      -- iterate on r, the inside private data, and not on the proxy.
      local quo = ""
      if len > 0 and type(raw[0]) == "string" then
         quo = '"'
      end
      
      for i, _ in pairs(raw) do
         s = s .. "["..tostring(i).."]" .. "= " .. quo..tostring(raw[i])..quo .. ", "
      end
      return s .. "}"
   end,

   __len = function(t)
      -- this does get called by the # operation
      --print("len called")
      local props = rawget(t, __giPrivateArrayProps)
      return props.__length
   end,

   __pairsDISABLED = function(t)
      print("__pairs called!")
      -- this makes a __giArray work in a for k,v in pairs() do loop.

      -- Iterator function takes the table and an index and returns the next index and associated value
      -- or nil to end iteration

      local function stateless_iter(t, k)
         local v
         --  Implement your own key,value selection logic in place of next
         local raw = rawget(t, __giPrivateRaw)
         print("raw from __giPrivateRaw is "..tostring(raw))
         k, v = next(raw, k)
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

function __gi_NewArray(x, typeKind, len, zeroVal)
   --print("__gi_NewArray constructor called with x=",x," and typeKind=", typeKind," len=", len)
   if typeKind == nil or typeKind == "" then
      error("must provide typeKind to __gi_NewArray")
   end
   
   if len == nil then
      error("must provide len to __gi_NewArray")
   end

   local proxy = {}
   proxy[__giPrivateRaw] = x

   -- zero any tail that is not set
   if zeroVal ~= nil then
      for i =0,len-1 do
         if x[i] == nil then
            x[i] = zeroVal
         end
      end
   end
   
   local props = {len=len, typeKind=typeKind}
   proxy[__giPrivateArrayProps] = props

   --print("upon init, len is ",proxy[__giPrivateArrayProps]["len"])
   
   setmetatable(proxy, __giPrivateArrayMt)
   return proxy
end;

--gohperjs:
--var $clone = function(src, type) {
--  var clone = type.zero();
--  type.copy(clone, src);
--  return clone;
--};

function __gi_clone2(src, typ)
   if typ == nil then
      print("__gi_clone2() called with nil typ!?!") -- typ='"..tostring(typ).."'")
      __st(typ)
      print(debug.traceback())
      error "don't call __gi_clone2 with a nil typ!"
   end
   
   local clone = typ.zero();
   typ.copy(clone, src);
   return clone;
end

-- jea: my earlier Proof of concept.
-- function __gi_clone(t, typ)
--     print("__gi_clone called with typ ", typ)
--     if type(t) ~= 'table' then
--        error "___gi_clone called on non-table"
--     end
--  
--     if typ == "kind_arrayType" then
--        local props = rawget(t, __giPrivateArrayProps)
--        if props == nil then
--           error "__gi_clone for arrayType could not get props" 
--        end
--        -- make a copy of the data
--        local src = rawget(t, __giPrivateRaw)
--        local dest = {}
--        for i,v in pairs(src) do
--           dest[i] = v
--        end
--        -- unpack ignores the [0] value, so less useful.
--        
--        local b = __gi_NewArray(dest, props.typeKind, props.__length)
--        return b
--     end
--     print("unimplemented typ in __gi_clone: '"..tostring(typ).."'") -- Beagle for 028
--     print(debug.traceback())
--     error("unimplemented typ in __gi_clone: '"..tostring(typ).."'")
--  end


-- __gi_UnpackArrayRaw is a helper, used in
-- generated Lua code.
--
-- This helper unpacks the raw __giArray
-- arguments. It returns non tables unchanged,
-- and non __giArray tables unpacked.
--
function __gi_UnpackArrayRaw(t)
   if type(t) ~= 'table' then
      return t
   end
   
   local raw = rawget(t, __giPrivateRaw)
   
   if raw == nil then
      -- unpack of empty table is ok. returns nil.
      return unpack(t) 
   end

   if #raw == 0 then
      return nil
   end
   return raw[0], unpack(raw)
end

function __gi_ArrayRaw(t)
   if type(t) ~= 'table' then
      return t
   end
   
   return rawget(t, __giPrivateRaw)
end
