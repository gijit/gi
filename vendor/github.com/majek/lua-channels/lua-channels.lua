----------------------------------------------------------------------------
-- Go style Channels for Lua
--
-- This code is derived from libtask library by Russ Cox, mainly from
-- channel.c. Semantically channels as implemented here are quite
-- similar to channels from the Go language.
--
-- Usage (we're using unbuffered channel here):
--
-- local task = require('task')
--
-- local function counter(channel)
--    local i = 1
--    while true do
--        channel:send(i)
--        i = i + 1
--    end
-- end
--
-- local function main()
--     local channel = task.Channel:new()
--     task.spawn(counter, channel)
--     assert(channel:recv() == 1)
--     assert(channel:recv() == 2)
--     assert(channel:recv() == 3)
-- end
--
-- task.spawn(main)
-- task.scheduler()
--
--
-- This module exposes:
--
--  task.spawn(fun, [...]) - run fun as a coroutine with given
--                         parameters. You should use this instead of
--                         coroutine.create()
--
--  task.scheduler() - can be run only from the main thread, executes
--                     all the stuff, resumes the coroutines that are
--                     blocked on channels that became available. You
--                     can only do non-blocking sends / receives from
--                     the main thread.
--
--  task.Channel:new([buffer size]) - create a new channel with given size
--
--  task.chanalt(alts, can_block) - run alt / select / multiplex over
--                                  the alts structure. For example:
--
-- task.chanalt({{c = channel_1, op = task.RECV},
--               {c = channel_2, op = task.SEND, p = "hello"}}, true)
--
-- This will block current coroutine until it's possible to receive
-- from channel_1 or send to channel_2. chanalt returns a number of
-- statement from alts that succeeded (1 or 2 here) and a received
-- value if executed statement was RECV.
--
-- Finally, if two alt statements can be fulfilled at the same time,
-- we use math.random() to decide which one should go first. So it
-- makes sense to initialize seed with something random. If you don't
-- have access to an entropy source you can do:
--   math.randomseed(os.time())
-- but beware, the results of random() will predictable to a attacker.
----------------------------------------------------------------------------

local _M = {}

-- Constants
local RECV = 0x1
local SEND = 0x2
local NOP  = 0x3
local TIMEOUT = {err = "TIMEOUT"}
local luajit = not not (package.loaded['jit'] and jit.version_num)

-- Global objects for scheduler
local tasks_runnable = {}       -- list of coroutines ready to be resumed
local tasks_to = {}             -- all the timeout tasks
local altexec

----------------------------------------------------------------------------
--- Helpers

local function random_choice(arr)
   if #arr > 1 then
      return arr[math.random(#arr)]
   else
      return arr[1]
   end
end

-- Specialised Set data structure (with random element selection)
local Set = {
   new = function(self)
      local o = {a = {}, l = {}}; setmetatable(o, self); self.__index = self
      return o
   end,

   add = function(self, v)
      local a, l = self.a, self.l
      if a[v] == nil then
         table.insert(l, v)
         a[v] = #l
         return true
      end
   end,

   remove = function(self, v)
      local a, l = self.a, self.l
      local i = a[v]
      if i > 0 then
         local t = l[#l]
         a[t], l[i] = i, t
         a[i], l[#l] = nil, nil
         return true
      end
   end,

   random = function(self, to)
      if to then
         local arr = {}
         for i = 1, #self.l do
            if self.l[i].to then table.insert(arr, self.l[i]) end
         end
         return random_choice(arr)
      end
      return random_choice(self.l)
   end,

   len = function(self)
      return #self.l
   end,
}



-- Circular Buffer data structure
local CircularBuffer = {
   new = function(self, size)
      local o = {b = {}, slots = size + 1, size = size, l = 0, r = 0}
      setmetatable(o, self); self.__index = self
      return o
   end,

   len = function(self)
      return (self.r - self.l) % self.slots
   end,

   pop = function(self)
      assert(self.l ~= self.r)
      local v = self.b[self.l]
      self.l = (self.l + 1) % self.slots
      return v
   end,

   push = function(self, v)
      self.b[self.r] = v
      self.r = (self.r + 1) % self.slots
      assert(self.l ~= self.r)
   end,
}

----------------------------------------------------------------------------
-- Scheduling
--
-- Tasks ready to be run are placed on a stack and it's possible to
-- starve a coroutine.
local function scheduler()
   local self_coro, is_main = coroutine.running()

   -- We actually don't care if scheduler is run from the main
   -- coroutine. But we do need to make sure that user doesn't do
   -- blocking operation from it, as it can't yield.

   -- Be compatible with 5.1 and 5.2
   assert(not(self_coro ~= nil and is_main ~= true),
          "Scheduler must be run from the main coroutine.")

   local i = 0
   while #tasks_runnable > 0 do
      local co = table.remove(tasks_runnable)
      tasks_to[co] = nil
      local okay, emsg = coroutine.resume(co)
      if not okay then
         error(emsg)
      end
      i = i + 1
   end

   local now = os.time()
   for co, alt in pairs(tasks_to) do
      if alt and now >= alt.to then
         altexec(alt)
         tasks_to[co] = nil
         alt.c:_get_alts(RECV):remove(alt)
      end
   end
   return i
end

local function task_ready(co)
   table.insert(tasks_runnable, co)
end

local function spawn(fun, ...)
   local args = {...}

   local f = function()
      -- In luajit we could use pcall() here to produce nicer
      -- tracebacks on errors, but that won't work on vanilla lua
      -- (can't yield from within pcall).
      if not luajit then
         fun(unpack(args))
      else
         local okay, emsg = pcall(fun, unpack(args))
         if not okay then
            print(debug.traceback(emsg))
            error(emsg)
         end
      end
   end
   local co = coroutine.create(f)
   task_ready(co)
end

----------------------------------------------------------------------------
-- Channels - chanalt and helpers

-- Given two Alts from a single channel exchange data between
-- them. It's implied that one is RECV and another is SEND. Channel
-- may be buffered.
local function altcopy(a, b)
   local r, s, c = a, b, a.c
   if r.op == SEND then
      r, s = s, r
   end

   assert(s == nil or s.op == SEND)
   assert(r == nil or r.op == RECV)

   -- Channel is empty or unbuffered, copy directly
   if s ~= nil and r and c._buf:len() == 0 then
      r.alt_array.value = s.p
      return
   end

   -- Otherwise it's always okay to receive and then send.
   if r ~= nil then
      if r.to then
         r.alt_array.value = TIMEOUT
         r.alt_array.resolved = 1
         return true
      elseif r.closed then
         r.alt_array.value = nil
         r.alt_array.resolved = 1
         return true
      else
         r.alt_array.value = c._buf:pop()
      end
   end
   if s ~= nil then
      c._buf:push(s.p)
   end
end

-- Given enqueued alt_array from a chanalt statement remove all alts
-- from the associated channels.
local function altalldequeue(alt_array)
   for i = 1, #alt_array do
      local a = alt_array[i]
      if a.op == RECV or a.op == SEND then
         a.c:_get_alts(a.op):remove(a)
      end
   end
end

-- Can this Alt be execed without blocking?
local function altcanexec(a)
   local c, op = a.c, a.op
   if c._buf.size == 0 then
      if op ~= NOP then
         return c:_get_other_alts(op):len() > 0
      end
   else
      if op == SEND then
         return c._buf:len() < c._buf.size
      elseif op == RECV then
         return c._buf:len() > 0
      end
   end
end

-- Alt can be execed so find a counterpart Alt and exec it!
altexec = function (a)
   local c, op = a.c, a.op
   local other_alts = c:_get_other_alts(op)
   local other_a = other_alts:random(a.to)
   -- other_a may be nil
   local isend = altcopy(a, other_a)

   if other_a ~= nil then
      -- Disengage from channels used by the other Alt and make it ready.
      altalldequeue(other_a.alt_array)
      other_a.alt_array.resolved = other_a.alt_index
      task_ready(other_a.alt_array.task)
   elseif isend then
      task_ready(a.alt_array.task)
   end
end

-- The main entry point. Call it `alt` or `select` or just a
-- multiplexing statement. This is user facing function so make sure
-- the parameters passed are sane.
local function chanalt(alt_array, canblock)
   assert(#alt_array)

   local list_of_canexec_i = {}
   for i = 1, #alt_array do
      local a = alt_array[i]
      a.alt_array = alt_array
      a.alt_index = i
      assert(type(a.op) == "number" and
                (a.op == RECV or a.op == SEND or a.op == NOP),
             "op field must be RECV, SEND or NOP in alt")
      assert(type(a.c) == "table" and a.c.__index == _M.Channel,
             "pass valid channel to a c field of alt")
      if altcanexec(a) == true then
         table.insert(list_of_canexec_i, i)
      elseif a.to then

         local sc = coroutine.running()
         if not tasks_to[sc] then
            tasks_to[sc] = a
         end
      end
   end

   if #list_of_canexec_i > 0 then
      local i = random_choice(list_of_canexec_i)
      altexec(alt_array[i])
      return i, alt_array.value, alt_array.closed == nil
   end

   if canblock ~= true then
      return nil
   end

   local self_coro, is_main = coroutine.running()
   alt_array.task = self_coro
   assert(self_coro ~= nil and is_main ~= true,
          "Unable to block from the main thread, run scheduler.")

   for i = 1, #alt_array do
      local a = alt_array[i]
      if a.op ~= NOP then
         a.c:_get_alts(a.op):add(a)
      end
   end

   -- Make sure we're not woken by someone who is not the scheduler.
   alt_array.resolved = nil
   coroutine.yield()
   assert(alt_array.resolved > 0)

   local r = alt_array.resolved
   return r, alt_array.value, alt_array.closed == nil
end


----------------------------------------------------------------------------
-- Channel object

local Channel = {
   new = function(self, buf_size)
      local o = {}; setmetatable(o, self); self.__index = self
      o._buf = CircularBuffer:new(buf_size or 0)
      o._recv_alts, o._send_alts = Set:new(), Set:new()
      return o
   end,

   send = function(self, msg)
      assert(chanalt({{c = self, op = SEND, p = msg}}, true) == 1)
      return true
   end,

   recv = function(self, to)
      local alts = {{c = self, op = RECV, to = to and os.time() + to or nil}}
      local s, msg = chanalt(alts, true)
      assert(s == 1)
      return msg, alts[1].closed == nil
   end,

   nbsend = function(self, msg)
      local s = chanalt({{c = self, op = SEND, p = msg}}, false)
      return s == 1
   end,

   nbrecv = function(self)
      local s, msg = chanalt({{c = self, op = RECV}}, false)
      return s == 1, msg
   end,

   close = function(self)
      local alts = self:_get_alts(RECV)
      for _, v in ipairs(alts.l) do
         v.closed = true
         altexec(v)
      end
   end,

   _get_alts = function(self, op)
      return (op == RECV) and self._recv_alts or self._send_alts
   end,

   _get_other_alts = function(self, op)
      return (op == SEND) and self._recv_alts or self._send_alts
   end,

   __tostring = function(self)
      return string.format("<Channel size=%i/%i send_alt=%i recv_alt=%i>",
                           self._buf:len(), self._buf.size, self._send_alts:len(),
                           self._recv_alts:len())
   end,

   __call = function(self)
      local function f(s, v)
         return true, self:recv()
      end
      return f, nil, nil
   end,
}

----------------------------------------------------------------------------
-- Public interface

_M.scheduler = scheduler
_M.spawn     = spawn
_M.Channel   = Channel
_M.chanalt   = chanalt
_M.RECV      = RECV
_M.SEND      = SEND
_M.NOP       = NOP
_M.Error     = {TIMEOUT = TIMEOUT}
----------------------------------------------------------------------------
----------------------------------------------------------------------------
-- Tests
--
-- To run:
--    $ lua task.lua

local task = _M

local tests = {
   counter = function ()
      local done
      local function counter(c)
         local i = 1
         while true do
            c:send(i)
            i = i + 1
         end
      end
      local function main()
         local c = task.Channel:new()
         task.spawn(counter, c)
         assert(c:recv() == 1)
         assert(c:recv() == 2)
         assert(c:recv() == 3)
         assert(c:recv() == 4)
         assert(c:recv() == 5)
         done = true
      end
      task.spawn(main)
      task.scheduler()
      assert(done)
   end,

   nonblocking_channel = function()
      local done
      local function main()
         local b = task.Channel:new()
         assert(b:nbsend(1) == false)
         assert(b:nbrecv() == false)

         local c = task.Channel:new(1)
         assert(c:nbrecv() == false)
         assert(c:nbsend(1) == true)
         assert(c:nbsend(1) == false)
         local r, v = c:nbrecv()
         assert(r == true)
         assert(v == 1)
         assert(c:nbrecv() == false)
         done = true
      end
      task.spawn(main)
      task.scheduler()
      assert(done)
   end,

   concurrent_send_and_recv = function()
      local l = {}
      local function a(c, name)
         -- Blocking send and recv from the same process
         local alt = {{c = c, op = task.SEND, p = 1},
            {c = c, op = task.RECV}}
         local i, v = task.chanalt(alt, true)
         local k = string.format('%s %s', name, i == 1 and "send" or "recv")
         l[k] = (l[k] or 0) + 1
      end

      for i = 0, 1000 do
         -- On Mac OS X in lua 5.1 initializing seed with a
         -- predictable value makes no sense. For all seeds from 1 to
         -- 1000 the result of math.random(1,3) is _exactly_ the same!
         -- So beware, when seeding!
         -- math.randomseed(i)
         local c = task.Channel:new()
         task.spawn(a, c, "a")
         task.spawn(a, c, "b")
         task.scheduler()
      end

      -- Make sure we have randomness, that is: events occur in both
      -- orders in 1000 runs
      assert(l['a recv'] > 0)
      assert(l['a send'] > 0)
      assert(l['b recv'] > 0)
      assert(l['b send'] > 0)
   end,

   channels_from_a_coroutine = function()
      local done
      local c = task.Channel:new()
      local function a()
         for i = 1, 100 do
            c:send(i)
         end
      end
      local function b()
         assert(c:recv() == 1)
         assert(c:recv() == 2)
         assert(c:recv() == 3)
         assert(c:recv() == 4)
         assert(c:recv() == 5)
         done = true
      end
      local a_co = coroutine.create(a)
      local b_co = coroutine.create(b)
      coroutine.resume(a_co)
      coroutine.resume(b_co)
      task.scheduler()
      assert(done)
   end,

   fibonacci = function()
      local done
      local function fib(c)
         local x, y = 0, 1
         while true do
            c:send(x)
            x, y = y, x + y
         end
      end
      local function main(c)
         assert(c:recv() == 0)
         assert(c:recv() == 1)
         assert(c:recv() == 1)
         assert(c:recv() == 2)
         assert(c:recv() == 3)
         assert(c:recv() == 5)
         assert(c:recv() == 8)
         assert(c:recv() == 13)
         assert(c:recv() == 21)
         assert(c:recv() == 34)
         done = true
      end

      local c = task.Channel:new()
      task.spawn(fib, c)
      task.spawn(main, c)
      task.scheduler()
      assert(done)
   end,

   non_blocking_chanalt = function()
      local done
      local function main()
         local c = task.Channel:new()
         local alts = {{c = c, op = task.RECV},
                       {c = c, op = task.NOP},
                       {c = c, op = task.SEND, p = 1}}
         assert(task.chanalt(alts, false) == nil)

         local c = task.Channel:new(1)
         local alts = {{c = c, op = task.RECV},
                       {c = c, op = task.NOP},
                       {c = c, op = task.SEND, p = 1}}
         assert(task.chanalt(alts, false) == 3)
         assert(task.chanalt(alts, false) == 1)

         local alts = {{c = c, op = task.NOP}}
         assert(task.chanalt(alts, false) == nil)

         done = true
      end
      task.spawn(main)
      task.scheduler()
      assert(done)
   end,

   -- Apparently it's not really a Sieve of Eratosthenes:
   --   http://www.cs.hmc.edu/~oneill/papers/Sieve-JFP.pdf
   eratosthenes_sieve = function()
      local done
      local function counter(c)
         local i = 2
         while true do
            c:send(i)
            i = i + 1
         end
      end

      local function filter(p, recv_ch, send_ch)
         while true do
            local i = recv_ch:recv()
            if i % p ~= 0 then
               send_ch:send(i)
            end
         end
      end

      local function sieve(primes_ch)
         local c = task.Channel:new()
         task.spawn(counter, c)
         while true do
            local p, newc = c:recv(), task.Channel:new()
            primes_ch:send(p)
            task.spawn(filter, p, c, newc)
            c = newc
         end
      end

      local function main()
         local primes = task.Channel:new()
         task.spawn(sieve, primes)
         assert(primes:recv() == 2)
         assert(primes:recv() == 3)
         assert(primes:recv() == 5)
         assert(primes:recv() == 7)
         assert(primes:recv() == 11)
         assert(primes:recv() == 13)
         done = true
      end

      task.spawn(main)
      task.scheduler()
      assert(done)
   end,

   channel_as_iterator = function()
      local done
      local function counter(c)
         local i = 2
         while true do
            c:send(i)
            i = i + 1
         end
      end

      local function main()
         local numbers = task.Channel:new()
         task.spawn(counter, numbers)
         for _, j in numbers() do
            if j == 100 then
               break
            end
            done = true
         end
      end
      if _VERSION == "Lua 5.1" and not luajit then
         -- sorry, this doesn't work in 5.1
         print('skipping... (5.1 unsupported)')
         done = true
      else
         task.spawn(main)
         task.scheduler()
      end
      assert(done)
   end,

   close_test = function()
      local values = {}
      for i = 1, 1000 do table.insert(values, i) end
      local chan = task.Channel:new()
      local done = task.Channel:new()
      task.spawn(function()
            local i = 0
            while true do
               local msg, more = chan:recv()
               if not more then
                  break
               else
                  i = i + 1
               end
            end
            assert(i == 200)
            done:send(1)
      end)

      task.spawn(function()
            for i =1, 200 do
               chan:send(i)
            end
            chan:close()
      end)

      task.scheduler()
   end,

   recv_timeout = function()
      print("\t testing... (this will block 8 seconds)")
      -- use copas to sleep
      local status, copas = pcall(require, "copas")
      if not status then
         print("\t no copas install skip testing...")
         return
      end

      -- testing data
      local values = {}
      for i = 1, 1001 do
         table.insert(values, i)
      end

      local exists = function(t, v)
         for _, val in ipairs(t) do
            if val == v then return true end
         end
         return false
      end


      local chan = task.Channel:new(1000)

      -- we have three reader, only one will timeout
      local r1 = function()
         for i = 1, 500 do
            local v = chan:recv()
            assert(exists(values, v))
         end
      end

      local r2 = function()
         for i =1, 500 do
            local v = chan:recv()
            assert(exists(values, v))
         end
      end

      local r3 = function()
         local v = chan:recv(1)
         assert(task.Error.TIMEOUT == v)
      end

      local r4 = function()
         local v = chan:recv(6)
         assert(exists(values, v))
      end

      local sender = function()
         copas.sleep(5)
         for _, v in ipairs(values) do
            chan:send(v)
         end
      end

      task.spawn(r1)
      task.spawn(r2)
      task.spawn(r3)
      task.spawn(r4)
      copas.addthread(sender)

      local done
      done = true

      for i = 1, 80 do
         copas.step(0.1)
         task.scheduler()
      end
      assert(done)
   end

}

-- No parameters: run tests
local args = {...}
if #args == 0 then
   print("[*] Running tests...")
   local ok, failed = 0, 0
   for k, v in pairs(tests) do
      print(string.format('  - %s', k))
      v()
      ok = ok + 1
   end
   print(string.format("[*] Successfully run %i tests", ok))
else
   return _M
end
