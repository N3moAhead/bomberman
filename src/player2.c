#include "player2.h"
#include "types.h"

// This function gets called on each tick
player_action_t get_player_2_action(block_t **map, players_t *players, int game_round) {
  return NONE;
}

void get_player2_bot_description(char *bot_name)
{
  strcpy(bot_name, "BOT 2");
}