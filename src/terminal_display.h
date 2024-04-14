#ifndef TERMINAL_DISPLAY_H
#define TERMINAL_DISPLAY_H
#include "types.h"

void display(block_t **map);

void display_player_lives(players_t *players);

void clear_display();

#endif
