-- zgoro.lua, named with a 'z' in order
-- to load last in prelude,
-- after tsys.lua, so we have our types.

-- For the hybrid/interacts with
-- native Go channels via reflect
-- version of goroutines and
-- channels: see reflect_goro.lua

-- turn off zgoro for now, using chan.lua presently.
-- just stub out __go as passthrough
function __go(f, ...)
   f(...)
end

--[==[

local ffi = require("ffi")


__noGoroutine = { asleep= false, exit= false, deferStack= {}, panicStack= {} };

__curGoroutine = __noGoroutine;
__totalGoroutines = 0;
__awakeGoroutines = 0;
__checkForDeadlock = true;
__mainFinished = false;

__go = function(fun, args) 
   __totalGoroutines=__totalGoroutines+1;
   __awakeGoroutines=__awakeGoroutines+1;
   local __goroutine = {
      __call = function() 
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
   }
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

-- gopherJs port: __recv
__recv = function(chan) 
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

__send = function(chan, value) 
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

__select = function(comms)
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
      -- jea NB lua's __builtin_math.random(n) returns in [1,n] to match its arrays.
      selection = ready[__builtin_math.random(#ready)]; 
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

--]==]
