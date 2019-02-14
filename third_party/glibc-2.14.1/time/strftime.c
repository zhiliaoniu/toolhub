/* Copyright (C) 1991-1999, 2000, 2001, 2002, 2003, 2004
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
   02111-1307 USA.  */

#include <time.h>
#include <locale/localeinfo.h>


size_t
strftime (char *s, size_t maxsize, const char *format, const struct tm *tp)
{
  return __strftime_l (s, maxsize, format, tp, _NL_CURRENT_LOCALE);
}
libc_hidden_def (strftime)
