-- int64 uint64 helpers

local ffi = require("ffi")

ffi.cdef[[
long long int atoll(const char *nptr);
]]
__atoll=ffi.C.atoll

-- assume 64-bit int and uint
int = ffi.typeof(0LL)
uint=ffi.typeof(0ULL)

int64=ffi.typeof("int64_t")
uint64=ffi.typeof("uint64_t")

int32 = ffi.typeof("int32_t")
uint32 = ffi.typeof("uint32_t")

int16 = ffi.typeof("int16_t")
uint16 = ffi.typeof("uint16_t")

int8 = ffi.typeof("int8_t")
uint8 = ffi.typeof("uint8_t")
byte = uint8

float64 = ffi.typeof("double")
float32 = ffi.typeof("float")

-- to display floats, use: tonumber() to convert to float64 that lua can print.

--MinInt64: -9223372036854775808
--MaxInt64: 9223372036854775807

-- a=int(-9223372036854775808LL)

-- to use cdata as hash keys... tostring() to make them strings first.
