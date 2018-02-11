
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
         --vals = __st(v,"",indent+1,quiet,methods_desc, seen)
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

__basicTypeMT = {
   __tostring = function(self, ...)
      if type(self.__val) == "string" then
         return '"'..self.__val..'"'
      end
      return tostring(self.__val)
   end
}


__tfunMT = {
   __name = "__tfunMT",
   __call = function(self, ...)
      --print("jea debug: __tfunMT.__call() invoked, self='"..tostring(self).."' with tfun = ".. tostring(self.tfun).. " and args=")
      
      --print("in __tfunMT, start __st on ...")
      --__st({...}, "__tfunMT.dots")
      --print("in __tfunMT,   end __st on ...")

      --print("in __tfunMT, start __st on self")
      --__st(self, "self")
      --print("in __tfunMT,   end __st on self")

      local newInstance = {}
      setmetatable(newInstance, __basicTypeMT)
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

__typeIDCounter = 0;

__newType = function(size, kind, str, named, pkg, exported, constructor)
  local typ ={};
  setmetatable(typ, __tfunMT)

  if kind ==  __kindBool or
  kind == __kindInt or 
  kind == __kindInt8 or 
  kind == __kindInt16 or 
  kind == __kindInt32 or 
  kind == __kindUint or 
  kind == __kindUint8 or 
  kind == __kindUint16 or 
  kind == __kindUint32 or 
  kind == __kindUintptr or 
  kind == __kindUnsafePointer then
     
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
__Int           = __newType( 4, __kindInt,           "int",            true, "", false, nil);
__Int8          = __newType( 1, __kindInt8,          "int8",           true, "", false, nil);
__Int16         = __newType( 2, __kindInt16,         "int16",          true, "", false, nil);
__Int32         = __newType( 4, __kindInt32,         "int32",          true, "", false, nil);
__Int64         = __newType( 8, __kindInt64,         "int64",          true, "", false, nil);
__Uint          = __newType( 4, __kindUint,          "uint",           true, "", false, nil);
__Uint8         = __newType( 1, __kindUint8,         "uint8",          true, "", false, nil);
__Uint16        = __newType( 2, __kindUint16,        "uint16",         true, "", false, nil);
__Uint32        = __newType( 4, __kindUint32,        "uint32",         true, "", false, nil);
__Uint64        = __newType( 8, __kindUint64,        "uint64",         true, "", false, nil);
__Uintptr       = __newType( 4, __kindUintptr,       "uintptr",        true, "", false, nil);
__Float32       = __newType( 4, __kindFloat32,       "float32",        true, "", false, nil);
__Float64       = __newType( 8, __kindFloat64,       "float64",        true, "", false, nil);
__Complex64     = __newType( 8, __kindComplex64,     "complex64",      true, "", false, nil);
__Complex128    = __newType(16, __kindComplex128,    "complex128",     true, "", false, nil);
__String        = __newType( 8, __kindString,        "string",         true, "", false, nil);
__UnsafePointer = __newType( 4, __kindUnsafePointer, "unsafe.Pointer", true, "", false, nil);
