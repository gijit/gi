
dofile 'deferinit.lua'

a = 0;
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
           print("first defer running, a=", a, " b=",b, " ret0=", ret0, " ret1=", ret1)
            b = b + 3
            ret0 = (ret0+1) * 3
            ret1 = ret1 + 1
            recov = recover()
            print("defer 1 recovered ", recov)
            
            if type(recov[1]) == "number" then
               r = recov[1]
               ret1 = ret1 + 9 + r*199
               ret0 = ret0 + 19 + r*3
            end
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
  return __actuallyCall("f", __actual, __namedNames, __zeroret, __defers)  
end

f1, f2 = f()
print("f1 = ",f1, " f2=", f2)

--[[

dofile 'flow2.lua'
panicHandler running with err[1] =	2	 and #defer = 	2
second defer running, a=	1	 b=	1	 ret1=	0
panic path defer call result: i=	1	  v=	false
panic path defer call result: i=	2	  v=	error in error handling
first defer running, a=	0	 b=	7	 ret0=	100	 ret1=	0
defer 1 recovered 	table: 0x0004e1b0
panic path defer call result: i=	1	  v=	true
checking for panic still un-caught...	nil
at f return point, ret0 = 	379	 ret1=	3791
f1 = 	379	 f2=	3791
> 

--]]

