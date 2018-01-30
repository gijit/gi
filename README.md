gijit: a go interpreter
=======
`gijit` aims at being a scalable Go REPL
for doing interactive coding and data analysis.
It is backed by LuaJIT, a tracing Just-in-time
compiler for Lua that provides quite
nice performance, sometimes even beating
ahead-of-time compiled Go code. The REPL
binary is called simply `gi`, as it is a
"go interpreter".

status
------

2018 Jan 29 update
-------
What is left to do: pointers, interfaces, go-routines,
and channels. The channels and go-routines will just wrap the
existing go runtime functionality, using reflection
and Luar.

Release v0.8.0 has command line history available with
the up-arrow/down-arrow. This is a nice supplement to
the `:h` history listing and replay facility.
Like traditional shell editing, we now have ctrl-a (move
to beginning of line); ctrl-e (move to end of line);
ctrl-k (cut); and ctrl-y (paste) at the REPL.

A quick demo of the `:h` history functionality, which
is also new (and distinct from the up-arrow/liner
functionality).

Line history is stored in `$HOME/.gijit.hist`, and is preserved
across `gi` restarts. It can be edited by removing
sets of lines using the `:rm a-b` command. The `:n`
command, where `n` is a number, replays history line `n`.
With a `-` dash, a range of commands to be replayed
is specified. `:10-` replays from 10 to the end of
history, while `:-10` replays everything from the
first line in the history, up to and
including line 10. `:-` replays everything in the
history, because the range endpoints have the intuitive
defaults.

Commands executed in the current session appear after the
`----- current session: -----` delineator.

~~~
$ gi -q
gi> :h
history:
001: greet := "hello gophers!"
----- current session: -----

gi>   ## notice that stuff from past sessions is above the `current session:` line.
gi> :reset
history cleared.
gi> :h
history: empty
----- current session: -----
gi> a := 1

elapsed: '17.614µs'
gi> :h
history:
----- current session: -----
001: a := 1

gi> b := a * 2

elapsed: '50.054µs'
gi> :h
history:
----- current session: -----
001: a := 1
002: b := a * 2

gi> :1-
replay history 001 - 002:
a := 1
b := a * 2


elapsed: '41.932µs'
gi> :h
history:
----- current session: -----
001: a := 1
002: b := a * 2
003: a := 1
004: b := a * 2

gi> :3-4
replay history 003 - 004:
a := 1
b := a * 2


elapsed: '76.605µs'
gi> :-2
replay history 001 - 002:
a := 1
b := a * 2


elapsed: '91.664µs'
gi> :-     ## replay everything in our history
replay history 001 - 008:
a := 1
b := a * 2
a := 1
b := a * 2
a := 1
b := a * 2
a := 1
b := a * 2


elapsed: '37.896µs'
gi> :h
history:
----- current session: -----
001: a := 1
002: b := a * 2
003: a := 1
004: b := a * 2
005: a := 1
006: b := a * 2
007: a := 1
008: b := a * 2
009: a := 1
010: b := a * 2
011: a := 1
012: b := a * 2
013: a := 1
014: b := a * 2
015: a := 1
016: b := a * 2

gi> :rm 3-    ## if a range leaves off an endpoint, it defaults to the beginning/end.
remove history 003 - 016.
gi> :h
history:
----- current session: -----
001: a := 1
002: b := a * 2

gi> :rm -    ## same as :reset or :clear
remove history 001 - 002.
gi> :h
history: empty
----- current session: -----
gi> 
~~~

Release v0.7.9 has slice `copy` and `append` working.
`copy` allows source and destination
slices to overlap, adjusting
the copy direction automatically.

In release v0.7.7, make applied to slices works. For example,
`make([]int, 4)` will create a new length four array filled
with the zero value for `int`, and then take a slice of that array.

In release v0.7.6, taking a slice of an array works.

2018 Jan 28 update
-------
In release v0.7.5, we added the ability for the `gi` REPL
to import "regexp" and "os". More generally, we added
the general ability
to quickly add existing native Go packages to be imported.

The process of making an existing native Go package available
is called shadowing.

The shadowing process is fairly straightforward.

A new utility, `gen-gijit-shadow-import`
is run, passing as an argument the package to be shadowed. The utility
produces a new directory and file under `pkg/compiler/shadow`. Then
a few lines must be added to `pkg/compiler/import.go` so
that the package's exported functions will be available
to the REPL at runtime, after import. The "regexp" and "os"
shadow packages provide examples of how to do this.

