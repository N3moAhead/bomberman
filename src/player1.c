#include "player1.h"
#include "types.h"
#include "constants.h"
#include "bots/n3mo_bot_v1.h"
#include <string.h>

// This function gets called on each tick
player_action_t get_player_1_action(block_t **map, players_t *players, int game_round)
{
  return get_bot_move(map, players, game_round, players->player1);
}

void get_player1_bot_description(char bot_name[50])
{
  get_bot_description(bot_name);
}