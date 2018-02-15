--------------------------------------------------------------------------------
-- Complex numbers.
--
-- Copyright (C) 2011-2016 Stefano Peluchetti. All rights reserved.
--------------------------------------------------------------------------------

local ffi = require 'ffi'

local sqrt = math.sqrt

local complex = ffi.typeof("complex")

local complex_mt = {
  __add = function(x, y)
    x, y = complex(x), complex(y)
    return complex(x.re + y.re, x.im + y.im)
  end,
  __sub = function(x, y)
    x, y = complex(x), complex(y)
    return complex(x.re - y.re, x.im - y.im)
  end,
  __mul = function(x, y)
    x, y = complex(x), complex(y)
    return complex(x.re*y.re - x.im*y.im, x.re*y.im + x.im*y.re)
  end,
  __div = function(x, y)
    x, y = complex(x), complex(y)
    local d = y.re^2 + y.im^2
    return complex((x.re*y.re + x.im*y.im)/d, (x.im*y.re - x.re*y.im)/d)
  end,
}

ffi.metatype(complex, complex_mt)

local function cabs(x)
  return sqrt(x.re^2 + x.im^2)
end

return {
  new = complex,

  abs = cabs,
}