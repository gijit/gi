--------------------------------------------------------------------------------
-- Statistical distributions module.
--
-- Copyright (C) 2011-2016 Stefano Peluchetti. All rights reserved.
--------------------------------------------------------------------------------

return {
  exponential = require("sci.dist._exponential").dist,
  normal      = require("sci.dist._normal").dist,
  lognormal   = require("sci.dist._lognormal").dist,
  gamma       = require("sci.dist._gamma").dist,
  beta        = require("sci.dist._beta").dist,
  student     = require("sci.dist._student").dist,
  uniform     = require("sci.dist._uniform").dist,
  
  mvuniform   = require("sci.dist._uniform").mvdist,
}
