#include <complex.h>
#include <math_ldbl_opt.h>
#include <math/cimag.c>
#if LONG_DOUBLE_COMPAT(libm, GLIBC_2_1)
compat_symbol (libm, __cimag, cimagl, GLIBC_2_1);
#endif
