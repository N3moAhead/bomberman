#include "player4.h"
#include "types.h"
#include <string.h>

// This function gets called on each tick
player_action_t get_player_4_action(block_t **map, players_t *players, int game_round) {
  return NONE;
}

void get_player4_bot_description(char bot_name[50])
{
  strcpy(bot_name, "BOT 4");
}
