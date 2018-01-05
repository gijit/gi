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

 _giPrivateMapMt = {

    __newindex = function(t, k, v)
      local len = t[_giPrivateMapProps]["len"]
      --print("newindex called for key", k, " len at start is ", len)
      if t[_giPrivateMapRaw][k] == nil then
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
      t[_giPrivateMapRaw][k] = v
      t[_giPrivateMapProps]["len"] = len
      --print("len at end of newidnex is ", len)
    end,

  -- __index allows us to have fields to access the count.
  --
    __index = function(t, k)
      --print("index called for key", k)
      return t[_giPrivateMapRaw][k]
    end,

    __tostring = function(t)
       local len = t[_giPrivateMapProps]["len"]
       local s = "map of length " .. tostring(len) .. " is _giMap{"
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
        local oper, key = ...
        print("oper is", oper)
        print("key is ", key)
        if oper == "delete" then
           -- the hash table delete operation
           if key == nil then
              return -- this is a no-op in Go.
           end
           -- forward the actual delete
           t[key]=nil
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

   local props = {len=len, keyType=keyType, valType=valType}
   proxy[_giPrivateMapProps] = props

   setmetatable(proxy, _giPrivateMapMt)
   return proxy
end;