Some non-functions in the shadow package
may need to be manually commented out in
the generated file to finish preparation. The Go compiler
will tell you which ones when you rebuild `gi`.

While shadowing (sandboxing) does
not allow arbitrary imports to be called from inside
the `gi` without prior preparation, this is
often a useful and desirable security feature. Moreover
we are able to provide these imports without using Go's
linux-only DLL loading system, so we remain portable/more
cross-platform compatible.

Interfaces are not yet implemented, and so are not yet imported.

2018 Jan 27 update
-------
In release v0.7.3, arrays are passed to Go native functions, and
array copy by value is implemented. Having arrays work well, now
we are well-positioned to get slices pointing to arrays working.
That will be the next focus.


2018 Jan 25 update
-------
In release v0.7.2, for range over strings produces utf8 runes.

Release v0.7.1 makes all integers work as map keys.
Go integers (int, int64, uint, etc) are represented with
cdata in LuaJIT, since otherwise all numbers default
to float64/doubles.

This was a problem because cdata are boxed, so equality
comparison on ints was comparing their addresses. Now we
translate all Go integer map keys to strings in Lua, which makes
key lookup for a Go integer key work as expected.


2018 Jan 23 update
------
Release v0.7.0 brings local variables inside
funtions, by default. The global package name space is not
changed by variables declared inside functions.
Function code will be much faster. Global level/package
level declarations are kept global, and are not
prefixed with the 'local' keyword in the Lua translation.

Release v0.6.8 brought the ability to update a method
definition at the REPL, replacing an old method
definition with a new one.  Functions and variables
could already be re-defined. Now all data types can
be updated as you work at the `gi` REPL.

2018 Jan 21 update
------
Release v0.6.7 has working `fmt.Printf`.

~~~
gi> import "fmt"
import "fmt"
gi> fmt.Printf("Hello World from gi!")
Hello World from gi!
gi> 
~~~

2018 Jan 20 update
------
Release v0.6.2 has a working `fmt.Sprintf`.

This needed handling vararg calls into
compiled Go code, which is working now.

A quick demo:

~~~
gi> import "fmt"
import "fmt"
gi> fmt.Sprintf("hi gi") // value not printed at the REPL (at present), so add println:

gi> println(fmt.Sprintf("hi gi"))
hi gi

gi> println(fmt.Sprintf("hi gi %v", 1))
hi gi 1

gi> println(fmt.Sprintf("hi gi %v", "hello"))
hi gi hello

gi> println(fmt.Sprintf("hi gi %v %v", "hello", 1))
hi gi hello 1

gi> 
~~~

2018 Jan 18 update
------
Release v0.6.0 was aimed at putting the infrastructure in place
to support package imports. Specifically, we aimed at getting
the first `import "fmt"` and the first use of `fmt.Sprintf`
from the REPL.

The varargs and int64 handling required in `fmt.Sprintf` made this extra
tricky. And so it's not quite done.

Nonetheless, I'm releasing v0.6.0 because there were big
refactorings that provide significant new internal
functionality, and contributors will want to leverage these.

The REPL will now accept `import "fmt"` and will wire in
the `fmt.Sprintf` for you. It's hardwired for now. Auto
loading of functions from already compiled packages will
come later. Standing on the shoulders of another giant,
the `luar` package, lets
us call pre-compiled Go packages from Lua through reflection.

A quick demo:
~~~
gi> import "fmt"
import "fmt"

gi> a:=fmt.Sprintf("hello gi!")
a:=fmt.Sprintf("hello gi!")

gi> a
a
hello gi!

gi>
~~~

