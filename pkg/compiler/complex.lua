--
-- complex number support
--

-- Portions of /usr/local/go/src/math/cmplx
-- in the Go language distribution/Go standard library, are
-- used under the following terms:
-- Copyright 2010 The Go Authors. All rights reserved.
-- Use of this source code is governed by a BSD-style
-- license that can be found in the LICENSE file.
--
-- See the top level LICENSE file for the full text.

local ffi = require("ffi")
local bit = require("bit")

-- complex128 and complex64 are Go predefined types
local complex128=ffi.typeof("complex double") -- aka "complex". re and im are each float64.
local complex64=ffi.typeof("complex float")   -- re and im are each float32

local ffiNew=ffi.new
local ffiIsType=ffi.istype

local function __truncateToInt(x)
   if x >= 0 then
       return x - (x % 1)
   end
   return x + (-x % 1)
end


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
local Inf=math.huge

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
local function cpolar(c)
   return cabs(c), carg(c)
end

-- convert from polar coordinates to a complex number
-- where the real and imag parts naturally provide rectangular
-- coordinates. e.g. r*exp(i*theta) -> x+iy, where
-- x is r*cos(theta), and y is r*sin(theta).
--
local function crect(r, theta)
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
      return crect(alensq^(rb/2)*exp(-ib*theta),ib*log(alensq)/2+rb*theta)
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
   Conj=conj,
   Abs=cabs,
   Arg=carg,
   Exp=cexp,
   Log=clog,
   Polar=cpolar,
   Rect=crect,
   Sqrt=csqrt
}


function cmath.Sin(c)
	local r,i=real(c),imag(c)
	return complex(sin(r)*cosh(i),cos(r)*sinh(i))
end
function cmath.Cos(c)
	local r,i=real(c),imag(c)
	return complex(cos(r)*cosh(i),sin(r)*sinh(i))
end
function cmath.Tan(c)
	local r,i=2*real(c),2*imag(c)
	local div=cos(r)+cosh(i)
	return complex(sin(r)/div,sinh(i)/div)
end

-- Program to subtract nearest integer multiple of PI
function reducePi(x) -- takes float64, returns float64
   -- extended precision value of PI:
   local DP1 = 3.14159265160560607910E0   -- ?? 0x400921fb54000000
   local DP2 = 1.98418714791870343106E-9  -- ?? 0x3e210b4610000000
   local DP3 = 1.14423774522196636802E-17 -- ?? 0x3c6a62633145c06e
	t = x / pi
	if t >= 0 then
		t = t + 0.5
	else
		t = t - 0.5
	end
    t = __truncateToInt(t) -- int64(t) = the multiple
	return ((x - t*DP1) - t*DP2) - t*DP3
end

-- Taylor series expansion for cosh(2y) - cos(2x)
function tanSeries(z) -- takes complex128, returns float64
   local MACHEP = 1.0 / tonumber(bit.lshift(1LL, 53))
   local x = abs(2 * real(z))
   local y = abs(2 * imag(z))
   x = reducePi(x)
   x = x * x
   y = y * y
   local x2=1
   local y2=1
   local f =1
   local rn = 0
   local d = 0
   while true do
      rn=rn+1
      f = f*rn
      rn=rn+1
      f=f*rn
      x2 = x2 * x
      y2 = y2 * y
      local t = y2 + x2
      t=t/f
      d=d+t
      
      rn=rn+1
      f=f*rn
      rn=rn+1
      f=f*rn
      x2 = x2 * x
      y2 = y2*y
      t = y2 - x2
      t = t/f
      d=d+t
      if not (abs(t/d) > MACHEP) then
         -- Caution: Use `not` and > instead of <= for correct behavior if t/d is NaN.
         -- See golang issue 17577.
         break
      end
   end
   return d
end

-- Complex circular cotangent
--
-- DESCRIPTION:
--
-- If
--     z = x + iy,
--
-- then
--
--           sin 2x  -  i sinh 2y
--     w  =  --------------------.
--            cosh 2y  -  cos 2x
--
-- On the real axis, the denominator has zeros at even
-- multiples of PI/2.  Near these points it is evaluated
-- by a Taylor series.
--
-- ACCURACY:
--
--                      Relative error:
-- arithmetic   domain     # trials      peak         rms
--    DEC       -10,+10      3000       6.5e-17     1.6e-17
--    IEEE      -10,+10     30000       9.2e-16     1.2e-16
-- Also tested by ctan * ccot = 1 + i0.

