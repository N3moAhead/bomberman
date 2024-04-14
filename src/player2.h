#ifndef PLAYER_2_H
#define PLAYER_2_H
#include "types.h"

player_action_t get_player_2_action(block_t **map, players_t *players, int game_round);
void get_player2_bot_description(char bot_name[50]);

#endif