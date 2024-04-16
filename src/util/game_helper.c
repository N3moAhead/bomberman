#include <stdlib.h>
#ifdef _WIN32
#include <Windows.h>
#else
#include <unistd.h>
#endif
#include "game_helper.h"

void delay()
{
#ifdef _WIN32
  Sleep(300);
#else
  usleep(300000);
#endif
}
