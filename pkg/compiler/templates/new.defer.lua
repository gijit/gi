a = 0;
b = 0;
function f(...)
   orig = {...}
   
   -- f uses defer, boilerplate:
   
   local __defers={}
   
   -- __recoverVal will be nil if no panic,
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
            --recov = recover()
            --print("first defer, recov is ", recov)
        end
     end
     __defer_func(a)
     
     local __defer_func = function()
        __defers[1+#__defers] = function()
           print("second defer running, a=", a, " b=",b)
           b = b * 7
           ret0 = ret0 + 100
           ret1 = ret0 + 100
           recov = recover()
           print("second defer, recov is ", recov)
           do return panic("new value") end
        end
     end
     __defer_func()

     a = 1
     b = 1

     do return panic("ouch") end
     
     return b, 58

  end -- end of __actual

  -- begin boilerplate part 2:

  local __panicStart = function(err) print("panic starting, err=",err) end
  
  -- all set. make the actual call.
  local __res = {xpcall(__actual, __panicStart, unpack(orig))}

  if not __res[1] then
     print("res[1] was false, panic must have happenned. __recoverVal =", __recoverVal)
  end
  
  if __res ~= nil and #__res > 0 then
     
     for k,v in pairs(__res) do print("__res k=", k, " val=", v) end

     -- explicit returns:  fill the named vals before defers see them.
     print("pre fill: ret0 = ", ret0, " and ret1=", ret1)

     ret0, ret1 = table.unpack(__res, 2)

     print("post fill: ret0 = ", ret0, " and ret1=", ret1)
  end

  for __i = #__defers, 1, -1 do
     local __defercall = {xpcall(__defers[__i], __panicStart)}
  end

  print("checking for panic still un-caught...", __recoverVal)
  -- is there an un-recovered panic that we need to rethrow?
  if __recoverVal ~= nil then
     print("un recovered error still exists, rethrowing ", __recoverVal)
     error(__recoverVal)
  end

  -- custom section number 2, the returns:
  print("at f return point, ret0 = ", ret0, " ret1=", ret1)
  return ret0, ret1
end
a1, a2 = f()
print("back from calling f(), a1=", a1, " a2=", a2)
