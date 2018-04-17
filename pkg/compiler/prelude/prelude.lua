-- prelude defines things that should
-- be available before any user code is run.

function __gi_GetRangeCheck(x, i)
   if i == nil then
      print(debug.traceback())
      error "where is i nil??"
   end
   if x == nil then
      print(debug.traceback())
      error "where is x nil??"
   end
   if x == nil or i < 0 or i >= #x then
      error("index out of range: i="..tostring(i).." vs #x is "..tostring(#x))
   end
   --print("range check on x = "..tostring(x).." at i = "..tostring(i).." with #x="..tostring(#x).." looks okay, returning value: ", x[i])
   --__st(x, "x")
  return x[i]
end;

function __gi_SetRangeCheck(x, i, val)
  --print("SetRangeCheck. x=".. __st(x) .." i="..tostring(i).." val=", val)
  if x == nil or i < 0 or i >= #x then
     error("index out of range; i="..tostring(i)..". #x is "..tostring(#x))
  end
  x[i] = val
  return val
end;

