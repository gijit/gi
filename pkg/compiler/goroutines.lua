-- goroutines.lua

__stackDepthOffset = 0;
__getStackDepth = function() 
  var err = new Error();
  if (err.stack == nil) then
    return nil;
  end
  return __stackDepthOffset + err.stack.split("\n").length;
end;

__panicStackDepth = null, __panicValue;
__callDeferred = function(deferred, jsErr, fromPanic) 
  if (!fromPanic && deferred ~= null && deferred.index >= __curGoroutine.deferStack.length) then
    throw jsErr;
  end
  if (jsErr ~= null) then
    var newErr = null;
    try {
      __curGoroutine.deferStack.push(deferred);
      __panic(new __jsErrorPtr(jsErr));
    } catch (err) {
      newErr = err;
    }
    __curGoroutine.deferStack.pop();
    __callDeferred(deferred, newErr);
    return;
  end
  if (__curGoroutine.asleep) then
    return;
  end

  __stackDepthOffset--;
  outerPanicStackDepth = __panicStackDepth;
  outerPanicValue = __panicValue;

  localPanicValue = __curGoroutine.panicStack.pop();
  if (localPanicValue ~= nil) then
    __panicStackDepth = __getStackDepth();
    __panicValue = localPanicValue;
  end

  try {
    while (true) do
      if (deferred == null) then
        deferred = __curGoroutine.deferStack[__curGoroutine.deferStack.length - 1];
        if (deferred == nil) then
          -- The panic reached the top of the stack. Clear it and throw it as a JavaScript error. --
          __panicStackDepth = null;
          if (localPanicValue.Object instanceof Error) then
            throw localPanicValue.Object;
          end
          var msg;
          if (localPanicValue.constructor == __String) then
            msg = localPanicValue.__val;
          elseif (localPanicValue.Error ~= nil) then
            msg = localPanicValue.Error();
          elseif (localPanicValue.String ~= nil) then
            msg = localPanicValue.String();
          else 
            msg = localPanicValue;
          end
          throw Error(msg); -- jea: was new Error
        end
      end
      var call = deferred.pop();
      if (call == nil) then
        __curGoroutine.deferStack.pop();
        if (localPanicValue ~= nil) then
          deferred = null;
          continue;
        end
        return;
      end
      var r = call[0].apply(call[2], call[1]);
      if (r && r.__blk ~= nil) then
        deferred.push([r.__blk, {}, r]);
        if (fromPanic) then
          throw null;
        end
        return;
      end

      if (localPanicValue ~= nil && __panicStackDepth == null) then
        throw null; -- error was recovered --
      end
    end
  } finally {
    if (localPanicValue ~= nil) then
      if (__panicStackDepth ~= null) then
        __curGoroutine.panicStack.push(localPanicValue);
      end
      __panicStackDepth = outerPanicStackDepth;
      __panicValue = outerPanicValue;
    end
    __stackDepthOffset++;
  }
end;

__panic = function(value) 
  __curGoroutine.panicStack.push(value);
  __callDeferred(null, null, true);
end;
__recover = function() 
  if (__panicStackDepth == null || (__panicStackDepth ~= nil && __panicStackDepth ~= __getStackDepth() - 2)) then
    return __ifaceNil;
  end
  __panicStackDepth = null;
  return __panicValue;
end;
__throw = function(err)  throw err; end;

__noGoroutine = { asleep: false, exit: false, deferStack: {}, panicStack: {} };
__curGoroutine = __noGoroutine, __totalGoroutines = 0, __awakeGoroutines = 0, __checkForDeadlock = true;
__mainFinished = false;
__go = function(fun, args) 
  __totalGoroutines++;
  __awakeGoroutines++;
  var __goroutine = function() 
    try {
      __curGoroutine = __goroutine;
      var r = fun.apply(nil, args);
      if (r && r.__blk ~= nil) then
        fun = function()  return r.__blk(); end;
        args = {};
        return;
      end
      __goroutine.exit = true;
    } catch (err) {
      if (!__goroutine.exit) then
        throw err;
      end
    } finally {
      __curGoroutine = __noGoroutine;
      if (__goroutine.exit) then -- also set by runtime.Goexit() --
        __totalGoroutines--;
        __goroutine.asleep = true;
      end
      if (__goroutine.asleep) then
        __awakeGoroutines--;
        if (!__mainFinished && __awakeGoroutines == 0 && __checkForDeadlock) then
          console.error("fatal error: all goroutines are asleep - deadlock!");
          if (__global.process ~= nil) then
            __global.process.exit(2);
          end
        end
      end
    }
  end;
  __goroutine.asleep = false;
  __goroutine.exit = false;
  __goroutine.deferStack = {};
  __goroutine.panicStack = {};
  __schedule(__goroutine);
end;

__scheduled = {};
__runScheduled = function() 
  try {
    var r;
    while ((r = __scheduled.shift()) ~= nil) do
      r();
    end
  } finally {
    if (__scheduled.length > 0) then
      setTimeout(__runScheduled, 0);
    end
  }
