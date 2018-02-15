--------------------------------------------------------------------------------
-- Automatic differentiation module.
--
-- Copyright (C) 2011-2016 Stefano Peluchetti. All rights reserved.
--------------------------------------------------------------------------------

-- TODO: Specialized gamma, loggamma, beta, logbeta.

-- TODO: Introduce specialized matrix multiply, inverse, ecc ecc using BLAS
-- TODO: (when available) based on the results summarized in Mike Giles paper.

-- PERF: We considered one-shot version with matrix-tape which computes the 
-- PERF: gradient in one function call but gains where not substantial and 
-- PERF: memory management is more complicated. Probably better to just 
-- PERF: introduce reverse-mode differentiation directly.

local ffi  = require "ffi"
local xsys = require "xsys"
local math = require "sci.math"

local type = type

local
abs, acos, asin, atan, atan2, ceil, cos, cosh, deg, exp, floor, fmod, frexp,
huge, ldexp, log, log10, max, min, modf, pi, pow, rad, random, randomseed, sin, 
sinh, sqrt, tan, tanh,
round, step, sign,
phi, iphi, gamma, loggamma, logbeta, beta
= xsys.from(math, [[
abs, acos, asin, atan, atan2, ceil, cos, cosh, deg, exp, floor, fmod, frexp,
huge, ldexp, log, log10, max, min, modf, pi, pow, rad, random, randomseed, sin, 
sinh, sqrt, tan, tanh,
round, step, sign,
phi, iphi, gamma, loggamma, logbeta, beta
]])

-- Forward mode, single directional derivative ---------------------------------
local dn -- Dual number: value + adjoint value.

-- Modified from sci.math, works *only* with dual number type.
local dngamma, dnloggamma
do
  -- r(10).
  local gamma_r10 = 10.900511
  -- dk[0], ..., dk[10].
  local gamma_dk = ffi.new("double[11]", 
    2.48574089138753565546e-5,
    1.05142378581721974210,
    -3.45687097222016235469,
    4.51227709466894823700,
    -2.98285225323576655721,
    1.05639711577126713077,
    -1.95428773191645869583e-1,
    1.70970543404441224307e-2,
    -5.71926117404305781283e-4,
    4.63399473359905636708e-6,
    -2.71994908488607703910e-9
  )
  local gamma_c = 2*sqrt(exp(1)/pi)

  -- Lanczos approximation, see:
  -- Pugh[2004]: AN ANALYSIS OF THE LANCZOS GAMMA APPROXIMATION
  -- http://bh0.physics.ubc.ca/People/matt/Doc/ThesesOthers/Phd/pugh.pdf
  -- page 116 for optimal formula and coefficients. Theoretical accuracy of 
  -- 16 digits is likely in practice to be around 14.
  -- Domain: R except 0 and negative integers.
  dngamma = function(z)  
    -- Reflection formula to handle negative z plane.
    -- Better to branch at z < 0 as some probabilistic use cases only consider 
    -- the case z >= 0.
    if z < 0 then 
      return pi/((pi*z):sin()*dngamma(1 - z)) 
    end  
    local sum = gamma_dk[0]
    sum = sum + gamma_dk[1]/(z + 0)
    sum = sum + gamma_dk[2]/(z + 1) 
    sum = sum + gamma_dk[3]/(z + 2) 
    sum = sum + gamma_dk[4]/(z + 3) 
    sum = sum + gamma_dk[5]/(z + 4) 
    sum = sum + gamma_dk[6]/(z + 5) 
    sum = sum + gamma_dk[7]/(z + 6) 
    sum = sum + gamma_dk[8]/(z + 7) 
    sum = sum + gamma_dk[9]/(z + 8) 
    sum = sum + gamma_dk[10]/(z + 9)  
    return gamma_c*((z  + gamma_r10 - 0.5)/exp(1))^(z - 0.5)*sum
  end

  -- Returns log(abs(gamma(z))).
  -- Domain: R except 0 and negative integers.
  dnloggamma = function(z)
    if z < 0 then 
      return log(pi) - (pi*z):sin():abs():log() - dnloggamma(1 - z) 
    end  
    local sum = gamma_dk[0]
    sum = sum + gamma_dk[1]/(z + 0)
    sum = sum + gamma_dk[2]/(z + 1) 
    sum = sum + gamma_dk[3]/(z + 2) 
    sum = sum + gamma_dk[4]/(z + 3) 
    sum = sum + gamma_dk[5]/(z + 4) 
    sum = sum + gamma_dk[6]/(z + 5) 
    sum = sum + gamma_dk[7]/(z + 6) 
    sum = sum + gamma_dk[8]/(z + 7) 
    sum = sum + gamma_dk[9]/(z + 8) 
    sum = sum + gamma_dk[10]/(z + 9) 
    -- For z >= 0 gamma function is positive, no abs() required.
    return log(gamma_c) + (z - 0.5)*(z  + gamma_r10 - 0.5):log()
      - (z - 0.5) + sum:log()
  end

