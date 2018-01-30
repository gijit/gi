package shadow_io

import "io"

var Pkg = make(map[string]interface{})
func init() {
    Pkg["ByteReader"] = __gi_ConvertTo2_ByteReader
    Pkg["ByteScanner"] = __gi_ConvertTo2_ByteScanner
    Pkg["ByteWriter"] = __gi_ConvertTo2_ByteWriter
    Pkg["Closer"] = __gi_ConvertTo2_Closer
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
    Pkg["ReadCloser"] = __gi_ConvertTo2_ReadCloser
    Pkg["ReadFull"] = io.ReadFull
    Pkg["ReadSeeker"] = __gi_ConvertTo2_ReadSeeker
    Pkg["ReadWriteCloser"] = __gi_ConvertTo2_ReadWriteCloser
    Pkg["ReadWriteSeeker"] = __gi_ConvertTo2_ReadWriteSeeker
    Pkg["ReadWriter"] = __gi_ConvertTo2_ReadWriter
    Pkg["Reader"] = __gi_ConvertTo2_Reader
    Pkg["ReaderAt"] = __gi_ConvertTo2_ReaderAt
    Pkg["ReaderFrom"] = __gi_ConvertTo2_ReaderFrom
    Pkg["RuneReader"] = __gi_ConvertTo2_RuneReader
    Pkg["RuneScanner"] = __gi_ConvertTo2_RuneScanner
    Pkg["SeekCurrent"] = io.SeekCurrent
    Pkg["SeekEnd"] = io.SeekEnd
    Pkg["SeekStart"] = io.SeekStart
    Pkg["Seeker"] = __gi_ConvertTo2_Seeker
    Pkg["TeeReader"] = io.TeeReader
    Pkg["WriteCloser"] = __gi_ConvertTo2_WriteCloser
    Pkg["WriteSeeker"] = __gi_ConvertTo2_WriteSeeker
    Pkg["WriteString"] = io.WriteString
    Pkg["Writer"] = __gi_ConvertTo2_Writer
    Pkg["WriterAt"] = __gi_ConvertTo2_WriterAt
    Pkg["WriterTo"] = __gi_ConvertTo2_WriterTo

}
func __gi_ConvertTo2_ByteReader(x interface{}) (y io.ByteReader, b bool) {
	y, b = x.(io.ByteReader)
	return
}

func __gi_ConvertTo1_ByteReader(x interface{}) io.ByteReader {
	return x.(io.ByteReader)
}


func __gi_ConvertTo2_ByteScanner(x interface{}) (y io.ByteScanner, b bool) {
	y, b = x.(io.ByteScanner)
	return
}

func __gi_ConvertTo1_ByteScanner(x interface{}) io.ByteScanner {
	return x.(io.ByteScanner)
}


func __gi_ConvertTo2_ByteWriter(x interface{}) (y io.ByteWriter, b bool) {
	y, b = x.(io.ByteWriter)
	return
}

func __gi_ConvertTo1_ByteWriter(x interface{}) io.ByteWriter {
	return x.(io.ByteWriter)
}


func __gi_ConvertTo2_Closer(x interface{}) (y io.Closer, b bool) {
	y, b = x.(io.Closer)
	return
}

func __gi_ConvertTo1_Closer(x interface{}) io.Closer {
	return x.(io.Closer)
}


func __gi_ConvertTo2_ReadCloser(x interface{}) (y io.ReadCloser, b bool) {
	y, b = x.(io.ReadCloser)
	return
}

func __gi_ConvertTo1_ReadCloser(x interface{}) io.ReadCloser {
	return x.(io.ReadCloser)
}


func __gi_ConvertTo2_ReadSeeker(x interface{}) (y io.ReadSeeker, b bool) {
	y, b = x.(io.ReadSeeker)
	return
}

func __gi_ConvertTo1_ReadSeeker(x interface{}) io.ReadSeeker {
	return x.(io.ReadSeeker)
}


func __gi_ConvertTo2_ReadWriteCloser(x interface{}) (y io.ReadWriteCloser, b bool) {
	y, b = x.(io.ReadWriteCloser)
	return
}

