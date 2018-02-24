--
-- tsys.lua provides the type system for gijit.
--
-- It started life as a port of the GopherJS type
-- system to LuaJIT, and still shows some
-- javascript vestiges.

-- We would typically assume these dofile imports
-- are already done by prelude loading.
-- For dev work, we'll load them if not already.
--

-- __minifs only has getcwd and chdir.
-- Just enough to bootstrap.
--
__minifs = {}
__ffi = require "ffi"
local __osname = __ffi.os == "Windows" and "windows" or "unix"

__built_in_starting_symbol_list={};
__built_in_starting_type_list={};

function __storeBuiltins()
   for k,_ in pairs(_G) do
      __built_in_starting_symbol_list[k]=true;
   end
   for k,_ in pairs(__type__) do
      __built_in_starting_type_list[k]=true;
   end
end

function __sorted_keys(t)
   local r = {}
   for k,_ in pairs(t) do
      table.insert(r, k)
   end
   table.sort(r)
   return r
end

-- list all global vars
function __gls()
   local i = 0
   for _,k in ipairs(__sorted_keys(_G)) do
      local v = _G[k]      
      i=i+1
      print("["..tostring(i).."] "..tostring(k).." : "..tostring(v))
   end
end

-- list user vars
function __ls()
   local i = 0
   for _,k in ipairs(__sorted_keys(_G)) do
      local uscore = 95 -- "_"
      if #k > 2 and string.byte(k,1,1)==uscore and string.byte(k,2,2) == uscore then
         -- we omit __ prefixed methods/values
      else
         if __built_in_starting_symbol_list == nil or 
         not __built_in_starting_symbol_list[k] then
            
            local v = _G[k]
            i=i+1
            print("["..tostring(i).."] "..tostring(k).." : "..tostring(v))
         end
      end
   end
end

-- list types
function __lst()
   local i = 0
   for _,k in ipairs(__sorted_keys(__type__)) do
      if __built_in_starting_type_list == nil or 
      not __built_in_starting_type_list[k] then
         local v = __type__[k]
         i=i+1
         print("["..tostring(i).."] "..tostring(k).." : "..tostring(v))
      end
   end
end

-- global list types
function __glst()
   local i = 0
   for _,k in ipairs(__sorted_keys(__type__)) do
      local v = __type__[k]
      i=i+1
      print("["..tostring(i).."] "..tostring(k).." : "..tostring(v))
   end
end


-- a __ namespace binding so __tostring is usable from Go
__tostring = tostring

local __dq = function(str)
   local ty = type(str)
   if ty == "string" then
      return '"'..str..'"'
   elseif ty ~= "table" then
      return tostring(str)
   end
   
   -- avoid infinite loops.
   local mt = getmetatable(str)
   setmetatable(str, nil)
   local s = tostring(str)
   setmetatable(str, mt)
   return s
end

local __system = ({
      windows	= {
         getcwd	= "_getcwd",
         chdir	= "_chdir",
         maxpath	= 260,
      },
      unix	= {
         getcwd	= "getcwd",
         chdir	= "chdir",
         maxpath	= 4096,
      }
                 })[__osname]

__ffi.cdef(
   [[
		char   *]] .. __system.getcwd .. [[ ( char *buf, size_t size );
		int		]] .. __system.chdir  .. [[ ( const char *path );
		]]
)

__minifs.getcwd = function ()
   local buff = __ffi.new("char[?]", __system.maxpath)
   __ffi.C[__system.getcwd](buff, __system.maxpath)
   return __ffi.string(buff)
end

__minifs.chdir = function (path)
   return __ffi.C[__system.chdir](path) == 0
end

-- Ugh, it renames.
-- So only use this on the "__gijit_prelude" marker file,
-- which is a file of little importance, only there
-- to verify our path is correct.
function __minifs.renameBasedFileExists(file)
   local ok, err, code = os.rename(file, file)
   if not ok then
      if code == 13 then
         -- denied, but it exists
         return true
      end
   end
   return ok, err
end

function __minifs.dirExists(path)
   -- "/" works on both Unix and Windows
   return __minifs.fileExists(path.."/")
end

-- The point of __minifs is so we can find
-- and set __preludePath if it is not set.
-- It will always be set by gijit, but this
-- allows standalone development and testing.
--
if __preludePath == nil then
   print("__preludePath is nil...")
   local origin=""
   local dir = os.getenv("GIJIT_PRELUDE_DIR")
   if dir ~= nil then
      origin = "__preludePath set from GIJIT_PRELUDE_DIR"
      __preludePath = dir .. "/"
   else
      local defaultPreludePath = "/src/github.com/gijit/gi/pkg/compiler"
      local gopath = os.getenv("GOPATH")
      if gopath ~= nil then
         origin = "__preludePath set from GOPATH"
         __preludePath = gopath .. defaultPreludePath .. "/"
      else
         -- try $HOME/go
         local home = os.getenv("HOME")
         if home ~= nil then
            origin = "__preludePath set from $HOME/go"
            __preludePath = home .. "/go" .. defaultPreludePath .. "/"
         else
            -- default to cwd
            origin = "__preludePath set from cwd"     
            __preludePath = __minifs.getcwd().."/"
         end
      end
   end
   -- check for our marker file.
   if not __minifs.renameBasedFileExists(__preludePath.."__gijit_prelude") then
      error("error in tsys.lua: could not find my prelude directory. Tried __preludePath='"..__preludePath.."'; "..origin)
   end
   print("using __preludePath = '"..__preludePath.."'")
end

if __min == nil then
   dofile(__preludePath..'math.lua') -- for __max, __min, __truncateToInt
end
if int8 == nil then
   dofile(__preludePath..'int64.lua') -- for integer types with Go naming.
end
if complex == nil then
   dofile(__preludePath..'complex.lua')
end

if __dfsOrder == nil then
   dofile(__preludePath..'dfs.lua')
end

-- global for now, later figure out to scope down.
__dfsGlobal = __NewDFSState()

-- tell Luar that it is running under gijit,
-- by setting this global flag.
__gijit_tsys = true

-- translation of javascript builtin 'prototype' -> typ.prototype
--                                   'constructor' -> typ.__constructor

__bit = require("bit")

-- ===========================
--
-- begin actual type system stuff
--
-- ===========================

__type__ ={}; -- global repo of types
__global ={};
__module ={};
__packages = {}
__idCounter = 0;
__pkg = {};

-- length of array, counting [0] if present,
-- but trusting __len if metamethod avail.
function __lenz(array)
   if array == nil then
      print(debug.traceback())
      error("cannot call __lenz with nil array")
   end
   local n = #array
   local mt = getmetatable(array)
   if mt ~= nil and mt.__len ~= nil then
      --print("__len was not nil")
      return n
   end
   --print("__len was nil")
   if array[0] ~= nil then
      --print("array[0] is not nil")
      n=n+1
   end
   return n
end

-- return an int64 as the value, not a double
function __lenzi(array)
   local n = #array
   local mt = getmetatable(array)
   if mt ~= nil and mt.__len ~= nil then
      --print("__len was not nil")
      return int(n)
   end
   --print("__len was nil")
   if array[0] ~= nil then
      --print("array[0] is not nil")
      n=n+1
   end
   return int(n)
end

function __ipairsZeroCheck(arr)
   if arr[0] ~= nil then error("ipairs will miss the [0] index of this array") end
end

__mod = function(y) return x % y; end;
__parseInt = parseInt;
__parseFloat = function(f)
   if f ~= nil  and  f ~= nil  and  f.constructor == Number then
      return f;
   end
   return parseFloat(f);
end;

-- __fround returns nearest float32
__fround = function(x)
   return float32(x)
end;

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

-- the __kind numbers must be kept in sync with rtyp.go.
-- returned by __basicValue2kind(v) on unrecognized kind.
__kindUnknown = -1;

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

