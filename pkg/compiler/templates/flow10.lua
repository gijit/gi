
-- flow10.lua

-- same go code as flow8.go, but here we factor out the lua funtions
--  into re-usable functions as opposed to being generated for
--  each function.

dofile 'deferinit.lua'

--

a = 0;
b = 0;

-- nested calls: can panic transfer
-- from lower in the stack to higher up?

global = 8


function deeper(...)
   orig = {...}
   
  -- 'deeper' uses defer or panic, boilerplate:
  local __defers={}  

  -- end boilerplate, begin custom:
  
  -- order here  must match declaration order.
  local __zeroret = {}
  
  -- order here must match declaration order.
  local __namedNames = {}

  local __actual=function(a)
       print("in deeper actual, global=", global)
       global = global + (a + 3)
       print("in deeper actual, after setting it, global=", global)     
       panic("panic-in-deeper")
       
  end
  return __actuallyCall("deeper", __actual, __namedNames, __zeroret, __defers)
end

function intermed(a)
   deeper(a)
end


function f(...)
   orig = {...}
   
   -- f uses defer or panic, boilerplate:
   local __defers={}
    
  -- end boilerplate, begin custom:

   -- named returns available by closure.
   
   -- The ordering of entries in the __zeroret and __namedNames arrays
   -- must match the declaration order of the return parameters.
   local __zeroret = {0, 0}
   -- even anonymous returns get names here.
   local __namedNames = {"ret0", "ret1"}

   local __actual=function()
  
     local __defer_func = function(a)
        -- capture any arguments at defer call point
        local a = a
        __defers[1+#__defers] = function() 
            print("first defer running, a=", a, " b=",b, " ret0=", ret0, " ret1=", ret1, " global=",global)
            b = b + 3
            ret0 = (ret0+1) * 3 + global
            print("debug, after ret0 = (ret0+1) * 3 + global: ret0=", ret0)
            ret1 = ret1 + 1
            recov = recover()
            print("\ndefer 1 recovered ", recov)
            if recov ~= nil then
               print("recov was not nil... ret0 is now=", ret0, " and ret1 is now=", ret1)
               ret1 = ret1 + 9 + global
               ret0 = ret0 + 19 + global
            end
            print("end of first defer, ret0=", ret0, " ret1=", ret1)
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

> dofile 'flow10.lua'
dofile 'flow10.lua'
in deeper actual, global=	8
in deeper actual, after setting it, global=	11
__panicHandler running with err =	a-panic-value:panic-in-deeper
__panicHandler running with defers:	table: 0x0004e970
__panicHandler: done with defer processing
deeper	: __processDefers top: __res[1] is: 	false
deeper	: __processDefers top: __namedNames is: 	<non-nil but empty table with 0 entries>: table: 0x0004e9c0
deeper	 __processDefers: checking for panic still un-caught...	a-panic-value:panic-in-deeper
deeper	__processDefers: un recovered error still exists, rethrowing 	a-panic-value:panic-in-deeper
__panicHandler running with err =	a-panic-value:panic-in-deeper
__panicHandler running with defers:	table: 0x0004f148
first defer running, a=	0	 b=	0	 ret0=	0	 ret1=	0	 global=	11
debug, after ret0 = (ret0+1) * 3 + global: ret0=	14

defer 1 recovered 	a-panic-value:panic-in-deeper
recov was not nil... ret0 is now=	14	 and ret1 is now=	1
end of first defer, ret0=	44	 ret1=	21
__panicHandler: panic path defer call result: i=	1	  v=	true
__panicHandler: done with defer processing
f	: __processDefers top: __res[1] is: 	false
f	: __processDefers top: __namedNames is: 	<non-nil table:>
key:1 -> val:ret0
key:2 -> val:ret1

f	 __processDefers: checking for panic still un-caught...	nil
f	 __processDefers: orderedReturns: i=	1	  v=	44
f	 __processDefers: orderedReturns: i=	2	  v=	21
f1 = 	44	 f2=	21
> 

--]]
