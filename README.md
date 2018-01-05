gi: a go interpreter
=======
`gi` aims at being a scalable Go REPL
for doing interactive coding and data analysis.
Currently it is backed by LuaJIT, a tracing Just-in-time
compiler for Lua that provides quite
nice performance, sometimes even beating
ahead-of-time compiled Go code.

status
------
Early stages, work in progress. Contribute!

Currently incremental type checking is applied
to all code. Slices are bounds-checked at runtime.
Functions, closures and slices, as well as
basic expressions compile and run. For-loops
including for-range loops compile and run.

Much is left to do: maps, structs, switch,
interfaces, imports. If this is exciting to
you, contribute! Existing TODOs/open issues and polite improvement
suggestions can be found here
https://github.com/go-interpreter/gi/issues

However, because we are bulding on the fantastic
front end provided by (Gopherjs)[https://github.com/gopherjs/gopherjs], and the fantastic
backend provided by (LuaJIT)[http://luajit.org/], progress has been
quite rapid.

# the dream

Go, if it only had a decent REPL, could be a great
language for exploratory data analysis.

# the rationale
Go has big advantages over python, R, and Matlab.
It has good type checking, reasonable compiled performance,
and excellent multicore support.

# the aim

`gi` aims to provide a library for many
interactive backends to be written against.
`gi` will provide the REPL and do the
initial lexical scanning, parsing, and type
checking. `gi` will then pass the
AST to a backend for codegen and/or immediate
interpretation.

# of course we need a backend to develop against

Considering possible backends for a
reference implementation,
I compared node.js, chez scheme, otto,
gopher-lua, and luajit.

# luajit did what?

Luajit in particular is an amazing
backend to target. In our quick and
dirty 500x500 random matrix multiplication
benchmark, luajit *beat even statically compiled go*
code by a factor of 3x. Go's time was 360 msec.
Luajit's time was 135 msec. Julia uses an optimized
BLAS library for this task and beats both Go
and luajit by multiplying in 6 msec, but
is too immature and too large to be
a viable embedded target.

# installation

~~~
$ go get -t -u -v github.com/go-interpreter/gi/...
$ cd $GOPATH/src/github.com/go-interpreter/gi/cmd/gi
$
$    # first time only, to install the latest luajit.
$    # it will sudo to install a luajit symlink in /usr/local/bin/luajit,
$    # so you may need to provide the root password if you
$    # want that installation, or ctrl-c out if you don't.
$    make onetime
$
$ make install
$
$ gi # start me up
====================
gi: a go interpreter
====================
https://github.com/go-interpreter/gi
Copyright (c) 2018, Jason E. Aten, Ph.D.
License: 3-clause BSD. See the LICENSE file at
https://github.com/go-interpreter/gi/blob/master/LICENSE
====================
  [gi is an interactive Golang environment,
   also known as a REPL or Read-Eval-Print-Loop.]
  [type ctrl-d to exit]
  [type :help for help]
  [gi -h for flag help]
  [gi -q to start quietly]
====================
built: '2018-01-04T22:30:28-0600'
last-git-commit-hash: '6f105a4b7a74509c6117105c908e9a6c38459119'
nearest-git-tag: 'v0.0.8'
git-branch: 'master'
go-version: 'go_version_go1.9_darwin/amd64'
luajit-version: 'LuaJIT_2.1.0-beta3_--_Copyright_(C)_2005-2017_Mike_Pall._http://luajit.org/'
==================
gi> a := []string{"howdy", "gophers!"}

gi> a
table of length 2 is _giSlice{[0]= howdy, [1]= gophers!, }

gi> a[0] = "you rock"

gi> a
table of length 2 is _giSlice{[0]= you rock, [1]= gophers!, }

gi> a[-1] = "compile-time-out-of-bounds-access"
oops: 'problem detected during Go static type checking: 'where error? err = '1:3: invalid argument: index -1 (constant of type int) must not be negative''' on input 'a[-1] = "compile-time-out-of-bounds-access"'

gi> a[100] = "runtime-out-of-bounds-access"
error from Lua vm.Pcall(0,0,0): 'run time error'. supplied lua with: '	_setRangeCheck(a, 100, "runtime-out-of-bounds-access");'
lua stack:
String : 	 [string "..."]:91: index out of range

gi> func myFirstGiFunc(a []string) int {
oops: 'problem detected during Go static type checking: '1:37: expected '}', found 'EOF''' on input 'func myFirstGiFunc(a []string) int {'

gi> // multiline not yet done! :)

gi> func myFirstGiFunc(a []string) int { for i := range a { println("our input is a[",i,"] = ", a[i]) }; return 43 }

gi> myFirstGiFunc(a)
our input is a[	0	] = 	you rock
our input is a[	1	] = 	gophers!
43

gi> // demo type checking

gi> b = []int{1,1}
oops: 'problem detected during Go static type checking: 'where error? err = '1:1: undeclared name: b''' on input 'b = []int{1,1}'

gi> b := []int{1,1}

gi> myFirstGiFunc(b)
oops: 'problem detected during Go static type checking: 'where error? err = '1:15: cannot use b (variable of type []int) as []string value in argument to myFirstGiFunc''' on input 'myFirstGiFunc(b)'

gi> 
~~~

# editor support

An emacs mode `gigo.el` can be found in the `emacs/` subdirectory
here https://github.com/go-interpreter/gi/tree/master/emacs/gigo.el

M-x `run-gi-golang` to start the interpreter. Pressing ctrl-n will
step through any file that is in `gi-golang` mode. 

Other editors: please contribute!

# origin

Author: Jason E. Aten, Ph.D.

License: 3-clause BSD.

Credits: some code here is dervied from the Go standard
libraries, the Go gc compiler,  and from Richard Musiol's excellent Gopherjs project.
This project and those are licensed under the 3-clause BSD license
found in the LICENSE file. The LuaJIT vm and compiler are statically linked
using CGO, and their MIT license can be found in their sub-directories
and online at http://luajit.org/ and https://github.com/LuaJIT/LuaJIT/blob/master/COPYRIGHT
