--------------------------------------------------------------------------------
-- A general purpose library that extends Lua standard libraries.
-- 
-- Copyright (C) 2011-2016 Stefano Peluchetti. All rights reserved.
--
-- Features, documentation and more: http://www.scilua.org .
-- 
-- This file is part of the Xsys library, which is released under the MIT 
-- license: full text in file LICENSE.TXT in the library's root folder.
--------------------------------------------------------------------------------

-- TODO: Design exec API logging so that files are generated (useful for 
-- TODO: debugging and profiling).

local ffi      = require "ffi"
local bit      = require "bit"
-- CREDIT: Peter Colberg's templet library:
local templet  = require "xsys._dep.templet"

local select, pairs, error, setmetatable = select, pairs, error, setmetatable
local type, loadstring, setfenv, unpack = type, loadstring, setfenv, unpack
local pcall = pcall
local insert, concat = table.insert, table.concat
local format = string.format
local abs = math.abs

-- Table -----------------------------------------------------------------------
-- TODO: Introduce optional trailing 'resv = onconflict(key, v, newv)'.
local function union(...)
  local o = {}
  local arg, n = { ... }, select("#", ...)
  for a=1,n do
    for k,v in pairs(arg[a]) do
      if type(o[k]) ~= "nil" then
        error("key '"..tostring(k).."' is not unique among tables to be merged")
      end
      o[k] = v
    end
  end
  return o
end

local function append(...)
  local o = { }
  local arg, n = { ... }, select("#", ...)
  local c = 0
  for a=1,n do
    local t = arg[a]
    for i=1,#t do
      c = c + 1
      local v = t[i]
      if type(v) == "nil" then
        error("argument #"..a.." is not a proper array: no nil values allowed")
      end
      o[c] = v
    end
  end
  return o
end

-- Another module might have modified the standard libraries, overwrite:
local table = union(table)
table.union = union
table.append = append

-- Tonumber --------------------------------------------------------------------
local function getton(x)
  return x.__tonumber
end

local function tonumberx(x)
  if type(x) ~= "table" and type(x) ~= "cdata" then
    return tonumber(x)
  else
    local haston, ton = pcall(getton, x)
    return (haston and ton) and ton(x) or tonumber(x)
  end
end

-- String ----------------------------------------------------------------------
-- CREDIT: Steve Dovan snippet.
-- TODO: Clarify corner cases, make more robust.
local function split(s, re)
  local i1, ls = 1, { }
  if not re then re = '%s+' end
  if re == '' then return { s } end
  while true do
    local i2, i3 = s:find(re, i1)
    if not i2 then
      local last = s:sub(i1)
      if last ~= '' then insert(ls, last) end
      if #ls == 1 and ls[1] == '' then
        return  { }
      else
        return ls
      end
    end
    insert(ls, s:sub(i1, i2 - 1))
    i1 = i3 + 1
  end
end

-- TODO: what = "lr"
local function trim(s)
  return (s:gsub("^%s*(.-)%s*$", "%1"))
end

local function adjustexp(s)
  if s:sub(-3, -3) == "+" or s:sub(-3, -3) == "-" then
    return s:sub(1, -3).."0"..s:sub(-2)
  else
    return s
  end
end

local function width(x, chars)
  chars = chars or 9
  if chars < 9 then
    error("at least 9 characters required")
  end
  if type(x) == "nil" then
    return (" "):rep(chars - 3).."nil"
  elseif type(x) == "boolean" then
    local s = tostring(x)
    return (" "):rep(chars - #s)..s
  elseif type(x) == "string" then
    if #x > chars then
      return x:sub(1, chars - 2)..".."
    else
      return (" "):rep(chars - #x)..x
    end
  else
    local formatf = "%+"..chars.."."..(chars - 3).."f"
    local formate = "%+."..(chars - 8).."e"
    x = tonumberx(x) -- Could be cdata.
    local s = format(formatf, x)
    if x ~= x or abs(x) == 1/0 then return s end
    if tonumberx(s:sub(2, chars)) == 0 then -- It's small.
      if abs(x) ~= 0 then -- And not zero.
        s = adjustexp(format(formate, x))
      end
    else
      s = s:sub(1, chars)
      if not s:sub(3, chars - 1):find('%.') then -- It's big.
        s = adjustexp(format(formate, x))
      end
    end
    return s
  end
end

-- Another module might have modified the standard libraries, overwrite:
local string = union(string)
string.split = split
string.trim = trim
string.width = width

-- Exec ------------------------------------------------------------------------
local function testexec(chunk, chunkname, fenv, ok, ...)
  if not ok then
    local err = select(1, ...)
    error("execution error: "..err)
  end
  return ...
end

local function exec(chunk, chunkname, fenv)
  chunkname = chunkname or chunk
  local f, err = loadstring(chunk, chunkname)
  if not f then
    error("parsing error: "..err)
  end
  if fenv then 
    setfenv(f, fenv)
  end
  return testexec(chunk, chunkname, fenv, pcall(f))
end

-- From ------------------------------------------------------------------------
local function from(what, keystr)
  local keys = split(keystr, ",")
  local o = { }
  for i=1,#keys do
    o[i] = "x."..trim(keys[i])
  end
  o = concat(o, ",")
  local s = "return function(x) return "..o.." end"
  return exec(s, "from<"..keystr..">")(what)
end

-- Bit -------------------------------------------------------------------------
local tobit, lshift, rshift, band = bit.tobit, bit.lshift, bit.rshift, bit.band 

-- 99 == not used.
local lsb_array = ffi.new("const int32_t[64]", {32, 0, 1, 12, 2, 6, 99, 13,
  3, 99, 7, 99, 99, 99, 99, 14, 10, 4, 99, 99, 8, 99, 99, 25, 99, 99, 99, 99,  
  99, 21, 27, 15, 31, 11, 5, 99, 99, 99, 99, 99, 9, 99, 99, 24, 99, 99, 20, 26, 
  30, 99, 99, 99, 99, 23, 99, 19, 29, 99, 22, 18, 28, 17, 16, 99})

-- Compute position of least significant bit, starting with 0 for the tail of
-- the bit representation and ending with 31 for the head of the bit
-- representation (right to left). If all bits are 0 then 32 is returned.
-- This corresponds to finding the i in 2^i if the 4-byte value is set to 2^i.
-- Branch free version.
local function lsb(x)
  x = band(x, -x)
  x = tobit(lshift(x, 4) + x)
  x = tobit(lshift(x, 6) + x)
  x = tobit(lshift(x, 16) - x)
  return lsb_array[rshift(x, 26)]
end

-- Another module might have modified the standard libraries, overwrite:
local bit = union(bit)
bit.lsb = lsb -- TODO: Document.

-- Export ----------------------------------------------------------------------

return {
  template = templet.loadstring,
  exec     = exec,
  from     = from,
  tonumber = tonumberx,
  table    = table,
  string   = string,
  bit      = bit,
}
