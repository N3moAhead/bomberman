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
  new_players->player1.lives = 3;

  // Player 2 starts top right
  new_players->player2.x = MAP_WIDTH - 2;
  new_players->player2.y = 1;
  new_players->player2.lives = 3;

  // Player 3 starts bottom left
  new_players->player3.x = 1;
  new_players->player3.y = MAP_HEIGHT - 2;
  new_players->player3.lives = 3;

  // Player 4 starts bottom right
  new_players->player4.x = MAP_WIDTH - 2;
  new_players->player4.y = MAP_HEIGHT - 2;
  new_players->player4.lives = 3;

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
  switch (player_action)
  {
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

  return 0;
}

/**
 * This function takes a player struct and a player action and
 * updates the player struct with the player action
 * So the given player struct will be modified!
 * Oh and this function also assumes that the player action is already validated!
 */
void update_player(player_t *player, player_action_t player_action, block_t **map)
{
  // Update the player lives
  if (map[player->y][player->x] == EXPLOSION) {
    player->lives--;
  }
  // Update the player position
  switch (player_action)
  {
  case MOVE_UP:
    player->y--;
    break;
  case MOVE_DOWN:
    player->y++;
    break;
  case MOVE_LEFT:
    player->x--;
    break;
  case MOVE_RIGHT:
    player->x++;
    break;
  case NONE:
    break;
  case PLANT_BOMB:
    break;
  }
}

int get_alive_player_count(players_t *players) {
  char alive_players = 0;
  if (players->player1.lives > 0)
    alive_players++;
  if (players->player2.lives > 0)
    alive_players++;
  if (players->player3.lives > 0)
    alive_players++;
  if (players->player4.lives > 0)
    alive_players++;
  return alive_players;
}
