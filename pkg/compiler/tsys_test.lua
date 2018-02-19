dofile 'tsys.lua'
dofile 'tutil.lua'

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
expectEq("0LL", tostring(__type__int()))
expectEq("0ULL", tostring(__type__uint()))
expectEq("0", tostring(__type__float64()))
expectEq('""', tostring(__type__string()))

-- basic types, non-zero values
expectEq("-43LL", tostring(__type__int(-43LL)))
expectEq("42ULL", tostring(__type__uint(42ULL)))
expectEq("-0.3", tostring(__type__float64(-0.3)))
expectEq('"hello world"', tostring(__type__string("hello world")))

                                                                                   
                                         
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

a = __type__int(4) -- currently, even integers are wrapped.

-- b := &a  -- gets translated as two parts:
ptrType = __ptrType(__type__int)

b = ptrType(function() return a; end, function(__v) a = __v; end, a);

-- arrays

arrayType = __arrayType(__type__int, 4);

a = arrayType()
expectEq(a[0], 0LL)
a[1] = 32LL
expectEq(a[1], 32LL)
expectEq(#a, 4LL)

b = arrayType()
a[0]=5LL
arrayType.copy(b, a)

-- verify that arrayType.copy() worked.
expectEq(b[0], 5LL)
expectEq(b[1], 32LL)

-- slices

slcInt = __sliceType(__type__int)

sl = __makeSlice(slcInt, 3, 4)

s0 = __subslice(sl, 2)
sl[2] = 45LL
expectEq(s0[0], 45LL)

-- copy, append

s2 = __makeSlice(slcInt, 3)
m = __copySlice(s2, sl)
expectEq(s2[2], 45LL)
expectEq(m, 3)

s0[0]=100LL
s2[0]=101LL
s2[1]=102LL
s2[2]=103LL

ap = __appendSlice(s0, s2)
expectEq(ap[0], 100LL)
expectEq(ap[1], 101LL)
expectEq(ap[2], 102LL)
expectEq(ap[3], 103LL)
expectEq(#ap, 4)


-- structs

--[[
package main
import (
	"fmt"
)
type WonderWoman struct {
	Bracelets int
    LassoPoints int
}
func (w *WonderWoman) Fight() {
   w.LassoPoints++
}
func main() {
	ww := WonderWoman{
		Bracelets: 2,
	}
    ww.Fight()
   fmt.Printf("ww=%#v\n", ww)
}
--]]

__type__WonderWoman = __newType(0, __kindStruct, "main.WonderWoman", true, "main", true, nil);

__type__WonderWoman.init("", {{__prop= "Bracelets", __name= "Bracelets", __anonymous= false, __exported= true, __typ= __type__int, __tag= ""}, {__prop= "LassoPoints", __name= "LassoPoints", __anonymous= false, __exported= true, __typ= __type__int, __tag= ""}});
	
__type__WonderWoman.__constructor = function(self, ...)
   print("DEBUG WonderWoman.ctor called! dots=")
   __st({...})   
   if self == nil then self = {}; end
   local Bracelets_, LassoPoints_ = ... ;
   self.Bracelets = Bracelets_ or 0LL;
   self.LassoPoints = LassoPoints_ or 0LL;
   return self; 
end;
;

__type__WonderWoman.ptr.prototype.Fight = function(w) 
   w.LassoPoints = w.LassoPoints + (1LL);
end;
__type__WonderWoman.prototype.Fight = function(this)  return this.__val.Fight(); end;

ww = __type__WonderWoman.ptr({}, 2LL, 0LL);

expectEq(ww.Bracelets, 2LL)

ww:Fight();
expectEq(ww.LassoPoints, 1LL)
ww:Fight()
expectEq(ww.LassoPoints, 2LL)
      

-- functions

--[[
package main
import (
	"fmt"
)
type Binop func(a, b int) int
func MyApply(bo Binop, x, y int) int {
	return bo(x, y)
}
func main() {
	res := MyApply(func(r, s int) int {
		return r + s
	}, 1, 2)
	fmt.Printf("res = %v\n", res)
}
--]]

sliceType = __sliceType(__type__emptyInterface);
MyApply = function(bo, x, y)
   return bo(x, y);
end
_r = MyApply(function(r, s)  return r + s; end, 1, 2);
expectEq(_r, 3)


--[[
package main
import (
	"fmt"
)
type Baggins interface {
	WearRing() bool
}
type Gollum interface {
	Scowl() int
}
type hobbit struct {
	hasRing bool
}
func (h *hobbit) WearRing() bool {
	h.hasRing = !h.hasRing
	return h.hasRing
}
type Wolf struct {
	Claw    int
	HasRing bool
}
func (w *Wolf) Scowl() int {
	w.Claw++
	return w.Claw
}
func battle(g Gollum, b Baggins) (int, bool) {
	return g.Scowl(), b.WearRing()
}
func tryTheTypeSwitch(i interface{}) int {
	switch x := i.(type) {
	case Gollum:
		return x.Scowl()
	case Baggins:
		if x.WearRing() {
			return 1
		}
	}
	return 0
}
func main() {
	w := &Wolf{}
	bilbo := &hobbit{}
	i0, b0 := battle(w, bilbo)
	i1, b1 := battle(w, bilbo)
	fmt.Printf("i0=%v, b0=%v\n", i0, b0)
	fmt.Printf("i1=%v, b1=%v\n", i1, b1)
	fmt.Printf("tried wolf=%v\n", tryTheTypeSwitch(w))
	fmt.Printf("tried bilbo=%v\n", tryTheTypeSwitch(bilbo))
}
/*
i0=1, b0=true
i1=2, b1=false
tried wolf=3
tried bilbo=1
*/
--]]


Baggins= __newType(8, __kindInterface, "main.Baggins", true, "github.com/gijit/gi/pkg/compiler/tmp", true, null);
Gollum = __newType(8, __kindInterface, "main.Gollum", true, "github.com/gijit/gi/pkg/compiler/tmp", true, null);

hobbit = __newType(0, __kindStruct, "main.hobbit", true, "github.com/gijit/gi/pkg/compiler/tmp", false, function(this, hasRing_)
		this.__val = this; -- signal a non-basic value?
		if hasRing_ == nil then
			this.hasRing = false;
			return;
		end
		this.hasRing = hasRing_;
end);

Wolf = __newType(0, __kindStruct, "main.Wolf", true, "github.com/gijit/gi/pkg/compiler/tmp", true, function(this, Claw_, HasRing_) 
                    this.__val = this;
                    
                    if HasRing_ == nil then
                       this.Claw = 0;
                       this.HasRing = false;
                       return;
                    end
                    this.Claw = Claw_ or 0;
                    this.HasRing = HasRing_ or false;
end);

sliceType = __sliceType(__type__emptyInterface);
ptrType = __ptrType(hobbit);
ptrType__1 = __ptrType(Wolf);
hobbit.ptr.prototype.WearRing = function(this)
   print("hobbit.ptr.prototype.WearRing called!")
   h = this;
   h.hasRing = not h.hasRing;
   return h.hasRing;
end

hobbit.prototype.WearRing = function(this)
   print("hobbit.prototype.WearRing called!")
   return this.__val.WearRing(this);
end

Wolf.ptr.prototype.Scowl = function(this)
   print("Wolf.ptr.prototype.Scowl called!")
   w = this;
   w.Claw = w.Claw + 1LL;
   return w.Claw;
end

Wolf.prototype.Scowl = function(this)
   print("Wolf.prototype.Scowl called!")
   return this.__val.Scowl(this);
end

battle = function(g, b) 
   return g:Scowl(), b:WearRing();
end
   
ptrType.methods = {{__prop= "WearRing", name= "WearRing", pkg= "", __typ= __funcType({}, {__type__bool}, false)}};

ptrType__1.methods = {{__prop= "Scowl", name= "Scowl", pkg= "", __typ= __funcType({}, {__type__int}, false)}};

Baggins.init({{__prop= "WearRing", name= "WearRing", pkg= "", __typ= __funcType({}, {__type__bool}, false)}});

Gollum.init({{__prop= "Scowl", name= "Scowl", pkg= "", __typ= __funcType({}, {__type__int}, false)}});

hobbit.init("github.com/gijit/gi/pkg/compiler/tmp", {{__prop= "hasRing", name= "hasRing", anonymous= false, exported= false, __typ= __type__bool, tag= ""}});

Wolf.init("", {{__prop= "Claw", name= "Claw", anonymous= false, exported= true, __typ= __type__int, tag= ""}, {__prop= "HasRing", name= "HasRing", anonymous= false, exported= true, __typ= __type__bool, tag= ""}});

tryTheTypeSwitch = function(i)
   print("top of tryTheTypeSwitch, with i=")
   __st(i)
   
   x, isG = __assertType(i, Gollum, true)
   if isG then
      print("yes, i satisfies Gollum interface")
      return x:Scowl()
   end
   print("i did not satisfy Gollum, trying Baggins...")
   
   x, isB = __assertType(i, Baggins, true)
   if isB then
      print("yes, i satisfies Baggins interface")
      if x:WearRing() then
         return 1LL
      end
   else
      print("i satisfied neither interface")
   end
   return 0LL
end

-- main
w = Wolf.ptr(0, false);
bilbo = hobbit.ptr(false);

-- problem:
-- hmm hobbit.methods and Wolf.methods are empty
-- but ptrType.methods is set above, as is
--     ptrType__1.methods

-- the Go spec says:
-- The method set of the corresponding pointer type *T is
-- the set of all methods declared with receiver *T or T
-- (that is, it also contains the method set of T).
--
-- So pointers should check their own and their elem methods.
--  but we'll need to clone the value before calling a value-receiver method with a pointer.

msWp = __methodSet(Wolf.ptr)
expectEq(#msWp, 1)
msW = __methodSet(Wolf)
expectEq(#msW, 0)

msHp = __methodSet(hobbit.ptr)
expectEq(#msHp, 1)
msH = __methodSet(hobbit)
expectEq(#msH, 0)

w2 = Wolf.ptr(0, false);
expectEq(getmetatable(w2).__name, "methodSet for *main.Wolf")

print("fin_test.lua: about to call battle(w, bilbo)")

i0, b0 = battle(w, bilbo);
i1, b1 = battle(w, bilbo);
try0 = tryTheTypeSwitch(w);
try1 = tryTheTypeSwitch(bilbo);

expectEq(i0, 1LL)
expectEq(b0, true)
expectEq(i1, 2LL)
expectEq(b1, false)
expectEq(try0, 3LL)
expectEq(try1, 1LL)

-- structs with pointers and slices

--[[
package main

import (
	"fmt"
)

type Hound struct {
	Name   string
	Id     int
	Mate   *Hound
	Litter []*Hound
	PtrLit *[]*Hound

	food int
	ate  bool
}

func (h *Hound) Eat(a int) int {
	if h.ate {
		return h.food
	}
	h.ate = true
	h.food += a

	for _, pup := range h.Litter {
		pup.Eat(a)
	}
	if h.Mate != nil {
		h.Mate.Eat(a)
	}
	return h.food
}

func main() {
	jake := &Hound{
		Name: "Jake",
		Id:   123,
	}
	joy := &Hound{
		Name: "Joy",
		Id:   456,
	}
	bubbles := &Hound{
		Name: "Bubbles",
		Id:   2,
	}
	barley := &Hound{
		Name: "Barley",
		Id:   3,
	}
	jake.Mate = joy
	joy.Mate = jake
	joy.Litter = []*Hound{bubbles, barley}
	jake.PtrLit = &(joy.Litter)

	got := joy.Eat(2)

	var clone Hound = *barley
	fmt.Printf("clone.food =%#v\n", clone.food)

	fmt.Printf("joy.Eat(2) returned =%#v\n", got)
	fmt.Printf("jake.food =%#v\n", jake.food)
	fmt.Printf("joy.food =%#v\n", joy.food)
	fmt.Printf("bubbles.food =%#v\n", bubbles.food)
	fmt.Printf("barley.food =%#v\n", barley.food)
}

/*
clone.food =2
joy.Eat(2) returned =2
jake.food =2
joy.food =2
bubbles.food =2
barley.food =2
*/

--]]


-- begin joy/jake puppies example

Hound = __newType(0, __kindStruct, "main.Hound", true, "github.com/gijit/gi/pkg/compiler/tmp", true, function(this, ...) 
                        this.__val = this;
                        local Name_, Id_, Mate_, Litter_, PtrLit_, food_, ate_ = ...
                        
			this.Name = Name_ or "";
			this.Id = Id_ or 0;
			this.Mate = Mate_ or ptrType.__nil;
			this.Litter = Litter_ or sliceType.__nil;
			this.PtrLit = PtrLit_ or ptrType__1.__nil;
			this.food = food_ or 0;
			this.ate = ate_ or false;
			return;
end);

Hound.init("github.com/gijit/gi/pkg/compiler/tmp", {{__prop= "Name", name= "Name", anonymous= false, exported= true, __typ= __type__string, tag= ""}, {__prop= "Id", name= "Id", anonymous= false, exported= true, __typ= __type__int, tag= ""}, {__prop= "Mate", name= "Mate", anonymous= false, exported= true, __typ= ptrType, tag= ""}, {__prop= "Litter", name= "Litter", anonymous= false, exported= true, __typ= sliceType, tag= ""}, {__prop= "PtrLit", name= "PtrLit", anonymous= false, exported= true, __typ= ptrType__1, tag= ""}, {__prop= "food", name= "food", anonymous= false, exported= false, __typ= __type__int, tag= ""}, {__prop= "ate", name= "ate", anonymous= false, exported= false, __typ= __type__bool, tag= ""}});

-- return .methodSet to  .prototype
Hound.ptr.prototype.Eat = function(this, a)
   print("Eat called, with a = ", a, " and this=")
   __st(this,"this-on-Hound.ptr")
   
		local _i, _ref, h, pup;
		h = this;
		if h.ate then
                   return h.food;
		end
		h.ate = true;
		h.food = h.food + a;
		_ref = h.Litter;
                __st(_ref, "_ref")
		_i = 0;
		while true do
                   if not (_i < #_ref) then break; end

                   if (_i < 0 or _i >= #_ref) then
                      __throwRuntimeError("index out of range")
                   end
                   __st(_ref.__array, "_ref.__array")
                   print("_i = ", _i)
                   pup = _ref.__array[_ref.__offset + _i + 1]; -- + 1 for lua's arrays
                   pup:Eat(a);
			_i=_i+1;
                        end
		if not (h.Mate == ptrType.__nil) then
			h.Mate:Eat(a);
		end
		return h.food;
	end;
	Hound.prototype.Eat = function(a)  return this.__val.Eat(a); end;

	ptrType = __ptrType(Hound);
	sliceType = __sliceType(ptrType);
	ptrType__1 = __ptrType(sliceType);
	sliceType__1 = __sliceType(__type__emptyInterface);        

		jake =  Hound.ptr("Jake", 123, ptrType.__nil, sliceType.__nil, ptrType__1.__nil, 0, false);
		joy =  Hound.ptr("Joy", 456, ptrType.__nil, sliceType.__nil, ptrType__1.__nil, 0, false);
		bubbles =  Hound.ptr("Bubbles", 2, ptrType.__nil, sliceType.__nil, ptrType__1.__nil, 0, false);
		barley =  Hound.ptr("Barley", 3, ptrType.__nil, sliceType.__nil, ptrType__1.__nil, 0, false);
		jake.Mate = joy;
		joy.Mate = jake;
		joy.Litter =  sliceType({bubbles, barley});
		jake.PtrLit = ptrType__1(function() return this.__target.Litter; end, function(__v)  this.__target.Litter = __v; end, joy);
		got = joy:Eat(2);


		clone = __clone(barley, Hound);
		print("clone.food =", clone.food)
		print("clone.Name =", clone.Name)
                
		print("joy:Eat(2) returned =", got)
		print("jake.food =",  jake.food)
		print("joy.food =",  joy.food)
		print("bubbles.food =",  bubbles.food)
		print("barley.food =",  barley.food)

-- end joy/jake puppies

                -- notice that structs have the __get, __set functions, and the __val table.
        -- what are these/do they work?/ should they live in the struct on on a related table?
        -- they are related to pointer read/writes, and conversions;
        -- StarExpr invokes __get, 
        
--[[
this-on-Hound.ptr: ============================ table: 0x000a8720
this-on-Hound.ptr:  1 key: 'ate' val: 'false'
this-on-Hound.ptr:  2 key: 'Mate' val: 'table: 0x000643e0'
this-on-Hound.ptr:  3 key: '__get' val: 'function: 0x000a84f8'
this-on-Hound.ptr:  4 key: 'Id' val: '123'
this-on-Hound.ptr:  5 key: '__set' val: 'function: 0x00064418'
this-on-Hound.ptr:  6 key: 'Litter' val: '<this.__val == this; avoid inf loop>'
this-on-Hound.ptr:  7 key: 'Name' val: 'Jake'
this-on-Hound.ptr:  8 key: 'food' val: '0'
this-on-Hound.ptr:  9 key: 'PtrLit' val: '<this.__val == this; avoid inf loop>'
this-on-Hound.ptr:  10 key: '__val' val: 'table: 0x000a8720'
--]]
                
print("done with fin_test.lua")