end;

__schedule = function(goroutine) 
  if (goroutine.asleep) then
    goroutine.asleep = false;
    __awakeGoroutines++;
  end
  __scheduled.push(goroutine);
  if (__curGoroutine == __noGoroutine) then
    __runScheduled();
  end
end;

__setTimeout = function(f, t) 
  __awakeGoroutines++;
  return setTimeout(function() 
    __awakeGoroutines--;
    f();
  end, t);
end;

__block = function() 
  if (__curGoroutine == __noGoroutine) then
    __throwRuntimeError("cannot block in JavaScript callback, fix by wrapping code in goroutine");
  end
  __curGoroutine.asleep = true;
end;

__send = function(chan, value) 
  if (chan.__closed) then
    __throwRuntimeError("send on closed channel");
  end
  var queuedRecv = chan.__recvQueue.shift();
  if (queuedRecv ~= nil) then
    queuedRecv([value, true]);
    return;
  end
  if (chan.__buffer.length < chan.__capacity) then
    chan.__buffer.push(value);
    return;
  end

  thisGoroutine = __curGoroutine;
  closedDuringSend;
  chan.__sendQueue.push(function(closed) 
    closedDuringSend = closed;
    __schedule(thisGoroutine);
    return value;
  end);
  __block();
  return {
    __blk: function() 
      if (closedDuringSend) then
        __throwRuntimeError("send on closed channel");
      end
    end
  };
end;
__recv = function(chan) 
  var queuedSend = chan.__sendQueue.shift();
  if (queuedSend ~= nil) then
    chan.__buffer.push(queuedSend(false));
  end
  var bufferedValue = chan.__buffer.shift();
  if (bufferedValue ~= nil) then
    return [bufferedValue, true];
  end
  if (chan.__closed) then
    return [chan.__elem.zero(), false];
  end

  var thisGoroutine = __curGoroutine;
  var f = { __blk: function()  return this.value; end };
  var queueEntry = function(v) 
    f.value = v;
    __schedule(thisGoroutine);
  end;
  chan.__recvQueue.push(queueEntry);
  __block();
  return f;
end;
__close = function(chan) 
  if (chan.__closed) then
    __throwRuntimeError("close of closed channel");
  end
  chan.__closed = true;
  while (true) do
    var queuedSend = chan.__sendQueue.shift();
    if (queuedSend == nil) then
      break;
    end
    queuedSend(true); -- will panic --
  end
  while (true) do
    var queuedRecv = chan.__recvQueue.shift();
    if (queuedRecv == nil) then
      break;
    end
    queuedRecv([chan.__elem.zero(), false]);
  end
end;
__select = function(comms)
   var ready = {};
   var selection = -1;
   --for i = 0; i < comms.length; i++) 
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
         if (#chan.__sendQueue ~= 0 or #chan.__buffer ~= 0 or chan.__closed then
                ready.push(i-1);
         end
         break;
      elseif comm_len == 2 then
         -- send --
         if chan.__closed then
            __throwRuntimeError("send on closed channel");
         end
         if (chan.__recvQueue.length ~= 0 or chan.__buffer.length < chan.__capacity) then
            ready.push(i-1);
         end
         break;
      end -- end switch
   end

  if (ready.length ~= 0) then
    selection = ready[Math.floor(Math.random() * ready.length)];
  end
  if (selection ~= -1) then
    var comm = comms[selection];
    switch (comm.length) {
    case 0: -- default --
      return [selection];
    case 1: -- recv --
      return [selection, __recv(comm[0])];
    case 2: -- send --
      __send(comm[0], comm[1]);
      return [selection];
    end
  end

  var entries = {};
  var thisGoroutine = __curGoroutine;
  var f = { __blk: function()  return this.selection; end };
  var removeFromQueues = function() 
    for (var i = 0; i < entries.length; i++) {
      var entry = entries[i];
      var queue = entry[0];
      var index = queue.indexOf(entry[1]);
      if (index ~= -1) then
        queue.splice(index, 1);
      end
    end
  end;
  for (var i = 0; i < comms.length; i++) {
    (function(i) 
      var comm = comms[i];
      switch (comm.length) {
      case 1: -- recv --
        var queueEntry = function(value) 
          f.selection = [i, value];
          removeFromQueues();
          __schedule(thisGoroutine);
        end;
        entries.push([comm[0].__recvQueue, queueEntry]);
        comm[0].__recvQueue.push(queueEntry);
        break;
      case 2: -- send --
        var queueEntry = function() 
          if (comm[0].__closed) then
            __throwRuntimeError("send on closed channel");
          end
          f.selection = [i];
          removeFromQueues();
          __schedule(thisGoroutine);
          return comm[1];
        end;
        entries.push([comm[0].__sendQueue, queueEntry]);
        comm[0].__sendQueue.push(queueEntry);
        break;
      end
    end)(i);
  end
  __block();
  return f;
end;

