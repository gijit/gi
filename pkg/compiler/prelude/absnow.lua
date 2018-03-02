-- native time queries

local ffi = require("ffi")

-- establish a fine-grained __abs_now()
-- that can be used for nanosecond timing.

if jit.os == "Windows" then
   ffi.cdef[[

    typedef uint8_t BYTE;
    typedef uint32_t DWORD;
    typedef int32_t LONG;
    typedef int64_t LONGLONG;
    
    typedef union _LARGE_INTEGER {
      struct {
        DWORD LowPart;
        LONG  HighPart;
      };
      struct {
        DWORD LowPart;
        LONG  HighPart;
      } u;
      LONGLONG QuadPart;
    } LARGE_INTEGER, *PLARGE_INTEGER;
    
    int __stdcall QueryPerformanceFrequency(
        LARGE_INTEGER *lpFrequency
    );
    int __stdcall QueryPerformanceCounter(
      LARGE_INTEGER *lpPerformanceCount
   );

   ]]

   local tmp = ffi.new("LARGE_INTEGER");
   ffi.C.QueryPerformanceFrequency(tmp);
   local countPerSec = tmp.QuadPart;
   local nanoSecPerCount = 1000000000ULL/countPerSec
   
   __abs_now=function()
      local now = ffi.new("LARGE_INTEGER")
      ffi.C.QueryPerformanceCounter(now)
      return int64(now.QuadPart * nanoSecPerCount)
   end
   
elseif jit.os == "OSX" then

   ffi.cdef[[
   uint64_t mach_absolute_time(void);
   struct mach_timebase_info {
	uint32_t	numer;
	uint32_t	denom;
   };
   typedef struct mach_timebase_info *mach_timebase_info_t;
   typedef struct mach_timebase_info mach_timebase_info_data_t;
   void mach_timebase_info(mach_timebase_info_t info);
   ]]
   local info = ffi.new("mach_timebase_info_data_t")
   ffi.C.mach_timebase_info(info);
   --print("info.numer =", info.numer) -- 1
   --print("info.denom =", info.denom) -- 1

   -- returns a nanosecond time stamp, but not
   -- since epoch of 1970. Maybe since last
   -- reboot? subtract two to get useful nanoseconds.
   __abs_now=function()
      return int64(ffi.C.mach_absolute_time())
   end
   
else
   -- for linux, clock_gettime(CLOCK_MONOTONIC)

   ffi.cdef[[
    typedef long time_t;
    typedef int clockid_t;

    typedef struct timespec {
            time_t   tv_sec;        /* seconds */
            long     tv_nsec;       /* nanoseconds */
    } nanotime;
    int clock_gettime(clockid_t clk_id, struct timespec *tp);
   ]]

   __abs_now=function()
      local pnano = assert(ffi.new("nanotime[?]", 1))

      -- CLOCK_MONOTONIC = 1
      ffi.C.clock_gettime(1, pnano)
      return int64(pnano[0].tv_sec * 1000000000 + pnano[0].tv_nsec)
   end

end
