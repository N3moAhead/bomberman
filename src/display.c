#include <stdlib.h>
#include <stdio.h>
#include "display.h"
#include "constants.h"
#include "types.h"

void display(block_t **map)
{
  char display[MAP_HEIGHT * MAP_WIDTH * 20];
  int write_index = 0;
  for (int row = 0; row < MAP_HEIGHT; row++)
  {
    for (int col = 0; col < MAP_WIDTH; col++)
    {
      switch (map[row][col])
      {
      case AIR:
        write_index += snprintf(
            display + write_index,
            sizeof(display) - write_index,
            "  ");
        break;
      case BOMB1:
      case BOMB2:
      case BOMB3:
        write_index += snprintf(
            display + write_index,
            sizeof(display) - write_index,
            "💣");
        break;
      case WALL:
        write_index += snprintf(
            display + write_index,
            sizeof(display) - write_index,
            "🧱");
        break;
      case EXPLOSION:
        write_index += snprintf(
            display + write_index,
            sizeof(display) - write_index,
            "💥");
        break;
      case PLAYER1:
        write_index += snprintf(
            display + write_index,
            sizeof(display) - write_index,
            "🕺");
        break;
      case PLAYER2:
        write_index += snprintf(
            display + write_index,
            sizeof(display) - write_index,
            "🏃");
        break;
      case PLAYER3:
        write_index += snprintf(
            display + write_index,
            sizeof(display) - write_index,
            "🧍");
        break;
      case PLAYER4:
        write_index += snprintf(
            display + write_index,
            sizeof(display) - write_index,
            "💃");
        break;
      }
    }
    write_index += snprintf(display + write_index, sizeof(display) - write_index, "\n");
  }
  printf("%s", display);
};

void clear_display()
{
#ifdef _WIN32
  system("cls");
#else
  system("clear");
#endif
}
