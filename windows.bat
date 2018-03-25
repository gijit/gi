@echo "windows.bat: Doing one time build of libluajit.a"
cd vendor\github.com\LuaJIT\LuaJIT\src
make libluajit.a
make gijit_luajit
cd ..\..\..\..\..
