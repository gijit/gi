package shadow_gonum.org/v1/gonum/unit

import "gonum.org/v1/gonum/unit"

var Pkg = make(map[string]interface{})
func init() {
    Pkg["Atto"] = unit.Atto
    Pkg["Centi"] = unit.Centi
    Pkg["Deca"] = unit.Deca
    Pkg["Deci"] = unit.Deci
    Pkg["DimensionsMatch"] = unit.DimensionsMatch
    Pkg["Exa"] = unit.Exa
    Pkg["Femto"] = unit.Femto
    Pkg["Giga"] = unit.Giga
    Pkg["Hecto"] = unit.Hecto
    Pkg["Kilo"] = unit.Kilo
    Pkg["Mega"] = unit.Mega
    Pkg["Micro"] = unit.Micro
    Pkg["Milli"] = unit.Milli
    Pkg["Nano"] = unit.Nano
    Pkg["New"] = unit.New
    Pkg["NewDimension"] = unit.NewDimension
    Pkg["Peta"] = unit.Peta
    Pkg["Pico"] = unit.Pico
    Pkg["SymbolExists"] = unit.SymbolExists
    Pkg["Tera"] = unit.Tera
    Pkg["Uniter"] = GijitShadow_InterfaceConvertTo2_Uniter
    Pkg["Yocto"] = unit.Yocto
    Pkg["Yotta"] = unit.Yotta
    Pkg["Zepto"] = unit.Zepto
    Pkg["Zetta"] = unit.Zetta

}
func GijitShadow_NewStruct_Unit() *unit.Unit {
	return &unit.Unit{}
}


func GijitShadow_InterfaceConvertTo2_Uniter(x interface{}) (y unit.Uniter, b bool) {
	y, b = x.(unit.Uniter)
	return
}

func GijitShadow_InterfaceConvertTo1_Uniter(x interface{}) unit.Uniter {
	return x.(unit.Uniter)
}

