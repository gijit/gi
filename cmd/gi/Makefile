

all:
	echo "type 'make luajit' to install luajit (only necessary once)."
	echo "You can skip 'make luajit' if you already have luajit-2.1.0-beta3 installed from source."
	go install


onetime: gitversion
	cd vendor/github.com/LuaJIT/LuaJIT && make install
	echo "installing this symlink may ask you for your password: "
	echo
	sudo ln -sf luajit-2.1.0-beta3 /usr/local/bin/luajit
	make all

# This how we want to name the binary output
BINARY=gi

BUILD_TIMESTAMP=`date +%FT%T%z`
LAST_GIT_COMMIT_HASH=`git rev-parse HEAD`
NEAREST_GIT_TAG=`git describe --abbrev=0 --tags`
GIT_BRANCH=`git rev-parse --abbrev-ref  HEAD`
GO_VERSION=`go version`


# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS=-ldflags "-w -s -X main.LastGitCommitHash=${LAST_GIT_COMMIT_HASH} -X main.BuildTimeStamp=${BUILD_TIMESTAMP} -X main.GitBranch=${GIT_BRANCH} -X main.GoVersion=${GO_VERSION} -X main.NearestGitTag=${NEAREST_GIT_TAG}"

# Builds the project
build:
	go build ${LDFLAGS} -o ${BINARY}_f1

# Installs our project: copies binaries
install:
	go install ${LDFLAGS}

# Cleans our project: deletes binaries
clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi


.PHONY: all onetime gitversion clean install
