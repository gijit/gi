all:
	cd gi && make

## allows testing under go1.9.4 using "make test -- -run 063" for example.
test:
	export CGO_LDFLAGS_ALLOW="${GOPATH}/src/github.com/gijit/gi/vendor/github.com/glycerine/golua/lua/../../../LuaJIT/LuaJIT/src/libluajit.a"; go test -v $(filter-out $@,$(MAKECMDGOALS))
%:
	@:
