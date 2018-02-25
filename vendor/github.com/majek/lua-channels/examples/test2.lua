local task = require("lua-channels")


function a(c)
   -- Blocking send and recv from the same process
   local alt = {{c = c, op = task.SEND, p = "from a"},
                {c = c, op = task.RECV}}

   print("a", task.chanalt(alt, true))
end

function b(c)
   local alt = {{c = c, op = task.SEND, p = "from b"},
                {c = c, op = task.RECV}}
   print("b", task.chanalt(alt, true))
end

local c = task.Channel:new()

task.spawn(a, c)
task.spawn(b, c)

math.randomseed(os.time())
task.scheduler()
