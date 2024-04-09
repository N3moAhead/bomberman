#include "player_helper.h"
#include "constants.h"
#include "types.h"

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
int get_distance(int from_x, int from_y, int to_x, int to_y)
{
  return (abs_int(to_y - from_y) + abs_int(to_x - from_x));
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

char is_bomb(block_t **map, int pos_x, int pos_y)
{
  int gated_x = gated_int(pos_x, MAP_WIDTH - 1, 0);
  int gated_y = gated_int(pos_y, MAP_HEIGHT - 1, 0);
  switch (map[gated_y][gated_x])
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

char is_wall(block_t **map, int pos_x, int pos_y)
{
  int gated_x = gated_int(pos_x, MAP_WIDTH - 1, 0);
  int gated_y = gated_int(pos_y, MAP_HEIGHT - 1, 0);
  if (map[gated_y][gated_x] == WALL)
  {
    return 1;
  }
  return 0;
}
