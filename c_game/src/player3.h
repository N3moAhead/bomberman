#ifndef PLAYER_3_H
#define PLAYER_3_H
#include "types.h"

player_action_t get_player_3_action(block_t **map, players_t *players, int game_round);
void get_player3_bot_description(char bot_name[50]);

#endif