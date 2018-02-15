--------------------------------------------------------------------------------
-- Special math functions module.
--
-- Copyright (C) 2011-2016 Stefano Peluchetti. All rights reserved.
--------------------------------------------------------------------------------

local ffi  = require "ffi"
local bit  = require "bit"
local xsys = require "xsys"

-- Cache all builtin math functions.
local
abs, acos, asin, atan, atan2, ceil, cos, cosh, deg, exp, floor, fmod, frexp,
huge, ldexp, log, log10, max, min, modf, pi, pow, rad, random, randomseed, sin, 
sinh, sqrt, tan, tanh 
= xsys.from(math, [[
abs, acos, asin, atan, atan2, ceil, cos, cosh, deg, exp, floor, fmod, frexp,
huge, ldexp, log, log10, max, min, modf, pi, pow, rad, random, randomseed, sin, 
sinh, sqrt, tan, tanh
]])

local type = type

-- Step-wise functions ---------------------------------------------------------
-- Halfway cases are rounded away from zero, regardless of the current rounding
-- direction.
local function round(x)
  return x < 0.0 and ceil(x - 0.5) or floor(x + 0.5)
end

local function step(x) -- 1 if x >= 0, 0 otherwise.
  return max(0, min(floor(x) + 1, 1))
end

local function sign(x) -- 1 if x >= 0, -1 otherwise.
  return -1 + step(x)*2
end

-- Special math functions ------------------------------------------------------
-- From Marsaglia, "Evaluating the Normal Distribution"
-- http://www.jstatsoft.org/v11/a05/paper
-- around 15 digits of absolute precision.
local function phi(x)
  -- Ensure 0 <= phi(x) <= 1 :
  if x <= -8 then
    return 0
  elseif x >= 8 then
    return 1
  else
  
  local s, b, q = x, x, x^2
  for i=3,1/0,2 do
    b = b*q/i
    local t = s
    s = t + b
    if s == t then break end
  end
  return 0.5 + s*exp(-0.5*q - 0.91893853320467274178)
  end
end

-- Inverse cdf for sampling based on Peter John Acklam research, see:
-- http://home.online.no/~pjacklam/notes/invnorm/ .
-- Maximum relative error of 1.15E-9 and machine accuracy with refinement.
-- In iphifast domain must be (0, 1) extremes excluded.
local iphifast, iphi
do 
  local a = ffi.new("double[7]", { 0,
  -3.969683028665376e+01,
  2.209460984245205e+02,
  -2.759285104469687e+02,
  1.383577518672690e+02,
  -3.066479806614716e+01,
  2.506628277459239e+00 })
  local b = ffi.new("double[6]", { 0,
  -5.447609879822406e+01,
  1.615858368580409e+02,
  -1.556989798598866e+02,
  6.680131188771972e+01,
  -1.328068155288572e+01 })
  local c = ffi.new("double[7]", { 0,
  -7.784894002430293e-03,
  -3.223964580411365e-01,
  -2.400758277161838e+00,
  -2.549732539343734e+00,
  4.374664141464968e+00,
  2.938163982698783e+00 })
  local d = ffi.new("double[5]", { 0,
  7.784695709041462e-03,
  3.224671290700398e-01,
  2.445134137142996e+00,
  3.754408661907416e+00 })

  -- PERF: just two branches, central with high prob.
  iphifast = function(p)
    -- Rational approximation for central region:
    if abs(p - 0.5) < 0.47575 then -- 95.14% of cases if p ~ U(0, 1).
      local q = p - 0.5
      local r = q^2
      return (((((a[1]*r+a[2])*r+a[3])*r+a[4])*r+a[5])*r+a[6])*q /
             (((((b[1]*r+b[2])*r+b[3])*r+b[4])*r+b[5])*r+1)
    -- Rational approximation for the two ends:
    else
      local iu = ceil(p - 0.97575)      -- 1 if p > 0.97575 (upper).
      local z = (1 - iu)*p + iu*(1 - p) -- p if lower, (1 - p) if upper.
      local sign = 1 - 2*iu             -- 1 if lower, -1 if upper.
      local q = sqrt(-2*log(z))
      return sign*(((((c[1]*q+c[2])*q+c[3])*q+c[4])*q+c[5])*q+c[6]) /
                   ((((d[1]*q+d[2])*q+d[3])*q+d[4])*q+1)
    end
  end
  
  iphi = function(p) 
    if p <= 0 then
      return -1/0
    elseif p >=1 then
      return 1/0
    else
      local x = iphifast(p)
      local e = phi(x) - p
      local u = e*sqrt(2*pi)*exp(x^2/2)
      return x - u/(1 + x*u/2)
    end
   end
  
end

local gamma, loggamma
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
  -- pag 116 for optimal formula and coefficients. Theoretical accuracy of 
  -- 16 digits is likely in practice to be around 14.
  -- Domain: R except 0 and negative integers.
  gamma = function(z)  
    -- Reflection formula to handle negative z plane.
    -- Better to branch at z < 0 as some use cases focus on z >= 0 only.
    if z < 0 then 
      return pi/(sin(pi*z)*gamma(1 - z)) 
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
  loggamma = function(z)
    if z < 0 then 
      return log(pi) - log(abs(sin(pi*z))) - loggamma(1 - z) 
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
    return log(gamma_c) + (z - 0.5)*log(z  + gamma_r10 - 0.5) 
      - (z - 0.5) + log(sum)
  end

end

