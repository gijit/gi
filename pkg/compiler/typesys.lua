-- dofile '../math.lua' -- for __max, __min, __truncateToInt

-- NB fin and fin_test use ___ triple underscores, to
-- avoid collision while integrating with struct.lua

-- translation of javascript builtin 'prototype' -> typ.methodSet
--                                   'constructor' -> typ.___constructor

-- dofile '../int64.lua'
___ffi = require("ffi")
___bit = require("bit")

___global ={};
___module ={};
___packages = {}
___idCounter = 0;

function ___ipairsZeroCheck(arr)
   if arr[0] ~= nil then error("ipairs will miss the [0] index of this array") end
end

___mod = function(y) return x % y; end;
___parseInt = parseInt;
___parseFloat = function(f)
  if f ~= nil  and  f ~= nil  and  f.constructor == Number then
    return f;
  end
  return parseFloat(f);
end;

--[[
 ___froundBuf = Float32Array(1);
___fround = Math.fround  or  function(f)
  ___froundBuf[0] = f;
  return ___froundBuf[0];
end;
--]]

--[[
___imul = Math.imul  or  function(b)
   local ah = ___bit.band(___bit.rshift(a, 16), 0xffff);
   local al = ___bit.band(a, 0xffff);
   local bh = ___bit.band(___bit.rshift(b, 16), 0xffff);
   local bl = ___bit.band(b, 0xffff);
   return ((al * bl) + ___bit.arshift((___bit.rshift(___bit.lshift(ah * bl + al * bh), 16), 0), 0);
end;
--]]

___floatKey = function(f)
  if f ~= f then
     ___idCounter=___idCounter+1;
    return "NaN___" .. tostring(___idCounter);
  end
  return tostring(f);
end;

___flatten64 = function(x)
  return x.___high * 4294967296 + x.___low;
end;


___Infinity = math.huge

-- returned by ___basicValue2kind(v) on unrecognized kind.
___kindUnknown = -1;

___kindBool = 1;
___kindInt = 2;
___kindInt8 = 3;
___kindInt16 = 4;
___kindInt32 = 5;
___kindInt64 = 6;
___kindUint = 7;
___kindUint8 = 8;
___kindUint16 = 9;
___kindUint32 = 10;
___kindUint64 = 11;
___kindUintptr = 12;
___kindFloat32 = 13;
___kindFloat64 = 14;
___kindComplex64 = 15;
___kindComplex128 = 16;
___kindArray = 17;
___kindChan = 18;
___kindFunc = 19;
___kindInterface = 20;
___kindMap = 21;
___kindPtr = 22;
___kindSlice = 23;
___kindString = 24;
___kindStruct = 25;
___kindUnsafePointer = 26;

-- jea: sanity check my assumption by comparing
-- length with #a
function ___assertIsArray(a)
   local n = 0
   for k,v in pairs(a) do
      n=n+1
   end
   if #a ~= n then
      error("not an array, ___assertIsArray failed")
   end
end


-- length of array, counting [0] if present.
function ___lenz(array)      
   local n = #array
   if array[0] ~= nil then
      n=n+1
   end
   return n
end

-- st or showtable, a debug print helper.
-- seen avoids infinite looping on self-recursive types.
function ___st(t, name, indent, quiet, methods_desc, seen)
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
         --vals = ___st(v,"",indent+1,quiet,methods_desc, seen)
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
      s = s .. "\n"..___st(mt, "mt.of."..name, indent+1, true, methods_desc, seen)
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
function ___mapAndJoinStrings(splice, arr, fun)
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
___keys = function(m)
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
___flushConsole = function() end;
___throwRuntimeError = function(...) error(...) end
___throwNilPointerError = function()  ___throwRuntimeError("invalid memory address or nil pointer dereference"); end;
___call = function(fn, rcvr, args)  return fn(rcvr, args); end;
___makeFunc = function(fn)
   return function()
      -- TODO: port this!
      print("jea TODO: port this, what is ___externalize doing???")
      error("NOT DONE: port this!")
      --return ___externalize(fn(this, (___sliceType({},___jsObjectPtr))(___global.Array.prototype.slice.call(arguments, {}))), ___emptyInterface);
   end;
end;
___unused = function(v) end;

--
___mapArray = function(arr, f)
   local newarr = {}
   -- handle a zero argument, if present.
   local bump = 0
   local zval = arr[0]
   if zval ~= nil then
      bump = 1
      newarr[1] = fun(zval)
   end
   ___ipairsZeroCheck(arr)
   for i,v in ipairs(arr) do
      newarr[i+bump] = fun(v)
   end
   return newarr
end;

___methodVal = function(recv, name) 
  local vals = recv.___methodVals  or  {};
  recv.___methodVals = vals; -- /* noop for primitives */
  local f = vals[name];
  if f ~= nil then
    return f;
  end
  local method = recv[name];
  f = function() 
     ___stackDepthOffset = ___stackDepthOffset-1;
     -- try
     local res = {pcall(function()
           return recv[method](arguments);
     end)}
        -- finally
     ___stackDepthOffset=___stackDepthOffset+1;
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

___methodExpr = function(typ, name) 
   local method = typ.methodSet[name];
   if method.___expr == nil then
      method.___expr = function() 
         ___stackDepthOffset=___stackDepthOffset-1;

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
         ___stackDepthOffset=___stackDepthOffset+1;
         -- no catch, so rethrow any exception
         if not ok then
            error(threw)
         end
         return table.remove(res, 1)
      end;
   end
   return method.___expr;
end;

___ifaceMethodExprs = {};
___ifaceMethodExpr = function(name) 
  local expr = ___ifaceMethodExprs["_"  ..  name];
  if expr == nil then
     expr = function()
        ___stackDepthOffset = ___stackDepthOffset-1;
        -- try
        local res = {pcall(
                        function()
                           return Function.call.apply(arguments[0][name], arguments);
        end)}
        -- finally
        ___stackDepthOffset = ___stackDepthOffset+1;
        -- no catch
        local ok, threw = unpack(res)
        if not ok then
           error(threw)
        else
           -- non-panic return from pcall
           return table.remove(res, 1)
        end   
     end;
     ___ifaceMethodExprs["_"  ..  name] = expr
  end
  return expr;
end;

--

___subslice = function(slice, low, high, max)
   if high == nil then
      
   end
   if low < 0  or  (high ~= nil and high < low)  or  (max ~= nil and high ~= nil and max < high)  or  (high ~= nil and high > slice.___capacity)  or  (max ~= nil and max > slice.___capacity) then
      ___throwRuntimeError("slice bounds out of range");
   end
   
   local s = {}
   slice.___constructor.tfun(s, slice.___array);
   s.___offset = slice.___offset + low;
   s.___length = slice.___length - low;
   s.___capacity = slice.___capacity - low;
   if high ~= nil then
      s.___length = high - low;
   end
   if max ~= nil then
      s.___capacity = max - low;
   end
   return s;
end;

___copySlice = function(dst, src)
   local n = __min(src.___length, dst.___length);
   ___copyArray(dst.___array, src.___array, dst.___offset, src.___offset, n, dst.___constructor.elem);
   return n;
end;

--

___copyArray = function(dst, src, dstOffset, srcOffset, n, elem)
   --print("___copyArray called with n = ", n, " dstOffset=", dstOffset, " srcOffset=", srcOffset)
   --print("___copyArray has dst:")
   --___st(dst)
   --print("___copyArray has src:")
   --___st(src)
   
   n = tonumber(n)
   if n == 0  or  (dst == src  and  dstOffset == srcOffset) then
      return;
   end

   local sw = elem.kind
   if sw == ___kindArray or sw == ___kindStruct then
      
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

--
___clone = function(src, typ)
  local clone = typ.zero();
  typ.copy(clone, src);
  return clone;
end;

___pointerOfStructConversion = function(obj, typ)
  if(obj.___proxies == nil) then
    obj.___proxies = {};
    obj.___proxies[obj.constructor.___str] = obj;
  end
  local proxy = obj.___proxies[typ.___str];
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
     
    proxy = Object.create(typ.methodSet, properties);
    proxy.___val = proxy;
    obj.___proxies[typ.___str] = proxy;
    proxy.___proxies = obj.___proxies;
  end
  return proxy;
end;

--


___append = function(...)
   local arguments = {...}
   local slice = arguments[1]
   return ___internalAppend(slice, arguments, 1, #arguments - 1);
end;

___appendSlice = function(slice, toAppend)
   if slice == nil then 
      error("error calling ___appendSlice: slice must be available")
   end
   if toAppend == nil then
      error("error calling ___appendSlice: toAppend must be available")      
   end
   if type(toAppend) == "string" then
      local bytes = ___stringToBytes(toAppend);
      return ___internalAppend(slice, bytes, 0, #bytes);
   end
   return ___internalAppend(slice, toAppend.___array, toAppend.___offset, toAppend.___length);
end;

___internalAppend = function(slice, array, offset, length)
   if length == 0 then
      return slice;
   end

   local newArray = slice.___array;
   local newOffset = slice.___offset;
   local newLength = slice.___length + length;
   --print("jea debug: ___internalAppend: newLength is "..tostring(newLength))
   local newCapacity = slice.___capacity;
   local elem = slice.___constructor.elem;

   if newLength > newCapacity then
      newOffset = 0;
      local tmpCap
      if slice.___capacity < 1024 then
         tmpCap = slice.___capacity * 2
      else
         tmpCap = __truncateToInt(slice.___capacity * 5 / 4)
      end
      newCapacity = __max(newLength, tmpCap);

      newArray = {}
      local w = slice.___offset
      for i = 0,slice.___length do
         newArray[i] = slice.___array[i + w]
      end
      for i = #slice,newCapacity-1 do
         newArray[i] = elem.zero();
      end
      
   end

   --print("jea debug, ___internalAppend, newOffset = ", newOffset, " and slice.___length=", slice.___length)

   ___copyArray(newArray, array, newOffset + slice.___length, offset, length, elem);
   --print("jea debug, ___internalAppend, after copying over array:")
   --___st(newArray)

   local newSlice ={}
   slice.___constructor.tfun(newSlice, newArray);
   newSlice.___offset = newOffset;
   newSlice.___length = newLength;
   newSlice.___capacity = newCapacity;
   return newSlice;
end;

--

___substring = function(str, low, high)
  if low < 0  or  high < low  or  high > #str then
    ___throwRuntimeError("string slice bounds out of range");
  end
  return string.sub(str, low+1, high); -- high is inclusive, so no +1 needed.
end;

___sliceToArray = function(slice)
   local cp = {}
   if slice.___length > 0 then
      local k = 0
      for i = slice.___offset, slice.___offset + slice.___length -1 do
         cp[k] = slice.array[i]
         k=k+1
      end
   end
   return cp
end;

--


--

___valueBasicMT = {
   __name = "___valueBasicMT",
   __tostring = function(self, ...)
      --print("__tostring called from ___valueBasicMT")
      if type(self.___val) == "string" then
         return '"'..self.___val..'"'
      end
      if self ~= nil and self.___val ~= nil then
         --print("___valueBasicMT.__tostring called, with self.___val set.")
         if self.___val == self then
            -- not a basic value, but a pointer, array, slice, or struct.
            return "<this.___val == this; avoid inf loop>"
         end
         --return tostring(self.___val)
      end
      if getmetatable(self.___val) == ___valueBasicMT then
         --print("avoid infinite loop")
         return "<avoid inf loop>"
      else
         return tostring(self.___val)
      end
   end
}

___valueArrayMT = {
   __name = "___valueArrayMT",
   
   __newindex = function(t, k, v)
      --print("___valueArrayMT.__newindex called, t is:")
      --___st(t)

      if k < 0 or k >= #t then
         error "read of array error: access out-of-bounds"
      end
      
      t.___val[k] = v
   end,
   
   __index = function(t, k)
      --print("___valueArrayMT.__index called, k='"..tostring(k).."'; t.___val is:")
      --___st(t.___val)
      if k < 0 or k >= #t then
         error("write to array error: access out-of-bounds; "..tostring(k).." is outside [0, "  .. tostring(#t) .. ")")
      end
      
      return t.___val[k]
   end,

   __len = function(t)
      return int(___lenz(t.___val))
   end,
   
   __tostring = function(self, ...)
      --print("__tostring called from ___valueArrayMT")
      if type(self.___val) == "string" then
         return '"'..self.___val..'"'
      end
      if self ~= nil and self.___val ~= nil then
         --print("___valueArrayMT.__tostring called, with self.___val set.")
         if self.___val == self then
            -- not a basic value, but a pointer, array, slice, or struct.
            return "<this.___val == this; avoid inf loop>"
         end
         --return tostring(self.___val)
      end
      if getmetatable(self.___val) == ___valueArrayMT then
         --print("avoid infinite loop")
         return "<avoid inf loop>"
      else
         return tostring(self.___val)
      end
   end
}

___valueSliceMT = {
   __name = "___valueSliceMT",
   
   __newindex = function(t, k, v)
      --print("___valueSliceMT.__newindex called, t is:")
      --___st(t)
      local w = t.___offset + k
      if k < 0 or k >= t.___capacity then
         error "slice error: write out-of-bounds"
      end
      t.___array[w] = v
   end,
   
   __index = function(t, k)
      --print("___valueSliceMT.__index called, k='"..tostring(k).."'; t.___val is:")
      --___st(t.___val)
      local w = t.___offset + k
      if k < 0 or k >= t.___capacity then
         error "slice error: access out-of-bounds"
      end
      return t.___array[w]
   end,

   __len = function(t)
      return t.___length
   end,
   
   __tostring = function(self, ...)
      --print("__tostring called from ___valueSliceMT")
      if type(self.___val) == "string" then
         return '"'..self.___val..'"'
      end
      if self ~= nil and self.___val ~= nil then
         --print("___valueSliceMT.__tostring called, with self.___val set.")
         if self.___val == self then
            -- not a basic value, but a pointer, array, slice, or struct.
            return "<this.___val == this; avoid inf loop>"
         end
         --return tostring(self.___val)
      end
      if getmetatable(self.___val) == ___valueSliceMT then
         --print("avoid infinite loop")
         return "<avoid inf loop>"
      else
         return tostring(self.___val)
      end
   end
}


___tfunBasicMT = {
   __name = "___tfunBasicMT",
   __call = function(self, ...)
      --print("jea debug: ___tfunBasicMT.__call() invoked") -- , self='"..tostring(self).."' with tfun = ".. tostring(self.tfun).. " and args=")
      --print(debug.traceback())
      
      --print("in ___tfunBasicMT, start ___st on ...")
      --___st({...}, "___tfunBasicMT.dots")
      --print("in ___tfunBasicMT,   end ___st on ...")

      --print("in ___tfunBasicMT, start ___st on self")
      --___st(self, "self")
      --print("in ___tfunBasicMT,   end ___st on self")

      local newInstance = {}
      setmetatable(newInstance, ___valueBasicMT)
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


function ___newAnyArrayValue(elem, len)
   local array = {}
   for i =0, len -1 do
      array[i]= elem.zero();
   end
   return array;
end


___methodSynthesizers = {};
___addMethodSynthesizer = function(f)
   if ___methodSynthesizers == nil then
      f();
      return;
   end
   table.insert(___methodSynthesizers, f);
end;

___synthesizeMethods = function()
   ___ipairsZeroCheck(___methodSynthesizers)
   for i,f in ipairs(___methodSynthesizers) do
      f();
   end
   ___methodSynthesizers = nil;
end;

___ifaceKeyFor = function(x)
  if x == ___ifaceNil then
    return 'nil';
  end
  local c = x.constructor;
  return c.___str .. '___' .. c.keyFor(x.___val);
end;

___identity = function(x) return x; end;

___typeIDCounter = 0;

___idKey = function(x)
   if x.___id == nil then
      ___idCounter=___idCounter+1;
      x.___id = ___idCounter;
   end
   return String(x.___id);
end;

___newType = function(size, kind, str, named, pkg, exported, constructor)
   local typ ={};
   setmetatable(typ, ___tfunBasicMT)

   if kind ==  ___kindBool or
      kind == ___kindInt or 
      kind == ___kindInt8 or 
      kind == ___kindInt16 or 
      kind == ___kindInt32 or 
      kind == ___kindInt64 or 
      kind == ___kindUint or 
      kind == ___kindUint8 or 
      kind == ___kindUint16 or 
      kind == ___kindUint32 or 
      kind == ___kindUint64 or 
      kind == ___kindUintptr or 
   kind == ___kindUnsafePointer then

      -- jea: I observe that
      -- primitives have: this.___val ~= v; and are the types are
      -- distinguished with typ.wrapped = true; versus
      -- all table based values, that have: this.___val == this;
      -- and no .wrapped field.
      --
      typ.tfun = function(this, v) this.___val = v; end;
      typ.wrapped = true;
      typ.keyFor = ___identity;

   elseif kind == ___kindString then
      
      typ.tfun = function(this, v)
         --print("strings' tfun called! with v='"..tostring(v).."' and this:")
         --___st(this)
         this.___val = v; end;
      typ.wrapped = true;
      typ.keyFor = function(x) return "_" .. x; end;

   elseif kind == ___kindFloat32 or
   kind == ___kindFloat64 then
      
      typ.tfun = function(this, v) this.___val = v; end;
      typ.wrapped = true;
      typ.keyFor = function(x) return ___floatKey(x); end;


  elseif kind ==  ___kindComplex64 then 
    typ.tfun = function(this, real, imag)
      this.___real = ___fround(real);
      this.___imag = ___fround(imag);
      this.___val = this;
    end;
    typ.keyFor = function(x) return x.___real .. "_" .. x.___imag; end;
    

  elseif kind ==  ___kindComplex128 then 
    typ.tfun = function(this, real, imag)
      this.___real = real;
      this.___imag = imag;
      this.___val = this;
    end;
    typ.keyFor = function(x) return x.___real .. "_" .. x.___imag; end;
    
      
   elseif kind ==  ___kindPtr then

      typ.tfun = constructor  or
         function(this, getter, setter, target)
            print("pointer typ.tfun which is same as constructor called! getter='"..tostring(getter).."'; setter='"..tostring(setter).."; target = '"..tostring(target).."'")
            this.___get = getter;
            this.___set = setter;
            this.___target = target;
            this.___val = this; -- seems to indicate a non-primitive value.
         end;
      typ.keyFor = ___idKey;
      typ.init = function(elem)
         typ.elem = elem;
         typ.wrapped = (elem.kind == ___kindArray);
         typ.___nil = typ(___throwNilPointerError, ___throwNilPointerError);
      end;

   elseif kind ==  ___kindSlice then
      
      typ.tfun = function(this, array)
         this.___array = array;
         this.___offset = 0;
         this.___length = ___lenz(array);
         this.___capacity = this.___length;
         --print("jea debug: slice tfun set ___length to ", this.___length)
         --print("jea debug: slice tfun set ___capacity to ", this.___capacity)
         --print("jea debug: slice tfun sees array: ")
         --for i,v in pairs(array) do
         --   print("array["..tostring(i).."] = ", v)
         --end
         
         this.___val = this;
         this.___constructor = typ
         setmetatable(this, ___valueSliceMT)
      end;
      typ.init = function(elem)
         typ.elem = elem;
         typ.comparable = false;
         typ.___nil = typ({},{});
      end;
      
   elseif kind ==  ___kindArray then

      typ.tfun = function(this, v)
         --print("in tfun ctor function for ___kindArray")
         this.___val = v;
         setmetatable(this, ___valueArrayMT)
      end;
      typ.wrapped = true;
      typ.ptr = ___newType(4, ___kindPtr, "*" .. str, false, "", false, function(this, array)
                             this.___get = function() return array; end;
                             this.___set = function(v) typ.copy(this, v); end;
                             this.___val = array;
      end);
      typ.init = function(elem, len)
         typ.elem = elem;
         typ.len = len;
         typ.comparable = elem.comparable;
         typ.keyFor = function(x)
            return ___mapAndJoinStrings("_", x, function(e)
                                          return string.gsub(tostring(elem.keyFor(e)), "\\", "\\\\")
            end)
         end
         typ.copy = function(dst, src)
            ___copyArray(dst, src, 0, 0, #src, elem);
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
         -- js: Object.defineProperty(typ.ptr.___nil, "nilCheck", { get= ___throwNilPointerError end);
      end;
      -- end ___kindArray

   
  elseif kind ==  ___kindChan then
     
    typ.tfun = function(this, v) this.___val = v; end;
    typ.wrapped = true;
    typ.keyFor = ___idKey;
    typ.init = function(elem, sendOnly, recvOnly)
      typ.elem = elem;
      typ.sendOnly = sendOnly;
      typ.recvOnly = recvOnly;
    end;
    

  elseif kind ==  ___kindFunc then 

     typ.tfun = function(this, v) this.___val = v; end;
     typ.wrapped = true;
     typ.init = function(params, results, variadic)
        typ.params = params;
        typ.results = results;
        typ.variadic = variadic;
        typ.comparable = false;
     end;
    

  elseif kind ==  ___kindInterface then 

     typ = { implementedBy= {}, missingMethodFor= {} };
     typ.keyFor = ___ifaceKeyFor;
     typ.init = function(methods)
        typ.methods = methods;
        for _, m in pairs(methods) do
           -- TODO:
           -- jea: why this? seems it would end up being a huge set?
           ___ifaceNil[m.prop] = ___throwNilPointerError;
        end;
     end;
     
     
   elseif kind ==  ___kindMap then 
      
      typ.tfun = function(this, v) this.___val = v; end;
      typ.wrapped = true;
      typ.init = function(key, elem)
         typ.key = key;
         typ.elem = elem;
         typ.comparable = false;
      end;
      
   elseif kind ==  ___kindStruct then
      
      typ.tfun = function(this, v) this.___val = v; end;
      typ.wrapped = true;

      -- the typ.methodSet will be the
      -- metatable for instances of the struct; this is
      -- equivalent to the prototype in js.
      --
      typ.methodSet = {___name="methodSet for "..str, ___typ = typ}
      typ.methodSet.__index = typ.methodSet
      
      local ctor = function(this, ...)
         this.___get = function() return this; end;
         this.___set = function(v) typ.copy(this, v); end;
         if constructor ~= nil then
            constructor(this, ...);
         end
         setmetatable(this, typ.ptr.methodSet)
      end
      typ.ptr = ___newType(4, ___kindPtr, "*" .. str, false, pkg, exported, ctor);
      -- ___newType sets typ.comparable = true

      -- pointers have their own method sets, but *T can call elem methods in Go.
      typ.ptr.elem = typ;
      typ.ptr.methodSet = {___name="methodSet for "..typ.ptr.___str, ___typ = typ.ptr}
      typ.ptr.methodSet.__index = typ.ptr.methodSet

      -- ___kindStruct.init is here:
      typ.init = function(pkgPath, fields)
         print("top of init() for struct, fields=")
         for i, f in pairs(fields) do
            ___st(f, "field #"..tostring(i))
            ___st(f.typ, "typ of field #"..tostring(i))
         end
         
         typ.pkgPath = pkgPath;
         typ.fields = fields;
         ___ipairsZeroCheck(fields)
         for i,f in ipairs(fields) do
            if not f.typ.comparable then
               typ.comparable = false;
               break;
            end
         end
         typ.keyFor = function(x)
            local val = x.___val;
            return ___mapAndJoinStrings("_", fields, function(f)
                                          return string.gsub(tostring(f.typ.keyFor(val[f.prop])), "\\", "\\\\")
            end)
         end;
         typ.copy = function(dst, src)
            print("top of typ.copy for structs, here is dst then src:")
            ___st(dst, "dst")
            ___st(src, "src")
            print("fields:")
            ___st(fields,"fields")
            ___ipairsZeroCheck(fields)
            for _, f in ipairs(fields) do
               local sw2 = f.typ.kind
               
               if sw2 == ___kindArray or
               sw2 ==  ___kindStruct then 
                  f.typ.copy(dst[f.prop], src[f.prop]);
               else
                  dst[f.prop] = src[f.prop];
               end
            end
         end;
         print("jea debug: on ___kindStruct: set .copy on typ to .copy=", typ.copy)
         -- /* nil value */
         local properties = {};
         ___ipairsZeroCheck(fields)
         for i,f in ipairs(fields) do
            properties[f.prop] = { get= ___throwNilPointerError, set= ___throwNilPointerError };
         end;
         typ.ptr.___nil = {} -- Object.create(constructor.prototype,s properties);
         --if constructor ~= nil then
         --   constructor(typ.ptr.___nil)
         --end
         typ.ptr.___nil.___val = typ.ptr.___nil;
         -- /* methods for embedded fields */
         ___addMethodSynthesizer(function()
               local synthesizeMethod = function(target, m, f)
                  if target.methodSet[m.prop] ~= nil then return; end
                  target.methodSet[m.prop] = function()
                     local v = this.___val[f.prop];
                     if f.typ == ___jsObjectPtr then
                        v = ___jsObjectPtr(v);
                     end
                     if v.___val == nil then
                        local w = {}
                        f.typ(w, v);
                        v = w
                     end
                     return v[m.prop](v, arguments);
                  end;
               end;
               for i,f in ipairs(fields) do
                  if f.anonymous then
                     for _, m in ipairs(___methodSet(f.typ)) do
                        synthesizeMethod(typ, m, f);
                        synthesizeMethod(typ.ptr, m, f);
                     end;
                     for _, m in ipairs(___methodSet(___ptrType(f.typ))) do
                        synthesizeMethod(typ.ptr, m, f);
                     end;
                  end
               end;
         end);
      end;
      
   else
      error("invalid kind: " .. tostring(kind));
   end
   
   -- set zero() method
   if kind == ___kindBool or
   kind ==___kindMap then
      typ.zero = function() return false; end;

   elseif kind == ___kindInt or
      kind ==  ___kindInt8 or
      kind ==  ___kindInt16 or
      kind ==  ___kindInt32 or
   kind ==  ___kindInt64 then
      typ.zero = function() return 0LL; end;
      
   elseif kind ==  ___kindUint or
      kind ==  ___kindUint8  or
      kind ==  ___kindUint16 or
      kind ==  ___kindUint32 or
      kind ==  ___kindUint64 or
      kind ==  ___kindUintptr or
   kind ==  ___kindUnsafePointer then
      typ.zero = function() return 0ULL; end;

   elseif   kind ==  ___kindFloat32 or
   kind ==  ___kindFloat64 then
      typ.zero = function() return 0; end;
      
   elseif kind ==  ___kindString then
      typ.zero = function() return ""; end;

   elseif kind == ___kindComplex64 or
   kind == ___kindComplex128 then
      local zero = typ(0, 0);
      typ.zero = function() return zero; end;
      
   elseif kind == ___kindPtr or
   kind == ___kindSlice then
      
      typ.zero = function() return typ.___nil; end;
      
   elseif kind == ___kindChan then
      typ.zero = function() return ___chanNil; end;
   
   elseif kind == ___kindFunc then
      typ.zero = function() return ___throwNilPointerError; end;
      
   elseif kind == ___kindInterface then
      typ.zero = function() return ___ifaceNil; end;
      
   elseif kind == ___kindArray then
      
      typ.zero = function()
         return ___newAnyArrayValue(typ.elem, typ.len)
      end;

   elseif kind == ___kindStruct then
      typ.zero = function()
         return typ.ptr();
      end;

   else
      error("invalid kind: " .. tostring(kind))
   end

   typ.id = ___typeIDCounter;
   ___typeIDCounter=___typeIDCounter+1;
   typ.size = size;
   typ.kind = kind;
   typ.___str = str;
   typ.named = named;
   typ.pkg = pkg;
   typ.exported = exported;
   typ.methods = typ.methods or {};
   typ.methodSetCache = nil;
   typ.comparable = true;
   return typ;
   
end

function ___methodSet(typ)
   
  --if typ.methodSetCache ~= nil then
  --return typ.methodSetCache;
  --end
  local base = {};

  local isPtr = (typ.kind == ___kindPtr);
  print("___methodSet called with typ=")
  ___st(typ)
  print("___methodSet sees isPtr=", isPtr)
  
  if isPtr  and  typ.elem.kind == ___kindInterface then
     -- jea: I assume this is because pointers to interfaces don't themselves have methods.
     typ.methodSetCache = {};
     return {};
  end

  local myTyp = typ
  if isPtr then
     myTyp = typ.elem
  end
  local current = {{typ= myTyp, indirect= isPtr}};

  -- the Go spec says:
  -- The method set of the corresponding pointer type *T is
  -- the set of all methods declared with receiver *T or T
  -- (that is, it also contains the method set of T).
  
  local seen = {};

  print("top of while, #current is", #current)
  while #current > 0 do
     local next = {};
     local mset = {};
     
     for _,e in pairs(current) do
        if seen[e.typ.___str] then
           break
        end
        seen[e.typ.___str] = true;
        
       if e.typ.named then
          for _, mthod in pairs(e.typ.methods) do
             print("adding to mset, mthod = ", mthod)
             table.insert(mset, mthod);
          end
          if e.indirect then
             for _, mthod in pairs(___ptrType(e.typ).methods) do
                print("adding to mset, mthod = ", mthod)
                table.insert(mset, mthod)
             end
          end
       end
       
       -- switch e.typ.kind
       local knd = e.typ.kind
       
       if knd == ___kindStruct then
          
          -- assume that e.typ.fields must be an array!
          -- TODO: remove this assert after confirmation.
          ___assertIsArray(e.typ.fields)
          ___ipairsZeroCheck(e.typ.fields)
          for i,f in ipairs(e.typ.fields) do
             if f.anonymous then
                local fTyp = f.typ;
                local fIsPtr = (fTyp.kind == ___kindPtr);
                local ty 
                if fIsPtr then
                   ty = fTyp.elem
                else
                   ty = fTyp
                end
                table.insert(next, {typ=ty, indirect= e.indirect or fIsPtr});
             end;
          end;
          
          
       elseif knd == ___kindInterface then
          
          for _, mthod in pairs(e.typ.methods) do
             print("adding to mset, mthod = ", mthod)
             table.insert(mset, mthod)
          end
       end
     end;

     -- above may have made duplicates, now dedup
     print("at dedup, #mset = " .. tostring(#mset))
     for _, m in pairs(mset) do
        if base[m.name] == nil then
           base[m.name] = m;
        end
     end;
     print("after dedup, base for typ '"..typ.___str.."' is ")
     ___st(base)
     
     current = next;
  end
  
  typ.methodSetCache = {};
  table.sort(base)
  for _, detail in pairs(base) do
     table.insert(typ.methodSetCache, detail)
  end;
  return typ.methodSetCache;
end;


___Bool          = ___newType( 1, ___kindBool,          "bool",           true, "", false, nil);
___Int           = ___newType( 8, ___kindInt,           "int",            true, "", false, nil);
___Int8          = ___newType( 1, ___kindInt8,          "int8",           true, "", false, nil);
___Int16         = ___newType( 2, ___kindInt16,         "int16",          true, "", false, nil);
___Int32         = ___newType( 4, ___kindInt32,         "int32",          true, "", false, nil);
___Int64         = ___newType( 8, ___kindInt64,         "int64",          true, "", false, nil);
___Uint          = ___newType( 8, ___kindUint,          "uint",           true, "", false, nil);
___Uint8         = ___newType( 1, ___kindUint8,         "uint8",          true, "", false, nil);
___Uint16        = ___newType( 2, ___kindUint16,        "uint16",         true, "", false, nil);
___Uint32        = ___newType( 4, ___kindUint32,        "uint32",         true, "", false, nil);
___Uint64        = ___newType( 8, ___kindUint64,        "uint64",         true, "", false, nil);
___Uintptr       = ___newType( 8, ___kindUintptr,       "uintptr",        true, "", false, nil);
___Float32       = ___newType( 8, ___kindFloat32,       "float32",        true, "", false, nil);
___Float64       = ___newType( 8, ___kindFloat64,       "float64",        true, "", false, nil);
--___Complex64     = ___newType( 8, ___kindComplex64,     "complex64",      true, "", false, nil);
--___Complex128    = ___newType(16, ___kindComplex128,    "complex128",     true, "", false, nil);
___String        = ___newType(16, ___kindString,        "string",         true, "", false, nil);
--___UnsafePointer = ___newType( 8, ___kindUnsafePointer, "unsafe.Pointer", true, "", false, nil);

--[[

___nativeArray = function(elemKind)

   if false then
      if elemKind ==  ___kindInt then 
         return Int32Array; -- in js, a builtin typed array
      elseif elemKind ==  ___kindInt8 then 
         return Int8Array;
      elseif elemKind ==  ___kindInt16 then 
         return Int16Array;
      elseif elemKind ==  ___kindInt32 then 
         return Int32Array;
      elseif elemKind ==  ___kindUint then 
         return Uint32Array;
      elseif elemKind ==  ___kindUint8 then 
         return Uint8Array;
      elseif elemKind ==  ___kindUint16 then 
         return Uint16Array;
      elseif elemKind ==  ___kindUint32 then 
         return Uint32Array;
      elseif elemKind ==  ___kindUintptr then 
         return Uint32Array;
      elseif elemKind ==  ___kindFloat32 then 
         return Float32Array;
      elseif elemKind ==  ___kindFloat64 then 
         return Float64Array;
      else
         return Array;
      end
   end
end;

___toNativeArray = function(elemKind, array)
  local nativeArray = ___nativeArray(elemKind);
  if nativeArray == Array {
    return array;
  end
  return nativeArray(array); -- new
end;

--]]


___ptrType = function(elem)
   local typ = elem.ptr;
   if typ == nil then
      typ = ___newType(4, ___kindPtr, "*" .. elem.___str, false, "", elem.exported, nil);
      elem.ptr = typ;
      typ.init(elem);
   end
   return typ;
end;

___newDataPointer = function(data, constructor)
   if constructor.elem.kind == ___kindStruct then
      return data;
   end
   return constructor(function() return data; end, function(v) data = v; end);
end;

___indexPtr = function(array, index, constructor)
   array.___ptr = array.___ptr  or  {};
   local a = array.___ptr[index]
   if a ~= nil then
      return a
   end
   a = constructor(function() return array[index]; end, function(v) array[index] = v; end);
   array.___ptr[index] = a
   return a
end;


___arrayTypes = {};
___arrayType = function(elem, len)
   local typeKey = elem.id .. "_" .. len;
   local typ = ___arrayTypes[typeKey];
   if typ == nil then
      typ = ___newType(24, ___kindArray, "[" .. len .. "]" .. elem.___str, false, "", false, nil);
      ___arrayTypes[typeKey] = typ;
      typ.init(elem, len);
   end
   return typ;
end;


___chanType = function(elem, sendOnly, recvOnly)
   
   local str
   local field
   if recvOnly then
      str = "<-chan " .. elem.___str
      field = "RecvChan"
   elseif sendOnly then
      str = "chan<- " .. elem.___str
      field = "SendChan"
   else
      str = "chan " .. elem.___str
      field = "Chan"
   end
   local typ = elem[field];
   if typ == nil then
      typ = ___newType(4, ___kindChan, str, false, "", false, nil);
      elem[field] = typ;
      typ.init(elem, sendOnly, recvOnly);
   end
   return typ;
end;

function ___Chan(elem, capacity)
   local this = {}
   if capacity < 0  or  capacity > 2147483647 then
      ___throwRuntimeError("makechan: size out of range");
   end
   this.elem = elem;
   this.___capacity = capacity;
   this.___buffer = {};
   this.___sendQueue = {};
   this.___recvQueue = {};
   this.___closed = false;
   return this
end;
___chanNil = ___Chan(nil, 0);
___chanNil.___recvQueue = { length= 0, push= function()end, shift= function() return nil; end, indexOf= function() return -1; end; };
___chanNil.___sendQueue = ___chanNil.___recvQueue


___funcTypes = {};
___funcType = function(params, results, variadic)
   local typeKey = ___mapAndJoinStrings(",", params, function(p) return p.id; end) .. "_" .. ___mapAndJoinStrings(",", results, function(r) return r.id; end) .. "_" .. tostring(variadic);
  local typ = ___funcTypes[typeKey];
  if typ == nil then
    local paramTypes = ___mapArray(params, function(p) return p.___str; end);
    if variadic then
      paramTypes[#paramTypes - 1] = "..." .. paramTypes[#paramTypes - 1].substr(2);
    end
    local str = "func(" .. table.concat(paramTypes, ", ") .. ")";
    if #results == 1 then
      str = str .. " " .. results[1].___str;
      end else if #results > 1 then
      str = str .. " (" .. ___mapAndJoinStrings(", ", results, function(r) return r.___str; end) .. ")";
    end
    typ = ___newType(4, ___kindFunc, str, false, "", false, nil);
    ___funcTypes[typeKey] = typ;
    typ.init(params, results, variadic);
  end
  return typ;
end;

--- interface types here

function ___interfaceStrHelper(m)
   local s = ""
   if m.pkg ~= "" then
      s = m.pkg .. "."
   end
   return s .. m.name .. string.sub(m.typ.___str, 6) -- sub for removing "___kind"
end

___interfaceTypes = {};
___interfaceType = function(methods)
   
   local typeKey = ___mapAndJoinStrings("_", methods, function(m)
                                          return m.pkg .. "," .. m.name .. "," .. m.typ.id;
   end)
   local typ = ___interfaceTypes[typeKey];
   if typ == nil then
      local str = "interface {}";
      if #methods ~= 0 then
         str = "interface { " .. ___mapAndJoinStrings("; ", methods, ___interfaceStrHelper) .. " }"
      end
      typ = ___newType(8, ___kindInterface, str, false, "", false, nil);
      ___interfaceTypes[typeKey] = typ;
      typ.init(methods);
   end
   return typ;
end;
___emptyInterface = ___interfaceType({});
___ifaceNil = {};
___error = ___newType(8, ___kindInterface, "error", true, "", false, nil);
___error.init({{prop= "Error", name= "Error", pkg= "", typ= ___funcType({}, {___String}, false) }});

___mapTypes = {};
___mapType = function(key, elem)
  local typeKey = key.id .. "_" .. elem.id;
  local typ = ___mapTypes[typeKey];
  if typ == nil then
    typ = ___newType(8, ___kindMap, "map[" .. key.___str .. "]" .. elem.___str, false, "", false, nil);
    ___mapTypes[typeKey] = typ;
    typ.init(key, elem);
  end
  return typ;
end;
___makeMap = function(keyForFunc, entries)
   local m = {};
   for i =0,#entries-1 do
    local e = entries[i];
    m[keyForFunc(e.k)] = e;
  end
  return m;
end;


-- ___basicValue2kind: identify type of basic value
--   or return ___kindUnknown if we don't recognize it.
function ___basicValue2kind(v)

   local ty = type(v)
   if ty == "cdata" then
      local cty = ___ffi.typeof(v)
      if cty == int64 then
         return ___kindInt
      elseif cty == int8 then
         return ___kindInt8
      elseif cty == int16 then
         return ___kindInt16
      elseif cty == int32 then
         return ___kindInt32
      elseif cty == int64 then
         return ___kindInt64
      elseif cty == uint then
         return ___kindUint
      elseif cty == uint8 then
         return ___kindUint8
      elseif cty == uint16 then
         return ___kindUint16
      elseif cty == uint32 then
         return ___kindUint32
      elseif cty == uint64 then
         return ___kindUint64
      elseif cty == float32 then
         return ___kindFloat32
      elseif cty == float64 then
         return ___kindFloat64         
      else
         return ___kindUnknown;
         --error("___basicValue2kind: unhandled cdata cty: '"..tostring(cty).."'")
      end      
   elseif ty == "boolean" then
      return ___kindBool;
   elseif ty == "number" then
      return ___kindFloat64
   elseif ty == "string" then
      return ___kindString
   end
   
   return ___kindUnknown;
   --error("___basicValue2kind: unhandled ty: '"..ty.."'")   
end

___sliceType = function(elem)
   local typ = elem.slice;
   if typ == nil then
      typ = ___newType(24, ___kindSlice, "[]" .. elem.___str, false, "", false, nil);
      elem.slice = typ;
      typ.init(elem);
   end
   return typ;
end;

___makeSlice = function(typ, length, capacity)
   length = tonumber(length)
   if capacity == nil then
      capacity = length
   else
      capacity = tonumber(capacity)
   end
   if length < 0  or  length > 9007199254740992 then -- 2^53
      ___throwRuntimeError("makeslice: len out of range");
   end
   if capacity < 0  or  capacity < length  or  capacity > 9007199254740992 then
      ___throwRuntimeError("makeslice: cap out of range");
   end
   local array = ___newAnyArrayValue(typ.elem, capacity)
   local slice = typ(array);
   slice.___length = length;
   return slice;
end;




function field2strHelper(f)
   local tag = ""
   if f.tag ~= "" then
      tag = string.gsub(f.tag, "\\", "\\\\")
      tag = string.gsub(tag, "\"", "\\\"")
   end
   return f.name .. " " .. f.typ.___str .. tag
end

function typeKeyHelper(f)
   return f.name .. "," .. f.typ.id .. "," .. f.tag;
end

___structTypes = {};
___structType = function(pkgPath, fields)
   local typeKey = ___mapAndJoinStrings("_", fields, typeKeyHelper)

   local typ = ___structTypes[typeKey];
   if typ == nil then
      local str
      if #fields == 0 then
         str = "struct {}";
      else
         str = "struct { " .. ___mapAndJoinStrings("; ", fields, field2strHelper) .. " }";
      end
      
      typ = ___newType(0, ___kindStruct, str, false, "", false, function()
                         local this = {}
                         this.___val = this;
                         for i = 0, #fields-1 do
                            local f = fields[i];
                            local arg = arguments[i];
                            if arg ~= nil then
                               this[f.prop] = arg
                            else
                               this[f.prop] = f.typ.zero();
                            end
                         end
                         return this
      end);
      ___structTypes[typeKey] = typ;
      typ.init(pkgPath, fields);
   end
   return typ;
end;


___equal = function(a, b, typ)
  if typ == ___jsObjectPtr then
    return a == b;
  end

  local sw = typ.kind
  if sw == ___kindComplex64 or
  sw == ___kindComplex128 then
     return a.___real == b.___real  and  a.___imag == b.___imag;
     
  elseif sw == ___kindInt64 or
         sw == ___kindUint64 then 
     return a.___high == b.___high  and  a.___low == b.___low;
     
  elseif sw == ___kindArray then 
    if #a ~= #b then
      return false;
    end
    for i=0,#a-1 do
      if  not ___equal(a[i], b[i], typ.elem) then
        return false;
      end
    end
    return true;
    
  elseif sw == ___kindStruct then
     
    for i = 0,#(typ.fields)-1 do
      local f = typ.fields[i];
      if  not ___equal(a[f.prop], b[f.prop], f.typ) then
        return false;
      end
    end
    return true;
  elseif sw == ___kindInterface then 
    return ___interfaceIsEqual(a, b);
  else
    return a == b;
  end
end;

___interfaceIsEqual = function(a, b)
  if a == ___ifaceNil  or  b == ___ifaceNil then
    return a == b;
  end
  if a.constructor ~= b.constructor then
    return false;
  end
  if a.constructor == ___jsObjectPtr then
    return a.object == b.object;
  end
  if  not a.constructor.comparable then
    ___throwRuntimeError("comparing uncomparable type "  ..  a.constructor.___str);
  end
  return ___equal(a.___val, b.___val, a.constructor);
end;


___assertType = function(value, typ, returnTuple)

   local isInterface = (typ.kind == ___kindInterface)
   local ok
   local missingMethod = "";
   if value == ___ifaceNil then
      ok = false;
   elseif  not isInterface then
      ok = value.___typ == typ;
   else
      local valueTypeString = value.___typ.___str;

      -- this caching doesn't get updated as methods
      -- are added, so disable it until fixed, possibly, in the future.
      --ok = typ.implementedBy[valueTypeString];
      ok = nil
      if ok == nil then
         ok = true;
         local valueMethodSet = ___methodSet(value.___typ);
         local interfaceMethods = typ.methods;
         print("valueMethodSet is")
         ___st(valueMethodSet)
         print("interfaceMethods is")
         ___st(interfaceMethods)

         ___ipairsZeroCheck(interfaceMethods)
         ___ipairsZeroCheck(valueMethodSet)
         for _, tm in ipairs(interfaceMethods) do            
            local found = false;
            for _, vm in ipairs(valueMethodSet) do
               print("checking vm against tm, where tm=")
               ___st(tm)
               print("and vm=")
               ___st(vm)
               
               if vm.name == tm.name  and  vm.pkg == tm.pkg  and  vm.typ == tm.typ then
                  print("match found against vm and tm.")
                  found = true;
                  break;
               end
            end
            if  not found then
               print("match *not* found for tm.name = '"..tm.name.."'")
               ok = false;
               typ.missingMethodFor[valueTypeString] = tm.name;
               break;
            end
         end
         typ.implementedBy[valueTypeString] = ok;
      end
      if not ok then
         missingMethod = typ.missingMethodFor[valueTypeString];
      end
   end
   print("___assertType: after matching loop, ok = ", ok)
   
   if not ok then
      if returnTuple then
         return typ.zero(), false
      end
      local msg = ""
      if value ~= ___ifaceNil then
         msg = value.___typ.___str
      end
      --___panic(___packages["runtime"].TypeAssertionError.ptr("", msg, typ.___str, missingMethod));
      error("type-assertion-error: could not '"..msg.."' -> '"..typ.___str.."', missing method '"..missingMethod.."'")
   end
   
   if not isInterface then
      value = value.___val;
   end
   if typ == ___jsObjectPtr then
      value = value.object;
   end
   if returnTuple then
      return value, true
   end
   return value
end;

___stackDepthOffset = 0;
___getStackDepth = function()
  local err = Error(); -- new
  if err.stack == nil then
    return nil;
  end
  return ___stackDepthOffset + #err.stack.split("\n");
end;

-- possible replacement for ipairs.
-- starts at a[0] if it is present.
function ___zipairs(a)
   local n = 0
   local s = #a
   if a[0] ~= nil then
      n = -1
   end
   return function()
      n = n + 1
      if n <= s then return n,a[n] end
   end
end
