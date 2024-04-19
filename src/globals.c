#include <SDL2/SDL.h>
#include <SDL2/SDL_ttf.h>
#include "types.h"

SDL_Window *window = NULL;
SDL_Renderer *renderer = NULL;
SDL_Texture *texture_atlas = NULL;
TTF_Font *font = NULL;
const atlas_positions_t texture_atlas_positions = {
  .player1 = {
    .x = 0,
    .y = 96
  }, 
  .player2 = {
    .x = 0,
    .y = 128
  }, 
  .player3 = {
    .x = 0,
    .y = 160
  }, 
  .player4 = {
    .x = 0,
    .y = 192
  }, 
  .bomb = {
    .x = 0,
    .y = 64
  } , 
  .wall = {
    .x = 0,
    .y = 32
  }, 
  .explosion = {
    .x = 0,
    .y = 224
  } , 
  .air = {
    .x = 64,
    .y = 0
  }, 
};