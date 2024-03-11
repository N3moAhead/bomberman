#ifndef MAP_H
#define MAP_H
#include "types.h"

block_t **init_map();

void copy_map(block_t **dest, block_t **map);

void free_map(block_t **map);

void update_map(block_t **map);

void add_players(block_t **map, players_t *players);

#endif