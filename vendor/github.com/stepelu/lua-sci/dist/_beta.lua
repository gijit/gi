--------------------------------------------------------------------------------
-- Beta statistical distribution.
--
-- Copyright (C) 2011-2016 Stefano Peluchetti. All rights reserved.
--------------------------------------------------------------------------------

local xsys   = require "xsys"
local ffi    = require "ffi"
local math   = require "sci.math"
local _gamma = require "sci.dist._gamma"

local exp, log, sqrt, min, beta, logbeta = xsys.from(math, 
     "exp, log, sqrt, min, beta, logbeta")

local gamma_sampleab = _gamma.sampleab
     
local beta_mt = {
  __new = function(ct, alpha, beta)
    if not alpha or not beta then
      error("distribution parameters must be set at construction")
    end
    if alpha <= 0 or beta <= 0 then
      error("alpha and beta must be positive, alpha is "..alpha..", beta is "
        ..beta)
    end
    return ffi.new(ct, alpha, beta)
  end,
  copy = function(self)
    return ffi.new(ffi.typeof(self), self)
  end,
  range = function(self)
    return 0, 1
  end,
  pdf = function(self, x)
    if x < 0 or x > 1 then return 0 end
    local a, b = self._a, self._b
    return x^(a - 1) * (1 - x)^(b - 1) / beta(a, b)
  end,
  logpdf = function(self, x)
    if x < 0 or x > 1 then return -1/0 end
    local a, b = self._a, self._b
    return (a - 1)*log(x) + (b - 1)*log(1 - x) - logbeta(a, b)
  end,
  mean = function(self)
    return self._a/(self._a + self._b)
  end,
  var = function(self)
    local a, b = self._a, self._b
    return (a*b)/((a + b)^2*(a + b + 1))
  end,
  sample = function(self, rng)
    local x = gamma_sampleab(self._a, 1, rng)
    local y = gamma_sampleab(self._b, 1, rng)
    return x/(x + y)
  end,
}
beta_mt.__index = beta_mt

local dist = ffi.metatype("struct { double _a, _b; }", beta_mt)

return {
  dist = dist,
}