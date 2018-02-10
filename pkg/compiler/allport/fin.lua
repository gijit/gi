
-- apply fun to each element of the array arr,
-- then concatenate them together with splice in
-- between each one. It arr is empty then we
-- return the empty string. arr can start at
-- [0] or [1].
function __mapAndJoinStrings(splice, arr, fun)
   local newarr = {}
   -- handle a zero argument, if present.
   local bump = 0
   local zval = arr[0]
   if zval ~= nil then
      bump = 1
      newarr[1] = fun(zval)
   end
   for i,v in ipairs(arr) do
      newarr[i+bump] = fun(v)
   end
   return table.concat(newarr, splice)
end
