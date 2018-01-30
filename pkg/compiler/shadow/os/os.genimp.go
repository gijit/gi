package shadow_os

import "os"

var Pkg = make(map[string]interface{})
func init() {
    Pkg["Args"] = os.Args
    Pkg["Chdir"] = os.Chdir
    Pkg["Chmod"] = os.Chmod
    Pkg["Chown"] = os.Chown
    Pkg["Chtimes"] = os.Chtimes
    Pkg["Clearenv"] = os.Clearenv
    Pkg["Create"] = os.Create
    Pkg["DevNull"] = os.DevNull
    Pkg["Environ"] = os.Environ
    Pkg["ErrClosed"] = os.ErrClosed
    Pkg["ErrExist"] = os.ErrExist
    Pkg["ErrInvalid"] = os.ErrInvalid
    Pkg["ErrNotExist"] = os.ErrNotExist
    Pkg["ErrPermission"] = os.ErrPermission
    Pkg["Executable"] = os.Executable
    Pkg["Exit"] = os.Exit
    Pkg["Expand"] = os.Expand
    Pkg["ExpandEnv"] = os.ExpandEnv
    Pkg["FileInfo"] = GijitShadow_InterfaceConvertTo2_FileInfo
    Pkg["FindProcess"] = os.FindProcess
    Pkg["Getegid"] = os.Getegid
    Pkg["Getenv"] = os.Getenv
    Pkg["Geteuid"] = os.Geteuid
    Pkg["Getgid"] = os.Getgid
    Pkg["Getgroups"] = os.Getgroups
    Pkg["Getpagesize"] = os.Getpagesize
    Pkg["Getpid"] = os.Getpid
    Pkg["Getppid"] = os.Getppid
    Pkg["Getuid"] = os.Getuid
    Pkg["Getwd"] = os.Getwd
    Pkg["Hostname"] = os.Hostname
    Pkg["Interrupt"] = os.Interrupt
    Pkg["IsExist"] = os.IsExist
    Pkg["IsNotExist"] = os.IsNotExist
    Pkg["IsPathSeparator"] = os.IsPathSeparator
    Pkg["IsPermission"] = os.IsPermission
    Pkg["Kill"] = os.Kill
    Pkg["Lchown"] = os.Lchown
    Pkg["Link"] = os.Link
    Pkg["LookupEnv"] = os.LookupEnv
    Pkg["Lstat"] = os.Lstat
    Pkg["Mkdir"] = os.Mkdir
    Pkg["MkdirAll"] = os.MkdirAll
    Pkg["NewFile"] = os.NewFile
    Pkg["NewSyscallError"] = os.NewSyscallError
    Pkg["O_APPEND"] = os.O_APPEND
    Pkg["O_CREATE"] = os.O_CREATE
    Pkg["O_EXCL"] = os.O_EXCL
    Pkg["O_RDONLY"] = os.O_RDONLY
    Pkg["O_RDWR"] = os.O_RDWR
    Pkg["O_SYNC"] = os.O_SYNC
    Pkg["O_TRUNC"] = os.O_TRUNC
    Pkg["O_WRONLY"] = os.O_WRONLY
    Pkg["Open"] = os.Open
    Pkg["OpenFile"] = os.OpenFile
    Pkg["PathListSeparator"] = os.PathListSeparator
    Pkg["PathSeparator"] = os.PathSeparator
    Pkg["Pipe"] = os.Pipe
    Pkg["Readlink"] = os.Readlink
    Pkg["Remove"] = os.Remove
    Pkg["RemoveAll"] = os.RemoveAll
    Pkg["Rename"] = os.Rename
    Pkg["SEEK_CUR"] = os.SEEK_CUR
    Pkg["SEEK_END"] = os.SEEK_END
    Pkg["SEEK_SET"] = os.SEEK_SET
    Pkg["SameFile"] = os.SameFile
    Pkg["Setenv"] = os.Setenv
    Pkg["Signal"] = GijitShadow_InterfaceConvertTo2_Signal
    Pkg["StartProcess"] = os.StartProcess
    Pkg["Stat"] = os.Stat
    Pkg["Stderr"] = os.Stderr
    Pkg["Stdin"] = os.Stdin
    Pkg["Stdout"] = os.Stdout
    Pkg["Symlink"] = os.Symlink
    Pkg["TempDir"] = os.TempDir
    Pkg["Truncate"] = os.Truncate
    Pkg["Unsetenv"] = os.Unsetenv

}
func GijitShadow_NewStruct_File() *os.File {
	return &os.File{}
}


func GijitShadow_InterfaceConvertTo2_FileInfo(x interface{}) (y os.FileInfo, b bool) {
	y, b = x.(os.FileInfo)
	return
}

func GijitShadow_InterfaceConvertTo1_FileInfo(x interface{}) os.FileInfo {
	return x.(os.FileInfo)
}


func GijitShadow_NewStruct_LinkError() *os.LinkError {
	return &os.LinkError{}
}


func GijitShadow_NewStruct_PathError() *os.PathError {
	return &os.PathError{}
}


func GijitShadow_NewStruct_ProcAttr() *os.ProcAttr {
	return &os.ProcAttr{}
}


func GijitShadow_NewStruct_Process() *os.Process {
	return &os.Process{}
}


func GijitShadow_NewStruct_ProcessState() *os.ProcessState {
	return &os.ProcessState{}
}


func GijitShadow_InterfaceConvertTo2_Signal(x interface{}) (y os.Signal, b bool) {
	y, b = x.(os.Signal)
	return
}

func GijitShadow_InterfaceConvertTo1_Signal(x interface{}) os.Signal {
	return x.(os.Signal)
}


func GijitShadow_NewStruct_SyscallError() *os.SyscallError {
	return &os.SyscallError{}
}

