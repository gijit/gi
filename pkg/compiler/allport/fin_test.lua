dofile 'fin.lua'

function expectEq(a, b)
   if a ~= b then
      error("expectEq failure: a="..tostring(a).." was not equal "..tostring(b))
   end
end

expectEq("", __mapAndJoinStrings("_", {}, function(x) return x end))
expectEq("a_b_c", __mapAndJoinStrings("_", {"a","b","c"}, function(x) return x end))
expectEq("a", __mapAndJoinStrings("_", {"a"}, function(x) return x end))
expectEq("a", __mapAndJoinStrings("_", {[0]="a"}, function(x) return x end))
expectEq("a_b", __mapAndJoinStrings("_", {[0]="a","b"}, function(x) return x end))
expectEq("1_2_3", __mapAndJoinStrings("_", {[0]=0,1,2}, function(x)  return x+1 end))
expectEq("1_2", __mapAndJoinStrings("_", {[0]=0,1}, function(x)  return x+1 end))
expectEq("1", __mapAndJoinStrings("_", {[0]=0}, function(x)  return x+1 end))
expectEq("", __mapAndJoinStrings("_", {}, function(x)  return x+1 end))
