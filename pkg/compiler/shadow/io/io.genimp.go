package shadow_io

import "io"

var Pkg = make(map[string]interface{})
func init() {
    Pkg["ByteReader"] = __gi_ConvertTo_ByteReader
    Pkg["ByteScanner"] = __gi_ConvertTo_ByteScanner
    Pkg["ByteWriter"] = __gi_ConvertTo_ByteWriter
    Pkg["Closer"] = __gi_ConvertTo_Closer
    Pkg["Copy"] = io.Copy
    Pkg["CopyBuffer"] = io.CopyBuffer
    Pkg["CopyN"] = io.CopyN
    Pkg["EOF"] = io.EOF
    Pkg["ErrClosedPipe"] = io.ErrClosedPipe
    Pkg["ErrNoProgress"] = io.ErrNoProgress
    Pkg["ErrShortBuffer"] = io.ErrShortBuffer
    Pkg["ErrShortWrite"] = io.ErrShortWrite
    Pkg["ErrUnexpectedEOF"] = io.ErrUnexpectedEOF
    Pkg["LimitReader"] = io.LimitReader
    Pkg["MultiReader"] = io.MultiReader
    Pkg["MultiWriter"] = io.MultiWriter
    Pkg["NewSectionReader"] = io.NewSectionReader
    Pkg["Pipe"] = io.Pipe
    Pkg["ReadAtLeast"] = io.ReadAtLeast
    Pkg["ReadCloser"] = __gi_ConvertTo_ReadCloser
    Pkg["ReadFull"] = io.ReadFull
    Pkg["ReadSeeker"] = __gi_ConvertTo_ReadSeeker
    Pkg["ReadWriteCloser"] = __gi_ConvertTo_ReadWriteCloser
    Pkg["ReadWriteSeeker"] = __gi_ConvertTo_ReadWriteSeeker
    Pkg["ReadWriter"] = __gi_ConvertTo_ReadWriter
    Pkg["Reader"] = __gi_ConvertTo_Reader
    Pkg["ReaderAt"] = __gi_ConvertTo_ReaderAt
    Pkg["ReaderFrom"] = __gi_ConvertTo_ReaderFrom
    Pkg["RuneReader"] = __gi_ConvertTo_RuneReader
    Pkg["RuneScanner"] = __gi_ConvertTo_RuneScanner
    Pkg["SeekCurrent"] = io.SeekCurrent
    Pkg["SeekEnd"] = io.SeekEnd
    Pkg["SeekStart"] = io.SeekStart
    Pkg["Seeker"] = __gi_ConvertTo_Seeker
    Pkg["TeeReader"] = io.TeeReader
    Pkg["WriteCloser"] = __gi_ConvertTo_WriteCloser
    Pkg["WriteSeeker"] = __gi_ConvertTo_WriteSeeker
    Pkg["WriteString"] = io.WriteString
    Pkg["Writer"] = __gi_ConvertTo_Writer
    Pkg["WriterAt"] = __gi_ConvertTo_WriterAt
    Pkg["WriterTo"] = __gi_ConvertTo_WriterTo

}
func __gi_ConvertTo_ByteReader(x interface{}) (y io.ByteReader, b bool) {
	y, b = x.(io.ByteReader)
	return
}


func __gi_ConvertTo_ByteScanner(x interface{}) (y io.ByteScanner, b bool) {
	y, b = x.(io.ByteScanner)
	return
}


func __gi_ConvertTo_ByteWriter(x interface{}) (y io.ByteWriter, b bool) {
	y, b = x.(io.ByteWriter)
	return
}


func __gi_ConvertTo_Closer(x interface{}) (y io.Closer, b bool) {
	y, b = x.(io.Closer)
	return
}


func __gi_ConvertTo_ReadCloser(x interface{}) (y io.ReadCloser, b bool) {
	y, b = x.(io.ReadCloser)
	return
}


func __gi_ConvertTo_ReadSeeker(x interface{}) (y io.ReadSeeker, b bool) {
	y, b = x.(io.ReadSeeker)
	return
}


func __gi_ConvertTo_ReadWriteCloser(x interface{}) (y io.ReadWriteCloser, b bool) {
	y, b = x.(io.ReadWriteCloser)
	return
}


func __gi_ConvertTo_ReadWriteSeeker(x interface{}) (y io.ReadWriteSeeker, b bool) {
	y, b = x.(io.ReadWriteSeeker)
	return
}


func __gi_ConvertTo_ReadWriter(x interface{}) (y io.ReadWriter, b bool) {
	y, b = x.(io.ReadWriter)
	return
}


func __gi_ConvertTo_Reader(x interface{}) (y io.Reader, b bool) {
	y, b = x.(io.Reader)
	return
}


func __gi_ConvertTo_ReaderAt(x interface{}) (y io.ReaderAt, b bool) {
	y, b = x.(io.ReaderAt)
	return
}


func __gi_ConvertTo_ReaderFrom(x interface{}) (y io.ReaderFrom, b bool) {
	y, b = x.(io.ReaderFrom)
	return
}


func __gi_ConvertTo_RuneReader(x interface{}) (y io.RuneReader, b bool) {
	y, b = x.(io.RuneReader)
	return
}


func __gi_ConvertTo_RuneScanner(x interface{}) (y io.RuneScanner, b bool) {
	y, b = x.(io.RuneScanner)
	return
}


func __gi_ConvertTo_Seeker(x interface{}) (y io.Seeker, b bool) {
	y, b = x.(io.Seeker)
	return
}


func __gi_ConvertTo_WriteCloser(x interface{}) (y io.WriteCloser, b bool) {
	y, b = x.(io.WriteCloser)
	return
}


func __gi_ConvertTo_WriteSeeker(x interface{}) (y io.WriteSeeker, b bool) {
	y, b = x.(io.WriteSeeker)
	return
}


func __gi_ConvertTo_Writer(x interface{}) (y io.Writer, b bool) {
	y, b = x.(io.Writer)
	return
}


func __gi_ConvertTo_WriterAt(x interface{}) (y io.WriterAt, b bool) {
	y, b = x.(io.WriterAt)
	return
}


func __gi_ConvertTo_WriterTo(x interface{}) (y io.WriterTo, b bool) {
	y, b = x.(io.WriterTo)
	return
}