-- Domain: a > 0 and b > 0.
local function logbeta(a, b)
  if a <= 0 or b <= 0 then return 0/0 end
  return loggamma(a) + loggamma(b) - loggamma(a + b)
end

-- Domain: a > 0 and b > 0.
local function beta(a, b)
  return exp(logbeta(a, b))
end

-- Support for generic arguments -----------------------------------------------
local function recmax(x, y, ...)
  return y and recmax(
    type(x) ~= "number" and x:max(y) or (type(y) ~= "number" and y.max(x, y) or 
      max(x, y)), ...) or x
end

local function recmin(x, y, ...)
  return y and recmin(
    type(x) ~= "number" and x:min(y) or (type(y) ~= "number" and y.min(x, y) or 
      min(x, y)), ...) or x
end

local disp2 = { 
  atan2 = function(x, y)
    return type(x) ~= "number" and x:atan2(y) 
       or (type(y) ~= "number" and y.atan2(x, y) or atan2(x, y))
  end,
  fmod = function(x, y)
    return type(x) ~= "number" and x:fmod(y) 
       or (type(y) ~= "number" and y.fmod(x, y) or fmod(x, y))
  end,
  ldexp = function(x, y)
    return type(x) ~= "number" and x:ldexp(y) 
       or (type(y) ~= "number" and y.ldexp(x, y) or ldexp(x, y))
  end,
  pow = function(x, y)
    return type(x) ~= "number" and x:pow(y) 
       or (type(y) ~= "number" and y.pow(x, y) or pow(x, y))
  end,
  
  beta = function(x, y)
    return type(x) ~= "number" and x:beta(y) 
       or (type(y) ~= "number" and y.beta(x, y) or beta(x, y))
  end,
  logbeta = function(x, y)
    return type(x) ~= "number" and x:logbeta(y) 
       or (type(y) ~= "number" and y.logbeta(x, y) or logbeta(x, y))
  end,
}

local generic = {
  -- Constants:
  pi   = pi,
  huge = huge,
  
  -- Random numbers:
  random     = random,
  randomseed = randomseed,
  
  -- Generic dispatch based on one variable:
  abs   = function(x) return type(x) == "number" and abs(x)   or x:abs() end,
  acos  = function(x) return type(x) == "number" and acos(x)  or x:acos() end,
  asin  = function(x) return type(x) == "number" and asin(x)  or x:asin() end,
  atan  = function(x) return type(x) == "number" and atan(x)  or x:atan() end,
  ceil  = function(x) return type(x) == "number" and ceil(x)  or x:ceil()  end,
  cos   = function(x) return type(x) == "number" and cos(x)   or x:cos() end,
  cosh  = function(x) return type(x) == "number" and cosh(x)  or x:cosh() end,
  deg   = function(x) return type(x) == "number" and deg(x)   or x:deg() end,
  exp   = function(x) return type(x) == "number" and exp(x)   or x:exp() end,
  floor = function(x) return type(x) == "number" and floor(x) or x:floor() end,
  frexp = function(x) return type(x) == "number" and frexp(x) or x:frexp() end,
  log   = function(x) return type(x) == "number" and log(x)   or x:log() end,
  log10 = function(x) return type(x) == "number" and log10(x) or x:log10() end,
  modf  = function(x) return type(x) == "number" and modf(x)  or x:modf() end,
  rad   = function(x) return type(x) == "number" and rad(x)   or x:rad() end,
  sin   = function(x) return type(x) == "number" and sin(x)   or x:sin() end,
  sinh  = function(x) return type(x) == "number" and sinh(x)  or x:sinh() end,
  sqrt  = function(x) return type(x) == "number" and sqrt(x)  or x:sqrt() end,
  tan   = function(x) return type(x) == "number" and tan(x)   or x:tan() end,
  tanh  = function(x) return type(x) == "number" and tanh(x)  or x:tanh() end,  
  
  -- Generic dispatch based on two variables:
  atan2 = disp2.atan2,
  fmod  = disp2.fmod,
  ldexp = disp2.ldexp,
  pow   = disp2.pow,
  
  -- Special dispatch for vararg max and min functions:
  max = recmax,
  min = recmin,  
    
  -- General dispatch based on one variable:
  gamma    = function(x) return type(x) == "number" and gamma(x)    or 
    x:gamma() end,
  iphi     = function(x) return type(x) == "number" and iphi(x)     or 
    x:iphi() end,
  loggamma = function(x) return type(x) == "number" and loggamma(x) or 
    x:loggamma() end,
  phi      = function(x) return type(x) == "number" and phi(x)      or 
    x:phi() end,
  round    = function(x) return type(x) == "number" and round(x)    or 
    x:round() end,
  sign     = function(x) return type(x) == "number" and sign(x)     or 
    x:sign() end,
  step     = function(x) return type(x) == "number" and step(x)     or 
    x:step() end,
  
  -- Generic dispatch based on two variables:
  beta    = disp2.beta,
  logbeta = disp2.logbeta,
  
  -- Note: no generic dispatch for private functions:
  _iphifast = iphifast,
}

--------------------------------------------------------------------------------

local M = xsys.table.union(math, {
  generic     = generic,  
  round       = round,
  step        = step,
  sign        = sign,  
  phi         = phi,
  iphi        = iphi,
  gamma       = gamma,
  loggamma    = loggamma,
  logbeta     = logbeta,
  beta        = beta,  
  _iphifast   = iphifast,
})

M.generic.std = M -- Gives access to builtin functions.

return M
