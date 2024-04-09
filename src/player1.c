#include <stdio.h>
#include "player1.h"
#include "types.h"
#include "constants.h"

/**
 * BOT IDEA
 *
 * The bot should be able to perform 3 moves
 * 1. Running away if the place is dangerous
 * 2. Walking Towards Enemies
 * 3. Plant a bomb close to enemies
 */

/**
 * Validate that an given Value is inside of given boundaries
 */
static int gated_int(int value, int max, int min)
{
  if (value >= max)
  {
    return max;
  }
  if (value <= min)
  {
    return min;
  }
  return value;
}

static char is_bomb(block_t **map, int pos_x, int pos_y)
{
  int gated_x = gated_int(pos_x, MAP_WIDTH - 1, 0);
  int gated_y = gated_int(pos_y, MAP_HEIGHT - 1, 0);
  if (map[gated_y][gated_x] == BOMB1 || map[gated_y][gated_x] == BOMB2 || map[gated_y][gated_x] == BOMB3)
  {
    return 1;
  }
  return 0;
}

/**
 * Returns 0 or 1 based on if the given position is dangerous
 */
static char is_pos_dangerous(block_t **map, int pos_x, int pos_y)
{
  char is_dangerous = 0;
  // Check the current field
  if (is_bomb(map, pos_x, pos_y))
  {
    // There is a dangerous bomb on the current field
    return 1;
  }
  // Check left
  if (is_bomb(map, pos_x - 1, pos_y) || is_bomb(map, pos_x - 2, pos_y))
  {
    return 1;
  }
  // Check right
  if (is_bomb(map, pos_x + 1, pos_y) || is_bomb(map, pos_x + 2, pos_y))
  {
    return 1;
  }
  // Check top
  if (is_bomb(map, pos_x, pos_y - 1) || is_bomb(map, pos_x, pos_y - 2))
  {
    return 1;
  }
  // Check bottom
  if (is_bomb(map, pos_x, pos_y + 1) || is_bomb(map, pos_x, pos_y + 2))
  {
    return 1;
  }
  return is_dangerous;
}

static char is_wall(block_t **map, int pos_x, int pos_y)
{
  int gated_x = gated_int(pos_x, MAP_WIDTH - 1, 0);
  int gated_y = gated_int(pos_y, MAP_HEIGHT - 1, 0);
  if (map[gated_y][gated_x] == WALL)
  {
    return 1;
  }
  return 0;
}

/**
 * Checks if a field is save to flee there
 * If so the function returns 1
 * If not the function returns 0
 */
static char is_flee_direction_safe(block_t **map, int pos_x, int pos_y)
{
  int gated_x = gated_int(pos_x, MAP_WIDTH - 1, 0);
  int gated_y = gated_int(pos_y, MAP_HEIGHT - 1, 0);
  if (!is_wall(map, gated_x, gated_y) && !is_pos_dangerous(map, gated_x, gated_y))
  {
    return 1;
  }
  return 0;
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
  if (is_flee_direction_safe(map, player.x, player.y - 1))
  {
    return MOVE_UP;
  }
  // down
  if (is_flee_direction_safe(map, player.x, player.y + 1))
  {
    return MOVE_DOWN;
  }
  // left
  if (is_flee_direction_safe(map, player.x - 1, player.y))
  {
    return MOVE_LEFT;
  }
  // right
  if (is_flee_direction_safe(map, player.x + 1, player.y))
  {
    return MOVE_RIGHT;
  }
  /**
   * If no field is save, we will choose the first free
   * field to move to that is not a wall
   */
  if (!is_wall(map, player.x, player.y - 1))
  {
    return MOVE_UP;
  }
  // down
  if (!is_wall(map, player.x, player.y + 1))
  {
    return MOVE_DOWN;
  }
  // left
  if (!is_wall(map, player.x - 1, player.y))
  {
    return MOVE_LEFT;
  }
  // right
  if (!is_wall(map, player.x + 1, player.y))
  {
    return MOVE_RIGHT;
  }
  // If no field is possible the bot will just do nothing
  return NONE;
}

