-- math helper functions

-- x == math.huge   -- test for +inf, inline

-- x == -math.huge  -- test for -inf, inline

-- x ~= x           -- test for nan, inline

-- x > -math.huge and x < math.huge  -- test for finite

-- or their slower counterparts:

math.isnan  = function(x) return x ~= x; end
math.finite = function(x) return x > -math.huge and x < math.huge; end

math.nan = math.huge * 0

__truncateToInt = function(x)
   if x >= 0 then
       return x - (x % 1)
   end
   return x + (-x % 1)
end

__integerByZeroCheck = function(x)
   if not math.finite(x) then
      error("integer divide by zero")
   end
   -- eliminate any fractional part
   if x >= 0 then
       return x - (x % 1)
   end
   return x + (-x % 1)
end

function __max(a,b)
   if a > b then
      return a
   end
   return b
end

function __min(a,b)
   if a < b then
      return a
   end
   return b
end
