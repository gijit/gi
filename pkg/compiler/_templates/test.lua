local jit = require('jit')
local format = string.format
local pcall = pcall

local function dofuss()
    return 42
end

local function w_() return pcall(dofuss) end

local function w() return w_() end
jit.off(w)

local ok, res = w()
print(format('ok=%s res=%s', ok, res))

-- let it JIT
for i = 1, 1000 do w() end

local ok, res = w()
print(format('ok=%s res=%s', ok, res))