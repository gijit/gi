-- like a Go slice, a lua slice needs to point
-- to a backing array


-- create private index
_giPrivateRaw = _giPrivateRaw or {}
_giPrivateSliceProps = _giPrivateSliceProps or {}
_giGo = _giGo or {}

_giPrivateSliceMt = {

    __newindex = function(t, k, v)
       --print("newindex called for key", k, " val=", v)
       local props = rawget(t, _giPrivateSliceProps)
       local len = props.len
       local beg = props.beg
       local raw = rawget(t, _giPrivateRaw)
       
       --print("newindex called for key", k, " with b=", b, " len at start is ", len)
       if raw[k+beg] == nil then
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
       raw[k+beg] = v
       props.len = len
       --print("len at end of newindex is ", len)
    end,

  -- __index allows us to have fields to access the count.
  --
    __index = function(t, k)
       --print("_gi_Slice: __index called for key", k)       
       local props = rawget(t, _giPrivateSliceProps)
       local raw = rawget(t, _giPrivateRaw)       
       local beg = props.beg
       local rawlen = #raw
       if raw[0] ~= nil then
          rawlen = rawlen + 1
       end
       --print("_gi_Slice: __index called for key", k, " with beg=", beg, " and rawlen=", rawlen)
       if k+beg >= rawlen then
          --print("out of bounds access " .. tostring(k+beg))
          error("out of bounds access " .. tostring(k+beg))
       end
       local res = rawget(t, _giPrivateRaw)[beg+k]
       --print("_gi_Slice __index returing res = ", res)
       return res
    end,

    __tostring = function(t)
       --print("_gi_Slice: tostring called")
       local props = rawget(t, _giPrivateSliceProps)
       local len = props.len
       local beg = props.beg
       local s = "slice <len=" .. tostring(len) .. "; beg=" .. beg .. "; cap=" .. props.cap ..  "> is _giSlice{"
       local raw = rawget(t, _giPrivateRaw)

       -- we want to skip both the _giPrivateRaw and the len
       -- when iterating, which happens automatically if we
       -- iterate on raw, the raw inside private data, and not on the proxy.
       local quo = ""
       if len > 0 and type(raw[beg]) == "string" then
          quo = '"'
       end
       for i = 0, len-1 do
          s = s .. "["..tostring(i).."]" .. "= " ..quo.. tostring(raw[beg+i]) .. quo .. ", "
       end
       
       return s .. "}"
    end,

    __len = function(t)
       -- __len does get called by the '#' operator, but IFF
       -- the XCFLAGS+= -DLUAJIT_ENABLE_LUA52COMPAT was used
       -- in the LuaJIT build. So use it!
       --
       local len = rawget(t, _giPrivateSliceProps)["len"]
       --print("len called for _gi_Slice, returning ", len)
       return len
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
           k, v = next(rawget(t, _giPrivateRaw), k)
           if v then return k,v end
       end

       -- Return an iterator function, the table, starting point
       return stateless_iter, t, nil
    end,

    __call = function(t, ...)
       --print("__call() invoked, with ... = ", ...)
       local oper, key = ...
       print("oper is", oper)
       print("key is ", key)
       if oper == "slice" then
          print("slice oper called!")
       end
    end
 }

