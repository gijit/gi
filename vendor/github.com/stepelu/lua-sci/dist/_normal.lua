--------------------------------------------------------------------------------
-- Normal statistical distribution.
--
-- Copyright (C) 2011-2016 Stefano Peluchetti. All rights reserved.
--------------------------------------------------------------------------------

-- Inverse cdf for sampling based on Peter John Acklam research, see:
-- http://home.online.no/~pjacklam/notes/invnorm/ .
-- Maximum relative error of 1.15E-9, fine for generation of random variates.
--
-- The following paper has some considerations on this topic:
-- http://epub.wu.ac.at/664/1/document.pdf .
--
-- Moreover, the following procedure has been employed as empirical validation 
-- of the sampling procedure: for a given prng the samples returned from a 
-- normal(0, 1) are converted back to uniform(0, 1) via the GSL cdf for the
-- normal(0, 1) which offers machine accuracy implementation. These numbers are
-- used as input to the small, normal, and big crush batteries from TestU01. 
-- All test passed aside from the followings suspect p-values. Also reported the
-- p-values obtained by repeating the suspect tests:
--
-- lfib4:
-- BIG: swalk_RandomWalkl test:
-- N = 1, n = 100000000, r = 0, s = 5, L0 = 50, L1 = 50
-- J: p-value of test : 6.8e-04 *****
-- Repeating the test yields p-values: 0.21, 0.51, 0.49, 0.07, 0.11 [OK]
--
-- mrg32k3a:
-- NORM: smarsa_MatrixRank test:
-- N = 1, n = 1000000, r = 0, s = 30, L = 60, k = 60
-- p-value of test : 9.4e-04 *****
-- Repeating the test yields p-values: 0.70, 0.24, 0.12, 0.40, 0.06 [OK]
--
-- The prng mrg32k3a has 2^53 accuracy and hence is of particular relevance
-- for the tails.

local ffi   = require "ffi"
local xsys  = require "xsys"
local math  = require "sci.math"

local exp, log, sqrt, pi, sin, cos, abs, ceil, _iphifast = xsys.from(math,
     "exp, log, sqrt, pi, sin, cos, abs, ceil, _iphifast")

local norm_mt = {
  __new = function(ct, mu, sigma)
    if not mu or not sigma then
      error("distribution parameters must be set at construction")
    end
    if sigma <= 0 then
      error("sigma must be positive, sigma is ", sigma)
    end
    return ffi.new(ct, mu, sigma)
  end,
  copy = function(self)
    return ffi.new(ffi.typeof(self), self)
  end,
  range = function(self)
    return -1/0, 1/0
  end,
  pdf = function(self, x)
    local mu, sigma = self._mu, self._sigma
    return exp(-0.5*((x - mu)/sigma)^2) / (sqrt(2*pi)*sigma)
  end,
  logpdf = function(self, x)
    local mu, sigma = self._mu, self._sigma
    return -0.5*((x - mu)/sigma)^2 - 0.5*log(2*pi) - log(sigma)
  end,
  mean = function(self)
    return self._mu
  end,
  var = function(self)
    return self._sigma^2
  end,
  absmoment = function(self, mm)
    if self._mu == 0 and self._sigma == 1 and mm == 1 then
      return sqrt(pi/2)
    else
      error("NYI: only first absolute moment currently implemented")
    end
  end,
  sample = function(self, rng)
    return _iphifast(rng:sample())*self._sigma + self._mu
  end,
}
norm_mt.__index = norm_mt

local dist = ffi.metatype("struct {double _mu, _sigma;}", norm_mt)

return {
  dist = dist,
}
