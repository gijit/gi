--------------------------------------------------------------------------------
-- Newton-type methods for root finding module.
--
-- Copyright (C) 2011-2016 Stefano Peluchetti. All rights reserved.
--------------------------------------------------------------------------------

-- TODO: Consider math.abs(y0) < epsilon safeguard instead of y0 == 0.

-- Make sure that the ordering is right.
local function rebracket(x1, y1, xl, xu, yl, yu)
  if yl*y1 <= 0 then
    return xl, x1, yl, y1
  else
    return x1, xu, y1, yu
  end
end

-- Required: f(x) --> y, f1.
local function newton(f, xl, xu, stop)
  if not(xl < xu) then
    error("xl < xu required: xl="..xl..", xu="..xu)
  end
  local yl, yu = f(xl), f(xu)
  if not (yl*yu <= 0) then
    error("root not bracketed by f(xl)="..yl..", f(xu)="..yu)
  end
  local x0 = xl + 0.5*(xu - xl) -- Avoid xm > xl or xm < xu.
  local y0, f10 = f(x0)
  while true do
    local x1 = x0 - y0/f10
    if x1 == x0 then
      if y0 == 0 then
        return x0, y0, x0, x0, y0, y0
      else
        return nil, "x1 == x0, f(x0)="..y0..", f1(x0)="..f10
      end
    end
    if not (xl <= x1 and x1 <= xu) then 
      return nil, "x1 outside bracket, f(x0)="..y0..", f1(x0)="..f10
    end
    local y1, f11 = f(x1)
    if stop(x1, y1, xl, xu, yl, yu) then
      return x1, y1, xl, xu, yl, yu
    end
    xl, xu, yl, yu = rebracket(x1, y1, xl, xu, yl, yu)
    x0, y0, f10 = x1, y1, f11 
  end
end

-- Required: f(x) --> y, f1, f2.
local function halley(f, xl, xu, stop)
  if not(xl < xu) then
    error("xl < xu required: xl="..xl..", xu="..xu)
  end
  local yl, yu = f(xl), f(xu)
  if not (yl*yu <= 0) then
    error("root not bracketed by f(xl)="..yl..", f(xu)="..yu)
  end
  local x0 = xl + 0.5*(xu - xl) -- Avoid xm > xl or xm < xu.
  local y0, f10, f20 = f(x0)  
  while true do
    local x1 = x0 - 2*y0*f10/(2*f10^2 - y0*f20)
    if x1 == x0 then
      if y0 == 0 then
        return x0, y0, x0, x0, y0, y0
      else
        return nil, "x1 == x0, f(x0)="..y0..", f1(x0)="..f10..", f2(x0)="..f20
      end
    end
    if not (xl <= x1 and x1 <= xu) then 
      return nil, "x1 outside bracket, f(x0)="..y0..", f1(x0)="..f10
                ..", f2(x0)="..f20
    end
    local y1, f11, f21 = f(x1)
    if stop(x1, y1, xl, xu, yl, yu) then
      return x1, y1, xl, xu, yl, yu
    end
    xl, xu, yl, yu = rebracket(x1, y1, xl, xu, yl, yu)
    x0, y0, f10, f20 = x1, y1, f11, f21
  end
end

return { 
  newton = newton,
  halley = halley,
}