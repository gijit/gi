
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


__tfunMT = {
   __name = "__tfunMT",
   __call = function(the_mt, self, ...)
      print("jea debug: __tfunMT.__call() invoked, self='",tostring(self),"', with tfun = ",self.tfun" and args=")
      
      print("in __tfunMT, start __st on ...")
      __st({...}, "__tfunMT.dots")
      print("in __tfunMT,   end __st on ...")

      print("in __tfunMT, start __st on self")
      __st(self, "self")
      print("in __tfunMT,   end __st on self")

      if self ~= nil and self.tfun ~= nil then
         print("calling tfun! -- let constructors set metatables if they wish to.")
         self.tfun({}, ...)
      else
         if self ~= nil then
            print("self.tfun was nil")
         end
      end
      return self
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
     
    typ.tfun = function(this, v) this.__val = v; end;
    typ.wrapped = true;
    typ.keyFor = function(x) return "_" .. x; end;

  elseif kind == __kindFloat32 or
  kind == __kindFloat64 then
       
       typ.tfun = function(this, v) this.__val = v; end;
       typ.wrapped = true;
       typ.keyFor = function(x) return __floatKey(x); end;
  end

  if kind == __kindBool or
  kind ==__kindMap then
    typ.zero = function() return false; end;

  elseif kind == __kindInt or
  kind ==  __kindInt8 or
  kind ==  __kindInt16 or
  kind ==  __kindInt32 or
  kind ==  __kindUint or
  kind ==  __kindUint8  or
  kind ==  __kindUint16 or
  kind ==  __kindUint32 or
  kind ==  __kindUintptr or
  kind ==  __kindUnsafePointer or
  kind ==  __kindFloat32 or
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
