local task = require("lua-channels")

function counter(c)
   local i = 1
   while true do
      local a = c:send(i)
      i = i + 1
   end
end

function a()
   local c = task.Channel:new()
   task.spawn(counter, c)

   for i = 1, 10 do
      local v = c:recv()
      print("recv", v)
   end
end

task.spawn(a)

math.randomseed(os.time())
task.scheduler()
