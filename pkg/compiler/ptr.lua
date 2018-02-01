
__giPrivatePointerProps = __giPrivatePointerProps or {}

-- metatable for pointers

__gi_PrivatePointerMt = {

    __newindex = function(t, k, v)
       print("__gi_Pointer: __newindex called, calling set() with val=", v)
       local props = rawget(t, __giPrivatePointerProps)
       return props.set(v)
    end,

    __index = function(t, k)
       print("__gi_Pointer: __index called, doing get()")       
       local props = rawget(t, __giPrivatePointerProps)
       return props.get()
    end,

    __tostring = function(t)
       --print("__gi_Pointer: tostring called")
       local props = rawget(t, __giPrivatePointerProps)
       local typ = props.typ or "&unknownType"
       return typ .. "{" .. tostring(props.get()) .. "}"
    end
 }


-- getter and setter are closures
function __gi_ptrType(getter, setter, typeName)
   if getter == nil then
      error "__gi_ptrType sees nil getter"
   end
   if setter == nil then
      error "__gi_ptrType sees nil setter"
   end
   local proxy = {}
   proxy[__giPrivatePointerProps] = {
      ["get"]=getter,
      ["set"]=setter,
      ["typ"]=typeName,
   }
   setmetatable(proxy, __gi_PrivatePointerMt)
   return proxy
end
