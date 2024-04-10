#include "n3mo_bot_v1.h"
#include "../types.h"
#include "../constants.h"
#include "../player_helper.h"

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
  if (!is_wall(map, gated_pos) && !is_pos_dangerous(map, gated_pos))
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

/**
 * The current field is dangerous so the bot has to flee from this field
 * This function should return a direction to run away
 */
static player_action_t get_flee_direction(block_t **map, player_t player)
{
  /**
   * Lets first check the sourounding fields
   * If one of the fields is save directly move there
   */
  if (is_flee_direction_safe(map, (cell_pos_t){
                                      .x = player.cell_pos.x,
                                      .y = player.cell_pos.y - 1}))
  {
    return MOVE_UP;
  }
  // down
  if (is_flee_direction_safe(map, (cell_pos_t){
                                      .x = player.cell_pos.x,
                                      .y = player.cell_pos.y + 1}))
  {
    return MOVE_DOWN;
  }
  // left
  if (is_flee_direction_safe(map, (cell_pos_t){
                                      .x = player.cell_pos.x - 1,
                                      .y = player.cell_pos.y}))
  {
    return MOVE_LEFT;
  }
  // right
  if (is_flee_direction_safe(map, (cell_pos_t){
                                      .x = player.cell_pos.x + 1,
                                      .y = player.cell_pos.y}))
  {
    return MOVE_RIGHT;
  }
  /**
   * If no field is save, we will choose the first free
   * field to move to that is not a wall
   */
  // Todo check try to move to the field that is the furthes away from all the placed bombs
  player_action_t possible_options[4];
  if (!is_wall(map, (cell_pos_t){
                        .x = player.cell_pos.x,
                        .y = player.cell_pos.y - 1}))
  {
    possible_options[0] = MOVE_UP;
  }
  else
  {
    possible_options[0] = NONE;
  }
  // down
  if (!is_wall(map, (cell_pos_t){
                        .x = player.cell_pos.x,
                        .y = player.cell_pos.y + 1}))
  {
    possible_options[1] = MOVE_DOWN;
  }
  else
  {
    possible_options[1] = NONE;
  }
  // left
  if (!is_wall(map, (cell_pos_t){
                        .x = player.cell_pos.x - 1,
                        .y = player.cell_pos.y}))
  {
    possible_options[2] = MOVE_LEFT;
  }
  else
  {
    possible_options[2] = NONE;
  }
  // right
  if (!is_wall(map, (cell_pos_t){
                        .x = player.cell_pos.x + 1,
                        .y = player.cell_pos.y}))
  {
    possible_options[3] = MOVE_RIGHT;
  }
  else
  {
    possible_options[3] = NONE;
  }

  int best_option = -1;
  int longest_bomb_distance;
  // Evaluate going up
  if (possible_options[0] != NONE)
  {
    int distance = get_distance_to_closest_bomb(map, (cell_pos_t){
                                                          .x = player.cell_pos.x,
                                                          .y = player.cell_pos.y - 1});
    if (best_option == -1 || distance > longest_bomb_distance)
    {
      best_option = 0;
      longest_bomb_distance = distance;
    }
  }
  // Evaluate going down
  if (possible_options[1] != NONE)
  {
    int distance = get_distance_to_closest_bomb(map, (cell_pos_t){
                                                          .x = player.cell_pos.x,
                                                          .y = player.cell_pos.y + 1});
    if (best_option == -1 || distance > longest_bomb_distance)
    {
      best_option = 1;
      longest_bomb_distance = distance;
    }
  }
  // Evaluate going left
  if (possible_options[2] != NONE)
  {
    int distance = get_distance_to_closest_bomb(map, (cell_pos_t){
                                                          .x = player.cell_pos.x - 1,
                                                          .y = player.cell_pos.y});
    if (best_option == -1 || distance > longest_bomb_distance)
    {
      best_option = 2;
      longest_bomb_distance = distance;
    }
  }
  // Evaluate going right
  if (possible_options[3] != NONE)
  {
    int distance = get_distance_to_closest_bomb(map, (cell_pos_t){
                                                          .x = player.cell_pos.x + 1,
                                                          .y = player.cell_pos.y});
    if (best_option == -1 || distance > longest_bomb_distance)
    {
      best_option = 3;
      longest_bomb_distance = distance;
    }
  }

  if (best_option != -1)
  {
    return possible_options[best_option];
  }

  // If no field is possible the bot will just do nothing
  return NONE;
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
      return PLANT_BOMB;
    }
    else
    {
      // No enemy is in range lets move towards the closest player
      return get_move_towards_enemy(map, *players, bot);
    };
  }
  return NONE;
}
