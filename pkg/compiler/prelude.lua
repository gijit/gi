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


-- length of array, counting [0] if present.
function __lenz(array)      
   local n = #array
   if array[0] ~= nil then
      n=n+1
   end
   return n
end


function __gijit_printQuoted(...)
   local a = {...}
   --print("__gijit_printQuoted called, a = " .. tostring(a), " len=", #a)
   if a[0] ~= nil then
      local v = a[0]
      if type(v) == "string" then
         print("`"..v.."`")
      else
         print(v)
      end
   end
   for _,v in ipairs(a) do
      if type(v) == "string" then
         print("`"..v.."`")
      else
         print(v)
      end
   end
end