__kind2str = {
   [1]="__kindBool",
   [2]="__kindInt",
   [3]="__kindInt8",
   [4]="__kindInt16",
   [5]="__kindInt32",
   [6]="__kindInt64",
   [7]="__kindUint",
   [8]="__kindUint8",
   [9]="__kindUint16",
   [10]="__kindUint32",
   [11]="__kindUint64",
   [12]="__kindUintptr",
   [13]="__kindFloat32",
   [14]="__kindFloat64",
   [15]="__kindComplex64",
   [16]="__kindComplex128",
   [17]="__kindArray",
   [18]="__kindChan",
   [19]="__kindFunc",
   [20]="__kindInterface",
   [21]="__kindMap",
   [22]="__kindPtr",
   [23]="__kindSlice",
   [24]="__kindString",
   [25]="__kindStruct",
   [26]="__kindUnsafePointer",
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

function __addressof(t)
   local mt = getmetatable(t)
   setmetatable(t, nil)
   local addr = tostring(t)
   setmetatable(t, mt)
   return addr
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
   -- restore the metatable just before returning!
   
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
         vals = __st(v,"",indent+3,quiet,methods_desc, seen)
      else
         local vmt = getmetatable(v)
         if type(v) == "table" and type(vmt) == "table" then
            setmetatable(v, nil)
            vals = tostring(v)
            setmetatable(v, vmt)
         else
            vals = tostring(v)
         end
         --vals = __st(v,"",indent+1,true,methods_desc, seen) or ""
      end
      local ty = type(i)
      if ty == "cdata" then
         ty = tostring(__ffi.typeof(i))
      end
      s = s..pre.." "..tostring(k).. " key ("..ty.."): '"..tostring(i).."' val: '"..vals.."'\n"
   end
   if k == 0 then
      s = pre.."<empty table> " .. addr
   end

   if mt ~= nil then
      s = s or ""
      local show = __st(mt, "mt.of."..name, indent+1, true, methods_desc, seen) or ""
      s = s .. "\n"..show
   end
   if not quiet then
      print(s)
   end
   -- restore metamethods
   setmetatable(t, mt)
   --print("__st returning '"..tostring(s).."'")
   return s or ""
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
   --print(debug.traceback())
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
      -- TODO: port this!
      print("jea TODO: port this, what is __externalize doing???")
      error("NOT DONE: port this!")
      --return __externalize(fn(this, (__sliceType(__jsObjectPtr))(__global.Array.prototype.slice.call(arguments, {}))), __type__.emptyInterface);
   end;
end;
__unused = function(v) end;

--
__mapArray = function(arr, fun)
   local newarr = {}
   -- handle a zero argument, if present.
   local bump = 0
   local zval = arr[0]
   if zval ~= nil then
      bump = 1
      newarr[1] = fun(zval)
   end
   __ipairsZeroCheck(arr)
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

-- low, high are 0-based slice [low,high). max is the
--  maximum capacity of the slice.
__subslice = function(slice, low, high, max)
   if high == nil then
      
   end
   if low < 0  or  (high ~= nil and high < low)  or  (max ~= nil and high ~= nil and max < high)  or  (high ~= nil and high > slice.__capacity)  or  (max ~= nil and max > slice.__capacity) then
      __throwRuntimeError("slice bounds out of range");
   end
   
   local s = slice.__constructor.tfun(slice.__array);
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
   return int(n);
end;

--

__copyArray = function(dst, src, dstOffset, srcOffset, n, elem)
   --print("__copyArray called with n = ", n, " dstOffset=", dstOffset, " srcOffset=", srcOffset)
   --print("__copyArray has dst:")
   --__st(dst)
   --print("__copyArray has src:")
   --__st(src)
   
   n = tonumber(n)
   if n == 0  or  (dst == src  and  dstOffset == srcOffset) then
      --setmetatable(dst, getmetatable(src))
      return;
   end

   local sw = elem.kind
   if sw == __kindArray or sw == __kindStruct then
      
      if dst == src  and  dstOffset > srcOffset then
         for i = n-1,0,-1 do
            elem.copy(dst[dstOffset + i], src[srcOffset + i]);
         end
         --setmetatable(dst, getmetatable(src))         
         return;
      end
      for i = 0,n-1 do
         elem.copy(dst[dstOffset + i], src[srcOffset + i]);
      end
      --setmetatable(dst, getmetatable(src))      
      return;
   end

   if dst == src  and  dstOffset > srcOffset then
      for i = n-1,0,-1 do
         dst[dstOffset + i] = src[srcOffset + i];
      end
      --setmetatable(dst, getmetatable(src))      
      return;
   end
   for i = 0,n-1 do
      dst[dstOffset + i] = src[srcOffset + i];
   end
   --setmetatable(dst, getmetatable(src))   
   --print("at end of array copy, src is:")
   --__st(src)
   --print("at end of array copy, dst is:")
   --__st(dst)
end;

--
__clone = function(src, typ)
   local clone = typ()
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
      for _,f in ipairs(typ.elem.fields) do
         helper(f.__prop);
      end
      
      proxy = Object.create(typ.prototype, properties);
      proxy.__val = proxy;
      obj.__proxies[typ.__str] = proxy;
      proxy.__proxies = obj.__proxies;
   end
   return proxy;
end;

--


__append = function(...)
   local arguments = {...}
   local slice = arguments[1]
   return __internalAppend(slice, arguments, 1, #arguments - 1);
end;

__appendSlice = function(slice, toAppend)

   -- recognize and resolve the ellipsis.
   if type(toAppend) == "table" then
      if toAppend.__name == "__lazy_ellipsis_instance" then
         --print("resolving lazy ellipsis.")
         toAppend = toAppend() -- resolve the lazy reference.
      end
   end
   --print("toAppend:")
   --__st(toAppend, "toAppend")
   --print("slice:")
   --__st(slice, "slice")
   
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

   local newSlice = slice.__constructor.tfun(newArray);
   newSlice.__offset = newOffset;
   newSlice.__length = newLength;
   newSlice.__capacity = newCapacity;
   return newSlice;
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
   cp.__length = k
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
   end,
}

-- use for slices and arrays
__valueSliceIpairs = function(t)
   
   --print("__ipairs called!")
   -- this makes a slice work in a for k,v in ipairs() do loop.
   local off = rawget(t, "__offset")
   local slcLen = rawget(t, "__length")
   local function stateless_iter(arr, k)
      k=k+1
      if k >= off + slcLen then
         return
      end
      return k, rawget(arr, off + k)
   end       
   -- Return an iterator function, the table, starting point
   local arr = rawget(t, "__array")
   --print("arr is "..tostring(arr))
   return stateless_iter, arr, -1
end

__valueArrayMT = {
   __name = "__valueArrayMT",

   __ipairs = __valueSliceIpairs,
   __pairs  = __valueSliceIpairs,
   
   __newindex = function(t, k, v)
      --print("__valueArrayMT.__newindex called, t is:")
      --__st(t)
      local w = tonumber(k)
      
      if w < 0 or w >= #t then
         error "write to array error: access out-of-bounds"
      end
      
      t.__val[w] = v
   end,
   
   __index = function(t, k)
      --print("__valueArrayMT.__index called, k='"..tostring(k).."'")
      local ktype = type(k)
      if ktype == "string" then
         print("ktype was string, doing rawget??? why?")
         print(debug.traceback())
         error("where is this used?")
         return rawget(t,k)
      elseif type(k) == "table" then
         print("callstack:"..tostring(debug.traceback()))
         error("table as key not supported in __valueArrayMT")
      elseif ktype == "cdata" then
         k = tonumber(k)
      end
      --__st(t.__val)
      if k < 0 or k >= #t then
         print(debug.traceback())
         error("read from array error: access out-of-bounds; "..tostring(k).." is outside [0, "  .. tostring(#t) .. ")")
      end
      --print("array access bounds check ok.")
      --__st(t.__array, "t.__array")
      return t.__array[k]
   end,

   __len = function(t)
      return __lenz(t.__val)
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

         local len = #self.__val
         if self.__val[0] ~= nil then
            len=len+1
         end
         local s = self.__constructor.__str.."{"
         local raw = self.__val
         local beg = 0

         local quo = ""
         if len > 0 and type(raw[beg]) == "string" then
            quo = '"'
         end
         for i = 0, len-1 do
            s = s .. "["..tostring(i).."]" .. "= " ..quo.. tostring(raw[beg+i]) .. quo .. ", "
         end
         
         return s .. "}"
      end
      
      if getmetatable(self.__val) == __valueArrayMT then
         --print("avoid infinite loop")
         return "<avoid inf loop>"
      else
         return tostring(self.__val)
      end
   end,
}


__valueSliceMT = {
   __name = "__valueSliceMT",
   
   __newindex = function(t, k, v)
      --print("__valueSliceMT.__newindex called, t is:")
      --__st(t)
      local w = tonumber(t.__offset + k)
      if k < 0 or k >= t.__capacity then
         error "slice error: write out-of-bounds"
      end
      t.__array[w] = v
   end,
   
   __index = function(t, k)
      
      --print("__valueSliceMT.__index called, k='"..tostring(k).."'")
      --__st(t.__val)
      --print("callstack:"..tostring(debug.traceback()))

      local ktype = type(k)
      if ktype == "string" then
         --print("we have string key, doing rawget on t")
         --__st(t, "t")
         return rawget(t,k)
      elseif ktype == "table" then
         print("callstack:"..tostring(debug.traceback()))
         error("table as key not supported in __valueSliceMT")
      elseif ktype == "cdata" then
         -- we may be called with 0LL, but arrays in Lua
         -- must be indexed with float64, so we'll convert
         --print("converting ktype cdata into number...k='"..tonumber(k).."'")
         k = tonumber(k)
      end
      local w = t.__offset + k
      --print("index slice with w = "..type(w).." value: "..tostring(w))
      if k < 0 or k >= t.__capacity then
         print(debug.traceback())
         error("slice error: access out-of-bounds, k="..tostring(k).."; cap="..tostring(t.__capacity))
      end
      --print("slice access bounds check ok: w = ", w)
      --__st(t.__array, "t.__array")
      local wv = t.__array[w]
      --__st(wv, "wv")
      --__st(getmetatable(wv), "metatable.for.wv")
      --print("wv back from t.__array[w] is: "..tostring(wv))
      return wv
   end,

   __len = function(t)
      --print("__valueSliceMT metamethod __len called, returning ", t.__length)
      return t.__length
   end,
   
   __tostring = function(self, ...)
      --print("__tostring called from __valueSliceMT: "..self.__typ.__str);

      local len = tonumber(self.__length) -- convert from LL int
      local off = tonumber(self.__offset)
      --print("__tostring sees self.__length of ", len, " __offset = ", off)
      local cap = self.__capacity
      --local s = "slice <len=" .. tostring(len) .. "; off=" .. off .. "; cap=" .. cap ..  "> is "..self.__constructor.__str.."{"
      --print("self.__constructor.__str = '"..self.__constructor.__str.."' and full display of self.__constructor:")
      --__st(self.__constructor, "self.__constructor")
      
      local s = self.__constructor.__str.."{"
      local raw = self.__array
      --__st(raw, "raw in valueSice tostring")
      
      -- we want to skip both the _giPrivateRaw and the len
      -- when iterating, which happens automatically if we
      -- iterate on raw, the raw inside private data, and not on the proxy.
      local quo = ""
      if len > 0 and type(raw[off]) == "string" then
         quo = '"'
      end
      for i = 0, len-1 do
         s = s .. "["..tostring(i).."]" .. "= " ..quo.. tostring(raw[off+i]) .. quo .. ", "
      end
      
      return s .. "}"
      
   end,
   __pairs = __valueSliceIpairs,
   __ipairs = __valueSliceIpairs,
}


__tfunBasicMT = {
   __name = "__tfunBasicMT",
   __call = function(self, ...)
      --print("jea debug: __tfunBasicMT.__call() invoked") -- , self='"..tostring(self).."' with tfun = ".. tostring(self.tfun).. " and args=")
      --print(debug.traceback())

      --print("debug, __tfunBasicMT itself as a table is:")
      --__st(__tfunBasicMT, "__tfunBasicMT")
      local args = {...}
      --print("in __tfunBasicMT, start __st on args")
      --__st(args, "args to __tfunBasicMT.__call")
      --print("in __tfunBasicMT,   end __st on args")

      --print("in __tfunBasicMT, start __st on self")
      --__st(self, "self")
      --print("in __tfunBasicMT,   end __st on self")

      -- args1 will have the empty instance if typ is invoked as, e.g.
      --    s = __type__.S(0LL);
      -- Since we just replace it with newInstance anyway, change that in expressions.go:289 and :349
      -- to just be
      --    s = __type__.S(0LL);
      
      if self ~= nil then
         if self.tfun ~= nil then
            --print("calling tfun! -- let constructors set metatables if they wish to")

            -- this makes a difference as to whether or
            -- not the ctor receives a nil 'this' or not...
            -- So *don't* set metatable here, let ctor do it.
            -- setmetatable(newInstance, __valueBasicMT)

            -- define (bloom) any lazy types we depend on
            if self.__dfsNode == nil then
               error("typ.__dfsNode was nil for type "..self.__str.." at "..__st(self, "self", 0, true))
            end
            if not self.__dfsNode.made then
               --turn off to bootstrap, TODO: enable this!
               --self.__dfsNode:makeRequiredTypes()
            end
            
            -- get zero value if no args
            if #{...} == 0 and self.zero ~= nil then
               local sz = self.zero()
               --print("tfun sees no args and we have a typ.zero() method, so invoking self.zero() got back sz=")
               --__st(sz, "sz")
               return self.tfun(sz)
            else
               return self.tfun(...)
            end
         end
      else
         local newInstance = {}
         --print("in __tfunBasicMT, made newInstance = ")
         --__st(newInstance,"newInstance")
         
         setmetatable(newInstance, __valueBasicMT)

         if self ~= nil then
            --print("self.tfun was nil")
         end
      end
      return newInstance
   end,
}

function __starToAsterisk(s)
   -- parenthesize to get rid of the
   -- substitution count.
   return (string.gsub(s,"*","&"))
end

__valuePointerMT = {
   __name = "__valuePointerMT",
   
   __newindex = function(t, k, v)
      --print("__valuePointerMT: __newindex called, calling set() with val=", v)
      return t.__set(v)
   end,

   __index = function(t, k)
      --print("__valuePointerMT: __index called, doing get()")       
      return t.__get()
   end,

   __tostring = function(t)
      --print("__valuePointerMT: tostring called")
      --__st(t, "t")
      if t == nil then return "<__valuePointer nil pointer>" end
      -- avoid getting messed up by metatable intercept, do a rawget.
      local typ = rawget(t, "__typ")
      if typ == nil then
         print(debug.traceback())
         error("must have __typ set on pointer target")
      end
      local ret = __starToAsterisk(typ.__str) .. "{" .. tostring(t.__get()) .. "}"
      return ret
   end
}

-- a __valueStructMT shouldn't be needed/used, instead the methodSet should
-- be the MT for any struct, even if it is emtpy of methods.
-- a.k.a. this is now called prototype


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
   --print("__synthesizeMethods called! we have #__methodSynthesizers = "..tostring(#__methodSynthesizers))
   __ipairsZeroCheck(__methodSynthesizers)
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
   --print("__newType called with str = '"..str.."'")
   local typ ={
      __str = str,
   };
   typ.__dfsNode = __dfsGlobal:newDfsNode(str, typ)

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
      typ.tfun = function(v)
         local this={};
         this.__val = v;
         this.__typ = typ;
         setmetatable(this, __valueBasicMT);
         return this;
      end;
      typ.wrapped = true;
      typ.keyFor = function(x) return tostring(x); end;

   elseif kind == __kindString then
      
      typ.tfun = function(v)
         local this={};
         --print("strings' tfun called! with v='"..tostring(v).."' and this:")
         --__st(this)
         this.__val = v;
         this.__typ = typ         
         setmetatable(this, __valueBasicMT)
         return this;
      end;
      typ.wrapped = true;
      typ.keyFor = __identity; -- function(x) return "_" .. x; end;

   elseif kind == __kindFloat32 or
   kind == __kindFloat64 then
      
      typ.tfun = function(v)
         local this={};
         this.__val = v;
         this.__typ = typ         
         setmetatable(this, __valueBasicMT)
         return this;
      end;
      typ.wrapped = true;
      typ.keyFor = function(x) return __floatKey(x); end;


   elseif kind ==  __kindComplex64 then

      typ.tfun = function(re, im)
         local this = {};
         this.__val = re + im*complex(0,1);
         this.__typ = typ;
         setmetatable(this, __valueBasicMT);
         return this;
      end;
      typ.wrapped = true;
      typ.keyFor = function(x) return tostring(x); end;
      
      --    typ.tfun = function(real, imag)
      --      local this={};
      --      this.__real = __fround(real);
      --      this.__imag = __fround(imag);
      --      this.__val = this;
      --      return this;
      --    end;
      --    typ.keyFor = function(x) return x.__real .. "_" .. x.__imag; end;

   elseif kind ==  __kindComplex128 then

      typ.tfun = function(re, im)
         local this = {}
         this.__val = re + im*complex(0,1);
         this.__typ = typ
         setmetatable(this, __valueBasicMT)
         return this
      end;
      typ.wrapped = true;
      typ.keyFor = function(x) return tostring(x); end;
      
      --     typ.tfun = function(real, imag)
      --        local this={};
      --        this.__real = real;
      --        this.__imag = imag;
      --        this.__val = this;
      --        this.__constructor = typ
      --        return this;
      --     end;
      --     typ.keyFor = __identity --function(x) return x.__real .. "_" .. x.__imag; end;
      --    
      
   elseif kind ==  __kindPtr then

      if constructor ~= nil then
         --print("in newType kindPtr, constructor is not-nil: "..tostring(constructor))
      end
      typ.tfun = constructor  or
         function(getter, setter, target)
            local this={};
            --print("in tfun for pointer: ",debug.traceback())
            --print("pointer typ.tfun which is same as constructor called! getter='"..tostring(getter).."'; setter='"..tostring(setter).."; target = '"..tostring(target).."'")
            -- sanity checks
            if setter ~= nil and type(setter) ~= "function" then
               error "setter must be function"
            end
            if getter ~= nil and type(getter) ~= "function" then
               error "getter must be function"
            end
            this.__get = getter;
            this.__set = setter;
            this.__target = target;
            this.__val = this; -- seems to indicate a non-primitive value.
            this.__typ = typ
            setmetatable(this, __valuePointerMT)
            return this;
         end;
      typ.keyFor = __idKey;
      
      typ.init = function(elem)
         --print("init(elem) for pointer type called.")
         __dfsGlobal:addChild(typ, elem)
         typ.elem = elem;
         typ.wrapped = (elem.kind == __kindArray);
         typ.__nil = typ(__throwNilPointerError, __throwNilPointerError);
      end;

   elseif kind ==  __kindSlice then
      
      typ.tfun = function(array)
         local this={};
         --print(debug.traceback())
         --print("slice tfun for type '"..__addressof(typ).."' called with array = ")
         --__st(array)
         this.__array = array;
         this.__offset = 0;
         this.__length = __lenz(array)
         --print("# of array returned ", this.__length)
         this.__capacity = this.__length;
         --print("jea debug: slice tfun set __length to ", this.__length)
         --print("jea debug: slice tfun set __capacity to ", this.__capacity)
         --print("jea debug: slice tfun sees array: ")
         --for i,v in pairs(array) do
         --print("array["..tostring(i).."] = ", v)
         --end
         
         this.__val = this;
         this.__constructor = typ
         this.__name = "__sliceValue"
         this.__typ = typ
         setmetatable(this, __valueSliceMT)
         return this
      end;
      typ.init = function(elem)
         typ.elem = elem;
         typ.comparable = false;
         typ.__nil = typ({});
      end;
      
   elseif kind ==  __kindArray then
      typ.tfun = function(v)
         local this={};
         --print("in tfun ctor function for __kindArray, this="..tostring(this).." and v="..tostring(v))
         this.__val = v;
         this.__array = v; -- like slice, to reuse ipairs method.
         this.__offset = 0; -- like slice.
         this.__constructor = typ
         this.__length = __lenz(v)
         this.__name = "__arrayValue"
         this.__typ = typ         
         setmetatable(this, __valueArrayMT)
         return this
      end;
      --print("in newType for array, and typ.tfun = "..tostring(typ.tfun))
      typ.wrapped = true;
      typ.ptr = __newType(4, __kindPtr, "*" .. str, false, "", false, function(array)
                             local this={};
                             this.__get = function() return array; end;
                             this.__set = function(v) typ.copy(this, v); end;
                             this.__val = array;
                             return this
      end);
      
      -- track the dependency between types
      __dfsGlobal:addChild(typ.ptr, typ)
      
      typ.init = function(elem, len)
         --print("init() called for array.")
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
         -- pointer. But perhaps this is javascript's prototypal inheritance in action.
         --
         -- gopherjs uses them in comma expressions. example, condensed:
         --     p$1 = new ptrType(...); sa$3.Port = (p$1.nilCheck, p$1[0])
         --
         -- Since comma expressions are not (efficiently) supported in Lua, lets
         -- implement the nil check in a different manner.
         -- js: Object.defineProperty(typ.ptr.__nil, "nilCheck", { get= __throwNilPointerError end);
      end;
      -- end __kindArray

      
   elseif kind ==  __kindChan then
      
      typ.tfun = function(v)
         local this={};
         this.__val = v;
         this.__typ = typ
         setmetatable(this, __valueBasicMT)
         return this
      end;
      typ.wrapped = true;
      typ.keyFor = __idKey;
      typ.init = function(elem, sendOnly, recvOnly)
         typ.elem = elem;
         typ.sendOnly = sendOnly;
         typ.recvOnly = recvOnly;
      end;
      

   elseif kind ==  __kindFunc then 

      typ.tfun = function(v)
         local this={};
         this.__val = v;
         this.__typ = typ
         setmetatable(this, __valueBasicMT)
         return this;
      end;
      typ.wrapped = true;
      typ.init = function(params, results, variadic)
         typ.params = params;
         typ.results = results;
         typ.variadic = variadic;
         typ.comparable = false;
      end;
      

   elseif kind ==  __kindInterface then 

      typ.implementedBy= {}
      typ.missingMethodFor= {}
      
      typ.keyFor = __ifaceKeyFor;
      typ.init = function(methods)
         --print("top of init() for kindInterface, methods= ")
         --__st(methods)
         --print("and also at top of init() for kindInterface, typ= ")
         --__st(typ)
         typ.methods = methods;
         for _, m in pairs(methods) do
            -- TODO:
            -- jea: why this? seems it would end up being a huge set?
            --print("about to index with m.__prop where m =")
            --__st(m)
            __ifaceNil[m.__prop] = __throwNilPointerError;
         end;
      end;
      
      
   elseif kind ==  __kindMap then 
      
      typ.tfun = function(entries)
        --print("map tfun called, entries = "..tostring(entries))
         local this={};
         this.__typ = typ         
         this.__val = {}; --no meta names, so clean. No accidental collisions.

         local kff = typ.elem.keyFor
         this.nilKeyStored = false
         
         local len=0
         for k, e in pairs(entries) do
            if k == nil then
               this.nilKeyStored = true
               this.nilValue = e or __intentionalNilValue
            else 
               local key = tostring(kff(k)) -- must be a string!
               --print("using key ", key, " for k=", k)
               this.__val[key] = e or __intentionalNilValue;
            end
            len=len+1;
         end
         
         this.len=len;
         this.keyType=typ.key
         this.elemType=typ.elem
         this.zeroValue = typ.elem.zero()
         
         setmetatable(this, __valueMapMT);
         return this;
                                end;
      typ.wrapped = true;
      typ.init = function(key, elem)
         typ.key = key;
         typ.elem = elem;
         typ.comparable = false;
      end;
      
   elseif kind ==  __kindStruct then

      -- provides a way to inject the fields from __structType(),
      -- centralizes the finishing up of struct value construction.
      typ.finishStructValueCreation=function(this, fields, args)
         this.__val = this;
         for i,fld in ipairs(fields) do
            this[fld.__prop] = args[i] or fld.__typ.zero();
         end         
         this.__name = "__structValue";
         this.__typ = typ;
         setmetatable(this, typ.prototype)
      end
      
      typ.tfun = function(...)
         local args = {...}
         --print("top of simple kindStruct tfun, args are:")
         --__st(args, "args")
         
         local this
         if typ.__constructor ~= nil then
            --print("simple kindStruct: typ.__constructor was not nil, is")
            --__st(typ.__constructor, "typ.__constructor")
            this = typ.__constructor(...);
         else
            --print("simple kindStruct: typ.__constructor was nil, typ.fields:")
            --__st(typ.fields, "typ.fields")
            --__st(typ.fields[1], "typ.fields[1]")
            --__st(typ.fields[1].__typ, "typ.fields[1].__typ")
            this={}
         end
         typ.finishStructValueCreation(this, typ.fields, args)
         return this
      end;
      typ.wrapped = true;

      -- the typ.prototype will be the
      -- metatable for instances of the struct; this is
      -- equivalent to the prototype in js.
      --
      typ.prototype = {__name="methodSet for "..str,
                       __typ = typ,
                       __tostring=function(instance)
                          --print("__tostring called for struct value with typ:")
                          --__st(typ)
                          --print("__tostring has instance:")
                          --__st(instance)

                          local s=typ.__str .. "{";
                          for _,fld in ipairs(typ.fields) do
                             s=s..fld.__name..": "..__dq(instance[fld.__name])..", ";
                          end
                          return s.."}";
                       end,
      }
      typ.prototype.__index = typ.prototype
      
      
      local ctor = function(structTarget, ...)
         local this={};
         --print("top of pointer-to-struct ctor, this="..tostring(this).."; typ.__constructor = "..tostring(typ.__constructor))
         --__st(structTarget, "structTarget")
         local args = {...}
         --__st(args, "args to ctor after structTarget")

         --print("callstack:")
         --print(debug.traceback())
         
         this.__get = function() return structTarget; end;
         this.__set = function(v) typ.copy(structTarget, v); end;
         this.__typ = typ.ptr
         this.__target = structTarget
         this.__val = structTarget -- or should this be this.__val = this?
         this.__name = "__pointerToStructValue"
         setmetatable(this, typ.ptr.prototype)
         return this;
      end
      typ.ptr = __newType(4, __kindPtr, "*" .. str, false, pkg, exported, ctor);
      -- __newType sets typ.comparable = true
      __dfsGlobal:addChild(typ.ptr, typ)
      
      -- pointers have their own method sets, but *T can call elem methods in Go.
      typ.ptr.elem = typ;
      
      typ.ptrToNewlyConstructed = function(...)
         -- built a new struct from scratch, return a pointer to it.
         local structValue = typ(...)
         return typ.ptr(structValue)
      end
      
      typ.ptr.prototype = {__name="methodSet for "..typ.ptr.__str,
                           __typ = typ.ptr,
                           
                           __tostring=function(instance)
                              --print("in pointer-to-struct __tostring, with instance=")
                              --__st(instance,"instance")
                              -- refer out to the value __tostring

                              -- avoid infinite loop...
                              if instance == instance.target then
                                 return("<avoid inf loop>");
                              end
                              return "&" .. typ.prototype.__tostring(instance.__target)
                           end,
                           __index = function(this, k)
                              --print("struct.ptr.prototype.__index called, k='"..k.."'")
                              --print(debug.traceback())
                              -- check methodsets first, then fields.
                              -- check *T:
                              local meth = typ.ptr.prototype[k]
                              if meth ~= nil then
                                 return meth
                              end
                              -- check T:
                              meth = typ.prototype[k]
                              if meth ~= nil then
                                 return meth
                              end
                              -- default to fields on __val
                              return this.__val[k]
                           end,
                           __newindex = function(this, k, v)
                              --print("struct.ptr.prototype.__newindex called, k='"..k.."'")
                              this.__val[k] = v
                           end,
                           
      }
      
      -- incrementally expand the method set. Full
      -- signature details are passed in det.
      
      -- a) for pointer
      typ.ptr.__addToMethods=function(det)
         --print("typ.ptr.__addToMethods called, existing methods:")
         --__st(typ.ptr.methods, "typ.ptr.methods")
         --__st(det, "det")
         if typ.ptr.methods == nil then
            typ.ptr.methods={}
         end
         table.insert(typ.ptr.methods, det)
      end

      -- b) for struct
      typ.__addToMethods=function(det)
         --print("typ.__addToMethods called, existing methods:")
         --__st(typ.methods, "typ.methods")
         --__st(det, "det")
         if typ.methods == nil then
            typ.methods={}
         end
         table.insert(typ.methods, det)
      end
      
      -- __kindStruct.init is here:
      typ.init = function(pkgPath, fields)
         --print("top of init() for struct, fields=")
         --for i, f in pairs(fields) do
         --__st(f, "field #"..tostring(i))
         --__st(f.__typ, "typ of field #"..tostring(i))
         --end
         
         typ.pkgPath = pkgPath;
         typ.fields = fields;
         for i,f in ipairs(fields) do
            if not f.__typ.comparable then
               typ.comparable = false;
               break;
            end
         end
         typ.keyFor = function(x)
            local val = x.__val;
            return __mapAndJoinStrings("_", fields, function(f)
                                          return string.gsub(tostring(f.__typ.keyFor(val[f.__prop])), "\\", "\\\\")
            end)
         end;
         typ.copy = function(dst, src)
            --print("top of typ.copy for structs, here is dst then src:")
            --__st(dst, "dst")
            --__st(src, "src")
            --print("fields:")
            --__st(fields,"fields")
            for _, f in ipairs(fields) do
               local sw2 = f.__typ.kind
               
               if sw2 == __kindArray or
               sw2 ==  __kindStruct then 
                  f.__typ.copy(dst[f.__prop], src[f.__prop]);
               else
                  --print("copying field '"..f.__prop.."'")
                  --__st(dst, "dst prior to copy")
                  dst[f.__prop] = src[f.__prop];
                  --__st(dst, "dst after copy")
               end
            end
         end;
         
         --print("jea debug: on __kindStruct: set .copy on typ to .copy=", typ.copy)
         -- /* nil value */
         local properties = {};

         for i,f in ipairs(fields) do
            properties[f.__prop] = { get= __throwNilPointerError, set= __throwNilPointerError };
         end;
         typ.ptr.__nil = {} -- Object.create(constructor.prototype,s properties);
         --if constructor ~= nil then
         --   constructor(typ.ptr.__nil)
         --end
         typ.ptr.__nil.__val = typ.ptr.__nil;
         -- /* methods for embedded fields */
         __addMethodSynthesizer(function()
               local synthesizeMethod = function(target, m, f)
                  if target.prototype[m.__prop] ~= nil then return; end
                  target.prototype[m.__prop] = function()
                     local v = this.__val[f.__prop];
                     if f.__typ == __jsObjectPtr then
                        v = __jsObjectPtr(v);
                     end
                     if v.__val == nil then
                        local w = {}
                        f.__typ(w, v);
                        v = w
                     end
                     return v[m.__prop](v, arguments);
                  end;
               end;
               for i,f in ipairs(fields) do
                  if f.anonymous then
                     for _, m in ipairs(__methodSet(f.__typ)) do
                        synthesizeMethod(typ, m, f);
                        synthesizeMethod(typ.ptr, m, f);
                     end;
                     for _, m in ipairs(__methodSet(__ptrType(f.__typ))) do
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
   if kind == __kindBool then
      typ.zero = function() return false; end;

   elseif kind ==__kindMap then
      typ.zero = function() return nil; end;

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

   elseif kind == __kindComplex64 or
   kind == __kindComplex128 then
      local zero = typ(0, 0);
      typ.zero = function() return zero; end;
      
   elseif kind == __kindPtr or
   kind == __kindSlice then
      
      typ.zero = function() return typ.__nil; end;
      
   elseif kind == __kindChan then
      typ.zero = function() return __chanNil; end;
      
   elseif kind == __kindFunc then
      typ.zero = function() return __throwNilPointerError; end;
      
   elseif kind == __kindInterface then
      typ.zero = function() return __ifaceNil; end;
      
   elseif kind == __kindArray then
      
      typ.zero = function()
         --print("in zero() for array...")
         return __newAnyArrayValue(typ.elem, typ.len)
      end;

   elseif kind == __kindStruct then
      typ.zero = function()
         return typ.ptr();
      end;

   else
      error("invalid kind: " .. tostring(kind))
   end

   typ.id = __typeIDCounter;
   __typeIDCounter=__typeIDCounter+1;
   typ.size = size;
   typ.kind = kind;
   typ.__str = str;
   typ.named = named;
   typ.pkg = pkg;
   typ.exported = exported;
   typ.methods = typ.methods or {};
   typ.methodSetCache = nil;
   typ.comparable = true;
   typ.bloom = function()
      print("bloom called for typ:")
      __st(typ)
   end
   --print("*** returning from __newType with typ=")
   --__st(typ)
   return typ;
end

function __methodSet(typ)
   
   --if typ.methodSetCache ~= nil then
   --return typ.methodSetCache;
   --end
   local base = {};

   local isPtr = (typ.kind == __kindPtr);
   --print("__methodSet called with typ=")
   --__st(typ)
   --print("__methodSet sees isPtr=", isPtr)
   
   if isPtr  and  typ.elem.kind == __kindInterface then
      -- jea: I assume this is because pointers to interfaces don't themselves have methods.
      typ.methodSetCache = {};
      return {};
   end

   local myTyp = typ
   if isPtr then
      myTyp = typ.elem
   end
   local current = {{__typ= myTyp, indirect= isPtr}};

   -- the Go spec says:
   -- The method set of the corresponding pointer type *T is
   -- the set of all methods declared with receiver *T or T
   -- (that is, it also contains the method set of T).
   
   local seen = {};

   --print("top of while, #current is", #current)
   while #current > 0 do
      local next = {};
      local mset = {};
      
      for _,e in pairs(current) do
         --print("e from pairs(current) is:")
         --__st(e,"e")
         --__st(e.__typ,"e.__typ")
         if seen[e.__typ.__str] then
            --print("already seen "..e.__typ.__str.." so breaking out of match loop")
            break
         end
         seen[e.__typ.__str] = true;
         
         if e.__typ.named then
            --print("have a named type, e.__typ.methods is:")
            --__st(e.__typ.methods, "e.__typ.methods")
            for _, mthod in pairs(e.__typ.methods) do
               --print("adding to mset, mthod = ", mthod)
               table.insert(mset, mthod);
            end
            if e.indirect then
               for _, mthod in pairs(__ptrType(e.__typ).methods) do
                  --print("adding to mset, mthod = ", mthod)
                  table.insert(mset, mthod)
               end
            end
         end
         
         -- switch e.__typ.kind
         local knd = e.__typ.kind
         
         if knd == __kindStruct then
            
            for i,f in ipairs(e.__typ.fields) do
               if f.anonymous then
                  local fTyp = f.__typ;
                  local fIsPtr = (fTyp.kind == __kindPtr);
                  local ty 
                  if fIsPtr then
                     ty = fTyp.elem
                  else
                     ty = fTyp
                  end
                  table.insert(next, {__typ=ty, indirect= e.indirect or fIsPtr});
               end;
            end;
            
            
         elseif knd == __kindInterface then
            
            for _, mthod in pairs(e.__typ.methods) do
               --print("adding to mset, mthod = ", mthod)
               table.insert(mset, mthod)
            end
         end
      end;

      -- above may have made duplicates, now dedup
      --print("at dedup, #mset = " .. tostring(#mset))
      for _, m in pairs(mset) do
         --print("m is ")
         --__st(m,"m")
         if base[m.__name] == nil then
            base[m.__name] = m;
         end
      end;
      --print("after dedup, base for typ '"..typ.__str.."' is ")
      --__st(base)
      
      current = next;
   end
   
   typ.methodSetCache = {};
   table.sort(base)
   for _, detail in pairs(base) do
      table.insert(typ.methodSetCache, detail)
   end;
   return typ.methodSetCache;
end;


__type__.bool    = __newType( 1, __kindBool,    "bool",     true, "", false, nil);
__type__.int = __newType( 8, __kindInt,     "int",   true, "", false, nil);
__type__.int8    = __newType( 1, __kindInt8,    "int8",     true, "", false, nil);
__type__.int16   = __newType( 2, __kindInt16,   "int16",    true, "", false, nil);
__type__.int32   = __newType( 4, __kindInt32,   "int32",    true, "", false, nil);
__type__.int64   = __newType( 8, __kindInt64,   "int64",    true, "", false, nil);
__type__.uint    = __newType( 8, __kindUint,    "uint",     true, "", false, nil);
__type__.uint8   = __newType( 1, __kindUint8,   "uint8",    true, "", false, nil);
__type__.uint16  = __newType( 2, __kindUint16,  "uint16",   true, "", false, nil);
__type__.uint32  = __newType( 4, __kindUint32,  "uint32",   true, "", false, nil);
__type__.uint64  = __newType( 8, __kindUint64,  "uint64",   true, "", false, nil);
__type__.uintptr = __newType( 8, __kindUintptr,    "uintptr",  true, "", false, nil);
__type__.float32 = __newType( 8, __kindFloat32,    "float32",  true, "", false, nil);
__type__.float64 = __newType( 8, __kindFloat64,    "float64",  true, "", false, nil);
__type__.complex64  = __newType( 8, __kindComplex64,  "complex64",   true, "", false, nil);
__type__.complex128 = __newType(16, __kindComplex128, "complex128",  true, "", false, nil);
__type__.string  = __newType(16, __kindString,  "string",   true, "", false, nil);
--__type__.unsafePointer = __newType( 8, __kindUnsafePointer, "unsafe.Pointer", true, "", false, nil);

__ptrType = function(elem)
   if elem == nil then
      error("internal error: cannot call __ptrType() will nil elem")
   end
   local typ = elem.ptr;
   if typ == nil then
      typ = __newType(4, __kindPtr, "*" .. elem.__str, false, "", elem.exported, nil);
      __dfsGlobal:addChild(typ, elem)
      elem.ptr = typ;
      typ.init(elem);
      
   end
   return typ;
end;

__newDataPointer = function(data, constructor)
  --print("__newDataPointer called")
   --   if constructor.elem.kind == __kindStruct then
   --      print("struct recognized in __newDataPointer")
   --      return data;
   --   end
   return constructor(function() return data; end, function(v) data = v; end, data);
end;

__indexPtr = function(array, index, constructor)
   array.__ptr = array.__ptr  or  {};
   local a = array.__ptr[index]
   if a ~= nil then
      return a
   end
   a = constructor(function() return array[index]; end, function(v) array[index] = v; end);
   array.__ptr[index] = a
   return a
end;


__arrayTypes = {};
__arrayType = function(elem, len)
   local typeKey = elem.id .. "_" .. len;
   local typ = __arrayTypes[typeKey];
   if typ == nil then
      typ = __newType(24, __kindArray, "[" .. len .. "]" .. elem.__str, false, "", false, nil);
      __arrayTypes[typeKey] = typ;
      __dfsGlobal:addChild(typ, elem)
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
      __dfsGlobal:addChild(typ, elem)
      typ.init(elem, sendOnly, recvOnly);
      
   end
   return typ;
end;

-- return the (un-named so as to be interoperable)
-- reflect Type that corresponds to tsys type 'typ'.
function __gijitTypeToGoType(typ)

   local kstring = __kind2str[typ.kind]
   local rtyp = __rtyp[kstring]
   if rtyp ~= nil then
      -- basic type, return straight away
      return rtyp
   end
   -- recurse to construct un-named/compound types

   if kind ==  __kindPtr then
      return reflect.PtrTo(__gijitTypeToGoType(typ.elem))
   
   elseif kind ==  __kindSlice then
      return reflect.SliceOf(__gijitTypeToGoType(typ.elem))   
      
   elseif kind ==  __kindArray then
      return reflect.ArrayOf(typ.len, __gijitTypeToGoType(typ.elem))
      
   elseif kind ==  __kindChan then
      local dir = 3 -- both by default
      if typ.sendOnly then
         dir = 2
      elseif typ.recvOnly then
         dir = 1
      end
      return reflect.ChanOf(dir, __gijitTypeToGoType(typ.elem))
      
   elseif kind ==  __kindMap then 
      return reflect.MapOf(__gijitTypeToGoType(typ.key), __gijitTypeToGoType(typ.elem))

      --- TODO: finish the rest

   elseif kind ==  __kindFunc then 
      error("TODO: finish func types")

   elseif kind ==  __kindInterface then 
      error("TODO: finish interface types")
            
   elseif kind ==  __kindStruct then
      error("TODO: finish struct types")
      for i,fld in ipairs(fields) do
         this[fld.__prop] = args[i] or fld.__typ.zero();
      end         
   else
      error("invalid kind: " .. tostring(kind));
   end  
end

__theNilChan={}

function __Chan(elem, capacity, elemReflectType)
   if elem == nil then
      return __theNilChan
   end
   --print("__Chan called")
   --print(debug.traceback())
   local dir = 3 -- direction: 1=recv, 2=send, 3=both.
   local elemty = __gijitTypeToGoType(elem)
   local chtype = reflect.ChanOf(dir, elemty)
   local ch = reflect.MakeChan(chtype, capacity)
   
   local this = {}
   this.__native = ch

   -- gopherJS stuff below
   if capacity < 0  or  capacity > 2147483647 then
      __throwRuntimeError("makechan: size out of range");
   end
   this.elem = elem;
   this.__capacity = capacity;
   this.__buffer = {};
   this.__sendQueue = {};
   this.__recvQueue = {};
   this.__closed = false;
   this.__val = this
   return this
end;


function __Chan_GopherJS(elem, capacity)
   local this = {}
   if capacity < 0  or  capacity > 2147483647 then
      __throwRuntimeError("makechan: size out of range");
   end
   this.elem = elem;
   this.__capacity = capacity;
   this.__buffer = {};
   this.__sendQueue = {};
   this.__recvQueue = {};
   this.__closed = false;
   this.__val = this -- jea add, should it be here?
   return this
end;
__chanNil = __Chan(nil, 0);
__chanNil.__recvQueue = { length= 0, push= function()end, shift= function() return nil; end, indexOf= function() return -1; end; };
__chanNil.__sendQueue = __chanNil.__recvQueue

-- parentTyp should be a typ, we will take parent
-- before calling __addChild.
function __addChildTypesHelper(parentTyp, array)
   __mapArray(array, function(ty)
                 __dfsGlobal:addChild(parentTyp, ty)
   end)
end


__funcTypes = {};
__funcType = function(params, results, variadic)

   -- example: func f(a int, b string) (string, uint32) {}
   --   would have typeKey:
   -- "parm_1,16__results_16,9__variadic_false"
   --
   local typeKey = "parm_" .. __mapAndJoinStrings(",", params, function(p)
                                                     if p.id == nil then
                                                        error("no id for p=",p);
                                                     end;
                                                     return p.id;
                                                 end) .. "__results_" .. __mapAndJoinStrings(",", results, function(r)
                                                                                                if r.id == nil then
                                                                                                   error("no id for r=",r);
                                                                                                end;
                                                                                                return r.id;
                                                                                            end) .. "__variadic_" .. tostring(variadic);
   --print("typeKey is '"..typeKey.."'")
   local typ = __funcTypes[typeKey];
   if typ == nil then
      local paramTypeNames = __mapArray(params, function(p) return p.__str; end);
      if variadic then
         paramTypeNames[#paramTypeNames - 1] = "..." .. paramTypeNames[#paramTypeNames - 1].substr(2);
      end
      local str = "func(" .. table.concat(paramTypeNames, ", ") .. ")";
      
      if #results == 1 then
         str = str .. " " .. results[1].__str;
      elseif #results > 1 then
         str = str .. " (" .. __mapAndJoinStrings(", ", results, function(r) return r.__str; end) .. ")";
      end
      
      typ = __newType(4, __kindFunc, str, false, "", false, nil);
      __funcTypes[typeKey] = typ;

      -- note the dependencies of the new function type
      __addChildTypesHelper(typ, params)
      __addChildTypesHelper(typ, results)

      typ.init(params, results, variadic);
      
   end
   return typ;
end;

--- interface types here

function __interfaceStrHelper(m)
   local s = ""
   if m.pkg ~= "" then
      s = m.pkg .. "."
   end
   return s .. m.__name .. string.sub(m.__typ.__str, 6) -- sub for removing "__kind"
end

__interfaceTypes = {};
__interfaceType = function(methods)
   
   local typeKey = __mapAndJoinStrings("_", methods, function(m)
                                          return m.pkg .. "," .. m.__name .. "," .. m.__typ.id;
   end)
   local typ = __interfaceTypes[typeKey];
   if typ == nil then
      local str = "interface {}";
      if #methods ~= 0 then
         str = "interface { " .. __mapAndJoinStrings("; ", methods, __interfaceStrHelper) .. " }"
      end
      typ = __newType(8, __kindInterface, str, false, "", false, nil);
      __interfaceTypes[typeKey] = typ;

      -- note dependencies
      __mapArray(methods, function(m)
                    __dfsGlobal:addChild(typ, m.__typ)
                    -- should be redundant b/c m.__typ already added these:
                    --__addChildTypesHelper(typ, m.__typ.params)
                    --__addChildTypesHelper(typ, m.__typ.results)
      end)
      
      typ.init(methods);
      
   end
   return typ;
end;
__type__.emptyInterface = __interfaceType({});
__ifaceNil = {};
__error = __newType(8, __kindInterface, "error", true, "", false, nil);
__error.init({{__prop= "Error", __name= "Error", __pkg= "", __typ= __funcType({}, {__type__.string}, false) }});

__mapTypes = {};
__mapType = function(key, elem, mType)
   if key.id == nil then
      print("key.id was nil in __mapType. trace:")
      print(debug.traceback())
   end
   local typeKey = key.id .. "_" .. elem.id;
   local typ = __mapTypes[typeKey];
   if typ == nil then

      -- moved inside __newType, unified. so pass nil to __newType ctor (last param)
--       local ctor =function(entries)
--          print("map ctor called, v = "..tostring(v))
--          local this={};
--          local typ = __mapTypes[typeKey]
--          this.__typ = typ         
--          this.__val = {}; --no meta names, so clean. No accidental collisions.
-- 
--          local kff = key.keyFor
--          this.nilKeyStored = false
--          
--          local len=0;
--          for k, e in pairs(entries) do
--             if k == nil then
--                this.nilKeyStored = true;
--                this.nilValue = e or __intentionalNilValue;
--             else
--                local key = kff(k);
--                --print("using key ", key, " for k=", k)
--                this.__val[key] = e or __intentionalNilValue;
--             end
--             len=len+1
--          end
--          
--          this.len=len
--          this.keyType=key
--          this.elemType=elem
--          this.zeroValue = elem.zero()
--          
--          setmetatable(this, __valueMapMT)
--          return this;
--       end;
      
      typ = __newType(8, __kindMap, "map[" .. key.__str .. "]" .. elem.__str, false, "", false, nil);
      __mapTypes[typeKey] = typ;

      __dfsGlobal:addChild(typ, key)
      __dfsGlobal:addChild(typ, elem)
      
      typ.init(key, elem);
      
   end
   return typ;
end;

-- stored as map value in place of nil, so
-- we can recognized stored nil values in maps.
__intentionalNilValue = {}
__start_new_map_iter_after_nil={}
__map_iter_with_nil_key={}

__valueMapMT = {
   __name = "__valueMapMT",

   __newindex = function(t, k, v)
      local len = t.len
      --print("map newindex called for key", k, " len at start is ", len)

      if k == nil then
         if t.nilKeyStored then
            -- replacement, no change in len.
         else
            -- new key
            t.len = len + 1
            t.nilKeyStored = true
         end
         t.nilValue = v
         return
      end

      -- invar: k is not nil

      local ks = t.__typ.elem.keyFor(k)
      --local ks = tostring(k)
      if v ~= nil then
         if t.__val[ks] == nil then
            -- new key
            t.len = len + 1
         end
         t.__val[ks] = v
         return

      else
         -- invar: k is not nil. v is nil.

         if t.__val[ks] == nil then
            -- new key
            t.len = len + 1
         end
         t.__val[ks] = __intentionalNilValue
         return
      end
      --print("len at end of newindex is ", len)
   end,
   
   __index = function(t, k)
      -- Instead of __index,
      -- use __call('get', ...) for two valued return and
      --  proper zero-value return upon not present.
      -- __index only ever returns one value[1].
      -- reference: [1] http://lua-users.org/lists/lua-l/2007-07/msg00182.html
      
      --print("__index called for key", k)
      if k == nil then
         if t.nilKeyStored then
            return t.nilValue
         else
            -- TODO: if zeroValue wasn't set (e.g. in __makeMap), then
            -- set it somehow...
            return t.zeroValue
         end
      end
      
      -- k is not nil.

      local ks = t.__typ.elem.keyFor(k)      
      --local ks = tostring(k)
      
      local val = t.__val[ks]
      if val == __intentionalNilValue then
         return nil
      end
      return val
   end,
   
   __tostring = function(t)
      --print("__tostring for map called")
      local len = t.len
      local s = "map["..t.keyType.__str.. "]"..t.elemType.__str.."{"
      local r = t.__val
      
      local vquo = ""
      if len > 0 and t.elemType.__str == "string" then
         vquo = '"'
      end
      local kquo = ""
      if len > 0 and t.keyType.__str == "string" then
         kquo = '"'
      end
      
      for i, v in pairs(r) do
         s = s .. kquo..tostring(i)..kquo.. ": " .. vquo..tostring(v) ..vquo.. ", "
      end
      return s .. "}"
   end,
   
   __len = function(t)
      -- this does get called by the # operation(!)
      -- print("len called")
      return t.len
   end,

   __pairs = function(t)
      print("map __pairs called!")
      -- this makes a map work in a for k,v in pairs() do loop.

      -- Iterator function takes the table and an index and returns the next index and associated value
      -- or nil to end iteration

      local function stateless_iter(t, k)
         
         print("map stateless_iter called, with k ="..tostring(k))
         __st(k, "k")

         --if k == __map_iter_with_nil_key then
         --   t.nilKeyStored
         --end
         if k == __start_new_map_iter_after_nil then
            return next(t.__val, nil)
         end
         
         --  Implement your own key,value selection logic in place of next
         
         -- we need to account for the fact that a nil key can be stored.
         -- a) how to start
         -- b) how to replace the __intentionalNilValue value and keep going
         -- c) how to finish.
         -- d) how to turn stored string keys back into their key type. Actually
         --    this is handled by front end generated code, so we don't need
         --    to worry about that here.

         -- early termination is a problem: if we have nil key in our map,
         -- then the Lua for loop will terminate. So we'll miss that key.
         -- So I think we need to write our while loop.
         
         -- for now, just get a basic iteration going.
         
         local nextKey
         local nextVal
         nextKey, nextVal = next(t.__val, k)
         print("map iter, back from next(t.__val, k), is")

         __st(nextKey, "nextKey")
         __st(nextVal, "nextVal")
         
         return nextKey, nextVal
      end
      
      -- Return an iterator function, the table, starting point (nil
      -- tells next to start).
      return stateless_iter, t, nil
   end,

   __call = function(t, ...)
      --print("__call() invoked, with ... = ", ...)
      local oper, k, zeroVal = ...
      --print("oper is", oper)
      --print("key is ", k)

      -- we use __call('get', k, zeroVal) instead of __index
      -- so that we can return multiple values
      -- to match Go's "a, ok := mymap[k]" call.
      
      if oper == "get" then

         --print("get called for key", k)
         if k == nil then
            if t.nilKeyStored then
               return t.nilValue, true;
            else
               -- key not present returns the zero value for the value.
               return zeroVal, false;
            end
         end
         
         -- k is not nil.
         local ks = tostring(k)      
         local val = t.__val[ks]
         if val == __intentionalNilValue then
            --print("val is the __intentionalNilValue")
            return nil, true;

         elseif val == nil then
            -- key not present
            --print("key not present, zeroVal=", zeroVal)
            --for i,v in pairs(t.__val) do
            --   print("debug: i=", i, "  v=", v)
            --end
            return zeroVal, false;
         end
         
         return val, true
         
      elseif oper == "delete" then

         -- the hash table delete operation
         
         if k == nil then

            if t.nilKeyStored then
               t.nilKeyStored = false
               t.nilVaue = nil
               t.len = t.len -1
            end

            --print("len at end of delete is ", t.len)              
            return
         end

         local ks = tostring(k)           
         if t.__val[ks] == nil then
            -- key not present
            return
         end
         
         -- key present and key is not nil
         t.__val[ks] = nil
         t.len = t.len - 1
         
         --print("len at end of delete is ", t.len)
      end
   end
   
}

__makeMap = function(entries, keyType, elemType, mType)
   local mty = __mapType(keyType, elemType, mType)
   local m = mty(entries);   
   --__st(m, "m in __makeMap")
   return m
end;


-- __basicValue2kind: identify type of basic value
--   or return __kindUnknown if we don't recognize it.
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
         return __kindUnknown;
         --error("__basicValue2kind: unhandled cdata cty: '"..tostring(cty).."'")
      end      
   elseif ty == "boolean" then
      return __kindBool;
   elseif ty == "number" then
      return __kindFloat64
   elseif ty == "string" then
      return __kindString
   end
   
   return __kindUnknown;
   --error("__basicValue2kind: unhandled ty: '"..ty.."'")   
end

__sliceType = function(elem)
   --print("__sliceType called with elem = ", elem)
   if elem == nil then
      print(debug.traceback())
      error "__sliceType called with nil elem!"
   end
   local typ = elem.slice;
   if typ == nil then
      typ = __newType(24, __kindSlice, "[]" .. elem.__str, false, "", false, nil);
      elem.slice = typ;
      __dfsGlobal:addChild(typ, elem)
      typ.init(elem);
      
   end
   return typ;
end;

__makeSlice = function(typ, length, capacity)
   --print("__makeSlice called with type length='"..type(length).."'")
   length = length or 0
   capacity = capacity or length
   
   length = tonumber(length)
   --print("in __makeSlice: after tonumber, length is: '"..tostring(length).."'")
   
   if capacity == nil then
      capacity = length
   else
      capacity = tonumber(capacity)
   end
   if length < 0  or  length > 9007199254740992 then -- 2^53
      __throwRuntimeError("makeslice: len out of range");
   end
   if capacity < 0  or  capacity < length  or  capacity > 9007199254740992 then
      __throwRuntimeError("makeslice: cap out of range: "..tostring(capcity));
   end
   local array = __newAnyArrayValue(typ.elem, capacity)
   local slice = typ(array);
   slice.__length = length;
   return slice;
end;




function __field2strHelper(f)
   local tag = ""
   if f.__tag ~= "" then
      tag = string.gsub(f.__tag, "\\", "\\\\")
      tag = string.gsub(tag, "\"", "\\\"")
   end
   return f.__name .. " " .. f.__typ.__str .. tag
end

function __typeKeyHelper(f)
   return f.__name .. "," .. f.__typ.id .. "," .. f.__tag;
end

__structTypes = {};
__structType = function(pkgPath, fields)
   local typeKey = __mapAndJoinStrings("_", fields, __typeKeyHelper)

   local typ = __structTypes[typeKey];
   if typ == nil then
      local str
      if #fields == 0 then
         str = "struct {}";
      else
         str = "struct { " .. __mapAndJoinStrings("; ", fields, __field2strHelper) .. " }";
      end
      
      typ = __newType(0, __kindStruct, str, false, "", false, function(...)
                         local this = {}
                         local args = {...}
                         local typ = __structTypes[typeKey]
                         typ.finishStructValueCreation(this, fields, args)
                         return this
      end);
      __structTypes[typeKey] = typ;

      __mapArray(fields, function(f)
                    __dfsGlobal:addChild(typ, f.__typ)
      end)
      
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
      
      for i,f in ipairs(typ.fields) do
         if  not __equal(a[f.__prop], b[f.__prop], f.__typ) then
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

__interfaceIsEqual = function(a, b)
   --print("top of __interfaceIsEqual! a is:")
   --__st(a,"a")
   --print("top of __interfaceIsEqual! b is:")   
   --__st(b,"b")
   if a == nil or b == nil then
      --print("one or both is nil")
      if a == nil and b == nil then
         --print("both are nil")
         return true
      else
         --print("one is nil, one is not")
         return false
      end
   end
   if a == __ifaceNil  or  b == __ifaceNil then
      --print("one or both is __ifaceNil")
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
      ok = value.__typ == typ;
   else
      local valueTypeString = value.__typ.__str;

      -- this caching doesn't get updated as methods
      -- are added, so disable it until fixed, possibly, in the future.
      --ok = typ.implementedBy[valueTypeString];
      ok = nil
      if ok == nil then
         ok = true;
         local valueMethodSet = __methodSet(value.__typ);
         local interfaceMethods = typ.methods;
         --print("valueMethodSet is")
         --__st(valueMethodSet)
         --print("interfaceMethods is")
         --__st(interfaceMethods)

         __ipairsZeroCheck(interfaceMethods)
         __ipairsZeroCheck(valueMethodSet)
         for _, tm in ipairs(interfaceMethods) do            
            local found = false;
            for _, vm in ipairs(valueMethodSet) do
               --print("checking vm against tm, where tm=")
               --__st(tm)
               --print("and vm=")
               --__st(vm)
               
               if vm.__name == tm.__name  and  vm.pkg == tm.pkg  and  vm.__typ == tm.__typ then
                  --print("match found against vm and tm.")
                  found = true;
                  break;
               end
            end
            if  not found then
               --print("match *not* found for tm.__name = '"..tm.__name.."'")
               ok = false;
               typ.missingMethodFor[valueTypeString] = tm.__name;
               break;
            end
         end
         typ.implementedBy[valueTypeString] = ok;
      end
      if not ok then
         missingMethod = typ.missingMethodFor[valueTypeString];
      end
   end
   --print("__assertType: after matching loop, ok = ", ok)
   
   if not ok then
      if returnTuple then
         if isInterface then
            return nil, false
         end
         return typ.zero(), false
      end
      local msg = ""
      if value ~= __ifaceNil then
         msg = value.__typ.__str
      end
      --__panic(__packages["runtime"].TypeAssertionError.ptr("", msg, typ.__str, missingMethod));
      error("type-assertion-error: could not '"..msg.."' -> '"..typ.__str.."', missing method '"..missingMethod.."'")
   end
   
   if not isInterface then
      value = value.__val;
   end
   if typ == __jsObjectPtr then
      value = value.object;
   end
   if returnTuple then
      return value, true
   end
   return value
end;

__stackDepthOffset = 0;
__getStackDepth = function()
   local err = Error(); -- new
   if err.stack == nil then
      return nil;
   end
   return __stackDepthOffset + #err.stack.split("\n");
end;

-- possible replacement for ipairs.
-- starts at a[0] if it is present.
function __zipairs(a)
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

-- __elim0 is a helper, to get rid of t[0], and
-- shift everything up in the returned array
-- that will start at [1], if non-empty/non-nil.
--
-- If t.__len is available, we assume that this
-- is a Go slices or arrays, that starts at 0
-- if len > 0.
--
-- Otherwise, we assume an array starting at either [0] or [1],
-- with no 'nil' holes in the middle.
--
function __elim0(t)
   if type(t) ~= 'table' then
      return t
   end

   if t == nil then
      return
   end

   -- is __len available?
   local mt = getmetatable(t)
   if mt ~= nil and rawget(mt, "__len") ~= nil then
      --print("__len found!")
      -- Go slice/array, from 0.
      local n = #t
      local r = {}
      for i=0,n-1 do
         table.insert(r, t[i])
      end
      return r
   end
   
   -- can we leave t unchanged?
   local z = t[0]
   if z == nil then
      return t
   end
   
   local r = {}
   table.insert(r, z)
   local i = 1
   while true do
      local v = t[i]
      if v == nil then
         break
      else
         table.insert(r, v)
      end
      i=i+1
   end
   return r
end

function __unpack0(t)
   if type(t) ~= 'table' then
      return t
   end
   if t == nil then
      return
   end
   return unpack(__elim0(t))
end

local __lazyEllipsisMT = {
   __call  = function(self)
      return self.__val
   end,
}

function __lazy_ellipsis(t)
   local r = {
      __name = "__lazy_ellipsis_instance",
      __val = t,
   }
   setmetatable(r, __lazyEllipsisMT)
   return r
end

function __printHelper(v)

   local tv = type(v)
   if tv == "string" then
      print("\""..v.."\"") -- used to be backticks
   elseif tv == "table" then
      if v.__name == "__lazy_ellipsis_instance" then
         local expand = v()
         for _,c in pairs(expand) do
            __printHelper(c)
         end
         return
      end
   end
   print(v)
end

function __gijit_printQuoted(...)
   local a = {...}
   --print("__gijit_printQuoted called, a = " .. tostring(a), " len=", #a)
   if a[0] ~= nil then
      __printHelper(a[0])
   end
   for _,v in ipairs(a) do
      __printHelper(v)
   end
end

-- last thing, so we store all types/vars defined so far
__storeBuiltins()


