--------------------------------------------------------------------------------
-- Exponential statistical distribution.
--
-- Copyright (C) 2011-2016 Stefano Peluchetti. All rights reserved.
--------------------------------------------------------------------------------

local ffi   = require "ffi"

local exp, log = math.exp, math.log

local expo_mt = {
  __new = function(ct, lambda)
    if not lambda then
      error("distribution parameters must be set at construction")
    end
    if lambda <= 0 then
      error("lambda must be positive, lambda is ", lambda)
    end
    return ffi.new(ct, lambda)
  end,
  copy = function(self)
    return ffi.new(ffi.typeof(self), self)
  end,
  range = function(self)
    return 0, 1/0
  end,
  pdf = function(self, x)
    if x < 0 then return 0 end
    return self._lambda*exp(-self._lambda*x)
  end,
  logpdf = function(self, x)
    if x < 0 then return -1/0 end
    return log(self._lambda) -self._lambda*x
  end,
  mean = function(self)
    return 1/self._lambda
  end,
  var = function(self)
    return 1/self._lambda^2
  end,  
  sample = function(self, rng)
    return -log(rng:sample())/self._lambda
  end,
}
expo_mt.__index = expo_mt

local dist = ffi.metatype("struct { double _lambda; }", expo_mt)

return {
  dist = dist,
}
