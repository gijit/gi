.PHONY: tags

install:
	git branch --set-upstream-to=origin/master
	CGO_LDFLAGS_ALLOW='.*\.a$$' go get -t -d -u ./pkg/... ./cmd/... && cd cmd/gi && make onetime && make install

minimal:
minimal:
	cd cmd/gi && make install

tags:
	find . -name "*.[chCH]" -o -name "*.lua" -print | etags -


