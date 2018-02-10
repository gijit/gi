package shadow_integrate

import "gonum.org/v1/gonum/integrate"

var Pkg = make(map[string]interface{})

func init() {
	Pkg["Trapezoidal"] = integrate.Trapezoidal

}
