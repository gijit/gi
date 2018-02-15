--------------------------------------------------------------------------------
-- Double exponential method for fast numerical integration of analytic real
-- functions, see:
-- http://crd-legacy.lbl.gov/~dhbailey/dhbpapers/dhb-tanh-sinh.pdf
-- and its references.
--
-- Copyright (C) 2011-2016 Stefano Peluchetti. All rights reserved.
--------------------------------------------------------------------------------

-- TODO: Unrolling to improve performance.
-- TODO: Use FFI arrays for abscissas and weights to improve performance.
-- TODO: Move change of variables code in separate module for re-use.

local data  = require "sci.quad._dblexp_precomputed"
local gmath = require "sci.math".generic

local abscissas = data.abscissas
local weigths   = data.weigths

local pi, cos, tan = math.pi, math.cos, math.tan
local abs = gmath.abs

-- Change of variables to bounded interval -------------------------------------
-- Finite case: x in (a, b) -> y in (u, v):
-- y    = u + (x - a)*(v - u)/(b - a)
-- x    = a + (y - u)*(b - a)/(v - u)
-- dxdy = (b - a)/(v - u)
local function chgedfinite(f, y, a, b, u, v)
  local dxdy = (b - a)/(v - u)
  return dxdy*f(a + (y - u)*dxdy)
end

-- Half open case: x in (a, +oo) -> y in (u, v):
-- y    = u + (v - u)*(x - a)/(1 + x - a)
-- x    = a + ((y - u)/(v - u))/(1 - (y - u)/(v - u))
-- dxdy = (v - u)/(1 - y)^2
local function chgedopenpos(f, y, a, b, u, v)
  local dxdy = (v - u)/(1 - y)^2
  local z = (y - u)/(v - u)
  return dxdy*f(a + z/(1 - z))
end

-- Half open case: x in (-oo, b) -> y in (u, v):
-- Use y = -x to transform to (-b, +oo) and apply chgopenpos.
local function chgedopenneg(f, y, a, b, u, v)
  local a = -b
  local dxdy = (v - u)/(1 - y)^2
  local z = (y - u)/(v - u)
  return dxdy*f(-(a + z/(1 - z)))
end

-- Open case: x in (-oo, +oo) -> y in (u, v):
-- y    = u + (atan(x)/pi + 0.5)*(v - u)
-- x    = tan(pi*(y - u)/(v - u) - 0.5)
-- dxdy = (sec(y)^2)*pi/(v - u)
local function chgedopen(f, y, a, b, u, v)
  local z = pi*(y - u)/(v - u) - 0.5
  local dxdy = (1/cos(z)^2)*pi/(v - u)
  return dxdy*f(tan(z))   
end

local function chgedintegrand(f, y, a, b, u, v)
  -- assert(a < b and u < v)
  local afinite, bfinite = a > -1/0, b < 1/0
  if afinite and bfinite then         -- Case (a, b).
    return chgedfinite (f, y, a, b, u, v)
  elseif afinite and not bfinite then -- Case (a, +oo).
    return chgedopenpos(f, y, a, b, u, v)
  elseif not afinite and bfinite then -- Case (-oo, b).
    return chgedopenneg(f, y, a, b, u, v)
  else                                -- Case (-oo, +oo).
    return chgedopen   (f, y, a, b, u, v)
  end
end

-- Integration -----------------------------------------------------------------
-- Change the integration interval to (-1, 1).
local function normf(f, y, a, b)
  return chgedintegrand(f, y, a, b, -1, 1)
end

-- Integration on interval (-1, 1).
local dblexp = function(f, a, b, abserror)
  if a > b then
    error("a <= b is required, a is "..tostring(a)..", b is "..tostring(b))
  end
  if a == b then
    return 0
  end
  abserror = abserror or 1e-6
  local absdelta = 0
  -- Level 1.
  local h = 1/2
  local x, w = abscissas[1], weigths[1]
  -- Threat special case of 0, no reflection.
  local sum = w[7]*normf(f, x[7], a, b)
  for i=1,6 do 
    sum = sum + w[i]*(normf(f, x[i], a, b) + normf(f, -x[i], a, b))
  end
  local integral = sum*h
  -- Levels 2,...,7.
  for level=2,7 do
    x, w = abscissas[level], weigths[level]
    sum = 0
    for i=1,#w do
      sum = sum + w[i]*(normf(f, x[i], a, b) + normf(f, -x[i], a, b))
    end
    h = h/2
    local newintegral = sum*h + integral/2
    absdelta = abs(newintegral - integral)
    integral = newintegral    
    if absdelta < abserror then break end
  end  
  return integral, absdelta
end

return {
  quad = dblexp,
}
