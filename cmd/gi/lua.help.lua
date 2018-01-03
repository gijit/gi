-- debugging locals notes:

-- Notice that in the lua REPL, each line is a separate chunk with separate locals. Also, internal variables are returned (names start with '(' if you want to remove them):

local a = 2; for x, v in pairs(locals()) do print(x, v) end



-- https://stackoverflow.com/questions/2834579/print-all-local-variables-accessible-to-the-current-scope-in-lua

function locals()
  local variables = {}
  local idx = 1
  while true do
    local ln, lv = debug.getlocal(2, idx)
    if ln ~= nil then
      variables[ln] = lv
    else
      break
    end
    idx = 1 + idx
  end
  return variables
end

local a = 2; for x, v in pairs(locals()) do print(x, v) end

-- Upvalues are local variables from outer scopes, that are used in the current function. They are neither in _G nor in locals()

function upvalues()
  local variables = {}
  local idx = 1
  local func = debug.getinfo(2, "f").func
  while true do
    local ln, lv = debug.getupvalue(func, idx)
    if ln ~= nil then
      variables[ln] = lv
    else
      break
    end
    idx = 1 + idx
  end
  return variables
end

Example (notice you have to use a for it to show up):

> local a= 2; function f() local b = a; for x,v in pairs(upvalues()) do print(x,v) end end; f()

-- See debug.getlocal:

local foobar = 1

local i = 0
repeat
    local k, v = debug.getlocal(1, i)
    if k then
        print(k, v)
        i = i + 1
    end
until nil == k

