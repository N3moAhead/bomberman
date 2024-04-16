#ifndef SDL_DISPLAY_H
#define SDL_DISPLAY_H
#include <SDL2/SDL.h>
#include "types.h"

char quitted_game();
void prepare_scene();
void present_scene();
void display_map(block_t **map, players_t players, int game_round);

#endif