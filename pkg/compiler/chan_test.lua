dofile 'channel.lua'

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
