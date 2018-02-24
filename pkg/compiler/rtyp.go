package compiler

import (
	golua "github.com/glycerine/golua/lua"
	"github.com/glycerine/luar"
	"reflect"
	"unsafe"
)

// generate the basic reflect types,
// for use in tsys.lua when constructing
// unnamed channel types to pass values on.
func registerBasicReflectTypes(vm *golua.State) {

	m := make(luar.Map)

	m["__kindBool"] = reflect.TypeOf(true)
	m["__kindInt"] = reflect.TypeOf(int(0))
	m["__kindInt8"] = reflect.TypeOf(int8(0))

	m["__kindInt16"] = reflect.TypeOf(int16(0))
	m["__kindInt32"] = reflect.TypeOf(int32(0))
	m["__kindInt64"] = reflect.TypeOf(int64(0))

	m["__kindUint"] = reflect.TypeOf(uint(0))
	m["__kindUint8"] = reflect.TypeOf(uint8(0))
	m["__kindUint16"] = reflect.TypeOf(uint16(0))

	m["__kindUint32"] = reflect.TypeOf(uint32(0))
	m["__kindUint64"] = reflect.TypeOf(uint64(0))
	m["__kindUintptr"] = reflect.TypeOf(uintptr(0))

	m["__kindUnsafePointer"] = reflect.TypeOf(unsafe.Pointer(&m))
	m["__kindString"] = reflect.TypeOf("")
	m["__kindFloat32"] = reflect.TypeOf(float32(0))

	m["__kindFloat64"] = reflect.TypeOf(float64(0))
	m["__kindComplex64"] = reflect.TypeOf(complex64(0))
	m["__kindComplex128"] = reflect.TypeOf(complex128(0))

	luar.Register(vm, "__rtypbasic", m)
}
