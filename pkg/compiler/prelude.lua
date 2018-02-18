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

-- complex numbers
-- "complex[?]" ?

ffi = require('ffi')
point = ffi.metatype("struct { double re, im; }", {
    __add = function(a, b)
     return point(a.re + b.re, a.im + b.im)
 end
})

-- 1+2i
-- if reloaded, we get this error, so comment out for now.
-- prelude.lua:34: cannot change a protected metatable
-- See also https://stackoverflow.com/questions/325323/is-there-anyway-to-avoid-this-security-issue-in-lua
--[[
point = ffi.metatype("complex", {
    __add = function(a, b)
     return point(a.re + b.re, a.im + b.im)
 end
})
--]]

function _gi_NewComplex128(real, imag)

end

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
   --print("__gijit_printQuoted called, a = " .. tostring(a), " len=", __lenz(a))
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
