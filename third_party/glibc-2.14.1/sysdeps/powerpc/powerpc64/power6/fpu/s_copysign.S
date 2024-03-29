/* copysign().  PowerPC64/POWER6 version.
   Copyright (C) 2010 Free Software Foundation, Inc.
   Contributed by Luis Machado <luisgpm@br.ibm.com>.
   This file is part of the GNU C Library.

   The GNU C Library is free software; you can redistribute it and/or
   modify it under the terms of the GNU Lesser General Public
   License as published by the Free Software Foundation; either
   version 2.1 of the License, or (at your option) any later version.

   The GNU C Library is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
   Lesser General Public License for more details.

   You should have received a copy of the GNU Lesser General Public
   License along with the GNU C Library; if not, write to the Free
   Software Foundation, Inc., 59 Temple Place, Suite 330, Boston, MA
   02111-1307 USA.  */

#include <sysdep.h>
#include <math_ldbl_opt.h>

/* double [f1] copysign (double [f1] x, double [f2] y);
   copysign(x,y) returns a value with the magnitude of x and
   with the sign bit of y.  */

	.section    ".text"
	.type	    __copysign, @function
	.machine    power6
EALIGN (__copysign, 4, 0)
	CALL_MCOUNT 0
	fcpsgn	fp1,fp2,fp1
	blr
END (__copysign)

hidden_def (__copysign)
weak_alias (__copysign, copysign)

/* It turns out that the 'double' version will also always work for
   single-precision.  */
strong_alias (__copysign, __copysignf)
hidden_def (__copysignf)
weak_alias (__copysignf, copysignf)

#ifdef NO_LONG_DOUBLE
strong_alias (__copysign, __copysignl)
weak_alias (__copysign, copysignl)
#endif

#ifndef IS_IN_libm
# if LONG_DOUBLE_COMPAT(libc, GLIBC_2_0)
compat_symbol (libc, __copysign, __copysignl, GLIBC_2_0);
compat_symbol (libc, copysign, copysignl, GLIBC_2_0);
# endif
#endif
