
_giPrivatePointerProps = _giPrivatePointerProps or {}

-- metatable for pointers

__gi_PrivatePointerMt = {

    __newindex = function(t, k, v)
       print("__gi_Pointer: __newindex called val=", v)
       local props = rawget(t, _giPrivatePointerProps)
       return props.set(v)
    end,

    __index = function(t, k)
       print("__gi_Pointer: __index called for key", k)       
       local props = rawget(t, _giPrivatePointerProps)
       return props.get()
    end,

    __tostring = function(t)
       print("__gi_Pointer: tostring called")
       local props = rawget(t, _giPrivatePointerProps)
    end
 }


-- getter and setter are closures
function __gi_ptrType(getter, setter)
   if getter == nil then
      error "__gi_ptrType sees nil getter"
   end
   if setter == nil then
      error "__gi_ptrType sees nil setter"
   end
   local proxy = {}
   proxy[_giPrivatePointerProps] = {["get"]=getter, ["set"]=setter}
   setmetatable(proxy, __gi_PrivatePointerMt)
   return proxy
end
