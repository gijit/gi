--------------------------------------------------------------------------------
-- Pseudo random number generators module.
--
-- Copyright (C) 2011-2016 Stefano Peluchetti. All rights reserved.
--------------------------------------------------------------------------------

local ffi       = require "ffi"
local marsaglia = require "sci.prng._marsaglia"
local mrg       = require "sci.prng._mrg"

local M = {
  std      = marsaglia.lfib4,
  lfib4    = marsaglia.lfib4,
  kiss99   = marsaglia.kiss99,
  mrg32k3a = mrg.mrg32k3a,
}

local function restore_unsafe(str)
  assert(type(str) == "string")
  local sep = str:find(" ")
  local rng = str:sub(1, sep-1)
  local arg = str:sub(sep+1)
  return ffi.new(assert(M[rng]), unpack(assert(loadstring("return "..arg))()))
end

local function restore(str)
  local ok, rng = pcall(restore_unsafe, str)
  if not ok then
    error("string is not a valid serialization of a prng")
  end
  return rng
end

M.restore = restore

local function new_parallel(rng, totalperiod)
  return function(n)
    local log2n = math.ceil(math.log(n)/math.log(2)) -- Log2.
    local rngperiod = totalperiod - log2n
    local r = rng()
    local out = { }
    for i=1,n do
      out[i] = r:copy()
      r:_sampleahead2pow(rngperiod)
    end
    return out
  end
end

M.parallel = {
  mrg32k3a = new_parallel(M.mrg32k3a, 191),
}

M.parallel.std = M.parallel.mrg32k3a

return M