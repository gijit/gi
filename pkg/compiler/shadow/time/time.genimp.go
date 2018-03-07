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
    Pkg["LoadLocationFromTZData"] = time.LoadLocationFromTZData
    Pkg["Local"] = time.Local
    Ctor["Location"] = GijitShadow_NewStruct_Location
    Pkg["NewTicker"] = time.NewTicker
    Pkg["NewTimer"] = time.NewTimer
    Pkg["Now"] = time.Now
    Pkg["Parse"] = time.Parse
    Pkg["ParseDuration"] = time.ParseDuration
    Ctor["ParseError"] = GijitShadow_NewStruct_ParseError
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
    Ctor["Ticker"] = GijitShadow_NewStruct_Ticker
    Ctor["Time"] = GijitShadow_NewStruct_Time
    Ctor["Timer"] = GijitShadow_NewStruct_Timer
    Pkg["UTC"] = time.UTC
    Pkg["Unix"] = time.Unix
    Pkg["UnixDate"] = time.UnixDate
    Pkg["Until"] = time.Until

}
func GijitShadow_NewStruct_Location(src *time.Location) *time.Location {
    if src == nil {
	   return &time.Location{}
    }
    a := *src
    return &a
}


func GijitShadow_NewStruct_ParseError(src *time.ParseError) *time.ParseError {
    if src == nil {
	   return &time.ParseError{}
    }
    a := *src
    return &a
}


func GijitShadow_NewStruct_Ticker(src *time.Ticker) *time.Ticker {
    if src == nil {
	   return &time.Ticker{}
    }
    a := *src
    return &a
}


func GijitShadow_NewStruct_Time(src *time.Time) *time.Time {
    if src == nil {
	   return &time.Time{}
    }
    a := *src
    return &a
}


func GijitShadow_NewStruct_Timer(src *time.Timer) *time.Timer {
    if src == nil {
	   return &time.Timer{}
    }
    a := *src
    return &a
}

func InitLua() string {
  return `
__type__.time ={};

-----------------
-- struct Location
-----------------

__type__.time.Location = {
 __name = "native_Go_struct_type_wrapper",
 __native_type = "Location",
 __call = function(t, src)
   return __ctor__time.Location(src)
 end,
};
setmetatable(__type__.time.Location, __type__.time.Location);


-----------------
-- struct ParseError
-----------------

__type__.time.ParseError = {
 __name = "native_Go_struct_type_wrapper",
 __native_type = "ParseError",
 __call = function(t, src)
   return __ctor__time.ParseError(src)
 end,
};
setmetatable(__type__.time.ParseError, __type__.time.ParseError);


-----------------
-- struct Ticker
-----------------

__type__.time.Ticker = {
 __name = "native_Go_struct_type_wrapper",
 __native_type = "Ticker",
 __call = function(t, src)
   return __ctor__time.Ticker(src)
 end,
};
setmetatable(__type__.time.Ticker, __type__.time.Ticker);


-----------------
-- struct Time
-----------------

__type__.time.Time = {
 __name = "native_Go_struct_type_wrapper",
 __native_type = "Time",
 __call = function(t, src)
   return __ctor__time.Time(src)
 end,
};
setmetatable(__type__.time.Time, __type__.time.Time);


-----------------
-- struct Timer
-----------------

__type__.time.Timer = {
 __name = "native_Go_struct_type_wrapper",
 __native_type = "Timer",
 __call = function(t, src)
   return __ctor__time.Timer(src)
 end,
};
setmetatable(__type__.time.Timer, __type__.time.Timer);


`}