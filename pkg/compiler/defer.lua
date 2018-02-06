-- deferinit.lua : global setup for defer handling

-- utility: table show
function ts(t)
   if t == nil then
      return "<nil>"
   end
   local s = "<non-nil table:>\n"
   local k = 0
   for i,v in pairs(t) do
      s = s .. "key:" .. tostring(i) .. " -> val:" .. tostring(v) .. "\n"
      k = k +1
   end
   if k > 0 then
      return s
   end
   return "<non-nil but empty table with 0 entries>: " .. tostring(t)
end

-- can we have one global definition of panic and recover?
-- This would be preferred to repeating them in every function.

-- string viewing of panic value
__recovMT = {__tostring = function(v) return 'a-panic-value:' .. tostring(v[1]) end}

-- __recoverVal will be nill if no panic,
--              or if panic happened and
--              then was recovered.
-- this is always a table to avoid
--  stringification problems. The real
--  panic value is inside at position [1].
--
-- NB __recoverVal  needs to be per-goroutine. As each
--  could be unwinding independently at any
--  point in time.

__recoverVal = nil

recover = function() 
    local cp = __recoverVal
    __recoverVal = nil
    return cp;
end

panic = function(err)
  -- wrap err in table to prevent conversion to string by error()
  __recoverVal = {err}
  -- but still allow it to be viewable in a stack trace:
  setmetatable(__recoverVal, __recovMT)
  error(__recoverVal)
end


  -- begin boilerplate part 2:
  
  -- prepare to handle panic/defer/recover
__handler2 = function(err)
     --print(" __handler2 running with err =", err)
     __recoverVal = err
     return err
end
  
__panicHandler = function(err, defers)
       --print("__panicHandler running with err =", err)
       -- print(debug.traceback())
       --print("__panicHandler running with defers:", tostring(defers))

     __recoverVal = err
     if defers ~= nil then

         --print(debug.traceback(), " __panicHandler running with err =", err, " and #defer = ", #defers)      
         --print(" __panicHandler running with err =", err, " and #defer = ", #defers)  
         for __i = #defers, 1, -1 do
             local dcall = {xpcall(defers[__i], __handler2)}
             --for i,v in pairs(dcall) do print("__panicHandler: panic path defer call result: i=",i, "  v=",v) end
         end
     else
         --print("debug: found no defers in __panicHandler")
     end
     --print("__panicHandler: done with defer processing")
     if __recoverVal ~= nil then
        return __recoverVal
     end
  end

  -- __processDefers represents the normal
  --    return path, without a panic.
  --
  --    We need to update the named return values if
  --    there were explicit return values from __actual,
  --    and then we need to call the defers.
  --
  --    __namedNames is an array of the variable names of the return values,
  --                 so we know how to update actEnv.
  --
__processDefers = function(who, defers, __res, __namedNames, actEnv)
  --print(who,": __processDefers top: __res[1] is: ", tostring(__res[1]))
  --print(who,": __processDefers top: __namedNames is: ", ts(__namedNames))

  if __res[1] then
      --print(who,": __processDefers: call had no panic")
      -- call had no panic. run defers with the nil recover

      if #__res > 1 then
         --for k,v in pairs(__res) do print(who, " __processDefers: __res k=", k, " val=", v) end

         -- explicit return, so fill the named vals before defers see them.
         local unp = {table.unpack(__res, 2)}
         --print("unp is: ", tostring(unp))
         for i, k in pairs(__namedNames) do
             actEnv[k] = unp[i]
         end

         --print(who, " __processDefers: post fill: ret0 = ", ret0, " and ret1=", ret1)
      end

      assert(recoverVal == nil)
      for __i = #defers, 1, -1 do
        local dcall = {xpcall(defers[__i], __handler2)}
        for i,v in pairs(dcall) do
            --print(who," __processDefers: normal path defer call result: i=",i, "  v=",v)
        end
      end
  else
      --print(who, " __processDefers: checking for panic still un-caught...", __recoverVal)
      -- is there an un-recovered panic that we need to rethrow?
      if __recoverVal ~= nil then
         --print(who, "__processDefers: un recovered error still exists, rethrowing ", __recoverVal)
         error(__recoverVal)
      end
  end

  if #__namedNames == 0 then
     --print("__processDefers: #__namedNames was 0, no returns")
     return nil
  end
  -- put the named return values in order
  local orderedReturns={}
  for i, k in pairs(__namedNames) do
     --print("debug: fetching from function env k=",k," which we see has value ", actEnv[k], "in actEnv", tostring(actEnv))
     orderedReturns[i] = actEnv[k]
  end
  for i,v in pairs(orderedReturns) do
       --print(who," __processDefers: orderedReturns: i=",i, "  v=",v)
  end 
  return unpack(orderedReturns)
end


__actuallyCall = function(who, __actual, __namedNames, __zeroret, __defers, __orig)

   --local actEnv = getfenv(__actual)
   -- So getfenv(__actual) showed that actEnv
   -- was the _G global env, not good.
   -- To fix this, we give f its own env,
   -- so that named return variables can
   -- be written/read from this env.
   
   local actEnv = {}
   local mt = {
      __index = _G, -- read through to globals.
      __newindex = _G, -- write to closure-capture globals too.
   }
   setmetatable(actEnv,mt)
   setfenv(__actual, actEnv)

  for i,k in pairs(__namedNames) do
     --print("filling actEnv[k='"..tostring(k).."'] = '"..tostring(actEnv[k]).."' with __zeroret[i='"..tostring(i).."']='",tostring(__zeroret[i]),"'")
     actEnv[k] = __zeroret[i]
  end  
  local myPanic = function(err) __panicHandler(err, __defers) end
  local __res = {xpcall(__actual, myPanic, unpack(__orig))}
  return __processDefers(who, __defers, __res,  __namedNames, actEnv)  
end