end

-- Domain: a > 0 and b > 0.
local function dnlogbeta(a, b)
  if a <= 0 or b <= 0 then return 0/0 end
  local lga = type(a) == "number" and loggamma(a) or dnloggamma(a)
  local lgb = type(b) == "number" and loggamma(b) or dnloggamma(b)
  return lga + lgb - dnloggamma(a + b)
end

-- Domain: a > 0 and b > 0.
local function dnbeta(a, b)
  return dnlogbeta(a, b):exp()
end

-- Derivative of phi function:
local function dphi(x)
  return (1/sqrt(2*pi))*exp(-0.5*x^2)
end

-- Use branchless optimization whenever finite values:
local function dnmax(x, y)
  x, y = dn(x), dn(y)
  if max(abs(x._v), abs(y._v)) == 1/0 then
    return x >=y and x or y
  else -- Branchless optimization.
    local z = step(y._v - x._v) -- 1 if y >= x, 0 otherwise.
    return dn(z*y._v + (1 - z)*x._v, z*y._a + (1 - z)*x._a)
  end
end
-- Use branchless optimization whenever finite values:
local function dnmin(x, y)
  x, y = dn(x), dn(y)
  if max(abs(x._v), abs(y._v)) == 1/0 then
    return x <= y and x or y
  else -- Branchless optimization.
    local z = step(x._v - y._v) -- 1 if x >= y, 0 otherwise.
    return dn(z*y._v + (1 - z)*x._v, z*y._a + (1 - z)*x._a)
  end
end

-- Note: dual numbers are immutable: new ones generated by operators.
local dn_mt = {
  __unm = function(x)
    return dn(-x._v, -x._a)
  end,
  __add = function(x, y) x, y = dn(x), dn(y)
    return dn(x._v + y._v, x._a + y._a)
  end,
  __sub = function(x, y) x, y = dn(x), dn(y)
    return dn(x._v - y._v, x._a - y._a)
  end,
  __mul = function(x, y) x, y = dn(x), dn(y)
    return dn(x._v*y._v, x._a*y._v + y._a*x._v)
  end,
  __div = function(x, y) x, y = dn(x), dn(y)
    return dn(x._v/y._v, (x._a*y._v - y._a*x._v)/y._v^2)
  end,
  __pow = function(x, y) -- Optimized version.
    if type(y) == "number" then
      return dn(x._v^y, y*x._v^(y-1)*x._a)
    elseif type(x) == "number" then
      return dn(x^y._v, x^y._v*log(x)*y._a)
    else
      return dn(x._v^y._v, x._v^y._v*(log(x._v)*y._a + y._v/x._v*x._a))
    end
  end,
  __eq = function(x, y) x, y = dn(x), dn(y)
    return x._v == y._v
  end,
  __lt = function(x, y) x, y = dn(x), dn(y)
    return x._v < y._v
  end,
  __le = function(x, y) x, y = dn(x), dn(y)
    return x._v <= y._v
  end,
  __tostring = function(x)
    return tostring(x._v) -- Better to mimic behavior of numbers.
  end,
  __tonumber = function(x) -- Honored only by xsys.string.width.
    return tonumber(x._v)
  end,
  
  copy = function(x)
    return dn(x)
  end,
  
  val = function(x)
    return x._v
  end,
  adj = function(x)
    return x._a
  end,
  
  sin  = function(x) return dn(sin(x._v),  x._a*cos(x._v)) end,
  cos  = function(x) return dn(cos(x._v),  x._a*(-sin(x._v))) end,
  tan  = function(x) return dn(tan(x._v),  x._a*(1 + tan(x._v)^2)) end,
  asin = function(x) return dn(asin(x._v), x._a/sqrt(1 - x._v^2)) end,
  acos = function(x) return dn(acos(x._v), -x._a/sqrt(1 - x._v^2)) end,
  atan = function(x) return dn(atan(x._v), x._a/(1 + x._v^2)) end,
  sinh = function(x) return dn(sinh(x._v), x._a*cosh(x._v)) end,
  cosh = function(x) return dn(cosh(x._v), x._a*sinh(x._v)) end,
  tanh = function(x) return dn(tanh(x._v), x._a*(1 - tanh(x._v)^2)) end,
  exp  = function(x) return dn(exp(x._v),  x._a*exp(x._v)) end,
  log  = function(x) return dn(log(x._v),  x._a/x._v) end,
  sqrt = function(x) return dn(sqrt(x._v), x._a/(2*sqrt(x._v))) end,
  abs  = function(x) return dn(abs(x._v),  x._a*sign(x._v)) end,
  -- Stick to dn type to improve type stability:
  floor = function(x) return dn(floor(x._v), 0) end,
  ceil  = function(x) return dn(ceil(x._v),  0) end,
  -- Stick to dn type to improve type stability:
  round = function(x) return dn(round(x._v), 0) end,
  step  = function(x) return dn(step(x._v),  0) end,
  sign  = function(x) return dn(sign(x._v),  0) end,
  
  gamma    = function(x) return    dngamma(x) end,
  loggamma = function(x) return dnloggamma(x) end,
  
  beta     = function(x, y) return    dnbeta(x, y) end,
  logbeta  = function(x, y) return dnlogbeta(x, y) end,

  phi = function(x) return dn(phi(x._v), x._a*dphi(x._v)) end,
  iphi = function(x) local y = iphi(x._v); return dn(y, x._a/dphi(y)) end,

  max   = dnmax,
  min   = dnmin,
}
dn_mt.__index = dn_mt

