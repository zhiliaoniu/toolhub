/*
 *  TEST SUITE FOR MB/WC FUNCTIONS IN C LIBRARY
 *
 *	 FILE:	dat_mbstowcs.c
 *
 *	 MBSTOWCS:  size_t  mbstowcs (wchar_t *ws, char *s, size_t n);
 */

#include <limits.h>

TST_MBSTOWCS tst_mbstowcs_loc [] = {
  {
    { Tmbstowcs, TST_LOC_de },
    {
      { /*----------------- #01 -----------------*/
	{
	  {
	    { 1,  1, "ABC",		   4			 },
	    { 1,  1, "ABC",		   3			 },
	    { 1,  1, "ABC",		   2			 },
	  }
	},
	{
	  {
	    { 0,1,3, { 0x0041,0x0042,0x0043,0x0000 }	 },
	    { 0,1,3, { 0x0041,0x0042,0x0043,0x0000 }	 },
	    { 0,1,2, { 0x0041,0x0042,0x0043,0x0000 }	 },
	  }
	}
      },
      { /*----------------- #02 -----------------*/
	{
	  {
	    { 1,  1, "ABC",		   4			 },
	    { 1,  1, "",		   1			 },
	    { 0,  1, "ABC",		   4			 },
	  }
	},
	{
	  {
	    { 0,1,3, { 0x0041,0x0042,0x0043,0x0000 }	 },
	    { 0,1,0, { 0x0000 }				 },
	    { 0,1,3, { 0x0000 }				 },
	  }
	}
      },
      { .is_last = 1 }
    }
  },
  {
    { Tmbstowcs, TST_LOC_enUS },
    {
      { /*----------------- #01 -----------------*/
	{
	  {
	    { 1,  1, "ABC",		   4			 },
	    { 1,  1, "ABC",		   3			 },
	    { 1,  1, "ABC",		   2			 },
	  }
	},
	{
	  {
	    { 0,1,3, { 0x0041,0x0042,0x0043,0x0000 }	 },
	    { 0,1,3, { 0x0041,0x0042,0x0043,0x0000 }	 },
	    { 0,1,2, { 0x0041,0x0042,0x0043,0x0000 }	 },
	  }
	}
      },
      { /*----------------- #02 -----------------*/
	{
	  {
	    { 1,  1, "ABC",		   4			 },
	    { 1,  1, "",		   1			 },
	    { 0,  1, "ABC",		   4			 },
	  }
	},
	{
	  {
	    { 0,1,3, { 0x0041,0x0042,0x0043,0x0000 }	 },
	    { 0,1,0, { 0x0000 }				 },
	    { 0,1,3, { 0x0000 }				 },
	  }
	}
      },
      { .is_last = 1 }
    }
  },
  {
    { Tmbstowcs, TST_LOC_eucJP },
    {
      { /*----------------- #01 -----------------*/
	{
	  {
	    { 1,  1, "\244\242\244\244\244\246ABC",      7 },
	    { 1,  1, "\244\242\244\244\244\246ABC",      6 },
	    { 1,  1, "\244\242\244\244\244\246ABC",      4 },
	  }
	},
	{
	  {
	    { 0,1,6, { 0x3042,0x3044,0x3046,0x0041,0x0042,0x0043,0x0000 }},
	    { 0,1,6, { 0x3042,0x3044,0x3046,0x0041,0x0042,0x0043,0x0000 }},
	    { 0,1,4, { 0x3042,0x3044,0x3046,0x0041,0x0000 }		 },
	  }
	}
      },
      { /*----------------- #02 -----------------*/
	{
	  {
#ifdef SHOJI_IS_RIGHT
	    /* XXX I really don't understand the first and third line.
	       the result of the first line is the same as the first
	       in the last test (i.e., returns 6).  Also, the third
	       test will simply convert everything.  */
	    { 1,  1, "\244\242\244\244\244\246ABC",      7 },
	    { 1,  1, "",                                 1 },
	    { 0,  1, "\244\242\244\244\244\246ABC",      7 },
#else
	    { 1,  1, "\244\242\244\244\244\246ABC",      4 },
	    { 1,  1, "",                                 1 },
	    { 0,  1, "\244\242\244\244\244\246ABC",      0 },
#endif
	  }
	},
	{
	  {
	    { 0,1,4, { 0x3042,0x3044,0x3046,0x0041,0x0000 }		 },
	    { 0,1,0, { 0x0000 }					 },
	    { 0,1,6, { 0x0000 }					 },
	  }
	}
      },
      { .is_last = 1 }
    }
  },
  {
    { Tmbstowcs, TST_LOC_end }
  }
};
