package shadow_regexp

import "regexp"

var Pkg = make(map[string]interface{})
var Ctor = make(map[string]interface{})

func init() {
	Pkg["Compile"] = regexp.Compile
	Pkg["CompilePOSIX"] = regexp.CompilePOSIX
	Pkg["Match"] = regexp.Match
	Pkg["MatchReader"] = regexp.MatchReader
	Pkg["MatchString"] = regexp.MatchString
	Pkg["MustCompile"] = regexp.MustCompile
	Pkg["MustCompilePOSIX"] = regexp.MustCompilePOSIX
	Pkg["QuoteMeta"] = regexp.QuoteMeta
	Ctor["Regexp"] = GijitShadow_NewStruct_Regexp

}
func GijitShadow_NewStruct_Regexp(src *regexp.Regexp) *regexp.Regexp {
	if src == nil {
		return &regexp.Regexp{}
	}
	a := *src
	return &a
}

func InitLua() string {
	return `
print("here is __type__ before our addition of regexp:")
__st(__type__, "__type__")
__type__.regexp ={};
print("here is __type__ after our addition of regexp:")
__st(__type__, "__type__")

-----------------
-- struct Regexp
-----------------

__type__.regexp.Regexp = {
 __name = "native_Go_struct_type_wrapper",
 __native_type = "Regexp",
 __str = "Regexp",
 exported = true,
 __call = function(t, src)
   return __ctor__regexp.Regexp(src)
 end,
};
setmetatable(__type__.regexp.Regexp, __type__.regexp.Regexp);

print("here is __type__.regexp after our addition of Regexp:")
__st(__type__.regexp, "__type__.regexp")

`
}
