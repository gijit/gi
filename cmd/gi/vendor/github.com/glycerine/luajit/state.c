#include <luajit-2.1/lua.h>
#include <luajit-2.1/lauxlib.h>
#include <stddef.h>
#include <stdlib.h>
#include <string.h>
#include "_cgo_export.h"

enum {
	Bufsz=	256
};

typedef struct Readbuf	Readbuf;
struct Readbuf {
	void*	reader;
	char*	buf;
	size_t	bufsz;
};

/* a lua_Reader */
static const char*
readchunk(lua_State *l, void *data, size_t *size)
{
	Readbuf *rb;
	size_t sz;
	
	rb = data;
	memset(rb->buf, 0, rb->bufsz);
	sz = goreadchunk(rb->reader, rb->buf, rb->bufsz);
	if(sz < 1){
		free(rb->buf);
		free(rb);
		return NULL;
	}
	*size = sz;
	return rb->buf;
}

/* a lua_Writer */
static int
writechunk(lua_State *l, const void *p, size_t sz, void *ud)
{
	if(gowritechunk(ud, (void*)p, sz) != sz)
		return 1;
	return 0;
}

lua_State*
newstate(void)
{
	return luaL_newstate();
}

int
load(lua_State *l, void *reader, const char *chunkname)
{
	char *buf;
	Readbuf *rb;
	
	buf = malloc(Bufsz);		/* both allocs are freed by readchunk */
	if(buf == NULL)
		return LUA_ERRMEM;
	rb = malloc(sizeof *rb);
	if(rb == NULL){
		free(buf);
		return LUA_ERRMEM;
	}
	rb->reader = reader;
	rb->buf = buf;
	rb->bufsz = Bufsz;
	return lua_load(l, readchunk, rb, chunkname);
}

int
dump(lua_State *l, void *ud)
{
	return lua_dump(l, writechunk, ud);
}

/* a lua_CFunction */
static int
bounce(lua_State* s)
{
	void *fn;

	fn = lua_touserdata(s, lua_upvalueindex(1));
	return docallback(fn, s);
}

void
pushclosure(lua_State *s, int n)
{
	lua_pushcclosure(s, bounce, n + 1);
}
