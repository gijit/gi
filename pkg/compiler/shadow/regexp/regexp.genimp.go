package shadow_regexp

import "regexp"

var Pkg = make(map[string]interface{})
func init() {
    Pkg["Compile"] = regexp.Compile
    Pkg["CompilePOSIX"] = regexp.CompilePOSIX
    Pkg["Match"] = regexp.Match
    Pkg["MatchReader"] = regexp.MatchReader
    Pkg["MatchString"] = regexp.MatchString
    Pkg["MustCompile"] = regexp.MustCompile
    Pkg["MustCompilePOSIX"] = regexp.MustCompilePOSIX
    Pkg["QuoteMeta"] = regexp.QuoteMeta

}
func GijitShadow_NewStruct_Regexp() *regexp.Regexp {
	return &regexp.Regexp{}
}

