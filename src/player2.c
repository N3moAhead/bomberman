#include "player2.h"
#include "types.h"
#include <string.h>
#include "bots/n3mo_bot_v2.h"

// This function gets called on each tick
player_action_t get_player_2_action(block_t **map, players_t *players, int game_round) {
  return get_bot_move_v2(map, players, game_round, players->player2);
}

void get_player2_bot_description(char bot_name[50])
{
  get_bot_description_v2(bot_name);
}