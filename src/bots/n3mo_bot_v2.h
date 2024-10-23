#ifndef N3MO_BOT_V1_H
#define N3MO_BOT_V1_H
#include "../types.h"

player_action_t get_bot_move_v2(block_t **map, players_t *players, int game_round, player_t bot);
void get_bot_description_v2(char *bot_name);

#endif