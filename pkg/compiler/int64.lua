-- int64 uint64 helpers

local ffi = require("ffi")

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

--
-- complex number support
--
--

-- complex128 and complex64 are Go predefined types
complex128=ffi.typeof("complex double") -- aka "complex". re and im are each float64.
complex64=ffi.typeof("complex float")   -- re and im are each float32

local ffiNew=ffi.new
local ffiIstype=ffi.istype

-- provide Go's builtin complex constructor.
function complex(re, im)
   return ffiNew("complex",re or 0,im or 0)
end

-- real is a Go builtin, returning the real part of z.
function real(z)
   if ffiIstype("complex", z) then
      return z.re
   end
   if type(z)=="number" then
      return z
   end
   return 0
end

-- imag is a Go builtin, returning the imaginary part of z.
function imag(z)
   if ffiIstype("complex", z) then
      return z.im
   end
   return 0
end

-- the metatable for complex number arithmetic.
local __cxMT={
   __add=function(a, b)
      return complex(real(a)+real(b),imag(a)+imag(b))
   end,
   
   __sub=function(a, b)
      return complex(real(a)-real(b),imag(a)-imag(b))
   end,
   
   __mul=function(a,b)
      local ra,ia=real(a),imag(a)
      local rb,ib=real(b),imag(b)
      return complex(ra*rb - ia*ib, ra*ib + rb*ia)
   end,
   
   __div=function(a,b)
      local ra,ia=real(a),imag(a)
      local rb,ib=real(b),imag(b)
      local denom=rb*rb + ib*ib
      return complex((ra*rb+ia*ib)/denom, (rb*ia-ra*ib)/denom)
   end
   
}

-- can only be done once, so we'll detect and skip
-- any 2nd import.
if not __cxMT_already then
   ffi.metatype(complex128, __cxMT)
   __cxMT_already = true
end
