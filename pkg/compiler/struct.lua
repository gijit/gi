-- structs and interfaces

__debug = false

-- general note:
-- the convention in translating gopherjs javascript's '$'
-- is to replace the '$' prefix with "__gi_"

-- port gopherjs system $clone()
--
function __NOTINUSE__gi_clone(src, typ, newType)
   print("__gi_clone() called with typ="..tostring(typ), " and newType="..tostring(newType))
   local clone = typ.__zero();
   typ.__copy(clone, src);
   return clone;
end

-- TODO: syncrhonize around this/deal with multi-threading?
--  may need to teach LuaJIT how to grab go mutexes or use sync.Atomics.
__gi_idCounter = 0;

__gi_PropsKey = {}
__gi_MethodsetKey = {}
__jsObjectPtr = {}

-- keep types and values separate; keep
-- packages distinct.
__curpkg = {
   path = "main",
   name = "main",
   types = {},
   vars  = {}
}

-- st or showtable, a debug print helper.
function __st(t, name, indent, quiet, methods_desc)
   if t == nil then
      local s = "<nil>"
      if not quiet then
         print(s)
      end
      return s
   end
   
   if type(t) ~= "table" then
      local s = tostring(t)
      if not quiet then
         if type(t) == "string" then
            print('"'..s..'"')
         else 
            print(s)
         end
      end
      return s
   end
   
   
   local k = 0
   local name = name or ""
   local namec = name
   if name ~= "" then
      namec = namec .. ": "
   end
   local indent = indent or 0
   local pre = string.rep(" ", 4*indent)..namec
   local s = pre .. "============================ "..tostring(t).."\n"
   for i,v in pairs(t) do
      k=k+1
      local vals = ""
      if methods_desc then
         --print("methods_desc is true")
         vals = __st(v,"",indent+1,quiet,methods_desc)
      else 
         vals = tostring(v)
      end
      s = s..pre.." "..tostring(k).. " key: '"..tostring(i).."' val: '"..vals.."'\n"
   end
   if k == 0 then
      s = pre.."<empty table>"
   end

   local mt = getmetatable(t)
   if mt ~= nil then
      s = s .. "\n"..__st(mt, "mt.of."..name, indent+1, true)
   end
   if not quiet then
      print(s)
   end
   return s
end

-- debug helper for .__methods_desc
function __show_methods_desc(t, name, indent, quiet)
   return __st(t,name,indent,quiet,true)
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

-- jea: can we delete this, or no?

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

   -- get address, avoiding infinite loop of self-calls.
   local mt = getmetatable(self)
   setmetatable(self, nil)
   local addr = tostring(self) 
   setmetatable(self, mt)

   local s = "self.__typename: "..self.__typename .."; "..addr.." {\n"   
   
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


function __ifacePrinter(self)
   --print("__ifacePrinter called")
   local s = "type " .. self.__typename .." interface{\n"

   local methset = self[__gi_MethodsetKey]
   
   local uscore = 95 -- "_"
   
   for i, v in pairs(methset) do
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


-- delete and fold into props:: __gi_structMT


-- common interface behavior
__gi_ifaceMT = {
   __name = "__gi_ifaceMT"
}

--
-- RegisterStruct is the first step in making a new struct.
-- It returns a methodset object.
-- Typically:
--
--   Bus   = __reg:RegisterStruct("Bus","main","main")
--   Train =  __reg:RegisterStruct("Train","main","main")
--
-- if we denote metatable with the
--  arrow from table -> metatable, then
--
--  instance ->  methodset -> props
--
function __reg:RegisterStruct(shortTypeName, pkgPath, shortPkg)
   local name = shortTypeName -- temporary fix
   print("RegisterStruct called, with shortTypeName='"..shortTypeName.."'")
   if shortTypeName == nil then
      error "error in __reg:RegisterStruct: shortTypeName cannot be nil"
   end
   
   local methodset = {
      __name="structMethodSet",

      -- make __tostring as local as possible,
      -- to avoid the infinite looping we got
      -- when it was higher up.

      -- essential for pretty-printing a struct
      __tostring = __structPrinter
   }
   methodset.__index = methodset
   
   local props = {__typename = name, __name="structProps", __nMethod=0}
   props[__gi_PropsKey] = props
   props[__gi_MethodsetKey] = methodset
   props.__index = props
   -- temp debug, but do we really need these?:
   --props.__structPairs = __structPairs
   --props.__pairs = __structPairs
   
   setmetatable(methodset, props)
   
   self.structs[name] = methodset
   print("__reg:RegisterStruct done, debug: new methodset is: ", methodset)
   __st(methodset, shortTypeName..".methodset")
   return methodset
end

function __reg:RegisterInterface(shortTypeName, pkgPath, shortPkg)
   local name = pkgPath.."."..shortTypeName
   if name == nil then
      error "error in RegisterInterface: name  cannot be nil"
   end
   --print("weee: __reg:RegisterInterface called with name='"..name.."'")
   --print(debug.traceback())
   
   local methodset = {
      __name="interfaceMethodSet",
      __tostring = __ifacePrinter
   }
   methodset.__index = methodset
   
   local props = {__typename = name, __name="interfaceProps"}
   props[__gi_PropsKey] = props
   props[__gi_MethodsetKey] = methodset
   props.__index = props

   setmetatable(methodset, props)
   -- jea: not sure we want or need this any more:
   --setmetatable(props, __gi_ifaceMT)
   
   self.interfaces[name] = methodset
   return methodset
