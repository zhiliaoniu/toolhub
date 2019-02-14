/* Copyright (C) 1993, 1995, 1997, 1998, 1999, 2002, 2003, 2007
   Free Software Foundation, Inc.
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
   02111-1307 USA.

   As a special exception, if you link the code in this file with
   files compiled with a GNU compiler to produce an executable,
   that does not cause the resulting executable to be covered by
   the GNU Lesser General Public License.  This exception does not
   however invalidate any other reasons why the executable file
   might be covered by the GNU Lesser General Public License.
   This exception applies to code released by its copyright holders
   in files containing the exception.  */

#include "libioP.h"
#include <stdio.h>

size_t
__fread_unlocked_chk (void *__restrict ptr, size_t ptrlen,
		      size_t size, size_t n, FILE *__restrict stream)
{
  size_t bytes_requested = size * n;
  if (__builtin_expect ((n | size)
			>= (((size_t) 1) << (8 * sizeof (size_t) / 2)), 0))
    {
      if (size != 0 && bytes_requested / size != n)
	__chk_fail ();
    }

  if (__builtin_expect (bytes_requested > ptrlen, 0))
    __chk_fail ();

  CHECK_FILE (stream, 0);
  if (bytes_requested == 0)
    return 0;

  size_t bytes_read
    = INTUSE(_IO_sgetn) (stream, (char *) ptr, bytes_requested);
  return bytes_requested == bytes_read ? n : bytes_read / size;
}
