gijit: a just-in-time trace-compiled golang interpreter
=======
`gijit` aims at being a scalable Go REPL
for doing interactive coding and data analysis.
It is backed by LuaJIT, a Just-in-Time
trace compiler. The REPL binary is called
simply `gi`, as it is a "go interpreter".

quick install
-------------

~~~
# use go1.10 or later.
$ go get -d github.com/gijit/gi/cmd/gi
$ cd $GOPATH/src/github.com/gijit/gi
$ (On posix/mac/linux: run `./posix.sh` to build libluajit.a)
$ (On windows: run `windows.bat` to build libluajit.a; see https://github.com/gijit/gi/issues/18
   for notes on installing both mingw64 and make, which are pre-requisites.)
$ make install
$ gi
~~~
See https://github.com/gijit/gi/issues/18 for windows install help.

For v1.3.2, there are pre-compiled binaries here https://github.com/gijit/gi/releases/tag/v1.3.2 They have the prelude compiled in now. They should run standalone, without needing to install the source. Note that the importing of other packages (e.g. `fmt`; use `println` instead) is not yet functional.

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

# Q: Can I embed `gijit` in my app?

A: Yes! That's how `gijit` is designed. In fact the `gi` command is just
a very thin wrapper around the `pkg/compiler` library. See

https://github.com/gijit/gi/blob/master/cmd/gi/repl.go#L9

and

https://github.com/gijit/gi/blob/master/pkg/compiler/repl_luajit.go#L57

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
a viable target.

Bonus: LuaJIT has Torch7 for machine learning.
And, Torch7 has GPU support. [1][2]

[1] https://github.com/torch/cutorch

[2] https://github.com/torch/torch7/wiki/Cheatsheet#gpu-support

Will golang (Go) run on GPUs?  It might be possible!

# installation from source

Release v1.3.0 and later needs go1.10 or later (a critical Windows timezone workaround was provided in go1.10). Works on Mac OSX and Linux and Windows. To build on Windows, you'll need to install mingw64 port of gcc first, if its not already installed, since that is what CGO on windows requires. See the notes in https://github.com/gijit/gi/issues/18

[Update: I put a v1.3.2 binaries here https://github.com/gijit/gi/releases/tag/v1.3.2 ; compiled on Windows10.]

~~~
$ go get -d github.com/gijit/gi/cmd/gi
$ cd $GOPATH/src/github.com/gijit/gi
$ (On posix/mac/linux: run `./posix.sh` to build libluajit.a)
$ (On windows: run `windows.bat` to build libluajit.a ; note pre-reqs below.)
   
$ make install
$ gi
~~~
See https://github.com/gijit/gi/issues/18 for windows install help.
Install both mingw64 and make before building gijit. These are prerequisites.

most recent status
------------------

2018 March 2 update
-------------------
Release v1.3.2 vendors the gonum libraries
used, as users were seeing some version skew.

There are binaries available. See https://github.com/gijit/gi/releases/tag/v1.3.2

2018 March 1 update
-------------------
A major milestone, the v1.2.x series and the
latest release bring
fully interactive goroutines
to the REPL. The REPL can perform receives
on unbuffered channels, and interact with background
goroutines. A background Lua coroutine
runs a scheduler that coordinates.

Full blocking at the REPL, on a select or
receive that cannot be finished at this
time, is not yet implemented.

Importantly, native Go imports are turned off while we
work on polishing the goroutine system.
Hence `import "fmt"` won't work.

As of v1.2.4, close() on channels is available.
As of v1.2.5, <-ch, for a channel ch, works
by itself at the repl and in functions.

2018 Feb 26 update
------------------
With release v1.1.0 we focused on establishing
an all-Lua goroutine functionality. Rather than
tackle both all-Lua and hybrid Lua/native
goroutines all at once, we focused on getting
goroutines working completely in Lua land.

Release v1.1.0 acheives that goal, making
full `go`, `select`, channel send, and channel
receive available. v1.1.1 provides a little
polish.

2018 Feb 24 update
------------------
Release v1.0.15 brings the ability to make native
Go channels via reflection for basic types, and then
use those channels in send, receive, and `select`
operations to communicate between interpreted and
pre-compiled Go code.

2018 Feb 23 update
------------------
Release v1.0.7 fixes a bunch of subtle corner
cases in the type system implementation, making
it much more robust.

2018 Feb 21 update
------------------
Release v1.0.2 brings a calculator mode, for
multiple lines of direct computation. Enter
calculator mode with `==` alone on a line,
and exit with `:` alone.

Release v1.0.0 marks these milestones:

* the full program _bench/mul.go (mat_test 501) builds and runs;
  with no imported packages used.
* `gijit` has much improved complex number support.
* `gijit` runs on Windows now, as well as OSX and Linux.

Limitations not yet addressed:

* import of source packages isn't working yet.
* imports of binary Go packages don't work, and in general need some help.
* chan/select/goroutines are not implemented; init()
  functions and full ahead-of-time compile of packages
  are not done.

# overview

    In the form of Q & A:

~~~
On Friday, February 9, 2018 at 12:48:17 PM UTC+7,
Christopher Sebastian wrote on golang-nuts:

gijit looks really interesting.  Do you
have a high-level description/diagram
anywhere that gives an overview of how
the system works? [0]

I took a brief look at the git repo,
and it seems like you're translating
Go into Javascript, and then to Lua,
and then using LuaJIT. [1] Is that right?
[2] How do you manage the complexity
of these seemingly-unwieldy pieces
and translation layers?

It seems like you have your own 'parser'
implementation (maybe copied from the Go src tree? [3]).
How will you keep your parser (as well as the
rest of the system) up-to-date as Go changes? [4]
Is it a very manual process? [5]

Thanks for working on this important project!
~Christopher Sebastian
~~~
My reply, condensed:
~~~
[0] Overview/how does the system work? How is `gijit` differ from GopherJS?

a) GopherJS does:  Go source  ->  Javascript source.

