stack = {}

function stack:new()
   return {}
end

function stack:pop()
   if #self == 0 then
     return nil
   end
   top = self[#self]
   self[#self] = nil
   return top
end

function stack:push(x)
   self[#self +1] = x
end

function stack:show()
  for i, v in ipairs(self) do
      print(i, v)
  end
end
