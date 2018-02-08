package shadow_gonum.org/v1/gonum/blas

import "gonum.org/v1/gonum/blas"

var Pkg = make(map[string]interface{})
func init() {
    Pkg["Complex128"] = GijitShadow_InterfaceConvertTo2_Complex128
    Pkg["Complex128Level1"] = GijitShadow_InterfaceConvertTo2_Complex128Level1
    Pkg["Complex128Level2"] = GijitShadow_InterfaceConvertTo2_Complex128Level2
    Pkg["Complex128Level3"] = GijitShadow_InterfaceConvertTo2_Complex128Level3
    Pkg["Complex64"] = GijitShadow_InterfaceConvertTo2_Complex64
    Pkg["Complex64Level1"] = GijitShadow_InterfaceConvertTo2_Complex64Level1
    Pkg["Complex64Level2"] = GijitShadow_InterfaceConvertTo2_Complex64Level2
    Pkg["Complex64Level3"] = GijitShadow_InterfaceConvertTo2_Complex64Level3
    Pkg["Float32"] = GijitShadow_InterfaceConvertTo2_Float32
    Pkg["Float32Level1"] = GijitShadow_InterfaceConvertTo2_Float32Level1
    Pkg["Float32Level2"] = GijitShadow_InterfaceConvertTo2_Float32Level2
    Pkg["Float32Level3"] = GijitShadow_InterfaceConvertTo2_Float32Level3
    Pkg["Float64"] = GijitShadow_InterfaceConvertTo2_Float64
    Pkg["Float64Level1"] = GijitShadow_InterfaceConvertTo2_Float64Level1
    Pkg["Float64Level2"] = GijitShadow_InterfaceConvertTo2_Float64Level2
    Pkg["Float64Level3"] = GijitShadow_InterfaceConvertTo2_Float64Level3

}
func GijitShadow_InterfaceConvertTo2_Complex128(x interface{}) (y blas.Complex128, b bool) {
	y, b = x.(blas.Complex128)
	return
}

func GijitShadow_InterfaceConvertTo1_Complex128(x interface{}) blas.Complex128 {
	return x.(blas.Complex128)
}


func GijitShadow_InterfaceConvertTo2_Complex128Level1(x interface{}) (y blas.Complex128Level1, b bool) {
	y, b = x.(blas.Complex128Level1)
	return
}

func GijitShadow_InterfaceConvertTo1_Complex128Level1(x interface{}) blas.Complex128Level1 {
	return x.(blas.Complex128Level1)
}


func GijitShadow_InterfaceConvertTo2_Complex128Level2(x interface{}) (y blas.Complex128Level2, b bool) {
	y, b = x.(blas.Complex128Level2)
	return
}

func GijitShadow_InterfaceConvertTo1_Complex128Level2(x interface{}) blas.Complex128Level2 {
	return x.(blas.Complex128Level2)
}


func GijitShadow_InterfaceConvertTo2_Complex128Level3(x interface{}) (y blas.Complex128Level3, b bool) {
	y, b = x.(blas.Complex128Level3)
	return
}

func GijitShadow_InterfaceConvertTo1_Complex128Level3(x interface{}) blas.Complex128Level3 {
	return x.(blas.Complex128Level3)
}


func GijitShadow_InterfaceConvertTo2_Complex64(x interface{}) (y blas.Complex64, b bool) {
	y, b = x.(blas.Complex64)
	return
}

func GijitShadow_InterfaceConvertTo1_Complex64(x interface{}) blas.Complex64 {
	return x.(blas.Complex64)
}


func GijitShadow_InterfaceConvertTo2_Complex64Level1(x interface{}) (y blas.Complex64Level1, b bool) {
	y, b = x.(blas.Complex64Level1)
	return
}

func GijitShadow_InterfaceConvertTo1_Complex64Level1(x interface{}) blas.Complex64Level1 {
	return x.(blas.Complex64Level1)
}


func GijitShadow_InterfaceConvertTo2_Complex64Level2(x interface{}) (y blas.Complex64Level2, b bool) {
	y, b = x.(blas.Complex64Level2)
	return
}

func GijitShadow_InterfaceConvertTo1_Complex64Level2(x interface{}) blas.Complex64Level2 {
	return x.(blas.Complex64Level2)
}


func GijitShadow_InterfaceConvertTo2_Complex64Level3(x interface{}) (y blas.Complex64Level3, b bool) {
	y, b = x.(blas.Complex64Level3)
	return
}

func GijitShadow_InterfaceConvertTo1_Complex64Level3(x interface{}) blas.Complex64Level3 {
	return x.(blas.Complex64Level3)
}


func GijitShadow_NewStruct_DrotmParams() *blas.DrotmParams {
	return &blas.DrotmParams{}
}


func GijitShadow_InterfaceConvertTo2_Float32(x interface{}) (y blas.Float32, b bool) {
	y, b = x.(blas.Float32)
	return
}

func GijitShadow_InterfaceConvertTo1_Float32(x interface{}) blas.Float32 {
	return x.(blas.Float32)
}


func GijitShadow_InterfaceConvertTo2_Float32Level1(x interface{}) (y blas.Float32Level1, b bool) {
	y, b = x.(blas.Float32Level1)
	return
}

func GijitShadow_InterfaceConvertTo1_Float32Level1(x interface{}) blas.Float32Level1 {
	return x.(blas.Float32Level1)
}


func GijitShadow_InterfaceConvertTo2_Float32Level2(x interface{}) (y blas.Float32Level2, b bool) {
	y, b = x.(blas.Float32Level2)
	return
}

func GijitShadow_InterfaceConvertTo1_Float32Level2(x interface{}) blas.Float32Level2 {
	return x.(blas.Float32Level2)
}


func GijitShadow_InterfaceConvertTo2_Float32Level3(x interface{}) (y blas.Float32Level3, b bool) {
	y, b = x.(blas.Float32Level3)
	return
}

func GijitShadow_InterfaceConvertTo1_Float32Level3(x interface{}) blas.Float32Level3 {
	return x.(blas.Float32Level3)
}


func GijitShadow_InterfaceConvertTo2_Float64(x interface{}) (y blas.Float64, b bool) {
	y, b = x.(blas.Float64)
	return
}

func GijitShadow_InterfaceConvertTo1_Float64(x interface{}) blas.Float64 {
	return x.(blas.Float64)
}


func GijitShadow_InterfaceConvertTo2_Float64Level1(x interface{}) (y blas.Float64Level1, b bool) {
	y, b = x.(blas.Float64Level1)
	return
}

func GijitShadow_InterfaceConvertTo1_Float64Level1(x interface{}) blas.Float64Level1 {
	return x.(blas.Float64Level1)
}


func GijitShadow_InterfaceConvertTo2_Float64Level2(x interface{}) (y blas.Float64Level2, b bool) {
	y, b = x.(blas.Float64Level2)
	return
}

func GijitShadow_InterfaceConvertTo1_Float64Level2(x interface{}) blas.Float64Level2 {
	return x.(blas.Float64Level2)
}


func GijitShadow_InterfaceConvertTo2_Float64Level3(x interface{}) (y blas.Float64Level3, b bool) {
	y, b = x.(blas.Float64Level3)
	return
}

func GijitShadow_InterfaceConvertTo1_Float64Level3(x interface{}) blas.Float64Level3 {
	return x.(blas.Float64Level3)
}


func GijitShadow_NewStruct_SrotmParams() *blas.SrotmParams {
	return &blas.SrotmParams{}
}