It translates a whole program at a time.

while

b) `gijit` does:

     i) Go source -> Lua source.

But in small chunks; for example, line-by-line.

Instead of whole program, each sequentially
entered expression/statement/declaration,
(or using :source, a larger fragment
if required for mutually recursive definitions)
is translated and evaluated. These can be
multiple lines, but as soon as the
gc syntax front-end tells us we have a
complete expression/statement/declaration, then
we immediately translate and evaluate it.

Yes, the semantics will be subtly different.
If you've used R/python/matlab
to develop your algorithms before
translating them to a compiled language
(a highly effective practice that I
recommend, and one that motivated gijit's development),
you'll know that this is of little concern.

    ii) Then run that Lua source on LuaJIT.

So while the front end of gijit is derived
from GopherJS, there is no javascript source generated.

There's still quite a bit of javascript
remnants in the source, however, as many corner
cases haven't been tackled. So that might
have misled a cursory inspection.

I have adopted the practice, in the GopherJS
derived code, of leaving in the
javascript generation until I hit the
problem with a new test. This lets me
visually see, and experience when run,
where I haven't handled a case yet.
There are many corner cases not yet handled.


[2]. Very carefully, with full Test-driven development. :)

Seriously though, TDD is fantastic (if
not essential) for compiler work.

I would add, I'm working with extremely
mature building blocks.
So if there's a bug, it is almost surely
in my code, not the libraries I'm using.

The three main libraries are already
themselves the highly
polished result of very high-quality
(German and Swiss!) engineering.

a) Mike Pall's LuaJIT is the work of 12
years by a brilliant designer and engineer.
Cloud Flare is now sponsoring King's
College London to continue LuaJIT's development.

b) Richard Musiol's GopherJS is the
work of some 5 years. It passes most of
the Go standard library tests.

c) I don't know how long it took Robert Griesemer
and the Go team to write the Go front-end
parser and type checker, but there is a ton
of time and effort and testing there.

And a significant 4th:

d) Luar (for Go <-> Lua binary object exchange,
so that binary imports of Go libraries work)
is a mature library. A great deal of work was
done on it by Steve Donovan and contributors.

This is all work that doesn't need to be
reinvented, and can be leveraged.

The slight changes to the above (vendored and modified)
libraries in order to
work at the REPL are light relaxations of
a small set of top-level checks.

[3]. I re-use the front-end syntax/parser
from the gc compiler only in order to quickly
determine if we have a complete
expression at the repl, or if we need to
ask for another line of input.

Then I throw away that parse and use
the standard library's go/ast, go/parser,
go/types to do the full parse and type-checking.

This is necessary because go/parser and
go/types are what GopherJS is built around.

[4]. Go 1 is very mature, and changes very
slowly, if at all, anymore, between releases.

[5]. Yes, it is manual. I note the changes
to the parser and typechecker are important
but quite small, and would rebase easily.

Best wishes,

Jason
~~~

I would add credit to Lua's originators
and contributors from around the world.

