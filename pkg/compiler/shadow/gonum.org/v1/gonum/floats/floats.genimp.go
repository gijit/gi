package shadow_floats

import "gonum.org/v1/gonum/floats"

var Pkg = make(map[string]interface{})

func init() {
	Pkg["Add"] = floats.Add
	Pkg["AddConst"] = floats.AddConst
	Pkg["AddScaled"] = floats.AddScaled
	Pkg["AddScaledTo"] = floats.AddScaledTo
	Pkg["AddTo"] = floats.AddTo
	Pkg["Argsort"] = floats.Argsort
	Pkg["Count"] = floats.Count
	Pkg["CumProd"] = floats.CumProd
	Pkg["CumSum"] = floats.CumSum
	Pkg["Distance"] = floats.Distance
	Pkg["Div"] = floats.Div
	Pkg["DivTo"] = floats.DivTo
	Pkg["Dot"] = floats.Dot
	Pkg["Equal"] = floats.Equal
	Pkg["EqualApprox"] = floats.EqualApprox
	Pkg["EqualFunc"] = floats.EqualFunc
	Pkg["EqualLengths"] = floats.EqualLengths
	Pkg["EqualWithinAbs"] = floats.EqualWithinAbs
	Pkg["EqualWithinAbsOrRel"] = floats.EqualWithinAbsOrRel
	Pkg["EqualWithinRel"] = floats.EqualWithinRel
	Pkg["EqualWithinULP"] = floats.EqualWithinULP
	Pkg["Find"] = floats.Find
	Pkg["HasNaN"] = floats.HasNaN
	Pkg["LogSpan"] = floats.LogSpan
	Pkg["LogSumExp"] = floats.LogSumExp
	Pkg["Max"] = floats.Max
	Pkg["MaxIdx"] = floats.MaxIdx
	Pkg["Min"] = floats.Min
	Pkg["MinIdx"] = floats.MinIdx
	Pkg["Mul"] = floats.Mul
	Pkg["MulTo"] = floats.MulTo
	Pkg["NaNPayload"] = floats.NaNPayload
	Pkg["NaNWith"] = floats.NaNWith
	Pkg["Nearest"] = floats.Nearest
	Pkg["NearestWithinSpan"] = floats.NearestWithinSpan
	Pkg["Norm"] = floats.Norm
	Pkg["ParseWithNA"] = floats.ParseWithNA
	Pkg["Prod"] = floats.Prod
	Pkg["Reverse"] = floats.Reverse
	Pkg["Round"] = floats.Round
	Pkg["RoundEven"] = floats.RoundEven
	Pkg["Same"] = floats.Same
	Pkg["Scale"] = floats.Scale
	Pkg["Span"] = floats.Span
	Pkg["Sub"] = floats.Sub
	Pkg["SubTo"] = floats.SubTo
	Pkg["Sum"] = floats.Sum
	Pkg["Within"] = floats.Within

}
