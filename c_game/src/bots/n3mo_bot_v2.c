#include "n3mo_bot_v2.h"
#include <string.h>
#include "../types.h"
#include "../constants.h"
#include "../util/player_helper.h"
#include "stdio.h"
#include "../util/debug_helper.h"
#include "../globals.h"

/**
 * C U S T O M   T Y P E S
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
 * H E L P E R   F U N C T I O N S
 */

/**
 * These functions take a cell_position and modify into one common direction.
 * Example: get_top(some_position) // will return the position above some_position
 */
cell_pos_t get_top(cell_pos_t pos) { return CELL_POS(pos.y - 1,pos.x); }
cell_pos_t get_bot(cell_pos_t pos) { return CELL_POS(pos.y + 1, pos.x); }
cell_pos_t get_left(cell_pos_t pos) { return CELL_POS(pos.y, pos.x - 1); }
cell_pos_t get_right(cell_pos_t pos) { return CELL_POS(pos.y, pos.x + 1); }

char cell_pos_equal(cell_pos_t pos1, cell_pos_t pos2) {
  if (pos1.x == pos2.x && pos1.y == pos2.y) {
    return 1;
  }
  return 0;
}

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
  if (is_bomb(map, get_top(pos)) || (!is_blocked(map, get_top(pos)) && is_bomb(map, get_top(get_top(pos))))) return 0;
  if (is_bomb(map, get_right(pos)) || (!is_blocked(map, get_right(pos)) && is_bomb(map, get_right(get_right(pos))))) return 0;
  if (is_bomb(map, get_bot(pos)) || (!is_blocked(map, get_bot(pos)) && is_bomb(map, get_bot(get_bot(pos))))) return 0;
  if (is_bomb(map, get_left(pos)) || (!is_blocked(map, get_left(pos)) && is_bomb(map, get_left(get_left(pos))))) return 0;
  return 1;
}

/**
 * Expects the current map and a position as parameter.
 * Returns a struct which holds the information in which directions from a certain
 * point are not blocked and SAFE!.
 */
turns_t get_possible_safe_turns(block_t **map, cell_pos_t pos) {
  cell_pos_t top = get_top(pos);
  cell_pos_t right = get_right(pos);
  cell_pos_t bot = get_bot(pos);
  cell_pos_t left = get_left(pos);
  turns_t possible_turns = {
    .top = ((is_blocked(map, top) == 0) && is_field_safe(map, top)) ? 1 : 0,
    .right = ((is_blocked(map, right) == 0) && is_field_safe(map, right)) ? 1 : 0,
    .bot = ((is_blocked(map, bot) == 0) && is_field_safe(map, bot)) ? 1 : 0,
    .left = ((is_blocked(map, left) == 0) && is_field_safe(map, left)) ? 1 : 0,
  };
  return possible_turns;
}

/**
 * Uses a recursive approach to walk towards the nearest safe field.
 */
turn_value_t recursive_search_safe_field(
  block_t **map,
  cell_pos_t pos,
  int depth,
  char visited[MAP_HEIGHT][MAP_WIDTH]
) {
  if (visited[pos.y][pos.x] == 1) {
    return (turn_value_t){.value = -1000};
  }
  if (is_field_safe(map, pos)) {
    return (turn_value_t){.value = depth};
  }
  if (depth <= 0) {
    return (turn_value_t){.value = -1000};
  }

  visited[pos.y][pos.x] = 1;
  turns_t possible_turns = get_possible_turns(map, pos);
  turn_value_t best_turn = {.turn = NONE, .value = -1000};
  if (possible_turns.top) {
    turn_value_t top_turn = recursive_search_safe_field(map, get_top(pos), depth - 1, visited);
    if (top_turn.value > best_turn.value) {
      best_turn.value = top_turn.value;
      best_turn.turn = MOVE_UP;
    };
  }
  if (possible_turns.bot) {
    turn_value_t bot_turn = recursive_search_safe_field(map, get_bot(pos), depth - 1, visited);
    if (bot_turn.value > best_turn.value) {
      best_turn.value = bot_turn.value;
      best_turn.turn = MOVE_DOWN;
    };
  }
  if (possible_turns.left) {
    turn_value_t left_turn = recursive_search_safe_field(map, get_left(pos), depth - 1, visited);
    if (left_turn.value > best_turn.value) {
      best_turn.value = left_turn.value;
      best_turn.turn = MOVE_LEFT;
    };
  }
  if (possible_turns.right) {
    turn_value_t right_turn = recursive_search_safe_field(map, get_right(pos), depth - 1, visited);
    if (right_turn.value > best_turn.value) {
      best_turn.value = right_turn.value;
      best_turn.turn = MOVE_RIGHT;
    };
  }
  return best_turn;
}

