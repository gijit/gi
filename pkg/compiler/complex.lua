--
-- complex number support
--

local ffi = require("ffi")

-- complex128 and complex64 are Go predefined types
local complex128=ffi.typeof("complex double") -- aka "complex". re and im are each float64.
local complex64=ffi.typeof("complex float")   -- re and im are each float32

local ffiNew=ffi.new
local ffiIsType=ffi.istype

-- provide Go's builtin complex constructor.
local function complex(re, im)
   if ffiIsType("complex", re) then
      if im ~= nil then
         error("bad input to complex: with first arg complex, 2nd arg must be nil")
      end
      return re
   end
   return ffiNew("complex",re or 0,im or 0)
end

-- real is a Go builtin, returning the real part of z.
local function real(z)
   if ffiIsType("complex", z) then
      return z.re
   end
   if type(z)=="number" then
      return z
   end
   return 0
end

-- imag is a Go builtin, returning the imaginary part of z.
local function imag(z)
   if ffiIsType("complex", z) then
      return z.im
   end
   return 0
end

-- for speed, make local versions

local type=type
local select=select
local tonumber=tonumber
local tostring=tostring

local e=math.exp(1)
local pi=math.pi
local abs=math.abs
local exp=math.exp
local log=math.log
local cos=math.cos
local sin=math.sin
local cosh=math.cosh
local sinh=math.sinh
local sqrt=math.sqrt
local atan2=math.atan2

local function cexp(a)
   return e^a
end

local function conj(c)
   return complex(real(c), -imag(c))
end

-- 
local function mod(a)
   local ra, ia = real(a), imag(a)
   if ia == 0 then
      return ra
   end
   return sqrt(ra*ra + ia*ia)
end

-- arg is the angle between the positive real
-- axis to the line joining the point to the origin;
-- also known as an argument of the point.
-- a.k.a phase
local function arg(a)
   a=complex(a)
   return atan2(imag(a), real(a))
end

-- returns two values: r, theta; giving the polar coordinates of c.
local function polar(c)
   return mod(c), arg(c)
end

-- convert from polar coordinates to a complex number
-- where the real and imag parts naturally provide rectangular
-- coordinates. e.g. r*exp(i*theta) -> x+iy, where
-- x is r*cos(theta), and y is r*sin(theta).
--
local function rect(r, theta)
   return complex(r*cos(theta), r*sin(theta))
end


-- clog computes the complex natural log, double precision,
-- with branch cut along the negative real axis.
-- The natural logarithm of a complex number z
-- with polar coordinate components (r, θ) equals
-- ln r + i(θ+2nπ), with the principal value ln r + iθ.
-- 
local function clog(a) -- 
   local ra, ia = real(a), imag(a)   
   return complex(log(ra*ra + ia*ia)/2, atan2(ia,ra))
end

-- clog computes the complex natural log, single precision,
-- with branch cut along the negative real axis.
local function clogf(a)
   local ra, ia = real(a), imag(a)
   return complex64(log(ra*ra + ia*ia)/2, atan2(ia,ra))
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
   end,

   __unm=function(a)
      return complex(-real(a),-imag(a))
   end,
   
   __tostring=function(c)
      return real(c).."+"..imag(c).."i"
   end,

   __pow=function(a,b)
      local ra,ia = real(a), imag(a)
      local rb,ib = real(b), imag(b)
      local alensq=ra*ra+ia*ia
      if alensq==0 then
         if rb==0 and ib==0 then
            return complex(1, 0)
         end
         return complex(0, 0)
      end
      local theta=atan2(ia, ra)
      return rect(alensq^(rb/2)*exp(-ib*theta),ib*log(alensq)/2+rb*theta)
   end
   
}

-- can only be done once, so we'll detect and skip
-- any 2nd import.
if not __cxMT_already then
   ffi.metatype(complex128, __cxMT)
   __cxMT_already = true
end

-- exports
_G.complex=complex
_G.complex128=complex128
_G.complex64=complex64
_G.real=real
_G.imag=imag
local cmplx = {
   conj=conj,
   mod=mod,
   arg=arg,
   theta=arg,
   polar=polar,
   rect=rect,
   cexp=cexp,
   clog=clog
}
_G.cmplx=cmplx
