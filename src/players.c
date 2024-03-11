#include <stdlib.h>
#include <stdio.h>
#include "players.h"
#include "types.h"
#include "constants.h"

players_t *allocate_players()
{
  players_t *new_players = (players_t *)malloc(sizeof(players_t));
  if (new_players == NULL)
  {
    printf("Could not create a new players object");
    exit(1);
  }
  return new_players;
}

players_t *init_players()
{
  players_t *new_players = allocate_players();

  // Player 1 starts top left
  new_players->player1.x = 1;
  new_players->player1.y = 1;

  // Player 2 starts top right
  new_players->player2.x = MAP_WIDTH - 2;
  new_players->player2.y = 1;

  // Player 3 starts bottom left
  new_players->player3.x = 1;
  new_players->player3.y = MAP_HEIGHT - 2;

  // Player 4 starts bottom right
  new_players->player4.x = MAP_WIDTH - 2;
  new_players->player4.y = MAP_HEIGHT - 2;

  return new_players;
}

/**
 * Copy a player object into another player object
 * This function does not allocate anything
 */
void copy_players(players_t *dest, players_t *players)
{
  // Player 1 starts top left
  dest->player1.x = players->player1.x;
  dest->player1.y = players->player1.y;

  // Player 2 starts top right
  dest->player2.x = players->player2.x;
  dest->player2.y = players->player2.y;

  // Player 3 starts bottom left
  dest->player3.x = players->player3.x;
  dest->player3.y = players->player3.y;

  // Player 4 starts bottom right
  dest->player4.x = players->player4.x;
  dest->player4.y = players->player4.y;
}

//! Depreacted dont use it!
void free_players(players_t *players)
{
  printf("The function free_players is deprecated and should not be used!");
  exit(1);
  free(players);
}

/**
 * Takes the current map, a player and a wanted action
 * if the action is valid 1 will be returned
 * if the action is invalid the function will return 0
*/
char validate_action(
    block_t **map,
    player_t *player,
    player_action_t player_action)
{
  switch (player_action) {
    case MOVE_UP:
      if (map[player->y - 1][player->x] != WALL)
        return 1;
      return 0;
    case MOVE_DOWN:
      if (map[player->y + 1][player->x] != WALL)
        return 1;
      return 0;
    case MOVE_LEFT:
      if (map[player->y][player->x - 1] != WALL)
        return 1;
      return 0;
    case MOVE_RIGHT:
      if (map[player->y][player->x + 1] != WALL)
        return 1;
      return 0;
    case NONE:
      return 1;
    case PLANT_BOMB:
      if (map[player->y][player->x] == AIR)
        return 1;
      return 0;
  }
}