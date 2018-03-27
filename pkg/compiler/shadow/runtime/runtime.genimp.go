package shadow_runtime

import "runtime"

var Pkg = make(map[string]interface{})
var Ctor = make(map[string]interface{})

func init() {
    Pkg["BlockProfile"] = runtime.BlockProfile
    Ctor["BlockProfileRecord"] = GijitShadow_NewStruct_BlockProfileRecord
    Pkg["Breakpoint"] = runtime.Breakpoint
    Pkg["CPUProfile"] = runtime.CPUProfile
    Pkg["Caller"] = runtime.Caller
    Pkg["Callers"] = runtime.Callers
    Pkg["CallersFrames"] = runtime.CallersFrames
    Pkg["Compiler"] = runtime.Compiler
    Pkg["Error"] = GijitShadow_InterfaceConvertTo2_Error
    Ctor["Frame"] = GijitShadow_NewStruct_Frame
    Ctor["Frames"] = GijitShadow_NewStruct_Frames
    Ctor["Func"] = GijitShadow_NewStruct_Func
    Pkg["FuncForPC"] = runtime.FuncForPC
    Pkg["GC"] = runtime.GC
    Pkg["GOARCH"] = runtime.GOARCH
    Pkg["GOMAXPROCS"] = runtime.GOMAXPROCS
    Pkg["GOOS"] = runtime.GOOS
    Pkg["GOROOT"] = runtime.GOROOT
    Pkg["Goexit"] = runtime.Goexit
    Pkg["GoroutineProfile"] = runtime.GoroutineProfile
    Pkg["Gosched"] = runtime.Gosched
    Pkg["KeepAlive"] = runtime.KeepAlive
    Pkg["LockOSThread"] = runtime.LockOSThread
    Pkg["MemProfile"] = runtime.MemProfile
    Pkg["MemProfileRate"] = runtime.MemProfileRate
    Ctor["MemProfileRecord"] = GijitShadow_NewStruct_MemProfileRecord
    Ctor["MemStats"] = GijitShadow_NewStruct_MemStats
    Pkg["MutexProfile"] = runtime.MutexProfile
    Pkg["NumCPU"] = runtime.NumCPU
    Pkg["NumCgoCall"] = runtime.NumCgoCall
    Pkg["NumGoroutine"] = runtime.NumGoroutine
    Pkg["ReadMemStats"] = runtime.ReadMemStats
    Pkg["ReadTrace"] = runtime.ReadTrace
    Pkg["SetBlockProfileRate"] = runtime.SetBlockProfileRate
    Pkg["SetCPUProfileRate"] = runtime.SetCPUProfileRate
    Pkg["SetCgoTraceback"] = runtime.SetCgoTraceback
    Pkg["SetFinalizer"] = runtime.SetFinalizer
    Pkg["SetMutexProfileFraction"] = runtime.SetMutexProfileFraction
    Pkg["Stack"] = runtime.Stack
    Ctor["StackRecord"] = GijitShadow_NewStruct_StackRecord
    Pkg["StartTrace"] = runtime.StartTrace
    Pkg["StopTrace"] = runtime.StopTrace
    Pkg["ThreadCreateProfile"] = runtime.ThreadCreateProfile
    Ctor["TypeAssertionError"] = GijitShadow_NewStruct_TypeAssertionError
    Pkg["UnlockOSThread"] = runtime.UnlockOSThread
    Pkg["Version"] = runtime.Version

}
func GijitShadow_NewStruct_BlockProfileRecord(src *runtime.BlockProfileRecord) *runtime.BlockProfileRecord {
    if src == nil {
	   return &runtime.BlockProfileRecord{}
    }
    a := *src
    return &a
}


func GijitShadow_InterfaceConvertTo2_Error(x interface{}) (y runtime.Error, b bool) {
	y, b = x.(runtime.Error)
	return
}

func GijitShadow_InterfaceConvertTo1_Error(x interface{}) runtime.Error {
	return x.(runtime.Error)
}


func GijitShadow_NewStruct_Frame(src *runtime.Frame) *runtime.Frame {
    if src == nil {
	   return &runtime.Frame{}
    }
    a := *src
    return &a
}


func GijitShadow_NewStruct_Frames(src *runtime.Frames) *runtime.Frames {
    if src == nil {
	   return &runtime.Frames{}
    }
    a := *src
    return &a
}


