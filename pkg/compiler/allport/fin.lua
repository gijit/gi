dofile '../math.lua' -- for __max, __min, __truncateToInt

dofile '../int64.lua'
__ffi = require("ffi")
__bit =require("bit")

__global ={};
__module ={};
__packages = {}
__idCounter = 0;


__mod = function(y) return x % y; end;
__parseInt = parseInt;
__parseFloat = function(f)
  if f ~= nil  and  f ~= nil  and  f.constructor == Number then
    return f;
  end
  return parseFloat(f);
end;

--[[
 __froundBuf = Float32Array(1);
__fround = Math.fround  or  function(f)
  __froundBuf[0] = f;
  return __froundBuf[0];
end;
--]]

--[[
__imul = Math.imul  or  function(b)
   local ah = __bit.band(__bit.rshift(a, 16), 0xffff);
   local al = __bit.band(a, 0xffff);
   local bh = __bit.band(__bit.rshift(b, 16), 0xffff);
   local bl = __bit.band(b, 0xffff);
   return ((al * bl) + __bit.arshift((__bit.rshift(__bit.lshift(ah * bl + al * bh), 16), 0), 0);
end;
--]]

__floatKey = function(f)
  if f ~= f then
     __idCounter=__idCounter+1;
    return "NaN__" .. tostring(__idCounter);
  end
  return tostring(f);
end;

__flatten64 = function(x)
  return x.__high * 4294967296 + x.__low;
end;


__Infinity = math.huge

__kindBool = 1;
__kindInt = 2;
__kindInt8 = 3;
__kindInt16 = 4;
__kindInt32 = 5;
__kindInt64 = 6;
__kindUint = 7;
__kindUint8 = 8;
__kindUint16 = 9;
__kindUint32 = 10;
__kindUint64 = 11;
__kindUintptr = 12;
__kindFloat32 = 13;
__kindFloat64 = 14;
__kindComplex64 = 15;
__kindComplex128 = 16;
__kindArray = 17;
__kindChan = 18;
__kindFunc = 19;
__kindInterface = 20;
__kindMap = 21;
__kindPtr = 22;
__kindSlice = 23;
__kindString = 24;
__kindStruct = 25;
__kindUnsafePointer = 26;

-- jea: sanity check my assumption by comparing
-- length with #a
function __assertIsArray(a)
   local n = 0
   for k,v in pairs(a) do
      n=n+1
   end
   if #a ~= n then
      error("not an array, __assertIsArray failed")
   end
end


-- length of array, counting [0] if present.
function __lenz(array)      
   local n = #array
   if array[0] ~= nil then
      n=n+1
   end
   return n
end

-- st or showtable, a debug print helper.
-- seen avoids infinite looping on self-recursive types.
function __st(t, name, indent, quiet, methods_desc, seen)
   if t == nil then
      local s = "<nil>"
      if not quiet then
         print(s)
      end
      return s
   end

   seen = seen or {}
   if seen[t] ~= nil then
      return
   end
   seen[t] =true   
   
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

   -- get address, avoiding infinite loop of self-calls.
   local mt = getmetatable(t)
   setmetatable(t, nil)
   local addr = tostring(t) 
   setmetatable(t, mt)
   
   local k = 0
   local name = name or ""
   local namec = name
   if name ~= "" then
      namec = namec .. ": "
   end
   local indent = indent or 0
   local pre = string.rep(" ", 4*indent)..namec
   local s = pre .. "============================ "..addr.."\n"
   for i,v in pairs(t) do
      k=k+1
      local vals = ""
      if methods_desc then
         --print("methods_desc is true")
         --vals = __st(v,"",indent+1,quiet,methods_desc, seen)
      else 
         vals = tostring(v)
      end
      s = s..pre.." "..tostring(k).. " key: '"..tostring(i).."' val: '"..vals.."'\n"
   end
   if k == 0 then
      s = pre.."<empty table>"
   end

   --local mt = getmetatable(t)
   if mt ~= nil then
      s = s .. "\n"..__st(mt, "mt.of."..name, indent+1, true, methods_desc, seen)
   end
   if not quiet then
      print(s)
   end
   return s
end


-- apply fun to each element of the array arr,
-- then concatenate them together with splice in
-- between each one. It arr is empty then we
-- return the empty string. arr can start at
-- [0] or [1].
function __mapAndJoinStrings(splice, arr, fun)
   local newarr = {}
   -- handle a zero argument, if present.
   local bump = 0
   local zval = arr[0]
   if zval ~= nil then
      bump = 1
      newarr[1] = fun(zval)
   end
   for i,v in ipairs(arr) do
      newarr[i+bump] = fun(v)
   end
   return table.concat(newarr, splice)
end

-- return sorted keys from table m
__keys = function(m)
   if type(m) ~= "table" then
      return {}
   end
   local r = {}
   for k in pairs(m) do
      local tyk = type(k)
      if tyk == "function" then
         k = tostring(k)
      end
      table.insert(r, k)
   end
   table.sort(r)
   return r
end

--
__flushConsole = function() end;
__throwRuntimeError = function(...) error(...) end
__throwNilPointerError = function()  __throwRuntimeError("invalid memory address or nil pointer dereference"); end;
__call = function(fn, rcvr, args)  return fn(rcvr, args); end;
__makeFunc = function(fn)
   return function()
      return __externalize(fn(this, (__sliceType({},__jsObjectPtr))(__global.Array.prototype.slice.call(arguments, {}))), __emptyInterface);
   end;
end;
__unused = function(v) end;

--
__mapArray = function(array, f)
   local newarr = {}
   -- handle a zero argument, if present.
   local bump = 0
   local zval = arr[0]
   if zval ~= nil then
      bump = 1
      newarr[1] = fun(zval)
   end
   for i,v in ipairs(arr) do
      newarr[i+bump] = fun(v)
   end
   return newarr
end;

__methodVal = function(recv, name) 
  local vals = recv.__methodVals  or  {};
  recv.__methodVals = vals; -- /* noop for primitives */
  local f = vals[name];
  if f ~= nil then
    return f;
  end
  local method = recv[name];
  f = function() 
     __stackDepthOffset = __stackDepthOffset-1;
     -- try
     local res = {pcall(function()
           return recv[method](arguments);
     end)}
        -- finally
     __stackDepthOffset=__stackDepthOffset+1;
     -- no catch, so either re-throw or return results
     local ok, err = unpack(res)
     if not ok then
        -- rethrow
        error(err)
     end
     -- return results (without the ok/not first value)
     return table.remove(res, 1)
  end;
  vals[name] = f;
  return f;
end;

__methodExpr = function(typ, name) 
   local method = typ.prototype[name];
   if method.__expr == nil then
      method.__expr = function() 
         __stackDepthOffset=__stackDepthOffset-1;

         -- try
         local res ={pcall(
            function()
               if typ.wrapped then
                  arguments[0] = typ(arguments[0]);
               end
               return method(arguments);
         end)}
         local ok, threw = unpack(res)
         -- finally
         __stackDepthOffset=__stackDepthOffset+1;
         -- no catch, so rethrow any exception
         if not ok then
            error(threw)
         end
         return table.remove(res, 1)
      end;
   end
   return method.__expr;
end;

__ifaceMethodExprs = {};
__ifaceMethodExpr = function(name) 
  local expr = __ifaceMethodExprs["_"  ..  name];
  if expr == nil then
     expr = function()
        __stackDepthOffset = __stackDepthOffset-1;
        -- try
        local res = {pcall(
                        function()
                           return Function.call.apply(arguments[0][name], arguments);
        end)}
        -- finally
        __stackDepthOffset = __stackDepthOffset+1;
        -- no catch
        local ok, threw = unpack(res)
        if not ok then
           error(threw)
        else
           -- non-panic return from pcall
           return table.remove(res, 1)
        end   
     end;
     __ifaceMethodExprs["_"  ..  name] = expr
  end
  return expr;
end;

--

__substring = function(str, low, high)
  if low < 0  or  high < low  or  high > #str then
    __throwRuntimeError("string slice bounds out of range");
  end
  return string.sub(str, low+1, high); -- high is inclusive, so no +1 needed.
end;

__sliceToArray = function(slice)
   local cp = {}
   if slice.__length > 0 then
      local k = 0
      for i = slice.__offset, slice.__offset + slice.__length -1 do
         cp[k] = slice.array[i]
         k=k+1
      end
   end
   return cp
end;

--


--

__valueBasicMT = {
   __name = "__valueBasicMT",
   __tostring = function(self, ...)
      --print("__tostring called from __valueBasicMT")
      if type(self.__val) == "string" then
         return '"'..self.__val..'"'
      end
      if self ~= nil and self.__val ~= nil then
         --print("__valueBasicMT.__tostring called, with self.__val set.")
         if self.__val == self then
            -- not a basic value, but a pointer, array, slice, or struct.
            return "<this.__val == this; avoid inf loop>"
         end
         --return tostring(self.__val)
      end
      if getmetatable(self.__val) == __valueBasicMT then
         --print("avoid infinite loop")
         return "<avoid inf loop>"
      else
         return tostring(self.__val)
      end
   end
}

__valueArrayMT = {
   __name = "__valueArrayMT",
   
   __newindex = function(t, k, v)
      --print("__valueArrayMT.__newindex called, t is:")
      --__st(t)

      if k < 0 or k >= #t then
         error "read of array error: access out-of-bounds"
      end
      
      t.__val[k] = v
   end,
   
   __index = function(t, k)
      --print("__valueArrayMT.__index called, k='"..tostring(k).."'; t.__val is:")
      --__st(t.__val)
      if k < 0 or k >= #t then
         error("write to array error: access out-of-bounds; "..tostring(k).." is outside [0, "  .. tostring(#t) .. ")")
      end
      
      return t.__val[k]
   end,

   __len = function(t)
      return int(__lenz(t.__val))
   end,
   
   __tostring = function(self, ...)
      --print("__tostring called from __valueArrayMT")
      if type(self.__val) == "string" then
         return '"'..self.__val..'"'
      end
      if self ~= nil and self.__val ~= nil then
         --print("__valueArrayMT.__tostring called, with self.__val set.")
         if self.__val == self then
            -- not a basic value, but a pointer, array, slice, or struct.
            return "<this.__val == this; avoid inf loop>"
         end
         --return tostring(self.__val)
      end
      if getmetatable(self.__val) == __valueArrayMT then
         --print("avoid infinite loop")
         return "<avoid inf loop>"
      else
         return tostring(self.__val)
      end
   end
}

__valueSliceMT = {
   __name = "__valueSliceMT",
   
   __newindex = function(t, k, v)
      --print("__valueSliceMT.__newindex called, t is:")
      --__st(t)
      local w = t.__offset + k
      if k < 0 or k >= t.__capacity then
         error "slice error: write out-of-bounds"
      end
      t.__array[w] = v
   end,
   
   __index = function(t, k)
      --print("__valueSliceMT.__index called, k='"..tostring(k).."'; t.__val is:")
      --__st(t.__val)
      local w = t.__offset + k
      if k < 0 or k >= t.__capacity then
         error "slice error: access out-of-bounds"
      end
      return t.__array[w]
   end,

   __len = function(t)
      return t.__length
   end,
   
   __tostring = function(self, ...)
      --print("__tostring called from __valueSliceMT")
      if type(self.__val) == "string" then
         return '"'..self.__val..'"'
      end
      if self ~= nil and self.__val ~= nil then
         --print("__valueSliceMT.__tostring called, with self.__val set.")
         if self.__val == self then
            -- not a basic value, but a pointer, array, slice, or struct.
            return "<this.__val == this; avoid inf loop>"
         end
         --return tostring(self.__val)
      end
      if getmetatable(self.__val) == __valueSliceMT then
         --print("avoid infinite loop")
         return "<avoid inf loop>"
      else
         return tostring(self.__val)
      end
   end
}


__tfunBasicMT = {
   __name = "__tfunBasicMT",
   __call = function(self, ...)
      --print("jea debug: __tfunBasicMT.__call() invoked") -- , self='"..tostring(self).."' with tfun = ".. tostring(self.tfun).. " and args=")
      --print(debug.traceback())
      
      --print("in __tfunBasicMT, start __st on ...")
      --__st({...}, "__tfunBasicMT.dots")
      --print("in __tfunBasicMT,   end __st on ...")

      --print("in __tfunBasicMT, start __st on self")
      --__st(self, "self")
      --print("in __tfunBasicMT,   end __st on self")

      local newInstance = {}
      setmetatable(newInstance, __valueBasicMT)
      if self ~= nil then
         if self.tfun ~= nil then
            --print("calling tfun! -- let constructors set metatables if they wish to.")

            -- get zero value if no args
            if #{...} == 0 and self.zero ~= nil then
               --print("tfun sees no args and we have a typ.zero() method, so invoking it")
               self.tfun(newInstance, self.zero())
            else
               self.tfun(newInstance, ...)
            end
         end
      else
         if self ~= nil then
            --print("self.tfun was nil")
         end
      end
      return newInstance
   end
}


function __newAnyArrayValue(elem, len)
   local array = {}
   for i =0, len -1 do
      array[i]= elem.zero();
   end
   return array;
end


__methodSynthesizers = {};
__addMethodSynthesizer = function(f)
   if __methodSynthesizers == nil then
      f();
      return;
   end
   table.insert(__methodSynthesizers, f);
end;

__synthesizeMethods = function()
   for i,f in ipairs(__methodSynthesizers) do
      f();
   end
   __methodSynthesizers = nil;
end;

__ifaceKeyFor = function(x)
  if x == __ifaceNil then
    return 'nil';
  end
  local c = x.constructor;
  return c.__str .. '__' .. c.keyFor(x.__val);
end;

__identity = function(x) return x; end;

__typeIDCounter = 0;

__idKey = function(x)
   if x.__id == nil then
      __idCounter=__idCounter+1;
      x.__id = __idCounter;
   end
   return String(x.__id);
end;

__newType = function(size, kind, str, named, pkg, exported, constructor)
   local typ ={};
   setmetatable(typ, __tfunBasicMT)

   if kind ==  __kindBool or
      kind == __kindInt or 
      kind == __kindInt8 or 
      kind == __kindInt16 or 
      kind == __kindInt32 or 
      kind == __kindInt64 or 
      kind == __kindUint or 
      kind == __kindUint8 or 
      kind == __kindUint16 or 
      kind == __kindUint32 or 
      kind == __kindUint64 or 
      kind == __kindUintptr or 
   kind == __kindUnsafePointer then

      -- jea: I observe that
      -- primitives have: this.__val ~= v; and are the types are
      -- distinguished with typ.wrapped = true; versus
      -- all table based values, that have: this.__val == this;
      -- and no .wrapped field.
      --
      typ.tfun = function(this, v) this.__val = v; end;
      typ.wrapped = true;
      typ.keyFor = __identity;

   elseif kind == __kindString then
      
      typ.tfun = function(this, v)
         --print("strings' tfun called! with v='"..tostring(v).."' and this:")
         --__st(this)
         this.__val = v; end;
      typ.wrapped = true;
      typ.keyFor = function(x) return "_" .. x; end;

   elseif kind == __kindFloat32 or
   kind == __kindFloat64 then
      
      typ.tfun = function(this, v) this.__val = v; end;
      typ.wrapped = true;
      typ.keyFor = function(x) return __floatKey(x); end;


  elseif kind ==  __kindComplex64 then 
    typ.tfun = function(this, real, imag)
      this.__real = __fround(real);
      this.__imag = __fround(imag);
      this.__val = this;
    end;
    typ.keyFor = function(x) return x.__real .. "_" .. x.__imag; end;
    

  elseif kind ==  __kindComplex128 then 
    typ.tfun = function(this, real, imag)
      this.__real = real;
      this.__imag = imag;
      this.__val = this;
    end;
    typ.keyFor = function(x) return x.__real .. "_" .. x.__imag; end;
    
      
   elseif kind ==  __kindPtr then
      
      typ.tfun = constructor  or
         function(this, getter, setter, target)
            --print("pointer typ.tfun which is same as constructor called! getter='"..tostring(getter).."'; setter='"..tostring(setter).."; target = '"..tostring(target).."'")
            this.__get = getter;
            this.__set = setter;
            this.__target = target;
            this.__val = this; -- seems to indicate a non-primitive value.
         end;
      typ.keyFor = __idKey;
      typ.init = function(elem)
         typ.elem = elem;
         typ.wrapped = (elem.kind == __kindArray);
         typ.__nil = typ(__throwNilPointerError, __throwNilPointerError);
      end;

   elseif kind ==  __kindSlice then
      
      typ.tfun = function(this, array)
         this.__array = array;
         this.__offset = 0;
         this.__length = __lenz(array);
         this.__capacity = this.__length;
         --print("jea debug: slice tfun set __length to ", this.__length)
         --print("jea debug: slice tfun set __capacity to ", this.__capacity)
         --print("jea debug: slice tfun sees array: ")
         --for i,v in pairs(array) do
         --   print("array["..tostring(i).."] = ", v)
         --end
         
         this.__val = this;
         this.__constructor = typ
         setmetatable(this, __valueSliceMT)
      end;
      typ.init = function(elem)
         typ.elem = elem;
         typ.comparable = false;
         typ.__nil = typ({},{});
      end;
      
   elseif kind ==  __kindArray then

      typ.tfun = function(this, v)
         --print("in tfun ctor function for __kindArray")
         this.__val = v;
         setmetatable(this, __valueArrayMT)
      end;
      typ.wrapped = true;
      typ.ptr = __newType(4, __kindPtr, "*" .. str, false, "", false, function(this, array)
                             this.__get = function() return array; end;
                             this.__set = function(v) typ.copy(this, v); end;
                             this.__val = array;
      end);
      typ.init = function(elem, len)
         typ.elem = elem;
         typ.len = len;
         typ.comparable = elem.comparable;
         typ.keyFor = function(x)
            return __mapAndJoinStrings("_", x, function(e)
                                          return string.gsub(tostring(elem.keyFor(e)), "\\", "\\\\")
            end)
         end
         typ.copy = function(dst, src)
            __copyArray(dst, src, 0, 0, #src, elem);
         end;
         typ.ptr.init(typ);

         -- TODO:
         -- jea: nilCheck allows asserting that a pointer is not nil before accessing it.
         -- jea: what seems odd is that the state of the pointer is
         -- here defined on the type itself, and not on the particular instance of the
         -- pointer. But perhaps this is javascript's prototypal inheritence in action.
         --
         -- gopherjs uses them in comma expressions. example, condensed:
         --     p$1 = new ptrType(...); sa$3.Port = (p$1.nilCheck, p$1[0])
         --
         -- Since comma expressions are not (efficiently) supported in Lua, let
         -- implement the nil check in a different manner.
         -- js: Object.defineProperty(typ.ptr.__nil, "nilCheck", { get= __throwNilPointerError end);
      end;
      -- end __kindArray

   elseif kind ==  __kindStruct then
      
      typ.tfun = function(this, v) this.__val = v; end;
      typ.wrapped = true;

      local ctor = function(this, ...)
         this.__get = function() return this; end;
         this.__set = function(v) typ.copy(this, v); end;
         constructor(this, ...);
      end
      typ.ptr = __newType(4, __kindPtr, "*" .. str, false, pkg, exported, ctor);
      -- __newType sets typ.comparable = true
      
      typ.ptr.elem = typ;
      typ.init = function(pkgPath, fields)
         typ.pkgPath = pkgPath;
         typ.fields = fields;
         for i,f in ipairs(fields) do
            if not f.typ.comparable then
               typ.comparable = false;
               break;
            end
         end
         typ.keyFor = function(x)
            local val = x.__val;
            return __mapAndJoinStrings("_", fields, function(f)
                                          return string.gsub(tostring(f.typ.keyFor(val[f.prop])), "\\", "\\\\")
            end)
         end;
         typ.copy = function(dst, src)
            for i=0,#fields-1 do
               local f = fields[i];
               local sw2 = f.typ.kind
               
               if sw2 == __kindArray or
               sw2 ==  __kindStruct then 
                  f.typ.copy(dst[f.prop], src[f.prop]);
               else
                  dst[f.prop] = src[f.prop];
               end
            end
         end;
         -- /* nil value */
         local properties = {};
         for i,f in ipairs(fields) do
            properties[f.prop] = { get= __throwNilPointerError, set= __throwNilPointerError };
         end;
         typ.ptr.__nil = {} -- Object.create(constructor.prototype, properties);
         --if constructor ~= nil then
         --   constructor(typ.ptr.__nil)
         --end
         typ.ptr.__nil.__val = typ.ptr.__nil;
         -- /* methods for embedded fields */
         __addMethodSynthesizer(function()
               local synthesizeMethod = function(target, m, f)
                  if target.prototype[m.prop] ~= nil then return; end
                  target.prototype[m.prop] = function()
                     local v = this.__val[f.prop];
                     if f.typ == __jsObjectPtr then
                        v = __jsObjectPtr(v);
                     end
                     if v.__val == nil then
                        local w = {}
                        f.typ(w, v);
                        v = w
                     end
                     return v[m.prop](v, arguments);
                  end;
               end;
               for i,f in ipairs(fields) do
                  if f.anonymous then
                     __methodSet(f.typ).forEach(function(m)
                           synthesizeMethod(typ, m, f);
                           synthesizeMethod(typ.ptr, m, f);
                                               end);
                     __methodSet(__ptrType(f.typ)).forEach(function(m)
                           synthesizeMethod(typ.ptr, m, f);
                                                          end);
                  end
               end;
         end);
      end;
      
   else
      error("invalid kind: " .. tostring(kind));
   end
   
   -- set zero() method
   if kind == __kindBool or
   kind ==__kindMap then
      typ.zero = function() return false; end;

   elseif kind == __kindInt or
      kind ==  __kindInt8 or
      kind ==  __kindInt16 or
      kind ==  __kindInt32 or
   kind ==  __kindInt64 then
      typ.zero = function() return 0LL; end;
      
   elseif kind ==  __kindUint or
      kind ==  __kindUint8  or
      kind ==  __kindUint16 or
      kind ==  __kindUint32 or
      kind ==  __kindUint64 or
      kind ==  __kindUintptr or
   kind ==  __kindUnsafePointer then
      typ.zero = function() return 0ULL; end;

   elseif   kind ==  __kindFloat32 or
   kind ==  __kindFloat64 then
      typ.zero = function() return 0; end;
      
   elseif kind ==  __kindString then
      typ.zero = function() return ""; end;

   elseif kind == __kindPtr or
   kind == __kindSlice then
      
      typ.zero = function() return typ.__nil; end;

   elseif kind == __kindArray then

      typ.zero = function()
         return __newAnyArrayValue(typ.elem, typ.len)
      end;

   elseif kind == __kindStruct then
      typ.zero = function()
         return typ.ptr();
      end;      
   end

   typ.id = __typeIDCounter;
   __typeIDCounter=__typeIDCounter+1;
   typ.size = size;
   typ.kind = kind;
   typ.__str = str;
   typ.named = named;
   typ.pkg = pkg;
   typ.exported = exported;
   typ.methods = {};
   typ.methodSetCache = nil;
   typ.comparable = true;
   return typ;
   
end

__Bool          = __newType( 1, __kindBool,          "bool",           true, "", false, nil);
__Int           = __newType( 8, __kindInt,           "int",            true, "", false, nil);
__Int8          = __newType( 1, __kindInt8,          "int8",           true, "", false, nil);
__Int16         = __newType( 2, __kindInt16,         "int16",          true, "", false, nil);
__Int32         = __newType( 4, __kindInt32,         "int32",          true, "", false, nil);
__Int64         = __newType( 8, __kindInt64,         "int64",          true, "", false, nil);
__Uint          = __newType( 8, __kindUint,          "uint",           true, "", false, nil);
__Uint8         = __newType( 1, __kindUint8,         "uint8",          true, "", false, nil);
__Uint16        = __newType( 2, __kindUint16,        "uint16",         true, "", false, nil);
__Uint32        = __newType( 4, __kindUint32,        "uint32",         true, "", false, nil);
__Uint64        = __newType( 8, __kindUint64,        "uint64",         true, "", false, nil);
__Uintptr       = __newType( 8, __kindUintptr,       "uintptr",        true, "", false, nil);
__Float32       = __newType( 8, __kindFloat32,       "float32",        true, "", false, nil);
__Float64       = __newType( 8, __kindFloat64,       "float64",        true, "", false, nil);
--__Complex64     = __newType( 8, __kindComplex64,     "complex64",      true, "", false, nil);
--__Complex128    = __newType(16, __kindComplex128,    "complex128",     true, "", false, nil);
__String        = __newType(16, __kindString,        "string",         true, "", false, nil);
--__UnsafePointer = __newType( 8, __kindUnsafePointer, "unsafe.Pointer", true, "", false, nil);


__ptrType = function(elem)
   local typ = elem.ptr;
   if typ == nil then
      typ = __newType(4, __kindPtr, "*" .. elem.__str, false, "", elem.exported, nil);
      elem.ptr = typ;
      typ.init(elem);
   end
   return typ;
end;

__arrayTypes = {};
__arrayType = function(elem, len)
   local typeKey = elem.id .. "_" .. len;
   local typ = __arrayTypes[typeKey];
   if typ == nil then
      typ = __newType(24, __kindArray, "[" .. len .. "]" .. elem.__str, false, "", false, nil);
      __arrayTypes[typeKey] = typ;
      typ.init(elem, len);
   end
   return typ;
end;

__copyArray = function(dst, src, dstOffset, srcOffset, n, elem)
   --print("__copyArray called with n = ", n, " dstOffset=", dstOffset, " srcOffset=", srcOffset)
   --print("__copyArray has dst:")
   --__st(dst)
   --print("__copyArray has src:")
   --__st(src)
   
   n = tonumber(n)
   if n == 0  or  (dst == src  and  dstOffset == srcOffset) then
      return;
   end

   local sw = elem.kind
   if sw == __kindArray or sw == __kindStruct then
      
      if dst == src  and  dstOffset > srcOffset then
         for i = n-1,0,-1 do
            elem.copy(dst[dstOffset + i], src[srcOffset + i]);
         end
         return;
      end
      for i = 0,n-1 do
         elem.copy(dst[dstOffset + i], src[srcOffset + i]);
      end
      return;
   end

   if dst == src  and  dstOffset > srcOffset then
      for i = n-1,0,-1 do
         dst[dstOffset + i] = src[srcOffset + i];
      end
      return;
   end
   for i = 0,n-1 do
      dst[dstOffset + i] = src[srcOffset + i];
   end
end;


-- __basicValue2kind: identify type of basic value

function __basicValue2kind(v)

   local ty = type(v)
   if ty == "cdata" then
      local cty = __ffi.typeof(v)
      if cty == int64 then
         return __kindInt
      elseif cty == int8 then
         return __kindInt8
      elseif cty == int16 then
         return __kindInt16
      elseif cty == int32 then
         return __kindInt32
      elseif cty == int64 then
         return __kindInt64
      elseif cty == uint then
         return __kindUint
      elseif cty == uint8 then
         return __kindUint8
      elseif cty == uint16 then
         return __kindUint16
      elseif cty == uint32 then
         return __kindUint32
      elseif cty == uint64 then
         return __kindUint64
      elseif cty == float32 then
         return __kindFloat32
      elseif cty == float64 then
         return __kindFloat64         
      else
         error("__basicValue2kind: unhandled cdata cty: '"..tostring(cty).."'")
      end      
   elseif ty == "boolean" then
      return __kindBool;
   elseif ty == "number" then
      return __kindFloat64
   elseif ty == "string" then
      return __kindString
   end
   error("__basicValue2kind: unhandled ty: '"..ty.."'")   
end

__sliceType = function(elem)
   local typ = elem.slice;
   if typ == nil then
      typ = __newType(24, __kindSlice, "[]" .. elem.__str, false, "", false, nil);
      elem.slice = typ;
      typ.init(elem);
   end
   return typ;
end;

__makeSlice = function(typ, length, capacity)
   length = tonumber(length)
   if capacity == nil then
      capacity = length
   else
      capacity = tonumber(capacity)
   end
   if length < 0  or  length > 9007199254740992 then -- 2^53
      __throwRuntimeError("makeslice: len out of range");
   end
   if capacity < 0  or  capacity < length  or  capacity > 9007199254740992 then
      __throwRuntimeError("makeslice: cap out of range");
   end
   local array = __newAnyArrayValue(typ.elem, capacity)
   local slice = typ(array);
   slice.__length = length;
   return slice;
end;


__subslice = function(slice, low, high, max)
   if high == nil then
      
   end
   if low < 0  or  (high ~= nil and high < low)  or  (max ~= nil and high ~= nil and max < high)  or  (high ~= nil and high > slice.__capacity)  or  (max ~= nil and max > slice.__capacity) then
      __throwRuntimeError("slice bounds out of range");
   end
   
   local s = {}
   slice.__constructor.tfun(s, slice.__array);
   s.__offset = slice.__offset + low;
   s.__length = slice.__length - low;
   s.__capacity = slice.__capacity - low;
   if high ~= nil then
      s.__length = high - low;
   end
   if max ~= nil then
      s.__capacity = max - low;
   end
   return s;
end;

__copySlice = function(dst, src)
   local n = __min(src.__length, dst.__length);
   __copyArray(dst.__array, src.__array, dst.__offset, src.__offset, n, dst.__constructor.elem);
   return n;
end;


__clone = function(src, typ)
  local clone = typ.zero();
  typ.copy(clone, src);
  return clone;
end;

__pointerOfStructConversion = function(obj, typ)
  if(obj.__proxies == nil) then
    obj.__proxies = {};
    obj.__proxies[obj.constructor.__str] = obj;
  end
  local proxy = obj.__proxies[typ.__str];
  if proxy == nil then
     local properties = {};
     
     local helper = function(p)
        properties[fieldProp] = {
           get= function() return obj[fieldProp]; end,
           set= function(value) obj[fieldProp] = value; end
        };
     end
     -- fields must be an array for this to work.
     for i=0,#typ.elem.fields-1 do
        helper(typ.elem.fields[i].prop);
     end
     
    proxy = Object.create(typ.prototype, properties);
    proxy.__val = proxy;
    obj.__proxies[typ.__str] = proxy;
    proxy.__proxies = obj.__proxies;
  end
  return proxy;
end;


__append = function(...)
   local arguments = {...}
   local slice = arguments[1]
   return __internalAppend(slice, arguments, 1, #arguments - 1);
end;

__appendSlice = function(slice, toAppend)
   if slice == nil then 
      error("error calling __appendSlice: slice must be available")
   end
   if toAppend == nil then
      error("error calling __appendSlice: toAppend must be available")      
   end
   if type(toAppend) == "string" then
      local bytes = __stringToBytes(toAppend);
      return __internalAppend(slice, bytes, 0, #bytes);
   end
   return __internalAppend(slice, toAppend.__array, toAppend.__offset, toAppend.__length);
end;

__internalAppend = function(slice, array, offset, length)
   if length == 0 then
      return slice;
   end

   local newArray = slice.__array;
   local newOffset = slice.__offset;
   local newLength = slice.__length + length;
   --print("jea debug: __internalAppend: newLength is "..tostring(newLength))
   local newCapacity = slice.__capacity;
   local elem = slice.__constructor.elem;

   if newLength > newCapacity then
      newOffset = 0;
      local tmpCap
      if slice.__capacity < 1024 then
         tmpCap = slice.__capacity * 2
      else
         tmpCap = __truncateToInt(slice.__capacity * 5 / 4)
      end
      newCapacity = __max(newLength, tmpCap);

      newArray = {}
      local w = slice.__offset
      for i = 0,slice.__length do
         newArray[i] = slice.__array[i + w]
      end
      for i = #slice,newCapacity-1 do
         newArray[i] = elem.zero();
      end
      
   end

   --print("jea debug, __internalAppend, newOffset = ", newOffset, " and slice.__length=", slice.__length)

   __copyArray(newArray, array, newOffset + slice.__length, offset, length, elem);
   --print("jea debug, __internalAppend, after copying over array:")
   --__st(newArray)

   local newSlice ={}
   slice.__constructor.tfun(newSlice, newArray);
   newSlice.__offset = newOffset;
   newSlice.__length = newLength;
   newSlice.__capacity = newCapacity;
   return newSlice;
end;



function field2strHelper(f)
   local tag = ""
   if f.tag ~= "" then
      tag = string.gsub(f.tag, "\\", "\\\\")
      tag = string.gsub(tag, "\"", "\\\"")
   end
   return f.name .. " " .. f.typ.__str .. tag
end

function typeKeyHelper(f)
   return f.name .. "," .. f.typ.id .. "," .. f.tag;
end

__structTypes = {};
__structType = function(pkgPath, fields)
   local typeKey = __mapAndJoinStrings("_", fields, typeKeyHelper)

   local typ = __structTypes[typeKey];
   if typ == nil then
      local str
      if #fields == 0 then
         str = "struct {}";
      else
         str = "struct { " .. __mapAndJoinStrings("; ", fields, field2strHelper) .. " }";
      end
      
      typ = __newType(0, __kindStruct, str, false, "", false, function()
                         local this = {}
                         this.__val = this;
                         for i = 0, #fields-1 do
                            local f = fields[i];
                            local arg = arguments[i];
                            if arg ~= nil then
                               this[f.prop] = arg
                            else
                               this[f.prop] = t.typ.zero();
                            end
                         end
                         return this
      end);
      __structTypes[typeKey] = typ;
      typ.init(pkgPath, fields);
   end
   return typ;
end;


__equal = function(a, b, typ)
  if typ == __jsObjectPtr then
    return a == b;
  end

  local sw = typ.kind
  if sw == __kindComplex64 or
  sw == __kindComplex128 then
     return a.__real == b.__real  and  a.__imag == b.__imag;
     
  elseif sw == __kindInt64 or
         sw == __kindUint64 then 
     return a.__high == b.__high  and  a.__low == b.__low;
     
  elseif sw == __kindArray then 
    if #a ~= #b then
      return false;
    end
    for i=0,#a-1 do
      if  not __equal(a[i], b[i], typ.elem) then
        return false;
      end
    end
    return true;
    
  elseif sw == __kindStruct then
     
    for i = 0,#(typ.fields)-1 do
      local f = typ.fields[i];
      if  not __equal(a[f.prop], b[f.prop], f.typ) then
        return false;
      end
    end
    return true;
  elseif sw == __kindInterface then 
    return __interfaceIsEqual(a, b);
  else
    return a == b;
  end
end;

__interfaceIsEqual = function(b)
  if a == __ifaceNil  or  b == __ifaceNil then
    return a == b;
  end
  if a.constructor ~= b.constructor then
    return false;
  end
  if a.constructor == __jsObjectPtr then
    return a.object == b.object;
  end
  if  not a.constructor.comparable then
    __throwRuntimeError("comparing uncomparable type "  ..  a.constructor.__str);
  end
  return __equal(a.__val, b.__val, a.constructor);
end;


__assertType = function(value, typ, returnTuple)

   local isInterface = (typ.kind == __kindInterface)
   local ok
   local missingMethod = "";
   if value == __ifaceNil then
      ok = false;
   elseif  not isInterface then
      ok = value.constructor == typ;
   else
      local valueTypeString = value.constructor.__str;
      ok = typ.implementedBy[valueTypeString];
      if ok == nil then
         ok = true;
         local valueMethodSet = __methodSet(value.constructor);
         local interfaceMethods = typ.methods;
         for i = 0,#interfaceMethods-1 do
            
            local tm = interfaceMethods[i];
            local found = false;
            for j = 0,#valueMethodSet-1 do
               
               local vm = valueMethodSet[j];
               if vm.name == tm.name  and  vm.pkg == tm.pkg  and  vm.typ == tm.typ then
                  found = true;
                  break;
               end
            end
            if  not found then
               ok = false;
               typ.missingMethodFor[valueTypeString] = tm.name;
               break;
            end
         end
         typ.implementedBy[valueTypeString] = ok;
      end
      if  not ok then
         missingMethod = typ.missingMethodFor[valueTypeString];
      end
   end
   
   if  not ok then
      if returnTuple then
         return {typ.zero(), false};
      end
      local msg = ""
      if value ~= __ifaceNil then
         msg = value.constructor.__str
      end
      --__panic(__packages["runtime"].TypeAssertionError.ptr("", msg, typ.__str, missingMethod));
      error("type-assertion-error: could not '"..msg.."' -> '"..typ.__str.."', missing method '"..missingMethod.."'")
   end
   
   if  not isInterface then
      value = value.__val;
   end
   if typ == __jsObjectPtr then
      value = value.object;
   end
   if returnTuple then
      return {value, true}
   end
   return value;
end;

__stackDepthOffset = 0;
__getStackDepth = function()
  local err = new Error();
  if err.stack == nil then
    return nil;
  end
  return __stackDepthOffset + #err.stack.split("\n");
end;