end


__gi_ifaceNil = __reg:RegisterInterface("nil","main","main")

function __reg:IsInterface(typ)
   local name = typ.__str
   return self.interfaces[name] ~= nil
end

function __reg:GetInterfaceMethods(typ)
   local name = typ.__str
   return self.interfaces[name]
end

function __reg:GetPointeeMethodset(shortTypeName, pkgPath, shortPkg)
   local goal = string.sub(shortTypeName, 2) -- remove leading dot.
   print("top of __reg:GetPointeeMethodset, goal='"..goal.."' are here are structs:")
   __st(self.structs, "__reg.structs")
   
   local strct = self.structs[goal]
   if strct ~= nil then
      print("__reg:GetPointeeMethodset: found in structs")
      return strct
   end
   
   local face = self.interfaces[goal]
   if face ~=nil then
      print("__reg:GetPointeeMethodset: found in interfaces")
      return face
   end

   print("__reg:GetPointeeMethodset: '"..goal .."' *not* found in structs or interfaces")   
   
   -- other types? well, they
   -- won't have methodsets, so nil seems appropriate.
   return nil
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


-- si should be "struct" or "iface"
-- siName is the name of the struct or interface.
--
function __reg:RemoveMethod(si, siName, methodName)
   -- known?
   local methodset
   if si == "struct" then
      methodset = self.structs[siName]
   else
      methodset = self.interfaces[siName]
   end

   if methodset == nil then
      error("unregistered "..si.." name '"..siName.."'")
   end

   if methodset[methodName] == nil then
      -- not known, don't adjust nMethod count
      error("error in RemoveMethod: '"..methodName.."' not found in "..si.." '"..siName.."'")
      return
   end
   -- delete from methoset, decrease nMethod count.
   methodset[methodName] = nil
   local props = methodset[__gi_PropsKey]
   props.__nMethod = props.__nMethod -1
end

-- AddMethod works for both structs and interfaces.
--
-- si should be "struct" or "iface", to say which type.
-- siName is the name of the struct or interface.
--
-- incCount is true if we need to increment the method
-- count to account for this just directly added method.
--
function __reg:AddMethod(si, siName, methodName, method, incCount)
   print("__reg:AddMethod for '"..si.."' called with methodName ", methodName)
   -- lookup the methodset
   local methodset
   if si == "struct" then
      methodset = self.structs[siName]
   else
      methodset = self.interfaces[siName]
   end
   if methodset == nil then
      error("unregistered "..si.." name '"..siName.."'")
   end

   print("prior to addition, methodset is:")
   __st(methodset, "methodset")
   
   -- new?
   if methodset[methodName] == nil or incCount then
      -- new, count it.
      local props = methodset[__gi_PropsKey]
      props.__nMethod = props.__nMethod + 1
   else      
      -- not new
      print("methodName "..methodName.." was not new, val is:", methodset[methodName])
   end
   
   -- add the method
   methodset[methodName] = method

   print("after addition, methodset is:")
   __st(methodset, "methodset")   
end

function __gi_methodVal(recvr, methodName, recvrType)
   print("__gi_methodVal with methodName ", methodName, " recvrType=", recvrType)

   -- try structs, then interfaces.
   
   local methodset = __reg.structs[recvrType]
   if methodset == nil then
      methodset = __reg.interfaces[recvrType]
   end
   
   if methodset == nil then
      error("error in __gi_methodVal: unregistered receiver type '"..recvrType.."'")
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
-- methods in the methodset of vi.
--
function __gi_count_methods(vi)
   local mset = vi[ __gi_MethodsetKey]
   if mset == nil then
      return 0
   end
   local n = 0
   local uscore = 95 -- "_"
   
   for i, v in pairs(mset) do
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

