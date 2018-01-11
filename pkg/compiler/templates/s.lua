-- recursive pcalls, do they work?

f1=function()
   error("ouch") end

f2=function()
   g={pcall(f1)}
   for k,v in pairs(g) do print("f2: g result of recursive pcall is k=",k," val=",v) end
   error("panic in f2")
end

f3=function()
   g={pcall(f2)}
   for k,v in pairs(g) do print("f3: g result of recursive pcall is k=",k," val=",v) end
   error("panic in f3")
end

f4=function()
   g={pcall(f3)}
   for k,v in pairs(g) do print("f4: g result of recursive pcall is k=",k," val=",v) end
end

f4()
