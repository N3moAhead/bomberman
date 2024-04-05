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

void add_explosion(block_t **map, int row, int col) {
  // Aplying the explosion to the center
  map[row][col] = EXPLOSION;
  // Top
  if (map[row - 1][col] != WALL) {
    map[row - 1][col] = EXPLOSION;
    if (map[row - 2][col] != WALL) {
      map[row - 2][col] = EXPLOSION;
    }
  }
  // Bottom
  if (map[row + 1][col] != WALL) {
    map[row + 1][col] = EXPLOSION;
    if (map[row + 2][col] != WALL) {
      map[row + 2][col] = EXPLOSION;
    }
  }
  // Right
  if (map[row][col + 1] != WALL) {
    map[row][col + 1] = EXPLOSION;
    if (map[row][col + 2] != WALL) {
      map[row][col + 2] = EXPLOSION;
    }
  }
  // Left
  if (map[row][col - 1] != WALL) {
    map[row][col - 1] = EXPLOSION;
    if (map[row][col - 2] != WALL) {
      map[row][col - 2] = EXPLOSION;
    }
  }
}

void update_map(block_t **map, int game_round)
{
  // Im not completly sure how good this solution is
  // It makes it kind of hard to predict bombs in the game
  // but it could also be fun though XD
  // It would be unclear with this change how old a bomb is and when it will explode
  char update_bombs = game_round % 2;
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
        if (update_bombs)
          map[row][col] = BOMB2;
        break;
      case BOMB2:
        if (update_bombs)
          map[row][col] = BOMB3;
        break;
      case BOMB3:
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
      if (map[row][col] == BOMB3 && update_bombs)
        add_explosion(map, row, col);
    }
  }
}

void add_players(block_t **map, players_t *players)
{
  // Player 1
  if (players->player1.lives > 0)
    map[players->player1.y][players->player1.x] = PLAYER1;
  // Player 2
  if (players->player2.lives > 0)
    map[players->player2.y][players->player2.x] = PLAYER2;
  // Player 3
  if (players->player3.lives > 0)
    map[players->player3.y][players->player3.x] = PLAYER3;
  // Player 4
  if (players->player4.lives > 0)
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
      map[player->y][player->x] = BOMB1;
      break;
  }
}
