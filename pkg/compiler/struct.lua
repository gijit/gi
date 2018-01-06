-- structs

-- create private index

_showStruct = function(props)
   r = props.structName .. " {\n"
   for k,v in pairs(props.keyValTypeMap) do
      r = r .. "   ".. k .. " ".. v .. "\n"
   end
   return r .. "}\n"
end

_giPrivateStructRaw = {}
_giPrivateStructProps = {}

 _giPrivateStructMt = {

    __newindex = function(t, k, v)
      local len = t[_giPrivateStructProps]["len"]
      --print("newindex called for key", k, " len at start is ", len)
      if t[_giPrivateStructRaw][k] == nil then
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
      t[_giPrivateStructRaw][k] = v
      t[_giPrivateStructProps]["len"] = len
      --print("len at end of newidnex is ", len)
    end,

  -- __index allows us to have fields to access the count.
  --
    __index = function(t, k)
      --print("index called for key", k)
      return t[_giPrivateStructRaw][k]
    end,
    
    __tostring = function(t)
       local props = t[_giPrivateStructProps]
       local len = props["len"]
       local s = "struct ".. _showStruct(props) .." of length " .. tostring(len) .. " is _giStruct{"
       local r = t[_giPrivateStructRaw]
       -- we want to skip both the _giPrivateStructRaw and the len
       -- when iterating, which happens automatically if we
       -- iterate on r, the inside private data, and not on the proxy.
       for i, _ in pairs(r) do s = s .. "["..tostring(i).."]" .. "= " .. tostring(r[i]) .. ", " end
       return s .. "}"
    end,

    __len = function(t)
       -- this does get called by the # operation(!)
       -- print("len called")
       return t[_giPrivateStructProps]["len"]
    end,

    __pairs = function(t)
       -- print("__pairs called!")
       -- this makes a _giStruct work in a for k,v in pairs() do loop.

       -- Iterator function takes the table and an index and returns the next index and associated value
       -- or nil to end iteration

       local function stateless_iter(t, k)
           local v
           --  Implement your own key,value selection logic in place of next
           k, v = next(t[_giPrivateStructRaw], k)
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

function _gi_NewStruct(structName, keyValTypeMap, x)
   assert(type(keyValTypeMap) == 'table', 'bad parameter #1: must be table')
   assert(type(x) == 'table', 'bad parameter #2: must be table')

   local proxy = {}
   proxy[_giPrivateStructRaw] = x

   -- get initial count
   local len = 0
   for k, v in pairs(x) do
      len = len + 1
   end

   local props = {len=len, keyValTypeMap=keyValTypeMap, structName=structName}
   proxy[_giPrivateStructProps] = props

   setmetatable(proxy, _giPrivateStructMt)
   return proxy
end;

