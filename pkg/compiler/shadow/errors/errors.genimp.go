package shadow_errors

import "errors"

var Pkg = make(map[string]interface{})
var Ctor = make(map[string]interface{})

func init() {
    Pkg["New"] = errors.New

}

 func InitLua() string {
  return `
__type__.errors ={};

`}