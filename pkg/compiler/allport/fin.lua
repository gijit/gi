dofile '../math.lua' -- for __max, __min, __truncateToInt

dofile '/Users/jaten/go/src/github.com/gijit/gi/pkg/compiler/int64.lua'
__ffi = require("ffi")



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

__throwNilPointerError = function() error("invalid memory address or nil pointer dereference"); end

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

__valueBasicMT = {
   __name = "__valueBasicMT",
   __tostring = function(self, ...)
      --print("__tostring called from __valueBasicMT")
      if type(self.__val) == "string" then
         return '"'..self.__val..'"'
      end
      if self ~= nil and self.__val ~= nil then
         print("__valueBasicMT.__tostring called, with self.__val set.")
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


__tfunBasicMT = {
   __name = "__tfunBasicMT",
   __call = function(self, ...)
      print("jea debug: __tfunBasicMT.__call() invoked") -- , self='"..tostring(self).."' with tfun = ".. tostring(self.tfun).. " and args=")
      
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
            print("calling tfun! -- let constructors set metatables if they wish to.")

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

__typeIDCounter = 0;

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

  elseif kind ==  __kindPtr then
     
     typ.tfun = constructor  or
        function(this, getter, setter, target)
           print("pointer typ.tfun which is same as constructor called! getter='"..tostring(getter).."' setter='"..tostring(setter).."target = '"..tostring(target).."'")
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
   
  if n == 0  or  (dst == src  and  dstOffset == srcOffset) then
    return;
  end

  if src.subarray then
    dst.set(src.subarray(srcOffset, srcOffset + n), dstOffset);
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

