
all:
	gen_static_prelude && cd ../../cmd/gi && make

## allows testing under go1.9.4 using "make test -- -run 063" for example.
test:
	export CGO_LDFLAGS_ALLOW='.*\.a$$'; go test -v $(filter-out $@,$(MAKECMDGOALS))
%:
	@:

testbuild:
	go test -c -gcflags "-N -l" -v -race
