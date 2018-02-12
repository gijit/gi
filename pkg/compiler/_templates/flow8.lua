dofile 'deferinit.lua'


a = 0;
b = 0;

-- nested calls: can panic transfer
-- from lower in the stack to higher up?

global = 0


function deeper(...)
   orig = {...}
   
   local __defers={}
   local __zeroret = {}
   local __namedNames = {}

  local __actual=function(a)
     global = global + (a + 3)
     panic("panic-in-deeper")
  end -- end of __actual
  return __actuallyCall("deeper", __actual, __namedNames, __zeroret, __defers)
end


function intermed(a)
   deeper(a)
end



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
            ret0 = (ret0+1) * 3 + global
            ret1 = ret1 + 1
            recov = recover()
            print("defer 1 recovered ", recov)
            if recov ~= nil then
               ret1 = ret1 + 9 + global
               ret0 = ret0 + 19 + global
            end
        end
     end
     __defer_func(a)

     deeper(a)
     return

   end -- end of __actual
   return __actuallyCall("f", __actual, __namedNames, __zeroret, __defers)
end

f1, f2 = f()
print("f1 = ",f1, " f2=", f2)

--[[
dofile 'flow8.lua'
deeper: panicHandler running with err =	flow8.lua:30: panic-in-deeper	 and #defer = 	0
deeper: checking for panic still un-caught...	flow8.lua:30: panic-in-deeper
deeper: un recovered error still exists, rethrowing 	flow8.lua:30: panic-in-deeper
panicHandler running with err =	flow8.lua:89: flow8.lua:30: panic-in-deeper	 and #defer = 	1
first defer running, a=	0	 b=	0	 ret0=	0	 ret1=	0
defer 1 recovered 	flow8.lua:89: flow8.lua:30: panic-in-deeper
panic path defer call result: i=	1	  v=	true
checking for panic still un-caught...	nil
at f return point, ret0 = 	28	 ret1=	13
f1 = 	28	 f2=	13

--]]
