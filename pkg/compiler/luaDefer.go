package compiler

/*
translation of templates/flow9.go, turned into templates.
See templates/flow9.lua for execuatble lua.
*/

// strings defined here:

// deferGlobalOneTimeInitForDefer
//
// deferCustom0template
// deferBoilerplate1
//
// custom2 : the main body of the function. The
//           compiler code will custom generate all of this.
//           There is no template.
//
// deferBoilerplate3
// custom4tempalate
// deferBoilerplate5
// custom6template

const deferGlobalOneTimeInitForDefer = `
__recovMT = {__tostring = function(v) return 'a-panic-value:' .. tostring(v[1]) end}
`

const deferCustom0 = `
a = 0;
b = 0;
function f(...)
   orig = {...}

`
const deferCustom0template = `
%v
function %v(...)
   orig = {...}

`

const deferBoilerplate1 = `
  -- deferBoilerplate1:

  local __defers={}

  -- __recoverVal will be nill if no panic,
  --              or if panic happened and
  --              then was recovered.
  -- this is always a table to avoid
  --  stringification problems. The real
  --  panic value is inside at position [1].
  --
  local __recoverVal

  local recover = function()
       local cp = __recoverVal
       __recoverVal = nil
       return cp;
  end

  local panic = function(err)
     -- must wrap err in table to prevent conversion to string by error()
     __recoverVal = {err}
     setmetatable(__recoverVal, __recovMT)
     error(__recoverVal)
  end

  -- end deferBoilerplate1
`

const custom2 = `
  -- named returns available by closure
  local ret0=0
  local ret1=0

  local __actual=function()

     local __defer_func = function(a)
        -- capture any arguments at defer call point
        local a = a
        __defers[1+#__defers] = function()
           print("first defer running, a=", a, " b=",b, " ret0=", ret0, " ret1=", ret1)
            b = b + 3
            ret0 = (ret0+1) * 3
            ret1 = ret1 + 1
        end
     end
     __defer_func(a)

     local __defer_func = function()
        __defers[1+#__defers] = function()
           print("second defer running, a=", a, " b=",b, " ret1=", ret1)
           b = b * 7
           ret0 = ret0 + 100
           recov = recover()
           if type(recov[1]) == "number" then
              panic(recov[1] + 17)
           end
           ret1 = ret1 + 100
           print("second defer just updated ret1 to ", ret1)
        end
     end
     __defer_func()

     a = 1
     b = 1

    panic(a+b)

     return b, 58

  end -- end of __actual
`
const deferBoilerplate3 = `
  -- begin boilerplate3:

  -- prepare to handle panic/defer/recover
  local __handler2 = function(err)
     print("__handler2 running with err=",tostring(err)," and err[1] =", err[1])
     __recoverVal = err
     return err
  end

  local __panicHandler = function(err)
     __recoverVal = err
     print("panicHandler running with err =", tostring(err), " and #defer = ", #__defers)
     for __i = #__defers, 1, -1 do
        local dcall = {xpcall(__defers[__i], __handler2)}
        for i,v in pairs(dcall) do print("panic path defer call result: i=",i, "  v=",v) end
     end
     if __recoverVal ~= nil then
        return __recoverVal
     end
  end

  -- all set. make the actual call.
  local __res = {xpcall(__actual, __panicHandler, unpack(orig))}

  if __res[1] then
     print("call had no panic")
     -- call had no panic. run defers with the nil recover

     if #__res > 1 then
        for k,v in pairs(__res) do print("__res k=", k, " val=", v) end
        -- explicit returns fill the named vals before defers see them.
        print("pre fill: ret0 = ", ret0, " and ret1=", ret1)
`

const custom4 = `
        ret0, ret1 = table.unpack(__res, 2)
`

const custom4template = `
        %v = table.unpack(__res, 2)
`

const deferBoilerplate5 = `
     end

      assert(recoverVal == nil)
      for __i = #__defers, 1, -1 do
        local dcall = {xpcall(__defers[__i], __handler2)}
        for i,v in pairs(dcall) do print("normal path defer call result: i=",i, "  v=",v) end
      end
  else
     --print("checking for panic still un-caught...", __recoverVal)
     -- is there an un-recovered panic that we need to rethrow?
     if __recoverVal ~= nil then
        --print("un recovered error still exists, rethrowing ", __recoverVal)
        error(__recoverVal)
     end
  end

  -- end deferBoilerplate5
`

const custom6 = `
  -- custom section number 6, the returns:
  return ret0, ret1
end
`

const custom6template = `
  -- custom section number 6, the returns:
  return %v
end
`