/**
 * Returns the next turn towards a specific point.
 * If the point is not accessible it will choose the point thats reachable and closest to it!
 * The value in the turn_value_t will be the distance to the specific point
 * Problem is, its never the best path, so the figure sometimes just runs in circles :sweat_smile:
 * I should implement a dijkstra for that or A* haha
 */
turn_value_t get_step_to_pos(block_t **map, cell_pos_t pos, cell_pos_t goal, char visited[MAP_HEIGHT][MAP_WIDTH]) {
  if (visited[pos.y][pos.x] == 1) {
    return (turn_value_t){.value = 1000};
  }
  visited[pos.y][pos.x] = 1;
  if (cell_pos_equal(pos, goal)) {
    return (turn_value_t){.value = 0};
  }

  // TODO Check if the value is fitting
  turn_value_t best_turn = {.value = get_distance(pos, goal), .turn = NONE};
  turns_t possible_turns = get_possible_safe_turns(map, pos);
  if (possible_turns.top) {
    turn_value_t top_turn = get_step_to_pos(map, get_top(pos), goal, visited);
    if (top_turn.value < best_turn.value) {
      best_turn.value = top_turn.value;
      best_turn.turn = MOVE_UP;
    }
  }
  if (possible_turns.right) {
    turn_value_t right_turn = get_step_to_pos(map, get_right(pos), goal, visited);
    if (right_turn.value < best_turn.value) {
      best_turn.value = right_turn.value;
      best_turn.turn = MOVE_RIGHT;
    }
  }
  if (possible_turns.bot) {
    turn_value_t bot_turn = get_step_to_pos(map, get_bot(pos), goal, visited);
    if (bot_turn.value < best_turn.value) {
      best_turn.value = bot_turn.value;
      best_turn.turn = MOVE_DOWN;
    }
  }
  if (possible_turns.left) {
    turn_value_t left_turn = get_step_to_pos(map, get_left(pos), goal, visited);
    if (left_turn.value < best_turn.value) {
      best_turn.value = left_turn.value;
      best_turn.turn = MOVE_LEFT;
    }
  }
  return best_turn;
}

cell_pos_t get_closest_player_pos(players_t *players, player_t bot) {
  int write_index = 0;
  player_t player_positions[3];
  if (players->player1.id != bot.id)
    player_positions[write_index++] = players->player1;
  if (players->player2.id != bot.id)
    player_positions[write_index++] = players->player2;
  if (players->player3.id != bot.id)
    player_positions[write_index++] = players->player3;
  if (players->player4.id != bot.id)
    player_positions[write_index++] = players->player4;

  int closest_player_distance = -1;
  cell_pos_t closest_pos;
  for (int i = 0; i < 3; i++)
  {
    int current_distance = get_distance(bot.cell_pos, player_positions[i].cell_pos);
    if (player_positions[i].lives > 0 && (closest_player_distance == -1 || current_distance < closest_player_distance))
    {
      closest_pos = player_positions[i].cell_pos;
      closest_player_distance = current_distance;
    }
  }
  return closest_pos;
}

/**
 * P L A Y E R   F U N C T I O N S
 */

/**
 * FLEE: Move towards the closest currently safe field
 */
player_action_t flee(block_t **map, player_t bot) {
  char visited[MAP_HEIGHT][MAP_WIDTH];
  for (int row = 0; row < MAP_HEIGHT; row++) {
    for (int col = 0; col < MAP_WIDTH; col++) {
      visited[row][col] = 0;
    }
  }
  turn_value_t best_flee_turn = recursive_search_safe_field(map, bot.cell_pos, 5, visited);
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
player_action_t move_to_enemy(block_t **map, players_t *players, player_t bot) {
  cell_pos_t closest_player_pos = {.y = 1, .x = 1};
  // cell_pos_t closest_player_pos = get_closest_player_pos(players, bot);
  char visited[MAP_HEIGHT][MAP_WIDTH];
  for (int row = 0; row < MAP_HEIGHT; row++) {
    for (int col = 0; col < MAP_WIDTH; col++) {
      visited[row][col] = 0;
    }
  }
  turn_value_t turn_value = get_step_to_pos(map, bot.cell_pos, closest_player_pos, visited);
  return turn_value.turn;
}

player_action_t get_bot_move_v2(block_t **map, players_t *players, int game_round, player_t bot)
{
  if (!is_field_safe(map, bot.cell_pos)) {
    return flee(map, bot);
  }
  return move_to_enemy(map, players, bot);
}

void get_bot_description_v2(char *bot_name)
{
  strcpy(bot_name, "N3moAhead v2");
}
