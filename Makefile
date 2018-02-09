.PHONY: tags

minimal:
	export CGO_CFLAGS_ALLOW="$(GOPATH)/github.com/gijit/gi/vendor/github.com/glycerine/golua/lua/../../../LuaJIT/LuaJIT/src/libluajit.a"; cd cmd/gi && make install

install:
	export CGO_CFLAGS_ALLOW="$(GOPATH)/github.com/gijit/gi/vendor/github.com/glycerine/golua/lua/../../../LuaJIT/LuaJIT/src/libluajit.a" cd cmd/gi && make onetime && make install

tags:
	find . -name "*.[chCH]" -o -name "*.lua" -print | etags -


