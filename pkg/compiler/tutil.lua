-- tutil.lua: test utilities


-- compare by value
function __ValEq(a,b)   
   local aty = type(a)
   local bty = type(b)
   if aty ~= bty then
      return false
   end
   if aty == "table" then
      -- compare two tables
      for ka,va in pairs(a) do
         vb = b[ka]
         if vb == nil then
            -- b doesn't have key ka in it.
            return false
         end
         if not __ValEq(vb, va) then
            return false
         end
      end
      return true
   end
   -- string, number, bool, userdata, functions
   return a == b
end

--[[
print(__ValEq(0,0))
print(__ValEq(0,1))
print(__ValEq({},{}))
print(__ValEq({a=1},{a=1}))
print(__ValEq({a=1},{a=2}))
print(__ValEq({a=1},{b=1}))
print(__ValEq({a=1,b=2},{a=1,b=2}))
print(__ValEq({a=1,b={c=2}},{a=1,b={c=2}}))
print(__ValEq({a=1,b={c=2}},{a=1,b={c=3}}))
print(__ValEq("hi","hi"))
print(__ValEq("he","hi"))
--]]

function __expectEq(a, b)
   if not __ValEq(a,b) then
      error("__expectEq failure: a='"..tostring(a).."' was not equal to b='"..tostring(b).."', of type "..type(b))
   end
end
