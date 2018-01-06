-- structs

-- __reg is a struct registry that associates
-- names to an  __index metatable
-- that holds the methods for the structs.
--
-- reference: https://www.lua.org/pil/16.html
-- reference: https://www.lua.org/pil/16.1.html

__reg={
   -- track the registered structs here
   structs = {},
}

-- create a new struct instance by
-- attaching the appropriate methodset
-- to data and returning it.
function __reg:NewInstance(name, data)
   
      local methodset = self.structs[name]
      if methodset == nil then
         error("error in _struct_registry.NewInstance:"..
                  "unknown struct '"..name.."'")
      end
      -- this is the essence. The
      -- methodset acts as the
      -- metatable for the struct.
      -- Thus unknown direct keys like method
      -- names are forwarded
      -- to the methodset.
      setmetatable(data,{__index = methodset})
      return data
end
   
function __reg:RegisterStruct(name)
      local methodset = {}
      self.structs[name] = methodset
      return methodset
end

function __reg:AddMethod(structName, methodName, method)

      -- instantiate a methodset if need be
      local methodset = self.structs[structName]
      if methodset == nil then
         error("unregistered struct name '"..structName.."'")
      end
      
      -- add the method
      methodset[methodName] = method
end


-- older stuff, do we need it at all any more?

-- create private index

-- helper
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
       args = (...)
       local method, key = ...
       print("method is", method)
       -- look up the method in the vtable
       -- list of methods available
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

   local props = {keyValTypeMap=keyValTypeMap, structName=structName,methods={}}
   proxy[_giPrivateStructProps] = props

   setmetatable(proxy, _giPrivateStructMt)
   return proxy
end;
