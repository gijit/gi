dofile 'fin.lua'

-- fin and fin_test use ___ triple underscores, to
-- avoid collision while integrating with struct.lua

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

-- ___mapAndJoinStrings(splice, arr, fun)
expectEq("", ___mapAndJoinStrings("_", {}, function(x) return x end))
expectEq("a_b_c", ___mapAndJoinStrings("_", {"a","b","c"}, function(x) return x end))
expectEq("a", ___mapAndJoinStrings("_", {"a"}, function(x) return x end))
expectEq("a", ___mapAndJoinStrings("_", {[0]="a"}, function(x) return x end))
expectEq("a_b", ___mapAndJoinStrings("_", {[0]="a","b"}, function(x) return x end))
expectEq("1_2_3", ___mapAndJoinStrings("_", {[0]=0,1,2}, function(x)  return x+1 end))
expectEq("1_2", ___mapAndJoinStrings("_", {[0]=0,1}, function(x)  return x+1 end))
expectEq("1", ___mapAndJoinStrings("_", {[0]=0}, function(x)  return x+1 end))
expectEq("", ___mapAndJoinStrings("_", {}, function(x)  return x+1 end))

-- ___keys = function(m)
expectEq({}, ___keys({}))
expectEq({"a","b"}, ___keys({b=2, a=4}))
expectEq({"a"}, ___keys({a=4}))
expectEq({0}, ___keys({[0]=4}))
expectEq({1}, ___keys({[1]=3}))
expectEq({7,11}, ___keys({[7]="seven", [11]="eleven"}))
local f = function() end
local g = function() end
expectEq({tostring(f)}, ___keys({[f]="seven"}))
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
expectEq(x, ___keys({[f]="seven", [g]="eight"}))

-- basic types, zero values
expectEq("0LL", tostring(___Int()))
expectEq("0ULL", tostring(___Uint()))
expectEq("0", tostring(___Float64()))
expectEq('""', tostring(___String()))

-- basic types, non-zero values
expectEq("-43LL", tostring(___Int(-43LL)))
expectEq("42ULL", tostring(___Uint(42ULL)))
expectEq("-0.3", tostring(___Float64(-0.3)))
expectEq('"hello world"', tostring(___String("hello world")))

                                                                                   
                                         
expectEq(___basicValue2kind("hi"), ___kindString)
expectEq(___basicValue2kind(""), ___kindString)
expectEq(___basicValue2kind(true), ___kindBool)
expectEq(___basicValue2kind(false), ___kindBool)

expectEq(___basicValue2kind(1LL), ___kindInt)
expectEq(___basicValue2kind(-1LL), ___kindInt)
expectEq(___basicValue2kind(int8(-3)), ___kindInt8)
expectEq(___basicValue2kind(int8(3)), ___kindInt8)
expectEq(___basicValue2kind(int16(0)), ___kindInt16)
expectEq(___basicValue2kind(int16(-1)), ___kindInt16)
expectEq(___basicValue2kind(int32(1)), ___kindInt32)
expectEq(___basicValue2kind(int32(-1)), ___kindInt32)

-- can't distinguish ___kindInt from ___kindInt64
-- they are both ctype<int64_t>
expectEq(___basicValue2kind(int64(1LL)), ___kindInt)
expectEq(___basicValue2kind(int64(-1LL)), ___kindInt)

expectEq(___basicValue2kind(uint(1)), ___kindUint)
expectEq(___basicValue2kind(uint(-1)), ___kindUint)
expectEq(___basicValue2kind(uint8(-3)), ___kindUint8)
expectEq(___basicValue2kind(uint8(3)), ___kindUint8)
expectEq(___basicValue2kind(uint16(0)), ___kindUint16)
expectEq(___basicValue2kind(uint16(-1)), ___kindUint16)
expectEq(___basicValue2kind(uint32(1)), ___kindUint32)
expectEq(___basicValue2kind(uint32(-1)), ___kindUint32)