`luar` (https://github.com/stevedonovan/luar) is a mature
library that gives us the basis for imports.
Since `LuaJIT` provides reasonable int64 handling, we've
extended the `luar` functionality to gracefully convey
`int64` and `uint64` values. There's more extension to
do for `int32`, `int16`, `int8`, `uint32`, etc but
this should be straightforward extension of the new
functionality.

While the vararg handling to make Sprintf
actually useful beyond just the format
string is missing, this should be done shortly.
The red 051 and 052 tests in imp_test.go track the
last bits of functionality needed to make Sprintf work.

API functions
luajit_push_cdata_uint64(), luajit_push_cdata_int64(),
and luajit_ctypeid() were added to the luajit API
to support passing int64/uint64 values from Go to Lua and
back to Go without the loss of data formerly associated
with casting to double (float64) and back.


2018 Jan 13 update
------
We've moved within github to make admin easier. We are now at https://github.com/gijit/gi

len(a) now displays at the REPL. Fixes #22.
This was a minor REPL nit, but well worth addressing.

We have a new contributor! Welcome to Malhar Vora.

As of v0.5.6, integer modulo and divide now
check for divide by zero, and panic like Go when found.
The math.lua library is supplemented with
math.isnan(), math.finite(), and __truncateToInt().

There's a new section of this readme,
https://github.com/gijit/gi#translation-hints
that gives specific hints for porting the javascript ternary
operator and other constructs.

jea: I'll be offline for a day or two.


2018 Jan 12 update
------
Switch statements now work. They work at the top level
and inside a function.


2018 Jan 11 update
------
With release v0.5.0 the inital test of the `defer`/
`panic`/`recover` mechanism passes. Woot!  There's more
to do here, but the design is solid so filling in
should be quick.

For a stack-unwinding `panic`, we use
what Lua offers, the `error`
mechanic -- to throw -- combined with the
`xpcall` mechanic to catch.

The only limitation I found here is on recursive `xpcall`: if
you are in a panic stack unwind, and then in a defer function,
and your code causes a second error that is *not* a deliberate panic,
then that error will be caught but recover won't return
that error value to the caller of recover. This is due to a wierd
corner case in the implementation of LuaJIT where
it does not like recursive `xpcall` invocations, and
reports "error in error handling".

I asked on the Lua and LuaJIT mailing lists about
this, and posted on Stack Overflow. So far no
replies. https://stackoverflow.com/questions/48202338/on-latest-luajit-2-1-0-beta3-is-recursive-xpcall-possible

It's a fairly minor limitation, and easy to work
around once you notice the bug: just call `panic`
directly rather than causing the error. Or don't
cause the error at all(!) Or simply use a
different side-band to pass around the value.
Lots of work arounds.


2018 Jan 10 update
------
With release v0.4.1 we have much improved map support.
In `gi`, maps now work as in Go. Multiple-valued queries
properly return the zero-value for the value-type
when the key is missing, and return the 2nd
value correctly. Nil keys and values are
handled properly. `delete` on a map works as
expected, and maintains the `len` property.

2018 Jan 9 update
------
Functions and methods can now be re-defined at the REPL. The
type checker was relaxed to allow this.

We changed over from one LuaJIT C binding to anther. The
new binding is the same one that LuaR uses, so this
enables LuaR exploration.

2018 Jan 8 update
------
Today we landed multiline support. We evalutate
Go expressions as they are entered, and these
can now span multiple lines. This lifts the
prior limitation that meant that functions
and types needed to be defined all on one line.

This was fun to put together. I used the actual gc front end that
parses regular go code. Since gc is written in Go,
why not leverage it! The advantage is that we know we
are building on correct parsing of the whole language.

Of course minor tweaks had to be made to allow statements and
expressions at global scope. Happily, from our experience
adding the same provisions to GopherJS, we knew these
were relatively minor changes. See the updated
demo transcript below in this readme for a multi-line
function definition taking effect.



2018 Jan 7: latest update
------
Today we acheived passing (light) tests for method definition and invocation!

Also a significant discovery for the object system: Steve Donovan's Luar
provides object exchange both ways between Go -> Lua and Lua -> Go.

That should influence our design of our Go source -> Lua source mapping. If we
map in a way that matches what Luar does when it translates from
Go binary -> Lua binary, then our objects will translate cleanly
into binary Go calls made by reflection.

Even more: Luar provides access to the full Go runtime and channels
via reflection. Nice! We don't have to reinvent the wheel, and we
get to use the high-performance multicore Go scheduler.


earlier summary
-----
Early stages, work in progress. Contribute!

Currently incremental type checking is applied
to all code. Slices are bounds-checked at runtime.
Functions, closures and slices, as well as
basic expressions compile and run. For-loops
including for-range loops compile and run.

If this is exciting to
you, contribute! Existing TODOs/open issues and polite improvement
suggestions can be found here
https://github.com/gijit/gi/issues

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

We want to provide one excellent integrated REPL for Go.
Exploratory data analysis should not be hampered
by weak type-checking or hard-to-refactor code,
and performance should not suffer just because
you require interaction with your data.


# of course we need a backend to develop against

Considering possible backends,
I compared node.js, chez scheme, otto,
gopher-lua, and LuaJIT.

# LuaJIT did what?  Will Golang run on GPU?

LuaJIT in particular is an amazing
backend to target. In our quick and
dirty 500x500 random matrix multiplication
benchmark, LuaJIT *beat even statically compiled go*
code by a factor of 3x. Go's time was 360 msec.
LuaJIT's time was 135 msec. Julia uses an optimized
BLAS library for this task and beats both Go
and LuaJIT by multiplying in 6 msec, but
is too immature and too large to be
a viable embedded target.

Bonus: LuaJIT has Torch7 for machine learning.
And, Torch7 has GPU support. [1][2]

[1] https://github.com/torch/cutorch

[2] https://github.com/torch/torch7/wiki/Cheatsheet#gpu-support

Will golang (Go) run on GPUs?  It might be possible!

# installation

Works on Mac OSX and Linux. On windows: theoretically it should work on windows, I have not worked out what flags are needed. One will need to install a C compiler on windows and work out the right compiler flags to make CGO build and link LuaJIT into the Go `gi` binary.

~~~
$ go get -t -u -v github.com/gijit/gi/...
$ cd $GOPATH/src/github.com/gijit/gi && make
$
$ ... wait for gi build to finish, it builds LuaJIT
$     using C, so it takes ~ 20 seconds to install `gi`.
$
$ gi # start me up (will be in $GOPATH/bin/gi now).

====================
gi: a go interpreter
====================
https://github.com/gijit/gi
Copyright (c) 2018, Jason E. Aten. All rights reserved.
License: 3-clause BSD. See the LICENSE file at
https://github.com/gijit/gi/blob/master/LICENSE
====================
  [ gi is an interactive Golang environment,
    also known as a REPL or Read-Eval-Print-Loop ]
  [ type ctrl-d to exit ]
  [ type :help for help ]
  [ gi -h for flag help ]
  [ gi -q to start quietly ]
====================
built: '2018-01-08T23:46:07-0600'
last-git-commit-hash: 'db302d2acb37d3c2ba2a0d376b6f233045928730'
nearest-git-tag: 'v0.3.3'
git-branch: 'master'
go-version: 'go_version_go1.9_darwin/amd64'
luajit-version: 'LuaJIT_2.1.0-beta3_--_Copyright_(C)_2005-2017_Mike_Pall._http://luajit.org/'
==================
using this prelude directory: '/Users/jaten/go/src/github.com/gijit/gi/pkg/compiler'
using these files as prelude: array.lua, map.lua, prelude.lua, slice.lua, struct.lua
gi>

gi> a := []string{"howdy", "gophers!"}

gi> a   // ^^ make data using Go's literals. inspect it by typing the variables name.
slice of length 2 is _giSlice{[0]= howdy, [1]= gophers!, }

gi> a[0]  = "you rock" // data can be changed

gi> a
slice of length 2 is _giSlice{[0]= you rock, [1]= gophers!, }

gi> // the Go type checker helps you quickly catch blunders, at compile time.

gi> a[-1] = "compile-time-out-of-bounds-access" 
oops: 'problem detected during Go static type checking: 'where error? err = '1:3: invalid argument: index -1 (constant of type int) must not be negative''' on input 'a[-1] = "compile-time-out-of-bounds-access" 
'

gi> // runtime bounds checks are compiled in too:

gi> a[100] = "runtime-out-of-bounds-access"
error from Lua vm.Pcall(0,0,0): 'run time error'. supplied lua with: '	_gi_SetRangeCheck(a, 100, "runtime-out-of-bounds-access");'
lua stack:
String : 	 ...rc/github.com/gijit/gi/pkg/compiler/prelude.lua:14: index out of range

gi> // We can define functions:

gi> func myFirstGiFunc(a []string) int {
>>>    for i := range a {
>>>      println("our input is a[",i,"] = ", a[i]) 
>>>    };
>>>    return 43
>>> }
func myFirstGiFunc(a []string) int {

	for i := range a {

		println("our input is a[", i, "] = ", a[i])

	}

	return 43

}
gi> myFirstGiFunc(a)
our input is a[	0	] = 	you rock
our input is a[	1	] = 	gophers!

gi> // ^^ and call them. They are tracing-JIT compiled on the LuaJIT vm.

gi> // more compile time type checking, because it rocks:

gi> b := []int{1,1}

gi> myFirstGiFunc(b)
oops: 'problem detected during Go static type checking: 'where error? err = '1:15: cannot use b (variable of type []int) as []string value in argument to myFirstGiFunc''' on input 'myFirstGiFunc(b)
'
gi>
~~~

# editor support

An emacs mode `gigo.el` can be found in the `emacs/` subdirectory
here https://github.com/gijit/gi/tree/master/emacs/gigo.el

M-x `run-gi-golang` to start the interpreter. Pressing ctrl-n will
step through any file that is in `gi-golang` mode. 

Other editors: please contribute!

# how to contribute: getting started

a) Pick an issue from here, https://github.com/gijit/gi/issues, and add a comment that
you are starting work on that feature. Make a branch for your feature, using `git checkout -b yourFeatureName`.

