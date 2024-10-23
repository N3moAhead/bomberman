#include "n3mo_bot_v2.h"
#include <string.h>
#include "../types.h"
#include "../constants.h"
#include "../util/player_helper.h"
#include "stdio.h"
#include "../util/debug_helper.h"

/**
 * H E L P E R   F U N C T I O N S
 */

typedef struct turns {
  char top;
  char bot;
  char left;
  char right;
} turns_t;

typedef struct turn_value {
  player_action_t turn;
  int value;
} turn_value_t;

/**
 * These functions take a cell_position and modify into one common direction.
 * Example: get_top(some_position) // will return the position above some_position
 */
cell_pos_t get_top(cell_pos_t pos) { return (cell_pos_t){.x = pos.x, .y = pos.y - 1}; }
cell_pos_t get_bot(cell_pos_t pos) { return (cell_pos_t){.x = pos.x, .y = pos.y + 1}; }
cell_pos_t get_left(cell_pos_t pos) { return (cell_pos_t){.x = pos.x - 1, .y = pos.y}; }
cell_pos_t get_right(cell_pos_t pos) { return (cell_pos_t){.x = pos.x + 1, .y = pos.y}; }

/**
 * Expects the current map and a position as parameter.
 * Returns a struct which holds the information in which directions from a certain
 * point are not blocked.
 */
turns_t get_possible_turns(block_t **map, cell_pos_t pos) {
  turns_t possible_turns = {
    .top = is_blocked(map, get_top(pos)) ? 0 : 1,
    .bot = is_blocked(map, get_bot(pos)) ? 0 : 1,
    .left = is_blocked(map, get_left(pos)) ? 0 : 1,
    .right = is_blocked(map, get_right(pos)) ? 0 : 1
  };
  return possible_turns;
}

/**
 * This function takes the current map and a position as a parameter.
 * 1 is returned if a field is safe and 0 is returned if not.
 */
char is_field_safe(block_t **map, cell_pos_t pos) {
  if (is_bomb(map, pos)) return 0;
  if (is_bomb(map, get_top(pos)) || !is_blocked(map, get_top(pos)) && is_bomb(map, get_top(get_top(pos)))) return 0;
  if (is_bomb(map, get_right(pos)) || !is_blocked(map, get_right(pos)) && is_bomb(map, get_right(get_right(pos)))) return 0;
  if (is_bomb(map, get_bot(pos)) || !is_blocked(map, get_bot(pos)) && is_bomb(map, get_bot(get_bot(pos)))) return 0;
  if (is_bomb(map, get_left(pos)) || !is_blocked(map, get_left(pos)) && is_bomb(map, get_left(get_left(pos)))) return 0;
  return 1;
}

/**
 * Uses a recursive approach to walk towards the nearest save field.
 */
turn_value_t recursive_search(block_t **map, cell_pos_t pos, player_action_t last_turn, int depth) {
  if (is_field_safe(map, pos)) {
    return (turn_value_t){.turn = last_turn, .value = depth};
  }
  if (depth <= 0) {
    return (turn_value_t){.turn = last_turn, .value = -1000};
  }
  turns_t possible_turns = get_possible_turns(map, pos);
  turn_value_t best_turn = {.turn = NONE, .value = -1000};
  if (possible_turns.top) {
    turn_value_t top_turn = recursive_search(map, get_top(pos), MOVE_UP, depth - 1);
    if (top_turn.value > best_turn.value) best_turn = top_turn;
  }
  if (possible_turns.bot) {
    turn_value_t bot_turn = recursive_search(map, get_bot(pos), MOVE_DOWN, depth - 1);
    if (bot_turn.value > best_turn.value) best_turn = bot_turn;
  }
  if (possible_turns.left) {
    turn_value_t left_turn = recursive_search(map, get_left(pos), MOVE_LEFT, depth - 1);
    if (left_turn.value > best_turn.value) best_turn = left_turn;
  }
  if (possible_turns.right) {
    turn_value_t right_turn = recursive_search(map, get_right(pos), MOVE_RIGHT, depth - 1);
    if (right_turn.value > best_turn.value) best_turn = right_turn;
  }
  return best_turn;
}

/**
 * P L A Y E R   F U N C T I O N S
 */

/**
 * FLEE: Move towards the closest currently safe field
 */
player_action_t flee(block_t **map, player_t bot) {
  turn_value_t best_flee_turn = recursive_search(map, bot.cell_pos, NONE, 5);
  return best_flee_turn.turn;
}

/**
 * ATTACK: If SAFE && POSSIBILITY TO RUN AWAY && ENEMY IN REACH
 */
player_action_t attack() {
  return NONE;
}

/**
 * MOVE_TO_ENEMY: Walk towards the closest enemy
 */
player_action_t move_to_enemy() {
  return NONE;
}

player_action_t get_bot_move_v2(block_t **map, players_t *players, int game_round, player_t bot)
{
  /**
   * TODO: Implement the other functions attack and move_to_enemy
   */
  return flee(map, bot);
}

void get_bot_description_v2(char *bot_name)
{
  strcpy(bot_name, "N3moAhead v2");
}
