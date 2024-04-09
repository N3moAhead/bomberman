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
      case BOMB4:
      case BOMB5:
      case BOMB6:
      case BOMB7:
      case BOMB8:
      case BOMB9:
      case BOMB10:
        #ifdef _WIN32
          write_index += snprintf(
            display + write_index,
            sizeof(display) - write_index,
            " O");
        #else
          write_index += snprintf(
            display + write_index,
            sizeof(display) - write_index,
            "💣");
        #endif
        break;
      case WALL:
        #ifdef _WIN32
          write_index += snprintf(
              display + write_index,
              sizeof(display) - write_index,
              "##");
        #else
          write_index += snprintf(
              display + write_index,
              sizeof(display) - write_index,
              "🧱");
        #endif
        break;
      case EXPLOSION:
        #ifdef _WIN32
          write_index += snprintf(
              display + write_index,
              sizeof(display) - write_index,
              "XX");
        #else
          write_index += snprintf(
              display + write_index,
              sizeof(display) - write_index,
              "💥");
        #endif
        break;
      case PLAYER1:
        #ifdef _WIN32
          write_index += snprintf(
              display + write_index,
              sizeof(display) - write_index,
              "P1");
        #else
          write_index += snprintf(
              display + write_index,
              sizeof(display) - write_index,
              "🕺");
        #endif
        break;
      case PLAYER2:
      #ifdef _WIN32
          write_index += snprintf(
              display + write_index,
              sizeof(display) - write_index,
              "P2");
        #else
          write_index += snprintf(
              display + write_index,
              sizeof(display) - write_index,
              "🏃");
        #endif
        break;
      case PLAYER3:
        #ifdef _WIN32
          write_index += snprintf(
              display + write_index,
              sizeof(display) - write_index,
              "P3");
        #else
          write_index += snprintf(
              display + write_index,
              sizeof(display) - write_index,
              "🧍");
        #endif
        break;
      case PLAYER4:
        #ifdef _WIN32
          write_index += snprintf(
              display + write_index,
              sizeof(display) - write_index,
              "P4");
        #else
          write_index += snprintf(
              display + write_index,
              sizeof(display) - write_index,
              "💃");
        #endif
        break;
      }
    }
    write_index += snprintf(display + write_index, sizeof(display) - write_index, "\n");
  }
  printf("%s", display);
};

static void display_lives(int lives) {
  for (int i = 0; i < MAX_LIVES; i++) {
    if (i < lives) {
      #ifdef _WIN32
        printf("<3 ");
      #else
        printf("❤️ ");
      #endif
    } else {
      if (lives <= 0 && i == 0) {
        #ifdef _WIN32
          printf("XX ");
        #else
          printf("💀 ");
        #endif
      } else {
        printf("   ");
      }
    }
  }
}

void display_player_lives(players_t *players) {
  #ifdef _WIN32
    printf("P1: ");
  #else
    printf("🕺: ");
  #endif
  display_lives(players->player1.lives);
  printf("\t");
  #ifdef _WIN32
    printf("P2: ");
  #else
    printf("🏃: ");
  #endif
  display_lives(players->player2.lives);
  printf("\n\n");

  #ifdef _WIN32
    printf("P3: ");
  #else
    printf("🧍: ");
  #endif
  display_lives(players->player3.lives);
  printf("\t");
  #ifdef _WIN32
    printf("P4: ");
  #else
    printf("💃: ");
  #endif
  display_lives(players->player4.lives);
  printf("\n\n\n");
};

void clear_display()
{
#ifdef _WIN32
  system("cls");
#else
  system("clear");
#endif
}