b) Write a test for your feature. Make sure it fails (the test is red), before
moving on to implementation. Tests are quite short. There are many examples are here,
which show the currently implemented features. Add your test at the end of
the compiler/repl_test.go file.

https://github.com/gijit/gi/blob/master/pkg/compiler/repl_test.go

Then simply implement your feature. (So simple! Yeah right!)

So this is hard part. It's too situational to give general advice, but do see the hints
https://github.com/gijit/gi#translation-hints below for some
specific Lua tricks for translating javascript idioms.

These are the main files you'll be adding to/updating:

https://github.com/gijit/gi/blob/master/pkg/compiler/incr.go

https://github.com/gijit/gi/blob/master/pkg/compiler/package.go

https://github.com/gijit/gi/blob/master/pkg/compiler/translate.go

https://github.com/gijit/gi/blob/master/pkg/compiler/statements.go

https://github.com/gijit/gi/blob/master/pkg/compiler/expressions.go

https://github.com/gijit/gi/blob/master/pkg/compiler/luaUtil.go

You wil find it necessary and useful to add print statements to the code. Do this using the `pp()` function, and feel free to leave those prints in during commit. It's a small matter later to take them out, and while you are adding functionality, the debug prints help immensely. You will see them scattered through the code as I've worked. Just leave them there; the verb.Verbose flag and VerboseVerbose can be used to mute them.

