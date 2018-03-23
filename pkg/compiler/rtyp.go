package compiler

import (
	golua "github.com/glycerine/golua/lua"
	"github.com/glycerine/luar"
	"reflect"
	"unsafe"
)

// must be kept in sync with tsys.lua
const __kindUnknown = -1
const __kindBool = 1
const __kindInt = 2
const __kindInt8 = 3
const __kindInt16 = 4
const __kindInt32 = 5
const __kindInt64 = 6
const __kindUint = 7
const __kindUint8 = 8
const __kindUint16 = 9
const __kindUint32 = 10
const __kindUint64 = 11
const __kindUintptr = 12
const __kindFloat32 = 13
const __kindFloat64 = 14
const __kindComplex64 = 15
const __kindComplex128 = 16
const __kindArray = 17
const __kindChan = 18
const __kindFunc = 19
const __kindInterface = 20
const __kindMap = 21
const __kindPtr = 22
const __kindSlice = 23
const __kindString = 24
const __kindStruct = 25
const __kindUnsafePointer = 26

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

	luar.Register(vm, "__rtyp", m)

}
