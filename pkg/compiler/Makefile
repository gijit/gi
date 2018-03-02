
all:
	gen_static_prelude && cd gi && make

## allows testing under go1.9.4 using "make test -- -run 063" for example.
test:
	export CGO_LDFLAGS_ALLOW='.*\.a$$'; go test -v $(filter-out $@,$(MAKECMDGOALS))
%:
	@:
