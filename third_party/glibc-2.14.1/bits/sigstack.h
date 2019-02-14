/* sigstack, sigaltstack definitions.
   Copyright (C) 1998, 1999 Free Software Foundation, Inc.
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

#ifndef _SIGNAL_H
# error "Never include this file directly.  Use <signal.h> instead"
#endif


/* Structure describing a signal stack (obsolete).  */
struct sigstack
  {
    __ptr_t ss_sp;		/* Signal stack pointer.  */
    int ss_onstack;		/* Nonzero if executing on this stack.  */
  };


/* Alternate, preferred interface.  */
typedef struct sigaltstack
  {
    __ptr_t ss_sp;
    size_t ss_size;
    int ss_flags;
  } stack_t;


/* Possible values for `ss_flags.'.  */
enum
{
  SS_ONSTACK = 0x0001,
#define SS_ONSTACK	SS_ONSTACK
  SS_DISABLE = 0x0004
#define SS_DISABLE	SS_DISABLE
};

/* Minumum stack size for a signal handler.  */
#define MINSIGSTKSZ	8192

/* System default stack size.  */
#define SIGSTKSZ	(MINSIGSTKSZ + 32768)