func __gi_ConvertTo1_ReadWriteCloser(x interface{}) io.ReadWriteCloser {
	return x.(io.ReadWriteCloser)
}


func __gi_ConvertTo2_ReadWriteSeeker(x interface{}) (y io.ReadWriteSeeker, b bool) {
	y, b = x.(io.ReadWriteSeeker)
	return
}

func __gi_ConvertTo1_ReadWriteSeeker(x interface{}) io.ReadWriteSeeker {
	return x.(io.ReadWriteSeeker)
}


func __gi_ConvertTo2_ReadWriter(x interface{}) (y io.ReadWriter, b bool) {
	y, b = x.(io.ReadWriter)
	return
}

func __gi_ConvertTo1_ReadWriter(x interface{}) io.ReadWriter {
	return x.(io.ReadWriter)
}


func __gi_ConvertTo2_Reader(x interface{}) (y io.Reader, b bool) {
	y, b = x.(io.Reader)
	return
}

func __gi_ConvertTo1_Reader(x interface{}) io.Reader {
	return x.(io.Reader)
}


func __gi_ConvertTo2_ReaderAt(x interface{}) (y io.ReaderAt, b bool) {
	y, b = x.(io.ReaderAt)
	return
}

func __gi_ConvertTo1_ReaderAt(x interface{}) io.ReaderAt {
	return x.(io.ReaderAt)
}


func __gi_ConvertTo2_ReaderFrom(x interface{}) (y io.ReaderFrom, b bool) {
	y, b = x.(io.ReaderFrom)
	return
}

func __gi_ConvertTo1_ReaderFrom(x interface{}) io.ReaderFrom {
	return x.(io.ReaderFrom)
}


func __gi_ConvertTo2_RuneReader(x interface{}) (y io.RuneReader, b bool) {
	y, b = x.(io.RuneReader)
	return
}

func __gi_ConvertTo1_RuneReader(x interface{}) io.RuneReader {
	return x.(io.RuneReader)
}


func __gi_ConvertTo2_RuneScanner(x interface{}) (y io.RuneScanner, b bool) {
	y, b = x.(io.RuneScanner)
	return
}

func __gi_ConvertTo1_RuneScanner(x interface{}) io.RuneScanner {
	return x.(io.RuneScanner)
}


func __gi_ConvertTo2_Seeker(x interface{}) (y io.Seeker, b bool) {
	y, b = x.(io.Seeker)
	return
}

func __gi_ConvertTo1_Seeker(x interface{}) io.Seeker {
	return x.(io.Seeker)
}


func __gi_ConvertTo2_WriteCloser(x interface{}) (y io.WriteCloser, b bool) {
	y, b = x.(io.WriteCloser)
	return
}

func __gi_ConvertTo1_WriteCloser(x interface{}) io.WriteCloser {
	return x.(io.WriteCloser)
}


func __gi_ConvertTo2_WriteSeeker(x interface{}) (y io.WriteSeeker, b bool) {
	y, b = x.(io.WriteSeeker)
	return
}

func __gi_ConvertTo1_WriteSeeker(x interface{}) io.WriteSeeker {
	return x.(io.WriteSeeker)
}


func __gi_ConvertTo2_Writer(x interface{}) (y io.Writer, b bool) {
	y, b = x.(io.Writer)
	return
}

func __gi_ConvertTo1_Writer(x interface{}) io.Writer {
	return x.(io.Writer)
}


func __gi_ConvertTo2_WriterAt(x interface{}) (y io.WriterAt, b bool) {
	y, b = x.(io.WriterAt)
	return
}

func __gi_ConvertTo1_WriterAt(x interface{}) io.WriterAt {
	return x.(io.WriterAt)
}


func __gi_ConvertTo2_WriterTo(x interface{}) (y io.WriterTo, b bool) {
	y, b = x.(io.WriterTo)
	return
}

func __gi_ConvertTo1_WriterTo(x interface{}) io.WriterTo {
	return x.(io.WriterTo)
}

