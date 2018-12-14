package shadow_rand

import "math/rand"

var Pkg = make(map[string]interface{})
var Ctor = make(map[string]interface{})

func init() {
    Pkg["ExpFloat64"] = rand.ExpFloat64
    Pkg["Float32"] = rand.Float32
    Pkg["Float64"] = rand.Float64
    Pkg["Int"] = rand.Int
    Pkg["Int31"] = rand.Int31
    Pkg["Int31n"] = rand.Int31n
    Pkg["Int63"] = rand.Int63
    Pkg["Int63n"] = rand.Int63n
    Pkg["Intn"] = rand.Intn
    Pkg["New"] = rand.New
    Pkg["NewSource"] = rand.NewSource
    Pkg["NewZipf"] = rand.NewZipf
    Pkg["NormFloat64"] = rand.NormFloat64
    Pkg["Perm"] = rand.Perm
    Ctor["Rand"] = GijitShadow_NewStruct_Rand
    Pkg["Read"] = rand.Read
    Pkg["Seed"] = rand.Seed
    Pkg["Shuffle"] = rand.Shuffle
    Pkg["Source"] = GijitShadow_InterfaceConvertTo2_Source
    Pkg["Source64"] = GijitShadow_InterfaceConvertTo2_Source64
    Pkg["Uint32"] = rand.Uint32
    Pkg["Uint64"] = rand.Uint64
    Ctor["Zipf"] = GijitShadow_NewStruct_Zipf

}
func GijitShadow_NewStruct_Rand(src *rand.Rand) *rand.Rand {
    if src == nil {
	   return &rand.Rand{}
    }
    a := *src
    return &a
}


func GijitShadow_InterfaceConvertTo2_Source(x interface{}) (y rand.Source, b bool) {
	y, b = x.(rand.Source)
	return
}

func GijitShadow_InterfaceConvertTo1_Source(x interface{}) rand.Source {
	return x.(rand.Source)
}


func GijitShadow_InterfaceConvertTo2_Source64(x interface{}) (y rand.Source64, b bool) {
	y, b = x.(rand.Source64)
	return
}

func GijitShadow_InterfaceConvertTo1_Source64(x interface{}) rand.Source64 {
	return x.(rand.Source64)
}


func GijitShadow_NewStruct_Zipf(src *rand.Zipf) *rand.Zipf {
    if src == nil {
	   return &rand.Zipf{}
    }
    a := *src
    return &a
}



 func InitLua() string {
  return `
__type__.rand ={};

-----------------
-- struct Rand
-----------------

__type__.rand.Rand = {
 id=0,
 __name = "native_Go_struct_type_wrapper",
 __native_type = "Rand",
 __str = "Rand",
 exported = true,
 __call = function(t, src)
   return __ctor__rand.Rand(src)
 end,
};
setmetatable(__type__.rand.Rand, __type__.rand.Rand);


-----------------
-- struct Zipf
-----------------

__type__.rand.Zipf = {
 id=0,
 __name = "native_Go_struct_type_wrapper",
 __native_type = "Zipf",
 __str = "Zipf",
 exported = true,
 __call = function(t, src)
   return __ctor__rand.Zipf(src)
 end,
};
setmetatable(__type__.rand.Zipf, __type__.rand.Zipf);


`}