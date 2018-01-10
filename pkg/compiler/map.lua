-- a Lua virtual table system suitable for use in arrays and maps

-- design using the _giPrivateMapRaw index was suggested
-- by https://www.lua.org/pil/13.4.4.html
-- To intercept all writes, the requirement is that
-- the table always be empty. Hence the user uses
-- and nearly empty proxy. The only thing the proxy
-- has in it are a pointer to the actual data,
-- and a len counter.

-- create private index
_giPrivateMapRaw = {}
_giPrivateMapProps = {}

-- stored as map value in place of nil, so
-- we can recognized stored nil values in maps.
_intentionalNilValue = {}

 _giPrivateMapMt = {

    __newindex = function(t, k, v)
       --print("newindex called for key", k, " len at start is ", len)

       local props = t[_giPrivateMapProps]
       local len = props.len

       if k == nil then
          if props.nilKeyStored then
             -- replacement, no change in len.
          else
             -- new key
             props.len = len + 1
             props.nilKeyStored = true
          end
          props.nilValue = v
          return
       end

       -- invar: k is not nil

       if v ~= nil then
          if t[_giPrivateMapRaw][k] == nil then
             -- new key
             props.len = len + 1
          end
          t[_giPrivateMapRaw][k] = v
          return

       else
          -- invar: k is not nil. v is nil.

          if t[_giPrivateMapRaw][k] == nil then
             -- new key
             props.len = len + 1
          end
          t[_giPrivateMapRaw][k] = _intentionalNilValue
          return
      end
      --print("len at end of newidnex is ", len)
    end,

    __index = function(t, k)
       -- apparently only the 1st value comes back, so
       -- we return a closure with the two values that
       -- must be called to get them both out.
       
       print("index called for key", k)
       if k == nil then
          local props = t[_giPrivateMapProps]
          if props.nilKeyStored then
             return function() return props.nilValue, true; end
          else
             -- TODO: replace nil with zero-value for the value type.
             return function() return nil, false; end
          end
       end

       -- k is not nil.
       
       local val = t[_giPrivateMapRaw][k]
       if val == _intentionalNilValue then
          return function() return nil, true; end
       end
       print("index returning 2nd value of ", val ~= nil)
       return function() return val, val ~= nil; end
    end,

    __tostring = function(t)
       local props = t[_giPrivateMapProps]
       local len = props["len"]
       local s = "map["..props["keyType"].. "]"..props["valType"].." of length " .. tostring(len) .. " is _giMap{"
       local r = t[_giPrivateMapRaw]
       -- we want to skip both the _giPrivateMapRaw and the len
       -- when iterating, which happens automatically if we
       -- iterate on r, the inside private data, and not on the proxy.
       for i, _ in pairs(r) do s = s .. "["..tostring(i).."]" .. "= " .. tostring(r[i]) .. ", " end
       return s .. "}"
    end,

    __len = function(t)
       -- this does get called by the # operation(!)
       -- print("len called")
       return t[_giPrivateMapProps]["len"]
    end,

    __pairs = function(t)
       -- print("__pairs called!")
       -- this makes a _giMap work in a for k,v in pairs() do loop.

       -- Iterator function takes the table and an index and returns the next index and associated value
       -- or nil to end iteration

       local function stateless_iter(t, k)
           local v
           --  Implement your own key,value selection logic in place of next
           k, v = next(t[_giPrivateMapRaw], k)
           if v then return k,v end
       end

       -- Return an iterator function, the table, starting point
       return stateless_iter, t, nil
    end,

    __call = function(t, ...)
        print("__call() invoked, with ... = ", ...)
        local oper, k = ...
        print("oper is", oper)
        print("key is ", k)
        if oper == "delete" then

           -- the hash table delete operation

           local props = t[_giPrivateMapProps]              
           local len = props.len
           print("delete called for key", k, " len at start is ", len)
                      
           if k == nil then

              if props.nilKeyStored then
                 props.nilKeyStored = false
                 props.nilVaue = nil
                 props.len = len -1
              end

              print("len at end of delete is ", props.len)              
              return
           end

           if t[_giPrivateMapRaw][k] == nil then
              -- key not present
              return
           end
           
           -- key present and key is not nil
           t[_giPrivateMapRaw][k] = nil
           props.len = len - 1
           
           print("len at end of delete is ", props.len)
        end
    end
 }
 
function _gi_NewMap(keyType, valType, x)
   assert(type(x) == 'table', 'bad parameter #1: must be table')

   local proxy = {}
   proxy[_giPrivateMapRaw] = x

   -- get initial count
   local len = 0
   for k, v in pairs(x) do
      len = len + 1
   end

   local props = {len=len, keyType=keyType, valType=valType, nilKeyStored=false}
   proxy[_giPrivateMapProps] = props

   setmetatable(proxy, _giPrivateMapMt)
   return proxy
end;