function _gi_NewSlice(typeKind, x, zeroVal, beg, endx, cap)
   --print("_gi_NewSlice called! beg=", beg, " endx=", endx, " cap=", cap)
   assert(type(x) == 'table', 'bad x parameter #1: must be table')

   local arrProp = rawget(x, _giPrivateArrayProps)
   local slcProp = rawget(x, _giPrivateSliceProps)

   local raw = x
   local xlen = #x
   
   if arrProp ~= nil then
      --print("_gi_NewSlice sees x is an array")
      raw = rawget(x, _giPrivateRaw)
      -- xlen is correct
   elseif slcProp ~= nil then
      --print("_gi_NewSlice sees x is a slice")
      raw = rawget(x, _giPrivateRaw)
      -- xlen is correct
   else
      --print("_gi_NewSlice sees x is not an array or slice. Hmm: raw input table")
      -- #x misses the [0] value, if present.
      if x[0] ~= nil then
         xlen = xlen + 1
      end      
   end
   
   --print("_gi_NewSlice: xlen is ", xlen)
   
   local proxy = {}
   proxy[_giPrivateRaw] = raw
   proxy["Typeof"]="_gi_Slice"

   -- this next is crashing
   --proxy[_giGo] = __lua2go(x)

   beg = beg or 0
   if endx == nil then
      len = xlen - beg
      endx = beg + len
   else
      len = endx - beg 
   end
   
   --print("_gi_NewSlice: beg=", beg, " endx=",endx," len of the new slice is ", len)
   
   -- TODO: cap not really implemented in any way, just stored.
   cap = cap or len

   --print("_gi_NewSlice debug: beg=", beg, " len=", len, " endx=", endx, " cap=", cap)
   
   local props = {beg=beg, len=len, cap=cap, endx=endx, typeKind=typeKind, zeroVal=zeroVal}
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
   --print("append running, t=", tostring(t), " and slc = ", ts(slc))
   if type(t) ~= 'table' then
      return t
   end

   local props = rawget(t, _giPrivateSliceProps)

   -- simpler?
   --local tot = #t + #slc
   --local arr = _gi_NewArray({}, props.typeKind, tot, props.zeroVal)
   --local res = _gi_NewSlice(props.typeKind, arr, props.zeroVal, 0, tot, tot)
   --copy(res, t)
   
   if props == nil then
      error "append() called with first value not a slice"
   end
   local len = props.len
   local raw = rawget(t, _giPrivateRaw)
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
   --print("done with slice addition, now len=", len)
   
   proxy[_giPrivateRaw] = raw2
   --proxy[_giGo] = __lua2go(raw2)
   proxy[_giPrivateSliceProps] = {beg=0, len=len, cap=len, typeKind=props.typeKind}

   setmetatable(proxy, _giPrivateSliceMt)

   --print("append returning new slice = ", tostring(proxy))
   return proxy

end

function appendSlice(...)
   --print("appendSlice called")
   return append(...)
end

function __copySlice(dest, src)
   --print("__copySlice called")
   local propsDest = rawget(dest, _giPrivateSliceProps)
   local propsSrc  = rawget(src,  _giPrivateSliceProps)
   if propsDest == nil then
      error "__copySlice() called with destination value not a slice"
   end
   if propsSrc == nil then
      error "__copySlice() called with source value not a slice"
   end

   local begDest = propsDest.beg
   local begSrc  = propsSrc.beg
   local rawDest = rawget(dest, _giPrivateRaw)
   local rawSrc  = rawget(src,  _giPrivateRaw)

   local dlen = #dest
   local slen = #src
   local len = dlen
   if slen < len then
      len = slen
   end
   if len == 0 then
      return 0LL
   end
   
   -- copy direction allows for overlap   
   if begSrc > begDest then
      --print("src.beg > dest.beg, copying forward, step=+1")
      for i = 0, len-1 do
         rawset(rawDest, i+begDest, rawget(rawSrc, i+begSrc))
      end
   else
      --print("src.beg <= dest.beg, copying backward, step=-1")
      for i = len-1, 0, -1 do
         rawset(rawDest, i+begDest,  rawget(rawSrc, i+begSrc))
      end
   end
   --print("done with __copySlice, returning len=", len)
   return int(len)
end

function __gi_clone(a, typ)
   error("__gi_clone called: not done! TODO: finish me!")
   return a
end

function __subslice(a, beg, endx)
   print("top of __subslice, beg=",beg, " endx=", endx)
   
   local arrProp = rawget(a, _giPrivateArrayProps)
   local slcProp = rawget(a, _giPrivateSliceProps)

   local raw = a
   if arrProp ~= nil then
      --print("__subslice sees x is an array")
      raw = rawget(a, _giPrivateRaw)
   elseif slcProp ~= nil then
      --print("__subslice sees x is a slice")
      raw = rawget(a, _giPrivateRaw)
   else
      --print("__subslice sees x is not an array or slice. Hmm?")
      error("must have slice or array in __subslice")
   end

   local props = slcProp or arrProp
   local typeKind = props.typeKind

   return _gi_NewSlice(typeKind, raw, props.zeroVal, beg, endx)
end

function __gi_makeSlice(typeKind, zeroVal, len, cap)
   print("__gi_makeSlice() called, typeKind=", typeKind, " zeroVal= ",zeroVal," len=", len, " cap=", cap)
   raw = {}
   cap = cap or len
   for i = 0, cap-1 do
      raw[i] = zeroVal
   end
   return _gi_NewSlice(typeKind, raw, zeroVal, 0, len)
end
