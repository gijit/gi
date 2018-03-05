#include "lua.h"
#include "lauxlib.h"
#include "lualib.h"
#define LJ_HASFFI 1
#include "luajit-ffi-ctypeid.h"
#include <stdint.h>
#include <stdlib.h> // _atoi64 on windows, atoll on posix.
#include  <stdio.h>
#include "_cgo_export.h"

#define MT_GOFUNCTION "GoLua.GoFunction"
#define MT_GOINTERFACE "GoLua.GoInterface"

#define GOLUA_DEFAULT_MSGHANDLER "golua_default_msghandler"

// golua registry key, main states only, non non-main coroutines.
// The address of this constant is used as a unique
// lightuserdata key.
static const char GoMainStatesKey = 'k';

// Within one main state, store the main state (at 1),
// plus all other coroutines at uniqArray[pos > 1].
// Stored in the per state Lua registry.
// The address of this constant is used as a unique
// lightuserdata key.
static const char GoWithinStateUniqArrayKey = 'u'; //golua registry key, uniq array.

// Within one main state's 
// the reverseUniqMap maps *lua_State -> uPos
// Stored in the per state Lua registry.
// The address of this constant is used as a unique
// lightuserdata key.
static const char GoWithinStateRevUniqMap = 'r'; //in golua registry, reverse uniq map

static const char PanicFIDRegistryKey = 'k';

/* makes sure we compile in atoll/_atoi64 if available.*/
long long int wrapAtoll(const char *nptr)
{
#if _WIN32 || _WIN64
  return _atoi64(nptr);
#else
  return atoll(nptr);
#endif
}

static lua_State* getMainThread(lua_State* L) {
    // experiment: can we get main from uniqArray[1]; it should be there.
    int top = lua_gettop(L);
    lua_pushlightuserdata(L, (void*)&GoWithinStateUniqArrayKey);
	lua_gettable(L, LUA_REGISTRYINDEX); // pushes value onto top of stack
    // stack: uniqArray

    if (lua_isnil(L,-1)) {
      //printf("\n debug: arg. nil back for UniqKey lookup.\n");
    }
    
    lua_pushnumber(L, 1);
    // stack: 1, uniqArray
    lua_gettable(L, -2);
    // stack: main *lua_State, uniqArray

    int ty = lua_type(L, -1);
    //printf("\ndebug: top on top : %d.\n", ty);    
    
    lua_State* mainThread = lua_tothread(L, -1);
    //printf("\ndebug: UniqKey -> mainThread : %p.\n", mainThread);
    
    lua_settop(L, top);

    return mainThread;
}

/* taken from lua5.2 source */
void *testudata(lua_State *L, int ud, const char *tname)
{
	void *p = lua_touserdata(L, ud);
	if (p != NULL)
	{  /* value is a userdata? */
		if (lua_getmetatable(L, ud))
		{  /* does it have a metatable? */
			luaL_getmetatable(L, tname);  /* get correct metatable */
			if (!lua_rawequal(L, -1, -2))  /* not the same? */
				p = NULL;  /* value is a userdata with wrong metatable */
			lua_pop(L, 2);  /* remove both metatables */
			return p;
		}
	}
	return NULL;  /* value is not a userdata with a metatable */
}

int clua_isgofunction(lua_State *L, int n)
{
	return testudata(L, n, MT_GOFUNCTION) != NULL;
}

int clua_isgostruct(lua_State *L, int n)
{
	return testudata(L, n, MT_GOINTERFACE) != NULL;
}

unsigned int* clua_checkgosomething(lua_State* L, int index, const char *desired_metatable)
{
	if (desired_metatable != NULL)
	{
		return testudata(L, index, desired_metatable);
	}
	else
	{
		unsigned int *sid = testudata(L, index, MT_GOFUNCTION);
		if (sid != NULL) return sid;
		return testudata(L, index, MT_GOINTERFACE);
	}
}

size_t clua_getgostate(lua_State* L)
{
	size_t gostateindex;
	//get gostate from registry entry
	lua_pushlightuserdata(L,(void*)&GoMainStatesKey);
	lua_gettable(L, LUA_REGISTRYINDEX); // pushes value onto top of stack

    // 'k' is now a map from lua_State* to index
    // push key
    lua_pushthread(L);
    // stack is now:
    //  key
    //  map
    lua_gettable(L, -2);
    // stack is now
    //  index value
    //  map

    // if nil, return 0
    if (lua_isnil(L, -1)) {
      gostateindex = (size_t)(0);
    } else {
      gostateindex = (size_t)lua_touserdata(L, -1);
    }
	lua_pop(L, 2);
	return gostateindex;
}


