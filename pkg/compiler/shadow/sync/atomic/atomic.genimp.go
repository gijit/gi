package shadow_atomic

import "sync/atomic"

var Pkg = make(map[string]interface{})
var Ctor = make(map[string]interface{})

func init() {
    Pkg["AddInt32"] = atomic.AddInt32
    Pkg["AddInt64"] = atomic.AddInt64
    Pkg["AddUint32"] = atomic.AddUint32
    Pkg["AddUint64"] = atomic.AddUint64
    Pkg["AddUintptr"] = atomic.AddUintptr
    Pkg["CompareAndSwapInt32"] = atomic.CompareAndSwapInt32
    Pkg["CompareAndSwapInt64"] = atomic.CompareAndSwapInt64
    Pkg["CompareAndSwapPointer"] = atomic.CompareAndSwapPointer
    Pkg["CompareAndSwapUint32"] = atomic.CompareAndSwapUint32
    Pkg["CompareAndSwapUint64"] = atomic.CompareAndSwapUint64
    Pkg["CompareAndSwapUintptr"] = atomic.CompareAndSwapUintptr
    Pkg["LoadInt32"] = atomic.LoadInt32
    Pkg["LoadInt64"] = atomic.LoadInt64
    Pkg["LoadPointer"] = atomic.LoadPointer
    Pkg["LoadUint32"] = atomic.LoadUint32
    Pkg["LoadUint64"] = atomic.LoadUint64
    Pkg["LoadUintptr"] = atomic.LoadUintptr
    Pkg["StoreInt32"] = atomic.StoreInt32
    Pkg["StoreInt64"] = atomic.StoreInt64
    Pkg["StorePointer"] = atomic.StorePointer
    Pkg["StoreUint32"] = atomic.StoreUint32
    Pkg["StoreUint64"] = atomic.StoreUint64
    Pkg["StoreUintptr"] = atomic.StoreUintptr
    Pkg["SwapInt32"] = atomic.SwapInt32
    Pkg["SwapInt64"] = atomic.SwapInt64
    Pkg["SwapPointer"] = atomic.SwapPointer
    Pkg["SwapUint32"] = atomic.SwapUint32
    Pkg["SwapUint64"] = atomic.SwapUint64
    Pkg["SwapUintptr"] = atomic.SwapUintptr
    Ctor["Value"] = GijitShadow_NewStruct_Value

}
func GijitShadow_NewStruct_Value(src *atomic.Value) *atomic.Value {
    if src == nil {
	   return &atomic.Value{}
    }
    a := *src
    return &a
}



 func InitLua() string {
  return `
__type__.atomic ={};

-----------------
-- struct Value
-----------------

__type__.atomic.Value = {
 __name = "native_Go_struct_type_wrapper",
 __native_type = "Value",
 __call = function(t, src)
   return __ctor__atomic.Value(src)
 end,
};
setmetatable(__type__.atomic.Value, __type__.atomic.Value);


`}