#include "n3mo_bot_v1.h"
#include <string.h>
#include "../types.h"
#include "../constants.h"
#include "../util/player_helper.h"

/**
 * Returns 0 or 1 based on if the given position is dangerous
 */
static char is_pos_dangerous(block_t **map, cell_pos_t pos)
{
  char is_dangerous = 0;
  // Check the current field
  if (is_bomb(map, pos))
  {
    // There is a dangerous bomb on the current field
    return 1;
  }
  // Check left
  if (is_bomb(map, (cell_pos_t){.x = pos.x - 1, .y = pos.y}) || is_bomb(map, (cell_pos_t){.x = pos.x - 2, .y = pos.y}))
  {
    return 1;
  }
  // Check right
  if (is_bomb(map, (cell_pos_t){.x = pos.x + 1, .y = pos.y}) || is_bomb(map, (cell_pos_t){.x = pos.x + 2, .y = pos.y}))
  {
    return 1;
  }
  // Check top
  if (is_bomb(map, (cell_pos_t){.x = pos.x, .y = pos.y - 1}) || is_bomb(map, (cell_pos_t){.x = pos.x, .y = pos.y - 2}))
  {
    return 1;
  }
  // Check bottom
  if (is_bomb(map, (cell_pos_t){.x = pos.x, .y = pos.y + 1}) || is_bomb(map, (cell_pos_t){.x = pos.x, .y = pos.y + 2}))
  {
    return 1;
  }
  return is_dangerous;
}

/**
 * Checks if a field is save to flee there
 * If so the function returns 1
 * If not the function returns 0
 */
static char is_flee_direction_safe(block_t **map, cell_pos_t pos)
{
  cell_pos_t gated_pos = get_gated_position(pos);
  if (!is_blocked(map, gated_pos) && !is_pos_dangerous(map, gated_pos))
  {
    return 1;
  }
  return 0;
}

static int get_distance_to_closest_bomb(block_t **map, cell_pos_t pos)
{
  // If no bomb on the field is found a ridiculous
  // far away distance will be returned as distance
  int closest_distance = MAP_HEIGHT * MAP_WIDTH;
  for (int row = 0; row < MAP_HEIGHT; row++)
  {
    for (int col = 0; col < MAP_WIDTH; col++)
    {
      if (is_bomb(map, (cell_pos_t){.x = col, .y = row}))
      {
        int distance = get_distance((cell_pos_t){.x = col, .y = row}, pos);
        if (closest_distance == MAP_HEIGHT * MAP_WIDTH || closest_distance > distance)
        {
          closest_distance = distance;
        }
      }
    }
  }
  return closest_distance;
}

static char could_a_bomb_reach_field(cell_pos_t from, cell_pos_t to)
{
  // TODO This function currently ignores walls which is pretty bad so i will have to refactor that later on!
  //  same field
  if (from.x == to.x && from.y == to.y)
    return 1;
  // up & down
  if (
      from.x == to.x && (from.y - 1 == to.y || from.y - 2 == to.y || from.y + 1 == to.y || from.y + 2 == to.y))
    return 1;
  // left & right
  if (
      from.y == to.y && (from.x - 1 == to.x || from.x - 2 == to.x || from.x + 1 == to.x || from.x + 2 == to.x))
    return 1;

  return 0;
}

static char is_player_in_range(players_t players, player_t bot)
{
  int write_index = 0;
  player_t player_positions[3];
  if (players.player1.id != bot.id)
    player_positions[write_index++] = players.player1;
  if (players.player2.id != bot.id)
    player_positions[write_index++] = players.player2;
  if (players.player3.id != bot.id)
    player_positions[write_index++] = players.player3;
  if (players.player4.id != bot.id)
    player_positions[write_index++] = players.player4;

  for (int i = 0; i < 3; i++)
  {
    if (player_positions[i].lives > 0 && could_a_bomb_reach_field(bot.cell_pos, player_positions[i].cell_pos))
    {
      return 1;
    }
  }

  return 0;
}

