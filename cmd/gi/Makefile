

all:
	@echo "type 'make onetime' to install luajit (only necessary once)."
	@echo "You can skip 'make luajit' if you already have luajit-2.1.0-beta3 installed from source."
	make install


onetime: gitversion
	cd vendor/github.com/LuaJIT/LuaJIT && make install
	echo "installing this symlink may ask you for your password: "
	echo
	sudo ln -sf luajit-2.1.0-beta3 /usr/local/bin/luajit
	make all

# This how we want to name the binary output
BINARY=gi

BUILD_TIMESTAMP=$(shell date +%FT%T%z)
LAST_GIT_COMMIT_HASH=$(shell git rev-parse HEAD)
NEAREST_GIT_TAG=$(shell git describe --abbrev=0 --tags)
GIT_BRANCH=$(shell git rev-parse --abbrev-ref  HEAD)
GOVER=$(shell go version | sed 's/\ /_/g')
LUAJIT_VER=$(shell luajit -v | sed 's/\ /_/g')

LDFLAGS=-ldflags "-X main.LastGitCommitHash=${LAST_GIT_COMMIT_HASH} -X main.BuildTimeStamp=${BUILD_TIMESTAMP} -X main.GitBranch=${GIT_BRANCH} -X main.NearestGitTag=${NEAREST_GIT_TAG}  -X main.GoVersion=${GOVER}  -X main.LuajitVersion=${LUAJIT_VER}"


SRCS=$(wildcard *.go)
OBJS=$(SRCS:.go=.o)


build:
	go build ${LDFLAGS} -o gi


install:
	make build
	cp -p gi ${GOPATH}/bin/

# Cleans our project: deletes binaries
clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi


.PHONY: all onetime gitversion clean install