-- can't distinguish ___kindUint from ___kindUint64
-- they are both ctype<uint64_t>
expectEq(___basicValue2kind(uint64(1)), ___kindUint)
expectEq(___basicValue2kind(uint64(-1)), ___kindUint)

expectEq(___basicValue2kind(float32(-1.0)), ___kindFloat32)
expectEq(___basicValue2kind(float64(-1.0)), ___kindFloat64)


-- pointers

a = ___Int(4) -- currently, even integers are wrapped.

-- b := &a  -- gets translated as two parts:
ptrType = ___ptrType(___Int)

b = ptrType(function() return a; end, function(___v) a = ___v; end, a);

-- arrays

arrayType = ___arrayType(___Int, 4);

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

slcInt = ___sliceType(___Int)

sl = ___makeSlice(slcInt, 3, 4)

s0 = ___subslice(sl, 2)
sl[2] = 45LL
expectEq(s0[0], 45LL)

-- copy, append

s2 = ___makeSlice(slcInt, 3)
m = ___copySlice(s2, sl)
expectEq(s2[2], 45LL)
expectEq(m, 3)

s0[0]=100LL
s2[0]=101LL
s2[1]=102LL
s2[2]=103LL

ap = ___appendSlice(s0, s2)
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

WonderWoman = ___newType(0, ___kindStruct, "main.WonderWoman", true, "github.com/gijit/gi/pkg/compiler/tmp", true, function(self, ...)
                           print("DEBUG WonderWoman.ctor called! dots=")
                           ___st({...})
                           if self == nil then self = {}; end
                           
                           local Bracelets_, LassoPoints_ = ... ;
                           -- for each zero value that is not a nil pointer:
                           self.Bracelets = Bracelets_ or (0LL);
                           self.LassoPoints = LassoPoints_ or (0LL);
                           return self; 
end);


WonderWoman.init("", {{prop= "Bracelets", name= "Bracelets", anonymous= false, exported= true, typ= ___Int, tag= ""}, {prop= "LassoPoints", name= "LassoPoints", anonymous= false, exported= true, typ= ___Int, tag= ""}});

ww = WonderWoman.ptr(2LL);

expectEq(ww.Bracelets, 2LL)

ptrType = ___ptrType(WonderWoman);
WonderWoman.ptr.methodSet.Fight = function(___self)
   local w = ___self;
   w.LassoPoints = w.LassoPoints + 1LL;
end

-- WonderWoman.ptr.methodSet.Fight = WonderWoman.methodSet.Fight

ww:Fight()
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

sliceType = ___sliceType(___emptyInterface);
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


Baggins= ___newType(8, ___kindInterface, "main.Baggins", true, "github.com/gijit/gi/pkg/compiler/tmp", true, null);
Gollum = ___newType(8, ___kindInterface, "main.Gollum", true, "github.com/gijit/gi/pkg/compiler/tmp", true, null);

hobbit = ___newType(0, ___kindStruct, "main.hobbit", true, "github.com/gijit/gi/pkg/compiler/tmp", false, function(this, hasRing_)
		this.___val = this; -- signal a non-basic value?
		if hasRing_ == nil then
			this.hasRing = false;
			return;
		end
		this.hasRing = hasRing_;
end);

Wolf = ___newType(0, ___kindStruct, "main.Wolf", true, "github.com/gijit/gi/pkg/compiler/tmp", true, function(this, Claw_, HasRing_) 
                    this.___val = this;
                    if HasRing_ == nil then
                       this.Claw = 0;
                       this.HasRing = false;
                       return;
                    end
                    this.Claw = Claw_;
                    this.HasRing = HasRing_;
end);

sliceType = ___sliceType(___emptyInterface);
ptrType = ___ptrType(hobbit);
ptrType___1 = ___ptrType(Wolf);
hobbit.ptr.methodSet.WearRing = function(this)
   print("hobbit.ptr.methodSet.WearRing called!")
   h = this;
   h.hasRing = not h.hasRing;
   return h.hasRing;
end

hobbit.methodSet.WearRing = function(this)
   print("hobbit.methodSet.WearRing called!")
   return this.___val.WearRing(this);
end

