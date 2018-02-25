-- reflect_goro.lua:
--
-- refactor from goroutines.lua into
-- their own file: implementation of
-- send/receive/select that interact
-- with native Go channels using reflect.


__recv_via_reflect = function(wchan)
   --print("__recv_via_reflect called! wchan=")
   --__st(wchan, "wchan")
   if wchan == nil then
      error("cannot read from nil channel")
   end

   -- unwrap
   if wchan.__native == nil then
      -- may not be wrapped; if native Go channel supplied
      chan = wchan
   else
      chan = wchan.__native.Interface()
   end

   local ch = reflect.ValueOf(chan)
   local rv, ok = ch.Recv();
   -- rv is userdata, a reflect.Value. Convert to
   -- interface{} for Luar, using Interface(), so
   -- luar can translate that to Lua for us.
   local v = rv.Interface();
   return {v, ok}
end

__send_via_reflect = function(wchan, value)
   --print("__send_via_reflect called! value=", value)
   --__st(value, "value")
   --__st(wchan, "wchan")
   if wchan == nil then
      error("cannot send on nil channel")
   end

   -- unwrap
   if wchan.__native == nil then
      -- may not be wrapped; if native Go channel supplied
      chan = wchan
   else
      chan = wchan.__native.Interface()
   end

   local ch = reflect.ValueOf(chan)
   local v = reflect.ValueOf(value)
   local cv = v.Convert(reflect.TypeOf(chan).Elem())
   ch.Send(cv);
end

__select_via_reflect = function(comms)
   --print("__select_via_reflect called!")
   --__st(comms, "comms")
   
   --__st(comms[1], "comms[1]")
   --__st(comms[2], "comms[2]")
   --__st(comms[3], "comms[3]")
   
   --__st(comms[1][1], "comms[1][1]")
   --__st(comms[2][1], "comms[2][1]")
   --__st(comms[3][1], "comms[3][1]")

   local c1 = reflect.ValueOf(comms[1][1])
   local c2 = reflect.ValueOf(comms[2][1])

   --print("c1 is "..type(c1))
   --print("c2 is "..type(c2))
      
   --print("c1 = "..tostring(c1))
   --print("c2 = "..tostring(c2))

   local cases = {}
   local casesType = {}
   local rty = reflect.TypeOf(__refSelCaseVal).Elem()
   
    for i, comm in ipairs(comms) do
      local chan = comm[1];
      --switch (comm.length)
      local comm_len = #comm

      local newCase = reflect.New(rty).Interface()
      
      if comm_len == 0 then
         -- default --
         --print("comm_len is 0/default at i =", i)
         newCase.Dir  = 3
         table.insert(cases, newCase)
         casesType[i-1]="d"
                      
      elseif comm_len == 1 then
         -- recv --
         --print("comm_len is 1/recv at i =", i)
         newCase.Chan = reflect.ValueOf(comm[1])
         newCase.Dir  = 2
         --print("newCase is ", newCase)
         --fmt.Printf("newCase is %#v\n", newCase)
         --fmt.Printf("newCase.dir is %#v\n", newCase.Dir)
         table.insert(cases, newCase)
         casesType[i-1]="r"
         
      elseif comm_len == 2 then
         -- send --
         --print("comm_len is 2/send at i =", i)

         newCase.Chan = reflect.ValueOf(comm[1])
         newCase.Dir  = 1
         newCase.Send = reflect.ValueOf(comm[2]) -- maybe comm[2]? or?
         --__st(comm, "comm in send case")
         --fmt.Printf("in send, newCase.Send is %#v\n", newCase.Send)
         
         table.insert(cases, newCase)
         casesType[i-1]="s"
         
      end -- end switch
   end
   
   local chosen, recv, recvOk = reflect.Select(cases)
   --print("back from reflect.Select, we got: chosen=", chosen)
   --print("back from reflect.Select, we got:   recv=", recv)
   --print("back from reflect.Select, we got: recvOk=", recvOk)

   local recvVal = nil
   --__st(casesType, "casesType")
   --print("chosen is ", chosen, " and casesType[chosen]= ", casesType[tonumber(chosen)])
   if casesType[tonumber(chosen)]=="r" then
      recvVal = recv.Interface()
      --print("recvVal in receive case is ", recvVal)
   end
   return {chosen, {recvVal, recvOk}};
end
