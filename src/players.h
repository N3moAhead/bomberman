#ifndef PLAYERS_H
#define PLAYERS_H
#include "types.h"

players_t *init_players();

void copy_players(players_t *dest, players_t *players);

//! deprecated
void free_players(players_t *players);

char validate_action(
    block_t **map,
    player_t *player,
    player_action_t player_action);

void update_player(
    player_t *player,
    player_action_t player_action,
    block_t **map
);

int get_alive_player_count(players_t *players);

#endif