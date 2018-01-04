## `make install` will do a full build and insteall,
##  which includes asking you for a sudo password
##  so that it can install luajit into /usr/local/bin.
install:
	go get github.com/glycerine/luajit
	cd cmd/gi && make onetime
