#include "player1.h"
#include "types.h"

// This function gets called on each tick
player_action_t get_player_1_action(block_t **map, players_t *players, int game_round) {
  player_action_t actions[3] = {PLANT_BOMB, MOVE_DOWN, MOVE_RIGHT};
  return actions[game_round % 3];
}
