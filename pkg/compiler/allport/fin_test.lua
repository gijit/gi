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
expectEq({0}, __keys({[0]=4}))
expectEq({1}, __keys({[1]=3}))
expectEq({7,11}, __keys({[7]="seven", [11]="eleven"}))
local f = function() end
local g = function() end
expectEq({tostring(f)}, __keys({[f]="seven"}))
-- no way to know if a function (as key) is
-- greater or smalller than another function.
-- We use tostring to get a string with the table address
-- before comparing.
--
local sf = tostring(f)
local sg = tostring(g)
x = {sg,sf}
if sf < sg then
   x = {sf,sg}
end
expectEq(x, __keys({[f]="seven", [g]="eight"}))

-- basic types, zero values
expectEq("0LL", tostring(__Int()))
expectEq("0ULL", tostring(__Uint()))
expectEq("0", tostring(__Float64()))
expectEq('""', tostring(__String()))

-- basic types, non-zero values
expectEq("-43LL", tostring(__Int(-43LL)))
expectEq("42ULL", tostring(__Uint(42ULL)))
expectEq("-0.3", tostring(__Float64(-0.3)))
expectEq('"hello world"', tostring(__String("hello world")))

                                                                                   
                                         
expectEq(__basicValue2kind("hi"), __kindString)
expectEq(__basicValue2kind(""), __kindString)
expectEq(__basicValue2kind(true), __kindBool)
expectEq(__basicValue2kind(false), __kindBool)

expectEq(__basicValue2kind(1LL), __kindInt)
expectEq(__basicValue2kind(-1LL), __kindInt)
expectEq(__basicValue2kind(int8(-3)), __kindInt8)
expectEq(__basicValue2kind(int8(3)), __kindInt8)
expectEq(__basicValue2kind(int16(0)), __kindInt16)
expectEq(__basicValue2kind(int16(-1)), __kindInt16)
expectEq(__basicValue2kind(int32(1)), __kindInt32)
expectEq(__basicValue2kind(int32(-1)), __kindInt32)

-- can't distinguish __kindInt from __kindInt64
-- they are both ctype<int64_t>
expectEq(__basicValue2kind(int64(1LL)), __kindInt)
expectEq(__basicValue2kind(int64(-1LL)), __kindInt)

expectEq(__basicValue2kind(uint(1)), __kindUint)
expectEq(__basicValue2kind(uint(-1)), __kindUint)
expectEq(__basicValue2kind(uint8(-3)), __kindUint8)
expectEq(__basicValue2kind(uint8(3)), __kindUint8)
expectEq(__basicValue2kind(uint16(0)), __kindUint16)
expectEq(__basicValue2kind(uint16(-1)), __kindUint16)
expectEq(__basicValue2kind(uint32(1)), __kindUint32)
expectEq(__basicValue2kind(uint32(-1)), __kindUint32)

-- can't distinguish __kindUint from __kindUint64
-- they are both ctype<uint64_t>
expectEq(__basicValue2kind(uint64(1)), __kindUint)
expectEq(__basicValue2kind(uint64(-1)), __kindUint)

expectEq(__basicValue2kind(float32(-1.0)), __kindFloat32)
expectEq(__basicValue2kind(float64(-1.0)), __kindFloat64)


-- pointers

a = __Int(4) -- currently, even integers are wrapped.

-- b := &a  -- gets translated as two parts:
ptrType = __ptrType(__Int)

b = ptrType(function() return a; end, function(__v) a = __v; end, a);

-- arrays

arrayType = __arrayType(__Int, 2);

a = arrayType()
expectEq(a[0], 0LL)
a[1] = 32LL
expectEq(a[1], 32LL)
expectEq(#a, 2LL)

