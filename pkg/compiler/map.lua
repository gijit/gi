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

 _giPrivateMapMt = {

    __newindex = function(t, k, v)
      print("newindex called for key", k)
      if t[_giPrivateMapRaw][k] ~= nil then
          -- replace or delete
          if v == nil then 
              t.len = t.len - 1 -- delete
          else
              -- replace, no count change              
          end
      else 
          t.len = t.len + 1
      end
      t[_giPrivateMapRaw][k] = v
    end,

  -- __index allows us to have fields to access the count.
  --
    __index = function(t, k)
      --print("index called for key", k)
      -- I don't think we need raw any more, now that
      -- __pairs works. It would be a problem for hash
      -- tables that want to store the key 'raw'.
      -- if k == 'raw' then return t[_giPrivateMapRaw] end
      return t[_giPrivateMapRaw][k]
    end,

    __tostring = function(t)
       local s = "slice of length " .. tostring(t.len) .. " is _giMap{"
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
       return t.len
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
               error("delete error: key to delete cannot be nil")
           end
           -- forward the actual delete
           t[key]=nil
        elseif oper == "slice" then

        end
    end
 }

function _gi_NewMap(x)
   assert(type(x) == 'table', 'bad parameter #1: must be table')

   -- get initial count
   local length = 0
   for k, v in pairs(x) do
      length = length + 1
   end

   local proxy = {len=length}
   proxy[_giPrivateMapRaw] = x
   setmetatable(proxy, _giPrivateMapMt)
   return proxy
end;

function _getRangeCheck(x, i)
  if x == nil or i < 0 or i >= x.len then
     error("index out of range")
  end
  return x[i]
end;

function _setRangeCheck(x, i, val)
  if x == nil or i < 0 or i >= x.len then
     error("index out of range")
  end
  x[i] = val
  return val
end;