/**
 * WARNING: This function only works for a depth of 3 fields
 */
static player_action_t any_save_path(block_t **map, cell_pos_t pos, player_action_t moved, int n)
{
  // Check if the current position is blocked
  if (is_blocked(map, pos))
  {
    return NONE;
  }
  // Check if we ran out of depth
  if (n == 0)
  {
    if (is_flee_direction_safe(map, pos)) {
      return moved;
    }
    return NONE;
  }
  // Check if the current position is already safe
  if (is_flee_direction_safe(map, pos))
  {
    return moved;
  }

  // MOVE TO DIRECTIONS

  // Move up
  if (
      moved != MOVE_DOWN && any_save_path(map, (cell_pos_t){.x = pos.x, .y = pos.y - 1}, MOVE_UP, n - 1) != NONE)
  {
    return MOVE_UP;
  }
  // Move down
  if (moved != MOVE_UP && any_save_path(map, (cell_pos_t){.x = pos.x, .y = pos.y + 1}, MOVE_DOWN, n - 1) != NONE)
  {
    return MOVE_DOWN;
  }
  // Move left
  if (moved != MOVE_RIGHT && any_save_path(map, (cell_pos_t){.x = pos.x - 1, .y = pos.y}, MOVE_LEFT, n - 1) != NONE)
  {
    return MOVE_LEFT;
  }
  // Move right
  if (moved != MOVE_LEFT && any_save_path(map, (cell_pos_t){.x = pos.x + 1, .y = pos.y}, MOVE_RIGHT, n - 1) != NONE)
  {
    return MOVE_RIGHT;
  }

  return NONE;
}

char box_around(block_t **map, cell_pos_t pos) {
  if (map[pos.y + 1][pos.x] == BOX) return 1;
  if (map[pos.y - 1][pos.x] == BOX) return 1;
  if (map[pos.y][pos.x + 1] == BOX) return 1;
  if (map[pos.y][pos.x - 1] == BOX) return 1;
  
  return 0;
}

/**
 * The current field is dangerous so the bot has to flee from this field
 * This function should return a direction to run away
 */
static player_action_t get_flee_direction(block_t **map, player_t player)
{
  player_action_t possible_move = any_save_path(map, player.cell_pos, NONE, 4);
  if (possible_move != NONE) {
    return possible_move;
  }
  return NONE;
}


/**
 * This function checks if its safe to plant a bomb
 *
 * WARNING: This function is only going to look
 * at the current field and is not going to predict the future field
 *
 * This function works by checking if there is a safe field which
 * is reachable in 3 or less turns. (USING BFS? nah just simple recursion :) lets keep it simple)
 * If so the function will return PLANT_BOMB.
 * If not the function will return NONE.
 */
static player_action_t plant_bomb(block_t **map, player_t bot)
{
  // I have to imagine that there is a bomb on the current field otherwise its probably always safe XD
  block_t before = map[bot.cell_pos.y][bot.cell_pos.x];
  map[bot.cell_pos.y][bot.cell_pos.x] = BOMB1;
  // CHECK FOR A SAVE ESCAPE PATH
  // top
  if (any_save_path(map, (cell_pos_t){.x = bot.cell_pos.x, .y = bot.cell_pos.y - 1}, MOVE_UP, 2) != NONE) {
    map[bot.cell_pos.y][bot.cell_pos.x] = before;
    return PLANT_BOMB;
  }
  // bottom
  if (any_save_path(map, (cell_pos_t){.x = bot.cell_pos.x, .y = bot.cell_pos.y + 1}, MOVE_DOWN, 2) != NONE) {
    map[bot.cell_pos.y][bot.cell_pos.x] = before;
    return PLANT_BOMB;
  }
  // left
  if (any_save_path(map, (cell_pos_t){.x = bot.cell_pos.x - 1, .y = bot.cell_pos.y}, MOVE_LEFT, 2) != NONE) {
    map[bot.cell_pos.y][bot.cell_pos.x] = before;
    return PLANT_BOMB;
  }
  // right
  if (any_save_path(map, (cell_pos_t){.x = bot.cell_pos.x + 1, .y = bot.cell_pos.y}, MOVE_RIGHT, 2) != NONE) {
    map[bot.cell_pos.y][bot.cell_pos.x] = before;
    return PLANT_BOMB;
  }
  map[bot.cell_pos.y][bot.cell_pos.x] = before;
  return NONE;
}