Wolf.ptr.methodSet.Scowl = function(this)
   print("Wolf.ptr.methodSet.Scowl called!")
   w = this;
   w.Claw = w.Claw + 1LL;
   return w.Claw;
end

Wolf.methodSet.Scowl = function(this)
   print("Wolf.methodSet.Scowl called!")
   return this.___val.Scowl(this);
end

battle = function(g, b) 
   return g:Scowl(), b:WearRing();
end
   
ptrType.methods = {{prop= "WearRing", name= "WearRing", pkg= "", typ= ___funcType({}, {___Bool}, false)}};

ptrType___1.methods = {{prop= "Scowl", name= "Scowl", pkg= "", typ= ___funcType({}, {___Int}, false)}};

Baggins.init({{prop= "WearRing", name= "WearRing", pkg= "", typ= ___funcType({}, {___Bool}, false)}});

Gollum.init({{prop= "Scowl", name= "Scowl", pkg= "", typ= ___funcType({}, {___Int}, false)}});

hobbit.init("github.com/gijit/gi/pkg/compiler/tmp", {{prop= "hasRing", name= "hasRing", anonymous= false, exported= false, typ= ___Bool, tag= ""}});

Wolf.init("", {{prop= "Claw", name= "Claw", anonymous= false, exported= true, typ= ___Int, tag= ""}, {prop= "HasRing", name= "HasRing", anonymous= false, exported= true, typ= ___Bool, tag= ""}});

tryTheTypeSwitch = function(i)
   print("top of tryTheTypeSwitch, with i=")
   ___st(i)
   
   x, isG = ___assertType(i, Gollum, true)
   if isG then
      print("yes, i satisfies Gollum interface")
      return x:Scowl()
   end
   print("i did not satisfy Gollum, trying Baggins...")
   
   x, isB = ___assertType(i, Baggins, true)
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
--     ptrType___1.methods

-- the Go spec says:
-- The method set of the corresponding pointer type *T is
-- the set of all methods declared with receiver *T or T
-- (that is, it also contains the method set of T).
--
-- So pointers should check their own and their elem methods.
--  but we'll need to clone the value before calling a value-receiver method with a pointer.

