dofile '../math.lua' -- for __max, __min, __truncateToInt

__Infinity = math.huge

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
         print("calling tfun! -- let constructors set metatables on instances if they wish to.")
         self.tfun({}, ...)
      else
         if self ~= nil then
            print("self.tfun was nil")
         end
      end
      return self
end

--from the __newType for __kindStruct:
--typ.ptr.prototype.__get = function() return this; end;
--typ.ptr.prototype.__set = function(v) typ.copy(this, v); end;

}

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

function __new(ctor, ...)
   return ctor({}, ...)
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

__global ={};
__module ={};

-- TODO: jea, needed?
--if module ~= nil then
--  __module = module;
--end

__packages = {}
__idCounter = 0;

-- return sorted keys from table m. Assumes keys are strings.
__keys = function(m)
   if type(m) ~= "table" then
      return {}
   end
   local r = {}
   for k in pairs(m) do
      local tyk = type(k)
      if tyk ~= "string" then
         error "__keys() assumption broken: key was not string, but rather "..tyk
      end
      if tyk == "function" then
         k = tostring(k)
      end
      table.insert(r, k)
   end
   table.sort(r)
   return r
end

__flushConsole = function() end;
__throwRuntimeError = function(...) error(...) end
__throwNilPointerError = function()  __throwRuntimeError("invalid memory address or nil pointer dereference"); end;
__call = function(fn, rcvr, args)  return fn(rcvr, args); end;
__makeFunc = function(fn)  return function()  return __externalize(fn(this, (__sliceType({},__jsObjectPtr))(__global.Array.prototype.slice.call(arguments, {}))), __emptyInterface); end; end;
__unused = function(v) end;

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
        error err
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
            error threw
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
    expr = __ifaceMethodExprs["_"  ..  name] = function()
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
          error threw
       else
          -- non-panic return from pcall
          return table.remove(res, 1)
       end       
    end;
  end
  return expr;
end;

__subslice = function(x)
  if low < 0  or  high < low  or  max < high  or  high > slice.__capacity  or  max > slice.__capacity then
    __throwRuntimeError("slice bounds out of range");
  end
  local s = new slice.constructor(slice.__array);
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

__substring = function(h)
  if low < 0  or  high < low  or  high > #str then
    __throwRuntimeError("slice bounds out of range");
  end
  return str.substring(low, high);
end;

__sliceToArray = function(slice)
  if slice.__array.constructor ~= Array then
    return slice.__array.subarray(slice.__offset, slice.__offset + slice.__length);
  end
  return slice.__array.slice(slice.__offset, slice.__offset + slice.__length);
end;

__decodeRune = function(s)
  local c0 = str.charCodeAt(pos);

  if c0 < 0x80 then
    return {c0, 1};
  end

  if c0 ~= c0  or  c0 < 0xC0 then
    return {0xFFFD, 1};
  end

  local c1 = str.charCodeAt(pos + 1);
  if c1 ~= c1  or  c1 < 0x80  or  0xC0 <= c1 then
    return {0xFFFD, 1};
  end

  if c0 < 0xE0 then
    local r = (c0 & 0x1F) << 6 | (c1 & 0x3F);
    if r <= 0x7F then
      return {0xFFFD, 1};
    end
    return {r, 2};
  end

  local c2 = str.charCodeAt(pos + 2);
  if c2 ~= c2  or  c2 < 0x80  or  0xC0 <= c2 then
    return {0xFFFD, 1};
  end

  if c0 < 0xF0 then
    local r = (c0 & 0x0F) << 12 | (c1 & 0x3F) << 6 | (c2 & 0x3F);
    if r <= 0x7FF then
      return {0xFFFD, 1};
    end
    if 0xD800 <= r  and  r <= 0xDFFF then
      return {0xFFFD, 1};
    end
    return {r, 3};
  end

  local c3 = str.charCodeAt(pos + 3);
  if c3 ~= c3  or  c3 < 0x80  or  0xC0 <= c3 then
    return {0xFFFD, 1};
  end

  if c0 < 0xF8 then
    local r = (c0 & 0x07) << 18 | (c1 & 0x3F) << 12 | (c2 & 0x3F) << 6 | (c3 & 0x3F);
    if r <= 0xFFFF  or  0x10FFFF < r then
      return {0xFFFD, 1};
    end
    return {r, 4};
  end

  return {0xFFFD, 1};
end;

__encodeRune = function(r)
  if r < 0  or  r > 0x10FFFF  or  (0xD800 <= r  and  r <= 0xDFFF) then
    r = 0xFFFD;
  end
  if r <= 0x7F then
    return String.fromCharCode(r);
  end
  if r <= 0x7FF then
    return String.fromCharCode(0xC0 | r >> 6, 0x80 | (r & 0x3F));
  end
  if r <= 0xFFFF then
    return String.fromCharCode(0xE0 | r >> 12, 0x80 | (r >> 6 & 0x3F), 0x80 | (r & 0x3F));
  end
  return String.fromCharCode(0xF0 | r >> 18, 0x80 | (r >> 12 & 0x3F), 0x80 | (r >> 6 & 0x3F), 0x80 | (r & 0x3F));
end;

