package shadow_reflect

import "reflect"

var Pkg = make(map[string]interface{})
var Ctor = make(map[string]interface{})

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
    Ctor["Method"] = GijitShadow_NewStruct_Method
    Pkg["New"] = reflect.New
    Pkg["NewAt"] = reflect.NewAt
    Pkg["PtrTo"] = reflect.PtrTo
    Pkg["Select"] = reflect.Select
    Ctor["SelectCase"] = GijitShadow_NewStruct_SelectCase
    Ctor["SliceHeader"] = GijitShadow_NewStruct_SliceHeader
    Pkg["SliceOf"] = reflect.SliceOf
    Ctor["StringHeader"] = GijitShadow_NewStruct_StringHeader
    Ctor["StructField"] = GijitShadow_NewStruct_StructField
    Pkg["StructOf"] = reflect.StructOf
    Pkg["Swapper"] = reflect.Swapper
    Pkg["Type"] = GijitShadow_InterfaceConvertTo2_Type
    Pkg["TypeOf"] = reflect.TypeOf
    Ctor["Value"] = GijitShadow_NewStruct_Value
    Ctor["ValueError"] = GijitShadow_NewStruct_ValueError
    Pkg["ValueOf"] = reflect.ValueOf
    Pkg["Zero"] = reflect.Zero

}
func GijitShadow_NewStruct_Method(src *reflect.Method) *reflect.Method {
    if src == nil {
	   return &reflect.Method{}
    }
    a := *src
    return &a
}


func GijitShadow_NewStruct_SelectCase(src *reflect.SelectCase) *reflect.SelectCase {
    if src == nil {
	   return &reflect.SelectCase{}
    }
    a := *src
    return &a
}


func GijitShadow_NewStruct_SliceHeader(src *reflect.SliceHeader) *reflect.SliceHeader {
    if src == nil {
	   return &reflect.SliceHeader{}
    }
    a := *src
    return &a
}


func GijitShadow_NewStruct_StringHeader(src *reflect.StringHeader) *reflect.StringHeader {
    if src == nil {
	   return &reflect.StringHeader{}
    }
    a := *src
    return &a
}


func GijitShadow_NewStruct_StructField(src *reflect.StructField) *reflect.StructField {
    if src == nil {
	   return &reflect.StructField{}
    }
    a := *src
    return &a
}


func GijitShadow_InterfaceConvertTo2_Type(x interface{}) (y reflect.Type, b bool) {
	y, b = x.(reflect.Type)
	return
}

func GijitShadow_InterfaceConvertTo1_Type(x interface{}) reflect.Type {
	return x.(reflect.Type)
}


func GijitShadow_NewStruct_Value(src *reflect.Value) *reflect.Value {
    if src == nil {
	   return &reflect.Value{}
    }
    a := *src
    return &a
}


func GijitShadow_NewStruct_ValueError(src *reflect.ValueError) *reflect.ValueError {
    if src == nil {
	   return &reflect.ValueError{}
    }
    a := *src
    return &a
}



 func InitLua() string {
  return `
__type__.reflect ={};

-----------------
-- struct Method
-----------------

__type__.reflect.Method = {
 id=0,
 __name = "native_Go_struct_type_wrapper",
 __native_type = "Method",
 __str = "Method",
 exported = true,
 __call = function(t, src)
   return __ctor__reflect.Method(src)
 end,
};
setmetatable(__type__.reflect.Method, __type__.reflect.Method);


-----------------
-- struct SelectCase
-----------------

__type__.reflect.SelectCase = {
 id=0,
 __name = "native_Go_struct_type_wrapper",
 __native_type = "SelectCase",
 __str = "SelectCase",
 exported = true,
 __call = function(t, src)
   return __ctor__reflect.SelectCase(src)
 end,
};
setmetatable(__type__.reflect.SelectCase, __type__.reflect.SelectCase);


-----------------
-- struct SliceHeader
-----------------

__type__.reflect.SliceHeader = {
 id=0,
 __name = "native_Go_struct_type_wrapper",
 __native_type = "SliceHeader",
 __str = "SliceHeader",
 exported = true,
 __call = function(t, src)
   return __ctor__reflect.SliceHeader(src)
 end,
};
setmetatable(__type__.reflect.SliceHeader, __type__.reflect.SliceHeader);


-----------------
-- struct StringHeader
-----------------

__type__.reflect.StringHeader = {
 id=0,
 __name = "native_Go_struct_type_wrapper",
 __native_type = "StringHeader",
 __str = "StringHeader",
 exported = true,
 __call = function(t, src)
   return __ctor__reflect.StringHeader(src)
 end,
};
setmetatable(__type__.reflect.StringHeader, __type__.reflect.StringHeader);


-----------------
-- struct StructField
-----------------

__type__.reflect.StructField = {
 id=0,
 __name = "native_Go_struct_type_wrapper",
 __native_type = "StructField",
 __str = "StructField",
 exported = true,
 __call = function(t, src)
   return __ctor__reflect.StructField(src)
 end,
};
setmetatable(__type__.reflect.StructField, __type__.reflect.StructField);


-----------------
-- struct Value
-----------------

__type__.reflect.Value = {
 id=0,
 __name = "native_Go_struct_type_wrapper",
 __native_type = "Value",
 __str = "Value",
 exported = true,
 __call = function(t, src)
   return __ctor__reflect.Value(src)
 end,
};
setmetatable(__type__.reflect.Value, __type__.reflect.Value);


-----------------
-- struct ValueError
-----------------

__type__.reflect.ValueError = {
 id=0,
 __name = "native_Go_struct_type_wrapper",
 __native_type = "ValueError",
 __str = "ValueError",
 exported = true,
 __call = function(t, src)
   return __ctor__reflect.ValueError(src)
 end,
};
setmetatable(__type__.reflect.ValueError, __type__.reflect.ValueError);


`}