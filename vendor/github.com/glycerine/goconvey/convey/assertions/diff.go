package assertions

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

func Diffb(a string, b string) []byte {

	dirpath := NewSimpleTempDir("diffdir_")
	defer os.RemoveAll(dirpath)

	fa := SimpleTempFile(dirpath)
	fmt.Fprintf(fa, "%s\n", a)
	fa.Close()

	fb := SimpleTempFile(dirpath)
	fmt.Fprintf(fb, "%s\n", b)
	fb.Close()

	co, err := exec.Command("diff", "-b", fa.Name(), fb.Name()).CombinedOutput()
	if err != nil {
		// don't panic, diff returns 2 on differences
	}
	return co
}

func NewSimpleTempDir(prefix string) string {
	dirpath, err := ioutil.TempDir(".", prefix)
	if err != nil {
		panic(err)
	}
	return dirpath
}

func SimpleTempFile(dirpath string) *os.File {

	f, err := ioutil.TempFile(dirpath, "diff_file_")
	if err != nil {
		panic(err)
	}
	return f
}
