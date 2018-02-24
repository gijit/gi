-- goroutines.lua

local ffi = require("ffi")

__stackDepthOffset = 0;
__getStackDepth = function() 
   -- javascript/mozilla gives a stack trace under .stack,
   -- which GopherJS was using if available.
   -- Inlining makes debug.traceback() pretty useless under LuaJIT.
   
   return __stackDepthOffset
end;

__panicStackDepth = nil;
__panicValue = nil;

__callDeferred = function(deferred, jsErr, fromPanic) 
   if not fromPanic and deferred ~= nil and deferred.index >= #__curGoroutine.deferStack then
      error( jsErr);
   end
   if jsErr ~= nil then
      local newErr = nil;
      --try
      local res = {pcall(function()
                         __curGoroutine.deferStack.push(deferred);
                         __panic(__jsErrorPtr(jsErr));
      end)}
      local ok, err = unpack(res)
      --catch (err)
      if not ok then
         newErr = err;
      end
      __curGoroutine.deferStack.pop();
      __callDeferred(deferred, newErr);
      return;
   end
   if __curGoroutine.asleep then
      return;
   end

   __stackDepthOffset=__stackDepthOffset-1;
   outerPanicStackDepth = __panicStackDepth;
   outerPanicValue = __panicValue;

   local localPanicValue = __curGoroutine.panicStack.pop();
   if localPanicValue ~= nil then
      __panicStackDepth = __getStackDepth();
      __panicValue = localPanicValue;
   end

   --try
   local res = {pcall(function()
                      ::top::                 
                      while true do
                         if deferred == nil then
                            deferred = __curGoroutine.deferStack[#__curGoroutine.deferStack - 1];
                            if deferred == nil then
                               -- The panic reached the top of the stack. Clear it and throw it as a Lua error. --
                               __panicStackDepth = nil;
                               error(localPanicValue)
                            end
                         end
                         local call = deferred.pop();
                         if call == nil then
                            __curGoroutine.deferStack.pop();
                            if localPanicValue ~= nil then
                               deferred = nil;
                               goto top; -- continue;
                            end
                            return;
                         end
                         local r = call[0](call[2], call[1]);
                         if r and r.__blk ~= nil then
                            deferred.push({r.__blk, {}, r});
                            if fromPanic then
                               error( nil);
                            end
                            return;
                         end

                         if localPanicValue ~= nil and __panicStackDepth == nil then
                            error( nil); -- error was recovered --
                         end
                      end
   end)}
   --finally, no catch
   
   if localPanicValue ~= nil then
      if __panicStackDepth ~= nil then
         __curGoroutine.panicStack.push(localPanicValue);
      end
      __panicStackDepth = outerPanicStackDepth;
      __panicValue = outerPanicValue;
   end
   __stackDepthOffset=__stackDepthOffset+1;

   -- end finally, no catch
   -- need to rethrow?
   local ok, err = unpack(res)
   if not ok then
      -- rethrow
      error(err)
   end
end;

__panic = function(value) 
   __curGoroutine.panicStack.push(value);
   __callDeferred(nil, nil, true);
end;

__recover = function() 
   if __panicStackDepth == nil or (__panicStackDepth ~= nil and __panicStackDepth ~= __getStackDepth() - 2) then
      return __ifaceNil;
   end
   __panicStackDepth = nil;
   return __panicValue;
end;

__throw = function(err)  error(err); end;

__noGoroutine = { asleep= false, exit= false, deferStack= {}, panicStack= {} };

__curGoroutine = __noGoroutine;
__totalGoroutines = 0;
__awakeGoroutines = 0;
__checkForDeadlock = true;
__mainFinished = false;

__go = function(fun, args) 
   __totalGoroutines=__totalGoroutines+1;
   __awakeGoroutines=__awakeGoroutines+1;
   local __goroutine = function() 
      --try
      local res = {pcall(function()
                         
                         __curGoroutine = __goroutine;
                         local r = fun(nil, args);
                         if r and r.__blk ~= nil then
                            fun = function()  return r.__blk(); end;
                            args = {};
                            return;
                         end
                         __goroutine.exit = true;
      end)}
      -- finally

      __curGoroutine = __noGoroutine;
      if __goroutine.exit then -- also set by runtime.Goexit() --
         __totalGoroutines=__totalGoroutines-1;
         __goroutine.asleep = true;
      end
      if __goroutine.asleep then
         __awakeGoroutines=__awakeGoroutines-1;
         if not __mainFinished and __awakeGoroutines == 0 and __checkForDeadlock then
            console.error("fatal error: all goroutines are asleep - deadlock!");
            if __global.process ~= nil then
               __global.process.exit(2);
            end
         end
      end
      -- end finally
      -- catch (err)
      local ok, err = unpack(res)
      if not ok then
         if not __goroutine.exit then
            -- rethrow
            error(err);
         end
      end
      
   end;
   __goroutine.asleep = false;
   __goroutine.exit = false;
   __goroutine.deferStack = {};
   __goroutine.panicStack = {};
   __schedule(__goroutine);
end;

__scheduled = {};

__runScheduled = function() 
   --try
   local res = {pcall(function()
                      while true do
                         local r = __scheduled.shift();
                         if r == nil then
                            break
                         end
                         r()
                      end
   end)}
   -- finally, no catch
   if #__scheduled > 0 then
      setTimeout(__runScheduled, 0);
   end
   -- end finally, no catch
   -- need to rethrow?
   local ok, err = unpack(res)
   if not ok then
      error(err)
   end
end;

__schedule = function(goroutine) 
   if goroutine.asleep then
      goroutine.asleep = false;
      __awakeGoroutines=__awakeGoroutines+1;
   end
   __scheduled.push(goroutine);
   if __curGoroutine == __noGoroutine then
      __runScheduled();
   end
end;

__setTimeout = function(f, t) 
   __awakeGoroutines=__awakeGoroutines+1;
   return setTimeout(function() 
         __awakeGoroutines=__awakeGoroutines-1;
         f();
                     end, t);
end;

__block = function() 
   if __curGoroutine == __noGoroutine then
      __throwRuntimeError("cannot block in JavaScript callback, fix by wrapping code in goroutine");
   end
   __curGoroutine.asleep = true;
end;

__send_GopherJS = function(chan, value) 
   if chan.__closed then
      __throwRuntimeError("send on closed channel");
   end
   local queuedRecv = chan.__recvQueue.shift();
   if queuedRecv ~= nil then
      queuedRecv({value, true});
      return;
   end
   if #chan.__buffer < chan.__capacity then
      chan.__buffer.push(value);
      return;
   end

   thisGoroutine = __curGoroutine;
   local closedDuringSend;
   chan.__sendQueue.push(function(closed) 
         closedDuringSend = closed;
         __schedule(thisGoroutine);
         return value;
   end);
   __block();
   return {
      __blk= function() 
         if closedDuringSend then
            __throwRuntimeError("send on closed channel");
         end
      end
   };
end;

__recv = function(chan)
   --print("__recv called!")

   local ch = reflect.ValueOf(chan)
   local rv, ok = ch.Recv();
   -- rv is userdata, a reflect.Value. Convert to
   -- interface{} for Luar, using Interface(), so
   -- luar can translate that to Lua for us.
   local v = rv.Interface();
   return {v, ok}
end

__send = function(chan, value)
   print("__send called! value=", value)
   local ch = reflect.ValueOf(chan)
   local v = reflect.ValueOf(value)
   local cv = v.Convert(reflect.TypeOf(chan).Elem())
   ch.Send(cv);
end

-- TODO: below is hardcoded for select with two receives.
--       We need to read comms and allocate cases based on
--       the content of comms.
--
__select = function(comms)
   print("__select called!")
   __st(comms, "comms")
   
   __st(comms[1], "comms[1]")
   __st(comms[2], "comms[2]")
   
   __st(comms[1][1], "comms[1][1]")
   __st(comms[2][1], "comms[2][1]")

   local c1 = reflect.ValueOf(comms[1][1])
   local c2 = reflect.ValueOf(comms[2][1])

   print("c1 is "..type(c1))
   print("c2 is "..type(c2))
      
   print("c1 = "..tostring(c1))
   print("c2 = "..tostring(c2))

 --[[
    for i, comm in ipairs(comms) do
      local chan = comm[1];
      --switch (comm.length)
      local comm_len = #comm
      if comm_len == 0 then
         -- default --
         selection = i-1;
         break;
      elseif comm_len == 1 then
         -- recv --

      elseif comm_len == 2 then
         -- send --

      end -- end switch
   end
 --]]
   
   __refSelCaseRecvVal0.Chan = c1
   __refSelCaseRecvVal1.Chan = c2

   local cases = {
      __refSelCaseRecvVal0,
      __refSelCaseRecvVal1,
   }
   
   local chosen, recv, recvOk = reflect.Select(cases)
   print("back from reflect.Select, we got: chosen=", chosen)
   print("back from reflect.Select, we got:   recv=", recv.Interface())
   print("back from reflect.Select, we got: recvOk=", recvOk)
   
   return {chosen, {recv.Interface(), recvOk}};
end

-- gopherJs port: __recv
__recv__GopherJS = function(chan) 
   local queuedSend = chan.__sendQueue.shift();
   if queuedSend ~= nil then
      chan.__buffer.push(queuedSend(false));
   end
   local bufferedValue = chan.__buffer.shift();
   if bufferedValue ~= nil then
      return bufferedValue, true;
   end
   if chan.__closed then
      return chan.__elem.zero(), false;
   end

   local thisGoroutine = __curGoroutine;
   local f = { __blk= function(self)  return self.value; end };
   local queueEntry = function(v) 
      f.value = v;
      __schedule(thisGoroutine);
   end;
   chan.__recvQueue.push(queueEntry);
   __block();
   return f;
end;

__close = function(chan) 
   if chan.__closed then
      __throwRuntimeError("close of closed channel");
   end
   chan.__closed = true;
   while true do
      local queuedSend = chan.__sendQueue.shift();
      if queuedSend == nil then
         break;
      end
      queuedSend(true); -- will panic --
   end
   while true do
      local queuedRecv = chan.__recvQueue.shift();
      if queuedRecv == nil then
         break;
      end
      queuedRecv({chan.__elem.zero(), false});
   end
end;

__select_GopherJS = function(comms)
   local ready = {};
   local selection = -1;
   for i, comm in ipairs(comms) do
      local chan = comm[1];
      --switch (comm.length)
      local comm_len = #comm
      if comm_len == 0 then
         -- default --
         selection = i-1;
         break;
      elseif comm_len == 1 then
         -- recv --
         if #chan.__sendQueue ~= 0 or #chan.__buffer ~= 0 or chan.__closed then
            ready.push(i-1);
         end
         break;
      elseif comm_len == 2 then
         -- send --
         if chan.__closed then
            __throwRuntimeError("send on closed channel");
         end
         if #chan.__recvQueue ~= 0 or #chan.__buffer < chan.__capacity then
            ready.push(i-1);
         end
         break;
      end -- end switch
   end

   if #ready  ~= 0 then
      -- jea NB lua's math.random(n) returns in [1,n] to match its arrays.
      selection = ready[math.random(#ready)]; 
   end
   if selection ~= -1 then
      local comm = comms[selection];
      -- switch (comm.length)
      local comm_len = #comm
      if comm_len == 0 then
         -- default --
         return {selection};
      elseif comm_len == 1 then
         -- recv --
         return {selection, __recv(comm[1])};
      elseif comm_len == 2 then
         -- send --
         __send(comm[1], comm[2]);
         return {selection};
      end
   end

   local entries = {};
   local thisGoroutine = __curGoroutine;
   local f = { __blk= function(self)  return self.selection; end };
   local removeFromQueues = function() 
      for i,entry in ipairs(entries) do
         local queue = entry[1];
         local index = queue.indexOf(entry[2]);
         if index ~= -1 then
            queue.splice(index, 2);
         end
      end
   end;
   for i,comm in ipairs(comms) do

      --switch on #comm
      local comm_len = #comm
      
      if comm_len == 1 then
         -- recv --
         local queueEntry = function(value) 
            f.selection = {i, value};
            removeFromQueues();
            __schedule(thisGoroutine);
         end;
         entries.push({comm[1].__recvQueue, queueEntry});
         comm[1].__recvQueue.push(queueEntry);

      elseif comm_len == 2 then
         -- send --
         local queueEntry = function() 
            if comm[1].__closed then
               __throwRuntimeError("send on closed channel");
            end
            f.selection = {i};
            removeFromQueues();
            __schedule(thisGoroutine);
            return comm[1];
         end;
         entries.push({comm[1].__sendQueue, queueEntry});
         comm[1].__sendQueue.push(queueEntry);

      end
   end
   __block();
   return f;
end;