//wrapper for callgofunction
int callback_function(lua_State* L)
{
	int r;
	unsigned int *fid = clua_checkgosomething(L, 1, MT_GOFUNCTION);
	size_t gostateindex = clua_getgostate(L);
    lua_State* mainThread = getMainThread(L);
	size_t mainIndex = clua_getgostate(mainThread);
    
    // jea: the metatable is on the stack.
    
	//remove the userdata metatable (go function??) from the stack (to present same behavior as lua_CFunctions)
	lua_remove(L, 1);
    
	return golua_callgofunction(L, gostateindex, mainIndex, mainThread, fid!=NULL ? *fid : -1);
}

//wrapper for gchook
int gchook_wrapper(lua_State* L)
{
	//printf("Garbage collection wrapper\n");
	unsigned int* fid = clua_checkgosomething(L, -1, NULL);
	size_t gostateindex = clua_getgostate(L);
	if (fid != NULL)
		return golua_gchook(gostateindex,*fid);
	return 0;
}

unsigned int clua_togofunction(lua_State* L, int index)
{
	unsigned int *r = clua_checkgosomething(L, index, MT_GOFUNCTION);
	return (r != NULL) ? *r : -1;
}

unsigned int clua_togostruct(lua_State *L, int index)
{
	unsigned int *r = clua_checkgosomething(L, index, MT_GOINTERFACE);
	return (r != NULL) ? *r : -1;
}

void clua_pushgofunction(lua_State* L, unsigned int fid)
{
	unsigned int* fidptr = (unsigned int *)lua_newuserdata(L, sizeof(unsigned int));
	*fidptr = fid;
	luaL_getmetatable(L, MT_GOFUNCTION);
	lua_setmetatable(L, -2);
}

static int callback_c (lua_State* L)
{
	int fid = clua_togofunction(L,lua_upvalueindex(1));
	size_t gostateindex = clua_getgostate(L);
	return golua_callgofunction(L, gostateindex, 0, getMainThread(L), fid);
}

void clua_pushcallback(lua_State* L)
{
	lua_pushcclosure(L,callback_c,1);
}

void clua_pushgostruct(lua_State* L, unsigned int iid)
{
	unsigned int* iidptr = (unsigned int *)lua_newuserdata(L, sizeof(unsigned int));
	*iidptr = iid;
	luaL_getmetatable(L, MT_GOINTERFACE);
	lua_setmetatable(L,-2);
}

int default_panicf(lua_State *L)
{
	const char *s = lua_tostring(L, -1);
	printf("Lua unprotected panic: %s\n", s);
	abort();
}


void create_uniqArray(lua_State* L) {
  // stack: ...
  
  lua_newtable(L);
  // stack: newtable, ...
  
  lua_pushlightuserdata(L,(void*)&GoWithinStateUniqArrayKey);
  // stack: ukey, newtable, ...
  
  lua_insert(L, -2);    
  // stack: newtable, ukey, ...
  
  lua_settable(L, LUA_REGISTRYINDEX);
  // stack: ...

  // and create the reverse uniq map: *lua_State -> uPos
  lua_newtable(L);
  // stack: newtable, ...
  
  lua_pushlightuserdata(L,(void*)&GoWithinStateRevUniqMap);
  // stack: urevkey, newtable, ...
  
  lua_insert(L, -2);    
  // stack: newtable, urevkey, ...
  
  lua_settable(L, LUA_REGISTRYINDEX);
  // stack: ...
}

// return 1 if created, 0 if already existed.
int clua_create_uniqArrayIfNotExists(lua_State* L) {
  // stack: ...
  
  lua_pushlightuserdata(L,(void*)&GoWithinStateUniqArrayKey);
  // stack: ukey, ...
  lua_gettable(L, LUA_REGISTRYINDEX);
  // stack: (uniqArray or nil), ...

  if (lua_isnil(L, -1)) {
      lua_pop(L, 1);
      // stack: ...
      create_uniqArray(L);
      // stack: ...
      return 1;
  }
  lua_pop(L, 1);
  // stack: ...
  return 0;
}

int clua_addThreadToUniqArrayAndRevUniq(lua_State* L);

