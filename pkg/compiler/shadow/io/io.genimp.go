package shadow_io

import "io"

var Pkg = make(map[string]interface{})
func init() {
    Pkg["ByteReader"] = GijitShadow_InterfaceConvertTo2_ByteReader
    Pkg["ByteScanner"] = GijitShadow_InterfaceConvertTo2_ByteScanner
    Pkg["ByteWriter"] = GijitShadow_InterfaceConvertTo2_ByteWriter
    Pkg["Closer"] = GijitShadow_InterfaceConvertTo2_Closer
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
    Pkg["ReadCloser"] = GijitShadow_InterfaceConvertTo2_ReadCloser
    Pkg["ReadFull"] = io.ReadFull
    Pkg["ReadSeeker"] = GijitShadow_InterfaceConvertTo2_ReadSeeker
    Pkg["ReadWriteCloser"] = GijitShadow_InterfaceConvertTo2_ReadWriteCloser
    Pkg["ReadWriteSeeker"] = GijitShadow_InterfaceConvertTo2_ReadWriteSeeker
    Pkg["ReadWriter"] = GijitShadow_InterfaceConvertTo2_ReadWriter
    Pkg["Reader"] = GijitShadow_InterfaceConvertTo2_Reader
    Pkg["ReaderAt"] = GijitShadow_InterfaceConvertTo2_ReaderAt
    Pkg["ReaderFrom"] = GijitShadow_InterfaceConvertTo2_ReaderFrom
    Pkg["RuneReader"] = GijitShadow_InterfaceConvertTo2_RuneReader
    Pkg["RuneScanner"] = GijitShadow_InterfaceConvertTo2_RuneScanner
    Pkg["SeekCurrent"] = io.SeekCurrent
    Pkg["SeekEnd"] = io.SeekEnd
    Pkg["SeekStart"] = io.SeekStart
    Pkg["Seeker"] = GijitShadow_InterfaceConvertTo2_Seeker
    Pkg["TeeReader"] = io.TeeReader
    Pkg["WriteCloser"] = GijitShadow_InterfaceConvertTo2_WriteCloser
    Pkg["WriteSeeker"] = GijitShadow_InterfaceConvertTo2_WriteSeeker
    Pkg["WriteString"] = io.WriteString
    Pkg["Writer"] = GijitShadow_InterfaceConvertTo2_Writer
    Pkg["WriterAt"] = GijitShadow_InterfaceConvertTo2_WriterAt
    Pkg["WriterTo"] = GijitShadow_InterfaceConvertTo2_WriterTo

}
func GijitShadow_InterfaceConvertTo2_ByteReader(x interface{}) (y io.ByteReader, b bool) {
	y, b = x.(io.ByteReader)
	return
}

func GijitShadow_InterfaceConvertTo1_ByteReader(x interface{}) io.ByteReader {
	return x.(io.ByteReader)
}


func GijitShadow_InterfaceConvertTo2_ByteScanner(x interface{}) (y io.ByteScanner, b bool) {
	y, b = x.(io.ByteScanner)
	return
}

func GijitShadow_InterfaceConvertTo1_ByteScanner(x interface{}) io.ByteScanner {
	return x.(io.ByteScanner)
}


func GijitShadow_InterfaceConvertTo2_ByteWriter(x interface{}) (y io.ByteWriter, b bool) {
	y, b = x.(io.ByteWriter)
	return
}

func GijitShadow_InterfaceConvertTo1_ByteWriter(x interface{}) io.ByteWriter {
	return x.(io.ByteWriter)
}


func GijitShadow_InterfaceConvertTo2_Closer(x interface{}) (y io.Closer, b bool) {
	y, b = x.(io.Closer)
	return
}

func GijitShadow_InterfaceConvertTo1_Closer(x interface{}) io.Closer {
	return x.(io.Closer)
}


func GijitShadow_NewStruct_LimitedReader() *io.LimitedReader {
	return &io.LimitedReader{}
}


func GijitShadow_NewStruct_PipeReader() *io.PipeReader {
	return &io.PipeReader{}
}


func GijitShadow_NewStruct_PipeWriter() *io.PipeWriter {
	return &io.PipeWriter{}
}


func GijitShadow_InterfaceConvertTo2_ReadCloser(x interface{}) (y io.ReadCloser, b bool) {
	y, b = x.(io.ReadCloser)
	return
}

func GijitShadow_InterfaceConvertTo1_ReadCloser(x interface{}) io.ReadCloser {
	return x.(io.ReadCloser)
}


func GijitShadow_InterfaceConvertTo2_ReadSeeker(x interface{}) (y io.ReadSeeker, b bool) {
	y, b = x.(io.ReadSeeker)
	return
}

func GijitShadow_InterfaceConvertTo1_ReadSeeker(x interface{}) io.ReadSeeker {
	return x.(io.ReadSeeker)
}


