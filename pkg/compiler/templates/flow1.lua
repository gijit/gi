
dofile 'deferinit.lua'

a = 7;
b = 0;

-- f uses defer, boilerplate:  
function f(...)
   orig = {...}
   
   local __defers={}
   
  
  -- end boilerplate, begin custom:

   -- The ordering of entries in the __zeroret and __namedNames arrays
   -- must match the declaration order of the return parameters.
   local __zeroret = {0,0}
  -- even anonymous returns get names here.
  local __namedNames = {"ret0", "ret1"}

     
  local __actual=function()

     local __defer_func = function(a)
        -- capture any arguments at defer call point
        local a = a
        __defers[1+#__defers] = function()
           print("first defer running, a=", a, " b=",b, " ret0=",ret0," ret1=", ret1)           
            ret0 = (ret0+1) * 3 + a
            ret1 = ret1 + 1 + a
        end
     end
     __defer_func(a)
     
     local __defer_func = function()
        __defers[1+#__defers] = function()
           print("second defer running, a=", a, " b=",b, " ret1=", ret1)
           ret0 = ret0 + 100
           ret1 = ret1 + 100
        end
     end
     __defer_func()

     a = 1
     b = 1
     
     return b, a + 58

  end -- end of __actual
  return __actuallyCall("f", __actual, __namedNames, __zeroret, __defers)

end

f1, f2 = f()
print("f1 = ",f1, " f2=", f2)

--[[
dofile 'flow1.lua'
call had no panic
__res k=	1	 val=	true
__res k=	2	 val=	1
__res k=	3	 val=	59
pre fill: ret0 = 	0	 and ret1=	0
post fill: ret0 = 	1	 and ret1=	59
second defer running, a=	1	 b=	1	 ret1=	59
normal path defer call result: i=	1	  v=	true
first defer running, a=	7	 b=	1	 ret0=	101	 ret1=	159
normal path defer call result: i=	1	  v=	true
at f return point, ret0 = 	313	 ret1=	167
f1 = 	313	 f2=	167
> 
--]]

