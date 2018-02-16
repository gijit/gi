
--[[

__nativeArray = function(elemKind)

   if false then
      if elemKind ==  __kindInt then 
         return Int32Array; -- in js, a builtin typed array
      elseif elemKind ==  __kindInt8 then 
         return Int8Array;
      elseif elemKind ==  __kindInt16 then 
         return Int16Array;
      elseif elemKind ==  __kindInt32 then 
         return Int32Array;
      elseif elemKind ==  __kindUint then 
         return Uint32Array;
      elseif elemKind ==  __kindUint8 then 
         return Uint8Array;
      elseif elemKind ==  __kindUint16 then 
         return Uint16Array;
      elseif elemKind ==  __kindUint32 then 
         return Uint32Array;
      elseif elemKind ==  __kindUintptr then 
         return Uint32Array;
      elseif elemKind ==  __kindFloat32 then 
         return Float32Array;
      elseif elemKind ==  __kindFloat64 then 
         return Float64Array;
      else
         return Array;
      end
   end
end;

__toNativeArray = function(elemKind, array)
  local nativeArray = __nativeArray(elemKind);
  if nativeArray == Array {
    return array;
  end
  return nativeArray(array); -- new
end;

--]]

