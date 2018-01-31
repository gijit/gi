-- face.lua has interface support functions

__gi_PrivateInterfaceProps = __gi_PrivateInterfaceProps or {}

__gi_ifaceNil = __gi_ifaceNil or {[__gi_PrivateInterfaceProps]={name="nil"}}

-- __gi_assertType is an interface type assertion.
--
--  either
--    a, ok := b.(Face)  ## the two value form (returnTupe==2)
--  or
--    a := b.(Face)      ## the one value form (returnTuple==0; can panic)
--  or
--    _, ok := b.(Face)  ## (returnTuple==1; does not panic)
--
-- returnTuple in {0, 1, 2},
--   0 returns just the interface-value, converted or nil/zero-value.
--   1 returns just the ok (2nd-value in a conversion, a bool)
--   2 returns both
--
--   if 0, then we may panic if the interface conversion fails.
--
function __gi_assertType(value, typ, returnTuple)

   print("__gi_assertType called, typ=", typ, " value=", value, " returnTuple=", returnTuple)

   local isInterface = false
   if typ == "kindInterface" then
      isInterface = true
   end
   
   local ok = false;
   local missingMethod = "";
   
  if value == __gi_ifaceNil then
     ok = false;
     
  elseif not isInterface then
     ok = value.constructor == typ;
     
  else
     local valueTypeString = value.constructor.string;
     ok = typ.implementedBy[valueTypeString];
     if ok == undefined then
        
        ok = true;
        local valueMethodSet = __gi_methodSet(value.constructor);
        local interfaceMethods = typ.methods;
        local li = interfaceMethods.length
        
        for i = 0, li-1 do
           
           local tm = interfaceMethods[i];
           local found = false;
           local msl = valueMethodSet.length
           
           for j = 0,msl-1 do
              local vm = valueMethodSet[j];
              if vm.name == tm.name and vm.pkg == tm.pkg and vm.typ == tm.typ then
                 found = true;
                 break;
              end
           end
           
           if not found then
              ok = false;
              typ.missingMethodFor[valueTypeString] = tm.name;
              break;
           end
        end
        
        typ.implementedBy[valueTypeString] = ok;
        
     end
     if not ok then
        missingMethod = typ.missingMethodFor[valueTypeString];
     end
  end
  
  if not ok then
     
     if returnTupe == 0 then
        __gi_panic(new __gi_packages["runtime"].TypeAssertionError.ptr("", (value === __gi_ifaceNil ? "" : value.constructor.string), typ.string, missingMethod));
        
     elseif returnTuple == 1 then
        return false
     else
        return zeroVal, false
     end
  end
  
  if not isInterface then
     value = value.$val;
  end
  
  if typ == $jsObjectPtr then
     value = value.object;
  end
  
  if returnTupe == 0 then
     return value
  elseif returnTuple == 1 then
     return true
  end
  return value, true
  
end
