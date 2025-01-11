#ifndef SDL_DISPLAY_H
#define SDL_DISPLAY_H
#include <SDL2/SDL.h>
#include "types.h"

char quitted_game();
void prepare_scene();
void present_scene();
void display_map(block_t **map, players_t players, player_action_t pl1_act, player_action_t pl2_act, player_action_t pl3_act, player_action_t pl4_act);

#endif