#ifndef PLAYER_HELPER_H
#define PLAYER_HELPER_H
#include "types.h"

int abs_int(int value);
int get_distance(cell_pos_t from, cell_pos_t to);
int gated_int(int value, int max, int min);
cell_pos_t get_gated_position(cell_pos_t pos);
char is_bomb(block_t **map, cell_pos_t pos);
char is_wall(block_t **map, cell_pos_t pos);

#endif