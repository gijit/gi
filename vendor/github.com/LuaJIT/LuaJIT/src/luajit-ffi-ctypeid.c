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
  luaL_loadstring(L, "return require('ffi').typeof");

  int err = lua_pcall(L, 0, 1, 0);
  if (err != 0) {
    goto fail;
  }
  
  if (lua_gettop(L) > 1 || !lua_isfunction(L, 1))
    goto fail;
  /* Push the first argument to ffi.typeof */
  lua_pushvalue(L, idx);
  /* Call ffi.typeof() */
  err = lua_pcall(L, 1, 1, 0);
  if (err != 0) {
    goto fail;
  }
  
  /* Returned type should be LUA_TCDATA with CTID_CTYPEID */
  if (lua_type(L, 1) != LUA_TCDATA)
    goto fail;
  cd = cdataV(L->base);
  if (cd->ctypeid != CTID_CTYPEID)
    goto fail;

  ctypeid = *(CTypeID *)cdataptr(cd);
  lua_settop(L, idx);
  return ctypeid;
 fail:
  lua_settop(L, idx);
  return luaL_error(L, "lua call to ffi.typeof failed");
}

#endif
