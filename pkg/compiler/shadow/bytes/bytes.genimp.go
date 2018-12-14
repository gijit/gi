package shadow_bytes

import "bytes"

var Pkg = make(map[string]interface{})
var Ctor = make(map[string]interface{})

func init() {
    Ctor["Buffer"] = GijitShadow_NewStruct_Buffer
    Pkg["Compare"] = bytes.Compare
    Pkg["Contains"] = bytes.Contains
    Pkg["ContainsAny"] = bytes.ContainsAny
    Pkg["ContainsRune"] = bytes.ContainsRune
    Pkg["Count"] = bytes.Count
    Pkg["Equal"] = bytes.Equal
    Pkg["EqualFold"] = bytes.EqualFold
    Pkg["ErrTooLarge"] = bytes.ErrTooLarge
    Pkg["Fields"] = bytes.Fields
    Pkg["FieldsFunc"] = bytes.FieldsFunc
    Pkg["HasPrefix"] = bytes.HasPrefix
    Pkg["HasSuffix"] = bytes.HasSuffix
    Pkg["Index"] = bytes.Index
    Pkg["IndexAny"] = bytes.IndexAny
    Pkg["IndexByte"] = bytes.IndexByte
    Pkg["IndexFunc"] = bytes.IndexFunc
    Pkg["IndexRune"] = bytes.IndexRune
    Pkg["Join"] = bytes.Join
    Pkg["LastIndex"] = bytes.LastIndex
    Pkg["LastIndexAny"] = bytes.LastIndexAny
    Pkg["LastIndexByte"] = bytes.LastIndexByte
    Pkg["LastIndexFunc"] = bytes.LastIndexFunc
    Pkg["Map"] = bytes.Map
    Pkg["MinRead"] = bytes.MinRead
    Pkg["NewBuffer"] = bytes.NewBuffer
    Pkg["NewBufferString"] = bytes.NewBufferString
    Pkg["NewReader"] = bytes.NewReader
    Ctor["Reader"] = GijitShadow_NewStruct_Reader
    Pkg["Repeat"] = bytes.Repeat
    Pkg["Replace"] = bytes.Replace
    Pkg["Runes"] = bytes.Runes
    Pkg["Split"] = bytes.Split
    Pkg["SplitAfter"] = bytes.SplitAfter
    Pkg["SplitAfterN"] = bytes.SplitAfterN
    Pkg["SplitN"] = bytes.SplitN
    Pkg["Title"] = bytes.Title
    Pkg["ToLower"] = bytes.ToLower
    Pkg["ToLowerSpecial"] = bytes.ToLowerSpecial
    Pkg["ToTitle"] = bytes.ToTitle
    Pkg["ToTitleSpecial"] = bytes.ToTitleSpecial
    Pkg["ToUpper"] = bytes.ToUpper
    Pkg["ToUpperSpecial"] = bytes.ToUpperSpecial
    Pkg["Trim"] = bytes.Trim
    Pkg["TrimFunc"] = bytes.TrimFunc
    Pkg["TrimLeft"] = bytes.TrimLeft
    Pkg["TrimLeftFunc"] = bytes.TrimLeftFunc
    Pkg["TrimPrefix"] = bytes.TrimPrefix
    Pkg["TrimRight"] = bytes.TrimRight
    Pkg["TrimRightFunc"] = bytes.TrimRightFunc
    Pkg["TrimSpace"] = bytes.TrimSpace
    Pkg["TrimSuffix"] = bytes.TrimSuffix

}
func GijitShadow_NewStruct_Buffer(src *bytes.Buffer) *bytes.Buffer {
    if src == nil {
	   return &bytes.Buffer{}
    }
    a := *src
    return &a
}


func GijitShadow_NewStruct_Reader(src *bytes.Reader) *bytes.Reader {
    if src == nil {
	   return &bytes.Reader{}
    }
    a := *src
    return &a
}



 func InitLua() string {
  return `
__type__.bytes ={};

-----------------
-- struct Buffer
-----------------

__type__.bytes.Buffer = {
 id=0,
 __name = "native_Go_struct_type_wrapper",
 __native_type = "Buffer",
 __str = "Buffer",
 exported = true,
 __call = function(t, src)
   return __ctor__bytes.Buffer(src)
 end,
};
setmetatable(__type__.bytes.Buffer, __type__.bytes.Buffer);


-----------------
-- struct Reader
-----------------

__type__.bytes.Reader = {
 id=0,
 __name = "native_Go_struct_type_wrapper",
 __native_type = "Reader",
 __str = "Reader",
 exported = true,
 __call = function(t, src)
   return __ctor__bytes.Reader(src)
 end,
};
setmetatable(__type__.bytes.Reader, __type__.bytes.Reader);


`}