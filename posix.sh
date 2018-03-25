#!/bin/sh
echo "posix.sh: Doing one time build of libluajit.a"
cd vendor/github.com/LuaJIT/LuaJIT && make clean && cd src && XCFLAGS=-DLUAJIT_ENABLE_GC64 make libluajit.a
make gijit_luajit && cp -p gijit_luajit ${GOPATH}/bin/