func GijitShadow_InterfaceConvertTo2_ReadWriteCloser(x interface{}) (y io.ReadWriteCloser, b bool) {
	y, b = x.(io.ReadWriteCloser)
	return
}

func GijitShadow_InterfaceConvertTo1_ReadWriteCloser(x interface{}) io.ReadWriteCloser {
	return x.(io.ReadWriteCloser)
}


func GijitShadow_InterfaceConvertTo2_ReadWriteSeeker(x interface{}) (y io.ReadWriteSeeker, b bool) {
	y, b = x.(io.ReadWriteSeeker)
	return
}

func GijitShadow_InterfaceConvertTo1_ReadWriteSeeker(x interface{}) io.ReadWriteSeeker {
	return x.(io.ReadWriteSeeker)
}


func GijitShadow_InterfaceConvertTo2_ReadWriter(x interface{}) (y io.ReadWriter, b bool) {
	y, b = x.(io.ReadWriter)
	return
}

func GijitShadow_InterfaceConvertTo1_ReadWriter(x interface{}) io.ReadWriter {
	return x.(io.ReadWriter)
}


func GijitShadow_InterfaceConvertTo2_Reader(x interface{}) (y io.Reader, b bool) {
	y, b = x.(io.Reader)
	return
}

func GijitShadow_InterfaceConvertTo1_Reader(x interface{}) io.Reader {
	return x.(io.Reader)
}


func GijitShadow_InterfaceConvertTo2_ReaderAt(x interface{}) (y io.ReaderAt, b bool) {
	y, b = x.(io.ReaderAt)
	return
}

func GijitShadow_InterfaceConvertTo1_ReaderAt(x interface{}) io.ReaderAt {
	return x.(io.ReaderAt)
}


func GijitShadow_InterfaceConvertTo2_ReaderFrom(x interface{}) (y io.ReaderFrom, b bool) {
	y, b = x.(io.ReaderFrom)
	return
}

func GijitShadow_InterfaceConvertTo1_ReaderFrom(x interface{}) io.ReaderFrom {
	return x.(io.ReaderFrom)
}


func GijitShadow_InterfaceConvertTo2_RuneReader(x interface{}) (y io.RuneReader, b bool) {
	y, b = x.(io.RuneReader)
	return
}

func GijitShadow_InterfaceConvertTo1_RuneReader(x interface{}) io.RuneReader {
	return x.(io.RuneReader)
}


func GijitShadow_InterfaceConvertTo2_RuneScanner(x interface{}) (y io.RuneScanner, b bool) {
	y, b = x.(io.RuneScanner)
	return
}

func GijitShadow_InterfaceConvertTo1_RuneScanner(x interface{}) io.RuneScanner {
	return x.(io.RuneScanner)
}


func GijitShadow_NewStruct_SectionReader() *io.SectionReader {
	return &io.SectionReader{}
}


func GijitShadow_InterfaceConvertTo2_Seeker(x interface{}) (y io.Seeker, b bool) {
	y, b = x.(io.Seeker)
	return
}

func GijitShadow_InterfaceConvertTo1_Seeker(x interface{}) io.Seeker {
	return x.(io.Seeker)
}


func GijitShadow_InterfaceConvertTo2_WriteCloser(x interface{}) (y io.WriteCloser, b bool) {
	y, b = x.(io.WriteCloser)
	return
}

func GijitShadow_InterfaceConvertTo1_WriteCloser(x interface{}) io.WriteCloser {
	return x.(io.WriteCloser)
}


func GijitShadow_InterfaceConvertTo2_WriteSeeker(x interface{}) (y io.WriteSeeker, b bool) {
	y, b = x.(io.WriteSeeker)
	return
}

func GijitShadow_InterfaceConvertTo1_WriteSeeker(x interface{}) io.WriteSeeker {
	return x.(io.WriteSeeker)
}


func GijitShadow_InterfaceConvertTo2_Writer(x interface{}) (y io.Writer, b bool) {
	y, b = x.(io.Writer)
	return
}

func GijitShadow_InterfaceConvertTo1_Writer(x interface{}) io.Writer {
	return x.(io.Writer)
}


func GijitShadow_InterfaceConvertTo2_WriterAt(x interface{}) (y io.WriterAt, b bool) {
	y, b = x.(io.WriterAt)
	return
}

func GijitShadow_InterfaceConvertTo1_WriterAt(x interface{}) io.WriterAt {
	return x.(io.WriterAt)
}


func GijitShadow_InterfaceConvertTo2_WriterTo(x interface{}) (y io.WriterTo, b bool) {
	y, b = x.(io.WriterTo)
	return
}

func GijitShadow_InterfaceConvertTo1_WriterTo(x interface{}) io.WriterTo {
	return x.(io.WriterTo)
}

