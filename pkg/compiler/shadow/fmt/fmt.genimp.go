package shadow_fmt

import "fmt"

var Pkg = make(map[string]interface{})
var Ctor = make(map[string]interface{})

func init() {
    Pkg["Errorf"] = fmt.Errorf
    Pkg["Formatter"] = GijitShadow_InterfaceConvertTo2_Formatter
    Pkg["Fprint"] = fmt.Fprint
    Pkg["Fprintf"] = fmt.Fprintf
    Pkg["Fprintln"] = fmt.Fprintln
    Pkg["Fscan"] = fmt.Fscan
    Pkg["Fscanf"] = fmt.Fscanf
    Pkg["Fscanln"] = fmt.Fscanln
    Pkg["GoStringer"] = GijitShadow_InterfaceConvertTo2_GoStringer
    Pkg["Print"] = fmt.Print
    Pkg["Printf"] = fmt.Printf
    Pkg["Println"] = fmt.Println
    Pkg["Scan"] = fmt.Scan
    Pkg["ScanState"] = GijitShadow_InterfaceConvertTo2_ScanState
    Pkg["Scanf"] = fmt.Scanf
    Pkg["Scanln"] = fmt.Scanln
    Pkg["Scanner"] = GijitShadow_InterfaceConvertTo2_Scanner
    Pkg["Sprint"] = fmt.Sprint
    Pkg["Sprintf"] = fmt.Sprintf
    Pkg["Sprintln"] = fmt.Sprintln
    Pkg["Sscan"] = fmt.Sscan
    Pkg["Sscanf"] = fmt.Sscanf
    Pkg["Sscanln"] = fmt.Sscanln
    Pkg["State"] = GijitShadow_InterfaceConvertTo2_State
    Pkg["Stringer"] = GijitShadow_InterfaceConvertTo2_Stringer

}
func GijitShadow_InterfaceConvertTo2_Formatter(x interface{}) (y fmt.Formatter, b bool) {
	y, b = x.(fmt.Formatter)
	return
}

func GijitShadow_InterfaceConvertTo1_Formatter(x interface{}) fmt.Formatter {
	return x.(fmt.Formatter)
}


func GijitShadow_InterfaceConvertTo2_GoStringer(x interface{}) (y fmt.GoStringer, b bool) {
	y, b = x.(fmt.GoStringer)
	return
}

func GijitShadow_InterfaceConvertTo1_GoStringer(x interface{}) fmt.GoStringer {
	return x.(fmt.GoStringer)
}


func GijitShadow_InterfaceConvertTo2_ScanState(x interface{}) (y fmt.ScanState, b bool) {
	y, b = x.(fmt.ScanState)
	return
}

func GijitShadow_InterfaceConvertTo1_ScanState(x interface{}) fmt.ScanState {
	return x.(fmt.ScanState)
}


func GijitShadow_InterfaceConvertTo2_Scanner(x interface{}) (y fmt.Scanner, b bool) {
	y, b = x.(fmt.Scanner)
	return
}

func GijitShadow_InterfaceConvertTo1_Scanner(x interface{}) fmt.Scanner {
	return x.(fmt.Scanner)
}


func GijitShadow_InterfaceConvertTo2_State(x interface{}) (y fmt.State, b bool) {
	y, b = x.(fmt.State)
	return
}

func GijitShadow_InterfaceConvertTo1_State(x interface{}) fmt.State {
	return x.(fmt.State)
}


func GijitShadow_InterfaceConvertTo2_Stringer(x interface{}) (y fmt.Stringer, b bool) {
	y, b = x.(fmt.Stringer)
	return
}

func GijitShadow_InterfaceConvertTo1_Stringer(x interface{}) fmt.Stringer {
	return x.(fmt.Stringer)
}



 func InitLua() string {
  return `
__type__.fmt ={};

`}