// clua_known_coro returns the index
// of L in uniqArray, adding L to
// uniqArray and revUniqMap if it
// is not already present in revUniqMap.
//
// POST INVAR:
// On return of value uPos, uniqArray[uPos] == L
// and revUniqMap[L] == uPos.
// 
// The returned value will be >= 1.
// 
int clua_dedup_coro(lua_State* L)
{
  if (1 == clua_create_uniqArrayIfNotExists(L)) {
    return clua_addThreadToUniqArrayAndRevUniq(L);
  }
  
  int top = lua_gettop(L);

  // use revUniqMap, for O(1) lookup.

  // store pos into revUniqMap too, for O(1) coroutine lookup.
  lua_pushlightuserdata(L, (void*)&GoWithinStateRevUniqMap);
  // stack: revkey

  lua_gettable(L, LUA_REGISTRYINDEX);
  // stack: revUniqMap

  int isMain = lua_pushthread(L);
  // stack: thread, revUniqMap
    
  lua_gettable(L, -2);
  // stack: (pos or nil), revUniqMap

  if (lua_isnil(L, -1)) {
    // stack: nil, revUniqMap
    
    // not previously known
    lua_settop(L, top);
    // stack clean

    // add it
    return clua_addThreadToUniqArrayAndRevUniq(L);    
  }
  
  // stack: pos, revUniqMap
  int res = (int)lua_tonumber(L, -1);
  
  lua_settop(L, top);
  // stack clean
  return res;
}

int clua_addThreadToUniqArrayAndRevUniq(lua_State* L) {
  // 
  // Do the equivalent of: table.insert(uniqueArray, L)
  //                   and then revUniq[L] = #uniqArray (==uPos)
    
  // stack: ...
  int top = lua_gettop(L);
  
  lua_pushlightuserdata(L,(void*)&GoWithinStateUniqArrayKey);
  // stack: ukey, ...

  lua_gettable(L, LUA_REGISTRYINDEX);    
  // stack: uniqArray, ...

  // jea: now store the actual thread in uniqArray
  //
  lua_pushthread(L);
  // stack: thread, uniqArray, ...

  // append to array at -2
  int pos = lua_objlen(L, -2) + 1; /* first empty element */

  lua_rawseti(L, -2, pos);  /* t[pos] = v; and pops v from the top of the stack*/
  // stack: uniqArray, ...

  lua_pop(L, 1);
  // stack: ...

  // Phase 2:
  // 
  // store pos into revUniqMap too, for O(1) coroutine lookup.
  lua_pushlightuserdata(L, (void*)&GoWithinStateRevUniqMap);
  // stack: revkey, ...

  lua_gettable(L, LUA_REGISTRYINDEX);
  // stack: revUniqMap, ...

  lua_pushthread(L);
  // stack: thread, revUniqMap, ...
    
  lua_pushnumber(L, (double)pos);
  // stack: pos, thread, revUniqMap, ...

  lua_settable(L, -3);
  // stack: revUniqMap, ...

  lua_settop(L, top);
  // stack clean

  return pos;
}

int clua_setgostate(lua_State* L, size_t gostateindex)
{
  int ret = 0;
  int top = lua_gettop(L);
      
  lua_atpanic(L, default_panicf);
  lua_pushlightuserdata(L,(void*)&GoMainStatesKey);

  //
  // store L into a table that maps lua_State* -> gostateindex.
  // Call the table the Lmap. It maps L to an index.
  //
  lua_gettable(L, LUA_REGISTRYINDEX); // pops the key
  // does it already exist, or did we get nil back?
  if (lua_isnil(L, -1)) {
    // doesn't exist yet, need to create it the first time.
    lua_pop(L, 1); // get rid of the nil
    lua_newtable(L);

    // save Lmap into lua registry under GoMainStatesKey.

    // stack:
    //   Lmap
    lua_pushvalue(L, -1);
    // stack:
    //   Lmap
    //   Lmap
    
    // get the key
    lua_pushlightuserdata(L,(void*)&GoMainStatesKey);
    // stack is now:
    //  key
    //  Lmap
    //  Lmap

    lua_insert(L, -2);
    // stack should now be ready for lua_settable:
    //  Lmap
    //  key
    //  Lmap

    //set into registry table
    lua_settable(L, LUA_REGISTRYINDEX);    
    // stack:
    //   Lmap

    // and create the uniq array too, at the same time.
    clua_create_uniqArrayIfNotExists(L);
    
    // stack: Lmap

    //printf("\ndebug: UniqKey setup in lua registry.\n");
  }
  // INVAR: our Lmap is at top of stack, -1 position.
  // stack: Lmap
  
  // does our key already exist in the Lmap?
  // If not, then append to the uniqArray.
  // key
  lua_pushthread(L);
  // stack: key, Lmap

  lua_gettable(L, -2);
  // stack: nil, Lmap  OR  priorValue, Lmap

  if (!lua_isnil(L, -1)) {
    // stack: priorValue, Lmap
    // no duplicate insert into Lmap needed
    lua_settop(L, top);
    // stack clean.

    // should be returning 1 always, indicating
    // a main coroutine.
    return clua_dedup_coro(L);
  }
  
  // stack: thread, Lmap
  lua_pop(L, 1);
  // stack: Lmap
  
  // This thread, L, is not known to Lmap.
  ret = clua_addThreadToUniqArrayAndRevUniq(L);

  // stack: Lmap
    
  // Finally, populate the Lmap
  
  // key
  lua_pushthread(L);
  // stack: key, Lmap
  
  // value
  lua_pushlightuserdata(L, (void*)gostateindex);
  // stack: value, key, Lmap
  
  // store key:value in map
  lua_settable(L, -3);
  // stack: Lmap
  
  // cleanup stack, remove Lmap
  lua_settop(L, top);

  return ret;
}


