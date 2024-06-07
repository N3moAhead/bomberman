#include "player_helper.h"
#include "../constants.h"
#include "../types.h"

/**
 * Returns an absoluted integer for an given int
 */
int abs_int(int value)
{
  if (value < 0)
  {
    return value * -1;
  }
  return value;
}

/**
 * Returns the distance between to points on the map
 * WARNING! Its not calculating a walkable path
 * Its also not calculating a diagonal path
 */
int get_distance(cell_pos_t from, cell_pos_t to)
{
  return (abs_int(to.y - from.y) + abs_int(to.x - from.x));
}

/**
 * Makes sure that a given value stays inside of given
 * boundaries. If the given value surpasses the boundaries
 * the maximum or the minimum value is going to be returned
 */
int gated_int(int value, int max, int min)
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

/**
 * Makes sure that a given position stays inside of
 * the the map boundaries
 */
cell_pos_t get_gated_position(cell_pos_t pos)
{
  cell_pos_t adjusted_pos = {
    .x = gated_int(pos.x, MAP_WIDTH - 1, 0),
    .y = gated_int(pos.y, MAP_HEIGHT - 1, 0)
  };
  return adjusted_pos;
}

char is_bomb(block_t **map, cell_pos_t pos)
{
  cell_pos_t gated_pos = get_gated_position(pos);
  switch (map[gated_pos.y][gated_pos.x])
  {
  case BOMB1:
  case BOMB2:
  case BOMB3:
  case BOMB4:
  case BOMB5:
  case BOMB6:
  case BOMB7:
  case BOMB8:
  case BOMB9:
  case BOMB10:
    return 1;
  default:
    return 0;
  }
}

char is_wall(block_t **map, cell_pos_t pos)
{
  cell_pos_t gated_pos = get_gated_position(pos);
  if (map[gated_pos.y][gated_pos.x] == WALL)
  {
    return 1;
  }
  return 0;
}

char is_box(block_t **map, cell_pos_t pos)
{
  cell_pos_t gated_pos = get_gated_position(pos);
  if (map[gated_pos.y][gated_pos.x] == BOX)
  {
    return 1;
  }
  return 0;
}

char is_blocked(block_t **map, cell_pos_t pos)
{
  if (is_wall(map, pos) || is_box(map, pos))
  {
    return 1;
  }
  return 0;
}

