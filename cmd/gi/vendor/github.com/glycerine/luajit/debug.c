#include <luajit-2.1/lua.h>
#include "_cgo_export.h"

static void
bouncehook(lua_State *s, lua_Debug *ar)
{
	hookevent(s, ar);
}

void
sethook(lua_State *s, int mask, int count)
{
	lua_sethook(s, bouncehook, mask, count);
}
