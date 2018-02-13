function __decodeRune(s, i)
   return {__utf8.sub(s, i+1, i+1), 1}
end

-- from gopherjs, ported to use bit ops.
__bit =require("bit")

--[[

-- js op precedence: higher precendence = tighter binding.
--
-- arshift/lshift: 12 left-to-right   
-- band: 9    left-to-right
-- bor : 7    left-to-right


__decodeRune = function(str, pos)
  local c0 = str.charCodeAt(pos);

  if c0 < 0x80 then
    return {c0, 1};
  end

  if c0 ~= c0  or  c0 < 0xC0 then
    return {0xFFFD, 1};
  end

  local c1 = str.charCodeAt(pos + 1);
  if c1 ~= c1  or  c1 < 0x80  or  0xC0 <= c1 then
    return {0xFFFD, 1};
  end

  if c0 < 0xE0 then
     local r = __bit.bor(__bit.lshift(__bit.band(c0, 0x1F), 6), __bit.band(c1, 0x3F));
    if r <= 0x7F then
      return {0xFFFD, 1};
    end
    return {r, 2};
  end

  local c2 = str.charCodeAt(pos + 2);
  if c2 ~= c2  or  c2 < 0x80  or  0xC0 <= c2 then
    return {0xFFFD, 1};
  end

  if c0 < 0xF0 then
   local r = __bit.bor(__bit.bor(__bit.lshift(__bit.band(c0, 0x0F), 12), __bit.lshift(__bit.band(c1, 0x3F), 6)), __bit.band(c2, 0x3F));
   
    if r <= 0x7FF then
      return {0xFFFD, 1};
    end
    if 0xD800 <= r  and  r <= 0xDFFF then
      return {0xFFFD, 1};
    end
    return {r, 3};
  end

  local c3 = str.charCodeAt(pos + 3);
  if c3 ~= c3  or  c3 < 0x80  or  0xC0 <= c3 then
    return {0xFFFD, 1};
  end

  if c0 < 0xF8 then
    local r = __bit.bor(__bit.bor(__bit.bor(__bit.lshift(__bit.band(c0, 0x07),18), __bit.lshift(__bit.band(c1, 0x3F), 12)), __bit.lshift(__bit.band(c2, 0x3F), 6), __bit.band(c3, 0x3F)));

    if r <= 0xFFFF  or  0x10FFFF < r then
      return {0xFFFD, 1};
    end
    return {r, 4};
  end

  return {0xFFFD, 1};
end;

__encodeRune = function(r)
  if r < 0  or  r > 0x10FFFF  or  (0xD800 <= r  and  r <= 0xDFFF) then
    r = 0xFFFD;
  end
  if r <= 0x7F then
    return String.fromCharCode(r);
  end
  if r <= 0x7FF then
    return String.fromCharCode(__bit.bor(0xC0, __bit.arshift(r,6)), __bit.bor(0x80, __bit.band(r, 0x3F)));
  end
  if r <= 0xFFFF then
   return String.fromCharCode(__bit.bor(0xE0, __bit.arshift(r,12)), __bit.bor(0x80, (__bit.band(__bit.arshift(r,6), 0x3F))), __bit.bor(0x80, __bit.band(r, 0x3F)));
  end
   return String.fromCharCode(__bit.bor(0xF0, __bit.arshift(r, 18)), __bit.bor(0x80, __bit.band(__bit.arshift(r,12),0x3F)), __bit.bor(0x80, __bit.band(__bit.arshift(r,6), 0x3F)), __bit.bor(0x80, __bit.band(r, 0x3F)));
end;

--]]