Roberto Ierusalimschy et al's design
and evolution of Lua over the last
25 years make it a (perhaps surprisingly)
great tool for this purpose.
It certainly surprised me.

Lua's primary goal of acting as
an embedded scripting language for
many disparate host languages with
different inheritance/prototype/overloading
semantics has sculpted it into
a power tool that fits the job like
a glove. It is very good at
language implementation.


status update history
------
See the top of this README for the latest update.

2018 Feb 20 update
------------------
v0.9.19 has great progress getting full program to run.
A crude matrix multiplication whole program (mat_test 501)
now executes correctly, albeit slowly and untuned.

~~~
`runMultiply(500,9,9) -> 301609258250` from `_bench/mul.go`/`_bench/mul.lua`
(a) took 342.456726ms  on Go
(b) took 5-7 seconds     on gijit (about 20x slower)
~~~

See https://github.com/stevedonovan/luar issue #23 to
follow some tuning -- we omit Luar's type() override
with `ProxyType()` to retain LuaJIT performance on
the matrix 501 example.

2018 Feb 18 update
------------------
As of `gigit` v0.9.15, Complex numbers
fully supported, with all of the `cmplx` library
functions available. The underlying complex.lua
library works at LuaJIT natively.

There's still some cruft to clear out
in the compiler from all the javascript
support that is no longer needed. So
complex number interactive use is still
a little awkward, but this can be
quickly improved.

The revamped type system supports the
matrix multiplication benchmark.

Building on Windows, alongside OSX
and Linux, now works.

2018 Feb 14 update
------------------
`gijit` was successfully built on Windows10,
and this same approach will probably work on
earlier Windows versions.

See https://github.com/gijit/gi/issues/18
for my notes on getting mingw64 installed, and
see the `windows` branch of this repo.

OpenBLAS and sci-lua are vendored for matrix
operations. The docs are http://scilua.org/.
Also from sci-lua, some benchmarks showing
LuaJIT can be as fast as C. LuaJIT is typically faster
than Julia 0.4.1.

2018 Feb 12 update
------------------
Release v0.9.12 works under both go1.9.4 and go1.9.3, but see the
new/revised installation instructions.

Actually we recommend avoiding go1.9.4. It is
pretty broken/borked.
Use go1.9.3 and wait for
go1.9.5 before upgrading.

The installation instructions are now slightly different, so
that all actual building is done under make, where we
can set the required environment variables.
~~~
$ go get -d github.com/gijit/gi/cmd/gi
$ cd $GOPATH/src/github.com/gijit/gi && make install
~~~

2018 Feb 09 update
------------------
v0.9.11 allows building under go1.9.4.

v0.9.9 restores building under go1.9.3 after attempts
to get go1.9.4 to work messed up the build.

v0.9.7 attempted to build under the newly release go1.9.4,
but failed to do so and has been replaced.

2018 Feb 08 update
------------------
v0.9.6 fixes #20. The type checker now allows
cleanly re-defining struct types at the REPL.


2018 Feb 08 update
------------------

`gijit` v0.9.4 feature summary:

* interactively code in Go. Just-in-time for Valentine's Day, a Go REPL that is not based on re-compiling everything after every line.

* the ability to import binary Go packages. Call into native Go code from the REPL.

* use Go as a calculator. Just start the line with `=`.
  Use `==` to enter calculator mode where multiple
  direct math expressions can be evaluated. Return
  to go mode with `:`.
 
* structs, interfaces, pointers, defer are all available.

* current limitation: no `go`/`select`/`chan` implementation.

* portable. Doesn't depend on Go's plugin system, or on
recompiling everything every time. We run on OSX and Linux and Windows.



# demo

~~~
$ gi -q
gi> import "math"

elapsed: '15.155µs'
gi> import "fmt"

elapsed: '16.507µs'
gi> fmt.Printf("hello expressions! %v\n", math.Exp(-0.5) * 3.3 - 2.2)
hello expressions! -0.1984488229483099

elapsed: '141.721µs'
gi> = 4.0 / 1.3 /* or, with '=', no fmt.Printf needed. */
3.0769230769231

elapsed: '112.18µs'
gi>  
~~~


2018 Feb 07 update
-------
In version v0.9.3, defers pass additional
tests (that found issues that were fixed),
and the repl in raw mode can
import Go binary libraries with
the call __go_import(path).

Raw mode factilitates system debugging and tests.
Raw mode isn't needed by
end users, unless one is developing gijit. Raw
mode allows direct LuaJIT commands to be entered. It is 
accessed with the ':r' enter-raw-mode command; and
':' returns one to Go mode.
~~~
gi> :r
Raw LuaJIT language mode.

