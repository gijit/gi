
dofile 'deferinit.lua'

a = 0;
b = 0;

-- f uses defer, boilerplate:  
function f(...)
   orig = {...}
   
   local __defers={}
   local __zeroret = {0,0}
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
        end
     end
     __defer_func(a)

     panic("in-between-defers-panic")
     
     local __defer_func = function()
        __defers[1+#__defers] = function()
           print("second defer running, a=", a, " b=",b, " ret1=", ret1)
           b = b * 7
           ret0 = ret0 + 100
           ret1 = ret1 + 100
           print("second defer just updated ret1 to ", ret1)
           recov = recover()
           print("second defer, recov is ", recov)
           panic("panic-in-defer-2")

           -- sadly, a raw error will result in loss of the "in-defer-2" value
           -- because of problems in luajit with recursive handling of
           -- xpcalls (it doesn't like them). So try to explicitly panic
           -- whenever possible instead of allowing an error to occur!
           -- https://stackoverflow.com/questions/48202338/on-latest-luajit-2-1-0-beta3-is-recursive-xpcall-possible
           --
           --error("error-in-defer-2") 
        end
     end
     __defer_func()

     a = 1
     b = 1

     return b, 58

  end -- end of __actual
  return __actuallyCall("f", __actual, __namedNames, __zeroret, __defers)
end

f1, f2 = f()
print("f1 = ",f1, " f2=", f2)

--[[

dofile 'flow5.lua'
panicHandler running with err =	flow5.lua:23: in-between-defers-panic	 and #defer = 	1
first defer running, a=	0	 b=	0	 ret0=	0	 ret1=	0
defer 1 recovered 	flow5.lua:23: in-between-defers-panic
panic path defer call result: i=	1	  v=	true
checking for panic still un-caught...	nil
at f return point, ret0 = 	3	 ret1=	1
f1 = 	3	 f2=	1
> 
--]]

