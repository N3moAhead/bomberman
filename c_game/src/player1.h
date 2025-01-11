#ifndef PLAYER_1_H
#define PLAYER_1_H
#include "types.h"

player_action_t get_player_1_action(block_t **map, players_t *players, int game_round);
void get_player1_bot_description(char bot_name[50]);

#endif