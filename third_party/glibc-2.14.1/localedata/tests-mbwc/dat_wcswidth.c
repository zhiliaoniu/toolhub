/*
 *  TEST SUITE FOR MB/WC FUNCTIONS IN C LIBRARY
 *
 *	 FILE:	dat_wcswidth.c
 *
 *	 WCSWIDTH:  int	 wcswidth (const wchar_t *ws, size_t n);
 */

TST_WCSWIDTH tst_wcswidth_loc [] = {
  {
    { Twcswidth, TST_LOC_de },
    {
      { /*input.*/ { { 0x00C1,0x00C2,0x00C3,0x0000 },	 0 },  /* 01 */
	/*expect*/ { 0,1,0				   },
      },
      { /*input.*/ { { 0x00C1,0x00C2,0x00C3,0x0000 },	 1 },  /* 02 */
	/*expect*/ { 0,1,1				   },
      },
      { /*input.*/ { { 0x00C1,0x00C2,0x00C3,0x0000 },	 2 },  /* 03 */
	/*expect*/ { 0,1,2				   },
      },
      { /*input.*/ { { 0x00C1,0x00C2,0x00C3,0x0000 },	 3 },  /* 04 */
	/*expect*/ { 0,1,3				   },
      },
      { /*input.*/ { { 0x00C1,0x00C2,0x00C3,0x0000 },	 4 },  /* 05 */
	/*expect*/ { 0,1,3				   },
      },
      { /*input.*/ { { 0x0000 },			 1 },  /* 06 */
	/*expect*/ { 0,1,0				   },
      },
      { /*input.*/ { { 0x00C1,0x0001,0x0000 },		 2 },  /* 07 */
	/*expect*/ { 0,1,-1				   },
      },
      { /*input.*/ { { 0x00C1,0x0001,0x0000 },		 1 },  /* 08 */
	/*expect*/ { 0,1,1				   },
      },
      { /*input.*/ { { 0x00C1,0x0001,0x0000 },		 2 },  /* 09 */
	/*expect*/ { 0,1,-1				   },
      },
      { /*input.*/ { { 0x00C1,0x0092,0x0000 },		 2 },  /* 10 */
	/*expect*/ { 0,1,-1				   },
      },
      { /*input.*/ { { 0x00C1,0x0020,0x0000 },		 2 },  /* 11 */
	/*expect*/ { 0,1,2				   },
      },
      { /*input.*/ { { 0x00C1,0x0021,0x0000 },		 2 },  /* 12 */
	/*expect*/ { 0,1,2				   },
      },
      { /*input.*/ { { 0x00C1,0x007E,0x0000 },		 2 },  /* 13 */
	/*expect*/ { 0,1,2				   },
      },
      { /*input.*/ { { 0x00C1,0x007F,0x0000 },		 2 },  /* 14 */
	/*expect*/ { 0,1,-1				   },
      },
      { /*input.*/ { { 0x00C1,0x0080,0x0000 },		 2 },  /* 15 */
	/*expect*/ { 0,1,-1				   },
      },
      { /*input.*/ { { 0x00C1,0x00A0,0x0000 },		 2 },  /* 16 */
#ifdef SHOJI_IS_RIGHT
	/*expect*/ { 0,1,-1				   },
#else
	/*expect*/ { 0,1,2				   },
#endif
      },
      { /*input.*/ { { 0x00C1,0x00A1,0x0000 },		 2 },  /* 17 */
	/*expect*/ { 0,1,2				   },
      },
      { /*input.*/ { { 0x00C1,0x00FF,0x0000 },		 2 },  /* 18 */
	/*expect*/ { 0,1,2				   },
      },
      { /*input.*/ { { 0x00C1,0x3042,0x0000 },		 2 },  /* 19 */
	/*expect*/ { 0,1,-1				   },
      },
      { /*input.*/ { { 0x00C1,0x3044,0x0000 },		 2 },  /* 20 */
	/*expect*/ { 0,1,-1				   },
      },
      { .is_last = 1 }
    }
  },
  {
    { Twcswidth, TST_LOC_enUS },
    {
      { /*input.*/ { { 0x0041,0x0042,0x00C3,0x0000 },	 0 },  /* 01 */
	/*expect*/ { 0,1,0				   },
      },
      { /*input.*/ { { 0x0041,0x0042,0x00C3,0x0000 },	 1 },  /* 02 */
	/*expect*/ { 0,1,1				   },
      },
      { /*input.*/ { { 0x0041,0x0042,0x00C3,0x0000 },	 2 },  /* 03 */
	/*expect*/ { 0,1,2				   },
      },
#ifdef SHOJI_IS_RIGHT
      { /*input.*/ { { 0x0041,0x0042,0x00C3,0x0000 },	 3 },  /* 04 */
	/*expect*/ { 0,1,3				   },
      },
      { /*input.*/ { { 0x0041,0x0042,0x00C3,0x0000 },	 4 },  /* 05 */
	/*expect*/ { 0,1,3				   },
      },
#else
      { /*input.*/ { { 0x0041,0x0042,0x00C3,0x0000 },	 3 },  /* 04 */
	/*expect*/ { 0,1,-1				   },
      },
      { /*input.*/ { { 0x0041,0x0042,0x0043,0x0000 },	 4 },  /* 05 */
	/*expect*/ { 0,1,3				   },
      },
#endif
      { /*input.*/ { { 0x0000 },			 1 },  /* 06 */
	/*expect*/ { 0,1,0				   },
      },
      { /*input.*/ { { 0x0041,0x0001,0x0000 },		 2 },  /* 07 */
	/*expect*/ { 0,1,-1				   },
      },
      { /*input.*/ { { 0x0041,0x0001,0x0000 },		 1 },  /* 08 */
	/*expect*/ { 0,1,1				   },
      },
      { /*input.*/ { { 0x0041,0x0001,0x0000 },		 2 },  /* 09 */
	/*expect*/ { 0,1,-1				   },
      },
      { /*input.*/ { { 0x0041,0x0092,0x0000 },		 2 },  /* 10 */
	/*expect*/ { 0,1,-1				   },
      },
      { /*input.*/ { { 0x0041,0x0020,0x0000 },		 2 },  /* 11 */
	/*expect*/ { 0,1,2				   },
      },
      { /*input.*/ { { 0x0041,0x0021,0x0000 },		 2 },  /* 12 */
	/*expect*/ { 0,1,2				   },
      },
      { /*input.*/ { { 0x0041,0x007E,0x0000 },		 2 },  /* 13 */
	/*expect*/ { 0,1,2				   },
      },
      { /*input.*/ { { 0x0041,0x007F,0x0000 },		 2 },  /* 14 */
	/*expect*/ { 0,1,-1				   },
      },
      { /*input.*/ { { 0x0041,0x0080,0x0000 },		 2 },  /* 15 */
	/*expect*/ { 0,1,-1				   },
      },
      { /*input.*/ { { 0x0041,0x00A0,0x0000 },		 2 },  /* 16 */
	/*expect*/ { 0,1,-1				   },
      },
#ifdef SHOJI_IS_RIGHT
      { /*input.*/ { { 0x0041,0x00A1,0x0000 },		 2 },  /* 17 */
	/*expect*/ { 0,1,2				   },
      },
      { /*input.*/ { { 0x0041,0x00FF,0x0000 },		 2 },  /* 18 */
	/*expect*/ { 0,1,2				   },
      },
#else
      { /*input.*/ { { 0x0041,0x007E,0x0000 },		 2 },  /* 17 */
	/*expect*/ { 0,1,2				   },
      },
      { /*input.*/ { { 0x0041,0x0020,0x0000 },		 2 },  /* 18 */
	/*expect*/ { 0,1,2				   },
      },
#endif
      { /*input.*/ { { 0x0041,0x3042,0x0000 },		 2 },  /* 19 */
	/*expect*/ { 0,1,-1				   },
      },
      { /*input.*/ { { 0x0041,0x3044,0x0000 },		 2 },  /* 20 */
	/*expect*/ { 0,1,-1				   },
      },
      { .is_last = 1 }
    }
  },
  {
    { Twcswidth, TST_LOC_eucJP },
    {
      { /*input.*/ { { 0x3041,0x3042,0x3043,0x0000 },	 0 },  /* 01 */
	/*expect*/ { 0,1,0				   },
      },
      { /*input.*/ { { 0x3041,0x3042,0x3043,0x0000 },	 1 },  /* 02 */
	/*expect*/ { 0,1,2				   },
      },
      { /*input.*/ { { 0x3041,0x3042,0x3043,0x0000 },	 2 },  /* 03 */
	/*expect*/ { 0,1,4				   },
      },
      { /*input.*/ { { 0x3041,0x3042,0x3043,0x0000 },	 3 },  /* 04 */
	/*expect*/ { 0,1,6				   },
      },
      { /*input.*/ { { 0x3041,0x3042,0x3043,0x0000 },	 4 },  /* 05 */
	/*expect*/ { 0,1,6				   },
      },
      { /*input.*/ { { 0x0000 },			 1 },  /* 06 */
	/*expect*/ { 0,1,0				   },
      },
      { /*input.*/ { { 0x008E,0x0001,0x0000 },		 2 },  /* 07 */
	/*expect*/ { 0,1,-1				   },
      },
      { /*input.*/ { { 0x3041,0x008E,0x0000 },		 1 },  /* 08 */
	/*expect*/ { 0,1,2				   },
      },
      { /*input.*/ { { 0x3041,0x008E,0x0000 },		 2 },  /* 09 */
	/*expect*/ { 0,1,-1				   },
      },
      { /*input.*/ { { 0x3041,0x0001,0x0000 },		 2 },  /* 10 */
	/*expect*/ { 0,1,-1				   },
      },
      { /*input.*/ { { 0x3041,0x3000,0x0000 },		 2 },  /* 11 */
	/*expect*/ { 0,1,4				   },
      },
      { /*input.*/ { { 0x0041,0x0021,0x0000 },		 2 },  /* 12 */
	/*expect*/ { 0,1,2				   },
      },
      { /*input.*/ { { 0x0041,0x007E,0x0000 },		 2 },  /* 13 */
	/*expect*/ { 0,1,2				   },
      },
      { /*input.*/ { { 0x0041,0x007F,0x0000 },		 2 },  /* 14 */
	/*expect*/ { 0,1,-1				   },
      },
      { /*input.*/ { { 0x0041,0x0080,0x0000 },		 2 },  /* 15 */
	/*expect*/ { 0,1,-1				   },
      },
      { /*input.*/ { { 0x0041,0x00A0,0x0000 },		 2 },  /* 16 */
	/*expect*/ { 0,1,-1				   },
      },
#ifdef NO_WAIVER
      /* <NO_WAIVER> */	 /* returns 3 */
      { /*input.*/ { { 0x0041,0x00A1,0x0000 },		 2 },  /* 17 */
	/*expect*/ { 0,1,-1				   },
      },
#else
      /* XXX U00A1 is valid -> /x8f/xa2/xc4 in JIS X 0212 */
      { /*input.*/ { { 0x0041,0x00A1,0x0000 },		 2 },  /* 17 */
	/*expect*/ { 0,1,3				   },
      },
#endif
      { /*input.*/ { { 0x0041,0xFF71,0x0000 },		 2 },  /* 18 */
	/*expect*/ { 0,1,2				   },
      },
      { /*input.*/ { { 0x0041,0x3042,0x0000 },		 2 },  /* 19 */
	/*expect*/ { 0,1,3				   },
      },
      { /*input.*/ { { 0x0041,0x3044,0x0000 },		 2 },  /* 20 */
	/*expect*/ { 0,1,3				   },
      },
      { .is_last = 1 }
    }
  },
  {
    { Twcswidth, TST_LOC_end }
  }
};
