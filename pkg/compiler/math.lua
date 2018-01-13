-- math helper functions

-- x == math.huge   -- test for +inf, inline

-- x == -math.huge  -- test for -inf, inline

-- x ~= x           -- test for nan, inline

-- x > -math.huge and x < math.huge  -- test for finite

-- or their slower counterparts:

math.isnan  = function(x) return x ~= x; end
math.finite = function(x) return x > -math.huge and x < math.huge; end
