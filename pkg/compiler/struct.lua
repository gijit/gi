-- structs

__gi_PrivateInterfaceProps = {}
__gi_ifaceNil = {[__gi_PrivateInterfaceProps]={name="nil"}}

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
   local methodset = {}
   methodset.__index = methodset
   self.interfaces[name] = methodset
   return methodset
end

function __reg:IsInterface(name)
   return self.interfaces[name] ~= nil
end

function __reg:GetInterface(name)
   return self.interfaces[name]
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


function __reg:AddStructMethod(structName, methodName, method)
   print("__reg:AddStructMethod called with methodName ", methodName)
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

-- face.lua merged into struct.lua, because we need _reg

-- __gi_assertType is an interface type assertion.
--
--  either
--    a, ok := b.(Face)  ## the two value form (returnTupe==2)
--  or
--    a := b.(Face)      ## the one value form (returnTuple==0; can panic)
--  or
--    _, ok := b.(Face)  ## (returnTuple==1; does not panic)
--
-- returnTuple in {0, 1, 2},
--   0 returns just the interface-value, converted or nil/zero-value.
--   1 returns just the ok (2nd-value in a conversion, a bool)
--   2 returns both
--
--   if 0, then we panic when the interface conversion fails.
--
function __gi_assertType(value, typ, returnTuple)

   print("__gi_assertType called, typ='", typ, "' value='", value, "', returnTuple='", returnTuple, "'")

   local isInterface = false
   local iMethodSet = nil
   if __reg:IsInterface(typ) then
      isInterface = true
      iMethodSet = __reg:GetInterface(typ)
   end
   
   local ok = false
   local missingMethod = ""
   
  if value == __gi_ifaceNil then
     ok = false;
     
  elseif not isInterface then
     ok = value.constructor == typ;
     
  else
     local valueTypeString = value.constructor.string;
     ok = typ.implementedBy[valueTypeString];
     if ok == nil then
        
        ok = true;
        local valueMethodSet = __gi_methodSet(value.constructor);
        local interfaceMethods = typ.methods;
        local ni = #interfaceMethods
        
        for i = 1, ni do
           
           local tm = interfaceMethods[i];
           local found = false;
           local msl = #valueMethodSet
           
           for j = 1,msl do
              local vm = valueMethodSet[j];
              if vm.name == tm.name and vm.pkg == tm.pkg and vm.typ == tm.typ then
                 found = true;
                 break;
              end
           end
           
           if not found then
              ok = false;
              -- cannot cache, as repl may add/subtract methods.
              missingMethod = tm.name;
              break;
           end
        end

        -- can't cache this, repl may change it.
        --typ.implementedBy[valueTypeString] = ok;
        
     end
  end
  
  if not ok then
     
     if returnTuple == 0 then
        __gi_panic(new __gi_packages["runtime"].TypeAssertionError.ptr("", (value === __gi_ifaceNil ? "" : value.constructor.string), typ.string, missingMethod));
        
     elseif returnTuple == 1 then
        return false
     else
        return zeroVal, false
     end
  end
  
  if not isInterface then
     value = value.$val;
  end
  
  if typ == $jsObjectPtr then
     value = value.object;
  end
   --]=====]
  
  if returnTupe == 0 then
     return value
  elseif returnTuple == 1 then
     return true
  end
  return value, true
   
end
