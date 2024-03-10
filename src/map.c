#include <stdlib.h>
#include "map.h"
#include "constants.h"
#include "types.h"

block_t **init_map()
{
  block_t **map = (block_t **)malloc(MAP_HEIGHT * sizeof(block_t *));
  for (int row = 0; row < MAP_HEIGHT; row++)
  {
    map[row] = (block_t *)malloc(MAP_WIDTH * sizeof(block_t));
  }

  // Init the whole map with air
  for (int row = 0; row < MAP_HEIGHT; row++)
  {
    for (int col = 0; col < MAP_WIDTH; col++) {
      if (col == 0 || col == MAP_WIDTH - 1 || row == 0 || row == MAP_HEIGHT - 1) {
        map[row][col] = WALL;
      } else {
        map[row][col] = AIR;
      }
    }
  }

  return map;
}