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

       local ks = tostring(k)
       if v ~= nil then
          if t[_giPrivateMapRaw][ks] == nil then
             -- new key
             props.len = len + 1
          end
          t[_giPrivateMapRaw][ks] = v
          return

       else
          -- invar: k is not nil. v is nil.

          if t[_giPrivateMapRaw][ks] == nil then
             -- new key
             props.len = len + 1
          end
          t[_giPrivateMapRaw][ks] = _intentionalNilValue
          return
      end
      --print("len at end of newidnex is ", len)
    end,

    __index = function(t, k)
       -- Instead of __index,
       -- use __call('get', ...) for two valued return and
       --  proper zero-value return upon not present.
       -- __index only ever returns one value[1].
       -- reference: [1] http://lua-users.org/lists/lua-l/2007-07/msg00182.html
              
       --print("__index called for key", k)
       if k == nil then
          local props = t[_giPrivateMapProps]
          if props.nilKeyStored then
             return props.nilValue
          else
             -- TODO: replace nil with zero-value for the value type.
             return nil
          end
       end

       -- k is not nil.

       local ks = tostring(k)       
       local val = t[_giPrivateMapRaw][ks]
       if val == _intentionalNilValue then
          return nil
       end
       return val
    end,

    __tostring = function(t)
       --print("__tostring for _gi_Map called")
       local props = t[_giPrivateMapProps]
       local len = props["len"]
       local s = "map["..props["keyType"].__str.. "]"..props["valType"].__str.."{"
       local r = t[_giPrivateMapRaw]
       
       local vquo = ""
       if len > 0 and props.valType.__str == "string" then
          vquo = '"'
       end
       local kquo = ""
       if len > 0 and props.keyType.__str == "string" then
          kquo = '"'
       end
       
       -- we want to skip both the _giPrivateMapRaw and the len
       -- when iterating, which happens automatically if we
       -- iterate on r, the inside private data, and not on the proxy.
       for i, _ in pairs(r) do

          -- lua style:
          -- s = s .. "["..kquo..tostring(i)..kquo.."]" .. "= " .. vquo..tostring(r[i]) ..vquo.. ", "
          -- Go style
          s = s .. kquo..tostring(i)..kquo.. ": " .. vquo..tostring(r[i]) ..vquo.. ", "
       end
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
           local ks = tostring(k)
           ks, v = next(t[_giPrivateMapRaw], tostring(k))
           if v then return ks,v end
       end

       -- Return an iterator function, the table, starting point
       return stateless_iter, t, nil
    end,

    __call = function(t, ...)
        --print("__call() invoked, with ... = ", ...)
        local oper, k, zeroVal = ...
        --print("oper is", oper)
        --print("key is ", k)

        -- we use __call('get', k, zeroVal) instead of __index
        -- so that we can return multiple values
        -- to match Go's "a, ok := mymap[k]" call.
        
        if oper == "get" then

           --print("get called for key", k)
           if k == nil then
              local props = t[_giPrivateMapProps]
              if props.nilKeyStored then
                 return props.nilValue, true;
              else
                 -- key not present returns the zero value for the value.
                 return zeroVal, false;
              end
           end
           
           -- k is not nil.
           local ks = tostring(k)      
           local val = t[_giPrivateMapRaw][ks]
           if val == _intentionalNilValue then
              --print("val is the _intentinoalNilValue")
              return nil, true;

           elseif val == nil then
              -- key not present
              --print("key not present, zeroVal=", zeroVal)
              --for i,v in pairs(t[_giPrivateMapRaw]) do
              --   print("debug: i=", i, "  v=", v)
              --end
              return zeroVal, false;
           end
           
           return val, true
           
        elseif oper == "delete" then

           -- the hash table delete operation

           local props = t[_giPrivateMapProps]              
           local len = props.len
           --print("delete called for key", k, " len at start is ", len)
                      
           if k == nil then

              if props.nilKeyStored then
                 props.nilKeyStored = false
                 props.nilVaue = nil
                 props.len = len -1
              end

              --print("len at end of delete is ", props.len)              
              return
           end

           local ks = tostring(k)           
           if t[_giPrivateMapRaw][ks] == nil then
              -- key not present
              return
           end
           
           -- key present and key is not nil
           t[_giPrivateMapRaw][ks] = nil
           props.len = len - 1
           
           --print("len at end of delete is ", props.len)
        end
    end
 }
 
function _gi_NewMap(keyType, valType, x)
   assert(type(x) == 'table', 'bad parameter #3: must be table')

   local proxy = {}
   proxy["Typeof"]="_gi_Map"
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

