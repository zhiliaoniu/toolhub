#include "nldbl-compat.h"
#include <complex.h>

double _Complex
attribute_hidden
catanhl (double _Complex x)
{
  return catanh (x);
}
