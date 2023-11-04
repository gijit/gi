#!/bin/sh
echo "posix.sh: Doing one time build of libluajit.a"

## recent macOS need the extra MACOSX_DEPLOYMENT_TARGET flag.
## See also https://github.com/LuaJIT/LuaJIT/issues/449
##
cd vendor/github.com/LuaJIT/LuaJIT && make clean && cd src && MACOSX_DEPLOYMENT_TARGET=10.14 XCFLAGS=-DLUAJIT_ENABLE_GC64 make libluajit.a
MACOSX_DEPLOYMENT_TARGET=10.14 make gijit_luajit && cp -p gijit_luajit ${GOPATH}/bin/

