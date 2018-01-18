#include "lj_obj.h"

#if LJ_HASFFI

#include "lj_state.h"
#include "lj_gc.h"
#include "lj_err.h"
#include "lj_tab.h"
#include "lj_ctype.h"
#include "lj_cconv.h"
#include "lj_cdata.h"
#include "lauxlib.h"

LUA_API uint32_t
luajit_ctypeid(struct lua_State *L)
{
  int idx = lua_gettop(L);
  CTypeID ctypeid;
  GCcdata *cd;

  /* Get ref to ffi.typeof */
  int err = luaL_loadstring(L, "return require('ffi').typeof");
  if (err != 0) {
    lua_settop(L, idx);
    return luaL_error(L, "luajit-ffi-ctypeid error: could not loadstring");
  }

  err = lua_pcall(L, 0, 1, 0);
  if (err != 0) {
    lua_settop(L, idx);
    return luaL_error(L, "lua pcall to require ffi.typeof failed.");
  }
  
  if (!lua_isfunction(L, -1)) {
    int new_top = lua_gettop(L);
    lua_settop(L, idx);
    return luaL_error(L, "luajit-ffi-ctypeid: !lua_isfunction() at top of stack; new_top=%d", new_top);
  }
  /* Push the first argument to ffi.typeof */
  lua_pushvalue(L, idx);
  /* Call ffi.typeof() */
 lua_call(L, 1, 1);  
  /*
  err = lua_pcall(L, 1, 1, 0);
  if (err != 0) {
    lua_settop(L, idx);
  return luaL_error(L, "lua call to ffi.typeof with duplicated top of stack failed.");
  }
  */
  
  /* Returned type should be LUA_TCDATA with CTID_CTYPEID */
  if (lua_type(L, -1) != LUA_TCDATA) {
    lua_settop(L, idx);
    return luaL_error(L, "lua call to ffi.typeof failed at lua_type(L,1) != LUA_TCDATA");
  }
  /*cd = cdataV(L->base);*/
  cd = cdataV(L->top);
  if (cd->ctypeid != CTID_CTYPEID) {
    lua_settop(L, idx);
    return luaL_error(L, "lua call to ffi.typeof failed at ctypeid != CTID_CTYPEID");
  }

  ctypeid = *(CTypeID *)cdataptr(cd);
  lua_settop(L, idx);
  return ctypeid;
}

#endif
