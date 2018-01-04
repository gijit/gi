package compiler

import (
	"fmt"
	"github.com/glycerine/luajit"
)

func DumpLuaStack(L *luajit.State) {
	var top int

	top = L.Gettop()
	for i := 1; i <= top; i++ {
		t := L.Type(i)
		switch t {
		case luajit.Tstring:
			fmt.Println("String : \t", L.Tostring(i))
		case luajit.Tboolean:
			fmt.Println("Bool : \t\t", L.Toboolean(i))
		case luajit.Tnumber:
			fmt.Println("Number : \t", L.Tonumber(i))
		default:
			fmt.Println("Type : \t\t", L.Typename(i))
		}
	}
	print("\n")
}

const GiLuaSliceMap = `
-- a Lua virtual table system suitable for use in slices and maps

-- create private index
_giPrivateSliceRaw = {}

 _giPrivateSliceMt = {

    __newindex = function(t, k, v)
      --print("newindex called for key", k)
      if t[_giPrivateSliceRaw][k] ~= nil then
          -- replace or delete
          if v == nil then 
              t.len = t.len - 1 -- delete
          else
              -- replace, no count change              
          end
      else 
          t.len = t.len + 1
      end
      t[_giPrivateSliceRaw][k] = v
    end,

  -- __index allows us to have fields to access these datacounts, *but*
  -- won't count as actual keys or indexes.
  --
    __index = function(t, k)
      --print("index called for key", k)
      if k == 'raw' then return t[_giPrivateSliceRaw] end
      return t[_giPrivateSliceRaw][k]
    end,

    __tostring = function(t)
       local s = "table of length" .. tostring(t.len) .. "{"
       local r = t[_giPrivateSliceRaw]
       -- we want to skip both the _giPrivateSliceRaw and the len
       -- when iterating.
       for i, _ in pairs(r) do s = s .. "["..tostring(i).."]" .. "= " .. tostring(r[i]) .. ", " end
       return s .. "}"
    end
 }

function _giSlice(x)
   assert(type(x) == 'table', 'bad parameter #1: must be table')

   local length = 0
   for k, v in pairs(x) do
      length = length + 1
   end

   local proxy = {len=length}
   proxy[_giPrivateSliceRaw] = x
   setmetatable(proxy, _giPrivateSliceMt)
   return proxy
end
`

const sampleUse = `
b=_giSlice {[0]=7,77,777}

assert(b.len == 3)
assert(b[0] == 7)
`