static player_action_t get_move_towards_enemy(block_t **map, players_t players, player_t bot)
{
  int write_index = 0;
  player_t player_positions[3];
  if (players.player1.id != bot.id)
    player_positions[write_index++] = players.player1;
  if (players.player2.id != bot.id)
    player_positions[write_index++] = players.player2;
  if (players.player3.id != bot.id)
    player_positions[write_index++] = players.player3;
  if (players.player4.id != bot.id)
    player_positions[write_index++] = players.player4;

  int closest_player_index = -1;
  int closest_player_distance;
  for (int i = 0; i < 3; i++)
  {
    int current_distance = get_distance(bot.cell_pos, player_positions[i].cell_pos);
    if (player_positions[i].lives > 0 && (closest_player_index == -1 || closest_player_distance > current_distance))
    {
      closest_player_index = i;
      closest_player_distance = current_distance;
    }
  }

  // TODO implement a proper path finding algorithm to get to the nearest player
  if (player_positions[closest_player_index].cell_pos.x != bot.cell_pos.x)
  {
    if (player_positions[closest_player_index].cell_pos.x < bot.cell_pos.x && is_flee_direction_safe(map, (cell_pos_t){.x = bot.cell_pos.x - 1, .y = bot.cell_pos.y}))
    {
      return MOVE_LEFT;
    }
    else if (player_positions[closest_player_index].cell_pos.x > bot.cell_pos.x && is_flee_direction_safe(map, (cell_pos_t){.x = bot.cell_pos.x + 1, .y = bot.cell_pos.y}))
    {
      return MOVE_RIGHT;
    }
  }

  if (player_positions[closest_player_index].cell_pos.y != bot.cell_pos.y)
  {
    if (player_positions[closest_player_index].cell_pos.y < bot.cell_pos.y && is_flee_direction_safe(map, (cell_pos_t){.x = bot.cell_pos.x, .y = bot.cell_pos.y - 1}))
    {
      return MOVE_UP;
    }
    else if (player_positions[closest_player_index].cell_pos.y > bot.cell_pos.y && is_flee_direction_safe(map, (cell_pos_t){.x = bot.cell_pos.x, .y = bot.cell_pos.y + 1}))
    {
      return MOVE_DOWN;
    }
  }

  if (box_around(map, bot.cell_pos)) {
    // Has to check if there is a way to escape when planting a bomb here
    // im going to give him 3 turns and then he has to be safe if that is
    // not possible im not allowing him to plant a bomb here
    // this has to be implemented everywhere :smile: hahaha
    // return plant_bomb(map, bot);
    return plant_bomb(map, bot);
  }
  return NONE;
}

player_action_t get_bot_move(block_t **map, players_t *players, int game_round, player_t bot)
{
  // Check if the current field is dangerous
  if (is_pos_dangerous(map, bot.cell_pos))
  {
    // Its dangerous we have to make it out of here!
    return get_flee_direction(map, bot);
  }
  else
  {
    // Its not dangerous lets go ahead
    // Check if an enemy is in reach
    if (is_player_in_range(*players, bot))
    {
      // An enemy is in reach we can plant a bomb
      return plant_bomb(map, bot);
    }
    else
    {
      // No enemy is in range lets move towards the closest player
      return get_move_towards_enemy(map, *players, bot);
    };
  }
  return NONE;
}

void get_bot_description(char *bot_name)
{
  strcpy(bot_name, "N3moAhead v1");
}