-- Cot returns the cotangent of x.
function cmath.Cot(x)
   local xr, xi = real(x), imag(x)
   local d = cosh(2*xi) - cos(2*xr)
	if abs(d) < 0.25 then
		d = tanSeries(x)
	end
	if d == 0 then
		return Inf
	end
	return complex(sin(2*xr)/d, -sinh(2*xi)/d)
end

function cmath.Sinh(c)
	local r,i=real(c),imag(c)
	return complex(cos(i)*sinh(r),sin(i)*cosh(r))
end
function cmath.Cosh(c)
	local r,i=real(c),imag(c)
	return complex(cos(i)*cosh(r),sin(i)*sinh(r))
end
function cmath.Tanh(c)
	local r,i=2*real(c),2*imag(c)
	local div=cos(i)+cosh(r)
	return complex(sinh(r)/div,sin(i)/div)
end

-- inverse trig functions

function cmath.Asin(c)
   return i*clog(i*c+(1-c^2)^0.5)
end
function cmath.Acos(c)
	return pi/2+i*clog(i*c+(1-c^2)^0.5)
end
function cmath.Atan(c)
	local r2,i2=real(c),imag(c)
	local c3,c4=complex(1-i2,r2),complex(1+r2^2-i2^2,2*r2*i2)
	return complex(arg(c3/c4^0.5),-clog(cmath.abs(c3)/cmath.abs(c4)^0.5))
end
function cmath.Atan2(c2,c1)--y,x
	local r1,i1,r2,i2=real(c1),imag(c1),real(c2),imag(c2)
	if r1==0 and i1==0 and r2==0 and i2==0 then
		return 0
	end
	local c3,c4=complex(r1-i2,i1+r2),complex(r1^2-i1^2+r2^2-i2^2,2*(r1*i1+r2*i2))
	return complex(arg(c3/c4^0.5),-clog(cmath.abs(c3)/cmath.Abs(c4)^0.5))
end

function cmath.Asinh(c)
	return clog(c+(1+c^2)^0.5)
end
function cmath.Acosh(c)
	return 2*clog((c-1)^0.5+(c+1)^0.5)-log(2)
end
function cmath.Atanh(c)
	return (clog(1+c)-clog(1-c))/2
end

-- complex base logarithm. log(b,z) gives log_b(z),
-- which is clog(z)/clog(b), with base b.
--
function cmath.Log(b, z)
   
	local br, bi = real(b), imag(b)
	local zr, zi = real(z), imag(z)
        
	local qr = log(br*br+bi*bi)/2
        local qi = atan2(bi,br)
        
	local sr = log(zr*zr+zi*zi)/2
        local si = atan2(zi,zr)
        
	local denom=qr*qr+qi*qi
	return complex((sr*qr+si*qi)/denom, (qr*si-sr*qi)/denom)
end

cmath.Pow = __cxMT.__pow


-- exports
_G.complex=complex
_G.complex128=complex128
_G.complex64=complex64
_G.real=real
_G.imag=imag
_G.cmath=cmath

-- test

local acosTestVals = {
   (1.0017679804707456328694569 - 2.9138232718554953784519807i),
   (0.03606427612041407369636057 + 2.7358584434576260925091256i),
   (1.6249365462333796703711823 + 2.3159537454335901187730929i),
   (2.0485650849650740120660391 - 3.0795576791204117911123886i),
   (0.29621132089073067282488147 - 3.0007392508200622519398814i),
   (1.0664555914934156601503632 - 2.4872865024796011364747111i),
   (0.48681307452231387690013905 - 2.463655912283054555225301i),
   (0.6116977071277574248407752 - 1.8734458851737055262693056i),
   (1.3649311280370181331184214 + 2.8793528632328795424123832i),
   (2.6189310485682988308904501 - 2.9956543302898767795858704i)
}

local expectedCot = {
   (0.005333839942314019+0.9975187770930108i),
   (0.0006110458528512657-1.0084212802655375i),
   (-0.002064200607561421-0.9808251268240298i),
   (-0.0034444657331057713+0.9975566188309715i),
   (0.002775424851723473+1.0041112277051858i),
   (0.011609912676285432+0.9925920872057117i),
   (0.012081615289134921+1.0081095190379168i),
   (0.04506207267921228+1.015185741448182i),
   (0.0025108438745688137-0.9942304877145374i),
   (-0.004336960053730652+1.0025022601350426i)
}

function cdiff(a,b)
   return cabs(a-b) > 0.0000001
end

for i, v in ipairs(acosTestVals) do
   --print("on i = ", i)
   if cdiff(cmath.Cot(v), expectedCot[i]) then
      error("difference at Cot case i="..i)
   end
end
