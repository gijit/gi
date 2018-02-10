dofile 'fin.lua'

-- compare by value
function ValEq(a,b)   
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
         if not ValEq(vb, va) then
            return false
         end
      end
      return true
   end
   -- string, number, bool, userdata, functions
   return a == b
end

--[[
print(ValEq(0,0))
print(ValEq(0,1))
print(ValEq({},{}))
print(ValEq({a=1},{a=1}))
print(ValEq({a=1},{a=2}))
print(ValEq({a=1},{b=1}))
print(ValEq({a=1,b=2},{a=1,b=2}))
print(ValEq({a=1,b={c=2}},{a=1,b={c=2}}))
print(ValEq({a=1,b={c=2}},{a=1,b={c=3}}))
print(ValEq("hi","hi"))
print(ValEq("he","hi"))
--]]

function expectEq(a, b)
   if not ValEq(a,b) then
      error("expectEq failure: a='"..tostring(a).."' was not equal to b='"..tostring(b).."', of type "..type(b))
   end
end

-- __mapAndJoinStrings(splice, arr, fun)
expectEq("", __mapAndJoinStrings("_", {}, function(x) return x end))
expectEq("a_b_c", __mapAndJoinStrings("_", {"a","b","c"}, function(x) return x end))
expectEq("a", __mapAndJoinStrings("_", {"a"}, function(x) return x end))
expectEq("a", __mapAndJoinStrings("_", {[0]="a"}, function(x) return x end))
expectEq("a_b", __mapAndJoinStrings("_", {[0]="a","b"}, function(x) return x end))
expectEq("1_2_3", __mapAndJoinStrings("_", {[0]=0,1,2}, function(x)  return x+1 end))
expectEq("1_2", __mapAndJoinStrings("_", {[0]=0,1}, function(x)  return x+1 end))
expectEq("1", __mapAndJoinStrings("_", {[0]=0}, function(x)  return x+1 end))
expectEq("", __mapAndJoinStrings("_", {}, function(x)  return x+1 end))

-- __keys = function(m)
expectEq({}, __keys({}))
expectEq({"a","b"}, __keys({b=2, a=4}))
expectEq({"a"}, __keys({a=4}))
