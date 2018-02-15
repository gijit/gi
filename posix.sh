#!/bin/sh
echo "posix.sh: Doing one time build of libluajit.a"
cd vendor/github.com/LuaJIT/LuaJIT && make clean && cd src && make libluajit.a