elapsed: '20.042µs'
raw luajit gi> __go_import "fmt"

elapsed: '2.048215ms'
raw luajit gi> fmt.Printf("hello Go!")
hello Go!
elapsed: '124.64µs'
raw luajit gi> :
Go language mode.
gi> 
~~~

In version v0.9.2, the REPL prints expressions that
produce multi-lines of Lua better. We only wrap
the final line with a print. This handles expressions
that generate anonymous pointer types gracefully.

In version v0.9.1, pointers inside structs work.

In version v0.9.0, defer handling of named
return values received some important correctness fixes.

2018 Feb 06 update
-------
In version v0.8.9, the repl received some refactoring
to make it easier to test.

In version v0.8.8 (quiet) and v0.8.7 (debug prints live),
pointer support is much improved. Cloning is restored,
and test 028 is green.

Basic assignment to pointers, and assignment through
pointers, work now. For example, the sequence

~~~
 a:= 1
 b := &a
 c := *b
 *b = 3
~~~

works as expected (ptr_test.go/Test099). As usual,
after the four statements, `c` ends as 1,
and `a` ends as 3.

Still TODO is supporting pointer members
within structs, but that should follow shortly.


2018 Feb 05 update
-------
Excellent progress.

In version v0.8.6, the majority of the type system
from GopherJS was ported over to LuaJIT in the
struct.lua file. Type assertions on interfaces
are working, cf face_test.go and tests 100, 102, 202.

The cloning of structs (test 028 in repl_test.go) is
temporarily broken while their infrastructure is being
refactored to use the new system.


2018 Feb 02 update
-------
In version v0.8.4, we began integrating the
GopherJS type system, in order to properly support
interfaces. It turns out that, because the Go
reflect system is incomplete, GopherJS and
now `gijit` both need a complete, stand-alone type system
implementation. So adding interface support
is a much bigger job than I orignally thought.

Nonetheless, with GopherJS lighting the way, we
are making progress with the port. Test 100
in face_test.go was red for a long time
while we began integrating our updated
Lua-metable based object model for types with the
GopherJS object model for types. The
integration isn't finished yet, as many
of the properties that live in the leaf
table of a new struct need to be moved
up in the properties table, but we've
got the basic machinery going and now
its simply a matter of fine tuning.

Test 100 finally went green, so we felt it
was time to mark progress with a release.

However other tests remain red, so v0.8.4 will be
an internal-only release.

In version v0.8.3, an internal release, we restrict `gijit` programs
to a subset of legal `go` programs, by imposing
some mild restrictions on the names of variables.

Minor restriction number one: variable names cannot start
with '__' two underscores.

Minor restriction number two: in `gijit`, you can't
have a variable named `int`, or `float64`, for example. These
are names of two of the pre-declared numeric types in Go.

So, while
~~~
func main() {
	var int int
        _ = int
}
~~~
is a legal `Go` program, `gijit` will reject it.

`gijit` won't let you re-use any of the basic,
pre-declared type names as variables. `uint`, `int`, `int8`, `int16`,
`int32`, `int64`, `uint8`, ... etc. are all off-limits.

Although in Go this is technically allowed, it can be highly confusing.
It is poor practice.

The technical reason for this restriction in `gijit` is that otherwise
the Go type checker can be corrupted by simple syntax errors
involving pre-declared identifiers. That's not an
issue for a full-recompile from the scratch each time, but for
a continuously online typechecker, it is a problem. To
stay online and functional after a syntax error like `var int w`,
(where w is unknown, provoking a syntax error yet shadowing
the pre-declared type), we disallow such variable names.


2018 Jan 31 update
-------
As of release v0.8.2, we support pointers, taking and de-referencing.

As of release v0.8.1, `gijit` works as a calculator. It will evaluate expressions at the command line.

In order to continue to detect syntax errors in Go code, we
adopt the same convention as Lua: the user must prepend an '=' equals sign to the
expression. For example:

~~~
gi> = 24/3
8LL

elapsed: '117.858µs'
gi>
~~~
Aside: the 'LL' suffix indicates a 64-bit, signed integer. Borrowed by LuaJIT from C/C++, it stands for "long long". There's also 'ULL' for uint64.

Notice, however, that without the '=', a syntax error is properly detected:
~~~
gi> 24/3
oops: 'problem detected during Go static type checking: '1:1: expected declaration, found 'INT' 24'' on input '24/3'