static char could_a_bomb_reach_field(int from_pos_x, int from_pos_y, int to_pos_x, int to_pos_y)
{
  // TODO This function currently ignores walls which is pretty bad so i will have to refactor that later on!
  //  same field
  if (from_pos_x == to_pos_x && from_pos_y == to_pos_y)
    return 1;
  // up & down
  if (
      from_pos_x == to_pos_x && (from_pos_y - 1 == to_pos_y || from_pos_y - 2 == to_pos_y || from_pos_y + 1 == to_pos_y || from_pos_y + 2 == to_pos_y))
    return 1;
  // left & right
  if (
      from_pos_y == to_pos_y && (from_pos_x - 1 == to_pos_x || from_pos_x - 2 == to_pos_x || from_pos_x + 1 == to_pos_x || from_pos_x + 2 == to_pos_x))
    return 1;

  return 0;
}

static char is_player_in_range(players_t players)
{
  int player_positions[3][3] = {
      {players.player2.x, players.player2.y, players.player2.lives},
      {players.player3.x, players.player3.y, players.player3.lives},
      {players.player4.x, players.player4.y, players.player4.lives},
  };

  for (int i = 0; i < 3; i++)
  {
    if (player_positions[i][2] > 0 && could_a_bomb_reach_field(players.player1.x, players.player1.y, player_positions[i][0], player_positions[i][1]))
    {
      return 1;
    }
  }

  return 0;
}

static int abs_int(int num)
{
  if (num < 0)
  {
    return num * -1;
  }
  return num;
}

static int get_distance(int from_x, int from_y, int to_x, int to_y)
{
  return (abs_int(to_y - from_y) + abs_int(to_x - from_x));
}

static player_action_t get_move_towards_enemy(block_t **map, players_t players)
{
  player_t player_positions[3] = {players.player2, players.player3, players.player4};

  int closest_player_index = -1;
  int closest_player_distance;
  for (int i = 0; i < 3; i++)
  {
    int current_distance = get_distance(players.player1.x, players.player1.y, player_positions[i].x, player_positions[i].y);
    if (player_positions[i].lives > 0 && (closest_player_index == -1 || closest_player_distance > current_distance))
    {
      closest_player_index = i;
      closest_player_distance = current_distance;
    }
  }

  // TODO implement a proper path finding algorithm to get to the nearest player
  if (player_positions[closest_player_index].x != players.player1.x)
  {
    if (player_positions[closest_player_index].x < players.player1.x && is_flee_direction_safe(map, players.player1.x - 1, players.player1.y))
    {
      return MOVE_LEFT;
    }
    else if (player_positions[closest_player_index].x > players.player1.x && is_flee_direction_safe(map, players.player1.x + 1, players.player1.y))
    {
      return MOVE_RIGHT;
    }
  }

  if (player_positions[closest_player_index].y != players.player1.y)
  {
    if (player_positions[closest_player_index].y < players.player1.y && is_flee_direction_safe(map, players.player1.x, players.player1.y - 1))
    {
      return MOVE_UP;
    }
    else if (player_positions[closest_player_index].y > players.player1.y && is_flee_direction_safe(map, players.player1.x, players.player1.y + 1))
    {
      return MOVE_DOWN;
    }
  }

  return NONE;
}

// This function gets called on each tick
player_action_t get_player_1_action(block_t **map, players_t *players, int game_round)
{
  // The current player
  player_t player = players->player1;
  // Check if the current field is dangerous
  if (is_pos_dangerous(map, player.x, player.y))
  {
    // Its dangerous we have to make it out of here!
    return get_flee_direction(map, player);
  }
  else
  {
    // Its not dangerous lets go ahead
    // Check if an enemy is in reach
    if (is_player_in_range(*players))
    {
      // An enemy is in reach we can plant a bomb
      return PLANT_BOMB;
    }
    else
    {
      // No enemy is in range lets move towards the closest player
      return get_move_towards_enemy(map, *players);
    };
  }
  return NONE;
}
