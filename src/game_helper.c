#include <stdlib.h>
#include "game_helper.h"

void sleep()
{
#ifdef _WIN32
  system("timeout 1");
#else
  system("sleep 0.7");
#endif
}
