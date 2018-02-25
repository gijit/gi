dofile 'tsys.lua'
dofile 'tutil.lua'
dofile 'chan.lua'

-- tests for 

local function counter(channel)
   local i = 1
   while true do
       channel:send(i)
       i = i + 1
   end
end

local function main()
    local channel = __task.Channel:new()
    __task.spawn(counter, channel)
    assert(channel:recv() == 1)
    assert(channel:recv() == 2)
    assert(channel:recv() == 3)
end

__task.spawn(main)
__task.scheduler()

__expectEq("", "")
