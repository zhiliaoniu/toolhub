#include <complex.h>
#include <math_ldbl_opt.h>
#include <math/s_cacosh.c>
#if LONG_DOUBLE_COMPAT(libm, GLIBC_2_1)
compat_symbol (libm, __cacosh, cacoshl, GLIBC_2_1);
#endif
