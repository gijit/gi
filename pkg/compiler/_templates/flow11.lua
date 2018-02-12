
-- flow11.lua

dofile("deferinit.lua")

---

a = 0;
b = 0;

-- nested calls: can panic transfer
-- from lower in the stack to higher up?


function f(...)
   orig = {...}
   
   -- f uses defer or panic, boilerplate:
   local __defers={}
    
  -- end boilerplate, begin custom:

  -- order here  must match declaration order.
  local __zeroret = {0} -- 0 is the zero value for int, the one return type of f.
   
  -- even anonymous returns get names here.
   __namedNames = {"__ret0"}
   
  local __actual=function()
  
     local __defer_func = function()
        -- capture any arguments at defer call point
        __defers[1+#__defers] = function()
           print("first defer running")
        end
     end
     __defer_func()

     local __defer_func = function()
        -- capture any arguments at defer call point
        __defers[1+#__defers] = function()
           print("second defer running")
           recov = recover()
        end
     end
     __defer_func()
     
     panic("ouch")
     return

  end -- end of __actual

  -- begin boilerplate part 2:

  local actEnv = getfenv(__actual)
  for i,k in pairs(__namedNames) do
     actEnv[k] = __zeroret[i]
  end  
  local myPanic = function(err) __panicHandler(err, __defers) end
  local __res = {xpcall(__actual, myPanic, unpack(orig))}
  return __processDefers("f", __defers, __res,  __namedNames, actEnv)  
end

f1 = f()
print("f1 = ",f1)

--[[

> dofile 'flow11.lua'
dofile 'flow11.lua'
f	: __processDefers top: __res[1] is: 	false
f	: __processDefers top: __namedNames is: 	<non-nil table:>
key:1 -> val:__ret0

f	 __processDefers: checking for panic still un-caught...	nil
__processDefers: generating 	1	default return values
i=	1	k=	__ret0	actEnv[k]=	0
orderedReturns has len	1
f	 __processDefers: orderedReturns: i=	1	  v=	0
f1 = 	0
> 


--]]
