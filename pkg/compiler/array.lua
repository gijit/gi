-- arrays

-- create private index
_giPrivateRaw = _giPrivateRaw or {}
_giPrivateArrayProps = _giPrivateArrayProps or {}

_giPrivateArrayMt = {

   __newindex = function(t, k, v)
      local props = rawget(t, _giPrivateArrayProps)
      local len = props.len
      --print("newindex called for key", k, " len at start is ", len)
      local raw = rawget(t, _giPrivateRaw)
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
      props.len = len
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
      return rawget(t, _giPrivateRaw)[k]
   end,

   __tostring = function(t)
      local props = rawget(t, _giPrivateArrayProps)
      local len =  props.len
      local s = "array <len= " .. tostring(len) .. "> is _gi_Array{"
      --print("t.len is", len)
      local raw = rawget(t, _giPrivateRaw)
      -- we want to skip both the _giPrivateRaw and the len
      -- when iterating, which happens automatically if we
      -- iterate on r, the inside private data, and not on the proxy.
      quo = ""
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
      local props = rawget(t, _giPrivateArrayProps)
      return props.len
   end,

   __pairs = function(t)
      -- print("__pairs called!")
      -- this makes a _giArray work in a for k,v in pairs() do loop.

      -- Iterator function takes the table and an index and returns the next index and associated value
      -- or nil to end iteration

      local function stateless_iter(t, k)
         local v
         --  Implement your own key,value selection logic in place of next
         k, v = next(rawget(t, _giPrivateRaw), k)
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

--gohperjs:
--var $clone = function(src, type) {
--  var clone = type.zero();
--  type.copy(clone, src);
--  return clone;
--};

function __clone(src, typ)
   print("__clone called with typ ", typ)
   local clone = typ.zero();
   typ.copy(clone, src);
   return clone;
end

-- jea: my earlier Proof of concept.
function __gi_clone(t, typ)
    print("_gi_clone called with typ ", typ)
    if type(t) ~= 'table' then
       error "__gi_clone called on non-table"
    end
 
    if typ == "kind_arrayType" then
       local props = rawget(t, _giPrivateArrayProps)
       if props == nil then
          error "__gi_clone for arrayType could not get props" 
       end
       -- make a copy of the data
       local src = rawget(t, _giPrivateRaw)
       local dest = {}
       for i,v in pairs(src) do
          dest[i] = v
       end
       -- unpack ignores the [0] value, so less useful.
       
       local b = __gi_NewArray(dest, props.typeKind, props.len)
       return b
    end
    print("unimplemented typ in __gi_clone: '"..tostring(typ).."'") -- Beagle for 028
    print(debug.traceback())
    error("unimplemented typ in __gi_clone: '"..tostring(typ).."'")
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
   
   raw = rawget(t, _giPrivateRaw)
   
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
   
   return rawget(t, _giPrivateRaw)
end
