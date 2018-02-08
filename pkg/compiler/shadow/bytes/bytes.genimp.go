package shadow_bytes

import "bytes"

var Pkg = make(map[string]interface{})
func init() {
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
func GijitShadow_NewStruct_Buffer() *bytes.Buffer {
	return &bytes.Buffer{}
}


func GijitShadow_NewStruct_Reader() *bytes.Reader {
	return &bytes.Reader{}
}

