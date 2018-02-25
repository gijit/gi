[![Build Status](https://travis-ci.org/majek/lua-channels.png)](https://travis-ci.org/majek/lua-channels)

Lua-Channels
============

*Go style Channels for Lua*

This code is derived from libtask library by Russ Cox, mainly from
channel.c. Semantically channels as implemented here are quite
similar to channels from the Go language.

Usage
-----

This is an example of using an unbuffered channel:

```
local task = require('lua-channels')

local function counter(channel)
   local i = 1
   while true do
       channel:send(i)
       i = i + 1
   end
end

local function main()
    local channel = task.Channel:new()
    task.spawn(counter, channel)
    assert(channel:recv() == 1)
    assert(channel:recv() == 2)
    assert(channel:recv() == 3)
end

task.spawn(main)
task.scheduler()
```

lua-channels exposes:

 * task.spawn(fun, [...]) - run fun as a coroutine with given
                        parameters. You should use this instead of
                        coroutine.create()

 * task.scheduler() - can be run only from the main thread, executes
                    all the stuff, resumes the coroutines that are
                    blocked on channels that became available. You
                    can only do non-blocking sends / receives from
                    the main thread.

 * task.Channel:new([buffer size]) - create a new channel with given size

 * task.chanalt(alts, can_block) - run alt / select / multiplex over
                                 the alts structure. For example:

 * task.chanalt({{c = channel_1, op = task.RECV},
              {c = channel_2, op = task.SEND, p = "hello"}}, true)

This will block current coroutine until it's possible to receive
from channel_1 or send to channel_2. chanalt returns a number of
statement from alts that succeeded (1 or 2 here) and a received
value if executed statement was RECV.

Finally, if two alt statements can be fulfilled at the same time,
we use math.random() to decide which one should go first. So it
makes sense to initialize seed with something random. If you don't
have access to an entropy source you can do:

```
  math.randomseed(os.time())
```

but beware, the results of random() will predictable to a attacker.

Installing
----------

You may simply require src/lua-channels.lua from the source, or install
`lua-channels` from [luarocks](http://luarocks.org).

