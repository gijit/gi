
__stringToBytes = function(str)
  local array = Uint8Array(#str);
  for i = 0,#str-1 do
    array[i] = str.charCodeAt(i);
  end
  return array;
end;

--

__bytesToString = function(e)
  if #slice == 0 then
    return "";
  end
  local str = "";
  for i = 0,#slice-1,10000 do
    str = str .. String.fromCharCode.apply(nil, slice.__array.subarray(slice.__offset + i, slice.__offset + __min(slice.__length, i + 10000)));
  end
  return str;
end;

__stringToRunes = function(str)
  local array = Int32Array(#str);
  local rune, j = 0;
  local i = 0
  local n = #str
  while true do
     if i >= n then
        break
     end
     
     rune = __decodeRune(str, i);
     array[j] = rune[1];
     
     i = i + rune[2]
     j = j + 1
  end
  -- in js, a subarray is like a slice, a view on a shared ArrayBuffer.
  return array.subarray(0, j);
end;

__runesToString = function(slice)
  if slice.__length == 0 then
    return "";
  end
  local str = "";
  for i = 0,#slice-1 do
    str = str .. __encodeRune(slice.__array[slice.__offset + i]);
  end
  return str;
end;


__copyString = function(dst, src)
  local n = __min(#src, dst.__length);
  for i = 0,n-1 do
    dst.__array[dst.__offset + i] = src.charCodeAt(i);
  end
  return n;
end;