-- face.lua merged into struct.lua, because we need _reg.
-- Thus the sequencing of these declarations is significant.

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

   print("__gi_assertType called, typ='", typ, "' value='", value, "', returnTuple='", returnTuple, "'. full value __st dump:")
   __st(value, "value")
   print("\n\n and typ is: ", type(typ))
   __st(typ, "typ")
   
   local isInterface = false
   local interfaceMethods = nil
   if __reg:IsInterface(typ) then
      print("__gi_assertType notes that typ is interface")
      isInterface = true
      
      --interfaceMethods = __reg:GetInterfaceMethods(typ)
      interfaceMethods = typ.__methods_desc
      
      if interfaceMethods == nil then
         print("interfaceMethods for typ was nil!?!")
      else
         __show_methods_desc(interfaceMethods, "interfaceMethods from typ.__methods_desc")
      end
   else
      print("__gi_assertType notes that typ is NOT an interface")
   end
   
   local ok = false
   local missingMethod = ""
   
   local valueMethods = value[__gi_MethodsetKey]
   local valueProps = value[__gi_PropsKey]
   
   -- not used any more
   --local nMethod = valueProps.__nMethod
   --local nvm = __gi_count_methods(valueMethods)
   
   if value == __gi_ifaceNil then
      ok = false;
      
   elseif not isInterface then
      ok = value.__constructor == typ;
      
   else
      local valueTypeString = value.__str
      
      print("__gi_assertType: valueTypeString='"..valueTypeString.."' and typ is: ")
      __st(typ)

      -- Unfortunately, we can't use the __implementedBy
      -- cache at the REPL.  The user might subtract
      -- or add methods to either type at any time.
      --
      -- Moreover, even if we wanted to, the updated type
      -- has no idea what other types to tell to update,
      -- and we certainly don't want to broadcast the
      -- update to all current types.
      --
      -- So for correctness, we can't use __implementedBy.
      --
      -- ok = typ.__implementedBy[valueTypeString];
      
      --if not ok then
         print("assertType: ")
         
         ok = true;

         local  valueMethodSet = value.__methods_desc
         
         --local  valueMethodSet = value[__gi_MethodsetKey]
         print("valueMethodSet is")
         __show_methods_desc(valueMethodSet, "valueMethodSet")
         
         --local valueMethodSet = __gi_methodSet(value.__str);

         local msl = __gi_count_methods(valueMethodSet)
         
         local ni = #interfaceMethods
         local uscore = 95 -- "_"

         print("ni = ", ni)
         
         for i = 1, ni do
            print("i = ",i," out of ni = ", ni)
            
            local tm = interfaceMethods[i];
            print("tm =")
            __st(tm)

            --if #tm >= 2 and string.byte(tm,1,1)==uscore and string.byte(tm,2,2) == uscore then
            -- print("skipping '__' prefixed method: "..tostring(tm))
            --   goto continue
            --end
            
            local found = false;

            print("i = ", i)
            __st(valueMethodSet, "valueMethodSet")
            
            for j, vm in pairs(valueMethodSet) do

               print("on j =", j, " vm =")
               __st(vm, "vm")
               
               --if type(j) == "string" and #j >= 2 and
               --   string.byte(j,1,1)==uscore and
               --string.byte(j,2,2) == uscore then
               --   
               --   print("skipping '__' prefixed method: "..tostring(j))
               --   goto continue2
               --end
               
               if vm.__name == tm.__name and vm.__pkg == tm.__pkg and vm.__typ == tm.__typ then
               -- temp debug: just match on the name until we
               -- figure out where the vm typ info lives.
               --if j == tm.__name then 
                  print("found 0000000000 method match, vm=")
                  __st(vm, "vm")
                  found = true;
                  break;
               else
                  -- debug prints:
                  print("not match, 111111111: tried to compare: vm=")
                  __st(vm, "vm")
                  print("not match, 111111111: tried to compare: tm=")
                  __st(tm, "tm")
               end
               
               --::continue2::
            end
            
            if not found then
               ok = false;
               -- cannot cache, as repl may add/subtract methods.
               missingMethod = tm.__name;
               print("7777777777 jea debug: missingMethod is '"..missingMethod.."'")
               break;
            end
            --::continue::            
         end
      
         -- but, note we can't cache this, repl may change it.
         -- typ.__implementedBy[valueTypeString] = ok;
      
      --end
   end

   if not ok then
      
      if returnTuple == 0 then
         
         local ctor
         if value == __gi_ifaceNil then 
            ctor = ""
         else
            ctor = value.str
         end
         error("runtime.TypeAssertionError."..typ.str.." is missing '"..missingMethod.."'")
         -- __gi_panic(new __gi_packages["runtime"].TypeAssertionError.__ptr("", ctor, typ.__str, missingMethod)
         
      elseif returnTuple == 1 then
         return false
      else
         return zeroVal, false
      end
   end
   
   if not isInterface then
      -- jea: value = value.__gi_val;
      -- jea: I think value should just be value, why not? no?
   end
   
   if typ == __gi_jsObjectPtr then
      value = value.object;
   end
   
   if returnTuple == 0 then
      return value
   elseif returnTuple == 1 then
      return true
   end
   return value, true

end


-- support for __gi_NewType

__gi_kind_bool = 1;
__gi_kind_int = 2;
__gi_kind_int8 = 3;
__gi_kind_int16 = 4;
__gi_kind_int32 = 5;
__gi_kind_int64 = 6;
__gi_kind_uint = 7;
__gi_kind_uint8 = 8;
__gi_kind_uint16 = 9;
__gi_kind_uint32 = 10;
__gi_kind_uint64 = 11;
__gi_kind_uintptr = 12;
__gi_kind_float32 = 13;
__gi_kind_float64 = 14;
__gi_kind_complex64 = 15;
__gi_kind_complex128 = 16;
__gi_kind_Array = 17;
__gi_kind_Chan = 18;
__gi_kind_Func = 19;
__gi_kind_Interface = 20;
__gi_kind_Map = 21;
__gi_kind_Ptr = 22;
__gi_kind_Slice = 23;
__gi_kind_String = 24;
__gi_kind_Struct = 25;
__gi_kind_UnsafePointer = 26;

__kind2str = {
[1]="__gi_kind_bool",
[2]="__gi_kind_int",
[3]="__gi_kind_int8",
[4]="__gi_kind_int16",
[5]="__gi_kind_int32",
[6]="__gi_kind_int64",
[7]="__gi_kind_uint",
[8]="__gi_kind_uint8",
[9]="__gi_kind_uint16",
[10]="__gi_kind_uint32",
[11]="__gi_kind_uint64",
[12]="__gi_kind_uintptr",
[13]="__gi_kind_float32",
[14]="__gi_kind_float64",
[15]="__gi_kind_complex64",
[16]="__gi_kind_complex128",
[17]="__gi_kind_Array",
[18]="__gi_kind_Chan",
[19]="__gi_kind_Func",
[20]="__gi_kind_Interface",
[21]="__gi_kind_Map",
[22]="__gi_kind_Ptr",
[23]="__gi_kind_Slice",
[24]="__gi_kind_String",
[25]="__gi_kind_Struct",
[26]="__gi_kind_UnsafePointer"
}

__gi_methodSynthesizers = {}
__gi_addMethodSynthesizer = function(f) 
   if __gi_methodSynthesizers == nil then
      f();
      return;
   end
   __gi_methodSynthesizers[#__gi_methodSynthesizers+1] = f;
end

__gi_synthesizeMethods = function()
   for i,f in pairs(__gi_methodSynthesizers) do
      f()
   end
   __gi_methodSynthesizers = nil;
end

__gi_ifaceKeyFor = function(x)
   if x == __gi_ifaceNil then
      return "nil"
   end
   local c = x.__constructor
   --return c.string .. "__gi_" .. c.__keyFor(x.__gi_val)
   return c.string .. "__gi_" .. c.__keyFor(x)
end

__gi_identity = function(x) return x; end

__gi_typeIDCounter = 0;

__gi_idKey = function(x) 
   if x.__id == nil then
      __gi_idCounter = __gi_idCounter + 1
      x.__id = __gi_idCounter;
   end
   return tostring(x.__id);
end

--__castableMT = {
--   __name = "__castableMT",
--    __call = function(t, ...)
--       print("__castableMT __call() invoked, with ... = ", ...)
--       local arg0 = ...
--       print("in __castableMT, arg0 is", arg0)
--       t.__gi_val = arg0
--    end
-- }

__gi_identity = function(x) return x; end

__gi_floatKey = function(f)
   if f ~= f then
      __gi_idCounter = __gi_idCounter+1
      return "NaN_" + __gi_idCounter;
   end
   return tostring(f);
end


-- metatable for __gi_NewType types
__gi_type_MT = {
   __name = "__gi_type_MT",
   __call = function(self, ...)
      local args = {...}
      print("jea debug: __git_type_MT.__call() invoked, self='",tostring(self),"', with args=")
      __st(args, "__gi_type_MT.args")
   end
}

-- ugh. too much javascript magic. avoid this in favor
-- of being more explicit.
--
-- for porting gopherjs' 'new Ctor' code where Ctor is a constructor
-- that must take 'self' as its first argument:
-- function __gi_new(ctor, ...)
--   local self = {}
--   ctor(self, ...)
--   return self
-- end


__gi_NewType_constructor_MT = {
   __name = "__gi_NewType_constructor_MT",
   __call = function(the_mt, self, ...)
      print("jea debug: __git_NewType_constructor_MT.__call() invoked, self='",tostring(self),"', with args=")
      print("start st")
      __st({...}, "__gi_NewType_constructor_MT.dots")
      print("end st")
      if self ~= nil and self.__constructor ~= nil then
         print("calling self.__constructor!")
         return self.__constructor(self, ...)
      end
      return self
   end
}


-- create new type. 
-- translate __gi_newType() in js,
-- from gopherjs/compiler/prelude/types.go#L64
--
-- sio \in {"struct", "iface", "other"}
--
function __gi_NewType(size, kind, shortPkg, shortTypeName, str, named, pkgPath, exported, constructor)

   if __debug then
      print("=====================")
      print("top of __gi_NewType()")
      print("=====================")
      
      print("size='"..tostring(size).."'")
      print("kind='"..tostring(kind).."'")
      print("kind2str='".. __kind2str[kind].."'")
      print("str='"..str.."'")
      print("shortTypeName='"..shortTypeName.."'")
      print("named='"..tostring(named).."'")
      print("shortPkg='".. shortPkg.."'")
      print("pkgPath='"..pkgPath.."'")
      print("exported='"..tostring(exported).."'")
      print("constructor='"..tostring(constructor).."'")
      print("elemTyp='"..tostring(elemTyp).."'")
   end
   
   -- we return typ at the end.
   local typ = {}

   if kind == __gi_kind_Struct then
      
      typ.__registered = __reg:RegisterStruct(shortTypeName, pkgPath, shortPkg)
      -- replace typ with the props for a struct
      typ = typ.__registered[__gi_PropsKey]
      
   elseif kind == __gi_kind_Interface then

      typ.__registered = __reg:RegisterInterface(shortTypeName, pkgPath, shortPkg)
      -- replace typ with the props for the interface
      typ = typ.__registered[__gi_PropsKey]
      
   elseif kind == __gi_kind_Ptr then
      
      typ.__registered = __reg:GetPointeeMethodset(shortTypeName, pkgPath, shortPkg)
      print("typ.__registered back from __reg:GetPointeeMethodset = ", typ.__registered)

   else
      setmetatable(typ, __gi_NewType_constructor_MT)
      --setmetatable(typ, __gi_type_MT) -- make it callable
   end
   
   if kind == __gi_kind_bool or
      kind == __gi_kind_int or
      kind == __gi_kind_int8 or
      kind == __gi_kind_int16 or
      kind == __gi_kind_int32 or
      kind == __gi_kind_int64 or
      kind == __gi_kind_uint or
      kind == __gi_kind_uint8 or
      kind == __gi_kind_uint16 or
      kind == __gi_kind_uint32 or
      kind == __gi_kind_uint64 or
      kind == __gi_kind_uintptr or
   kind == __gi_kind_UnsafePointer then
      
      typ.__constructor= function(self, v)
         --self.__gi_val = v;
      end
      --typ.__gi_val = 0LL
      typ.__wrapped = true;
      typ.__keyFor = __gi_identity;
      
      
   elseif kind == __gi_kind_String then

      typ.__constructor = function(self, v)
         --self.__gi_val = v;
      end
      typ.__wrapped = true;
      typ.__keyFor = function(x) return "__gi_"..x; end

   elseif kind ==  __gi_kind_float32 or
   kind == __gi_kind_float64 then

      --typ.__gi_val = 0
      typ.__constructor = function(self, v)
         --self.__gi_val = v;
      end
      typ.__wrapped = true;      
      typ.__keyFor = function(x) return __gi_floatKey(x) end
      
   elseif kind == __gi_kind_complex64 then
      typ.__wrapped = true;      
      typ.__constructor = function(self, real, imag)
         self.__gi_real = __gi_fround(real);
         self.__gi_imag = __gi_fround(imag);
         --self.__gi_val = self;
      end
      typ.__keyFor = function(x)  return x.__gi_real .. "__gi_" .. x.__gi_imag; end
      

   elseif kind == __gi_kind_complex128 then
      typ.__wrapped = true;      
      typ.__constructor = function(self, real, imag)      
         self.__gi_real = real;
         self.__gi_imag = imag;
         --self.__gi_val = self;
      end
      typ.__keyFor = function(x)  return x.__gi_real .. "_" .. x.__gi_imag; end

   elseif kind == __gi_kind_Array then
      
      typ.__constructor = function(v)
         return __gi_NewArray(v, typ.__elem, #v, typ.__elem.__zero())
      end      
      typ.__wrapped = true;
      typ.__ptr = __gi_NewType(8, __gi_kind_Ptr, shortPkg, "*"..shortTypeName, "*" .. str, false, "", false, function(self, array) 
                                self.__gi_get = function() return array; end;
                                self.__gi_set = function(v) typ.__copy(self, v); end
                                --self.__gi_val = array;
      end);
      typ.__init = function(elem, len) 
         typ.__elem = elem;
         typ.__len = len;
         typ.__comparable = elem.__comparable;
         typ.__keyFor = function(x)
            local ma = __gi_mapArray(x, function(e)
                                        return tostring(elem.__keyFor(e))
            end)
            return table.concat(ma, "_")
         end
         typ.__copy = function(dst, src) 
            __gi_copyArray(dst, src, 0, 0, src.length, elem);
         end
         typ.__ptr.__init(typ);
         --jea: what to do with this? define a __call somewhere?
         --jea: Object.defineProperty(typ.__ptr.nil, "nilCheck", { get: __gi_throwNilPointerError });
      end
      

   elseif kind == __gi_kind_Chan then
      typ = function(self, v)
         --self.__gi_val = v;
      end
      typ.__wrapped = true;
      typ.__keyFor = __gi_idKey;
      typ.__init = function(elem, sendOnly, recvOnly)
         typ.__elem = elem;
         typ.__sendOnly = sendOnly;
         typ.__recvOnly = recvOnly;
      end
      

   elseif kind == __gi_kind_Func then
      typ.__constructor = function(v)
         local newb = {__gi_val = v}
         setmetatable(newb, {
                         __call = function(the_mt, me, ...) return me.__gi_val(...) end
         })
         return newb
      end
      typ.__wrapped = true;
      typ.__init = function(params, results, variadic)
         typ.__params = params;
         typ.__results = results;
         typ.__variadic = variadic;
         typ.__comparable = false;
      end
      

   elseif kind == __gi_kind_Interface then
      --typ.__implementedBy= {}
      typ.__missingMethodFor= {}
      typ.__keyFor = __gi_ifaceKeyFor;
      typ.__init = function(methods)
         print("in __init function for interface, is typ == self? -> "..tostring((typ == self)))
         typ.__methods_desc = methods;
         __st(methods, "methods")
         for i,m in pairs(methods) do
            __gi_ifaceNil[m.__prop] = __gi_throwNilPointerError;
         end
      end
      
   elseif kind == __gi_kind_Map then
      typ.__constructor = function(v)
         local newb = {__gi_val = v}
         return newb
      end
      typ.__wrapped = true;
      typ.__init = function(key, elem)
         typ.__key = key;
         typ.__elem = elem;
         typ.__comparable = false;
      end
      
   elseif kind == __gi_kind_Slice then
      typ.__constructor = function(self, array)
         self.__gi_array = array;
         self.__gi_offset = 0;
         self.__gi_length = #array
         self.__gi_capacity = self.__gi_length
         --self.__gi_val = self;
      end
      typ.__init = function(elem)
         typ.__elem = elem;
         typ.__comparable = false;
         typ.__nativeArray = __gi_nativeArray(elem.__kind);
         typ.__nil = typ({});
      end

   --------------------------------------------
   --------------------------------------------
   --------------------------------------------
      
   elseif kind == __gi_kind_Ptr then
      print("jea debug: at kind == __gi_kind_Ptr in __gi_NewType()")

      local mt = {
         __name = "Ptr type constructed mt",
         __call = function(the_mt, self, ...)
            print("jea debug: per-ptr-type ctor_mt.__call() invoked, the_mt='"..tostring(the_mt).."', self='"..tostring(self).."', with args=")
            print("start st")
            __st({...}, "Ptr.mt.__call.per-ptr-type-ctor.args")
            print("end st")

            -- typ captured by closure.
            if typ ~= nil and typ.__constructor ~= nil then
               print("calling ptr self.__constructor!")
               local newb = typ.__constructor(self, ...)
               print("done calling ptr typ.__constructor!")               
               if typ.__registered ~= nil then
                  print("after ptr self.ctor, setting typ.__registered as metatable.")
                  setmetatable(newb, typ.__registered)
               else
                  print("after ptr self.ctor, setting typ.__registered was nil")
                  setmetatable(newb, __gi_NewType_constructor_MT)
               end
               return newb
            end
            --is this why our .__ptr changing after instantiation? nope.
            setmetatable(self, typ.__registered)
            return self
         end
      }
      setmetatable(typ, mt)
      print("setting Ptr typ.__constructor to constructor: "..tostring(constructor))
      -- typ.__constructor = constructor
      
      typ.__constructor = constructor or
         function(getter, setter, target)
            print("jea debug: top of kind_Ptr constructor")
            local newb = {__gi_get = getter, __gi_set = setter}
            newb.__gi_val = newb
            return newb
         end
      
      typ.__keyFor = __gi_idKey;
      typ.__init = function(elem) 
         typ.__elem = elem;
         typ.__wrapped = (elem.__kind == __gi_kind_Array);-- key insight: what __wrapped means!
         typ.__nil = __gi_ptrType(__gi_throwNilPointerError, __gi_throwNilPointerError, "nil");
      end

   --------------------------------------------
   --------------------------------------------
   --------------------------------------------
      
   elseif kind == __gi_kind_Struct then
      print("jea debug: at kind == __gi_kind_Struct in __gi_NewType()")

      local mt = {
         __name = "Struct type constructed mt",
         __call = function(the_mt, self, ...)
            print("jea debug: per-struct-type ctor_mt.__call() invoked, self='",tostring(self),"', with args=")
            print("start st")
            __st({...},"Struct.mt.__call.dots")
            print("end st")
            if self ~= nil and self.__constructor ~= nil then
               print("calling self.__constructor!")
               local newb = self.__constructor(self, ...)
               if typ.__registered ~= nil then
                  setmetatable(newb, typ.__registered)
               else
                  setmetatable(newb, __gi_NewType_constructor_MT)                  
               end
               return newb
            end
            setmetatable(self, typ.__registered)
            return self
         end
      }
      setmetatable(typ, mt)
      typ.__constructor = constructor
      
      typ.__wrapped = true;
      
      -- NB, we are currently in kind == __gi_kind_Struct
      -- the typ.__ptr gets built before typ itself is finished.
      
      typ.__ptr = __gi_NewType(8, __gi_kind_Ptr, shortPkg, "*"..shortTypeName, "*" .. str, false, pkgPath, exported, constructor);
      typ.__ptr.__elem = typ;
      typ.__ptr.prototype = {}
      typ.__ptr.prototype.__gi_get = function()  return this; end
      typ.__ptr.prototype.__gi_set = function(v) typ.__copy(this, v); end
      typ.__init = function(pkgPath, fields)
         typ.__pkg = pkgPath;
         typ.__fields = fields;
         for i,f in pairs(fields) do

            print("jea debug, f =")
            print("jea debug, type(f.__typ) =", type(f.__typ))
            __st(f, "f in __init() for struct")
            
            if type(f.__typ) == "cdata" then
               -- cdata should be comparable.
            elseif f ~= nil and f.__typ ~= nil and not f.__typ.__comparable then
               typ.__comparable = false;
            end
         end
         typ.__keyFor = function(x) 
            return __gi_mapArray(fields, function(f)
                                    return tostring(f.__typ.__keyFor(val[f.__prop])) end)
         end
         typ.__copy = function(dst, src) 
            for i = 0,fields.length-1 do
               local f = fields[i];
               local knd = f.__typ.__kind
               if knd ==  __gi_kind_Array then
                  -- do nothing
               elseif knd == __gi_kind_Struct then
                  f.__typ.__copy(dst[f.__prop], src[f.__prop]);
               else
                  -- default:
                  dst[f.__prop] = src[f.__prop];
               end
            end
         end
         -- nil value
         local properties = {};
         for i,f in pairs(fields) do
               properties[f.__prop] = { get= __gi_throwNilPointerError, set= __gi_throwNilPointerError }
         end
         typ.__ptr.__nil = {} -- jea what here? Object.create(constructor.prototype, properties);
         typ.__ptr.__nil.__gi_val = typ.__ptr.__nil;
         -- methods for embedded fields
         __gi_addMethodSynthesizer(function()
               local synthesizeMethod = function(target, m, f)
                  
                  if target.prototype[m.__prop] ~= nil then return end
                  
                  target.prototype[m.__prop] = function(self)
                     -- jea, temp comment out to figure spurious __gi_val source.
                     --local v = self.__gi_val[f.__prop];
                     if f.__typ == __gi_jsObjectPtr then
                        --v = new __gi_jsObjectPtr(v);
                        v = __gi_jsObjectPtr(v);
                     end
                     -- jea, temp comment out to figure where spurious __gi_val
                     -- is comfing from
                     --if v.__gi_val == nil then
                     --   --v = new f.__typ(v);
                     --   v = f.__typ(v);
                     --end
                     return (v[m.__prop])(v, arguments);
                  end
               end
               for i,f in pairs(fields) do
                  if (f.__anonymous) then
                     for i,m in pairs(__gi_methodSet(f.__typ)) do
                              synthesizeMethod(typ, m, f);
                              synthesizeMethod(typ.__ptr, m, f);
                     end
                     for i,m in pairs(__gi_methodSet(__gi_ptrType(f.__typ))) do
                        synthesizeMethod(typ.__ptr, m, f);
                     end
                  end
               end
         end);
      end

   --------------------------------------------
   --------------------------------------------
   --------------------------------------------
      
   else
      -- __gi_panic(new __gi_String("invalid kind: " .. kind));
      kind = kind or "<nil>"
      error("error at struct.lua:(maybe line 1187?): invalid kind: "..kind);
   end

   --big switch (kind) in js.
   if kind == __gi_kind_bool or
   kind == __gi_kind_Map then
      
      typ.__zero = function() return false; end
      
      
   elseif kind == __gi_kind_int or
      kind == __gi_kind_int8 or
      kind == __gi_kind_int16 or
      kind == __gi_kind_int32 or
      kind == __gi_kind_int64 or
      kind == __gi_kind_uint or
      kind == __gi_kind_uint8  or
      kind == __gi_kind_uint16 or
      kind == __gi_kind_uint32 or
      kind == __gi_kind_uint64 or
      kind == __gi_kind_uintptr or
   kind == __gi_kind_UnsafePointer then
      
      typ.__zero = function() return 0LL; end
      

   elseif kind == __gi_kind_float32 or
   kind == __gi_kind_float64 then

      typ.__zero = function() return 0; end
      
      
   elseif kind ==  __gi_kind_String then
      typ.__zero = function() return ""; end
      

   elseif kind ==  __gi_kind_complex64 or
   kind ==  __gi_kind_complex128 then
      
      -- hmm... how to translate this new typ(0, 0)from javascript?
      -- local __zero = new typ(0, 0);
      typ.__zero = function() return 0,0; end
      
      
   elseif kind ==  __gi_kind_Ptr or
   kind ==  __gi_kind_Slice then
      
      typ.__zero = function() return typ.__nil; end
      

   elseif kind ==  __gi_kind_Chan then
      
      typ.__zero = function() return __gi_chanNil; end
      

   elseif kind ==  __gi_kind_Func then
      
      typ.__zero = function() return __gi_throwNilPointerError; end
      

   elseif kind ==  __gi_kind_Interface then
      
      typ.__zero = function() return __gi_ifaceNil; end
      

   elseif kind ==  __gi_kind_Array then
      
      typ.__zero = function() 
         local arrayClass = __gi_nativeArray(typ.__elem.__kind);
         if arrayClass ~= Array then
            --return new arrayClass(typ.__len)
            return arrayClass(typ.__len)
         end
         --Local array = new Array(typ.__len)
         return  __gi_NewArray({}, typ.__elem.__kind, typ.__len, typ.__elem.__zero())
      end
      
      

   elseif kind ==  __gi_kind_Struct then

      --typ.__zero = function() return new typ.__ptr(); end
      typ.__zero = function() return typ.__ptr(); end
      

   else
      --__gi_panic("invalid kind: "..kind)
      error("invalid kind: "..kind)
   end

   typ.__id = __gi_typeIDCounter;
   __gi_typeIDCounter = __gi_typeIDCounter+1;
   typ.__size = size;
   typ.__kind = kind;
   typ.__str = str;
   typ.__named = named;
   typ.__pkg = pkgPath;
   typ.__exported = exported;
   typ.__methods_desc = {};
   typ.__methodsetCache = nil;
   typ.__comparable = true;
   typ.__shortPkg = shortPkg;
   typ.__shortTypeName = shortTypeName;
   
   return typ;
end

-------------------

-- distinct from __gi_ptrTyp in ptr.lua.
-- port of javascript $ptrType() function
-- for top level structs (the elem).

__ptrTypeFromElem = function(elem)
  local typ = elem.__ptr;
  if typ == nil then
     typ = __gi_NewType(8, __gi_kind_Ptr, elem.__shortPkg, "*"..elem.__shortTypeName, "*"..elem.__str, false, elem.__pkg, elem.__exported, nil);
     elem.__ptr = typ;

     -- this is where we set __elem on the typ.
     typ.__init(elem);
  end
  return typ;
end

-------------------

function __gi_mapArray(array, f)
   if array == nil or #array == 0 then
      return {}
   end
   local na = #array
   if array.__constructor == nil then
      print("--------  begin stacktrace  ----------")
      print(debug.traceback())
      print("--------  end stacktrace  ----------")
      error "why no __constructor on array??"
   end
   local newArray = array.__constructor(na);
   for i = 0,na-1 do
      newArray[i] = f(array[i]);
   end
   return newArray;
end

function __mapFuncOverTable(tbl, func)
   local newtbl = {}
   for i,v in pairs(tbl) do
      newtbl[i] = func(v)
   end
   return newtbl
end

-- straight port from gohperjs, not done or tested, yet.
-- It seems to be building from text a type signature...
-- then making a new type.

--helper
__type2str = function(t)
   if type(t) == "table" then
      local s = t.__str;
      if s == nil then
         s = tostring(t)
      end
      return s
   end
   return tostring(t)
end

__gi_funcTypes = {};
__gi_funcType = function(params, results, variadic)
   print("debug: __gi_funcType called")
   __st(params, "params")
   __st(results, "results")
   
   local paramsM = __mapFuncOverTable(params, __type2str)
   local resultsM = __mapFuncOverTable(results, __type2str)
   
   print("debug: paramsM = ")
   __st(paramsM)

   print("debug: resultsM = ")
   __st(resultsM)
   
   local typeKey = "params:"..table.concat(paramsM, ",") .. "_results:" .. table.concat(resultsM,",") .. "_variadic:" .. tostring(variadic);

   print("debug: typeKey = '".. typeKey.."'")
   
  local typ = __gi_funcTypes[typeKey];
  if typ == nil then
     local paramTypes = __mapFuncOverTable(params, function(ty) __type2str(ty); end);
     if variadic then

       -- jea: Hmm, I haven't figured why the substr(2) wants to chop off the first 2 char.
       -- print to see the difference:
       print("jea debug: paramTypes[paramTypes.length - 1].substr(2) = '"..paramTypes[paramTypes.length - 1].substr(2).."'   versus without the substr: '"..paramTypes[paramTypes.length - 1] .. "'")
       
       paramTypes[paramTypes.length - 1] = "..." .. paramTypes[paramTypes.length - 1]
       --paramTypes[paramTypes.length - 1] = "..." .. paramTypes[paramTypes.length - 1].substr(2)
       
    end
    local str = "func(" .. table.concat(paramTypes, ", ") .. ")";
    if #results == 1 then
       str = str.. " " .. __type2str(results[1])
    elseif #results > 1 then
       str = str.. " (" .. table.concat(__mapFuncOverTable(results, __type2str),  ", ") .. ")";
    end

    print("jea debug: final func signature is: '"..str.."'")
    
    typ = __gi_NewType(4, __gi_kind_Func, str, false, "", false, nil);
    __gi_funcTypes[typeKey] = typ;
    typ.__init(params, results, variadic);
  end
  return typ;
end

--
-- basic types
--
__type__bool = __gi_NewType(1, __gi_kind_bool, "", "bool", "bool", true, "", false, nil);
__type__int = __gi_NewType(8, __gi_kind_int, "", "int", "int", true, "", false, nil);
__type__int8 = __gi_NewType(1, __gi_kind_int8, "", "int8", "int8", true, "", false, nil);
__type__int16 = __gi_NewType(2, __gi_kind_int16, "", "int16", "int16", true, "", false, nil);
__type__int32 = __gi_NewType(4, __gi_kind_int32, "", "int32", "int32", true, "", false, nil);
__type__int64 = __gi_NewType(8, __gi_kind_int64, "", "int64", "int64", true, "", false, nil);
__type__uint = __gi_NewType(8, __gi_kind_uint, "", "uint", "uint", true, "", false, nil);
__type__uint8 = __gi_NewType(1, __gi_kind_uint8, "", "uint8", "uint8", true, "", false, nil);
__type__uint16 = __gi_NewType(2, __gi_kind_uint16, "", "uint16", "uint16", true, "", false, nil);
__type__uint32 = __gi_NewType(4, __gi_kind_uint32, "", "uint32", "uint32", true, "", false, nil);
__type__uint64 = __gi_NewType(8, __gi_kind_uint64, "", "uint64", "uint64", true, "", false, nil);
__type__uintptr = __gi_NewType(8, __gi_kind_uintptr, "", "uintptr", "uintptr", true, "", false, nil);
__type__float32 = __gi_NewType(4, __gi_kind_float32, "", "float32", "float32", true, "", false, nil);
__type__float64 = __gi_NewType(8, __gi_kind_float64, "", "float64", "float64", true, "", false, nil);
__type__complex64 = __gi_NewType(8, __gi_kind_complex64, "", "complex64", "complex64", true, "", false, nil);
__type__complex128 = __gi_NewType(16, __gi_kind_complex128, "", "complex128", "complex128", true, "", false, nil);
__type__String = __gi_NewType(8, __gi_kind_String, "", "string", "string", true, "", false, nil);
__type__UnsafePointer = __gi_NewType(8, __gi_kind_UnsafePointer, "", "unsafe.Pointer", "unsafe.Pointer", true, "", false, nil);


-- 

__equal = function(a, b, typ)
   if typ == __jsObjectPtr then
    return a == b;
    end
   
   local k = typ.__kind
   if k ==  __gi_kind_complex64 or
   k == __gi_kind_complex128 then

      return a.re == b.re and a.im == b.im;
      
   elseif
      k == __gi_kind_int64 or
      k == __gi_kind_int or
      k == __gi_kind_uint or
   k == __gi_kind_uint64 then
      
      return a == b;
      
   elseif k == __gi_kind_Array then

      local na = #a
      if #b ~= na then
         return false;
      end
      for i= 0,na-1 do
         if not __equal(a[i], b[i], typ.elem) then
            return false;
         end
      end
      return true;
      
   elseif k ==  __gi_kind_Struct then

      for i,f in pairs(typ.fields) do
         if not __equal(a[f.__prop], b[f.__prop], f.__typ) then
            return false;
         end
      end
      return true;

   elseif k == __gi_kind_Interface then
      return __interfaceIsEqual(a, b);

   else
      return a == b;
   end
end

__interfaceIsEqual = function(a, b) 
  if a == __ifaceNil or b == __ifaceNil then
         return a == b;
  end
  if a.__constructor ~= b.__constructor then
    return false;
  end
  if a.__constructor == __jsObjectPtr then
    return a.object == b.object;
  end
  if not a.__comparable then
    error("comparing uncomparable typ='" .. a.str .. "'");
  end
  return __equal(a, b, a.__constructor);
end
