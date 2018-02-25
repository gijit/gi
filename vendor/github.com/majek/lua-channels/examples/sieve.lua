local task = require('lua-channels')

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
   for i = 1, 10 do
      print(primes:recv())
   end
end

task.spawn(main)
task.scheduler()
