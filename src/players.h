#ifndef PLAYERS_H
#define PLAYERS_H
#include "types.h"

players_t *init_players();

void copy_players(players_t *dest, players_t *players);

void free_players(players_t *players);

#endif