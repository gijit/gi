dofile 'channel.lua'


local __chan_tests = {
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

--  run tests

   print("[*] Running tests...")
   local ok, failed = 0, 0
   for k, v in pairs(__chan_tests) do
      print(string.format('  - %s', k))
      v()
      ok = ok + 1
   end
   print(string.format("[*] Successfully run %i tests", ok))

