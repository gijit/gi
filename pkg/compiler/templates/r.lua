-- recursive xpcalls, do they work? hmm... looks like we get "error in error handling"

ouch=function() error("ouch") end

ouch2=function() error("ouch2") end

ok=function() return "ok" end

h2 = function(err)
   print("panicHandler2 running with err =", err) -- hmm, can't get this to call?
end

h = function(err)
   print("panicHandler running with err =", err)
   g = {xpcall(ouch2, h2)}
   for k,v in pairs(g) do print("g result of recursive xpcall is k=",k," val=",v) end
   -- g result of recursive xpcall is k=	1	 val=	false
   -- g result of recursive xpcall is k=	2	 val=	error in error handling
end

r={xpcall(ouch, h)}

for k,v in pairs(r) do print("r result of top xpcall is k=",k," val=",v) end
