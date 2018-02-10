package shadow_fd

import "gonum.org/v1/gonum/diff/fd"

var Pkg = make(map[string]interface{})

func init() {
	Pkg["Backward"] = fd.Backward
	Pkg["Backward2nd"] = fd.Backward2nd
	Pkg["Central"] = fd.Central
	Pkg["Central2nd"] = fd.Central2nd
	Pkg["CrossLaplacian"] = fd.CrossLaplacian
	Pkg["Derivative"] = fd.Derivative
	Pkg["Forward"] = fd.Forward
	Pkg["Forward2nd"] = fd.Forward2nd
	Pkg["Gradient"] = fd.Gradient
	Pkg["Hessian"] = fd.Hessian
	Pkg["Jacobian"] = fd.Jacobian
	Pkg["Laplacian"] = fd.Laplacian

}
func GijitShadow_NewStruct_Formula() *fd.Formula {
	return &fd.Formula{}
}

func GijitShadow_NewStruct_JacobianSettings() *fd.JacobianSettings {
	return &fd.JacobianSettings{}
}

func GijitShadow_NewStruct_Point() *fd.Point {
	return &fd.Point{}
}

func GijitShadow_NewStruct_Settings() *fd.Settings {
	return &fd.Settings{}
}
