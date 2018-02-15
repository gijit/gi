--------------------------------------------------------------------------------
-- Uniform statistical distribution.
--
-- Copyright (C) 2011-2016 Stefano Peluchetti. All rights reserved.
--------------------------------------------------------------------------------

local ffi   = require "ffi"

local log = math.log

local uni_mt = {
  __new = function(ct, a, b)
    if not a or not b then
      error("distribution parameters must be set at construction")
    end
    if not (a < b) then
      error("a < b is required, a is "..a..", b is "..b)
    end
    return ffi.new(ct, a, b)
  end,
  copy = function(self)
    return ffi.new(ffi.typeof(self), self)
  end,
  range = function(self)
    return self._a, self._b
  end,
  pdf = function(self, x)
    return 1/(self._b - self._a)
  end,
  logpdf = function(self, x)
    return -log(self._b - self._a)
  end,
  mean = function(self)
    return 0.5*(self._a + self._b)
  end,
  var = function(self)
    return (self._b - self._a)^2/12
  end,
  sample = function(self, rng)
    return self._a + (self._b - self._a)*rng:sample()
  end,
}
uni_mt.__index = uni_mt

local dist = ffi.metatype("struct { double _a, _b; }", uni_mt)

-- Multi variate uniform distribution:
local mvuni_mt = {
  sample = function(self, rng, x)
    local a, b = self._a, self._b
    for i=1,#x do 
      x[i] = a[i] + (b[i] - a[i])*rng:sample()
    end
  end,
}
mvuni_mt.__index = mvuni_mt

local function mvdist(a, b)
  if not a or not b then
    error("distribution parameters must be set at construction")
  end
  if #a ~= #b then
    error("a and b must have the same size: #a="..#a..", #b="..#b)
  end
  for i=1,#a do
    if not (a[i] < b[i]) then 
      error("a < b is required: a[i]="..a[i]..", b[i]="..b[i].." for i="..i)
    end
  end
  return setmetatable({ _a = a:copy(), _b = b:copy() }, mvuni_mt)
end

return {
  dist   = dist, 
  mvdist = mvdist,
}