-- Note: _v comes first. This allows to construct a dual number from a number
-- using dn(number) => adjoint part correctly initialized to 0.
dn = ffi.metatype("struct { double _v, _a; }", dn_mt)

-- To improve type stability we always pass all arguments wrt differentiation
-- will take place as dual numbers, only one of which will have adj part 
-- equal to 1.
local pderf_template = xsys.template[[
local f, dn = f, dn
| local args = { }
| for i=1,n do
|   args[i] = "x"..i
| end
| args = concat(args, ",")
return function(${args})
| local dargs = args..","
| for i=1,#dxi do
|   local dx = dxi[i]
|   dargs = dargs:gsub("x"..dx..",", "dn(x"..dx..",0),")
| end
| dargs = dargs:gsub(",$", "")
| local retadj = { } 
| for i=1,#dxi do
| local from, to = "dn%(x"..dxi[i]..",0%)", "dn%(x"..dxi[i]..",1%)"
local y${i} = dn(f(${dargs:gsub(from, to)}))
| retadj[i] = "y"..i..":adj()"
| end
return y1:val(),${concat(retadj, ",")}
end
]]

local function derivativef(f, n, ...)
  assert(1 <= n, "function's argument # must be positive")
  local dxi
  if select("#", ...) == 0 then
    dxi = { }
    for i=1,n do 
      dxi[i] = i
    end
  else
    dxi = { ... }
    for i=1,#dxi do
      local dx = dxi[i]
      assert(dx <= n, "differentiating variable outside function's argument #")
    end
  end
  local src = pderf_template({ n = n, dxi = dxi, concat = table.concat })
  return xsys.exec(src, "<derivative>", { f = f, dn = dn })
end

-- For forward mode differentiation we could just use grad(f, x, y) to return
-- f(x) and set y to the gradient of f(x), but gradients are best computed in 
-- reverse mode differentiation, and a stack will need to be allocated for f, 
-- the length of which (it's growth-able) likely depends on the f itself hence
-- it's best stored via a closure.
local function gradientf(f, n)
  local algdn = require("sci.alg").typeof(dn, dn)
  local xd = algdn.vec(n)
  return function(x, grad)
    local n = #x
    if n ~= #grad then
      error("value length must be equal to gradient length")
    end
    if n ~= #xd then
      error("value and gradient lengths must be as per initialization")
    end
    local val = 0
    for i=1,n do
      for j=1,n do xd[j] = x[j] end
      xd[i] = xd[i] + dn(0, 1)
      local vd = dn(f(xd))
      val = vd:val()
      grad[i] = vd:adj()
    end
    return val
  end
end

return {
  dn          = dn,
  derivativef = derivativef,
  gradientf   = gradientf,
}
