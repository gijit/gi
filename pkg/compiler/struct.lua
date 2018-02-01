-- structs and interfaces

-- TODO: syncrhonize around this/deal with multi-threading?
--  may need to teach LuaJIT how to grab go mutexes or use sync.Atomics.
__gi_idCounter = 0;

__gi_PropsKey = {}
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

      -- make __tostring as local as possible,
      -- to avoid the infinite looping we got
      -- when it was higher up.
      __tostring = __structPrinter
   }
   methodset.__index = methodset
   
   local props = {__typename = name, __name="structProps", __nMethod=0}
   props[__gi_PropsKey] = props
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
   local methodset = {
      __name="ifaceMethodSet",
   }
   methodset.__index = methodset
   
   local props = {__typename = name, __name="ifaceProps"}
   props[__gi_PropsKey] = props
   props[__gi_MethodsetKey] = methodset
   props[__gi_BaseKey] = __gi_ifaceMT
   props.__tostring = __ifacePrinter
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
function __reg:AddMethod(si, siName, methodName, method)
   --print("__reg:AddMethod for '"..si.."' called with methodName ", methodName)
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

   -- new?
   if methodset[methodName] ~= nil then
      -- not new
   else
      -- new, count it.
      local props = methodset[__gi_PropsKey]
      props.__nMethod = props.__nMethod + 1
   end
   
   -- add the method
   methodset[methodName] = method
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

   print("__gi_assertType called, typ='", typ, "' value='", value, "', returnTuple='", returnTuple, "'")

   local isInterface = false
   local interfaceMethods = nil
   if __reg:IsInterface(typ) then
      isInterface = true
      interfaceMethods = __reg:GetInterface(typ)
   end
   
   local ok = false
   local missingMethod = ""
   
   local valueMethods = value[__gi_MethodsetKey]
   local valueProps = value[__gi_PropsKey]
   
   local nMethod = valeProps.__nMethod
   
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
         
         local ctor
         -- (value == __gi_ifaceNil ? "" : value.constructor.string)
         if value == __gi_ifaceNil then 
            ctor = ""
         else
            ctor = value.constructor.string
         end
         error("runtime.TypeAssertionError."..typ.str.." is missing '"..missingMethod.."'")
         -- __gi_panic(new __gi_packages["runtime"].TypeAssertionError.ptr("", ctor, typ.string, missingMethod)
                    
      elseif returnTuple == 1 then
         return false
      else
         return zeroVal, false
      end
   end
   
   if not isInterface then
      value = value.__gi_val;
   end
   
   if typ == __gi_jsObjectPtr then
      value = value.object;
   end
   
   if returnTupe == 0 then
      return value
   elseif returnTuple == 1 then
      return true
   end
   return value, true

end

-- support for __gi_NewType

__gi_kindBool = 1;
__gi_kindInt = 2;
__gi_kindInt8 = 3;
__gi_kindInt16 = 4;
__gi_kindInt32 = 5;
__gi_kindInt64 = 6;
__gi_kindUint = 7;
__gi_kindUint8 = 8;
__gi_kindUint16 = 9;
__gi_kindUint32 = 10;
__gi_kindUint64 = 11;
__gi_kindUintptr = 12;
__gi_kindFloat32 = 13;
__gi_kindFloat64 = 14;
__gi_kindComplex64 = 15;
__gi_kindComplex128 = 16;
__gi_kindArray = 17;
__gi_kindChan = 18;
__gi_kindFunc = 19;
__gi_kindInterface = 20;
__gi_kindMap = 21;
__gi_kindPtr = 22;
__gi_kindSlice = 23;
__gi_kindString = 24;
__gi_kindStruct = 25;
__gi_kindUnsafePointer = 26;

__gi_methodSynthesizers = {}
__gi_addMethodSynthesizer = function(f) 
  if __gi_methodSynthesizers == nil then
    f();
    return;
  end
  __gi_methodSynthesizers.push(f);
end

__gi_synthesizeMethods = function() 
  __gi_methodSynthesizers.forEach(function(f) f(); end);
  __gi_methodSynthesizers = nil;
end

__gi_ifaceKeyFor = function(x)
  if x == __gi_ifaceNil then
    return "nil"
  end
  local c = x.constructor
  return c.string .. "__gi_" .. c.keyFor(x.__gi_val)
end

__gi_identity = function(x) return x; end

__gi_typeIDCounter = 0;

__gi_idKey = function(x) 
  if x.__gi_id == nil then
     __gi_idCounter = __gi_idCounter + 1
     x.__gi_id = __gi_idCounter;
  end
  return tostring(x.__gi_id);
end

__castableMT = {
   __call = function(t, ...)
      print("__castableMT __call() invoked, with ... = ", ...)
      local arg0 = ...
      print("arg0 is", arg0)
      t.__gi_val = arg0
   end
}

__gi_identity = function(x) return x; end

__gi_floatKey = function(f)
   if f ~= f then
      __gi_idCounter = __gi_idCounter+1
         return "NaN_" + __gi_idCounter;
   end
   return tostring(f);
end


-- create new type. 
-- translate __gi_newType() in js,
-- from gopherjs/compiler/prelude/types.go#L64
--
function __gi_NewType(size, kind, str, named, pkg, exported, constructor)

      print("size=",size,", kind=",kind,", str=",str)
      print("named=",named,", pkg=",pkg)
      print("exported=",exported,", constructor=",constructor)
      
      local typ = {}
      if kind == __gi_kindBool or
         kind == __gi_kindInt or
         kind == __gi_kindInt8 or
         kind == __gi_kindInt16 or
         kind == __gi_kindInt32 or
         kind == __gi_kindInt64 or
         kind == __gi_kindUint or
         kind == __gi_kindUint8 or
         kind == __gi_kindUint16 or
         kind == __gi_kindUint32 or
         kind == __gi_kindUint64 or
         kind == __gi_kindUintptr or
      kind == __gi_kindUnsafePointer then
         
         typ = {__gi_val, wrapped=true, keyFor=__gi_identity};
         setmetatable(typ, __castableMT);
         -- gopherjs:
         -- typ = function(v) this.__gi_val = v; end
         -- typ.wrapped = true;
         -- typ.keyFor = __gi_identity;
         break;
         
      elseif kind == __gi_kindString then

         typ = {__gi_val, wrapped=true};
         setmetatable(typ, __castableMT);
         -- typ = function(v) this.__gi_val = v; end
         -- typ.wrapped = true;
         typ.keyFor = function(x) return "__gi_"..x; end
         break;

      elseif kind ==  __gi_kindFloat32 or
      kind == __gi_kindFloat64 then

         typ = {__gi_val, wrapped=true};
         setmetatable(typ, __castableMT);         
         -- typ = function(v) { this.__gi_val = v; };
         typ.keyFor = function(x) return __gi_floatKey(x) end
         break;
               
  elseif kind == __gi_kindComplex64 then
     typ = function(real, imag)
      this.__gi_real = __gi_fround(real);
      this.__gi_imag = __gi_fround(imag);
      this.__gi_val = this;
    end
    typ.keyFor = function(x)  return x.__gi_real .. "__gi_" .. x.__gi_imag; end
    break;

  elseif kind == __gi_kindComplex128 then
    typ = function(real, imag) 
      this.__gi_real = real;
      this.__gi_imag = imag;
      this.__gi_val = this;
    end
    typ.keyFor = function(x)  return x.__gi_real .. "__gi_" .. x.__gi_imag; end
    break;

      elseif kind == __gi_kindArray then
         setmetatable(typ, __castableMT)
         --typ = function(v) this.__gi_val = v; end
         typ.wrapped = true;
         typ.ptr = __gi_NewType(8, __gi_kindPtr, "*" .. str, false, "", false, function(array) 
                                   this.__gi_get = function() return array; end;
                                   this.__gi_set = function(v) typ.copy(this, v); end
                                   this.__gi_val = array;
         end);
    typ.init = function(elem, len) 
       typ.elem = elem;
       typ.len = len;
       typ.comparable = elem.comparable;
       typ.keyFor = function(x) 
          return Array.prototype.join.call(__gi_mapArray(x, function(e) {
          return tostring(elem.keyFor(e)).replace(/\\/g, "\\\\").replace(/\__gi_/g, "\\__gi_");
        end), "__gi_");
       end
       typ.copy = function(dst, src) 
          __gi_copyArray(dst, src, 0, 0, src.length, elem);
       end
       typ.ptr.init(typ);
       Object.defineProperty(typ.ptr.nil, "nilCheck", { get: __gi_throwNilPointerError });
    end
    break;

  elseif kind == __gi_kindChan then
    typ = function(v) this.__gi_val = v; end
    typ.wrapped = true;
    typ.keyFor = __gi_idKey;
    typ.init = function(elem, sendOnly, recvOnly)
      typ.elem = elem;
      typ.sendOnly = sendOnly;
      typ.recvOnly = recvOnly;
    end
    break;

  elseif kind == __gi_kindFunc then
    typ = function(v) this.__gi_val = v; end
    typ.wrapped = true;
    typ.init = function(params, results, variadic)
      typ.params = params;
      typ.results = results;
      typ.variadic = variadic;
      typ.comparable = false;
    end
    break;

  elseif kind == __gi_kindInterface then
    typ = { implementedBy= {}, missingMethodFor= {} };
    typ.keyFor = __gi_ifaceKeyFor;
    typ.init = function(methods) 
      typ.methods = methods;
      methods.forEach(function(m) 
        __gi_ifaceNil[m.prop] = __gi_throwNilPointerError;
      end);
    end
    break;

  elseif kind == __gi_kindMap then
    typ = function(v) this.__gi_val = v; end
    typ.wrapped = true;
    typ.init = function(key, elem)
      typ.key = key;
      typ.elem = elem;
      typ.comparable = false;
    end
    break;

  elseif kind == __gi_kindPtr then
    typ = constructor  or  function(getter, setter, target)
      this.__gi_get = getter;
      this.__gi_set = setter;
      this.__gi_target = target;
      this.__gi_val = this;
    end
    typ.keyFor = __gi_idKey;
    typ.init = function(elem) 
      typ.elem = elem;
      typ.wrapped = (elem.kind == __gi_kindArray);
      typ.nil = new typ(__gi_throwNilPointerError, __gi_throwNilPointerError);
    end
    break;

  elseif kind == __gi_kindSlice then
    typ = function(array)
      if array.constructor ~= typ.nativeArray then
        array = new typ.nativeArray(array);
      end
      this.__gi_array = array;
      this.__gi_offset = 0;
      this.__gi_length = array.length;
      this.__gi_capacity = array.length;
      this.__gi_val = this;
    end
    typ.init = function(elem)
      typ.elem = elem;
      typ.comparable = false;
      typ.nativeArray = __gi_nativeArray(elem.kind);
      typ.nil = new typ([]);
    end
    break;

  elseif kind == __gi_kindStruct then
    typ = function(v)  this.__gi_val = v; end
    typ.wrapped = true;
    typ.ptr = __gi_NewType(8, __gi_kindPtr, "*" .. str, false, pkg, exported, constructor);
    typ.ptr.elem = typ;
    typ.ptr.prototype.__gi_get = function()  return this; end
    typ.ptr.prototype.__gi_set = function(v) typ.copy(this, v); end
    typ.init = function(pkgPath, fields)
      typ.pkgPath = pkgPath;
      typ.fields = fields;
      fields.forEach(function(f) 
        if not f.typ.comparable then
          typ.comparable = false;
        end
      end);
      typ.keyFor = function(x) 
         local val = x.__gi_val;
         return __gi_mapArray(fields, function(f)
                                 return tostring(f.typ.keyFor(val[f.prop])).replace(/\\/g, "\\\\").replace(/\__gi_/g, "\\__gi_");
         end).join("__gi_");
      end
      typ.copy = function(dst, src) 
         for i = 0,fields.length-1 do
            local f = fields[i];
            local knd = f.typ.kind
            if knd ==  __gi_kindArray then
               -- do nothing
            elseif knd == __gi_kindStruct then
               f.typ.copy(dst[f.prop], src[f.prop]);
            else
               -- default:
               dst[f.prop] = src[f.prop];
            end
         end
      end
      -- nil value
      local properties = {};
      fields.forEach(function(f) 
            properties[f.prop] = { get: __gi_throwNilPointerError, set: __gi_throwNilPointerError };
      end);
      typ.ptr.nil = Object.create(constructor.prototype, properties);
      typ.ptr.nil.__gi_val = typ.ptr.nil;
      -- methods for embedded fields
      __gi_addMethodSynthesizer(function()
            local synthesizeMethod = function(target, m, f)
               
               if target.prototype[m.prop] ~= nil then return end
               
               target.prototype[m.prop] = function()
                  local v = this.__gi_val[f.prop];
                  if f.typ == __gi_jsObjectPtr then
                     v = new __gi_jsObjectPtr(v);
                  end
                  if v.__gi_val == nil then
                     v = new f.typ(v);
                  end
                  return v[m.prop].apply(v, arguments);
               end
            end
            fields.forEach(function(f)
                  if (f.anonymous) then
                     __gi_methodSet(f.typ).forEach(function(m) 
                           synthesizeMethod(typ, m, f);
                           synthesizeMethod(typ.ptr, m, f);
                                                  end)
                     __gi_methodSet(__gi_ptrType(f.typ)).forEach(function(m)
                           synthesizeMethod(typ.ptr, m, f);
                                                                end);
                  end
            end);
      end);
      end
break;

  else
    __gi_panic(new __gi_String("invalid kind: " .. kind));
  end

--switch (kind) {
if kind == __gi_kindBool or
kind == __gi_kindMap then
   
   typ.zero = function() return false; end
   break;
    
elseif kind == __gi_kindInt or
   kind == __gi_kindInt8 or
   kind == __gi_kindInt16 or
   kind == __gi_kindInt32 or
   kind == __gi_kindInt64 or
   kind == __gi_kindUint or
   kind == __gi_kindUint8  or
   kind == __gi_kindUint16 or
   kind == __gi_kindUint32 or
   kind == __gi_kindUint64 or
   kind == __gi_kindUintptr or
kind == __gi_kindUnsafePointer then
   
   typ.zero = function() return 0LL; end
   break;

elseif kind == __gi_kindFloat32 or
kind == __gi_kindFloat64 then

   typ.zero = function() return 0; end
   break;
      
elseif kind ==  __gi_kindString then
    typ.zero = function() return ""; end
    break;

elseif kind ==  __gi_kindComplex64 or
kind ==  __gi_kindComplex128 then
   
   -- hmm... how to translate this new typ(0, 0)from javascript?
   -- local zero = new typ(0, 0);
   typ.zero = function() return 0,0; end
   break;
   
elseif kind ==  __gi_kindPtr or
kind ==  __gi_kindSlice then
   
    typ.zero = function() return typ.nil; end
    break;

elseif kind ==  __gi_kindChan then
   
    typ.zero = function() return __gi_chanNil; end
    break;

elseif kind ==  __gi_kindFunc then
   
    typ.zero = function() return __gi_throwNilPointerError; end
    break;

elseif kind ==  __gi_kindInterface then
   
    typ.zero = function() return __gi_ifaceNil; end
    break;

elseif kind ==  __gi_kindArray then
   
   typ.zero = function() 
      local arrayClass = __gi_nativeArray(typ.elem.kind);
      if arrayClass ~= Array then
         return new arrayClass(typ.len)
      end
      local array = new Array(typ.len)
      for i = 0, typ.len-1 do
         array[i] = typ.elem.zero()
      end
      return array;
   end
   
   break;

elseif kind ==  __gi_kindStruct then

   typ.zero = function() return new typ.ptr(); end
   break;

else
   __gi_panic("invalid kind: "..kind)
end

  typ.id = __gi_typeIDCounter;
  __gi_typeIDCounter = __gi_typeIDCounter+1;
  typ.size = size;
  typ.kind = kind;
  typ.str = str;
  typ.named = named;
  typ.pkg = pkg;
  typ.exported = exported;
  typ.methods = {};
  typ.methodSetCache = nil;
  typ.comparable = true;
  
  return typ;
end

