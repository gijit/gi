package shadow_encoding

import "encoding"

var Pkg = make(map[string]interface{})
var Ctor = make(map[string]interface{})

func init() {
    Pkg["BinaryMarshaler"] = GijitShadow_InterfaceConvertTo2_BinaryMarshaler
    Pkg["BinaryUnmarshaler"] = GijitShadow_InterfaceConvertTo2_BinaryUnmarshaler
    Pkg["TextMarshaler"] = GijitShadow_InterfaceConvertTo2_TextMarshaler
    Pkg["TextUnmarshaler"] = GijitShadow_InterfaceConvertTo2_TextUnmarshaler

}
func GijitShadow_InterfaceConvertTo2_BinaryMarshaler(x interface{}) (y encoding.BinaryMarshaler, b bool) {
	y, b = x.(encoding.BinaryMarshaler)
	return
}

func GijitShadow_InterfaceConvertTo1_BinaryMarshaler(x interface{}) encoding.BinaryMarshaler {
	return x.(encoding.BinaryMarshaler)
}


func GijitShadow_InterfaceConvertTo2_BinaryUnmarshaler(x interface{}) (y encoding.BinaryUnmarshaler, b bool) {
	y, b = x.(encoding.BinaryUnmarshaler)
	return
}

func GijitShadow_InterfaceConvertTo1_BinaryUnmarshaler(x interface{}) encoding.BinaryUnmarshaler {
	return x.(encoding.BinaryUnmarshaler)
}


func GijitShadow_InterfaceConvertTo2_TextMarshaler(x interface{}) (y encoding.TextMarshaler, b bool) {
	y, b = x.(encoding.TextMarshaler)
	return
}

func GijitShadow_InterfaceConvertTo1_TextMarshaler(x interface{}) encoding.TextMarshaler {
	return x.(encoding.TextMarshaler)
}


func GijitShadow_InterfaceConvertTo2_TextUnmarshaler(x interface{}) (y encoding.TextUnmarshaler, b bool) {
	y, b = x.(encoding.TextUnmarshaler)
	return
}

func GijitShadow_InterfaceConvertTo1_TextUnmarshaler(x interface{}) encoding.TextUnmarshaler {
	return x.(encoding.TextUnmarshaler)
}



 func InitLua() string {
  return `
__type__.encoding ={};

`}