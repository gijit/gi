package shadow_io

import "io"

var Pkg = make(map[string]interface{})

func init() {
	//Pkg["ByteReader"] = io.ByteReader
	//Pkg["ByteScanner"] = io.ByteScanner
	//Pkg["ByteWriter"] = io.ByteWriter
	//Pkg["Closer"] = io.Closer
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
	//Pkg["LimitedReader"] = io.LimitedReader
	Pkg["MultiReader"] = io.MultiReader
	Pkg["MultiWriter"] = io.MultiWriter
	Pkg["NewSectionReader"] = io.NewSectionReader
	Pkg["Pipe"] = io.Pipe
	//Pkg["PipeReader"] = io.PipeReader
	//Pkg["PipeWriter"] = io.PipeWriter
	Pkg["ReadAtLeast"] = io.ReadAtLeast
	//Pkg["ReadCloser"] = io.ReadCloser
	Pkg["ReadFull"] = io.ReadFull
	//Pkg["ReadSeeker"] = io.ReadSeeker
	//Pkg["ReadWriteCloser"] = io.ReadWriteCloser
	//Pkg["ReadWriteSeeker"] = io.ReadWriteSeeker
	//Pkg["ReadWriter"] = io.ReadWriter
	//Pkg["Reader"] = io.Reader
	//Pkg["ReaderAt"] = io.ReaderAt
	//Pkg["ReaderFrom"] = io.ReaderFrom
	//Pkg["RuneReader"] = io.RuneReader
	//Pkg["RuneScanner"] = io.RuneScanner
	//Pkg["SectionReader"] = io.SectionReader
	Pkg["SeekCurrent"] = io.SeekCurrent
	Pkg["SeekEnd"] = io.SeekEnd
	Pkg["SeekStart"] = io.SeekStart
	//Pkg["Seeker"] = io.Seeker
	Pkg["TeeReader"] = io.TeeReader
	//Pkg["WriteCloser"] = io.WriteCloser
	//Pkg["WriteSeeker"] = io.WriteSeeker
	Pkg["WriteString"] = io.WriteString
	//Pkg["Writer"] = io.Writer
	//Pkg["WriterAt"] = io.WriterAt
	//Pkg["WriterTo"] = io.WriterTo

}
