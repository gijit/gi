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
                    this.Claw = Claw_;
                    this.HasRing = HasRing_;
end);

sliceType = __sliceType(__emptyInterface);
ptrType = __ptrType(hobbit);
ptrType__1 = __ptrType(Wolf);
hobbit.ptr.methodSet.WearRing = function(this)
   h = this;
   h.hasRing = not h.hasRing;
   return h.hasRing;
end

hobbit.methodSet.WearRing = hobbit.ptr.methodSet.WearRing 


Wolf.ptr.methodSet.Scowl = function(this) 
   w = this;
   w.Claw = w.Claw + 1LL;
   return w.Claw;
end

Wolf.methodSet.Scowl = Wolf.ptr.methodSet.Scowl

battle = function(g, b) 
   return g:Scowl(), b:WearRing();
end
   
ptrType.methods = {{prop= "WearRing", name= "WearRing", pkg= "", typ= __funcType({}, {__Bool}, false)}};

ptrType__1.methods = {{prop= "Scowl", name= "Scowl", pkg= "", typ= __funcType({}, {__Int}, false)}};

Baggins.init({{prop= "WearRing", name= "WearRing", pkg= "", typ= __funcType({}, {__Bool}, false)}});

Gollum.init({{prop= "Scowl", name= "Scowl", pkg= "", typ= __funcType({}, {__Int}, false)}});

hobbit.init("github.com/gijit/gi/pkg/compiler/tmp", {{prop= "hasRing", name= "hasRing", anonymous= false, exported= false, typ= __Bool, tag= ""}});

Wolf.init("", {{prop= "Claw", name= "Claw", anonymous= false, exported= true, typ= __Int, tag= ""}, {prop= "HasRing", name= "HasRing", anonymous= false, exported= true, typ= __Bool, tag= ""}});

tryTheTypeSwitch = function(i)
   x, isG = __assertType(i, Gollum, true)
   if isG then
      return x.Scowl()
   end
   
   x, isB = __assertType(i, Baggins, true)
   if isB then
      if x.WearRing() then
         return 1
      end
   end
   return 0
end

-- main
w = Wolf.ptr(0, false);
bilbo = hobbit.ptr(false);
i0, b0 = battle(w, bilbo);
i1, b1 = battle(w, bilbo);
try0 = tryTheTypeSwitch(w);
try1 = tryTheTypeSwitch(bilbo);

