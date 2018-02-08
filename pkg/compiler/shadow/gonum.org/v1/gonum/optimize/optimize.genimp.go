package shadow_gonum.org/v1/gonum/optimize

import "gonum.org/v1/gonum/optimize"

var Pkg = make(map[string]interface{})
func init() {
    Pkg["ArmijoConditionMet"] = optimize.ArmijoConditionMet
    Pkg["CGVariant"] = GijitShadow_InterfaceConvertTo2_CGVariant
    Pkg["DefaultSettings"] = optimize.DefaultSettings
    Pkg["DefaultSettingsGlobal"] = optimize.DefaultSettingsGlobal
    Pkg["ErrLinesearcherBound"] = optimize.ErrLinesearcherBound
    Pkg["ErrLinesearcherFailure"] = optimize.ErrLinesearcherFailure
    Pkg["ErrNoProgress"] = optimize.ErrNoProgress
    Pkg["ErrNonDescentDirection"] = optimize.ErrNonDescentDirection
    Pkg["ErrZeroDimensional"] = optimize.ErrZeroDimensional
    Pkg["Global"] = optimize.Global
    Pkg["GlobalMethod"] = GijitShadow_InterfaceConvertTo2_GlobalMethod
    Pkg["Linesearcher"] = GijitShadow_InterfaceConvertTo2_Linesearcher
    Pkg["Local"] = optimize.Local
    Pkg["Method"] = GijitShadow_InterfaceConvertTo2_Method
    Pkg["Needser"] = GijitShadow_InterfaceConvertTo2_Needser
    Pkg["NewPrinter"] = optimize.NewPrinter
    Pkg["NewStatus"] = optimize.NewStatus
    Pkg["NextDirectioner"] = GijitShadow_InterfaceConvertTo2_NextDirectioner
    Pkg["Recorder"] = GijitShadow_InterfaceConvertTo2_Recorder
    Pkg["Statuser"] = GijitShadow_InterfaceConvertTo2_Statuser
    Pkg["StepSizer"] = GijitShadow_InterfaceConvertTo2_StepSizer
    Pkg["StrongWolfeConditionsMet"] = optimize.StrongWolfeConditionsMet
    Pkg["WeakWolfeConditionsMet"] = optimize.WeakWolfeConditionsMet

}
func GijitShadow_NewStruct_BFGS() *optimize.BFGS {
	return &optimize.BFGS{}
}


func GijitShadow_NewStruct_Backtracking() *optimize.Backtracking {
	return &optimize.Backtracking{}
}


func GijitShadow_NewStruct_Bisection() *optimize.Bisection {
	return &optimize.Bisection{}
}


func GijitShadow_NewStruct_CG() *optimize.CG {
	return &optimize.CG{}
}


func GijitShadow_InterfaceConvertTo2_CGVariant(x interface{}) (y optimize.CGVariant, b bool) {
	y, b = x.(optimize.CGVariant)
	return
}

func GijitShadow_InterfaceConvertTo1_CGVariant(x interface{}) optimize.CGVariant {
	return x.(optimize.CGVariant)
}


func GijitShadow_NewStruct_CmaEsChol() *optimize.CmaEsChol {
	return &optimize.CmaEsChol{}
}


func GijitShadow_NewStruct_ConstantStepSize() *optimize.ConstantStepSize {
	return &optimize.ConstantStepSize{}
}


func GijitShadow_NewStruct_DaiYuan() *optimize.DaiYuan {
	return &optimize.DaiYuan{}
}


func GijitShadow_NewStruct_ErrGrad() *optimize.ErrGrad {
	return &optimize.ErrGrad{}
}


func GijitShadow_NewStruct_FirstOrderStepSize() *optimize.FirstOrderStepSize {
	return &optimize.FirstOrderStepSize{}
}


func GijitShadow_NewStruct_FletcherReeves() *optimize.FletcherReeves {
	return &optimize.FletcherReeves{}
}


func GijitShadow_NewStruct_FunctionConverge() *optimize.FunctionConverge {
	return &optimize.FunctionConverge{}
}


func GijitShadow_InterfaceConvertTo2_GlobalMethod(x interface{}) (y optimize.GlobalMethod, b bool) {
	y, b = x.(optimize.GlobalMethod)
	return
}

func GijitShadow_InterfaceConvertTo1_GlobalMethod(x interface{}) optimize.GlobalMethod {
	return x.(optimize.GlobalMethod)
}


func GijitShadow_NewStruct_GlobalTask() *optimize.GlobalTask {
	return &optimize.GlobalTask{}
}