msWp = ___methodSet(Wolf.ptr)
expectEq(#msWp, 1)
msW = ___methodSet(Wolf)
expectEq(#msW, 0)

msHp = ___methodSet(hobbit.ptr)
expectEq(#msHp, 1)
msH = ___methodSet(hobbit)
expectEq(#msH, 0)

w2 = Wolf.ptr(0, false);
expectEq(getmetatable(w2).___name, "methodSet for *main.Wolf")

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

Hound = ___newType(0, ___kindStruct, "main.Hound", true, "github.com/gijit/gi/pkg/compiler/tmp", true, function(this, ...) 
                        this.___val = this;
                        local Name_, Id_, Mate_, Litter_, PtrLit_, food_, ate_ = ...
                        local args = {...}
		if #args == 0  then
			this.Name = "";
			this.Id = 0;
			this.Mate = ptrType.___nil;
			this.Litter = sliceType.___nil;
			this.PtrLit = ptrType___1.___nil;
			this.food = 0;
			this.ate = false;
			return;
		end
		this.Name = Name_ or "";
		this.Id = Id_ or 0LL;
		this.Mate = Mate_ or ptrType.___nil;
		this.Litter = Litter_;
		this.PtrLit = PtrLit_;
		this.food = food_ or 0LL;
		this.ate = ate_ or false;
end);

Hound.init("github.com/gijit/gi/pkg/compiler/tmp", {{prop= "Name", name= "Name", anonymous= false, exported= true, typ= ___String, tag= ""}, {prop= "Id", name= "Id", anonymous= false, exported= true, typ= ___Int, tag= ""}, {prop= "Mate", name= "Mate", anonymous= false, exported= true, typ= ptrType, tag= ""}, {prop= "Litter", name= "Litter", anonymous= false, exported= true, typ= sliceType, tag= ""}, {prop= "PtrLit", name= "PtrLit", anonymous= false, exported= true, typ= ptrType___1, tag= ""}, {prop= "food", name= "food", anonymous= false, exported= false, typ= ___Int, tag= ""}, {prop= "ate", name= "ate", anonymous= false, exported= false, typ= ___Bool, tag= ""}});

-- replace .prototype with .methodSet
Hound.ptr.methodSet.Eat = function(this, a)
   print("Eat called, with a = ", a, " and this=")
   ___st(this,"this-on-Hound.ptr")
   
		local _i, _ref, h, pup;
		h = this;
		if h.ate then
                   return h.food;
		end
		h.ate = true;
		h.food = h.food + a;
		_ref = h.Litter;
                ___st(_ref, "_ref")
		_i = 0;
		while true do
                   if not (_i < #_ref) then break; end

                   if (_i < 0 or _i >= #_ref) then
                      ___throwRuntimeError("index out of range")
                   end
                   ___st(_ref.___array, "_ref.___array")
                   print("_i = ", _i)
                   pup = _ref.___array[_ref.___offset + _i + 1]; -- + 1 for lua's arrays
                   pup:Eat(a);
			_i=_i+1;
                        end
		if not (h.Mate == ptrType.___nil) then
			h.Mate:Eat(a);
		end
		return h.food;
	end;
	Hound.methodSet.Eat = function(a)  return this.___val.Eat(a); end;

	ptrType = ___ptrType(Hound);
	sliceType = ___sliceType(ptrType);
	ptrType___1 = ___ptrType(sliceType);
	sliceType___1 = ___sliceType(___emptyInterface);        

		jake =  Hound.ptr("Jake", 123, ptrType.___nil, sliceType.___nil, ptrType___1.___nil, 0, false);
		joy =  Hound.ptr("Joy", 456, ptrType.___nil, sliceType.___nil, ptrType___1.___nil, 0, false);
		bubbles =  Hound.ptr("Bubbles", 2, ptrType.___nil, sliceType.___nil, ptrType___1.___nil, 0, false);
		barley =  Hound.ptr("Barley", 3, ptrType.___nil, sliceType.___nil, ptrType___1.___nil, 0, false);
		jake.Mate = joy;
		joy.Mate = jake;
		joy.Litter =  sliceType({bubbles, barley});
		jake.PtrLit = ptrType___1(function() return this.___target.Litter; end, function(___v)  this.___target.Litter = ___v; end, joy);
		got = joy:Eat(2);


		clone = ___clone(barley, Hound);
		print("clone.food =", clone.food)
		print("clone.Name =", clone.Name)
                
		print("joy:Eat(2) returned =", got)
		print("jake.food =",  jake.food)
		print("joy.food =",  joy.food)
		print("bubbles.food =",  bubbles.food)
		print("barley.food =",  barley.food)

-- end joy/jake puppies

                -- notice that structs have the ___get, ___set functions, and the ___val table.
        -- what are these/do they work?/ should they live in the struct on on a related table?
        -- they are related to pointer read/writes, and conversions;
        -- StarExpr invokes ___get, 
        
--[[
this-on-Hound.ptr: ============================ table: 0x000a8720
this-on-Hound.ptr:  1 key: 'ate' val: 'false'
this-on-Hound.ptr:  2 key: 'Mate' val: 'table: 0x000643e0'
this-on-Hound.ptr:  3 key: '___get' val: 'function: 0x000a84f8'
this-on-Hound.ptr:  4 key: 'Id' val: '123'
this-on-Hound.ptr:  5 key: '___set' val: 'function: 0x00064418'
this-on-Hound.ptr:  6 key: 'Litter' val: '<this.___val == this; avoid inf loop>'
this-on-Hound.ptr:  7 key: 'Name' val: 'Jake'
this-on-Hound.ptr:  8 key: 'food' val: '0'
this-on-Hound.ptr:  9 key: 'PtrLit' val: '<this.___val == this; avoid inf loop>'
this-on-Hound.ptr:  10 key: '___val' val: 'table: 0x000a8720'
--]]
                
print("done with fin_test.lua")
