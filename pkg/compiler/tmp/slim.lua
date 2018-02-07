
__type__Ragdoll = __gi_NewType(0, __gi_kind_Struct, "main", "Ragdoll", "main.Ragdoll", true, "main", true, nil)

print("just prior to defining anon_ptrType")
anon_ptrType = __ptrType(__type__Ragdoll); -- utils.go:490 immediate anon type printing.
__st(anon_ptrType, "anon_ptrType")

__type__Ragdoll.__init("", {{__prop= "Andy", __name= "Andy", __anonymous= false, __exported= true, __typ= anon_ptrType, __tag= ""}});

__type__Ragdoll.__constructor = function(self, ...) 
   if self == nil then self = {}; end
   local args={...};
   if #args == 0 then
      self.Andy = anon_ptrType.__nil;
   else 
      local Andy_ = ... ;
      self.Andy = Andy_;
   end
   print("Ragdoll ctor returing self=")
   __st(self, "self")
   return self; 
end

--__st(anon_ptrType, "anon_ptrType *prior* to __ctor set")
--rawset(anon_ptrType, "__constructor", __type__Ragdoll.__constructor)
--__st(anon_ptrType, "anon_ptrType *after* to __ctor set")

doll = __type__Ragdoll.__ptr({}, anon_ptrType.__nil);

doll.Andy = doll;
same = doll.Andy == doll;
