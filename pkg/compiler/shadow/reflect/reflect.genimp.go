package shadow_reflect

import "reflect"

var Pkg = make(map[string]interface{})
func init() {
    Pkg["Append"] = reflect.Append
    Pkg["AppendSlice"] = reflect.AppendSlice
    Pkg["ArrayOf"] = reflect.ArrayOf
    Pkg["ChanOf"] = reflect.ChanOf
    Pkg["Copy"] = reflect.Copy
    Pkg["DeepEqual"] = reflect.DeepEqual
    Pkg["FuncOf"] = reflect.FuncOf
    Pkg["Indirect"] = reflect.Indirect
    Pkg["MakeChan"] = reflect.MakeChan
    Pkg["MakeFunc"] = reflect.MakeFunc
    Pkg["MakeMap"] = reflect.MakeMap
    Pkg["MakeMapWithSize"] = reflect.MakeMapWithSize
    Pkg["MakeSlice"] = reflect.MakeSlice
    Pkg["MapOf"] = reflect.MapOf
    Pkg["New"] = reflect.New
    Pkg["NewAt"] = reflect.NewAt
    Pkg["PtrTo"] = reflect.PtrTo
    Pkg["Select"] = reflect.Select
    Pkg["SliceOf"] = reflect.SliceOf
    Pkg["StructOf"] = reflect.StructOf
    Pkg["Swapper"] = reflect.Swapper
    Pkg["Type"] = GijitShadow_InterfaceConvertTo2_Type
    Pkg["TypeOf"] = reflect.TypeOf
    Pkg["ValueOf"] = reflect.ValueOf
    Pkg["Zero"] = reflect.Zero

}
func GijitShadow_NewStruct_Method() *reflect.Method {
	return &reflect.Method{}
}


func GijitShadow_NewStruct_SelectCase() *reflect.SelectCase {
	return &reflect.SelectCase{}
}


func GijitShadow_NewStruct_SliceHeader() *reflect.SliceHeader {
	return &reflect.SliceHeader{}
}


func GijitShadow_NewStruct_StringHeader() *reflect.StringHeader {
	return &reflect.StringHeader{}
}


func GijitShadow_NewStruct_StructField() *reflect.StructField {
	return &reflect.StructField{}
}


func GijitShadow_InterfaceConvertTo2_Type(x interface{}) (y reflect.Type, b bool) {
	y, b = x.(reflect.Type)
	return
}

func GijitShadow_InterfaceConvertTo1_Type(x interface{}) reflect.Type {
	return x.(reflect.Type)
}


func GijitShadow_NewStruct_Value() *reflect.Value {
	return &reflect.Value{}
}


func GijitShadow_NewStruct_ValueError() *reflect.ValueError {
	return &reflect.ValueError{}
}

