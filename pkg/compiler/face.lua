-- face.lua has interface support functions

__gi_PrivateInterfaceProps = __gi_PrivateInterfaceProps or {}

__gi_ifaceNil = __gi_ifaceNil or {[__gi_PrivateInterfaceProps]={name="nil"}}

function __gi_assertType(value, typ, returnTuple)

  var isInterface = (typ.kind === __gi_kindInterface), ok, missingMethod = "";
  if (value === __gi_ifaceNil) {
    ok = false;
  } else if (!isInterface) {
    ok = value.constructor === typ;
  } else {
    var valueTypeString = value.constructor.string;
    ok = typ.implementedBy[valueTypeString];
    if (ok === undefined) {
      ok = true;
      var valueMethodSet = __gi_methodSet(value.constructor);
      var interfaceMethods = typ.methods;
      for (var i = 0; i < interfaceMethods.length; i++) {
        var tm = interfaceMethods[i];
        var found = false;
        for (var j = 0; j < valueMethodSet.length; j++) {
          var vm = valueMethodSet[j];
          if (vm.name === tm.name && vm.pkg === tm.pkg && vm.typ === tm.typ) {
            found = true;
            break;
          }
        }
        if (!found) {
          ok = false;
          typ.missingMethodFor[valueTypeString] = tm.name;
          break;
        }
      }
      typ.implementedBy[valueTypeString] = ok;
    }
    if (!ok) {
      missingMethod = typ.missingMethodFor[valueTypeString];
    }
  }

  if (!ok) {
    if (returnTuple) {
      return [typ.zero(), false];
    }
    __gi_panic(new __gi_packages["runtime"].TypeAssertionError.ptr("", (value === __gi_ifaceNil ? "" : value.constructor.string), typ.string, missingMethod));
  }

  if (!isInterface) {
    value = value.$val;
  }
  if (typ === $jsObjectPtr) {
    value = value.object;
  }
  return returnTuple ? [value, true] : value;
};

end
