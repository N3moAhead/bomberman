#ifndef PLAYER_HELPER_H
#define PLAYER_HELPER_H
#include "types.h"

int abs_int(int value);
int get_distance(int from_x, int from_y, int to_x, int to_y);
int gated_int(int value, int max, int min);
char is_bomb(block_t **map, int pos_x, int pos_y);
char is_wall(block_t **map, int pos_x, int pos_y);

#endif