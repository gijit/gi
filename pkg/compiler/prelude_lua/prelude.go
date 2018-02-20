package prelude

const Prelude = `

-------- ../complex.lua -------
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
local asin=math.asin
local acos=math.acos
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


-- phase returns the phase (also called the argument) of x.
-- The returned value is in the range [-Pi, Pi].
--
-- It is the angle between the positive real
-- axis to the line joining the point to the origin;
-- also known as an argument of the point.
--
-- If no errors occur, returns
-- the phase angle of z in the interval [−π; π].
--
local function phase(x)
   x=complex(x)
   return atan2(imag(x), real(x))
end

-- returns two values: r, theta; giving the polar coordinates of c.
local function polar(c)
   return cabs(c), phase(c)
end

-- rect returns the complex number x with polar coordinates r, θ.
-- i.e.
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

-- clogf computes the complex natural log, single precision,
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
   Conj=conj,
   Abs=cabs,
   Phase=phase,
   Exp=cexp,
   Log=clog,
   Polar=polar,
   Rect=rect,
   Sqrt=csqrt
}


function cmath.Sin(c)
	local r,i=real(c),imag(c)
	return complex(sin(r)*cosh(i),cos(r)*sinh(i))
end
function cmath.Cos(c)
	local r,i=real(c),imag(c)
	return complex(cos(r)*cosh(i),-sin(r)*sinh(i))
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
         -- Caution: Use "not" and > instead of <= for correct behavior if t/d is NaN.
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

-- Complex circular arc sine
--
-- DESCRIPTION:
--
-- Inverse complex sine:
--                               2
-- w = -i clog( iz + csqrt( 1 - z ) ).
--
-- casin(z) = -i casinh(iz)
--

-- Asin returns the inverse sine of x.
function cmath.Asin(x)
   local xr = real(x)
   local xi = imag(x)
	if xi == 0 then
		if abs(xr) > 1 then
			return complex(pi/2, 0) -- DOMAIN error
		end
		return complex(asin(xr), 0)
	end
	local ct = complex(-xi, xr) -- i * x
	local xx = x * x
	local x1 = complex(1-real(xx), -imag(xx)) -- 1 - x*x
	local x2 = csqrt(x1)                       -- x2 = sqrt(1 - x*x)
	local w = clog(ct + x2)
	return complex(imag(w), -real(w)) -- -i * w
end

function cmath.Acos(c)
	return pi/2+i*clog(i*c+(1-c^2)^0.5)
end
function cmath.Atan(c)
	local r2,i2=real(c),imag(c)
	local c3,c4=complex(1-i2,r2),complex(1+r2^2-i2^2,2*r2*i2)
	return complex(phase(c3/c4^0.5),-clog(cabs(c3)/cabs(c4)^0.5))
end
function cmath.Atan2(c2,c1) -- y, x
	local r1,i1,r2,i2=real(c1),imag(c1),real(c2),imag(c2)
	if r1==0 and i1==0 and r2==0 and i2==0 then
		return 0
	end
	local c3,c4=complex(r1-i2,i1+r2),complex(r1*r1 - i1*i1 + r2*r2 - i2*i2, 2*(r1*i1 + r2*i2))
	return complex(phase(c3/c4^0.5),-clog(cabs(c3)/cabs(c4)^0.5))
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
function cmath.ComplexLog(b, z)
   
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

-- Specifically do not import as cmplx, so that we can 
-- allow the Go library to exist side-by-side for testing/comparison.
-- _G.cmplx=cmath

--[[
-- tests

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

function check(a, b)
   if cabs(a-b) > 0.0000001 then
      error("difference!")
   end
end

for i, v in ipairs(acosTestVals) do
   --print("on i = ", i)
   check(cmath.Cot(v), expectedCot[i])
end

check(cmath.Exp((1.0017679804707456-2.9138232718554953i)), (-2.652761299626814-0.6148879956088891i))
print("cmath.Exp checks")

check(cmath.Conj((1.0017679804707456-2.9138232718554953i)), (1.0017679804707456+2.9138232718554953i))
check(cmath.Abs((1.0017679804707456-2.9138232718554953i)), 3.081218127024294)
check(cmath.Phase((1.0017679804707456-2.9138232718554953i)), -1.2396569035824907)
check(cmath.Log((1.0017679804707456-2.9138232718554953i)), (1.1253250145847473-1.2396569035824907i))
check(cmath.Sqrt((1.0017679804707456-2.9138232718554953i)), (1.4288082634655777-1.019669099893085i))
check(cmath.Sin((1.0017679804707456-2.9138232718554953i)), (7.784589065071272-4.949771659620112i))
check(cmath.Cos((1.0017679804707456-2.9138232718554953i)), (4.979011924883674+7.738872474578103i))
check(cmath.Tan((1.0017679804707456-2.9138232718554953i)), (0.005360254415745048-1.0024587328321108i))
check(cmath.Cot((1.0017679804707456-2.9138232718554953i)), (0.005333839942314019+0.9975187770930108i))
check(cmath.Sinh((1.0017679804707456-2.9138232718554953i)), (-1.1475081545902999-0.3489051537979257i))
check(cmath.Cosh((1.0017679804707456-2.9138232718554953i)), (-1.5052531450365143-0.26598284181096343i))
check(cmath.Tanh((1.0017679804707456-2.9138232718554953i)), (0.7789713818473163+0.0941450495765785i))
check(cmath.Asin((1.0017679804707456-2.9138232718554953i)), (0.315902142176508-1.8389632402722567i))
check(cmath.Acos((1.0017679804707456-2.9138232718554953i)), (1.2548941846183885+1.8389632402722567i))
check(cmath.Atan((1.0017679804707456-2.9138232718554953i)), (1.4549738117535138-0.3130322073097648i))
check(cmath.Asinh((1.0017679804707456-2.9138232718554953i)), (1.797481526353204-1.2223981995549138i))
check(cmath.Acosh((1.0017679804707456-2.9138232718554953i)), (1.8389632402722567-1.2548941846183885i))
check(cmath.Atanh((1.0017679804707456-2.9138232718554953i)), (0.0966478569558242-1.2701291676466084i))
print("cmath.Atanh checks")

local b = (0.48681307452231387690013905 - 2.463655912283054555225301i)
local z = (2.6189310485682988308904501 - 2.9956543302898767795858704i)
local l = cmath.ComplexLog(b, z)
check(b^l, z)
print("cmath.ComplexLog checks")

-- end tests
--]]
-----------------------
-------- ../defer.lua -------
-- deferinit.lua : global setup for defer handling

-- utility: table show
function __ts(t)
   if t == nil then
      return "<nil>"
   end
   local s = "<non-nil table:>\n"
   local k = 0
   for i,v in pairs(t) do
      s = s .. "key:" .. tostring(i) .. " -> val:" .. tostring(v) .. "\n"
      k = k +1
   end
   if k > 0 then
      return s
   end
   return "<non-nil but empty table with 0 entries>: " .. tostring(t)
end

-- can we have one global definition of panic and recover?
-- This would be preferred to repeating them in every function.

-- string viewing of panic value
__recovMT = {__tostring = function(v) return 'a-panic-value:' .. tostring(v[1]) end}

-- __recoverVal will be nill if no panic,
--              or if panic happened and
--              then was recovered.
-- this is always a table to avoid
--  stringification problems. The real
--  panic value is inside at position [1].
--
-- NB __recoverVal  needs to be per-goroutine. As each
--  could be unwinding independently at any
--  point in time.

__recoverVal = nil

recover = function() 
    local cp = __recoverVal
    __recoverVal = nil
    return cp;
end

panic = function(err)
  -- wrap err in table to prevent conversion to string by error()
  __recoverVal = {err}
  -- but still allow it to be viewable in a stack trace:
  setmetatable(__recoverVal, __recovMT)
  error(__recoverVal)
end


  -- begin boilerplate part 2:
  
  -- prepare to handle panic/defer/recover
__handler2 = function(err)
     --print(" __handler2 running with err =", err)
     __recoverVal = err
     return err
end
  
__panicHandler = function(err, defers)
       --print("__panicHandler running with err =", err)
       -- print(debug.traceback())
       --print("__panicHandler running with defers:", tostring(defers))

     __recoverVal = err
     if defers ~= nil then

         --print(debug.traceback(), " __panicHandler running with err =", err, " and #defer = ", #defers)      
         --print(" __panicHandler running with err =", err, " and #defer = ", #defers)  
         for __i = #defers, 1, -1 do
             local dcall = {xpcall(defers[__i], __handler2)}
             --for i,v in pairs(dcall) do print("__panicHandler: panic path defer call result: i=",i, "  v=",v) end
         end
     else
         --print("debug: found no defers in __panicHandler")
     end
     --print("__panicHandler: done with defer processing")
     if __recoverVal ~= nil then
        return __recoverVal
     end
  end

  -- __processDefers represents the normal
  --    return path, without a panic.
  --
  --    We need to update the named return values if
  --    there were explicit return values from __actual,
  --    and then we need to call the defers.
  --
  --    __namedNames is an array of the variable names of the return values,
  --                 so we know how to update actEnv.
  --
__processDefers = function(who, defers, __res, __namedNames, actEnv)
  --print(who,": __processDefers top: __res[1] is: ", tostring(__res[1]))
  print(who,": __processDefers top: __namedNames is: ", __ts(__namedNames))

  if __res[1] then
      --print(who,": __processDefers: call had no panic")
      -- call had no panic. run defers with the nil recover

      if #__res > 1 then
         --for k,v in pairs(__res) do print(who, " __processDefers: __res k=", k, " val=", v) end

         -- explicit return, so fill the named vals before defers see them.
         local unp = {table.unpack(__res, 2)}
         --print("unp is: ", tostring(unp))
         for i, k in pairs(__namedNames) do
             actEnv[k] = unp[i]
         end

         --print(who, " __processDefers: post fill: ret0 = ", ret0, " and ret1=", ret1)
      end

      assert(recoverVal == nil)
      for __i = #defers, 1, -1 do
        local dcall = {xpcall(defers[__i], __handler2)}
        for i,v in pairs(dcall) do
            --print(who," __processDefers: normal path defer call result: i=",i, "  v=",v)
        end
      end
  else
      print(who, " __processDefers: checking for panic still un-caught...", __recoverVal)
      -- is there an un-recovered panic that we need to rethrow?
      if __recoverVal ~= nil then
         --print(who, "__processDefers: un recovered error still exists, rethrowing ", __recoverVal)
         error(__recoverVal)
      end
  end

  if #__namedNames == 0 then
     print("__processDefers: #__namedNames was 0, no returns")
     return nil
  end
  -- put the named return values in order
  local orderedReturns={}
  for i, k in pairs(__namedNames) do
     print("debug: fetching from function env k=",k," which we see has value ", actEnv[k], "in actEnv", tostring(actEnv))
     orderedReturns[i] = actEnv[k]
  end
  local debug= true
  if debug then
     print("orderedReturns is len ", #orderedReturns)
     for i,v in pairs(orderedReturns) do
        print(who," __processDefers: orderedReturns: i=",i, "  v=",v)
     end
  end
  return unpack(orderedReturns)
end


__actuallyCall = function(who, __actual, __namedNames, __zeroret, __defers, __orig)

   --local actEnv = getfenv(__actual)
   -- So getfenv(__actual) showed that actEnv
   -- was the _G global env, not good.
   -- To fix this, we give f its own env,
   -- so that named return variables can
   -- be written/read from this env.
   
   local actEnv = {}
   local mt = {
      __index = _G, -- read through to globals.
      __newindex = _G, -- write to closure-capture globals too.
   }
   setmetatable(actEnv,mt)
   setfenv(__actual, actEnv)

  for i,k in pairs(__namedNames) do
     --print("filling actEnv[k='"..tostring(k).."'] = '"..tostring(actEnv[k]).."' with __zeroret[i='"..tostring(i).."']='",tostring(__zeroret[i]),"'")
     actEnv[k] = __zeroret[i]
  end  
  local myPanic = function(err) __panicHandler(err, __defers) end
  local __res = {xpcall(__actual, myPanic, unpack(__orig))}
  return __processDefers(who, __defers, __res,  __namedNames, actEnv)  
end
-----------------------
-------- ../dfs.lua -------
--  depend.lua:
--
--  Implement Depth-First-Search (DFS)
--  on the graph of depedencies
--  between types. A pre-order
--  traversal will print
--  leaf types before the compound
--  types that need them defined.

local __dfsTestMode = false

function __newDfsNode(self, name, typ)
   if typ == nil then
      error "typ cannot be nil in __newDfsNode"
   end
   if not __dfsTestMode then
      if typ.__str == nil then
         print(debug.traceback())
         error "typ must be typ, in __newDfsNode"
      end
      -- but we won't know the kind until
      -- later, since this may be early in
      -- typ construction.
   end
   
   local nd = self.dfsDedup[typ]
   if nd ~= nil then
      return nd
   end
   local node= {
      visited=false,
      children=false, -- lazily put in table, better printing.
      dedupChildren={},
      id = self.dfsNextID,
      name=name,
      typ=typ,
   }
   self.dfsNextID=self.dfsNextID+1
   self.dfsDedup[typ] = node
   table.insert(self.dfsNodes, node)

   --print("just added to dfsNodes node "..name)
   --__st(typ, "typ, in __newDfsNode")
   
   self.stale = true
   
   return node
end

function __isBasicTyp(typ)
   if typ == nil or
      typ.kind == nil or
   typ.named then
      return false
   end
   
   -- we can skip all basic types,
   -- as they are already defined.
   --
   if typ.kind <= 16 or -- __kindComplex128
      typ.kind == 24 or -- __kindString
   typ.kind == 26 then  -- __kindUnsafePointer
      return
   end
end

-- par should be a node; e.g. typ.__dfsNode
function __addChild(self, parTyp, chTyp)

   if parTyp == nil then
      error "parTyp cannot be nil in __addChild"
   end
   if chTyp == nil then
      print(debug.traceback())
      error "chTyp cannot be nil in __addChild"
   end
   if not __dfsTestMode then
      if parTyp.__str == nil then
         print(debug.traceback())
         error "parTyp must be typ, in __addChild"
      end
      if chTyp.__str == nil then
         print(debug.traceback())
         error "chTyp must be typ, in __addChild"
      end
   end

   -- we can skip all basic types,
   -- as they are already defined.   
   if __isBasicTyp(chTyp) then
      return
   end
   if __isBasicTyp(parTyp) then
      error("__addChild error: parent was basic type. "..
               "cannot add child to basic typ ".. parType.__str)
   end

   local chNode = self.dfsDedup[chTyp]
   if chNode == nil then
      -- child was previously generated, so
      -- we don't need to worry about this
      -- dependency
      return
   end
   
   local parNode = self.dfsDedup[parTyp]
   if parNode == nil then
      parNode = self:newDfsNode(parTyp.__str, parTyp)
   end
   
   if parNode.dedupChildren[ch] ~= nil then
      -- avoid adding same child twice.
      return
   end

   -- In Lua both nil and the boolean
   -- value false represent false in
   -- a logical expression
   --
   if not parNode.children then
      -- we lazily instantiate children
      -- for better diagnostics. Its
      -- much clearer to see "children = false"
      -- than "children = 'table: 0x04e99af8'"
      parNode.children = {}
   end
   
   local pnc = #parNode.children
   parNode.dedupChildren[chNode]= pnc+1
   table.insert(parNode.children, chNode)
   self.stale = true
end

function __markGraphUnVisited(self)
   self.dfsOrder = {}
   for _,n in ipairs(self.dfsNodes) do
      n.visited = false
   end
   self.stale = false
end

function __emptyOutGraph(self)
   self.dfsOrder = {}
   self.dfsNodes = {} -- node stored in value.
   self.dfsDedup = {} -- payloadTyp key -> node value.
   self.dfsNextID = 0
   self.stale = false
end

function __dfsHelper(self, node)
   if node == nil then
      return
   end
   if node.visited then
      return
   end
   node.visited = true
   __st(node,"node, in __dfsHelper")
   if node.children then
      for _, ch in ipairs(node.children) do
         self:dfsHelper(ch)
      end
   end
   print("post-order visit sees node "..tostring(node.id).." : "..node.name)
   table.insert(self.dfsOrder, node)
end

function __showDFSOrder(self)
   if self.stale then
      self:doDFS()
   end
   for i, n in ipairs(self.dfsOrder) do
      print("dfs order "..i.." is "..tostring(n.id).." : "..n.name)
   end
end

function __doDFS(self)
   __markGraphUnVisited(self)
   for _, n in ipairs(self.dfsNodes) do
      self:dfsHelper(n)
   end
   self.stale = false
end

function __hasTypes(self)
   return self.dfsNextID ~= 0
end


function __NewDFSState()
   return {
      dfsNodes = {},
      dfsOrder = {},
      dfsDedup = {},
      dfsNextID = 0,

      doDFS = __doDFS,
      dfsHelper = __dfsHelper,
      reset = __emptyOutGraph,
      newDfsNode = __newDfsNode,
      addChild = __addChild,
      markGraphUnVisited = __markGraphUnVisited,
      hasTypes = __hasTypes,
      showDFSOrder=__showDFSOrder,
   }
end

--[[
-- test. To test, change the --[[ above to ---[[
--       and issue dofile('dfs.lua')
dofile 'tutil.lua' -- must be in prelude dir to test.

function __testDFS()
   __dfsTestMode = true
   local s = __NewDFSState()

   -- verify that reset()
   -- works by doing two passes.
   
   for i =1,2 do
      s:reset()
      
      local aPayload = {}
      local a = s:newDfsNode("a", aPayload)
   
      local adup = s:newDfsNode("a", aPayload)
      if adup ~= a then
          error "dedup failed."
      end

      local b = s:newDfsNode("b", {})
      local c = s:newDfsNode("c", {})
      local d = s:newDfsNode("d", {})
      local e = s:newDfsNode("e", {})
      local f = s:newDfsNode("f", {})

      -- separate island:
      local g = s:newDfsNode("g", {})
      
      s:addChild(a, b)

      -- check dedup of child add
      local startCount = #a.children
      s:addChild(a, b)
      if #a.children ~= startCount then
          error("child dedup failed.")
      end

      s:addChild(b, c)
      s:addChild(b, d)
      s:addChild(d, e)
      s:addChild(d, f)

      s:doDFS()

      s:showDFSOrder()

      expectEq(s.dfsOrder[1], c)
      expectEq(s.dfsOrder[2], e)
      expectEq(s.dfsOrder[3], f)
      expectEq(s.dfsOrder[4], d)
      expectEq(s.dfsOrder[5], b)
      expectEq(s.dfsOrder[6], a)
      expectEq(s.dfsOrder[7], g)
   end
   
end
__testDFS()
__testDFS()
--]]
-----------------------
-------- ../int64.lua -------
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
-----------------------
-------- ../map.lua -------
-- a Lua virtual table system suitable for use in arrays and maps

-- design using the _giPrivateMapRaw index was suggested
-- by https://www.lua.org/pil/13.4.4.html
-- To intercept all writes, the requirement is that
-- the table always be empty. Hence the user uses
-- and nearly empty proxy. The only thing the proxy
-- has in it are a pointer to the actual data,
-- and a len counter.

-- create private index
_giPrivateMapRaw = {}
_giPrivateMapProps = {}

-- stored as map value in place of nil, so
-- we can recognized stored nil values in maps.
_intentionalNilValue = {}

 _giPrivateMapMt = {

    __newindex = function(t, k, v)
       --print("newindex called for key", k, " len at start is ", len)

       local props = t[_giPrivateMapProps]
       local len = props.len

       if k == nil then
          if props.nilKeyStored then
             -- replacement, no change in len.
          else
             -- new key
             props.len = len + 1
             props.nilKeyStored = true
          end
          props.nilValue = v
          return
       end

       -- invar: k is not nil

       local ks = tostring(k)
       if v ~= nil then
          if t[_giPrivateMapRaw][ks] == nil then
             -- new key
             props.len = len + 1
          end
          t[_giPrivateMapRaw][ks] = v
          return

       else
          -- invar: k is not nil. v is nil.

          if t[_giPrivateMapRaw][ks] == nil then
             -- new key
             props.len = len + 1
          end
          t[_giPrivateMapRaw][ks] = _intentionalNilValue
          return
      end
      --print("len at end of newidnex is ", len)
    end,

    __index = function(t, k)
       -- Instead of __index,
       -- use __call('get', ...) for two valued return and
       --  proper zero-value return upon not present.
       -- __index only ever returns one value[1].
       -- reference: [1] http://lua-users.org/lists/lua-l/2007-07/msg00182.html
              
       --print("__index called for key", k)
       if k == nil then
          local props = t[_giPrivateMapProps]
          if props.nilKeyStored then
             return props.nilValue
          else
             -- TODO: replace nil with zero-value for the value type.
             return nil
          end
       end

       -- k is not nil.

       local ks = tostring(k)       
       local val = t[_giPrivateMapRaw][ks]
       if val == _intentionalNilValue then
          return nil
       end
       return val
    end,

    __tostring = function(t)
       --print("__tostring for _gi_Map called")
       local props = t[_giPrivateMapProps]
       local len = props["len"]
       local s = "map["..props["keyType"].__str.. "]"..props["valType"].__str.."{"
       local r = t[_giPrivateMapRaw]
       
       local vquo = ""
       if len > 0 and props.valType.__str == "string" then
          vquo = '"'
       end
       local kquo = ""
       if len > 0 and props.keyType.__str == "string" then
          kquo = '"'
       end
       
       -- we want to skip both the _giPrivateMapRaw and the len
       -- when iterating, which happens automatically if we
       -- iterate on r, the inside private data, and not on the proxy.
       for i, _ in pairs(r) do

          -- lua style:
          -- s = s .. "["..kquo..tostring(i)..kquo.."]" .. "= " .. vquo..tostring(r[i]) ..vquo.. ", "
          -- Go style
          s = s .. kquo..tostring(i)..kquo.. ": " .. vquo..tostring(r[i]) ..vquo.. ", "
       end
       return s .. "}"
    end,

    __len = function(t)
       -- this does get called by the # operation(!)
       -- print("len called")
       return t[_giPrivateMapProps]["len"]
    end,

    __pairs = function(t)
       -- print("__pairs called!")
       -- this makes a _giMap work in a for k,v in pairs() do loop.

       -- Iterator function takes the table and an index and returns the next index and associated value
       -- or nil to end iteration

       local function stateless_iter(t, k)
           local v
           --  Implement your own key,value selection logic in place of next
           local ks = tostring(k)
           ks, v = next(t[_giPrivateMapRaw], tostring(k))
           if v then return ks,v end
       end

       -- Return an iterator function, the table, starting point
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
              local props = t[_giPrivateMapProps]
              if props.nilKeyStored then
                 return props.nilValue, true;
              else
                 -- key not present returns the zero value for the value.
                 return zeroVal, false;
              end
           end
           
           -- k is not nil.
           local ks = tostring(k)      
           local val = t[_giPrivateMapRaw][ks]
           if val == _intentionalNilValue then
              --print("val is the _intentinoalNilValue")
              return nil, true;

           elseif val == nil then
              -- key not present
              --print("key not present, zeroVal=", zeroVal)
              --for i,v in pairs(t[_giPrivateMapRaw]) do
              --   print("debug: i=", i, "  v=", v)
              --end
              return zeroVal, false;
           end
           
           return val, true
           
        elseif oper == "delete" then

           -- the hash table delete operation

           local props = t[_giPrivateMapProps]              
           local len = props.len
           --print("delete called for key", k, " len at start is ", len)
                      
           if k == nil then

              if props.nilKeyStored then
                 props.nilKeyStored = false
                 props.nilVaue = nil
                 props.len = len -1
              end

              --print("len at end of delete is ", props.len)              
              return
           end

           local ks = tostring(k)           
           if t[_giPrivateMapRaw][ks] == nil then
              -- key not present
              return
           end
           
           -- key present and key is not nil
           t[_giPrivateMapRaw][ks] = nil
           props.len = len - 1
           
           --print("len at end of delete is ", props.len)
        end
    end
 }
 
function _gi_NewMap(keyType, valType, x)
   assert(type(x) == 'table', 'bad parameter #1: must be table')

   local proxy = {}
   proxy["Typeof"]="_gi_Map"
   proxy[_giPrivateMapRaw] = x

   -- get initial count
   local len = 0
   for k, v in pairs(x) do
      len = len + 1
   end

   local props = {len=len, keyType=keyType, valType=valType, nilKeyStored=false}
   proxy[_giPrivateMapProps] = props

   setmetatable(proxy, _giPrivateMapMt)
   return proxy
end;

-----------------------
-------- ../math.lua -------
-- math helper functions

-- x == math.huge   -- test for +inf, inline

-- x == -math.huge  -- test for -inf, inline

-- x ~= x           -- test for nan, inline

-- x > -math.huge and x < math.huge  -- test for finite

-- or their slower counterparts:

math.isnan  = function(x) return x ~= x; end
math.finite = function(x) return x > -math.huge and x < math.huge; end

math.nan = math.huge * 0

__truncateToInt = function(x)
   if x >= 0 then
       return x - (x % 1)
   end
   return x + (-x % 1)
end

__integerByZeroCheck = function(x)
   if not math.finite(x) then
      error("integer divide by zero")
   end
   -- eliminate any fractional part
   if x >= 0 then
       return x - (x % 1)
   end
   return x + (-x % 1)
end

function __max(a,b)
   if a > b then
      return a
   end
   return b
end

function __min(a,b)
   if a < b then
      return a
   end
   return b
end
-----------------------
-------- ../prelude.lua -------
-- prelude defines things that should
-- be available before any user code is run.

function __gi_GetRangeCheck(x, i)
   if i == nil then
      print(debug.traceback())
      error "where is i nil??"
   end
   if x == nil then
      print(debug.traceback())
      error "where is x nil??"
   end
   if x == nil or i < 0 or i >= #x then
      error("index out of range: i="..tostring(i).." vs #x is "..tostring(#x))
  end
  return x[i]
end;

function __gi_SetRangeCheck(x, i, val)
  --print("SetRangeCheck. x=", x, " i=", i, " val=", val)
  if x == nil or i < 0 or i >= #x then
     error("index out of range")
  end
  x[i] = val
  return val
end;

-----------------------
-------- ../rune.lua -------
function __decodeRune(s, i)
   return {__utf8.sub(s, i+1, i+1), 1}
end

-- from gopherjs, ported to use bit ops.
__bit =require("bit")

--[[

-- js op precedence: higher precendence = tighter binding.
--
-- arshift/lshift: 12 left-to-right   
-- band: 9    left-to-right
-- bor : 7    left-to-right


__decodeRune = function(str, pos)
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
     local r = __bit.bor(__bit.lshift(__bit.band(c0, 0x1F), 6), __bit.band(c1, 0x3F));
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
   local r = __bit.bor(__bit.bor(__bit.lshift(__bit.band(c0, 0x0F), 12), __bit.lshift(__bit.band(c1, 0x3F), 6)), __bit.band(c2, 0x3F));
   
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
    local r = __bit.bor(__bit.bor(__bit.bor(__bit.lshift(__bit.band(c0, 0x07),18), __bit.lshift(__bit.band(c1, 0x3F), 12)), __bit.lshift(__bit.band(c2, 0x3F), 6), __bit.band(c3, 0x3F)));

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
    return String.fromCharCode(__bit.bor(0xC0, __bit.arshift(r,6)), __bit.bor(0x80, __bit.band(r, 0x3F)));
  end
  if r <= 0xFFFF then
   return String.fromCharCode(__bit.bor(0xE0, __bit.arshift(r,12)), __bit.bor(0x80, (__bit.band(__bit.arshift(r,6), 0x3F))), __bit.bor(0x80, __bit.band(r, 0x3F)));
  end
   return String.fromCharCode(__bit.bor(0xF0, __bit.arshift(r, 18)), __bit.bor(0x80, __bit.band(__bit.arshift(r,12),0x3F)), __bit.bor(0x80, __bit.band(__bit.arshift(r,6), 0x3F)), __bit.bor(0x80, __bit.band(r, 0x3F)));
end;

--]]
-----------------------
-------- ../string.lua -------

__stringToBytes = function(str)
  local array = Uint8Array(#str);
  for i = 0,#str-1 do
    array[i] = str.charCodeAt(i);
  end
  return array;
end;

--

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

__stringToRunes = function(str)
  local array = Int32Array(#str);
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

__runesToString = function(slice)
  if slice.__length == 0 then
    return "";
  end
  local str = "";
  for i = 0,#slice-1 do
    str = str .. __encodeRune(slice.__array[slice.__offset + i]);
  end
  return str;
end;


__copyString = function(dst, src)
  local n = __min(#src, dst.__length);
  for i = 0,n-1 do
    dst.__array[dst.__offset + i] = src.charCodeAt(i);
  end
  return n;
end;
-----------------------
-------- ../tsys.lua -------
--
-- tsys.lua the type system for gijit.
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

__global ={};
__module ={};
__packages = {}
__idCounter = 0;
__pkg = {};

-- length of array, counting [0] if present,
-- but trusting __len if metamethod avail.
function __lenz(array)
   local n = #array
   local lenmeth = array.__len
   if lenmeth ~= nil then
      return n
   end
   if array[0] ~= nil then
      n=n+1
   end
   return n
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
         --vals = __st(v,"",indent+1,quiet,methods_desc, seen)
      else 
         vals = tostring(v)
      end
      s = s..pre.." "..tostring(k).. " key: '"..tostring(i).."' val: '"..vals.."'\n"
   end
   if k == 0 then
      s = pre.."<empty table> " .. addr
   end

   --local mt = getmetatable(t)
   if mt ~= nil then
      s = s .. "\n"..__st(mt, "mt.of."..name, indent+1, true, methods_desc, seen)
   end
   if not quiet then
      print(s)
   end
   -- restore metamethods
   setmetatable(t, mt)
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
      --return __externalize(fn(this, (__sliceType({},__jsObjectPtr))(__global.Array.prototype.slice.call(arguments, {}))), __type__emptyInterface);
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

__subslice = function(slice, low, high, max)
   if high == nil then
      
   end
   if low < 0  or  (high ~= nil and high < low)  or  (max ~= nil and high ~= nil and max < high)  or  (high ~= nil and high > slice.__capacity)  or  (max ~= nil and max > slice.__capacity) then
      __throwRuntimeError("slice bounds out of range");
   end
   
   local s = {}
   slice.__constructor.tfun(s, slice.__array);
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
      -- fields must be an array for this to work.
      for i=0,#typ.elem.fields-1 do
         helper(typ.elem.fields[i].__prop);
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

   local newSlice ={}
   slice.__constructor.tfun(newSlice, newArray);
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

      if k < 0 or k >= #t then
         error "write to array error: access out-of-bounds"
      end
      
      t.__val[k] = v
   end,
   
   __index = function(t, k)
     --print("__valueArrayMT.__index called, k='"..tostring(k).."'")
      if type(k) == "string" then
            return rawget(t,k)
      elseif type(k) == "table" then
         print("callstack:"..tostring(debug.traceback()))
         error("table as key not supported in __valueArrayMT")
      else
         --__st(t.__val)
         if k < 0 or k >= #t then
            print(debug.traceback())
            error("read from array error: access out-of-bounds; "..tostring(k).." is outside [0, "  .. tostring(#t) .. ")")
         end
        --print("array access bounds check ok.")
         return t.__val[k]
      end
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
      local w = t.__offset + k
      if k < 0 or k >= t.__capacity then
         error "slice error: write out-of-bounds"
      end
      t.__array[w] = v
   end,
   
   __index = function(t, k)
      
     --print("__valueSliceMT.__index called, k='"..tostring(k).."'");
      --__st(t.__val)
     --print("callstack:"..tostring(debug.traceback()))

      if type(k) == "string" then
         --print("we have string key, doing rawget on t")
         --__st(t, "t")
         return rawget(t,k)
      elseif type(k) == "table" then
         print("callstack:"..tostring(debug.traceback()))
         error("table as key not supported in __valueSliceMT")
      else
         local w = t.__offset + k
         if k < 0 or k >= t.__capacity then
            print(debug.traceback())
            error("slice error: access out-of-bounds, k="..tostring(k).."; cap="..tostring(t.__capacity))
         end
        --print("slice access bounds check ok.")
         return t.__array[w]
      end
   end,

   __len = function(t)
     --print("__valueSliceMT metamethod __len called, returning ", t.__length)
      return t.__length
   end,
   
   __tostring = function(self, ...)
     --print("__tostring called from __valueSliceMT")

      local len = tonumber(self.__length) -- convert from LL int
      local off = tonumber(self.__offset)
     --print("__tostring sees self.__length of ", len, " __offset = ", off)
      local cap = self.__capacity
      --local s = "slice <len=" .. tostring(len) .. "; off=" .. off .. "; cap=" .. cap ..  "> is "..self.__constructor.__str.."{"
      local s = self.__constructor.__str.."{"
      local raw = self.__array

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
      
      --print("in __tfunBasicMT, start __st on ...")
      --__st({...}, "__tfunBasicMT.dots")
      --print("in __tfunBasicMT,   end __st on ...")

      --print("in __tfunBasicMT, start __st on self")
      --__st(self, "self")
      --print("in __tfunBasicMT,   end __st on self")

      local newInstance = {}
      if self ~= nil then
         if self.tfun ~= nil then
            --print("calling tfun! -- let constructors set metatables if they wish to. our newInstance is an empty table="..tostring(newInstance))

            -- this makes a difference as to whether or
            -- not the ctor receives a nil 'this' or not...
            -- So *don't* set metatable here, let ctor do it.
            -- setmetatable(newInstance, __valueBasicMT)
            
            -- get zero value if no args
            if #{...} == 0 and self.zero ~= nil then
              --print("tfun sees no args and we have a typ.zero() method, so invoking it")
               self.tfun(newInstance, self.zero())
            else
               self.tfun(newInstance, ...)
            end
         end
      else
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
      local str = t.__str or ""
      return str .. "{" .. tostring(t.__get()) .. "}"
   end
}



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
      typ.tfun = function(this, v)
         this.__val = v;
         setmetatable(this, __valueBasicMT)
      end;
      typ.wrapped = true;
      typ.keyFor = function(x) return tostring(x); end;

   elseif kind == __kindString then
      
      typ.tfun = function(this, v)
         --print("strings' tfun called! with v='"..tostring(v).."' and this:")
         --__st(this)
         this.__val = v;
         setmetatable(this, __valueBasicMT)
      end;
      typ.wrapped = true;
      typ.keyFor = __identity; -- function(x) return "_" .. x; end;

   elseif kind == __kindFloat32 or
   kind == __kindFloat64 then
      
      typ.tfun = function(this, v)
         this.__val = v;
         setmetatable(this, __valueBasicMT)
      end;
      typ.wrapped = true;
      typ.keyFor = function(x) return __floatKey(x); end;


   elseif kind ==  __kindComplex64 then

      typ.tfun = function(this, re, im)
         this.__val = re + im*complex(0,1);
         setmetatable(this, __valueBasicMT)
      end;
      typ.wrapped = true;
      typ.keyFor = function(x) return tostring(x); end;
      
      --    typ.tfun = function(this, real, imag)
      --      this.__real = __fround(real);
      --      this.__imag = __fround(imag);
      --      this.__val = this;
      --    end;
      --    typ.keyFor = function(x) return x.__real .. "_" .. x.__imag; end;

   elseif kind ==  __kindComplex128 then

      typ.tfun = function(this, re, im)
         this.__val = re + im*complex(0,1);
         setmetatable(this, __valueBasicMT)
      end;
      typ.wrapped = true;
      typ.keyFor = function(x) return tostring(x); end;
      
      --     typ.tfun = function(this, real, imag)
      --        this.__real = real;
      --        this.__imag = imag;
      --        this.__val = this;
      --        this.__constructor = typ
      --     end;
      --     typ.keyFor = __identity --function(x) return x.__real .. "_" .. x.__imag; end;
      --    
      
   elseif kind ==  __kindPtr then

      if constructor ~= nil then
        --print("in newType kindPtr, constructor is not-nil: "..tostring(constructor))
      end
      typ.tfun = constructor  or
         function(this, wat, getter, setter, target)
           --print("pointer typ.tfun which is same as constructor called! getter='"..tostring(getter).."'; setter='"..tostring(setter).."; target = '"..tostring(target).."'; this:")
            --__st(this, "this")
            --print("wat, 2nd arg to ctor, is:")
            --__st(wat, "wat")
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
            setmetatable(this, __valuePointerMT)
         end;
      typ.keyFor = __idKey;
      
      typ.init = function(elem)
         --print("init(elem) for pointer type called.")
         __dfsGlobal:addChild(typ, elem)
         typ.elem = elem;
         typ.wrapped = (elem.kind == __kindArray);
         typ.__nil = typ({}, __throwNilPointerError, __throwNilPointerError);
      end;

   elseif kind ==  __kindSlice then
      
      typ.tfun = function(this, array)
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
         -- TODO: come back and fix up Luar.
         -- must set these for Luar (binary Go translation) to work.
         --this[__giPrivateRaw] = array
         --this[__giPrivateSliceProps] = this
         setmetatable(this, __valueSliceMT)
      end;
      typ.init = function(elem)
         typ.elem = elem;
         typ.comparable = false;
         typ.__nil = typ({},{});
      end;
      
   elseif kind ==  __kindArray then
      typ.tfun = function(this, v)
         --print("in tfun ctor function for __kindArray, this="..tostring(this).." and v="..tostring(v))
         this.__val = v;
         this.__array = v; -- like slice, to reuse ipairs method.
         this.__offset = 0; -- like slice.
         this.__constructor = typ
         this.__length = __lenz(v)
         this.__name = "__arrayValue"
         -- TODO: come back and fix up Luar
         -- must set these keys for Luar to work:
         --this[__giPrivateRaw] = v
         --this[__giPrivateArrayProps] = this
         setmetatable(this, __valueArrayMT)
      end;
     --print("in newType for array, and typ.tfun = "..tostring(typ.tfun))
      typ.wrapped = true;
      typ.ptr = __newType(4, __kindPtr, "*" .. str, false, "", false, function(this, array)
                             this.__get = function() return array; end;
                             this.__set = function(v) typ.copy(this, v); end;
                             this.__val = array;
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
         -- pointer. But perhaps this is javascript's prototypal inheritence in action.
         --
         -- gopherjs uses them in comma expressions. example, condensed:
         --     p$1 = new ptrType(...); sa$3.Port = (p$1.nilCheck, p$1[0])
         --
         -- Since comma expressions are not (efficiently) supported in Lua, let
         -- implement the nil check in a different manner.
         -- js: Object.defineProperty(typ.ptr.__nil, "nilCheck", { get= __throwNilPointerError end);
      end;
      -- end __kindArray

      
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
      
      typ.tfun = function(this, v) this.__val = v; end;
      typ.wrapped = true;
      typ.init = function(key, elem)
         typ.key = key;
         typ.elem = elem;
         typ.comparable = false;
      end;
      
   elseif kind ==  __kindStruct then
      
      typ.tfun = function(this, v)
         --print("top of simple kindStruct tfun")
         this.__val = v;
      end;
      typ.wrapped = true;

      -- the typ.prototype will be the
      -- metatable for instances of the struct; this is
      -- equivalent to the prototype in js.
      --
      typ.prototype = {__name="methodSet for "..str, __typ = typ}
      typ.prototype.__index = typ.prototype

      local ctor = function(this, ...)
         --print("top of struct ctor, this="..tostring(this).."; typ.__constructor = "..tostring(typ.__constructor))
         local args = {...}
         --__st(args, "args to ctor")
         --__st(args[1], "args[1]")

         --print("callstack:")
         --print(debug.traceback())
         
         this.__get = function() return this; end;
         this.__set = function(v) typ.copy(this, v); end;
         if typ.__constructor ~= nil then
            -- have to skip the first empty table...
            local skipFirst = {}
            for i,v in ipairs(args) do
               if i > 1 then table.insert(skipFirst, v) end
            end
            typ.__constructor(this, unpack(skipFirst));
         end
         setmetatable(this, typ.ptr.prototype)
      end
      typ.ptr = __newType(4, __kindPtr, "*" .. str, false, pkg, exported, ctor);
      -- __newType sets typ.comparable = true
      __dfsGlobal:addChild(typ.ptr, typ)
      
      -- pointers have their own method sets, but *T can call elem methods in Go.
      typ.ptr.elem = typ;
      typ.ptr.prototype = {__name="methodSet for "..typ.ptr.__str, __typ = typ.ptr}
      typ.ptr.prototype.__index = typ.ptr.prototype

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
         __ipairsZeroCheck(fields)
         for i,f in ipairs(fields) do
            __st(f,"f")
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
            __ipairsZeroCheck(fields)
            for _, f in ipairs(fields) do
               local sw2 = f.__typ.kind
               
               if sw2 == __kindArray or
               sw2 ==  __kindStruct then 
                  f.__typ.copy(dst[f.__prop], src[f.__prop]);
               else
                  dst[f.__prop] = src[f.__prop];
               end
            end
         end;
         --print("jea debug: on __kindStruct: set .copy on typ to .copy=", typ.copy)
         -- /* nil value */
         local properties = {};
         __ipairsZeroCheck(fields)
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
            
            -- assume that e.__typ.fields must be an array!
            -- TODO: remove this assert after confirmation.
            __assertIsArray(e.__typ.fields)
            __ipairsZeroCheck(e.__typ.fields)
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
         print("m is ")
         __st(m,"m")
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


__type__bool    = __newType( 1, __kindBool,    "bool",     true, "", false, nil);
__type__int = __newType( 8, __kindInt,     "int",   true, "", false, nil);
__type__int8    = __newType( 1, __kindInt8,    "int8",     true, "", false, nil);
__type__int16   = __newType( 2, __kindInt16,   "int16",    true, "", false, nil);
__type__int32   = __newType( 4, __kindInt32,   "int32",    true, "", false, nil);
__type__int64   = __newType( 8, __kindInt64,   "int64",    true, "", false, nil);
__type__uint    = __newType( 8, __kindUint,    "uint",     true, "", false, nil);
__type__uint8   = __newType( 1, __kindUint8,   "uint8",    true, "", false, nil);
__type__uint16  = __newType( 2, __kindUint16,  "uint16",   true, "", false, nil);
__type__uint32  = __newType( 4, __kindUint32,  "uint32",   true, "", false, nil);
__type__uint64  = __newType( 8, __kindUint64,  "uint64",   true, "", false, nil);
__type__uintptr = __newType( 8, __kindUintptr,    "uintptr",  true, "", false, nil);
__type__float32 = __newType( 8, __kindFloat32,    "float32",  true, "", false, nil);
__type__float64 = __newType( 8, __kindFloat64,    "float64",  true, "", false, nil);
__type__complex64  = __newType( 8, __kindComplex64,  "complex64",   true, "", false, nil);
__type__complex128 = __newType(16, __kindComplex128, "complex128",  true, "", false, nil);
__type__string  = __newType(16, __kindString,  "string",   true, "", false, nil);
--__type__unsafePointer = __newType( 8, __kindUnsafePointer, "unsafe.Pointer", true, "", false, nil);

__ptrType = function(elem, selfDefnSrcCode)
   if elem == nil then
      error("internal error: cannot call __ptrType() will nil elem")
   end
   local typ = elem.ptr;
   if typ == nil then
      typ = __newType(4, __kindPtr, "*" .. elem.__str, false, "", elem.exported, nil);
      __dfsGlobal:addChild(typ, elem)
      elem.ptr = typ;
      typ.init(elem);
      typ.src = selfDefnSrcCode
   end
   return typ;
end;

__newDataPointer = function(data, constructor)
   if constructor.elem.kind == __kindStruct then
      return data;
   end
   return constructor(function() return data; end, function(v) data = v; end);
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
__arrayType = function(elem, len, selfDefnSrcCode)
   local typeKey = elem.id .. "_" .. len;
   local typ = __arrayTypes[typeKey];
   if typ == nil then
      typ = __newType(24, __kindArray, "[" .. len .. "]" .. elem.__str, false, "", false, nil);
      __arrayTypes[typeKey] = typ;
      __dfsGlobal:addChild(typ, elem)
      typ.init(elem, len);
      typ.src = selfDefnSrcCode
   end
   return typ;
end;


__chanType = function(elem, sendOnly, recvOnly, selfDefnSrcCode)
   
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
      typ.src = selfDefnSrcCode
   end
   return typ;
end;

function __Chan(elem, capacity)
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
__funcType = function(params, results, variadic, selfDefnSrcCode)

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
      typ.src = selfDefnSrcCode
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
__interfaceType = function(methods, selfDefnSrcCode)
   
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
      typ.src = selfDefnSrcCode
   end
   return typ;
end;
__type__emptyInterface = __interfaceType({});
__ifaceNil = {};
__error = __newType(8, __kindInterface, "error", true, "", false, nil);
__error.init({{__prop= "Error", __name= "Error", __pkg= "", __typ= __funcType({}, {__String}, false) }});

__mapTypes = {};
__mapType = function(key, elem, selfDefnSrcCode)
   local typeKey = key.id .. "_" .. elem.id;
   local typ = __mapTypes[typeKey];
   if typ == nil then
      typ = __newType(8, __kindMap, "map[" .. key.__str .. "]" .. elem.__str, false, "", false, nil);
      __mapTypes[typeKey] = typ;

      __dfsGlobal:addChild(typ, key)
      __dfsGlobal:addChild(typ, elem)
      
      typ.init(key, elem);
      typ.src = selfDefnSrcCode
   end
   return typ;
end;

__makeMap = function(keyForFunc, entries, keyType, elemType, mapType)
   local m={};
   for k, e in pairs(entries) do
      local key = keyForFunc(k)
      --print("using key ", key, " for k=", k)
      m[key] = e;
   end
   local mp = _gi_NewMap(keyType, elemType, m);
   --setmetatable(mp, mapType)
   return mp
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

__sliceType = function(elem, selfDefnSrcCode)
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
      typ.src = selfDefnSrcCode
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




function field2strHelper(f)
   local tag = ""
   if f.tag ~= "" then
      tag = string.gsub(f.tag, "\\", "\\\\")
      tag = string.gsub(tag, "\"", "\\\"")
   end
   return f.__name .. " " .. f.__typ.__str .. tag
end

function typeKeyHelper(f)
   return f.__name .. "," .. f.__typ.id .. "," .. f.tag;
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
                               this[f.__prop] = arg
                            else
                               this[f.__prop] = f.__typ.zero();
                            end
                         end
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
      
      for i = 0,#(typ.fields)-1 do
         local f = typ.fields[i];
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
-----------------------
-------- ../utf8.lua -------
-- utf8.lua
--
-- from https://github.com/Stepets/utf8.lua
--
-- Provides UTF-8 aware string functions implemented in pure lua:
-- * utf8len(s)
-- * utf8sub(s, i, j)
-- * utf8reverse(s)
-- * utf8char(unicode)
-- * utf8unicode(s, i, j)
-- * utf8gensub(s, sub_len)
-- * utf8find(str, regex, init, plain)
-- * utf8match(str, regex, init)
-- * utf8gmatch(str, regex, all)
-- * utf8gsub(str, regex, repl, limit)
--
-- If utf8data.lua (containing the lower<->upper case mappings) is loaded, these
-- additional functions are available:
-- * utf8upper(s)
-- * utf8lower(s)
--
-- All functions behave as their non UTF-8 aware counterparts with the exception
-- that UTF-8 characters are used instead of bytes for all units.

--[[
Copyright (c) 2006-2007, Kyle Smith
All rights reserved.

Contributors:
	Alimov Stepan

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

    * Redistributions of source code must retain the above copyright notice,
      this list of conditions and the following disclaimer.
    * Redistributions in binary form must reproduce the above copyright
      notice, this list of conditions and the following disclaimer in the
      documentation and/or other materials provided with the distribution.
    * Neither the name of the author nor the names of its contributors may be
      used to endorse or promote products derived from this software without
      specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
--]]

-- ABNF from RFC 3629
--
-- UTF8-octets = *( UTF8-char )
-- UTF8-char   = UTF8-1 / UTF8-2 / UTF8-3 / UTF8-4
-- UTF8-1      = %x00-7F
-- UTF8-2      = %xC2-DF UTF8-tail
-- UTF8-3      = %xE0 %xA0-BF UTF8-tail / %xE1-EC 2( UTF8-tail ) /
--               %xED %x80-9F UTF8-tail / %xEE-EF 2( UTF8-tail )
-- UTF8-4      = %xF0 %x90-BF 2( UTF8-tail ) / %xF1-F3 3( UTF8-tail ) /
--               %xF4 %x80-8F 2( UTF8-tail )
-- UTF8-tail   = %x80-BF
--

local byte    = string.byte
local char    = string.char
local dump    = string.dump
local find    = string.find
local format  = string.format
local len     = string.len
local lower   = string.lower
local rep     = string.rep
local sub     = string.sub
local upper   = string.upper

-- returns the number of bytes used by the UTF-8 character at byte i in s
-- also doubles as a UTF-8 character validator
local function utf8charbytes (s, i)
	-- argument defaults
	i = i or 1

	-- argument checking
	if type(s) ~= "string" then
		error("bad argument #1 to 'utf8charbytes' (string expected, got ".. type(s).. ")")
	end
	if type(i) ~= "number" then
		error("bad argument #2 to 'utf8charbytes' (number expected, got ".. type(i).. ")")
	end

	local c = byte(s, i)

	-- determine bytes needed for character, based on RFC 3629
	-- validate byte 1
	if c > 0 and c <= 127 then
		-- UTF8-1
		return 1

	elseif c >= 194 and c <= 223 then
		-- UTF8-2
		local c2 = byte(s, i + 1)

		if not c2 then
			error("UTF-8 string terminated early")
		end

		-- validate byte 2
		if c2 < 128 or c2 > 191 then
			error("Invalid UTF-8 character")
		end

		return 2

	elseif c >= 224 and c <= 239 then
		-- UTF8-3
		local c2 = byte(s, i + 1)
		local c3 = byte(s, i + 2)

		if not c2 or not c3 then
			error("UTF-8 string terminated early")
		end

		-- validate byte 2
		if c == 224 and (c2 < 160 or c2 > 191) then
			error("Invalid UTF-8 character")
		elseif c == 237 and (c2 < 128 or c2 > 159) then
			error("Invalid UTF-8 character")
		elseif c2 < 128 or c2 > 191 then
			error("Invalid UTF-8 character")
		end

		-- validate byte 3
		if c3 < 128 or c3 > 191 then
			error("Invalid UTF-8 character")
		end

		return 3

	elseif c >= 240 and c <= 244 then
		-- UTF8-4
		local c2 = byte(s, i + 1)
		local c3 = byte(s, i + 2)
		local c4 = byte(s, i + 3)

		if not c2 or not c3 or not c4 then
			error("UTF-8 string terminated early")
		end

		-- validate byte 2
		if c == 240 and (c2 < 144 or c2 > 191) then
			error("Invalid UTF-8 character")
		elseif c == 244 and (c2 < 128 or c2 > 143) then
			error("Invalid UTF-8 character")
		elseif c2 < 128 or c2 > 191 then
			error("Invalid UTF-8 character")
		end

		-- validate byte 3
		if c3 < 128 or c3 > 191 then
			error("Invalid UTF-8 character")
		end

		-- validate byte 4
		if c4 < 128 or c4 > 191 then
			error("Invalid UTF-8 character")
		end

		return 4

	else
		error("Invalid UTF-8 character")
	end
end

-- returns the number of characters in a UTF-8 string
local function utf8len (s)
	-- argument checking
	if type(s) ~= "string" then
		for k,v in pairs(s) do print('"',tostring(k),'"',tostring(v),'"') end
		error("bad argument #1 to 'utf8len' (string expected, got ".. type(s).. ")")
	end

	local pos = 1
	local bytes = len(s)
	local length = 0

	while pos <= bytes do
		length = length + 1
		pos = pos + utf8charbytes(s, pos)
	end

	return length
end

-- functions identically to string.sub except that i and j are UTF-8 characters
-- instead of bytes
local function utf8sub (s, i, j)
	-- argument defaults
	j = j or -1

	local pos = 1
	local bytes = len(s)
	local length = 0

	-- only set l if i or j is negative
	local l = (i >= 0 and j >= 0) or utf8len(s)
	local startChar = (i >= 0) and i or l + i + 1
	local endChar   = (j >= 0) and j or l + j + 1

	-- can't have start before end!
	if startChar > endChar then
		return ""
	end

	-- byte offsets to pass to string.sub
	local startByte,endByte = 1,bytes

	while pos <= bytes do
		length = length + 1

		if length == startChar then
			startByte = pos
		end

		pos = pos + utf8charbytes(s, pos)

		if length == endChar then
			endByte = pos - 1
			break
		end
	end

	if startChar > length then startByte = bytes+1   end
	if endChar   < 1      then endByte   = 0         end

	return sub(s, startByte, endByte)
end

--[[
-- replace UTF-8 characters based on a mapping table
local function utf8replace (s, mapping)
	-- argument checking
	if type(s) ~= "string" then
		error("bad argument #1 to 'utf8replace' (string expected, got ".. type(s).. ")")
	end
	if type(mapping) ~= "table" then
		error("bad argument #2 to 'utf8replace' (table expected, got ".. type(mapping).. ")")
	end

	local pos = 1
	local bytes = len(s)
	local charbytes
	local newstr = ""

	while pos <= bytes do
		charbytes = utf8charbytes(s, pos)
		local c = sub(s, pos, pos + charbytes - 1)

		newstr = newstr .. (mapping[c] or c)

		pos = pos + charbytes
	end

	return newstr
end


-- identical to string.upper except it knows about unicode simple case conversions
local function utf8upper (s)
	return utf8replace(s, utf8_lc_uc)
end

-- identical to string.lower except it knows about unicode simple case conversions
local function utf8lower (s)
	return utf8replace(s, utf8_uc_lc)
end
]]

-- identical to string.reverse except that it supports UTF-8
local function utf8reverse (s)
	-- argument checking
	if type(s) ~= "string" then
		error("bad argument #1 to 'utf8reverse' (string expected, got ".. type(s).. ")")
	end

	local bytes = len(s)
	local pos = bytes
	local charbytes
	local newstr = ""

	while pos > 0 do
		local c = byte(s, pos)
		while c >= 128 and c <= 191 do
			pos = pos - 1
			c = byte(s, pos)
		end

		charbytes = utf8charbytes(s, pos)

		newstr = newstr .. sub(s, pos, pos + charbytes - 1)

		pos = pos - 1
	end

	return newstr
end

-- http://en.wikipedia.org/wiki/Utf8
-- http://developer.coronalabs.com/code/utf-8-conversion-utility
local function utf8char(unicode)
	if unicode <= 0x7F then return char(unicode) end

	if (unicode <= 0x7FF) then
		local Byte0 = 0xC0 + math.floor(unicode / 0x40);
		local Byte1 = 0x80 + (unicode % 0x40);
		return char(Byte0, Byte1);
	end;

	if (unicode <= 0xFFFF) then
		local Byte0 = 0xE0 +  math.floor(unicode / 0x1000);
		local Byte1 = 0x80 + (math.floor(unicode / 0x40) % 0x40);
		local Byte2 = 0x80 + (unicode % 0x40);
		return char(Byte0, Byte1, Byte2);
	end;

	if (unicode <= 0x10FFFF) then
		local code = unicode
		local Byte3= 0x80 + (code % 0x40);
		code       = math.floor(code / 0x40)
		local Byte2= 0x80 + (code % 0x40);
		code       = math.floor(code / 0x40)
		local Byte1= 0x80 + (code % 0x40);
		code       = math.floor(code / 0x40)
		local Byte0= 0xF0 + code;

		return char(Byte0, Byte1, Byte2, Byte3);
	end;

	error 'Unicode cannot be greater than U+10FFFF!'
end

local shift_6  = 2^6
local shift_12 = 2^12
local shift_18 = 2^18

local utf8unicode
utf8unicode = function(str, i, j, byte_pos)
	i = i or 1
	j = j or i

	if i > j then return end

	local ch,bytes

	if byte_pos then
		bytes = utf8charbytes(str,byte_pos)
		ch  = sub(str,byte_pos,byte_pos-1+bytes)
	else
		ch,byte_pos = utf8sub(str,i,i), 0
		bytes       = #ch
	end

	local unicode

	if bytes == 1 then unicode = byte(ch) end
	if bytes == 2 then
		local byte0,byte1 = byte(ch,1,2)
		local code0,code1 = byte0-0xC0,byte1-0x80
		unicode = code0*shift_6 + code1
	end
	if bytes == 3 then
		local byte0,byte1,byte2 = byte(ch,1,3)
		local code0,code1,code2 = byte0-0xE0,byte1-0x80,byte2-0x80
		unicode = code0*shift_12 + code1*shift_6 + code2
	end
	if bytes == 4 then
		local byte0,byte1,byte2,byte3 = byte(ch,1,4)
		local code0,code1,code2,code3 = byte0-0xF0,byte1-0x80,byte2-0x80,byte3-0x80
		unicode = code0*shift_18 + code1*shift_12 + code2*shift_6 + code3
	end

	return unicode,utf8unicode(str, i+1, j, byte_pos+bytes)
end

-- Returns an iterator which returns the next substring and its byte interval
local function utf8gensub(str, sub_len)
	sub_len        = sub_len or 1
	local byte_pos = 1
	local length   = #str
	return function(skip)
		if skip then byte_pos = byte_pos + skip end
		local char_count = 0
		local start      = byte_pos
		repeat
			if byte_pos > length then return end
			char_count  = char_count + 1
			local bytes = utf8charbytes(str,byte_pos)
			byte_pos    = byte_pos+bytes

		until char_count == sub_len

		local last  = byte_pos-1
		local slice = sub(str,start,last)
		return slice, start, last
	end
end

local function binsearch(sortedTable, item, comp)
	local head, tail = 1, #sortedTable
	local mid = math.floor((head + tail)/2)
	if not comp then
		while (tail - head) > 1 do
			if sortedTable[tonumber(mid)] > item then
				tail = mid
			else
				head = mid
			end
			mid = math.floor((head + tail)/2)
		end
	end
	if sortedTable[tonumber(head)] == item then
		return true, tonumber(head)
	elseif sortedTable[tonumber(tail)] == item then
		return true, tonumber(tail)
	else
		return false
	end
end
local function classMatchGenerator(class, plain)
	local codes = {}
	local ranges = {}
	local ignore = false
	local range = false
	local firstletter = true
	local unmatch = false

	local it = utf8gensub(class)

	local skip
	for c, _, be in it do
		skip = be
		if not ignore and not plain then
			if c == "%" then
				ignore = true
			elseif c == "-" then
				table.insert(codes, utf8unicode(c))
				range = true
			elseif c == "^" then
				if not firstletter then
					error('!!!')
				else
					unmatch = true
				end
			elseif c == ']' then
				break
			else
				if not range then
					table.insert(codes, utf8unicode(c))
				else
					table.remove(codes) -- removing '-'
					table.insert(ranges, {table.remove(codes), utf8unicode(c)})
					range = false
				end
			end
		elseif ignore and not plain then
			if c == 'a' then -- %a: represents all letters. (ONLY ASCII)
				table.insert(ranges, {65, 90}) -- A - Z
				table.insert(ranges, {97, 122}) -- a - z
			elseif c == 'c' then -- %c: represents all control characters.
				table.insert(ranges, {0, 31})
				table.insert(codes, 127)
			elseif c == 'd' then -- %d: represents all digits.
				table.insert(ranges, {48, 57}) -- 0 - 9
			elseif c == 'g' then -- %g: represents all printable characters except space.
				table.insert(ranges, {1, 8})
				table.insert(ranges, {14, 31})
				table.insert(ranges, {33, 132})
				table.insert(ranges, {134, 159})
				table.insert(ranges, {161, 5759})
				table.insert(ranges, {5761, 8191})
				table.insert(ranges, {8203, 8231})
				table.insert(ranges, {8234, 8238})
				table.insert(ranges, {8240, 8286})
				table.insert(ranges, {8288, 12287})
			elseif c == 'l' then -- %l: represents all lowercase letters. (ONLY ASCII)
				table.insert(ranges, {97, 122}) -- a - z
			elseif c == 'p' then -- %p: represents all punctuation characters. (ONLY ASCII)
				table.insert(ranges, {33, 47})
				table.insert(ranges, {58, 64})
				table.insert(ranges, {91, 96})
				table.insert(ranges, {123, 126})
			elseif c == 's' then -- %s: represents all space characters.
				table.insert(ranges, {9, 13})
				table.insert(codes, 32)
				table.insert(codes, 133)
				table.insert(codes, 160)
				table.insert(codes, 5760)
				table.insert(ranges, {8192, 8202})
				table.insert(codes, 8232)
				table.insert(codes, 8233)
				table.insert(codes, 8239)
				table.insert(codes, 8287)
				table.insert(codes, 12288)
			elseif c == 'u' then -- %u: represents all uppercase letters. (ONLY ASCII)
				table.insert(ranges, {65, 90}) -- A - Z
			elseif c == 'w' then -- %w: represents all alphanumeric characters. (ONLY ASCII)
				table.insert(ranges, {48, 57}) -- 0 - 9
				table.insert(ranges, {65, 90}) -- A - Z
				table.insert(ranges, {97, 122}) -- a - z
			elseif c == 'x' then -- %x: represents all hexadecimal digits.
				table.insert(ranges, {48, 57}) -- 0 - 9
				table.insert(ranges, {65, 70}) -- A - F
				table.insert(ranges, {97, 102}) -- a - f
			else
				if not range then
					table.insert(codes, utf8unicode(c))
				else
					table.remove(codes) -- removing '-'
					table.insert(ranges, {table.remove(codes), utf8unicode(c)})
					range = false
				end
			end
			ignore = false
		else
			if not range then
				table.insert(codes, utf8unicode(c))
			else
				table.remove(codes) -- removing '-'
				table.insert(ranges, {table.remove(codes), utf8unicode(c)})
				range = false
			end
			ignore = false
		end

		firstletter = false
	end

	table.sort(codes)

	local function inRanges(charCode)
		for _,r in ipairs(ranges) do
			if r[1] <= charCode and charCode <= r[2] then
				return true
			end
		end
		return false
	end
	if not unmatch then
		return function(charCode)
			return binsearch(codes, charCode) or inRanges(charCode)
		end, skip
	else
		return function(charCode)
			return charCode ~= -1 and not (binsearch(codes, charCode) or inRanges(charCode))
		end, skip
	end
end

--[[
-- utf8sub with extra argument, and extra result value
local function utf8subWithBytes (s, i, j, sb)
	-- argument defaults
	j = j or -1

	local pos = sb or 1
	local bytes = len(s)
	local length = 0

	-- only set l if i or j is negative
	local l = (i >= 0 and j >= 0) or utf8len(s)
	local startChar = (i >= 0) and i or l + i + 1
	local endChar   = (j >= 0) and j or l + j + 1

	-- can't have start before end!
	if startChar > endChar then
		return ""
	end

	-- byte offsets to pass to string.sub
	local startByte,endByte = 1,bytes

	while pos <= bytes do
		length = length + 1

		if length == startChar then
			startByte = pos
		end

		pos = pos + utf8charbytes(s, pos)

		if length == endChar then
			endByte = pos - 1
			break
		end
	end

	if startChar > length then startByte = bytes+1   end
	if endChar   < 1      then endByte   = 0         end

	return sub(s, startByte, endByte), endByte + 1
end
]]

local cache = setmetatable({},{
	__mode = 'kv'
})
local cachePlain = setmetatable({},{
	__mode = 'kv'
})
local function matcherGenerator(regex, plain)
	local matcher = {
		functions = {},
		captures = {}
	}
	if not plain then
		cache[regex] =  matcher
	else
		cachePlain[regex] = matcher
	end
	local function simple(func)
		return function(cC)
			if func(cC) then
				matcher:nextFunc()
				matcher:nextStr()
			else
				matcher:reset()
			end
		end
	end
	local function star(func)
		return function(cC)
			if func(cC) then
				matcher:fullResetOnNextFunc()
				matcher:nextStr()
			else
				matcher:nextFunc()
			end
		end
	end
	local function minus(func)
		return function(cC)
			if func(cC) then
				matcher:fullResetOnNextStr()
			end
			matcher:nextFunc()
		end
	end
	local function question(func)
		return function(cC)
			if func(cC) then
				matcher:fullResetOnNextFunc()
				matcher:nextStr()
			end
			matcher:nextFunc()
		end
	end

	local function capture(id)
		return function(_)
			local l = matcher.captures[id][2] - matcher.captures[id][1]
			local captured = utf8sub(matcher.string, matcher.captures[id][1], matcher.captures[id][2])
			local check = utf8sub(matcher.string, matcher.str, matcher.str + l)
			if captured == check then
				for _ = 0, l do
					matcher:nextStr()
				end
				matcher:nextFunc()
			else
				matcher:reset()
			end
		end
	end
	local function captureStart(id)
		return function(_)
			matcher.captures[id][1] = matcher.str
			matcher:nextFunc()
		end
	end
	local function captureStop(id)
		return function(_)
			matcher.captures[id][2] = matcher.str - 1
			matcher:nextFunc()
		end
	end

	local function balancer(str)
		local sum = 0
		local bc, ec = utf8sub(str, 1, 1), utf8sub(str, 2, 2)
		local skip = len(bc) + len(ec)
		bc, ec = utf8unicode(bc), utf8unicode(ec)
		return function(cC)
			if cC == ec and sum > 0 then
				sum = sum - 1
				if sum == 0 then
					matcher:nextFunc()
				end
				matcher:nextStr()
			elseif cC == bc then
				sum = sum + 1
				matcher:nextStr()
			else
				if sum == 0 or cC == -1 then
					sum = 0
					matcher:reset()
				else
					matcher:nextStr()
				end
			end
		end, skip
	end

	matcher.functions[1] = function(_)
		matcher:fullResetOnNextStr()
		matcher.seqStart = matcher.str
		matcher:nextFunc()
		if (matcher.str > matcher.startStr and matcher.fromStart) or matcher.str >= matcher.stringLen then
			matcher.stop = true
			matcher.seqStart = nil
		end
	end

	local lastFunc
	local ignore = false
	local skip = nil
	local it = (function()
		local gen = utf8gensub(regex)
		return function()
			return gen(skip)
		end
	end)()
	local cs = {}
	for c, bs, be in it do
		skip = nil
		if plain then
			table.insert(matcher.functions, simple(classMatchGenerator(c, plain)))
		else
			if ignore then
				if find('123456789', c, 1, true) then
					if lastFunc then
						table.insert(matcher.functions, simple(lastFunc))
						lastFunc = nil
					end
					table.insert(matcher.functions, capture(tonumber(c)))
				elseif c == 'b' then
					if lastFunc then
						table.insert(matcher.functions, simple(lastFunc))
						lastFunc = nil
					end
					local b
					b, skip = balancer(sub(regex, be + 1, be + 9))
					table.insert(matcher.functions, b)
				else
					lastFunc = classMatchGenerator('%' .. c)
				end
				ignore = false
			else
				if c == '*' then
					if lastFunc then
						table.insert(matcher.functions, star(lastFunc))
						lastFunc = nil
					else
						error('invalid regex after ' .. sub(regex, 1, bs))
					end
				elseif c == '+' then
					if lastFunc then
						table.insert(matcher.functions, simple(lastFunc))
						table.insert(matcher.functions, star(lastFunc))
						lastFunc = nil
					else
						error('invalid regex after ' .. sub(regex, 1, bs))
					end
				elseif c == '-' then
					if lastFunc then
						table.insert(matcher.functions, minus(lastFunc))
						lastFunc = nil
					else
						error('invalid regex after ' .. sub(regex, 1, bs))
					end
				elseif c == '?' then
					if lastFunc then
						table.insert(matcher.functions, question(lastFunc))
						lastFunc = nil
					else
						error('invalid regex after ' .. sub(regex, 1, bs))
					end
				elseif c == '^' then
					if bs == 1 then
						matcher.fromStart = true
					else
						error('invalid regex after ' .. sub(regex, 1, bs))
					end
				elseif c == '$' then
					if be == len(regex) then
						matcher.toEnd = true
					else
						error('invalid regex after ' .. sub(regex, 1, bs))
					end
				elseif c == '[' then
					if lastFunc then
						table.insert(matcher.functions, simple(lastFunc))
					end
					lastFunc, skip = classMatchGenerator(sub(regex, be + 1))
				elseif c == '(' then
					if lastFunc then
						table.insert(matcher.functions, simple(lastFunc))
						lastFunc = nil
					end
					table.insert(matcher.captures, {})
					table.insert(cs, #matcher.captures)
					table.insert(matcher.functions, captureStart(cs[#cs]))
					if sub(regex, be + 1, be + 1) == ')' then matcher.captures[#matcher.captures].empty = true end
				elseif c == ')' then
					if lastFunc then
						table.insert(matcher.functions, simple(lastFunc))
						lastFunc = nil
					end
					local cap = table.remove(cs)
					if not cap then
						error('invalid capture: "(" missing')
					end
					table.insert(matcher.functions, captureStop(cap))
				elseif c == '.' then
					if lastFunc then
						table.insert(matcher.functions, simple(lastFunc))
					end
					lastFunc = function(cC) return cC ~= -1 end
				elseif c == '%' then
					ignore = true
				else
					if lastFunc then
						table.insert(matcher.functions, simple(lastFunc))
					end
					lastFunc = classMatchGenerator(c)
				end
			end
		end
	end
	if #cs > 0 then
		error('invalid capture: ")" missing')
	end
	if lastFunc then
		table.insert(matcher.functions, simple(lastFunc))
	end

	table.insert(matcher.functions, function()
		if matcher.toEnd and matcher.str ~= matcher.stringLen then
			matcher:reset()
		else
			matcher.stop = true
		end
	end)

	matcher.nextFunc = function(self)
		self.func = self.func + 1
	end
	matcher.nextStr = function(self)
		self.str = self.str + 1
	end
	matcher.strReset = function(self)
		local oldReset = self.reset
		local str = self.str
		self.reset = function(s)
			s.str = str
			s.reset = oldReset
		end
	end
	matcher.fullResetOnNextFunc = function(self)
		local oldReset = self.reset
		local func = self.func +1
		local str = self.str
		self.reset = function(s)
			s.func = func
			s.str = str
			s.reset = oldReset
		end
	end
	matcher.fullResetOnNextStr = function(self)
		local oldReset = self.reset
		local str = self.str + 1
		local func = self.func
		self.reset = function(s)
			s.func = func
			s.str = str
			s.reset = oldReset
		end
	end

	matcher.process = function(self, str, start)

		self.func = 1
		start = start or 1
		self.startStr = (start >= 0) and start or utf8len(str) + start + 1
		self.seqStart = self.startStr
		self.str = self.startStr
		self.stringLen = utf8len(str) + 1
		self.string = str
		self.stop = false

		self.reset = function(s)
			s.func = 1
		end

		-- local lastPos = self.str
		-- local lastByte
		local ch
		while not self.stop do
			if self.str < self.stringLen then
				--[[ if lastPos < self.str then
					print('last byte', lastByte)
					ch, lastByte = utf8subWithBytes(str, 1, self.str - lastPos - 1, lastByte)
					ch, lastByte = utf8subWithBytes(str, 1, 1, lastByte)
					lastByte = lastByte - 1
				else
					ch, lastByte = utf8subWithBytes(str, self.str, self.str)
				end
				lastPos = self.str ]]
				ch = utf8sub(str, self.str,self.str)
				--print('char', ch, utf8unicode(ch))
				self.functions[self.func](utf8unicode(ch))
			else
				self.functions[self.func](-1)
			end
		end

		if self.seqStart then
			local captures = {}
			for _,pair in pairs(self.captures) do
				if pair.empty then
					table.insert(captures, pair[1])
				else
					table.insert(captures, utf8sub(str, pair[1], pair[2]))
				end
			end
			return self.seqStart, self.str - 1, unpack(captures)
		end
	end

	return matcher
end

-- string.find
local function utf8find(str, regex, init, plain)
	local matcher = cache[regex] or matcherGenerator(regex, plain)
	return matcher:process(str, init)
end

-- string.match
local function utf8match(str, regex, init)
	init = init or 1
	local found = {utf8find(str, regex, init)}
	if found[1] then
		if found[3] then
			return unpack(found, 3)
		end
		return utf8sub(str, found[1], found[2])
	end
end

-- string.gmatch
local function utf8gmatch(str, regex, all)
	regex = (utf8sub(regex,1,1) ~= '^') and regex or '%' .. regex
	local lastChar = 1
	return function()
		local found = {utf8find(str, regex, lastChar)}
		if found[1] then
			lastChar = found[2] + 1
			if found[all and 1 or 3] then
				return unpack(found, all and 1 or 3)
			end
			return utf8sub(str, found[1], found[2])
		end
	end
end

local function replace(repl, args)
	local ret = ''
	if type(repl) == 'string' then
		local ignore = false
		local num
		for c in utf8gensub(repl) do
			if not ignore then
				if c == '%' then
					ignore = true
				else
					ret = ret .. c
				end
			else
				num = tonumber(c)
				if num then
					ret = ret .. args[num]
				else
					ret = ret .. c
				end
				ignore = false
			end
		end
	elseif type(repl) == 'table' then
		ret = repl[args[1] or args[0]] or ''
	elseif type(repl) == 'function' then
		if #args > 0 then
			ret = repl(unpack(args, 1)) or ''
		else
			ret = repl(args[0]) or ''
		end
	end
	return ret
end
-- string.gsub
local function utf8gsub(str, regex, repl, limit)
	limit = limit or -1
	local ret = ''
	local prevEnd = 1
	local it = utf8gmatch(str, regex, true)
	local found = {it()}
	local n = 0
	while #found > 0 and limit ~= n do
		local args = {[0] = utf8sub(str, found[1], found[2]), unpack(found, 3)}
		ret = ret .. utf8sub(str, prevEnd, found[1] - 1)
		.. replace(repl, args)
		prevEnd = found[2] + 1
		n = n + 1
		found = {it()}
	end
	return ret .. utf8sub(str, prevEnd), n
end

local utf8 = {}
utf8.len = utf8len
utf8.sub = utf8sub
utf8.reverse = utf8reverse
utf8.char = utf8char
utf8.unicode = utf8unicode
utf8.gensub = utf8gensub
utf8.byte = utf8unicode
utf8.find    = utf8find
utf8.match   = utf8match
utf8.gmatch  = utf8gmatch
utf8.gsub    = utf8gsub
utf8.dump    = dump
utf8.format = format
utf8.lower = lower
utf8.upper = upper
utf8.rep     = rep
utf8.charbytes = utf8charbytes
return utf8
-----------------------

`