__stringToBytes = function(r)
  local array = new Uint8Array(#str);
  for i = 0,#str-1 do
    array[i] = str.charCodeAt(i);
  end
  return array;
end;

__bytesToString = function(e)
  if #slice == 0 then
    return "";
  end
  local str = "";
  for i = 0,#slice-1,10000 do
    str = str .. String.fromCharCode.apply(nil, slice.__array.subarray(slice.__offset + i, slice.__offset + __min(slice.__length, i + 10000)));
  end
  return str;
end;

__stringToRunes = function(r)
  local array = new Int32Array(#str);
  local rune, j = 0;
  local i = 0
  local n = #str
  while true do
     if i >= n then
        break
     end
     
     rune = __decodeRune(str, i);
     array[j] = rune[1];
     
     i = i + rune[2]
     j = j + 1
  end
  -- in js, a subarray is like a slice, a view on a shared ArrayBuffer.
  return array.subarray(0, j);
end;

__runesToString = function(e)
  if slice.__length == 0 then
    return "";
  end
  local str = "";
  for i = 0,#slice-1 do
    str = str .. __encodeRune(slice.__array[slice.__offset + i]);
  end
  return str;
end;

__copyString = function(c)
  local n = __min(#src, dst.__length);
  for i = 0,n-1 do
    dst.__array[dst.__offset + i] = src.charCodeAt(i);
  end
  return n;
end;

__copySlice = function(c)
  local n = __min(src.__length, dst.__length);
  __copyArray(dst.__array, src.__array, dst.__offset, src.__offset, n, dst.constructor.elem);
  return n;
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

__append = function(slice)
  return __internalAppend(slice, arguments, 1, #arguments - 1);
end;

__appendSlice = function(slice, toAppend)
  if toAppend.constructor == String then
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
  local newCapacity = slice.__capacity;

  if newLength > newCapacity then
     local newOffset = 0;
     local tmpCap
     if slice.__capacity < 1024 then
        tmpCap = slice.__capacity * 2
     else
        tmpCap = __truncateToInt(slice.__capacity * 5 / 4)
     end
     newCapacity = __max(newLength, tmpCap);

    if slice.__array.constructor == Array then
       newArray = slice.__array.slice(slice.__offset, slice.__offset + slice.__length);
       #newArray = newCapacity;
       local zero = slice.constructor.elem.zero;
       for i = #slice,newCapacity-1 do
          newArray[i] = zero();
       end
  else 
          newArray = new slice.__array.constructor(newCapacity);
          newArray.set(slice.__array.subarray(slice.__offset, slice.__offset + slice.__length));
  end
  end

  __copyArray(newArray, array, newOffset + slice.__length, offset, length, slice.constructor.elem);

  local newSlice = new slice.constructor(newArray);
  newSlice.__offset = newOffset;
  newSlice.__length = newLength;
  newSlice.__capacity = newCapacity;
  return newSlice;
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

__mod = function(y) return x % y; end;
__parseInt = parseInt;
__parseFloat = function(f)
  if f ~= nil  and  f ~= nil  and  f.constructor == Number then
    return f;
  end
  return parseFloat(f);
end;

__froundBuf = new Float32Array(1);
__fround = Math.fround  or  function(f)
  __froundBuf[0] = f;
  return __froundBuf[0];
end;

__imul = Math.imul  or  function(b)
  local ah = (a >>> 16) & 0xffff;
  local al = a & 0xffff;
  local bh = (b >>> 16) & 0xffff;
  local bl = b & 0xffff;
  return ((al * bl) + (((ah * bl + al * bh) << 16) >>> 0) >> 0);
end;

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

__shiftLeft64 = function(y)
  if y == 0 then
    return x;
  end
  if y < 32 then
    return new x.constructor(x.__high << y | x.__low >>> (32 - y), (x.__low << y) >>> 0);
  end
  if y < 64 then
    return new x.constructor(x.__low << (y - 32), 0);
  end
  return new x.constructor(0, 0);
end;

__shiftRightInt64 = function(y)
  if y == 0 then
    return x;
  end
  if y < 32 then
    return new x.constructor(x.__high >> y, (x.__low >>> y | x.__high << (32 - y)) >>> 0);
  end
  if y < 64 then
    return new x.constructor(x.__high >> 31, (x.__high >> (y - 32)) >>> 0);
  end
  if x.__high < 0 then
    return new x.constructor(-1, 4294967295);
  end
  return new x.constructor(0, 0);
end;

__shiftRightUint64 = function(y)
  if y == 0 then
    return x;
  end
  if y < 32 then
    return new x.constructor(x.__high >>> y, (x.__low >>> y | x.__high << (32 - y)) >>> 0);
  end
  if y < 64 then
    return new x.constructor(0, x.__high >>> (y - 32));
  end
  return new x.constructor(0, 0);
end;

__mul64 = function(y)
  local high = 0, low = 0;
  if (y.__low & 1) ~= 0 then
    high = x.__high;
    low = x.__low;
  end
  for i = 1,31 do
    if (y.__low & 1<<i) ~= 0 then
      high = high + x.__high << i | x.__low >>> (32 - i);
      low = low + (x.__low << i) >>> 0;
    end
  end
  for i = 0,31 do
    if (y.__high & 1<<i) ~= 0 then
      high = high + x.__low << i;
    end
  end
  return new x.constructor(high, low);
end;

__div64 = function(r)
  if y.__high == 0  and  y.__low == 0 then
    __throwRuntimeError("integer divide by zero");
  end

  local s = 1;
  local rs = 1;

  local xHigh = x.__high;
  local xLow = x.__low;
  if xHigh < 0 then
    s = -1;
    rs = -1;
    xHigh = -xHigh;
    if xLow ~= 0 then
      xHigh=xHigh-1;
      xLow = 4294967296 - xLow;
    end
  end

  local yHigh = y.__high;
  local yLow = y.__low;
  if y.__high < 0 then
    s *= -1;
    yHigh = -yHigh;
    if yLow ~= 0 then
       yHigh=yHigh-1;
      yLow = 4294967296 - yLow;
    end
  end

  local high = 0, low = 0, n = 0;
  while (yHigh < 2147483648  and  ((xHigh > yHigh)  or  (xHigh == yHigh  and  xLow > yLow))) {
    yHigh = (yHigh << 1 | yLow >>> 31) >>> 0;
    yLow = (yLow << 1) >>> 0;
    n=n+1;
end
  for i = 0, n do
    high = high << 1 | low >>> 31;
    low = (low << 1) >>> 0;
    if (xHigh > yHigh)  or  (xHigh == yHigh  and  xLow >= yLow) then
      xHigh = xHigh - yHigh;
      xLow = xLow - yLow;
      if xLow < 0 then
         xHigh=xHigh-1;
        xLow=xLow +4294967296;
      end
      low=low+1;
      if low == 4294967296 then
         high=high+1;
         low = 0;
      end
    end
    yLow = (yLow >>> 1 | yHigh << (32 - 1)) >>> 0;
    yHigh = yHigh >>> 1;
  end

  if returnRemainder then
    return new x.constructor(xHigh * rs, xLow * rs);
  end
  return new x.constructor(high * s, low * s);
end;

__divComplex = function(d)
  local ninf = n.__real == __Infinity  or  n.__real == -__Infinity  or  n.__imag == __Infinity  or  n.__imag == -__Infinity;
  local dinf = d.__real == __Infinity  or  d.__real == -__Infinity  or  d.__imag == __Infinity  or  d.__imag == -__Infinity;
  local nnan =  not ninf  and  (n.__real ~= n.__real  or  n.__imag ~= n.__imag);
  local dnan =  not dinf  and  (d.__real ~= d.__real  or  d.__imag ~= d.__imag);
  if(nnan  or  dnan) then
    return new n.constructor(NaN, NaN);
  end
  if ninf  and   not dinf then
    return new n.constructor(__Infinity, __Infinity);
  end
  if  not ninf  and  dinf then
    return new n.constructor(0, 0);
  end
  if d.__real == 0  and  d.__imag == 0 then
    if n.__real == 0  and  n.__imag == 0 then
      return new n.constructor(NaN, NaN);
    end
    return new n.constructor(__Infinity, __Infinity);
  end
  local a = Math.abs(d.__real);
  local b = Math.abs(d.__imag);
  if a <= b then
    local ratio = d.__real / d.__imag;
    local denom = d.__real * ratio + d.__imag;
    return new n.constructor((n.__real * ratio + n.__imag) / denom, (n.__imag * ratio - n.__real) / denom);
  end
  local ratio = d.__imag / d.__real;
  local denom = d.__imag * ratio + d.__real;
  return new n.constructor((n.__imag * ratio + n.__real) / denom, (n.__imag - n.__real * ratio) / denom);
end;

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

__methodSynthesizers = {};
__addMethodSynthesizer = function(f)
  if __methodSynthesizers == nil then
    f();
    return;
  end
  __methodSynthesizers.push(f);
end;
__synthesizeMethods = function()
  __methodSynthesizers.forEach(function(f) f(); end);
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
       

  elseif kind ==  __kindInt64 then 
    typ.tfun = function(this, high, low)
      this.__high = (high + Math.floor(Math.ceil(low) / 4294967296)) >> 0;
      this.__low = low >>> 0;
      this.__val = this;
    end;
    typ.keyFor = function(x) return x.__high .. "_" .. x.__low; end;
    

  elseif kind ==  __kindUint64 then 
    typ.tfun = function(this, high, low)
      this.__high = (high + Math.floor(Math.ceil(low) / 4294967296)) >>> 0;
      this.__low = low >>> 0;
      this.__val = this;
    end;
    typ.keyFor = function(x) return x.__high .. "_" .. x.__low; end;
    

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
    

elseif kind ==  __kindArray then

   typ.tfun = function(this, v) this.__val = v; end;
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
          local ma = __mapArray(x, function(e)
               return tostring(elem.keyFor(e))
            end)
            return table.concat(ma, "_")
         end

      --typ.keyFor = function(x)
        --return Array.prototype.join.call(__mapArray(x, function(e)
        --  return tostring(elem.keyFor(e)).replace(/\\/g, "\\\\").replace(/\__/g, "\\__");
        --end), "_");
      --end;

      typ.copy = function(dst, src)
        __copyArray(dst, src, 0, 0, #src, elem);
      end;
      typ.ptr.init(typ);
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
   
  elseif kind ==  __kindChan then
     
    typ.tfun = function(this, v) this.__val = v; end;
    typ.wrapped = true;
    typ.keyFor = __idKey;
    typ.init = function(elem, sendOnly, recvOnly)
      typ.elem = elem;
      typ.sendOnly = sendOnly;
      typ.recvOnly = recvOnly;
    end;
    

  elseif kind ==  __kindFunc then 
     typ.tfun = function(this, v) this.__val = v; end;
     typ.wrapped = true;
     typ.init = function(params, results, variadic)
        typ.params = params;
        typ.results = results;
        typ.variadic = variadic;
        typ.comparable = false;
     end;
    

  elseif kind ==  __kindInterface then 
     typ = { implementedBy= {}, missingMethodFor= {} };
    typ.keyFor = __ifaceKeyFor;
    typ.init = function(methods)
      typ.methods = methods;
      methods.forEach(function(m)
        __ifaceNil[m.prop] = __throwNilPointerError;
      end);
    end;
    

  elseif kind ==  __kindMap then 
    typ.tfun = function(this, v) this.__val = v; end;
    typ.wrapped = true;
    typ.init = function(key, elem)
      typ.key = key;
      typ.elem = elem;
      typ.comparable = false;
    end;
    

  elseif kind ==  __kindPtr then

    typ = constructor  or  function(getter, setter, target)
      this.__get = getter;
      this.__set = setter;
      this.__target = target;
      this.__val = this;
    end;
    typ.keyFor = __idKey;
    typ.init = function(elem)
      typ.elem = elem;
      typ.wrapped = (elem.kind == __kindArray);
      typ.__nil = new typ(__throwNilPointerError, __throwNilPointerError);
    end;
    

  elseif kind ==  __kindSlice then
    typ.tfun = function(this, array)
      if array.constructor ~= typ.nativeArray then
         array = new typ.nativeArray(array);
      end
      this.__array = array;
      this.__offset = 0;
      this.__length = #array;
      this.__capacity = #array;
      this.__val = this;
    end;
    typ.init = function(elem)
      typ.elem = elem;
      typ.comparable = false;
      typ.nativeArray = __nativeArray(elem.kind);
      typ.__nil = new typ({});
    end;
    

    elseif kind ==  __kindStruct then
       
    typ.tfun = function(this, v) this.__val = v; end;
    typ.wrapped = true;
    typ.ptr = __newType(4, __kindPtr, "*" .. str, false, pkg, exported, constructor);
    typ.ptr.elem = typ;
    typ.ptr.prototype.__get = function() return this; end;
    typ.ptr.prototype.__set = function(v) typ.copy(this, v); end;
    typ.init = function(pkgPath, fields)
      typ.pkgPath = pkgPath;
      typ.fields = fields;
      fields.forEach(function(f)
        if  not f.typ.comparable then
          typ.comparable = false;
        end
      end);
      typ.keyFor = function(x)
        local val = x.__val;
        return __mapArray(fields, function(f)
          return tostring(f.typ.keyFor(val[f.prop])).replace(/\\/g, "\\\\").replace(/\__/g, "\\__");
        end).join("_");
      end;
      typ.copy = function(dst, src)
         for i=0,#fields-1 do
            local f = fields[i];
            local sw2 = f.typ.kind
            
            if sw2 == __kindArray or
            sw2 ==  __kindStruct then 
               f.typ.copy(dst[f.prop], src[f.prop]);
               continue;
            else
               dst[f.prop] = src[f.prop];
               continue;
          end
        end
      end;
      -- /* nil value */
      local properties = {};
      fields.forEach(function(f)
            properties[f.prop] = { get= __throwNilPointerError, set= __throwNilPointerError };
      end);
      typ.ptr.__nil = Object.create(constructor.prototype, properties);
      typ.ptr.__nil.__val = typ.ptr.__nil;
      -- /* methods for embedded fields */
      __addMethodSynthesizer(function()
        local synthesizeMethod = function(target, m, f)
          if target.prototype[m.prop] ~= nil then return; end
          target.prototype[m.prop] = function()
            local v = this.__val[f.prop];
            if f.typ == __jsObjectPtr then
              v = new __jsObjectPtr(v);
            end
            if v.__val == nil then
              v = new f.typ(v);
            end
            return v[m.prop].apply(v, arguments);
          end;
        end;
        fields.forEach(function(f)
          if f.anonymous then
            __methodSet(f.typ).forEach(function(m)
              synthesizeMethod(typ, m, f);
              synthesizeMethod(typ.ptr, m, f);
            end);
            __methodSet(__ptrType(f.typ)).forEach(function(m)
              synthesizeMethod(typ.ptr, m, f);
            end);
          end
        end);
      end);
    end;
    

  else
     __panic("invalid kind: " .. tostring(kind));
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
    
  elseif k ==  __kindInt64 or
  kind == __kindUint64 or
  kind == __kindComplex64 or
  kind == __kindComplex128 then
    local zero = new typ(0, 0);
    typ.zero = function() return zero; end;
    

  elseif kind == __kindPtr or
  kind == __kindSlice then
    typ.zero = function() return typ.__nil; end;
    

    elseif kind == __kindChan:
    typ.zero = function() return __chanNil; end;
    

  elseif kind == __kindFunc then
    typ.zero = function() return __throwNilPointerError; end;
    

  elseif kind == __kindInterface then
    typ.zero = function() return __ifaceNil; end;
    

  elseif kind == __kindArray then
    typ.zero = function()
       --local arrayClass = __nativeArray(typ.elem.kind);
       --if arrayClass ~= Array then
       --   return new arrayClass(typ.len);
       --end
       --local array = new Array(typ.len);
       local array = {} -- new Array(typ.len);
       for i =0, typ.len -1 do
          table.insert(array, typ.elem.zero());
       end
       return array;
    end;
    

  elseif kind == __kindStruct then
    typ.zero = function() return new typ.ptr(); end;
    

  else
    __panic(__String("invalid kind: " .. tostring(kind)));
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
end;

__methodSet = function(typ)
  if typ.methodSetCache ~= nil then
    return typ.methodSetCache;
  end
  local base = {};

  local isPtr = (typ.kind == __kindPtr);
  if isPtr  and  typ.elem.kind == __kindInterface then
    typ.methodSetCache = {};
    return {};
  end

  local current = [{typ= isPtr ? typ.elem : typ, indirect= isPtrend];

  local seen = {};

  while (#current > 0) {
    local next = {};
    local mset = {};

    current.forEach(function(e)
      if seen[e.typ.__str] then
        return;
      end
      seen[e.typ.__str] = true;

      if e.typ.named then
        mset = mset.concat(e.typ.methods);
        if e.indirect then
          mset = mset.concat(__ptrType(e.typ).methods);
        end
      end

      -- switch e.typ.kind
      local knd = e.typ.kind
      if knd == __kindStruct then

         -- assume that e.typ.fields must be an array!
         __assertIsArray(e.typ.fields)
         for i,f in ipairs(e.typ.fields) do
            if f.anonymous then
               local fTyp = f.typ;
               local fIsPtr = (fTyp.kind == __kindPtr);
               local ty 
               if fIsPtr then
                  ty = fTyp.elem
               else
                  ty = fTyp
               end
               next.push({typ=ty, indirect= e.indirect  or  fIsPtr});
            end;
         end;
         

      elseif knd == __kindInterface then
        mset = mset.concat(e.typ.methods);
        
      end
    end);

    mset.forEach(function(m)
      if base[m.name] == nil then
        base[m.name] = m;
      end
    end);

    current = next;
  end

  typ.methodSetCache = {};
  Object.keys(base).sort().forEach(function(name)
    typ.methodSetCache.push(base[name]);
  end);
  return typ.methodSetCache;
end;

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

__nativeArray = function(elemKind)
  
  if elemKind ==  __kindInt then 
    return Int32Array;
  elseif elemKind ==  __kindInt8 then 
    return Int8Array;
  elseif elemKind ==  __kindInt16 then 
    return Int16Array;
  elseif elemKind ==  __kindInt32 then 
    return Int32Array;
  elseif elemKind ==  __kindUint then 
    return Uint32Array;
  elseif elemKind ==  __kindUint8 then 
    return Uint8Array;
  elseif elemKind ==  __kindUint16 then 
    return Uint16Array;
  elseif elemKind ==  __kindUint32 then 
    return Uint32Array;
  elseif elemKind ==  __kindUintptr then 
    return Uint32Array;
  elseif elemKind ==  __kindFloat32 then 
    return Float32Array;
  elseif elemKind ==  __kindFloat64 then 
    return Float64Array;
  else
    return Array;
  end
end;

__toNativeArray = function(elemKind, array)
  local nativeArray = __nativeArray(elemKind);
  if nativeArray == Array {
    return array;
  end
  return new nativeArray(array);
end;
__arrayTypes = {};
__arrayType = function(elem, len)
  local typeKey = elem.id .. "_" .. len;
  local typ = __arrayTypes[typeKey];
  if typ == nil then
     typ = __newType(24, __kindArray, "[" .. tostring(len) .. "]" .. elem.__str, false, "", false, nil);
    __arrayTypes[typeKey] = typ;
    typ.init(elem, len);
  end
  return typ;
end;

__chanType = function(elem, sendOnly, recvOnly)

   local str
   local field
   if recvOnly then
      str = "<-chan " .. elem.__str
      field = "RecvChan"
   elseif sendOnly then
      str = "chan<- " .. elem.__str
      field = "SendChan"
   else
      str = "chan " .. elem.__str
      field = "Chan"
   end
  local typ = elem[field];
  if typ == nil then
    typ = __newType(4, __kindChan, str, false, "", false, nil);
    elem[field] = typ;
    typ.init(elem, sendOnly, recvOnly);
  end
  return typ;
end;
__Chan = function(this, elem, capacity)
  if capacity < 0  or  capacity > 2147483647 then
     __throwRuntimeError("makechan: size out of range");
  end
  this.elem = elem;
  this.__capacity = capacity;
  this.__buffer = {};
  this.__sendQueue = {};
  this.__recvQueue = {};
  this.__closed = false;
end;
__chanNil = new __Chan(nil, 0);
__chanNil.__sendQueue = __chanNil.__recvQueue = { length= 0, push= function()end, shift= function() return nil; end, indexOf= function() return -1; end; };

__funcTypes = {};
__funcType = function(params, results, variadic)
  local typeKey = __mapAndJoinStrings(",", params, function(p) return p.id; end) .. "_" .. __mapAndJoinStrings(",", results, function(r) return r.id; end) .. "_" .. variadic;
  local typ = __funcTypes[typeKey];
  if typ == nil then
    local paramTypes = __mapArray(params, function(p) return p.__str; end);
    if variadic then
      paramTypes[#paramTypes - 1] = "..." .. paramTypes[#paramTypes - 1].substr(2);
    end
    local string = "func(" .. table.concat(paramTypes, ", ") .. ")";
    if #results == 1 then
      str = str .. " " .. results[0].__str;
      end else if #results > 1 then
      str = str .. " (" .. __mapAndJoinStrings(", ", results, function(r) return r.__str; end) .. ")";
    end
    typ = __newType(4, __kindFunc, str, false, "", false, nil);
    __funcTypes[typeKey] = typ;
    typ.init(params, results, variadic);
  end
  return typ;
end;

function __interfaceStrHelper(m)
   local s = ""
   if m.pkg ~= "" then
      s = m.pkg .. "."
   end
   return s .. m.name .. string.sub(m.typ.__str, 6) -- sub for removing "__kind"
end   

__interfaceTypes = {};
__interfaceType = function(methods)
  local typeKey = __mapArray(methods, function(m) return m.pkg .. "," .. m.name .. "," .. m.typ.id; end).join("_");
  local typ = __interfaceTypes[typeKey];
  if typ == nil then
    local str = "interface {}";
    if #methods ~= 0 then
       str = "interface { " .. __mapAndJoinStrings("; ", methods, __interfaceStrHelper) .. " }"
    end
    typ = __newType(8, __kindInterface, str, false, "", false, nil);
    __interfaceTypes[typeKey] = typ;
    typ.init(methods);
  end
  return typ;
end;
__emptyInterface = __interfaceType({});
__ifaceNil = {};
__error = __newType(8, __kindInterface, "error", true, "", false, nil);
__error.init([{prop= "Error", name= "Error", pkg= "", typ= __funcType({}, [__String], false)end]);

__mapTypes = {};
__mapType = function(key, elem)
  local typeKey = key.id .. "_" .. elem.id;
  local typ = __mapTypes[typeKey];
  if typ == nil then
    typ = __newType(8, __kindMap, "map[" .. key.__str .. "]" .. elem.__str, false, "", false, nil);
    __mapTypes[typeKey] = typ;
    typ.init(key, elem);
  end
  return typ;
end;
__makeMap = function(keyForFunc, entries)
   local m = {};
   for i =0,#entries-1 do
    local e = entries[i];
    m[keyForFunc(e.k)] = e;
  end
  return m;
end;

__ptrType = function(elem)
  local typ = elem.ptr;
  if typ == nil then
    typ = __newType(4, __kindPtr, "*" .. elem.__str, false, "", elem.exported, nil);
    elem.ptr = typ;
    typ.init(elem);
  end
  return typ;
end;

__newDataPointer = function(data, constructor)
  if constructor.elem.kind == __kindStruct then
    return data;
  end
  return new constructor(function() return data; end, function(v) data = v; end);
end;

__indexPtr = function(array, index, constructor)
  array.__ptr = array.__ptr  or  {};
  return array.__ptr[index]  or  (array.__ptr[index] = new constructor(function() return array[index]; end, function(v) array[index] = v; end));
end;

__sliceType = function(elem)
  local typ = elem.slice;
  if typ == nil then
    typ = __newType(12, __kindSlice, "[]" .. elem.__str, false, "", false, nil);
    elem.slice = typ;
    typ.init(elem);
  end
  return typ;
end;
  
__makeSlice = function(typ, length, capacity)
  capacity = capacity  or  length;
  if length < 0  or  length > 2147483647 then
    __throwRuntimeError("makeslice: len out of range");
  end
  if capacity < 0  or  capacity < length  or  capacity > 2147483647 then
    __throwRuntimeError("makeslice: cap out of range");
  end
  local array = new typ.nativeArray(capacity);
  if typ.nativeArray == Array then
     for i = 0,capacity-1 do
      array[i] = typ.elem.zero();
    end
  end
  local slice = new typ(array);
  slice.__length = length;
  return slice;
end;

function field2strHelper(f)
   local tag = ""
   if f.tag ~= "" then
   tag = tag
   -- TODO, translate from js, instead of the above tag = tag:
   -- seems to be escaping backslashes, then escaping double quotes.
   -- tag = " \"" .. f.tag.replace(/\\/g, "\\\\").replace(/"/g, "\\\"") .. "\""
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

__assertType = function(value, type, returnTuple)
  local isInterface = (typ.kind == __kindInterface), ok, missingMethod = "";
  if value == __ifaceNil then
    ok = false;
  end else if  not isInterface then
    ok = value.constructor == type;
  end else {
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
    __panic(new __packages["runtime"].TypeAssertionError.ptr("", (value == __ifaceNil ? "" : value.constructor.__str), typ.__str, missingMethod));
  end

  if  not isInterface then
    value = value.__val;
  end
  if type == __jsObjectPtr then
    value = value.object;
  end
  return returnTuple ? [value, true] : value;
end;

__stackDepthOffset = 0;
__getStackDepth = function()
  local err = new Error();
  if err.stack == nil then
    return nil;
  end
  return __stackDepthOffset + #err.stack.split("\n");
end;

__panicStackDepth = nil, __panicValue;
__callDeferred = function(deferred, jsErr, fromPanic) 
  if  not fromPanic  and  deferred ~= nil  and  deferred.index >= __curGoroutine.#deferStack then
    throw jsErr;
  end
  if jsErr ~= nil then
     local newErr = nil;

       -- try
     local res = {pcall(
                     function()
                        __curGoroutine.deferStack.push(deferred);
                        __panic(new __jsErrorPtr(jsErr));
     end)}
     
     -- catch
     local ok, err = unpack(res)
     if not ok then
        newErr = err;
     end

     __curGoroutine.deferStack.pop();
     __callDeferred(deferred, newErr);
     return;
  end
  if __curGoroutine.asleep then
     return;
  end
  
  __stackDepthOffset = __stackDepthOffset-1;
  local outerPanicStackDepth = __panicStackDepth;
  local outerPanicValue = __panicValue;

  local localPanicValue = __curGoroutine.panicStack.pop();
  if localPanicValue ~= nil then
    __panicStackDepth = __getStackDepth();
    __panicValue = localPanicValue;
  end

  -- try
  local res = {pcall(function()

    while (true) do 
      if deferred == nil then
        deferred = __curGoroutine.deferStack[__curGoroutine.#deferStack - 1];
        if deferred == nil then
          -- /* The panic reached the top of the stack. Clear it and throw it as a JavaScript error. */
          __panicStackDepth = nil;
          if localPanicValue.Object instanceof Error then
            error localPanicValue.Object;
          end
          local msg;
          if localPanicValue.constructor == __String then
            msg = localPanicValue.__val;
          end else if localPanicValue.Error ~= nil then
            msg = localPanicValue.Error();
          end else if localPanicValue.__str ~= nil then
            msg = localPanicValue.__str();
          end else then
            msg = localPanicValue;
          end
          error Error(msg);
        end
      end
      local call = deferred.pop();
      if call == nil then
        __curGoroutine.deferStack.pop();
        if localPanicValue ~= nil then
          deferred = nil;
          continue;
        end
        return;
      end
      local r = call[0].apply(call[2], call[1]);
      if r  and  r.__blk ~= nil then
        deferred.push([r.__blk, {}, r]);
        if fromPanic then
          throw nil;
        end
        return;
      end

      if localPanicValue ~= nil  and  __panicStackDepth == nil then
        throw nil; -- /* error was recovered */
        end
      end
end)}

    -- finally 
    if (localPanicValue ~= nil) {
       if (__panicStackDepth ~= nil) {
          __curGoroutine.panicStack.push(localPanicValue);
       end
       __panicStackDepth = outerPanicStackDepth;
       __panicValue = outerPanicValue;
       end
       __stackDepthOffset = __stackDepthOffset+1;
    end
   -- end finally
end;

__panic = function(value)
  __curGoroutine.panicStack.push(value);
  __callDeferred(nil, nil, true);
end;
__recover = function()
  if __panicStackDepth == nil  or  (__panicStackDepth ~= nil  and  __panicStackDepth ~= __getStackDepth() - 2) then
    return __ifaceNil;
  end
  __panicStackDepth = nil;
  return __panicValue;
end;
__throw = function(err) error err; end;

__noGoroutine = { asleep= false, exit= false, deferStack= {}, panicStack= {} };
__curGoroutine = __noGoroutine, __totalGoroutines = 0, __awakeGoroutines = 0, __checkForDeadlock = true;
__mainFinished = false;
__go = function(fun, args)
  __totalGoroutines=__totalGoroutines+1;
  __awakeGoroutines=__awakeGoroutines+1;
  local __goroutine = function()
     --try
     local ok, err = pcall(function()
                        __curGoroutine = __goroutine;
                        local r = fun(nil, args);
                        if r  and  r.__blk ~= nil then
                           fun = function() return r.__blk(); end;
                           args = {};
                           return;
                        end
                        __goroutine.exit = true;
     end)
     -- finally
     __curGoroutine = __noGoroutine;
     if __goroutine.exit then -- /* also set by runtime.Goexit() */
        __totalGoroutines=__totalGoroutines-1;
        __goroutine.asleep = true;
     end
     if __goroutine.asleep then
        __awakeGoroutines=__awakeGoroutines-1;
        if  not __mainFinished  and  __awakeGoroutines == 0  and  __checkForDeadlock then
           error("fatal error: all goroutines are asleep - deadlock!");
           if __global.process ~= nil then
              __global.process.exit(2);
           end
        end
     end
     -- catch(err)
     if not ok then
        if not __goroutine.exit then
           error err;
        end
     end
  end;
  __goroutine.asleep = false;
  __goroutine.exit = false;
  __goroutine.deferStack = {};
  __goroutine.panicStack = {};
  __schedule(__goroutine);
end;

__scheduled = {};
__runScheduled = function()
   -- try
   local res = {pcall(
                   function()
                      local r;
                      while ((r = __scheduled.shift()) ~= nil) do
                         r();
                      end
   end)}
   -- finally
   if __#scheduled > 0 then
      setTimeout(__runScheduled, 0);
   end
end;

__schedule = function(goroutine)
  if goroutine.asleep then
    goroutine.asleep = false;
    __awakeGoroutines=__awakeGoroutines+1;
  end
  __scheduled.push(goroutine);
  if __curGoroutine == __noGoroutine then
    __runScheduled();
  end
end;

__setTimeout = function(f, t)
  __awakeGoroutines=__awakeGoroutines+1;
  return setTimeout(function()
    __awakeGoroutines=__awakeGoroutines-1;
    f();
  end, t);
end;

__block = function()
  if __curGoroutine == __noGoroutine then
    __throwRuntimeError("cannot block in JavaScript callback, fix by wrapping code in goroutine");
  end
  __curGoroutine.asleep = true;
end;

__send = function(chan, value)
  if chan.__closed then
    __throwRuntimeError("send on closed channel");
  end
  local queuedRecv = chan.__recvQueue.shift();
  if queuedRecv ~= nil then
    queuedRecv([value, true]);
    return;
  end
  if chan.__#buffer < chan.__capacity then
    chan.__buffer.push(value);
    return;
  end

  local thisGoroutine = __curGoroutine;
  local closedDuringSend;
  chan.__sendQueue.push(function(closed)
    closedDuringSend = closed;
    __schedule(thisGoroutine);
    return value;
  end);
  __block();
  return {
    __blk: function()
      if closedDuringSend then
        __throwRuntimeError("send on closed channel");
      end
    end
  end;
end;
__recv = function(chan)
  local queuedSend = chan.__sendQueue.shift();
  if queuedSend ~= nil then
    chan.__buffer.push(queuedSend(false));
  end
  local bufferedValue = chan.__buffer.shift();
  if bufferedValue ~= nil then
    return {bufferedValue, true};
  end
  if chan.__closed then
    return {chan.__elem.zero(), false};
  end

  local thisGoroutine = __curGoroutine;
  local f = { __blk= function() return this.value; end };
  local queueEntry = function(v)
    f.value = v;
    __schedule(thisGoroutine);
  end;
  chan.__recvQueue.push(queueEntry);
  __block();
  return f;
end;
__close = function(chan)
  if chan.__closed then
    __throwRuntimeError("close of closed channel");
  end
  chan.__closed = true;
  while (true) do 
    local queuedSend = chan.__sendQueue.shift();
    if queuedSend == nil then
      break;
    end
    queuedSend(true); -- /* will panic */
  end
  while (true) do 
    local queuedRecv = chan.__recvQueue.shift();
    if queuedRecv == nil then
      break;
    end
    queuedRecv([chan.__elem.zero(), false]);
  end
end;
__select = function(comms)
  local ready = {};
  local selection = -1;
  for i = 0,#comms-1 do

    local comm = comms[i];
    local chan = comm[0];
    local ncomm = #comm
    if ncomm == 0 then -- default
      selection = i;
      break;

    elseif ncomm == 1 then -- recv
      if #chan.__sendQueue ~= 0  or  #chan.__buffer ~= 0  or  chan.__closed then
        ready.push(i);
      end
      break;

    elseif ncomm == 2 then -- send
      if chan.__closed then
        __throwRuntimeError("send on closed channel");
      end
      if #chan.__recvQueue ~= 0  or  #chan.__buffer < chan.__capacity {
        ready.push(i);
      end
      break;
    end
  end

  if #ready ~= 0 then
    selection = ready[Math.floor(Math.random() * #ready)];
  end
  if selection ~= -1 then
    local comm = comms[selection];
    local ncomm = #comm
    if ncomm == 0 then -- default
      return {selection};

    elseif ncomm == 1 then -- recv
      return {selection, __recv(comm[0])};

    elseif ncomm == 2 then -- send
      __send(comm[0], comm[1]);
      return {selection};

    end
  end

  local entries = {};
  local thisGoroutine = __curGoroutine;
  local f = { __blk= function() return this.selection; end };
  local removeFromQueues = function()
    for i = 0,#entries-1 do

      local entry = entries[i];
      local queue = entry[0];
      local index = queue.indexOf(entry[1]);
      if index ~= -1 then
        queue.splice(index, 1);
      end
    end
  end;
  for i = 0,#comms-1 do

    (function(i)
      local comm = comms[i];
      local ncomm = #comm
      if ncomm == 1 then -- recv
        local queueEntry = function(value)
          f.selection = [i, value];
          removeFromQueues();
          __schedule(thisGoroutine);
        end;
        entries.push([comm[0].__recvQueue, queueEntry]);
        comm[0].__recvQueue.push(queueEntry);
        break;
      elseif ncomm == 2 then -- send
        local queueEntry = function()
          if (comm[0].__closed) then
            __throwRuntimeError("send on closed channel");
          end
          f.selection = [i];
          removeFromQueues();
          __schedule(thisGoroutine);
          return comm[1];
        end;
        entries.push([comm[0].__sendQueue, queueEntry]);
        comm[0].__sendQueue.push(queueEntry);
        break;
      end
    end)(i);
  end
  __block();
  return f;
end;

__jsObjectPtr, __jsErrorPtr;

__needsExternalization = function(t)

  local k = t.kind
    if k ==  __kindBool or
    k == __kindInt or
    k == __kindInt8 or
    k == __kindInt16 or
    k == __kindInt32 or
    k == __kindUint or
    k == __kindUint8 or
    k == __kindUint16 or
    k == __kindUint32 or
    k == __kindUintptr or
    k == __kindFloat32 or
    k == __kindFloat64 then
      return false;
    else
      return t ~= __jsObjectPtr;
  end
end;

__externalize = function(v, t)
  if t == __jsObjectPtr then
    return v;
  end
  local sw = t.kind
   if sw ==  __kindBool or
  sw ==kindInt or
  sw ==kindInt8 or
  sw ==kindInt16 or
  sw ==kindInt32 or
  sw ==kindUint or
  sw ==kindUint8 or
  sw ==kindUint16 or
  sw ==kindUint32 or
  sw ==kindUintptr or
  sw ==kindFloat32 or
  sw ==kindFloat64 then
    return v;

  elseif sw ==__kindInt64 or
   sw == __kindUint64 then

    return __flatten64(v);

  elseif we == __kindArray then

    if __needsExternalization(t.elem) then
      return __mapArray(v, function(e) return __externalize(e, t.elem); end);
    end
    return v;

  elseif sw ==  __kindFunc then

    return __externalizeFunction(v, t, false);

  elseif sw ==  __kindInterface then 

    if v == __ifaceNil then
      return nil;
    end
    if v.constructor == __jsObjectPtr then
      return v.__val.object;
    end
    return __externalize(v.__val, v.constructor);

  elseif sw ==  __kindMap then 

    local m = {};
    local keys = __keys(v);
    for i = 0,#keys-1 do
      local entry = v[keys[i]];
      m[__externalize(entry.k, t.key)] = __externalize(entry.v, t.elem);
    end
    return m;
  elseif sw ==  __kindPtr then 

    if v == t.__nil then
      return nil;
    end
    return __externalize(v.__get(), t.elem);

  elseif sw ==  __kindSlice then 

    if __needsExternalization(t.elem) then
      return __mapArray(__sliceToArray(v), function(e) return __externalize(e, t.elem); end);
    end
    return __sliceToArray(v);

  elseif sw ==  __kindString then 

    if __isASCII(v) then
      return v;
    end
    local s = "", r;
    local i = 0
    while(true) do
       if i >= #v then
          break
       end
       r = __decodeRune(v, i);
       local c = r[1];
       if c > 0xFFFF then
          local h = Math.floor((c - 0x10000) / 0x400) + 0xD800;
          local l = (c - 0x10000) % 0x400 + 0xDC00;
          s = s.. String.fromCharCode(h, l);
          continue;
      end
      s = s .. String.fromCharCode(c);
      i = i + r[2]
    end
    return s;

  elseif sw ==  __kindStruct then

    local timePkg = __packages["time"];
    if timePkg ~= nil  and  v.constructor == timePkg.Time.ptr then
      local milli = __div64(v.UnixNano(), new __Int64(0, 1000000));
      return new Date(__flatten64(milli));
    end

    local noJsObject = {};
    local searchJsObject = function(v, t)
      if t == __jsObjectPtr then
        return v;
      end
      local sw2 = t.kind
      if sw2 == __kindPtr then
        if v == t.__nil then
          return noJsObject;
        end
        return searchJsObject(v.__get(), t.elem);
      elseif sk2 ==  __kindStruct then
        local f = t.fields[0];
        return searchJsObject(v[f.prop], f.typ);
      elseif sk2 ==  __kindInterface then
        return searchJsObject(v.__val, v.constructor);
      else
        return noJsObject;
      end
    end;
    local o = searchJsObject(v, t);
    if o ~= noJsObject then
      return o;
    end

    o = {};
    for i = 0,#(t.fields)-1 do

      local f = t.fields[i];
      if  not f.exported then
        continue;
      end

      o[f.name] = __externalize(v[f.prop], f.typ);
    end

    return o;
  end
  __throwRuntimeError("cannot externalize " .. t.__str);

end;

__externalizeFunction = function(v, t, passThis)

  if v == __throwNilPointerError then
    return nil;
  end

  if v.__externalizeWrapper == nil then
    __checkForDeadlock = false;
    v.__externalizeWrapper = function()

      local args = {};
      for i = 0,#params-1 do 
        if t.variadic  and  i == t.#params - 1 then
          local vt = t.params[i].__elem
          local varargs = {};
          for j=i,#arguments-1 do
            varargs.push(__internalize(arguments[j], vt));
          end
          args.push(new (t.params[i])(varargs));
          break;
        end
        args.push(__internalize(arguments[i], t.params[i]));
      end
      local result = v.apply(passThis ? this : nil, args);
      local sw = #(t.result)
      if sw == 0 then
        return;
      elseif sw ==  1 then
        return __externalize(result, t.results[0]);
      else
        for i = 0,#results-1 do
          result[i] = __externalize(result[i], t.results[i]);
        end;
        return result;
      end;
    end;
  end
  return v.__externalizeWrapper;
end;

__internalize = function(v, t, recv)
  if t == __jsObjectPtr then
    return v;
  end
  if t == __jsObjectPtr.elem then
    __throwRuntimeError("cannot internalize js.Object, use *js.Object instead");
  end
  if v  and  v.__internal_object__ ~= nil then
    return __assertType(v.__internal_object__, t, false);
  end
  local timePkg = __packages["time"];
  if timePkg ~= nil  and  t == timePkg.Time then
    if  not (v ~= nil  and  v ~= nil  and  v.constructor == Date) then
      __throwRuntimeError("cannot internalize time.Time from " .. type(v) .. ", must be Date");
    end
    return timePkg.Unix(new __Int64(0, 0), new __Int64(0, v.getTime() * 1000000));
  end
  local sw = t.kind

  if sw == __kindBool then
    return  not  not v;
  elseif sw == __kindInt then
    return parseInt(v);
  elseif sw == __kindInt8 then
    return parseInt(v) << 24 >> 24;
  elseif sw == __kindInt16 then
    return parseInt(v) << 16 >> 16;
  elseif sw == __kindInt32 then
    return parseInt(v) >> 0;
  elseif sw == __kindUint then
    return parseInt(v);
  elseif sw == __kindUint8 then
    return parseInt(v) << 24 >>> 24;
  elseif sw == __kindUint16 then
    return parseInt(v) << 16 >>> 16;
  elseif sw == __kindUint32 or
 sw == __kindUintptr then
    return parseInt(v) >>> 0;
  elseif sw == __kindInt64 or
  elseif sw == __kindUint64 then
    return new t(0, v);
  elseif sw == __kindFloat32 or
  elseif sw == __kindFloat64 then
    return parseFloat(v);
  elseif sw == __kindArray then
    if #v ~= t.len then
      __throwRuntimeError("got array with wrong size from JavaScript native");
    end
    return __mapArray(v, function(e) return __internalize(e, t.elem); end);
  elseif sw == __kindFunc then
    return function()
      local args = [];
      for i = 0,t.#params-1 do
        if t.variadic  and  i == t.#params - 1 then
          local vt = t.params[i].__elem
          local varargs = arguments[i];
          for j = 0, #varargs-1 do
            args.push(__externalize(varargs.__array[varargs.__offset + j], vt));
          end
          break;
        end
        args.push(__externalize(arguments[i], t.params[i]));
      end
      local result = v.apply(recv, args);
      local sw2 = #(t.results)
      if sw2 == 0 then
        return;
      elseif sw2 == 1 then
        return __internalize(result, t.results[0]);
      else
        for i=0,#(t.results)-1 do
          result[i] = __internalize(result[i], t.results[i]);
        end
        return result;
      end
    end;
  elseif sw == __kindInterface then
    if #t.methods ~= 0 then
      __throwRuntimeError("cannot internalize " .. t.__str);
    end
    if v == nil then
      return __ifaceNil;
    end
    if v == nil then
      return new __jsObjectPtr(nil);
    end
    local vc = v.constructor

    if vc == Int8Array then
      return new (__sliceType(__Int8))(v);
    elseif vc ==  Int16Array then
      return new (__sliceType(__Int16))(v);
    elseif vc ==  Int32Array then
      return new (__sliceType(__Int))(v);
    elseif vc ==  Uint8Array then
      return new (__sliceType(__Uint8))(v);
    elseif vc ==  Uint16Array then
      return new (__sliceType(__Uint16))(v);
    elseif vc ==  Uint32Array then
      return new (__sliceType(__Uint))(v);
    elseif vc ==  Float32Array then
      return new (__sliceType(__Float32))(v);
    elseif vc ==  Float64Array then
      return new (__sliceType(__Float64))(v);
    elseif vc ==  Array then
      return __internalize(v, __sliceType(__emptyInterface));
    elseif vc ==  Boolean then
      return new __Bool( not  not v);
    elseif vc ==  Date then
      if timePkg == nil then
        -- /* time package is not present, internalize as &js.Object{Dateend so it can be externalized into original Date. */
        return new __jsObjectPtr(v);
      end
      return new timePkg.Time(__internalize(v, timePkg.Time));
    elseif vc ==  Function then
      local funcType = __funcType([__sliceType(__emptyInterface)], [__jsObjectPtr], true);
      return new funcType(__internalize(v, funcType));
    elseif vc ==  Number then
      return new __Float64(parseFloat(v));
    elseif vc ==  String then
      return new __String(__internalize(v, __String));
    else
      if __global.Node  and  v instanceof __global.Node then
        return new __jsObjectPtr(v);
      end
      local mapType = __mapType(__String, __emptyInterface);
      return new mapType(__internalize(v, mapType));
    end
  elseif sw ==  __kindMap then
    local m = {};
    local keys = __keys(v);
    for i = 0,#keys-1 do
      local k = __internalize(keys[i], t.key);
      m[t.key.keyFor(k)] = { k= k, v= __internalize(v[keys[i]], t.elem) };
    end
    return m;
  elseif sw ==  __kindPtr then
    if t.elem.kind == __kindStruct then
      return __internalize(v, t.elem);
    end
  elseif sw ==  __kindSlice then
    return new t(__mapArray(v, function(e) return __internalize(e, t.elem); end));
  elseif sw ==  __kindString then
    v = String(v);
    if __isASCII(v) then
      return v;
    end
    local s = "";
    local i = 0;
    while i < #v do
      local h = v.charCodeAt(i);
      if 0xD800 <= h  and  h <= 0xDBFF then
        local l = v.charCodeAt(i + 1);
        local c = (h - 0xD800) * 0x400 + l - 0xDC00 + 0x10000;
        s = s..__encodeRune(c);
        i =i+ 2;
        continue;
      end
      s = s..__encodeRune(h);
      i=i+1
    end
    return s;
  elseif sw ==  __kindStruct then
    local noJsObject = {};
    local searchJsObject = function(t)
      if t == __jsObjectPtr then
        return v;
      end
      if t == __jsObjectPtr.elem then
        __throwRuntimeError("cannot internalize js.Object, use *js.Object instead");
      end

      if t.kind ==  __kindPtr then
        return searchJsObject(t.elem);

      elseif t.kind == __kindStruct then
        local f = t.fields[0];
        local o = searchJsObject(f.typ);
        if (o ~= noJsObject) {
          local n = new t.ptr();
          n[f.prop] = o;
          return n;
        end
        return noJsObject;
     else
        return noJsObject;
      end
    end;
    local o = searchJsObject(t);
    if (o ~= noJsObject) {
      return o;
    end
  end
  __throwRuntimeError("cannot internalize " .. t.__str);
end;

-- /* __isASCII reports whether string s contains only ASCII characters. */
__isASCII = function(s)
  for i=0,#s-1 do
    if s.charCodeAt(i) >= 128 then
      return false;
    end
  end
  return true;
end;

-- packages were here.

-- at end of file, after all the packages.

__synthesizeMethods();
__mainPkg = __packages["github.com/gijit/gi/pkg/compiler/tmp"];
-- __packages["runtime"].__init();
__go(__mainPkg.__init, {});
__flushConsole();

