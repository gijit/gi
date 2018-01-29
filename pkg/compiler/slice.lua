-- like a Go slice, a lua slice needs to point
-- to a backing array


-- create private index
_giPrivateRaw = _giPrivateRaw or {}
_giPrivateSliceProps = _giPrivateSliceProps or {}
_giGo = _giGo or {}

_giPrivateSliceMt = {

    __newindex = function(t, k, v)
      --print("newindex called for key", k, " val=", v)
      local len = t[_giPrivateSliceProps]["len"]
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
      t[_giPrivateSliceProps]["len"] = len
      --print("len at end of newidnex is ", len)
    end,

  -- __index allows us to have fields to access the count.
  --
    __index = function(t, k)
       --print("_gi_Slice: index called for key", k)
      return t[_giPrivateRaw][k]
    end,

    __tostring = function(t)
       --print("_gi_Slice: tostring called")
       local len = t[_giPrivateSliceProps]["len"]
       local s = "slice of length " .. tostring(len) .. " is _giSlice{"
       local r = t[_giPrivateRaw]
       -- we want to skip both the _giPrivateRaw and the len
       -- when iterating, which happens automatically if we
       -- iterate on r, the inside private data, and not on the proxy.
       for i, _ in pairs(r) do s = s .. "["..tostring(i).."]" .. "= " .. tostring(r[i]) .. ", " end
       return s .. "}"
    end,

    __len = function(t)
       -- __len does get called by the '#' operator, but IFF
       -- the XCFLAGS+= -DLUAJIT_ENABLE_LUA52COMPAT was used
       -- in the LuaJIT build. So use it!
       --
       print("len called for _gi_Slice")
       return t[_giPrivateSliceProps]["len"]
    end,

    __pairs = function(t)
       -- print("__pairs called!")
       -- this makes a _giSlice work in a for k,v in pairs() do loop.

       -- Iterator function takes the table and
       -- an index and returns the next index and
       -- associated value,
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
        if oper == "slice" then
           print("slice oper called!")
        end
    end
 }

function _gi_NewSlice(typeKind, x)
   print("_gi_NewSlice called!")
   assert(type(x) == 'table', 'bad x parameter #1: must be table')

   local arrProp = x[_giPrivateArrayProps]
   local slcProp = x[_giPrivateSliceProps]

   local raw = x
   if arrProp ~= nil then
      print("_gi_NewSlice sees x is an array")
      raw = x[_giPrivateRaw]
   elseif slcProp ~= nil then
      print("_gi_NewSlice sees x is a slice")
      raw = x[_giPrivateRaw]
   else
      print("_gi_NewSlice sees x is not an array or slice. Hmm?")      
   end
   
   -- get initial count
   local len = #x
   -- #x misses the [0] value, if present.
   if x[0] ~= nil then
      len = len + 1
   end
   --print("len is ", len)
   
   local proxy = {}
   proxy[_giPrivateRaw] = raw
   proxy["Typeof"]="_gi_Slice"

   -- this next is crashing
   --proxy[_giGo] = __lua2go(x)

   local props = {beg=0, len=len, cap=len, typeKind=typeKind}
   proxy[_giPrivateSliceProps] = props

   setmetatable(proxy, _giPrivateSliceMt)

   return proxy
end;

-- _gi_UnpackSliceRaw is a helper, used in
-- generated Lua code,
-- for calling into vararg ... Go functions.
-- This helper unpacks the raw _gi_giSlice
-- arguments. It returns non tables unchanged,
-- and non _giSlice tables unpacked.
--
function _gi_UnpackSliceRaw(t)
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

function _gi_Raw(t)
   if type(t) ~= 'table' then
      return t
   end
   
   return t[_giPrivateRaw]
end


-- append slice, as raw lua array, indexed from 1.
function append(t, ...)
   slc = {...}
   --print("append running")
   if type(t) ~= 'table' then
      return t
   end
   
   local props = t[_giPrivateSliceProps]
   if props == nil then
      error "append() called with first value not a slice"
   end
   local len = props["len"]
   local raw = t[_giPrivateRaw]
   if raw == nil then
      error "could not get raw table from slice, internal error?"
   end

   -- make copy
   local proxy = {}
   proxy["Typeof"]="_gi_Slice"

   local raw2 = {}
   if len > 0 then
      raw2[0] = raw[0]
      --print("copied raw[0] ==", raw[0])
   end
   for i,v in ipairs(raw) do
      raw2[i] = v
      --print("copied raw2[i=",i,"] ==", v)
   end
   
   --print("append at slc addition")
   local k = 0
   for i,v in ipairs(slc) do
      --print("append: i=",i," next slc element is at len+i-1=",len+i-1,"  is v=",v)
      raw2[len+i-1]=v
      k=k+1
   end
   len=len+k

   proxy[_giPrivateRaw] = raw2
   proxy[_giGo] = __lua2go(raw2)
   proxy[_giPrivateSliceProps] = {len=len, typeKind=t.typeKind}
   
   setmetatable(proxy, _giPrivateSliceMt)
   
   return proxy

end

function appendSlice(...)
   --print("appendSlice called")
   return append(...)
end

function copySlice(dest, src)
   local dlen = #dest
   local slen = #src
   local len = dlen
   if slen < len then
      len = slen
   end
   if len == 0 then
      return 0
   end
   error "not done, use wrapped Go arrays/slices"
end

function __gi_clone(a, typ)
   print("__gi_clone called: not done! TODO: finish me!")
   return a
end

function __subslice(a, beg, endx)

end

