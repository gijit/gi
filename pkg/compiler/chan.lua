-- chan.lua
-- Derived from lua-channels.lua. Portions
-- Copyright (c) 2013 Marek Majkowski
-- used under the MIT license and similar from libtask upstream, see
-- github.com/gijit/gi/vendor/github.com/majek/lua-channels/LICENSE-MIT-lua-channels
--
----------------------------------------------------------------------------
-- Go style Channels for Lua
--
-- This code is derived from libtask library by Russ Cox, mainly from
-- channel.c. Semantically channels as implemented here are quite
-- similar to channels from the Go language.
--
-- Usage (we're using unbuffered channel here):
--
-- local __task = require('chan')
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
--     local channel = __task.Channel:new()
--     __task.spawn(counter, channel)
--     assert(channel:recv() == 1)
--     assert(channel:recv() == 2)
--     assert(channel:recv() == 3)
-- end
--
-- __task.spawn(main)
-- __task.scheduler()
--
--
-- This module exposes:
--
--  __task.spawn(fun, [...]) - run fun as a coroutine with given
--                         parameters. You should use this instead of
--                         coroutine.create()
--
--  __task.scheduler() - can be run only from the main thread, executes
--                     all the stuff, resumes the coroutines that are
--                     blocked on channels that became available. You
--                     can only do non-blocking sends / receives from
--                     the main thread.
--
--  __task.Channel:new([buffer size]) - create a new channel with given size
--
--  __task.select(alts, can_block) - run alt / select / multiplex over
--                                  the alts structure. For example:
--
-- __task.select({{c = channel_1, op = __task.RECV},
--               {c = channel_2, op = __task.SEND, p = "hello"}}, true)
--
-- This will block current coroutine until it's possible to receive
-- from channel_1 or send to channel_2. select returns a number of
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


local __M = {}

-- Constants
local RECV = "recv" -- 0x1
local SEND = "send" -- 0x2
local NOP  = "nop"  -- 0x3
local TIMEOUT = {err = "TIMEOUT"}

-- Global objects for scheduler
local tasks_runnable = {}       -- list of coroutines ready to be resumed
local tasks_to = {}             -- all the timeout tasks
local altexec

local main_coro, is_main = coroutine.running()
if not is_main then
   error("must be loaded, for now, by main coroutine")
end

local scheduler_co

local parkedForeverPool = coroutine.create(function() end)

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



-- Circular Buffer data structure, used for channels with buffers.
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
   print("top of scheduler")
   
   -- returns nil if running on the main coroutine, otherwise
   -- returns the running coro.
   local self_coro, is_main = coroutine.running()

   -- We actually don't care if scheduler is run from the main
   -- coroutine. But we do need to make sure that user doesn't do
   -- blocking operation from it, as it can't yield.

   -- jea: I think we want to add a call to scheduler
   --      from the main coroutine, after each repl command,
   --      so as to let all background goroutines run.

   -- Be compatible with 5.1 and 5.2
   --assert(not(self_coro ~= nil and is_main ~= true),
   --      "Scheduler must be run from the main coroutine.")
   
   print("scheduler: self_coro is ", self_coro)
   print("scheduler: is_main is ", is_main)
   
   print("scheduler: past assert")
   
   local i = 0
   while true do
      if #tasks_runnable == 0 then
         print("scheduler: no more runnable tasks")
         break
      end
      -- table.remove takes the last by default.
      local co = table.remove(tasks_runnable)
      tasks_to[co] = nil

      -- and resume co
      print("scheduler: coroutine.resume about to be called on co="..tostring(co))
      local back = {coroutine.resume(co)}
      print("scheduler: got back from resume: ", unpack(back))
      
      okay, emsg = unpack(back)
      if not okay then
         print(debug.traceback(emsg))
         error(emsg)
      end
      i = i + 1
      print("scheduler: resume was okay, i is now = ", i)      
   end

   local now = __abs_now()
   print("scheduler: checking for timeouts, here is tasks_to: "..type(tasks_to))
   
   local k = 0
   print("scheduler: just before pairs(tasks_to)")
   for co, alt in pairs(tasks_to) do
      print("scheduler: top of tasks_to loop")
      print("scheduler: on tasks_to, on k=",k,"  we have co = ", co, " and alt=", alt)
      if alt and now >= alt.to then
         altexec(alt)
         tasks_to[co] = nil
         alt.c:_get_alts(RECV):remove(alt)
      end
   end
   
   print("end of scheduler, we ran i tasks, returning i=", i)   
   return i
end

local function task_ready(co)
   table.insert(tasks_runnable, co)
end

local function spawn(fun, args)
   --local args = {...}

   local f = function()
      local okay, emsg = pcall(fun, unpack(args))
      if not okay then
         print(debug.traceback(emsg))
         error(emsg)
      end
   end
   local co = coroutine.create(f)
   task_ready(co)
   print("spawn added to ready queue: co = ", co)
end

----------------------------------------------------------------------------
-- Channels - select and helpers

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

-- Given enqueued alt_array from a select statement remove all alts
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

local function __fldcnt(t)
   if type(t) ~= "table" then
      return 0
   end
   local k = 0
   for _,_ in pairs(t) do k=k+1; end
   return k
end

-- The main entry point. Call it `alt` or `select` or just a
-- multiplexing statement. This is user facing function so make sure
-- the parameters passed are sane.
local function select(alt_array)
   print("top of select, alt_array is")
   __st(alt_array, "alt_array")
   for i,_ in ipairs(alt_array) do
      __st(alt_array[i], "alt_array["..i.."]", 6)
      if type(alt_array[i]) == "table" and type(alt_array[i][1]) == "table" then
         __st(alt_array[i][1], "alt_array["..i.."][1]", 10)
      end
   end
   
   local defaultPresent = nil
   local canblock = true
   
   -- detect the default only {{}} case.
   if #alt_array==1 and __fldcnt(alt_array[1]) == 0 then
      -- select{ default: }
      -- only the default channel
      defaultPresent = true
      defaultNum = 0LL
      canblock =false
      print("select: {{}} default only case seen!")
      
   elseif #alt_array==0 then
      -- no default, no cases: "select{}".
      -- block this goroutine forever
      local self_coro, is_main = coroutine.running()
      print("warning: select{} is blocking the goroutine forever... ", self_coro)
      -- just set ourselves to state normal? call the scheduler?

      -- jea: resume puts us in "normal" mode
      -- where we cannot be scheduled again without yielding
      -- back to us. And the scheduler should never yield.
      -- So this should effectively park us forever, just
      -- like select{} does. We pause forever, without
      -- running any defers, just as select{} would.
      coroutine.resume(scheduler_co)
   end

   print("select: loop through the alt_array...")   
   local list_of_canexec_i = {}
   for i = 1, #alt_array do
      print("top of alt_array loop, i = ", i)
      local a = alt_array[i]

      if __fldcnt(a) == 0 then
         print("select: default case observed to be present.")
         --default: option
         defaultPresent = true
         canblock = false
         defaultNum = int(i)-1
         i=i+1
         goto zcontinue
      end
      print("select: non default case! i=",i)
      
      a.alt_array = alt_array
      a.alt_index = i
      print("a.op is ", a.op, " of type "..type(a.op))
      assert(type(a.op) == "string" and
                (a.op == RECV or a.op == SEND or a.op == NOP),
             "op field must be RECV, SEND or NOP in alt")
      assert(type(a.c) == "table" and a.c.__index == __M.Channel,
             "pass valid channel to a c field of alt")
      if altcanexec(a) == true then
         table.insert(list_of_canexec_i, i)
      elseif a.to then

         local sc = coroutine.running()
         if not tasks_to[sc] then
            print("select is adding sc to the tasks_to table as key. value a")
            tasks_to[sc] = a
         end
      end
      ::zcontinue::
   end
   print("select: done with alt_array loop") 

   if #list_of_canexec_i > 0 then
      if #list_of_canexec_i > 1 then
         print("select: multiple choices from alt_array, can proceed... choosing one at random")
      else
         print("select: one choice from alt_array can proceed.")
      end
      local i = random_choice(list_of_canexec_i)
      altexec(alt_array[i])
      return i, alt_array.value, alt_array.closed == nil
   else
      print("select: no cases to execute.")
   end

   print("select: defaultPresent = ", defaultPresent)   
   if defaultPresent then
      print("select choosing our default: ")
      return {defaultNum}
   end
   print("select: no default present, going on to block...")
   
   local self_coro, is_main = coroutine.running()
   alt_array.task = self_coro
   if not (self_coro ~= nil and is_main ~= true) then
      local err = "Unable to block from the main thread, run scheduler."
      print(debug.traceback(err))
      error(err)
   end

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
   new = function(self, buf_size, elemTyp)
      local o = {__elemTyp=elemTyp, __name="__valChannel"};
      setmetatable(o, self);
      self.__index = self
      o._buf = CircularBuffer:new(buf_size or 0)
      o._recv_alts, o._send_alts = Set:new(), Set:new()
      return o
   end,

   send = function(self, msg)
      assert(select({{c = self, op = SEND, p = msg}}, true) == 1)
      return true
   end,

   recv = function(self, to)
      local alts = {{c = self, op = RECV, to = to and __abs_now() + to or nil}}
      local s, msg = select(alts, true)
      assert(s == 1)
      return msg, alts[1].closed == nil
   end,

   nbsend = function(self, msg)
      local s = select({{c = self, op = SEND, p = msg}}, false)
      return s == 1
   end,

   nbrecv = function(self)
      local s, msg = select({{c = self, op = RECV}}, false)
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
      if op == RECV then
         return self._recv_alts
      end
      return self._send_alts
   end,

   _get_other_alts = function(self, op)
      if op == SEND then
         return self._recv_alts
      end
      return self._send_alts
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

scheduler_co = coroutine.create(scheduler)

----------------------------------------------------------------------------
-- Public interface

__task = __M

__task.resume_scheduler = function()
   print("__task.resume_scheduler called! scheduler_co is:")
   __st(scheduler_co, "scheduler_co")
   local ok, err = coroutine.resume(scheduler_co)
   -- crashes before we get to next line.
   print("__task.resume_scheduler back from coroutine.resume(scheduler_co)")
   if not ok then
      print("error detected in __task.resume_scheduler!")
      print(debug.traceback(err))
      error(err)
   end
end

__task.scheduler = scheduler
__task.spawn     = spawn
__task.Channel   = Channel
__task.select    = select
__task.RECV      = RECV
__task.SEND      = SEND
__task.NOP       = NOP
__task.Error     = {TIMEOUT = TIMEOUT}
----------------------------------------------------------------------------
----------------------------------------------------------------------------

__send = function(chan, value)
   return chan:send(value)
end

__recv = function(chan)
   -- return both, wrapped in a tuple for now, since that's what js had to do.
   return {chan:recv()}
end
