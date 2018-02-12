a = 0;
b = 0;
function f(...)
   orig = {...}
   
   -- f uses defer, boilerplate:
   
   local __defers={}
   
   -- __recoverVal will be nill if no panic,
   --              or if panic happened and
   --              then was recovered.   
   local __recoverVal
   
  local recover = function() 
       local cp = __recoverVal
       __recoverVal = nil
       return cp;
  end

  local panic = function(err)
     __recoverVal = err
     error(err)
  end

  
  -- end boilerplate, begin custom:
  
  -- named returns available by closure
  local ret0=0
  local ret1=0

  local __actual=function()

     local __defer_func = function(a)
        -- capture any arguments at defer call point
        local a = a
        __defers[1+#__defers] = function() 
            print("first defer running, a=", a, " b=",b)
            b = b + 3
            ret0 = (ret0+1) * 3
            ret1 = ret1 + 1
            recov = recover()
            print("defer 1 recovered ", recov)
        end
     end
     __defer_func(a)
     
     local __defer_func = function()
        __defers[1+#__defers] = function()
           print("second defer running, a=", a, " b=",b)
           b = b * 7
           ret0 = ret0 + 100
           ret1 = ret0 + 100
           --recov = recover()
           --print("second defer, recov is ", recov)
           panic("panic-in-defer-2")
        end
     end
     __defer_func()

     a = 1
     b = 1

    panic("ouch")
     
     return b, 58

  end -- end of __actual

  -- begin boilerplate part 2:
  
  -- prepare to handle panic/defer/recover
  local __handler2 = function(err)
     print("__handler2 running with err =", err)
     __recoverVal = err
     return err
  end
  
  local __panicHandler = function(err)
     __recoverVal = err
     print("panicHandler running with err =", err, " and #defer = ", #__defers)
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
        ret0, ret1 = table.unpack(__res, 2)
        print("post fill: ret0 = ", ret0, " and ret1=", ret1)
     end

      assert(recoverVal == nil)
      for __i = #__defers, 1, -1 do
        local dcall = {xpcall(__defers[__i], __handler2)}
        for i,v in pairs(dcall) do print("normal path defer call result: i=",i, "  v=",v) end
      end
  else
     print("checking for panic still un-caught...", __recoverVal)
     -- is there an un-recovered panic that we need to rethrow?
     if __recoverVal ~= nil then
        print("un recovered error still exists, rethrowing ", __recoverVal)
        error(__recoverVal)
     end
  end

  -- custom section number 2, the returns:
  print("at f return point, ret0 = ", ret0, " ret1=", ret1)
  return ret0, ret1
end
f()
