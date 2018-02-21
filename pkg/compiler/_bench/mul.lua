dofile 'tsys.lua'

--__go_import "math/rand"
--__go_import "fmt"
--__go_import "time"

__type__.anon_sliceType = __sliceType(__type__.float64); -- 'IMMEDIATE' anon type printing. utils.go:506

__type__.Matrix = __newType(0, __kindStruct, "main.Matrix", true, "main", true, nil);

__type__.anon_sliceType_1 = __sliceType(__type__.anon_sliceType); -- 'DELAYED' anon type printing. utils.go:506

__type__.Matrix.__methods_desc = {{prop= "Set", __name= "Set", __pkg="", __typ= __funcType({__type__.int, __type__.int, __type__.float64}, {}, false)}, {prop= "Get", __name= "Get", __pkg="", __typ= __funcType({__type__.int, __type__.int}, {__type__.float64}, false)}}; -- incr.go:817 for methods


__type__.anon_ptrType = __ptrType(__type__.Matrix); -- 'IMMEDIATE' anon type printing. utils.go:506

__type__.Matrix.ptr.__methods_desc = {{prop= "Set", __name= "Set", __pkg="", __typ= __funcType({__type__.int, __type__.int, __type__.float64}, {}, false)}, {prop= "Get", __name= "Get", __pkg="", __typ= __funcType({__type__.int, __type__.int}, {__type__.float64}, false)}}; -- incr.go:827 for ptr_methods

__type__.Matrix.init("", {{__prop= "A", __name= "A", __anonymous= false, __exported= true, __typ= __type__.anon_sliceType_1, __tag= ""}, {__prop= "Nrow", __name= "Nrow", __anonymous= false, __exported= true, __typ= __type__.int, __tag= ""}, {__prop= "Ncol", __name= "Ncol", __anonymous= false, __exported= true, __typ= __type__.int, __tag= ""}}); -- incr.go:873

__type__.Matrix.__constructor = function(self, ...) 
   if self == nil then self = {}; end
   local A_, Nrow_, Ncol_ = ... ;
   self.A = A_ or __type__.anon_sliceType_1.__nil;
   self.Nrow = Nrow_ or 0LL;
   self.Ncol = Ncol_ or 0LL;
   return self; 
end;
;

NewMatrix = function(nrow, ncol, fill) 
   local m = __type__.Matrix.ptr({}, __makeSlice(__type__.anon_sliceType_1, nrow), nrow, ncol);
   do  local i = 0; 
      local __gensym_2_i = 0; local __gensym_1__lim = __lenz(m.A);
      while __gensym_2_i < __gensym_1__lim do
         
         i = __gensym_2_i;

         m.A[i] = __makeSlice(__type__.anon_sliceType, ncol);
         
         __gensym_2_i=__gensym_2_i+1;

   end end;
   
   local next = 2;
   if (fill) then 
      do  local i = 0; 
         local __gensym_5_i = 0; local __gensym_4__lim = __lenz(m.A);
         while __gensym_5_i < __gensym_4__lim do
            
            i = __gensym_5_i;

			do  local j = 0; 
               local __gensym_8_i = 0; local __gensym_7__lim = __lenz(m.A[i]);
               while __gensym_8_i < __gensym_7__lim do
                  
                  j = __gensym_8_i;

                  m.A[i][j] = next
                  --__gi_SetRangeCheck(__gi_GetRangeCheck(m.A, i), j, next);
                  next = next + (1);
                  
                  __gensym_8_i=__gensym_8_i+1;

            end end;
            
			
            __gensym_5_i=__gensym_5_i+1;

      end end;
      
   end 
   return  m ;
end;
__pkg.NewMatrix = NewMatrix;
mult = function(m1, m2) 
   local r = __type__.anon_ptrType.__nil;
   if ( not ((m1.Ncol == m2.Nrow))) then 
      panic("incompatible: dimensions")
   end 
   r = NewMatrix(m1.Nrow, m2.Ncol, false);
   local nr1 = m1.Nrow;
   local nr2 = m2.Nrow;
   local nc2 = m2.Ncol;
   local i = 0LL;
   while (true) do
      if (not (i < nr1)) then break; end
      local k = 0LL;
      while (true) do
         if (not (k < nr2)) then break; end
         local j = 0LL;
         while (true) do
            if (not (j < nc2)) then break; end
            local a = r:Get(i, j);
            a = a + (m1:Get(i, k) * m2:Get(k, j));
            r:Set(i, j, a);
            j = j + (1LL);
         end 
         k = k + (1LL);
      end 
      i = i + (1LL);
   end 
   return  r ;
end;
__type__.Matrix.ptr.prototype.Set = function(m,i, j, val) 
   m.A[i][j] = val;
end;
__type__.Matrix.prototype.Set = function(this , i, j, val)  return this.__val.Set(this , i, j, val); end;

__type__.Matrix.__addToMethods({prop= "Set", __name= "Set", __pkg="", __typ= __funcType({__type__.int, __type__.int, __type__.float64}, {}, false)}); -- package.go:344

__type__.Matrix.ptr.__addToMethods({prop= "Set", __name= "Set", __pkg="", __typ= __funcType({__type__.int, __type__.int, __type__.float64}, {}, false)}); -- package.go:346
__type__.Matrix.ptr.prototype.Get = function(m,i, j) 
   return  m.A[i][j] ;
end;
__type__.Matrix.prototype.Get = function(this , i, j)  return this.__val.Get(this , i, j); end;

__type__.Matrix.__addToMethods({prop= "Get", __name= "Get", __pkg="", __typ= __funcType({__type__.int, __type__.int}, {__type__.float64}, false)}); -- package.go:344

__type__.Matrix.ptr.__addToMethods({prop= "Get", __name= "Get", __pkg="", __typ= __funcType({__type__.int, __type__.int}, {__type__.float64}, false)}); -- package.go:346
MatScaMul = function(m, x) 
   local r = __type__.anon_ptrType.__nil;
   r = NewMatrix(m.Nrow, m.Ncol, false);
   local i = 0LL;
   while (true) do
      if (not (i < m.Nrow)) then break; end
      local j = 0LL;
      while (true) do
         if (not (j < m.Ncol)) then break; end
         r:Set(i, j, x * m:Get(i, j));
         j = j + (1LL);
      end 
      i = i + (1LL);
   end 
   return  r ;
end;
__pkg.MatScaMul = MatScaMul;

done = false;
runMultiply = function(sz, i, j) 
   local mu = __type__.anon_ptrType.__nil;
   local k = 0LL;
   while (true) do
      if (not (k < 1LL)) then break; end
      local a = NewMatrix(sz, sz, true);
      local b = NewMatrix(sz, sz, true);
      mu = mult(a, b);
      k = k + (1LL);
   end 
   done = true;
   return  mu.A[i][j] ;
end;

main = function() 
   --local r = runMultiply(10LL, 9LL, 9LL);
   --print("r=", r);
   --t0 = __clone(time.Now(), __type__.time.Time);
   print("runMultiply(100,9,9) -> ", runMultiply(100LL, 9LL, 9LL))
   --local elap = time.Since(__clone(t0, __type__.time.Time));
   --fmt.Printf("compiled Go elap = %v\n", elap);
end;
