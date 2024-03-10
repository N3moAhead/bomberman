#include <stdlib.h>
#include <stdio.h>
#include "constants.h"
#include "types.h"
#include "display.h"
#include "map.h"

int main() {
  block_t **map = init_map();
  display(map);
  return 0;
}