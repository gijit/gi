--------------------------------------------------------------------------------
-- Function minimization module.
--
-- Copyright (C) 2011-2016 Stefano Peluchetti. All rights reserved.
--------------------------------------------------------------------------------

local function tofmin(optim)
  return function(...)
    return optim(1, ...)
  end
end

return {
  de    = tofmin(require("sci.fmin._de").optim),
  lbfgs = tofmin(require("sci.fmin._lbfgs").optim),
}