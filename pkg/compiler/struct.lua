-- structs and interfaces

__gi_InterfacePropsKey = {}
__gi_StructPropsKey = {}
__gi_MethodsetKey = {}
__gi_BaseKey = {}

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

function __structPrinter(self)
   --print("__structPrinter called")
   local s = self.__typename .." {\n"

   local uscore = 95 -- "_"
   
   for i, v in pairs(self) do
      if #i < 2 or string.byte(i,1,1)~=uscore or string.byte(i,2,2) ~= uscore then
           -- skip __ prefixed methods when printing;
           -- since most of these live in the metatable anyway.
         sv = ""
         if type(v) == "string" then
            sv = string.format("%q", v)
         else
            sv = tostring(v)
         end
         s = s .. "    "..tostring(i).. ":\t" .. sv .. ",\n"
      end
   end
   return s .. "}"
end


-- common struct behavior in this metatable
__gi_structMT = {
   __structPairs = __structPairs,
   __pairs = __structPairs,
   __name = "__gi_structMT"
}

-- common interface behavior
__gi_ifaceMT = {
   __name = "__gi_ifaceMT"
}

--
-- RegisterStruct is the first step in making a new struct.
-- It returns a methodset object.
-- Typically:
--
--   Bus   = __reg:RegisterStruct("Bus")
--   Train =  __reg:RegisterStruct("Train")
--
-- let -> point to the metatable:
--    methodset -> props -> __gi_structMT
--
function __reg:RegisterStruct(name)
   local methodset = {
      __name="structMethodSet",
      __tostring = __structPrinter
   }
   methodset.__index = methodset
   
   local props = {__typename = name, __name="structProps"}
   props[__gi_StructPropsKey] = props
   props[__gi_MethodsetKey] = methodset
   props[__gi_BaseKey] = __gi_structMT
   props.__index = props -- __gi_structMT
   
   setmetatable(props, __gi_structMT)
   setmetatable(methodset, props)
      
   self.structs[name] = methodset
   --print("debug: new methodset is: ", methodset)
   --st(methodset)
   return methodset
end

function __reg:RegisterInterface(name)
   local methodset = {__name="ifaceMethodSet"}
   methodset.__index = methodset
   
   local props = {__typename = name, __name="ifaceProps"}
   props[__gi_InterfacePropsKey] = props
   props[__gi_MethodsetKey] = methodset
   props[__gi_BaseKey] = __gi_ifaceMT
   props.__index = props

   setmetatable(props, __gi_ifaceMT)
   setmetatable(methodset, props)
   
   self.interfaces[name] = methodset
   return methodset
end



__gi_ifaceNil = __reg:RegisterInterface("nil")

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



function __reg:AddStructMethod(structName, methodName, method)
   --print("__reg:AddStructMethod called with methodName ", methodName)
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

-- __gi_count_methods
--
-- vi can be struct value or interface value;
-- We count the number of non "__" prefixed
-- methods in the metatable of vi.
--
function __gi_count_methods(vi)
   local mt = getmetatable(vi)
   if mt == nil then
      return 0
   end
   local n = 0
   local uscore = 95 -- "_"
   
   for i, v in pairs(mt) do
      -- we omit __ prefixed methods/values
      if #i < 2 or string.byte(i,1,1)~=uscore or string.byte(i,2,2) ~= uscore then
         if type(v) == "function" then
            --print("we see a function! '"..tostring(i).."'")
            n = n + 1
         end
      end
   end
   return n
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
 --[=[
   print("__gi_assertType called, typ='", typ, "' value='", value, "', returnTuple='", returnTuple, "'")
   
   local isInterface = false
   local interfaceMethods = nil
   if __reg:IsInterface(typ) then
      isInterface = true
      interfaceMethods = __reg:GetInterface(typ)
   end
   
   local ok = false
   local missingMethod = ""

   local valueMethods = getmetatable(value)

   local nvm = __gi_count_methods(valueMethods)
   
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
        
        local ni = #interfaceMethods
        
      for i, v in pairs(interfaceMethods) do
          if #i >=2 and i[1]=="_" and i[2]=="_" then
             -- skip __ prefixed methods when printing; atypical
             -- since most of these live in the metatable anyway.
             goto continue
        end

        ::continue::
     end

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
  
  if returnTupe == 0 then
     return value
  elseif returnTuple == 1 then
     return true
  end
  return value, true
 --]=]   
end
