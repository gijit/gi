dofile 'tsys.lua'
dofile 'tutil.lua'

-- __mapAndJoinStrings(splice, arr, fun)
__expectEq("", __mapAndJoinStrings("_", {}, function(x) return x end))
__expectEq("a_b_c", __mapAndJoinStrings("_", {"a","b","c"}, function(x) return x end))
__expectEq("a", __mapAndJoinStrings("_", {"a"}, function(x) return x end))
__expectEq("a", __mapAndJoinStrings("_", {[0]="a"}, function(x) return x end))
__expectEq("a_b", __mapAndJoinStrings("_", {[0]="a","b"}, function(x) return x end))
__expectEq("1_2_3", __mapAndJoinStrings("_", {[0]=0,1,2}, function(x)  return x+1 end))
__expectEq("1_2", __mapAndJoinStrings("_", {[0]=0,1}, function(x)  return x+1 end))
__expectEq("1", __mapAndJoinStrings("_", {[0]=0}, function(x)  return x+1 end))
__expectEq("", __mapAndJoinStrings("_", {}, function(x)  return x+1 end))

-- __keys = function(m)
__expectEq({}, __keys({}))
__expectEq({"a","b"}, __keys({b=2, a=4}))
__expectEq({"a"}, __keys({a=4}))
__expectEq({0}, __keys({[0]=4}))
__expectEq({1}, __keys({[1]=3}))
__expectEq({7,11}, __keys({[7]="seven", [11]="eleven"}))
local f = function() end
local g = function() end
__expectEq({tostring(f)}, __keys({[f]="seven"}))
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
__expectEq(x, __keys({[f]="seven", [g]="eight"}))

-- basic types, zero values
__expectEq("0LL", tostring(__type__.int()))
__expectEq("0ULL", tostring(__type__.uint()))
__expectEq("0", tostring(__type__.float64()))
__expectEq('""', tostring(__type__.string()))

-- basic types, non-zero values
__expectEq("-43LL", tostring(__type__.int(-43LL)))
__expectEq("42ULL", tostring(__type__.uint(42ULL)))
__expectEq("-0.3", tostring(__type__.float64(-0.3)))
__expectEq('"hello world"', tostring(__type__.string("hello world")))

                                                                                   
                                         
__expectEq(__basicValue2kind("hi"), __kindString)
__expectEq(__basicValue2kind(""), __kindString)
__expectEq(__basicValue2kind(true), __kindBool)
__expectEq(__basicValue2kind(false), __kindBool)

__expectEq(__basicValue2kind(1LL), __kindInt)
__expectEq(__basicValue2kind(-1LL), __kindInt)
__expectEq(__basicValue2kind(int8(-3)), __kindInt8)
__expectEq(__basicValue2kind(int8(3)), __kindInt8)
__expectEq(__basicValue2kind(int16(0)), __kindInt16)
__expectEq(__basicValue2kind(int16(-1)), __kindInt16)
__expectEq(__basicValue2kind(int32(1)), __kindInt32)
__expectEq(__basicValue2kind(int32(-1)), __kindInt32)

-- can't distinguish __kindInt from __kindInt64
-- they are both ctype<int64_t>
__expectEq(__basicValue2kind(int64(1LL)), __kindInt)
__expectEq(__basicValue2kind(int64(-1LL)), __kindInt)

__expectEq(__basicValue2kind(uint(1)), __kindUint)
__expectEq(__basicValue2kind(uint(-1)), __kindUint)
__expectEq(__basicValue2kind(uint8(-3)), __kindUint8)
__expectEq(__basicValue2kind(uint8(3)), __kindUint8)
__expectEq(__basicValue2kind(uint16(0)), __kindUint16)
__expectEq(__basicValue2kind(uint16(-1)), __kindUint16)
__expectEq(__basicValue2kind(uint32(1)), __kindUint32)
__expectEq(__basicValue2kind(uint32(-1)), __kindUint32)

-- can't distinguish __kindUint from __kindUint64
-- they are both ctype<uint64_t>
__expectEq(__basicValue2kind(uint64(1)), __kindUint)
__expectEq(__basicValue2kind(uint64(-1)), __kindUint)

__expectEq(__basicValue2kind(float32(-1.0)), __kindFloat32)
__expectEq(__basicValue2kind(float64(-1.0)), __kindFloat64)


-- pointers

a = __type__.int(4) -- currently, even integers are wrapped.

-- b := &a  -- gets translated as two parts:
ptrType = __ptrType(__type__.int)

b = ptrType(function() return a; end, function(__v) a = __v; end, a);

-- arrays

arrayType = __arrayType(__type__.int, 4);

a = arrayType()
__expectEq(a[0], 0LL)
a[1] = 32LL
__expectEq(a[1], 32LL)
__expectEq(#a, 4)

b = arrayType()
a[0]=5LL
arrayType.copy(b, a)

-- verify that arrayType.copy() worked.
__expectEq(b[0], 5LL)
__expectEq(b[1], 32LL)

-- slices

slcInt = __sliceType(__type__.int)

sl = __makeSlice(slcInt, 3, 4)

s0 = __subslice(sl, 2)
sl[2] = 45LL
__expectEq(s0[0], 45LL)

-- copy, append

s2 = __makeSlice(slcInt, 3)
m = __copySlice(s2, sl)
__expectEq(s2[2], 45LL)
__expectEq(m, 3LL)

s0[0]=100LL
s2[0]=101LL
s2[1]=102LL
s2[2]=103LL

ap = __appendSlice(s0, s2)
__expectEq(ap[0], 100LL)
__expectEq(ap[1], 101LL)
__expectEq(ap[2], 102LL)
__expectEq(ap[3], 103LL)
__expectEq(#ap, 4)



local indirectCallExampleStack = [=[
	r2:47: in function 'doubleRecover'
	r2:7: in function 'mustRecover'
	r2:64: in function 'f'
	[string "--..."]:23: in function '__top_of_defer'
]=]

local directCallExampleStack = [=[
	r2:63: in function 'f'
	[string "--..."]:23: in function '__top_of_defer'
]=]

__expectEq(__isDirectDefer(directCallExampleStack), true)
__expectEq(__isDirectDefer(indirectCallExampleStack), false)

print("done with fin_test.lua")