elapsed: '24.166µs'
gi> 
~~~

Keeping Go's type checking intact at the REPL preserves one of the most important
advantages of Go. We catch typos early, at compile time.

Multiple expressions at once also work, and each is printed on its own line.
~~~
gi> = 2+4, "gophers" + " " + "rock",  7-3
6LL
`gophers rock`
4LL

elapsed: '155.55µs'
gi> 
~~~

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
$ gi -q   ## start quietly, omit the banner for the demo.
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

# binary imports

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

As an example of shadwoing the io/ioutil package, we ran:
~~~
$ gen-gijit-shadow-import io/ioutil
writing to odir '/Users/jaten/go/src/github.com/gijit/gi/pkg/compiler/shadow/io/ioutil'
$ 
~~~

While shadowing (sandboxing) does
not allow arbitrary imports to be called from inside
the `gi` without prior preparation, this is
often a useful and desirable security feature. Moreover
we are able to provide these imports without using Go's
linux-only DLL loading system, so we remain portable/more
cross-platform compatible.

Interfaces are not yet implemented, and so are not yet imported.
Update: interfaces are represented in the imports, but
are not well tested. Please file issues as you find them.

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
replies. https://stackoverflow.com/questions/48202338/on-latest-luajit-2-1-0-beta3-is-recursive-xpcall-possible.

[Update: Mike Pall replied, (yay!) -- but (sadly),
it doesn't sound like he'll be fixing
this himself.

> Error handlers shouldn't throw errors themselves. The semantics would be too messy.
 -- Mike Pall, on https://github.com/LuaJIT/LuaJIT/issues/383
 ]

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


# special note on installing under the borked go1.9.4:

Very important: if you want to build with go1.9.4, well, just don't. go1.9.4 is a very borked release with respect to CGO projects; we recommend you avoid it. go1.9.5 will fix the most glaring problems; see https://github.com/golang/go/issues/23749 and just use go1.9.3, it works fine. Should you masochistically attempt to use go1.9.4, poor soul, then prior to building, you must add
~~~
export CGO_LDFLAGS_ALLOW='.*\.a$'
~~~
to your ~/.bashrc or equivalent, and restart your shell before building, so that `CGO_LDFLAGS_ALLOW` is defined in your environment prior to building.

Verify that `CGO_LDFLAGS_ALLOW` has been set before proceeding (only under the not-recommended go1.9.4):
~~~
$ env | grep CGO_LDFLAGS_ALLOW
CGO_LDFLAGS_ALLOW=.*\.a$
$
~~~
If you don't see `CGO_LDFLAGS_ALLOW` defined as the above, then fix your environment first.

Then, (or just start here for go1.9.3) to install:
~~~
$ go get -d github.com/gijit/gi/cmd/gi
$ cd $GOPATH/src/github.com/gijit/gi && make install
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
here https://github.com/gijit/gi/tree/master/emacs/gijit.el

M-x `run-gijit` to start the interpreter. Pressing ctrl-n will
step through any file that is in `gijit` mode. 

Other editors: please contribute!

# how to contribute: getting started

a) Pick an issue from here, https://github.com/gijit/gi/issues, and add a comment that
you are starting work on that feature. Make a branch for your feature, using `git checkout -b yourFeatureName`.

b) Write a test for your feature. Make sure it fails (the test is red), before
moving on to implementation. Tests are quite short. There are many examples here
in the pkg/compiler/*_test.go files. These show the currently implemented
features. Add your test to a new _test.go file.

https://github.com/gijit/gi/blob/master/pkg/compiler/repl_test.go

Then simply implement your feature. (So simple! Yeah right!)

So this is fun part. It's too situational to give general advice, but do see the hints
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

Specific javascript to Lua translation hints are below. Note
that `gijit` doesn't generate javascript. However, the
compiler package is derived from GopherJS, which did.

The job of making the full transition
from Javascript to Lua within the GopherJS-derived
code base is half-done/still in progress.

So when deciding how to change the output in the GopherJS derived code
of a particular javascript idiom, we
note the following are helpful hints.

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

e) join of strings:

table.concat({"a", "b", "c"}, ",")

f splitting strings, see string:split(sep)

NB: Lua has a limit of a few thousand return values.

f) s.substr(n): return the substring of s starting at n (0-indexed)

Javascript s.substr(4) is a zero-indexed substring from 4 to end of string `s`.

The Lua equivalent is string.sub(s, 5)

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
