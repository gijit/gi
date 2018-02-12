package shadow_ioutil

import "io/ioutil"

var Pkg = make(map[string]interface{})
func init() {
    Pkg["Discard"] = ioutil.Discard
    Pkg["NopCloser"] = ioutil.NopCloser
    Pkg["ReadAll"] = ioutil.ReadAll
    Pkg["ReadDir"] = ioutil.ReadDir
    Pkg["ReadFile"] = ioutil.ReadFile
    Pkg["TempDir"] = ioutil.TempDir
    Pkg["TempFile"] = ioutil.TempFile
    Pkg["WriteFile"] = ioutil.WriteFile

}