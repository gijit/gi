
stack = {}

function stack.new()
   return {}
end

function stack.pop(self)
   if #self == 0 then
     return nil
   end
   top = self[#self]
   self[#self] = nil
   return top
end

function stack.push(self, x)
   self[#self +1] = x
end

function stack.show(self)
  for i, v in ipairs(self) do
      print(i, v)
  end
end

test = function() 
 q = stack:new()
 stack.show(q)
 stack.push(q, "hi")
 stack.push(q, "there")
 stack.push(q, "jason")
 stack.show(q)
 rrr=stack.pop(q)
 print(rrr)
 stack.show(q)
 rrr=stack.pop(q)
 print(rrr)
 stack.show(q)
 rrr=stack.pop(q)
 print(rrr)
 stack.show(q)
 rrr=stack.pop(q)
 print(rrr)
end