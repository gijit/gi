-- prelude defines things that should
-- be available before any user code is run.

function __gi_GetRangeCheck(x, i)
   if x == nil or i < 0 or i >= #x then
      error("index out of range")
  end
  return x[i]
end;

function __gi_SetRangeCheck(x, i, val)
  --print("SetRangeCheck. x=", x, " i=", i, " val=", val)
  if x == nil or i < 0 or i >= #x then
     error("index out of range")
  end
  x[i] = val
  return val
end;

