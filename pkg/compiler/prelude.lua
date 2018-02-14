-- prelude defines things that should
-- be available before any user code is run.

function _gi_GetRangeCheck(x, i)
  if x == nil or i < 0 or i >= #x then
     error("index out of range")
  end
  return x[i]
end;

function _gi_SetRangeCheck(x, i, val)
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

function __gijit_printQuoted(...)
   local a = {...}
   --print("__gijit_printQuoted called, a = " .. tostring(a), " len=", #a)
   if a[0] ~= nil then
      if type(v) == "string" then
         print("`"..v.."`")
      else
         print(v)
      end
   end
   for i,v in ipairs(a) do
      if type(v) == "string" then
         print("`"..v.."`")
      else
         print(v)
      end
   end
end
