#include <stdio.h>
#include <stdlib.h>
#include "map.h"
#include "constants.h"
#include "types.h"

block_t **init_map()
{
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
      case PLAYER1:
        break;
      case PLAYER2:
        break;
      case PLAYER3:
        break;
      case PLAYER4:
        break;
      case BOMB1:
        map[row][col] = BOMB2;
        break;
      case BOMB2:
        map[row][col] = BOMB3;
        break;
      case BOMB3:
        // TODO Add a real big explosion!
        // A normal explosion should be a star 2 in each direction from the center
        // until it hits a wall
        map[row][col] = EXPLOSION;
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
}

void add_players(block_t **map, players_t *players)
{
  // Player 1
  map[players->player1.y][players->player1.x] = PLAYER1;
  // Player 2
  map[players->player2.y][players->player2.x] = PLAYER2;
  // Player 3
  map[players->player3.y][players->player3.x] = PLAYER3;
  // Player 4
  map[players->player4.y][players->player4.x] = PLAYER4;
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
  switch(player_action) {
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
      map[player->y][player->y] = BOMB1;
      break;
  }
}