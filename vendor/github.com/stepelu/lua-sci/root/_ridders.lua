--------------------------------------------------------------------------------
-- Ridders root finding module.
--
-- Copyright (C) 2011-2016 Stefano Peluchetti. All rights reserved.
--------------------------------------------------------------------------------

local math = require "sci.math"

local max, abs, sqrt, sign = math.max, math.abs, math.sqrt, math.sign

local function root(f, xl, xu, stop)
  if not(xl < xu) then
    error("xl < xu required: xl="..xl..", xu="..xu)
  end
  local yl, yu = f(xl), f(xu)
  if not (yl*yu <= 0) then
    error("root not bracketed by f(xl)="..yl..", f(xu)="..yu)
  end
  while true do
    local xm = xl + 0.5*(xu - xl) -- Avoid xm > xl or xm < xu.
    local ym = f(xm)
    local x1, y1
    do
      local d = ym^2 - yl*yu
      if d == 0 then -- Function is flat on xl, xm, xu.
        x1, y1 = xm, ym
      else -- Exponential inversion.
        x1 = xm + (xm - xl)*sign(yl - yu)*ym/sqrt(d)
        y1 = f(x1)
      end
    end
    if stop(x1, y1, xl, xu, yl, yu) then
      return x1, y1, xl, xu, yl, yu
    end
    if y1*ym <= 0 then
      if xm < x1 then
        xl, xu = xm, x1
        yl, yu = ym, y1
      else
        xl, xu = x1, xm
        yl, yu = y1, ym
      end
    elseif ym*yl <= 0 then
      xu, yu = xm, ym
    else
      xl, yl = xm, ym
    end
  end
end

return {
  root = root,
}