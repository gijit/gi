gi: go interactive : a front-end for Go REPLs written in Golang
=======

Status: early stages, work in progress. More
ambition than reality at present.

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
$ gi  # start me up

~~~

# origin

Author: Jason E. Aten, Ph.D.

License: 3-clause BSD.

Credits: some code here is dervied from the Go gc compiler
and from Richard Musiol's excellent Gopherjs project.
This project and those are licensed under the 3-clause BSD license
found in the LICENSE file.
