--------------------------------------------------------------------------------
-- Student-t statistical distribution.
--
-- Copyright (C) 2011-2016 Stefano Peluchetti. All rights reserved.
--------------------------------------------------------------------------------

local xsys  = require "xsys"
local ffi   = require "ffi"
local math  = require "sci.math"

local sqrt, log, pi, cos, gamma, loggamma, beta, logbeta = xsys.from(math,
     "sqrt, log, pi, cos, gamma, loggamma, beta, logbeta")

local stud_mt = {
  __new = function(ct, nu)
    if not nu then
      error("distribution parameters must be set at construction")
    end
    if nu <= 0 then
      error("nu must be positive, nu is "..nu)
    end
    return ffi.new(ct, nu)
  end,
  copy = function(self)
    return ffi.new(ffi.typeof(self), self)
  end,
  range = function(self)
    return -1/0, 1/0
  end,
  pdf = function(self, x)
    local nu = self._nu
    return (1 + x^2/nu)^(-0.5*(nu + 1)) / (sqrt(nu)*beta(0.5, 0.5*nu))
  end,
  logpdf = function(self, x)
    local nu = self._nu
    return -(0.5*(nu + 1))*log(1 + x^2/nu) - 0.5*log(nu) - logbeta(0.5, 0.5*nu)
  end,
  mean = function(self)
    if self._nu <= 1 then 
      return 0/0
    else 
      return 0
    end
  end,
  var = function(self)
    local nu = self._nu
    if nu <= 1 then 
      return 0/0
    elseif nu <= 2 then 
      return 1/0
    else 
      return nu/(nu - 2) 
    end
  end,
  absmoment = function(self, mm)
    local nu = self._nu
    local num = nu^(0.5*mm)*gamma(0.5*(mm+1))*gamma(0.5*(nu)) 
    local den = sqrt(pi)*gamma(0.5*nu)
    return num/den
  end,
  sample = function(self, rng)
    local nu, u1, u2 = self._nu, rng:sample(), rng:sample()
    return sqrt(nu*(u1^(-2/nu) - 1))*cos(2*pi*u2)
  end,
}
stud_mt.__index = stud_mt
     
local dist = ffi.metatype("struct { double _nu; }", stud_mt)

return {
  dist = dist
}
