#ifndef _LUAJIT_FFI_CTYPEID_H
#define _LUAJIT_FFI_CTYPEID_H

#if LJ_HASFFI
#include "lj_obj.h"

/* return the ctype of the cdata at the top of the stack*/
LUA_API uint32_t luajit_ctypeid(struct lua_State *L, int idx);

LUA_API void luajit_push_cdata_int64(struct lua_State *L, int64_t n);
LUA_API void luajit_push_cdata_uint64(struct lua_State *L, int64_t n);

#endif
#endif
