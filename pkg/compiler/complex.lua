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
   elseif ffiIsType("complex float", z) then
      return float32(z.re)
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
   elseif ffiIsType("complex float", z) then
      return float32(z.im)
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
local i=complex(0,1)

local function cexp(a)
   return e^a
end

local function conj(c)
   return complex(real(c), -imag(c))
end

-- complex absolute value, also known
-- as modulus, magnitude, or norm.
local function cabs(a)
   local ra, ia = real(a), imag(a)
   if ia == 0 then
      return ra
   end
   return sqrt(ra*ra + ia*ia)
end

-- carg is the angle between the positive real
-- axis to the line joining the point to the origin;
-- also known as an argument of the point.
-- a.k.a phase; If no errors occur, returns
-- the phase angle of z in the interval [−π; π].
local function carg(z)
   z=complex(z)
   return atan2(imag(z), real(z))
end

-- returns two values: r, theta; giving the polar coordinates of c.
local function polar(c)
   return cabs(c), carg(c)
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
local function clog(a)
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

local function csqrt(c)
   return complex(c)^0.5
end

-- can only be done once, so we'll detect and skip
-- any 2nd import.
if not __cxMT_already then
   ffi.metatype(complex128, __cxMT)
   __cxMT_already = true
end

-- cmath library functions
local cmath = {
   conj=conj,
   cabs=cabs,
   carg=carg,
   cexp=cexp,
   clog=clog,
   polar=polar,
   rect=rect,
   csqrt=csqrt
}


function cmath.sin(c)
	local r,i=real(c),imag(c)
	return complex(sin(r)*cosh(i),cos(r)*sinh(i))
end
function cmath.cos(c)
	local r,i=real(c),imag(c)
	return complex(cos(r)*cosh(i),sin(r)*sinh(i))
end
function cmath.tan(c)
	local r,i=2*real(c),2*imag(c)
	local div=cos(r)+cosh(i)
	return complex(sin(r)/div,sinh(i)/div)
end

function cmath.sinh(c)
	local r,i=real(c),imag(c)
	return complex(cos(i)*sinh(r),sin(i)*cosh(r))
end
function cmath.cosh(c)
	local r,i=real(c),imag(c)
	return complex(cos(i)*cosh(r),sin(i)*sinh(r))
end
function cmath.tanh(c)
	local r,i=2*real(c),2*imag(c)
	local div=cos(i)+cosh(r)
	return complex(sinh(r)/div,sin(i)/div)
end

-- inverse trig functions

function cmath.asin(c)
   return i*clog(i*c+(1-c^2)^0.5)
end
function cmath.acos(c)
	return pi/2+i*clog(i*c+(1-c^2)^0.5)
end
function cmath.atan(c)
	local r2,i2=re(c),im(c)
	local c3,c4=complex(1-i2,r2),complex(1+r2^2-i2^2,2*r2*i2)
	return complex(arg(c3/c4^0.5),-clog(cmath.abs(c3)/cmath.abs(c4)^0.5))
end
function cmath.atan2(c2,c1)--y,x
	local r1,i1,r2,i2=re(c1),im(c1),re(c2),im(c2)
	if r1==0 and i1==0 and r2==0 and i2==0 then--Indeterminate
		return 0
	end
	local c3,c4=complex(r1-i2,i1+r2),complex(r1^2-i1^2+r2^2-i2^2,2*(r1*i1+r2*i2))
	return complex(arg(c3/c4^0.5),-clog(cmath.abs(c3)/cmath.abs(c4)^0.5))
end

function cmath.asinh(c)
	return clog(c+(1+c^2)^0.5)
end
function cmath.acosh(c)
	return 2*clog((c-1)^0.5+(c+1)^0.5)-log(2)
end
function cmath.atanh(c)
	return (clog(1+c)-clog(1-c))/2
end

-- complex base logarithm. log(b,z) gives log_b(z),
-- which is clog(z)/clog(b), with base b.
--
function cmath.log(b, z)
   
	local br, bi = real(b), imag(b)
	local zr, zi = real(z), imag(z)
        
	local qr = log(br*br+bi*bi)/2
        local qi = atan2(bi,br)
        
	local sr = log(zr*zr+zi*zi)/2
        local si = atan2(zi,zr)
        
	local denom=qr*qr+qi*qi
	return complex((sr*qr+si*qi)/denom, (qr*si-sr*qi)/denom)
end

cmath.pow = __cxMT.__pow


-- exports
_G.complex=complex
_G.complex128=complex128
_G.complex64=complex64
_G.real=real
_G.imag=imag
_G.cmath=cmath