The files above derive from GopherJS which compiles Go into Javascript; whereas `gi` translates Go into
Lua. This makes implementation usually very fast, since mostly it is just
above figuring out how to re-write javascript into Lua. You are typically just
checking the syntax of the source-to-source translation. Sometimes some
Lua support functions will be needed. Add them to a new .lua file in `compile/`
directory.

By default, `gi` looks in `./prelude/` relative to its current directory, and this is symlinked to `pkg/compile/` if you are running in `cmd/gi`. Otherwise use the `-prelude` flag to `gi` to tell it where to find its prelude files. All
.lua files found the prelude directory will be sourced during `gi` startup. The default prelude is the `pkg/compile` directory. These files are required for `gi` to work.

c) When you are done, make sure all the tests are green `go test -v` in the compile/ directory.
Run `go fmt` on your code.

d) submit your pull request! (Rebase against master first, please).

# Lua resources

LuaJIT targets Lua 5.1 with some 5.2 extensions.

a) main web site

https://www.lua.org/

b) Programming in Lua by by Roberto Ierusalimschy, the chief architect of Lua.

1st edition. html format (Lua 5.0) https://www.lua.org/pil/contents.html

2nd edition. pdf format (Lua 5.1) https://doc.lagout.org/programmation/Lua/Programming%20in%20Lua%20Second%20Edition.pdf

c) Lua 5.1 Reference Manual, by R. Ierusalimschy, L. H. de Figueiredo, W. Celes
Lua.org, August 2006 

Lua 5.1 https://www.lua.org/manual/5.1/ 

# translation hints

specific javascript to Lua translation hints:

d1) the ternary operator
~~~
x ? y : z
~~~
should be translated as
~~~
( x and {y} or {z} )[1]
~~~

d) the comma operator
~~~
x = (a, b, c) // return value of c after executing `a` and `b`
~~~
doesn't have a direct equivalent in Lua. Try
to see if you can't define a new function in
the prelude to take care of the same processing
that a,b,c does.

If `b` doesn't refer to `c` directly, and `a` doesn't
refer to `b` directly, then
~~~
x = {a, b, c}[3]
~~~
comes close. Rarely does such a construct arise,
since `a` and `b` are typically helper computations
to compute `c`. However that boxing-unboxing construct is
helpful in some tight corners, and may be your
fastest alternative.

Compared to defining and then calling a new closure,
boxing and unboxing is 100x faster.

# origin

Author: Jason E. Aten

License: 3-clause BSD.

Credits: some code here is dervied from the Go standard
libraries, the Go gc compiler,  and from Richard Musiol's excellent Gopherjs project.
This project and those are licensed under the 3-clause BSD license
found in the LICENSE file. The LuaJIT vm and compiler are statically linked
using CGO, and their MIT license can be found in their sub-directories
and online at http://luajit.org/ and https://github.com/LuaJIT/LuaJIT/blob/master/COPYRIGHT
See the subdirectories of vendored and utilized libraries for their
license details.
