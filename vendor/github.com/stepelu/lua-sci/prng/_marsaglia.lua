--------------------------------------------------------------------------------
-- George Marsaglia pseudo rngs module.
--
-- Credit: George Marsaglia Newsgroups posted code:
-- http://www.math.niu.edu/~rusin/known-math/99/RNG .
--
-- Copyright (C) 2011-2016 Stefano Peluchetti. All rights reserved.
--------------------------------------------------------------------------------

-- These specific implementations have been tested against small, normal and big
-- crush batteries of TestU01, the following suspect p-value has been observed:
--
-- kiss99:
-- BIG crush: smarsa_MatrixRank test:
-- N = 10, n = 1000000, r = 0, s = 5, L = 30, k = 30
-- Test on the sum of all N observations: p-value of test : 7.6e-04 *****
-- Repeating the test yields p-values: 0.58, 0.21, 0.70, 0.08, 0.11 [OK]

local ffi  = require "ffi"
local bit  = require "bit"
local xsys = require "xsys"

local tobit, band, bxor, lshift, rshift =  xsys.from(bit,
     "tobit, band, bxor, lshift, rshift")
      
local function sarg(...)
  return "{"..table.concat({ ... }, ",").."}"
end

-- Guarantee range (0, 1) extremes excluded.
local function sample_double(self)
  local b = self:_bitsample()
  return (bxor(b, 0x80000000) + (0x80000000+1)) * (1/(2^32+1))
end

local kiss99_mt = {
  __new = function(ct)
    return ffi.new(ct, tobit(12345), tobit(34221), tobit(12345), tobit(65435))
  end,  
  __tostring = function(self)
    return "kiss99 "..sarg(self._s1 , self._s2 ,self._s3 , self._s4)
  end,
  copy = function(self)
    return ffi.new(ffi.typeof(self), self)
  end,
  _bitsample = function(self)
    local r = self
    r._s1 = tobit(tobit(69069*r._s1) + 1234567)   
    local b = bxor(r._s2, lshift(r._s2, 17))
    b = bxor(b, rshift(b, 13))
    r._s2 = bxor(b, lshift(b, 5))   
    r._s3 = tobit(tobit(36969*band(r._s3, 0xffff)) + rshift(r._s3, 16))
    r._s4 = tobit(tobit(18000*band(r._s4, 0xffff)) + rshift(r._s4, 16))
    b = tobit(lshift(r._s3, 16) + r._s4)
    return tobit(r._s2 + bxor(r._s1, b))
  end,
  sample = sample_double,
}
kiss99_mt.__index = kiss99_mt

local kiss99 = ffi.metatype("struct { int32_t _s1, _s2, _s3, _s4; }", kiss99_mt)

local lfib4_mt = {
  __new = function(ct)
    -- Follow Marsaglia initialization.
    local o = ffi.new(ct) -- Zero filled => _i is 0.
    local r = kiss99()
    for i=0,255 do o._s[i] = r:_bitsample() end
    return o
  end,  
  __tostring = function(self)
    local t = { }
    for i=1,256 do t[i] = self._s[i-1] end
    t = "{"..table.concat(t, ",").."}"
    return "lfib4 "..sarg(t, self._i)
  end,
  copy = function(self)
    return ffi.new(ffi.typeof(self), self)
  end,
  _bitsample = function(self)
    local r = self
    r._i = band(r._i + 1, 255)
    r._s[r._i] = tobit(tobit(r._s[r._i] + r._s[band(r._i+58, 255)]) 
      + tobit(r._s[band(r._i+119, 255)] + r._s[band(r._i+178, 255)]))
    return r._s[r._i]
  end,
  sample = sample_double,
}
lfib4_mt.__index = lfib4_mt

local lfib4 = ffi.metatype("struct { int32_t _s[256]; int32_t _i; } ", lfib4_mt)

return {
  kiss99    = kiss99,
  lfib4     = lfib4,
}