func GijitShadow_NewStruct_Func(src *runtime.Func) *runtime.Func {
    if src == nil {
	   return &runtime.Func{}
    }
    a := *src
    return &a
}


func GijitShadow_NewStruct_MemProfileRecord(src *runtime.MemProfileRecord) *runtime.MemProfileRecord {
    if src == nil {
	   return &runtime.MemProfileRecord{}
    }
    a := *src
    return &a
}


func GijitShadow_NewStruct_MemStats(src *runtime.MemStats) *runtime.MemStats {
    if src == nil {
	   return &runtime.MemStats{}
    }
    a := *src
    return &a
}


func GijitShadow_NewStruct_StackRecord(src *runtime.StackRecord) *runtime.StackRecord {
    if src == nil {
	   return &runtime.StackRecord{}
    }
    a := *src
    return &a
}


func GijitShadow_NewStruct_TypeAssertionError(src *runtime.TypeAssertionError) *runtime.TypeAssertionError {
    if src == nil {
	   return &runtime.TypeAssertionError{}
    }
    a := *src
    return &a
}



 func InitLua() string {
  return `
__type__.runtime ={};

-----------------
-- struct BlockProfileRecord
-----------------

__type__.runtime.BlockProfileRecord = {
 __name = "native_Go_struct_type_wrapper",
 __native_type = "BlockProfileRecord",
 __str = "BlockProfileRecord",
 exported = true,
 __call = function(t, src)
   return __ctor__runtime.BlockProfileRecord(src)
 end,
};
setmetatable(__type__.runtime.BlockProfileRecord, __type__.runtime.BlockProfileRecord);


-----------------
-- struct Frame
-----------------

__type__.runtime.Frame = {
 __name = "native_Go_struct_type_wrapper",
 __native_type = "Frame",
 __str = "Frame",
 exported = true,
 __call = function(t, src)
   return __ctor__runtime.Frame(src)
 end,
};
setmetatable(__type__.runtime.Frame, __type__.runtime.Frame);


-----------------
-- struct Frames
-----------------

__type__.runtime.Frames = {
 __name = "native_Go_struct_type_wrapper",
 __native_type = "Frames",
 __str = "Frames",
 exported = true,
 __call = function(t, src)
   return __ctor__runtime.Frames(src)
 end,
};
setmetatable(__type__.runtime.Frames, __type__.runtime.Frames);


-----------------
-- struct Func
-----------------

__type__.runtime.Func = {
 __name = "native_Go_struct_type_wrapper",
 __native_type = "Func",
 __str = "Func",
 exported = true,
 __call = function(t, src)
   return __ctor__runtime.Func(src)
 end,
};
setmetatable(__type__.runtime.Func, __type__.runtime.Func);


-----------------
-- struct MemProfileRecord
-----------------

__type__.runtime.MemProfileRecord = {
 __name = "native_Go_struct_type_wrapper",
 __native_type = "MemProfileRecord",
 __str = "MemProfileRecord",
 exported = true,
 __call = function(t, src)
   return __ctor__runtime.MemProfileRecord(src)
 end,
};
setmetatable(__type__.runtime.MemProfileRecord, __type__.runtime.MemProfileRecord);


-----------------
-- struct MemStats
-----------------

__type__.runtime.MemStats = {
 __name = "native_Go_struct_type_wrapper",
 __native_type = "MemStats",
 __str = "MemStats",
 exported = true,
 __call = function(t, src)
   return __ctor__runtime.MemStats(src)
 end,
};
setmetatable(__type__.runtime.MemStats, __type__.runtime.MemStats);


-----------------
-- struct StackRecord
-----------------

__type__.runtime.StackRecord = {
 __name = "native_Go_struct_type_wrapper",
 __native_type = "StackRecord",
 __str = "StackRecord",
 exported = true,
 __call = function(t, src)
   return __ctor__runtime.StackRecord(src)
 end,
};
setmetatable(__type__.runtime.StackRecord, __type__.runtime.StackRecord);


-----------------
-- struct TypeAssertionError
-----------------

__type__.runtime.TypeAssertionError = {
 __name = "native_Go_struct_type_wrapper",
 __native_type = "TypeAssertionError",
 __str = "TypeAssertionError",
 exported = true,
 __call = function(t, src)
   return __ctor__runtime.TypeAssertionError(src)
 end,
};
setmetatable(__type__.runtime.TypeAssertionError, __type__.runtime.TypeAssertionError);


`}