func GijitShadow_NewStruct_GradientDescent() *optimize.GradientDescent {
	return &optimize.GradientDescent{}
}


func GijitShadow_NewStruct_GuessAndCheck() *optimize.GuessAndCheck {
	return &optimize.GuessAndCheck{}
}


func GijitShadow_NewStruct_HagerZhang() *optimize.HagerZhang {
	return &optimize.HagerZhang{}
}


func GijitShadow_NewStruct_HestenesStiefel() *optimize.HestenesStiefel {
	return &optimize.HestenesStiefel{}
}


func GijitShadow_NewStruct_LBFGS() *optimize.LBFGS {
	return &optimize.LBFGS{}
}


func GijitShadow_NewStruct_LinesearchMethod() *optimize.LinesearchMethod {
	return &optimize.LinesearchMethod{}
}


func GijitShadow_InterfaceConvertTo2_Linesearcher(x interface{}) (y optimize.Linesearcher, b bool) {
	y, b = x.(optimize.Linesearcher)
	return
}

func GijitShadow_InterfaceConvertTo1_Linesearcher(x interface{}) optimize.Linesearcher {
	return x.(optimize.Linesearcher)
}


func GijitShadow_NewStruct_Location() *optimize.Location {
	return &optimize.Location{}
}


func GijitShadow_InterfaceConvertTo2_Method(x interface{}) (y optimize.Method, b bool) {
	y, b = x.(optimize.Method)
	return
}

func GijitShadow_InterfaceConvertTo1_Method(x interface{}) optimize.Method {
	return x.(optimize.Method)
}


func GijitShadow_NewStruct_MoreThuente() *optimize.MoreThuente {
	return &optimize.MoreThuente{}
}


func GijitShadow_InterfaceConvertTo2_Needser(x interface{}) (y optimize.Needser, b bool) {
	y, b = x.(optimize.Needser)
	return
}

func GijitShadow_InterfaceConvertTo1_Needser(x interface{}) optimize.Needser {
	return x.(optimize.Needser)
}


func GijitShadow_NewStruct_NelderMead() *optimize.NelderMead {
	return &optimize.NelderMead{}
}


func GijitShadow_NewStruct_Newton() *optimize.Newton {
	return &optimize.Newton{}
}


func GijitShadow_InterfaceConvertTo2_NextDirectioner(x interface{}) (y optimize.NextDirectioner, b bool) {
	y, b = x.(optimize.NextDirectioner)
	return
}

func GijitShadow_InterfaceConvertTo1_NextDirectioner(x interface{}) optimize.NextDirectioner {
	return x.(optimize.NextDirectioner)
}


func GijitShadow_NewStruct_PolakRibierePolyak() *optimize.PolakRibierePolyak {
	return &optimize.PolakRibierePolyak{}
}


func GijitShadow_NewStruct_Printer() *optimize.Printer {
	return &optimize.Printer{}
}


func GijitShadow_NewStruct_Problem() *optimize.Problem {
	return &optimize.Problem{}
}


func GijitShadow_NewStruct_QuadraticStepSize() *optimize.QuadraticStepSize {
	return &optimize.QuadraticStepSize{}
}


func GijitShadow_InterfaceConvertTo2_Recorder(x interface{}) (y optimize.Recorder, b bool) {
	y, b = x.(optimize.Recorder)
	return
}

func GijitShadow_InterfaceConvertTo1_Recorder(x interface{}) optimize.Recorder {
	return x.(optimize.Recorder)
}


func GijitShadow_NewStruct_Result() *optimize.Result {
	return &optimize.Result{}
}


func GijitShadow_NewStruct_Settings() *optimize.Settings {
	return &optimize.Settings{}
}


func GijitShadow_NewStruct_Stats() *optimize.Stats {
	return &optimize.Stats{}
}


func GijitShadow_InterfaceConvertTo2_Statuser(x interface{}) (y optimize.Statuser, b bool) {
	y, b = x.(optimize.Statuser)
	return
}

func GijitShadow_InterfaceConvertTo1_Statuser(x interface{}) optimize.Statuser {
	return x.(optimize.Statuser)
}


func GijitShadow_InterfaceConvertTo2_StepSizer(x interface{}) (y optimize.StepSizer, b bool) {
	y, b = x.(optimize.StepSizer)
	return
}

func GijitShadow_InterfaceConvertTo1_StepSizer(x interface{}) optimize.StepSizer {
	return x.(optimize.StepSizer)
}

