#include "sdl_display.h"
#include "constants.h"
#include <stdlib.h>
#include <SDL2/SDL.h>
#include <SDL2/SDL_image.h>
#include "util/sdl_helper.h"
#include "globals.h"
#include "types.h"

//! A very bad try to process the input
// But it should be okay here because im only using the events to
// quit the game! But i really have to change if i want to add a real
// user to control a character in the game
char quitted_game()
{
  SDL_Event event;
  char quitted_game = 0;

  // Polling events until we have handeld each input
  while (SDL_PollEvent(&event))
  {
    switch (event.type)
    {
    case SDL_QUIT:
      quitted_game = 1;
      break;
    case SDL_KEYDOWN:
      if (event.key.keysym.sym == SDLK_ESCAPE)
      {
        quitted_game = 1;
        break;
      }
    default:
      break;
    }
  }

  return quitted_game;
}

void prepare_scene()
{
  SDL_SetRenderDrawColor(renderer, 96, 128, 255, 255);
  SDL_RenderClear(renderer);
}

void present_scene()
{
  SDL_RenderPresent(renderer);
}

/** 
 * This function takes a position and checks for the sourounding
 * blocks if the current position is the same block
 * depending on the outcome this function will return the equivalent
 * block variant.
 * WARNING! This function does not work on the AIR TYPE!
 */
static block_variant_type_t get_block_variant(block_t **map, cell_pos_t pos) {
  block_t block = map[pos.y][pos.x];
  // The sourounding blocks
  block_t top = pos.y - 1 >= 0 ? map[pos.y - 1][pos.x] : AIR;
  block_t right = pos.x + 1 < MAP_WIDTH ? map[pos.y][pos.x + 1] : AIR;
  block_t bottom = pos.y + 1 < MAP_HEIGHT ? map[pos.y + 1][pos.x] : AIR;
  block_t left = pos.x - 1 >= 0 ? map[pos.y][pos.x - 1] : AIR;
  // CENTER_VARIANT
  if (top != block && right != block && bottom != block && left != block) {
    return CENTER_VARIANT;
  }
  // LEFT_RIGHT_VARIANT
  if (top != block && right == block && bottom != block && left == block) {
    return LEFT_RIGHT_VARIANT;
  }
  // TOP_BOTTOM_VARIANT
  if (top == block && right != block && bottom == block && left != block) {
    return TOP_BOTTOM_VARIANT;
  }
  // BOTTOM_RIGHT_VARIANT
  if (top != block && right == block && bottom == block && left != block) {
    return BOTTOM_RIGHT_VARIANT;
  }
  // TOP_RIGHT_VARIANT
  if (top == block && right == block && bottom != block && left != block) {
    return TOP_RIGHT_VARIANT;
  }
  // BOTTOM_LEFT_VARIANT
  if (top != block && right != block && bottom == block && left == block) {
    return BOTTOM_LEFT_VARIANT;
  }
  // TOP_LEFT_VARIANT
  if (top == block && right != block && bottom != block && left == block) {
    return TOP_LEFT_VARIANT;
  }

  printf("Could not make a decision in get_block_variant");
  exit(EXIT_FAILURE);
}

static void draw_wall(block_t **map, cell_pos_t pos, vector_2d_t draw_pos)
{
  block_variant_type_t variant = get_block_variant(map, pos);
  vector_2d_t atlas_pos = texture_atlas_positions.wall;
  switch (variant) {
    case CENTER_VARIANT:
      atlas_pos.x = ASSET_SPRITE_SIZE * 6;
      break;
    case LEFT_RIGHT_VARIANT:
      atlas_pos.x = 0;
      break;
    case TOP_BOTTOM_VARIANT:
      atlas_pos.x = ASSET_SPRITE_SIZE * 1;
      break;
    case BOTTOM_RIGHT_VARIANT:
      atlas_pos.x = ASSET_SPRITE_SIZE * 2;
      break;
    case TOP_RIGHT_VARIANT:
      atlas_pos.x = ASSET_SPRITE_SIZE * 4;
      break;
    case BOTTOM_LEFT_VARIANT:
      atlas_pos.x = ASSET_SPRITE_SIZE * 3;
      break;
    case TOP_LEFT_VARIANT:
      atlas_pos.x = ASSET_SPRITE_SIZE * 5;
      break;
  }
  blit_from_atlas(atlas_pos, draw_pos);
}

void display_map(block_t **map) {
  for (int row = 0; row < MAP_HEIGHT; row++)
  {
    for (int col = 0; col < MAP_WIDTH; col++)
    {
      vector_2d_t draw_pos = {
          .y = row * DISPLAY_SPRITE_SIZE,
          .x = col * DISPLAY_SPRITE_SIZE,
      };
      blit_from_atlas(texture_atlas_positions.air, draw_pos);
      switch (map[row][col])
      {
      case AIR:
        break;
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
        blit_from_atlas(texture_atlas_positions.bomb, draw_pos);
        break;
      case WALL:
        draw_wall(map, (cell_pos_t){.y = row, .x = col}, draw_pos);
        break;
      case EXPLOSION:
        blit_from_atlas(texture_atlas_positions.explosion, draw_pos);
        break;
      case PLAYER1:
        blit_from_atlas(texture_atlas_positions.player1, draw_pos);
        break;
      case PLAYER2:
        blit_from_atlas(texture_atlas_positions.player2, draw_pos);
        break;
      case PLAYER3:
        blit_from_atlas(texture_atlas_positions.player3, draw_pos);
        break;
      case PLAYER4:
        blit_from_atlas(texture_atlas_positions.player4, draw_pos);
        break;
      }
    }
  }
}
