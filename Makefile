.PHONY: tags

install:
	cd cmd/gi && make onetime && make install

tags:
	find . -name "*.[chCH]" -o -name "*.lua" -print | etags -


