#include <math_ldbl_opt.h>
#undef weak_alias
#define weak_alias(n,a)
#include <sysdeps/ieee754/ldbl-128/s_sinl.c>
long_double_symbol (libm, __sinl, sinl);
