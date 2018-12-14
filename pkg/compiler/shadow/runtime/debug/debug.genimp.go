package shadow_debug

import "runtime/debug"

var Pkg = make(map[string]interface{})
var Ctor = make(map[string]interface{})

func init() {
    Pkg["FreeOSMemory"] = debug.FreeOSMemory
    Ctor["GCStats"] = GijitShadow_NewStruct_GCStats
    Pkg["PrintStack"] = debug.PrintStack
    Pkg["ReadGCStats"] = debug.ReadGCStats
    Pkg["SetGCPercent"] = debug.SetGCPercent
    Pkg["SetMaxStack"] = debug.SetMaxStack
    Pkg["SetMaxThreads"] = debug.SetMaxThreads
    Pkg["SetPanicOnFault"] = debug.SetPanicOnFault
    Pkg["SetTraceback"] = debug.SetTraceback
    Pkg["Stack"] = debug.Stack
    Pkg["WriteHeapDump"] = debug.WriteHeapDump

}
func GijitShadow_NewStruct_GCStats(src *debug.GCStats) *debug.GCStats {
    if src == nil {
	   return &debug.GCStats{}
    }
    a := *src
    return &a
}



 func InitLua() string {
  return `
__type__.debug ={};

-----------------
-- struct GCStats
-----------------

__type__.debug.GCStats = {
 id=0,
 __name = "native_Go_struct_type_wrapper",
 __native_type = "GCStats",
 __str = "GCStats",
 exported = true,
 __call = function(t, src)
   return __ctor__debug.GCStats(src)
 end,
};
setmetatable(__type__.debug.GCStats, __type__.debug.GCStats);


`}