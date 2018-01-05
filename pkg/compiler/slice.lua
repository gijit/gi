-- like a Go slice, a lua slice needs to point
-- to a backing array


-- create private index
_giPrivateSliceRaw = {}

 _giPrivateSliceMt = {

    __newindex = function(t, k, v)
      print("newindex called for key", k)
      if t[_giPrivateSliceRaw][k] ~= nil then
          -- replace or delete
          if v == nil then 
              t.len = t.len - 1 -- delete
          else
              -- replace, no count change              
          end
      else 
          t.len = t.len + 1
      end
      t[_giPrivateSliceRaw][k] = v
    end,

  -- __index allows us to have fields to access the count.
  --
    __index = function(t, k)
      --print("index called for key", k)
      -- I don't think we need raw any more, now that
      -- __pairs works. It would be a problem for hash
      -- tables that want to store the key 'raw'.
      -- if k == 'raw' then return t[_giPrivateSliceRaw] end
      return t[_giPrivateSliceRaw][k]
    end,

    __tostring = function(t)
       local s = "slice of length " .. tostring(t.len) .. " is _giSlice{"
       local r = t[_giPrivateSliceRaw]
       -- we want to skip both the _giPrivateSliceRaw and the len
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
       -- this makes a _giSlice work in a for k,v in pairs() do loop.

       -- Iterator function takes the table and an index and returns the next index and associated value
       -- or nil to end iteration

       local function stateless_iter(t, k)
           local v
           --  Implement your own key,value selection logic in place of next
           k, v = next(t[_giPrivateSliceRaw], k)
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

function _gi_NewSlice(x)
   assert(type(x) == 'table', 'bad parameter #1: must be table')

   -- get initial count
   local length = 0
   for k, v in pairs(x) do
      length = length + 1
   end

   local proxy = {len=length}
   proxy[_giPrivateSliceRaw] = x
   setmetatable(proxy, _giPrivateSliceMt)
   return proxy
end;
