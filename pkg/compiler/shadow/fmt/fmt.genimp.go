package shadow_fmt

import "fmt"

var Pkg = make(map[string]interface{})

func init() {
	Pkg["Errorf"] = fmt.Errorf
	Pkg["Fprint"] = fmt.Fprint
	Pkg["Fprintf"] = fmt.Fprintf
	Pkg["Fprintln"] = fmt.Fprintln
	Pkg["Fscan"] = fmt.Fscan
	Pkg["Fscanf"] = fmt.Fscanf
	Pkg["Fscanln"] = fmt.Fscanln
	Pkg["Print"] = fmt.Print
	Pkg["Printf"] = fmt.Printf
	Pkg["Println"] = fmt.Println
	Pkg["Scan"] = fmt.Scan
	Pkg["Scanf"] = fmt.Scanf
	Pkg["Scanln"] = fmt.Scanln
	Pkg["Sprint"] = fmt.Sprint
	Pkg["Sprintf"] = fmt.Sprintf
	Pkg["Sprintln"] = fmt.Sprintln
	Pkg["Sscan"] = fmt.Sscan
	Pkg["Sscanf"] = fmt.Sscanf
	Pkg["Sscanln"] = fmt.Sscanln

}
