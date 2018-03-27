package shadow_sync

import "sync"

var Pkg = make(map[string]interface{})
var Ctor = make(map[string]interface{})

func init() {
    Ctor["Cond"] = GijitShadow_NewStruct_Cond
    Pkg["Locker"] = GijitShadow_InterfaceConvertTo2_Locker
    Ctor["Map"] = GijitShadow_NewStruct_Map
    Ctor["Mutex"] = GijitShadow_NewStruct_Mutex
    Pkg["NewCond"] = sync.NewCond
    Ctor["Once"] = GijitShadow_NewStruct_Once
    Ctor["Pool"] = GijitShadow_NewStruct_Pool
    Ctor["RWMutex"] = GijitShadow_NewStruct_RWMutex
    Ctor["WaitGroup"] = GijitShadow_NewStruct_WaitGroup

}
func GijitShadow_NewStruct_Cond(src *sync.Cond) *sync.Cond {
    if src == nil {
	   return &sync.Cond{}
    }
    a := *src
    return &a
}


func GijitShadow_InterfaceConvertTo2_Locker(x interface{}) (y sync.Locker, b bool) {
	y, b = x.(sync.Locker)
	return
}

func GijitShadow_InterfaceConvertTo1_Locker(x interface{}) sync.Locker {
	return x.(sync.Locker)
}


func GijitShadow_NewStruct_Map(src *sync.Map) *sync.Map {
    if src == nil {
	   return &sync.Map{}
    }
    a := *src
    return &a
}


func GijitShadow_NewStruct_Mutex(src *sync.Mutex) *sync.Mutex {
    if src == nil {
	   return &sync.Mutex{}
    }
    a := *src
    return &a
}


func GijitShadow_NewStruct_Once(src *sync.Once) *sync.Once {
    if src == nil {
	   return &sync.Once{}
    }
    a := *src
    return &a
}


func GijitShadow_NewStruct_Pool(src *sync.Pool) *sync.Pool {
    if src == nil {
	   return &sync.Pool{}
    }
    a := *src
    return &a
}


func GijitShadow_NewStruct_RWMutex(src *sync.RWMutex) *sync.RWMutex {
    if src == nil {
	   return &sync.RWMutex{}
    }
    a := *src
    return &a
}


func GijitShadow_NewStruct_WaitGroup(src *sync.WaitGroup) *sync.WaitGroup {
    if src == nil {
	   return &sync.WaitGroup{}
    }
    a := *src
    return &a
}



 func InitLua() string {
  return `
__type__.sync ={};

-----------------
-- struct Cond
-----------------

__type__.sync.Cond = {
 __name = "native_Go_struct_type_wrapper",
 __native_type = "Cond",
 __str = "Cond",
 exported = true,
 __call = function(t, src)
   return __ctor__sync.Cond(src)
 end,
};
setmetatable(__type__.sync.Cond, __type__.sync.Cond);


-----------------
-- struct Map
-----------------

__type__.sync.Map = {
 __name = "native_Go_struct_type_wrapper",
 __native_type = "Map",
 __str = "Map",
 exported = true,
 __call = function(t, src)
   return __ctor__sync.Map(src)
 end,
};
setmetatable(__type__.sync.Map, __type__.sync.Map);


-----------------
-- struct Mutex
-----------------

__type__.sync.Mutex = {
 __name = "native_Go_struct_type_wrapper",
 __native_type = "Mutex",
 __str = "Mutex",
 exported = true,
 __call = function(t, src)
   return __ctor__sync.Mutex(src)
 end,
};
setmetatable(__type__.sync.Mutex, __type__.sync.Mutex);


-----------------
-- struct Once
-----------------

__type__.sync.Once = {
 __name = "native_Go_struct_type_wrapper",
 __native_type = "Once",
 __str = "Once",
 exported = true,
 __call = function(t, src)
   return __ctor__sync.Once(src)
 end,
};
setmetatable(__type__.sync.Once, __type__.sync.Once);


-----------------
-- struct Pool
-----------------

__type__.sync.Pool = {
 __name = "native_Go_struct_type_wrapper",
 __native_type = "Pool",
 __str = "Pool",
 exported = true,
 __call = function(t, src)
   return __ctor__sync.Pool(src)
 end,
};
setmetatable(__type__.sync.Pool, __type__.sync.Pool);


-----------------
-- struct RWMutex
-----------------

__type__.sync.RWMutex = {
 __name = "native_Go_struct_type_wrapper",
 __native_type = "RWMutex",
 __str = "RWMutex",
 exported = true,
 __call = function(t, src)
   return __ctor__sync.RWMutex(src)
 end,
};
setmetatable(__type__.sync.RWMutex, __type__.sync.RWMutex);


-----------------
-- struct WaitGroup
-----------------

__type__.sync.WaitGroup = {
 __name = "native_Go_struct_type_wrapper",
 __native_type = "WaitGroup",
 __str = "WaitGroup",
 exported = true,
 __call = function(t, src)
   return __ctor__sync.WaitGroup(src)
 end,
};
setmetatable(__type__.sync.WaitGroup, __type__.sync.WaitGroup);


`}