.PHONY: tags

install:
	git branch --set-upstream-to=origin/master
	CGO_LDFLAGS_ALLOW='.*\.a$$' go get -t -d -u ./pkg/... ./cmd/...
	## LuaJIT compilation is now done manually, with windows.bat or posix.sh
	##cd vendor/github.com/LuaJIT/LuaJIT/src && make libluajit.a
	CGO_LDFLAGS_ALLOW='.*\.a$$' cd cmd/gi && make install

minimal:
minimal:
	cd cmd/gi && make install

tags:
	find . -name "*.[chCH]" -o -name "*.lua" -print | etags -