/* called when lua code attempts to access a field of a published go object */
int interface_index_callback(lua_State *L)
{
	unsigned int *iid = clua_checkgosomething(L, 1, MT_GOINTERFACE);
	if (iid == NULL)
	{
		lua_pushnil(L);
		return 1;
	}

	char *field_name = (char *)lua_tostring(L, 2);
	if (field_name == NULL)
	{
		lua_pushnil(L);
		return 1;
	}

	size_t gostateindex = clua_getgostate(L);

	int r = golua_interface_index_callback(gostateindex, *iid, field_name);

	if (r < 0)
	{
		lua_error(L);
		return 0;
	}
	else
	{
		return r;
	}
}

/* called when lua code attempts to set a field of a published go object */
int interface_newindex_callback(lua_State *L)
{
	unsigned int *iid = clua_checkgosomething(L, 1, MT_GOINTERFACE);
	if (iid == NULL)
	{
		lua_pushnil(L);
		return 1;
	}

	char *field_name = (char *)lua_tostring(L, 2);
	if (field_name == NULL)
	{
		lua_pushnil(L);
		return 1;
	}

	size_t gostateindex = clua_getgostate(L);

	int r = golua_interface_newindex_callback(gostateindex, *iid, field_name);

	if (r < 0)
	{
		lua_error(L);
		return 0;
	}
	else
	{
		return r;
	}
}

int panic_msghandler(lua_State *L)
{
	size_t gostateindex = clua_getgostate(L);
	go_panic_msghandler(gostateindex, (char *)lua_tolstring(L, -1, NULL));
	return 0;
}

void clua_hide_pcall(lua_State *L)
{
	lua_getglobal(L, "pcall");
	lua_setglobal(L, "unsafe_pcall");
	lua_pushnil(L);
	lua_setglobal(L, "pcall");

	lua_getglobal(L, "xpcall");
	lua_setglobal(L, "unsafe_xpcall");
	lua_pushnil(L);
	lua_setglobal(L, "xpcall");
}

void clua_initstate(lua_State* L)
{
	/* create the GoLua.GoFunction metatable */
	luaL_newmetatable(L, MT_GOFUNCTION);

	// gofunction_metatable[__call] = &callback_function
	lua_pushliteral(L,"__call");
	lua_pushcfunction(L,&callback_function);
	lua_settable(L,-3);

	// gofunction_metatable[__gc] = &gchook_wrapper
	lua_pushliteral(L,"__gc");
	lua_pushcfunction(L,&gchook_wrapper);
	lua_settable(L,-3);
	lua_pop(L,1);

	luaL_newmetatable(L, MT_GOINTERFACE);

	// gointerface_metatable[__gc] = &gchook_wrapper
	lua_pushliteral(L, "__gc");
	lua_pushcfunction(L, &gchook_wrapper);
	lua_settable(L, -3);

	// gointerface_metatable[__index] = &interface_index_callback
	lua_pushliteral(L, "__index");
	lua_pushcfunction(L, &interface_index_callback);
	lua_settable(L, -3);

	// gointerface_metatable[__newindex] = &interface_newindex_callback
	lua_pushliteral(L, "__newindex");
	lua_pushcfunction(L, &interface_newindex_callback);
	lua_settable(L, -3);

	lua_register(L, GOLUA_DEFAULT_MSGHANDLER, &panic_msghandler);
	lua_pop(L, 1);
}


