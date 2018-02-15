--------------------------------------------------------------------------------
-- Function maximization module.
--
-- Copyright (C) 2011-2016 Stefano Peluchetti. All rights reserved.
--------------------------------------------------------------------------------

local function tofmax(optim)
  return function(...)
    return optim(-1, ...)
  end
end

return {
  de    = tofmax(require("sci.fmin._de").optim),
  lbfgs = tofmax(require("sci.fmin._lbfgs").optim),
}