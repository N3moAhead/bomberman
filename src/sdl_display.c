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

void display_map(block_t **map)
{
  for (int row = 0; row < MAP_HEIGHT; row++)
  {
    for (int col = 0; col < MAP_WIDTH; col++)
    {
      vector_2d_t draw_pos = {
          .y = row * DISPLAY_SPRITE_SIZE,
          .x = col * DISPLAY_SPRITE_SIZE,
      };
      switch (map[row][col])
      {
      case AIR:
        blit_from_atlas(texture_atlas_positions.air, draw_pos);
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
        blit_from_atlas(texture_atlas_positions.wall, draw_pos);
        break;
      case EXPLOSION:
        blit_from_atlas(texture_atlas_positions.explosion, draw_pos);
        break;
      case PLAYER1:
        blit_from_atlas(texture_atlas_positions.air, draw_pos);
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
