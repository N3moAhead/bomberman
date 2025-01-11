#include "stdio.h"
#include "debug_helper.h"
#include "../types.h"

void print_block_t(block_t block) {
  switch (block) {
    case BOMB1:
      printf("BOMB1\n");
      break; 
    case BOMB2:
      printf("BOMB2\n");
      break;
    case BOMB3:
      printf("BOMB3\n");
      break;
    case BOMB4:
      printf("BOMB4\n");
      break;
    case BOMB5:
      printf("BOMB5\n");
      break;
    case BOMB6:
      printf("BOMB6\n");
      break;
    case BOMB7:
      printf("BOMB7\n");
      break;
    case BOMB8:
      printf("BOMB8\n");
      break;
    case BOMB9:
      printf("BOMB9\n");
      break;
    case BOMB10:
      printf("BOMB10\n");
      break;
    case WALL:
      printf("WALL\n");
      break;
    case EXPLOSION:
      printf("EXPLOSION\n");
      break; 
    case AIR:
      printf("AIR\n");
      break;
    case BOX:
      printf("BOX\n");
      break;
  }
}

void print_player_action(player_action_t action) {
  switch (action) {
    case MOVE_UP:
      printf("MOVE_UP\n");
      break;
    case MOVE_DOWN:
      printf("MOVE_DOWN\n");
      break;
    case MOVE_LEFT:
      printf("MOVE_LEFT\n");
      break;
    case MOVE_RIGHT:
      printf("MOVE_RIGHT\n");
      break;
    case PLANT_BOMB:
      printf("PLANT_BOMB\n");
      break;
    case NONE:
      printf("NONE\n");
      break;
  }
}

void print_cell_pos(cell_pos_t pos) {
  printf("CELL_POS: row: %d, col: %d\n", pos.y, pos.x);
}