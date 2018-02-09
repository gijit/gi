.PHONY: tags

install:
	export CGO_LDFLAGS_ALLOW=".*\.a"; cd cmd/gi && make onetime && make install

minimal:
	export CGO_LDFLAGS_ALLOW=".*\.a"; cd cmd/gi && make install

tags:
	find . -name "*.[chCH]" -o -name "*.lua" -print | etags -


