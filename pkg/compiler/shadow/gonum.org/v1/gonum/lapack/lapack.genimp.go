package shadow_gonum.org/v1/gonum/lapack

import "gonum.org/v1/gonum/lapack"

var Pkg = make(map[string]interface{})
func init() {
    Pkg["Complex128"] = GijitShadow_InterfaceConvertTo2_Complex128
    Pkg["Float64"] = GijitShadow_InterfaceConvertTo2_Float64
    Pkg["None"] = lapack.None

}
func GijitShadow_InterfaceConvertTo2_Complex128(x interface{}) (y lapack.Complex128, b bool) {
	y, b = x.(lapack.Complex128)
	return
}

func GijitShadow_InterfaceConvertTo1_Complex128(x interface{}) lapack.Complex128 {
	return x.(lapack.Complex128)
}


func GijitShadow_InterfaceConvertTo2_Float64(x interface{}) (y lapack.Float64, b bool) {
	y, b = x.(lapack.Float64)
	return
}

func GijitShadow_InterfaceConvertTo1_Float64(x interface{}) lapack.Float64 {
	return x.(lapack.Float64)
}