int callback_panicf(lua_State* L)
{
	lua_pushlightuserdata(L,(void*)&PanicFIDRegistryKey);
	lua_gettable(L,LUA_REGISTRYINDEX);
	unsigned int fid = lua_tointeger(L,-1);
	lua_pop(L,1);
	size_t gostateindex = clua_getgostate(L);
	return golua_callpanicfunction(gostateindex,fid);

}

//TODO: currently setting garbage when panicf set to null
GoInterface clua_atpanic(lua_State* L, unsigned int panicf_id)
{
	//get old panicfid
	unsigned int old_id;
	lua_pushlightuserdata(L, (void*)&PanicFIDRegistryKey);
	lua_gettable(L,LUA_REGISTRYINDEX);
	if(lua_isnil(L, -1) == 0)
		old_id = lua_tointeger(L,-1);
	lua_pop(L, 1);

	//set registry key for function id of go panic function
	lua_pushlightuserdata(L, (void*)&PanicFIDRegistryKey);
	//push id value
	lua_pushinteger(L, panicf_id);
	//set into registry table
	lua_settable(L, LUA_REGISTRYINDEX);

	//now set the panic function
	lua_CFunction pf = lua_atpanic(L,&callback_panicf);
	//make a GoInterface with a wrapped C panicf or the original go panicf
	if(pf == &callback_panicf)
	{
		return golua_idtointerface(old_id);
	}
	else
	{
		//TODO: technically UB, function ptr -> non function ptr
		return golua_cfunctiontointerface((GoUintptr *)pf);
	}
}

int clua_callluacfunc(lua_State* L, lua_CFunction f)
{
	return f(L);
}

void* allocwrapper(void* ud, void *ptr, size_t osize, size_t nsize)
{
	return (void*)golua_callallocf((GoUintptr)ud,(GoUintptr)ptr,osize,nsize);
}

lua_State* clua_newstate(void* goallocf)
{
	return lua_newstate(&allocwrapper,goallocf);
}

void clua_setallocf(lua_State* L, void* goallocf)
{
	lua_setallocf(L,&allocwrapper,goallocf);
}

void clua_openbase(lua_State* L)
{
	lua_pushcfunction(L,&luaopen_base);
	lua_pushstring(L,"");
	lua_call(L, 1, 0);
	clua_hide_pcall(L);
}

void clua_openio(lua_State* L)
{
	lua_pushcfunction(L,&luaopen_io);
	lua_pushstring(L,"io");
	lua_call(L, 1, 0);
}

void clua_openmath(lua_State* L)
{
	lua_pushcfunction(L,&luaopen_math);
	lua_pushstring(L,"math");
	lua_call(L, 1, 0);
}

void clua_openpackage(lua_State* L)
{
	lua_pushcfunction(L,&luaopen_package);
	lua_pushstring(L,"package");
	lua_call(L, 1, 0);
}

void clua_openstring(lua_State* L)
{
	lua_pushcfunction(L,&luaopen_string);
	lua_pushstring(L,"string");
	lua_call(L, 1, 0);
}

void clua_opentable(lua_State* L)
{
	lua_pushcfunction(L,&luaopen_table);
	lua_pushstring(L,"table");
	lua_call(L, 1, 0);
}

void clua_openos(lua_State* L)
{
	lua_pushcfunction(L,&luaopen_os);
	lua_pushstring(L,"os");
	lua_call(L, 1, 0);
}

void clua_hook_function(lua_State *L, lua_Debug *ar)
{
	lua_checkstack(L, 2);
	lua_pushstring(L, "Lua execution quantum exceeded");
	lua_error(L);
}

void clua_setexecutionlimit(lua_State* L, int n)
{
	lua_sethook(L, &clua_hook_function, LUA_MASKCOUNT, n);
}

/*return the ctype of the cdata at the top of the stack*/
uint32_t clua_luajit_ctypeid(lua_State *L, int idx)
{
  return luajit_ctypeid(L, idx);
}

void clua_luajit_push_cdata_int64(lua_State *L, int64_t n)
{
  return luajit_push_cdata_int64(L, n);
}

void clua_luajit_push_cdata_uint64(lua_State *L, uint64_t u)
{
  return luajit_push_cdata_uint64(L, u);
}

