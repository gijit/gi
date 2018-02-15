--------------------------------------------------------------------------------
-- Sobol quasi random number generator module.
--
-- Credit: this implementation is based on the code published at:
-- http://web.maths.unsw.edu.au/~fkuo/sobol/ .
-- Please notice that the code written in this file is NOT endorsed in any way 
-- by the authors of the original C++ code (on which this implementation is
-- based), S. Joe and F. Y. Kuo, nor they participated in the development of 
-- this Lua implementation.
-- Any bug / problem introduced in this port is my sole responsibility.
--
-- Copyright (C) 2011-2016 Stefano Peluchetti. All rights reserved.
--------------------------------------------------------------------------------

local ffi      = require "ffi"
local dirndata = require "sci.qrng._new-joe-kuo-6-21201"
local xsys     = require "xsys"

local UNROLL = 10

local bit = xsys.bit
local tobit, lshift, rshift = bit.tobit, bit.lshift, bit.rshift
local band, bor, bxor, lsb =  bit.band, bit.bor, bit.bxor, bit.lsb

local sobol_t = ffi.typeof([[
struct {
  int32_t  _n;        // Counter.
  int32_t  _o;        // Offset.
  int32_t  _s;        // Status.
  int32_t  _d;        // Dimension.
  int32_t  _x[21202]; // State.
} ]])

-- For background see:
-- http://web.maths.unsw.edu.au/~fkuo/sobol/joe-kuo-notes.pdf .
local m = dirndata.m -- Sequence of positive integers.
local a = dirndata.a -- Primitive polynomial coefficients.

-- Direction numbers: 32 bits * 21201 dimensions.
-- Maximum number of samples is 2^32-1 (the origin, i=0, is discarded).
local v = ffi.new("int32_t[33][21202]")

-- Fill direction numbers for first dimension, all m = 1.
for i=1,32 do v[i][1] = lshift(1, 32-i) end

local maxdim = 1 -- Current maximum dimension.

local function compute_dn(dim) -- Fill direction numbers up to dimension dim.
  if dim > 21201 then
    error("maximum dimensionality of 21201 exceeded, requested "..dim)
  end
  if dim > maxdim then -- Maxdim is max dimension computed up to now.
    for j=maxdim+1, dim do -- Compute missing dimensions.
      local s = #m[j]
      for i=1,s do 
        v[i][j] = lshift(m[j][i], 32-i)
      end
      for i=s+1,32 do
        v[i][j] = bxor(v[i-s][j], rshift(v[i-s][j], s))
        for k=1,s-1 do 
          v[i][j] = bxor(v[i][j], band(rshift(a[j], s-1-k), 1) * v[i-k][j])
        end
      end
    end
    maxdim = dim
  end
end

local function fill_state(self, c)
  for i=1,21201 do self._x[i] = c end
end

local next_state_template = xsys.template([[
return function(self)
  local n = self._d
  local c = lsb(self._n) + 1
  |for n=1,UNROLL do
  ${n==1 and 'if' or 'elseif'} n == ${n} then
    |for i=1,n do  
    self._x[${i}] = bxor(self._x[${i}], v[c][${i}])
    |end
  |end
  ${UNROLL>=1 and 'else'}
    for i=1,n do
      self._x[i] = bxor(self._x[i], v[c][i])
    end
  ${UNROLL>=1 and 'end'}
  self._o = 0
end
]])

local next_state = xsys.exec(next_state_template({ UNROLL = UNROLL }), 
  "next_state", { lsb = lsb, bxor = bxor, v = v })

local sobol_mt = {
  __new = function(ct)
    -- -2^31 is state at first iteration, which is precomputed.
    local o = ffi.new(ct, -1, 0, 0, 21201)
    fill_state(o, -2^31)
    return o
  end,
  -- Move rng to next state (exactly all dimensions must have been used).
  nextstate = function(self)
    self._n = tobit(self._n + 1)
    if self._n <= 0 then
      if self._s == 0 then -- Zero iterations up to now.
        if self._o ~= 0 then
          error(":sample() called before :nextstate()")
        end
        self._n = -1 -- Get back to self._n = 0 condition next iteration.
        self._s = 1
        return
      elseif self._s == 1 then -- One iteration up to now, initializing.
        self._s = 2
        self._d = self._o
        fill_state(self, 0)        
        compute_dn(self._d)
        -- Recover missed computations.
        self._n = 1
        next_state(self)
        self._o = self._d
        self._n = 2
      else
        error("limit of 2^32-1 states exceeded")
      end
    end
    -- Usual operation.
    if self._o ~= self._d then
      error("not enough samples generated for current state, dimensionality is "
          ..self._d)
    end
    next_state(self)
  end,
  -- Result between (0, 1) extremes excluded.
  sample = function(self)
    self._o = self._o + 1
    if self._o > self._d then
      error("too many samples generated for current state, dimensionality is "
           ..self._d)
    end
    return (bxor(self._x[self._o], 0x80000000) + 0x80000000)*(1/2^32)
  end,
}
sobol_mt.__index = sobol_mt

local qrng = ffi.metatype(sobol_t, sobol_mt)

return {
  qrng = qrng,
}