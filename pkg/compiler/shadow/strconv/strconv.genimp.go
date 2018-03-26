package shadow_strconv

import "strconv"

var Pkg = make(map[string]interface{})
var Ctor = make(map[string]interface{})

func init() {
    Pkg["AppendBool"] = strconv.AppendBool
    Pkg["AppendFloat"] = strconv.AppendFloat
    Pkg["AppendInt"] = strconv.AppendInt
    Pkg["AppendQuote"] = strconv.AppendQuote
    Pkg["AppendQuoteRune"] = strconv.AppendQuoteRune
    Pkg["AppendQuoteRuneToASCII"] = strconv.AppendQuoteRuneToASCII
    Pkg["AppendQuoteRuneToGraphic"] = strconv.AppendQuoteRuneToGraphic
    Pkg["AppendQuoteToASCII"] = strconv.AppendQuoteToASCII
    Pkg["AppendQuoteToGraphic"] = strconv.AppendQuoteToGraphic
    Pkg["AppendUint"] = strconv.AppendUint
    Pkg["Atoi"] = strconv.Atoi
    Pkg["CanBackquote"] = strconv.CanBackquote
    Pkg["ErrRange"] = strconv.ErrRange
    Pkg["ErrSyntax"] = strconv.ErrSyntax
    Pkg["FormatBool"] = strconv.FormatBool
    Pkg["FormatFloat"] = strconv.FormatFloat
    Pkg["FormatInt"] = strconv.FormatInt
    Pkg["FormatUint"] = strconv.FormatUint
    Pkg["IntSize"] = strconv.IntSize
    Pkg["IsGraphic"] = strconv.IsGraphic
    Pkg["IsPrint"] = strconv.IsPrint
    Pkg["Itoa"] = strconv.Itoa
    Ctor["NumError"] = GijitShadow_NewStruct_NumError
    Pkg["ParseBool"] = strconv.ParseBool
    Pkg["ParseFloat"] = strconv.ParseFloat
    Pkg["ParseInt"] = strconv.ParseInt
    Pkg["ParseUint"] = strconv.ParseUint
    Pkg["Quote"] = strconv.Quote
    Pkg["QuoteRune"] = strconv.QuoteRune
    Pkg["QuoteRuneToASCII"] = strconv.QuoteRuneToASCII
    Pkg["QuoteRuneToGraphic"] = strconv.QuoteRuneToGraphic
    Pkg["QuoteToASCII"] = strconv.QuoteToASCII
    Pkg["QuoteToGraphic"] = strconv.QuoteToGraphic
    Pkg["Unquote"] = strconv.Unquote
    Pkg["UnquoteChar"] = strconv.UnquoteChar

}
func GijitShadow_NewStruct_NumError(src *strconv.NumError) *strconv.NumError {
    if src == nil {
	   return &strconv.NumError{}
    }
    a := *src
    return &a
}



 func InitLua() string {
  return `
__type__.strconv ={};

-----------------
-- struct NumError
-----------------

__type__.strconv.NumError = {
 __name = "native_Go_struct_type_wrapper",
 __native_type = "NumError",
 __call = function(t, src)
   return __ctor__strconv.NumError(src)
 end,
};
setmetatable(__type__.strconv.NumError, __type__.strconv.NumError);


`}