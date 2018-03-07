package shadow_time

import "time"

var Pkg = make(map[string]interface{})
var Ctor = make(map[string]interface{})

func init() {
	Pkg["ANSIC"] = time.ANSIC
	Pkg["After"] = time.After
	Pkg["AfterFunc"] = time.AfterFunc
	Pkg["Date"] = time.Date
	Pkg["FixedZone"] = time.FixedZone
	Pkg["Kitchen"] = time.Kitchen
	Pkg["LoadLocation"] = time.LoadLocation
	Pkg["Local"] = time.Local
	Pkg["NewTicker"] = time.NewTicker
	Pkg["NewTimer"] = time.NewTimer
	Pkg["Now"] = time.Now
	Pkg["Parse"] = time.Parse
	Pkg["ParseDuration"] = time.ParseDuration
	Pkg["ParseInLocation"] = time.ParseInLocation
	Pkg["RFC1123"] = time.RFC1123
	Pkg["RFC1123Z"] = time.RFC1123Z
	Pkg["RFC3339"] = time.RFC3339
	Pkg["RFC3339Nano"] = time.RFC3339Nano
	Pkg["RFC822"] = time.RFC822
	Pkg["RFC822Z"] = time.RFC822Z
	Pkg["RFC850"] = time.RFC850
	Pkg["RubyDate"] = time.RubyDate
	Pkg["Since"] = time.Since
	Pkg["Sleep"] = time.Sleep
	Pkg["Stamp"] = time.Stamp
	Pkg["StampMicro"] = time.StampMicro
	Pkg["StampMilli"] = time.StampMilli
	Pkg["StampNano"] = time.StampNano
	Pkg["Tick"] = time.Tick
	Pkg["UTC"] = time.UTC
	Pkg["Unix"] = time.Unix
	Pkg["UnixDate"] = time.UnixDate
	Pkg["Until"] = time.Until

	// __ctor__time.Time
	Ctor["Time"] = GijitShadow_CtorStruct_Time
}
func GijitShadow_NewStruct_Location() *time.Location {
	return &time.Location{}
}

func GijitShadow_NewStruct_ParseError() *time.ParseError {
	return &time.ParseError{}
}

func GijitShadow_NewStruct_Ticker() *time.Ticker {
	return &time.Ticker{}
}

func GijitShadow_NewStruct_Time() *time.Time {
	return &time.Time{}
}

func GijitShadow_NewStruct_Timer() *time.Timer {
	return &time.Timer{}
}

var Types = map[string]interface{}{"Time": GijitShadow_NewStruct_Time}

func GijitShadow_CtorStruct_Time(src *time.Time) *time.Time {
	if src == nil {
		return &time.Time{}
	}
	a := *src
	return &a
}

func InitLua() []byte {
	return []byte(`
__type__.time ={};

-----------------
-- struct Time 
-----------------

__type__.time.Time = {
 __name = "native_Go_struct_type_wrapper",
 __native_type = "Time",
 __call = function(t, src)
   -- With nil src, return zero value for type.
   -- Else act as copy constructor, copying from src.
   return __ctor__time.Time(src)
 end,
};
setmetatable(__type__.time.Time, __type__.time.Time);

`)
}
