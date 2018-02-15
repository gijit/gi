--------------------------------------------------------------------------------
-- Pierre L'Ecuyer MRG pseudo rngs module.
--
-- Copyright (C) 2011-2016 Stefano Peluchetti. All rights reserved.
--------------------------------------------------------------------------------

-- This specific implementation has been tested against small, normal and big
-- crush batteries of TestU01 passing all tests.

-- TODO: Replace tables with matrices of cdata<uint64_t> (and their operators).

local ffi  = require "ffi"

-- Constants defining the rng.
local a12   = 1403580
local a13   = -810728
local m1    = 2^32 - 209
local a21   = 527612
local a23   = -1370589
local m2    = 2^32 - 22853
local scale = 1/(m1 + 1)
local y0    = 12345
local MAX_PERIOD_LOG2 = 191

local ull = ffi.typeof("uint64_t")

local function ulmat()
  return {{0ULL, 0ULL, 0ULL}, {0ULL, 0ULL, 0ULL}, {0ULL, 0ULL, 0ULL}}
end

-- Return modular matrix product: X1*X2 % m. All matrices (3, 3).
-- Require X1, X2 of uint64_t if m = 2^32 as products can go up almost 2^64.
local function modmul(X1, X2, m) 
  local Y = ulmat()
  for r=1,3 do
    for c=1,3 do
      local v = 0ULL
      for i=1,3 do
        local prod = (X1[r][i]*X2[i][c]) % m -- prod is uint64_t.
        v = (v + prod) % m -- v is uint64_t.
      end
      Y[r][c] = v
    end
  end
  return Y
end


local function vvmodmul(A, i, x1, x2, x3, m)
  return tonumber(((A[i][1]*x1) % m + (A[i][2]*x2) % m + (A[i][3]*x3) % m) % m)
end

-- Skip ahead matrices valid for A^p with p = 2^i, i >= 1 (so p even).
local aheadA1, aheadA2 = { }, { }
do 
  local A1, A2 = ulmat(), ulmat() -- Initialized to 0ULL.
  A1[1][2] = ull(a12 % m1); A1[1][3] = ull(a13 % m1)
  A1[2][1] = 1ULL
  A1[3][2] = 1ULL 
  A2[1][1] = ull(a21 % m2); A2[1][3] = ull(a23 % m2)
  A2[2][1] = 1ULL
  A2[3][2] = 1ULL
  for i=1,MAX_PERIOD_LOG2 do
    A1 = modmul(A1, A1, m1)
    aheadA1[i] = A1
    A2 = modmul(A2, A2, m2)
    aheadA2[i] = A2
  end
end

local function sarg(...)
  return "{"..table.concat({ ... }, ",").."}"
end

local mrg_mt = {
  __new = function(ct, self)
    return ffi.new(ct, y0, y0, y0, y0, y0, y0)
  end,
  __tostring = function(self)
    local o = self
    return "mrg32k3a "..sarg(o._y11, o._y12, o._y13, o._y21, o._y22, o._y23)
  end,
  copy = function(self)
    return ffi.new(ffi.typeof(self), self)
  end,
  -- Sampling algorithm (combine excluded), see pag. 11 of:
  -- http://citeseerx.ist.psu.edu/viewdoc/summary?doi=10.1.1.48.1341 .
  -- This rng, with this parameters set, allows to keep the state info in 
  -- double precision numbers (i.e. Lua numbers).
  sample = function(self)
    -- assert(math.abs(a12*self._y12 + a13*self._y13) < 2^53)
    local p1 = (a12*self._y12 + a13*self._y13) % m1
    self._y13 = self._y12; self._y12 = self._y11; self._y11 = p1  
    -- assert(math.abs(a21*self._y21 + a23*self._y23) < 2^53)
    local p2 = (a21*self._y21 + a23*self._y23) % m2
    self._y23 = self._y22; self._y22 = self._y21; self._y21 = p2  
    -- This branchless version is faster, shift by 1 instead of branching:
    return ((p1 - p2) % m1 + 1)*scale
  end,
  -- Skip ahead 2^k samples and returns last sample.
  -- Notice n = 2^k, k >= 1.
  _sampleahead2pow = function(self, k)
    local A1 = aheadA1[k]
    local A2 = aheadA2[k]
    local y11, y12, y13 = self._y11, self._y12, self._y13
    local y21, y22, y23 = self._y21, self._y22, self._y23
    self._y11 = vvmodmul(A1, 1, y11, y12, y13, m1)
    self._y12 = vvmodmul(A1, 2, y11, y12, y13, m1)
    self._y13 = vvmodmul(A1, 3, y11, y12, y13, m1)
    self._y21 = vvmodmul(A2, 1, y21, y22, y23, m2)
    self._y22 = vvmodmul(A2, 2, y21, y22, y23, m2)
    self._y23 = vvmodmul(A2, 3, y21, y22, y23, m2)
    -- This branchless version is faster, shift by 1 instead of branching:
    return ((self._y11 - self._y21) % m1 + 1)*scale
  end
}
mrg_mt.__index = mrg_mt

local mrg32k3a = ffi.metatype(
  "struct { double _y11, _y12, _y13, _y21, _y22, _y23; }", mrg_mt)

return {
  mrg32k3a = mrg32k3a,  
}
