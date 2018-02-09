.PHONY: tags

install:
	export CGO_LDFLAGS_ALLOW="${GOPATH}/src/github.com/gijit/gi/vendor/github.com/glycerine/golua/lua/../../../LuaJIT/LuaJIT/src/libluajit.a"; cd cmd/gi && make onetime && make install

minimal:
	export CGO_LDFLAGS_ALLOW="${GOPATH}/src/github.com/gijit/gi/vendor/github.com/glycerine/golua/lua/../../../LuaJIT/LuaJIT/src/libluajit.a"; cd cmd/gi && make install

tags:
	find . -name "*.[chCH]" -o -name "*.lua" -print | etags -


