-- int64 uint64 helpers

local ffi = require("ffi")

if jit.os == "Windows" then
   ffi.cdef[[
   long long int _atoi64(const char *nptr);
   ]]
   __atoll=ffi.C._atoi64
else
   ffi.cdef[[
   long long int atoll(const char *nptr);
   ]]
   __atoll=ffi.C.atoll   
end

-- assume 64-bit int and uint
int = ffi.typeof(0LL)
uint=ffi.typeof(0ULL)

int64=ffi.typeof("int64_t")
uint64=ffi.typeof("uint64_t")

int32 = ffi.typeof("int32_t")
uint32 = ffi.typeof("uint32_t")

int16 = ffi.typeof("int16_t")
uint16 = ffi.typeof("uint16_t")

int8 = ffi.typeof("int8_t")
uint8 = ffi.typeof("uint8_t")
byte = uint8

float64 = ffi.typeof("double")
float32 = ffi.typeof("float")

-- to display floats, use: tonumber() to convert to float64 that lua can print.

--MinInt64: -9223372036854775808
--MaxInt64: 9223372036854775807

-- a=int(-9223372036854775808LL)

-- to use cdata as hash keys... tostring() to make them strings first.

----------------------------------
----------------------------------
--
-- ffi dependent string/byte stuff
--
----------------------------------
----------------------------------

-- __newByteArray is created from vals, which
-- can be a table, an empty table, or
-- a string. If it is a string, then to
-- quote from the ffi docs at http://luajit.org/ext_ffi_semantics.html
--
-- "Byte arrays may also be initialized with
--  a Lua string. This copies the whole string
--  plus a terminating zero-byte. The copy stops
--  early only if the array has a known, fixed size."
--
function __newByteArray(vals)
   vals = vals or {}
   local sz = #vals
   local res= {
      __bytes = ffi.new("char["..sz.."]", vals),
      __sz=sz,
      __name="__valueByteArray",
   }
   setmetatable(res, {
                   __index = function(me, i)
                      return me.__bytes[i]
                   end,
                   __len=function(me)
                      --print("__length on byteArray called")
                      return me.__sz
                   end,
                   __tostring=function(me)
                      --print("__tostring on byteArray called")
                      return string.sub(ffi.string(me.__bytes), 1, me.__sz)
                   end,
   })
   return res
end

__stringToBytes = function(str)
   return __newByteArray(str)
end;

__bytesToString = function(ba)
   return tostring(ba)
end;
