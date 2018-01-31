-- structs

__gi_PrivateInterfaceProps = __gi_PrivateInterfaceProps or {}

-- st or showtable, a helper.
function st(t)
   local k = 1
   for i,v in pairs(t) do
      print("num ",k, "key:",i,"val:",v)
      k=k+1
   end
end

-- don't think we're going to use these/the slice and map approach for structs.
-- _giPrivateStructRaw = {}
-- _giPrivateStructProps = {}

-- __reg is a struct registry that associates
-- names to an  __index metatable
-- that holds the methods for the structs.
--
-- reference: https://www.lua.org/pil/16.html
-- reference: https://www.lua.org/pil/16.1.html

__reg={
   -- track the registered structs here
   structs = {},
   interfaces={}
}

-- helper for iterating over structs
__structPairs = function(t)
   -- print("__pairs called!")
   -- Iterator function takes the table and an index and returns the next index and associated value
   -- or nil to end iteration
   local function stateless_iter(t, k)
      local v
      --  Implement your own key,value selection logic in place of next
      k, v = next(t, k)
      if v then return k,v end
   end
   
   -- Return an iterator function, the table, starting point
   return stateless_iter, t, nil
end

--
-- RegisterStruct is the first step in making a new struct.
-- It returns a methodset object.
-- Typically:
--
--   Bus   = __reg:RegisterStruct("Bus")
--   Train =  __reg:RegisterStruct("Train")
--
function __reg:RegisterStruct(name)
      local methodset = {}
      methodset.__tostring = __structPrinter
      methodset.__index = methodset -- is its own metatable, saving a table. (p151 / Ch 16.1 Classes, PIL 2nd ed.)
      methodset.__typename = name
      methodset.__pairs = __structPairs
      
      self.structs[name] = methodset
      --print("debug: new methodset is: ", methodset)
      --st(methodset)
      return methodset
end

function __reg:RegisterInterface(name)
   local methodset = {[__gi_PrivateInterfaceProps] = {name=name}}
   methodset.__index = methodset
   self.interfaces[name] = methodset
   return methodset
end


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
      -- methodset is the
      -- metatable for the struct.
      -- 
      -- Thus unknown direct keys like method
      -- names are forwarded
      -- to the methodset.
      setmetatable(data, methodset)
      return data
end

function __structPrinter(self)
     local s = self.__typename .." {\n"
     for i, v in pairs(self) do
        if #i >=2 and i[1]=="_" and i[2]=="_" then
           -- skip __ prefixed methods when printing; atypical
           -- since most of these live in the metatable anyway.
           goto continue
        end
        sv = ""
        if type(v) == "string" then
           sv = string.format("%q", v)
        else
           sv = tostring(v)
        end
        s = s .. "    "..tostring(i).. ":\t" .. sv .. ",\n"
        ::continue::
     end
     return s .. "}"
end


function __reg:AddMethod(structName, methodName, method)
   print("__reg:AddMethod called with methodName ", methodName)
      -- lookup the methodset
      local methodset = self.structs[structName]
      if methodset == nil then
         error("unregistered struct name '"..structName.."'")
      end
      
      -- add the method
      methodset[methodName] = method
end


function __gi_methodVal(recvr, methodName, recvrType)
   print("__gi_methodVal with methodName ", methodName, " recvrType=", recvrType)
   local methodset = __reg.structs[recvrType]
   if methodset == nil then
      error("error in __gi_methodVal: unregistered struct name '"..recvrType.."'")
   end
   local method = methodset[methodName]
   if method == nil then
      error("error in __gi_methodVal: method '"..methodName .."' not found for type '"..recvrType.."'")
   end
   return method
end

