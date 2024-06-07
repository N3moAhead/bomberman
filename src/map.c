#include <stdio.h>
#include <stdlib.h>
#include <time.h>
#include "map.h"
#include "constants.h"
#include "types.h"

static char place_box(int row, int col)
{
  // No boxes around the player spawns
  if (
      // Player 1
      (col == 1 && row == 1) || (col == 2 && row == 1) || (col == 1 && row == 2)
      // Player 2
      || (col == MAP_WIDTH - 3 && row == 1) || (col == MAP_WIDTH - 2 && row == 1) || (col == MAP_WIDTH - 2 && row == 2)
      // Player 3
      || (col == 1 && row == MAP_HEIGHT - 3) || (col == 1 && row == MAP_HEIGHT - 2) || (col == 2 && row == MAP_HEIGHT - 2)
      // Player 4
      || (col == MAP_WIDTH - 3 && row == MAP_HEIGHT - 2) || (col == MAP_WIDTH - 2 && row == MAP_HEIGHT - 2) || (col == MAP_WIDTH - 2 && row == MAP_HEIGHT - 3))
    return 0;

  return rand() % 100 < BOX_SPAWN_RATE;
}

block_t **init_map()
{
  srand(time(NULL));
  block_t **map = (block_t **)malloc(MAP_HEIGHT * sizeof(block_t *));
  for (int row = 0; row < MAP_HEIGHT; row++)
  {
    map[row] = (block_t *)malloc(MAP_WIDTH * sizeof(block_t));
  }

  // Init the whole map with air
  for (int row = 0; row < MAP_HEIGHT; row++)
  {
    for (int col = 0; col < MAP_WIDTH; col++)
    {
      if (col == 0 || col == MAP_WIDTH - 1 || row == 0 || row == MAP_HEIGHT - 1)
      {
        map[row][col] = WALL;
      }
      else if (row % 2 == 0 && col % 2 == 0 && row > 0 && row < MAP_HEIGHT - 2 && col > 0 && col < MAP_WIDTH - 2)
      {
        map[row][col] = WALL;
      }
      else if (place_box(row, col))
      {
        map[row][col] = BOX;
      }
      else
      {
        map[row][col] = AIR;
      }
    }
  }

  return map;
}

/**
 * Copies a given map object into another map array
 * This function won't allocate memory
 */
void copy_map(block_t **dest, block_t **map)
{
  for (int row = 0; row < MAP_HEIGHT; row++)
  {
    for (int col = 0; col < MAP_WIDTH; col++)
    {
      dest[row][col] = map[row][col];
    }
  }
}

//! Deprecated
void free_map(block_t **map)
{
  printf("The function free_map is deprecated and should not be used!");
  exit(1);
  for (int row = 0; row < MAP_HEIGHT; row++)
  {
    free(map[row]);
  }
  free(map);
};

void add_explosion(block_t **map, int row, int col)
{
  // Aplying the explosion to the center
  map[row][col] = EXPLOSION;
  // Top
  if (map[row - 1][col] != WALL)
  {
    if (map[row - 2][col] != WALL && map[row - 1][col] != BOX)
    {
      map[row - 2][col] = EXPLOSION;
    }
    map[row - 1][col] = EXPLOSION;
  }
  // Bottom
  if (map[row + 1][col] != WALL)
  {
    if (map[row + 2][col] != WALL && map[row + 1][col] != BOX)
    {
      map[row + 2][col] = EXPLOSION;
    }
    map[row + 1][col] = EXPLOSION;
  }
  // Right
  if (map[row][col + 1] != WALL)
  {
    if (map[row][col + 2] != WALL && map[row][col + 1] != BOX)
    {
      map[row][col + 2] = EXPLOSION;
    }
    map[row][col + 1] = EXPLOSION;
  }
  // Left
  if (map[row][col - 1] != WALL)
  {
    if (map[row][col - 2] != WALL && map[row][col - 1] != BOX)
    {
      map[row][col - 2] = EXPLOSION;
    }
    map[row][col - 1] = EXPLOSION;
  }
}

void update_map(block_t **map)
{
  for (int row = 0; row < MAP_HEIGHT; row++)
  {
    for (int col = 0; col < MAP_WIDTH; col++)
    {
      /**
       * What i basically need is a switch case that performs diffrent kind of
       * actions based on the type of block it currently is
       * Some do nothing like wall or air they just stay but others can change!
       */
      switch (map[row][col])
      {
      case BOMB1:
        map[row][col] = BOMB2;
        break;
      case BOMB2:
        map[row][col] = BOMB3;
        break;
      case BOMB3:
        map[row][col] = BOMB4;
        break;
      case BOMB4:
        map[row][col] = BOMB5;
        break;
      case BOMB5:
        map[row][col] = BOMB6;
        break;
      case BOMB6:
        map[row][col] = BOMB7;
        break;
      case BOMB7:
        map[row][col] = BOMB8;
        break;
      case BOMB8:
        map[row][col] = BOMB9;
        break;
      case BOMB9:
        map[row][col] = BOMB10;
        break;
      case BOMB10:
        break;
      case WALL:
        break;
      case EXPLOSION:
        map[row][col] = AIR;
        break;
      case AIR:
        break;
      }
    }
  }
  // Thats a super stupid solution to add bombs
  // but i kind of want to get it done so it will stay for
  // the moment like that yikes :(
  for (int row = 0; row < MAP_HEIGHT; row++)
  {
    for (int col = 0; col < MAP_WIDTH; col++)
    {
      if (map[row][col] == BOMB10)
        add_explosion(map, row, col);
    }
  }
}

/**
 * Adding the player input to the map and updating the given map pointer
 * This function will not update the given player struct.
 * This function will also assume that the action it handles is already
 * validated and just executes it!
 */
void apply_player_input(
    block_t **map,
    player_t *player,
    player_action_t player_action)
{
  switch (player_action)
  {
  /**
   * Players are not displayed in the map so
   * we wont have to update them here
   */
  case MOVE_UP:
  case MOVE_DOWN:
  case MOVE_LEFT:
  case MOVE_RIGHT:
  case NONE:
    break;
  case PLANT_BOMB:
    map[player->cell_pos.y][player->cell_pos.x] = BOMB1;
    break;
  }
}
