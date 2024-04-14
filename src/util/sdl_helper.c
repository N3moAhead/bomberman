#include "sdl_helper.h"
#include "../globals.h"
#include "../constants.h"
#include "../types.h"
#include <stdlib.h>
#include <SDL2/SDL.h>
#include <SDL2/SDL_image.h>

SDL_Texture *load_texture(char *filename) {
  SDL_Texture *texture;
  printf("Loading Image: %s\n", filename);
  texture = IMG_LoadTexture(renderer, filename);
  if (!texture) {
    printf("Error while trying to load Image %s. SDL_Error: %s", filename, SDL_GetError());
    exit(EXIT_FAILURE);
  }
  return texture;
}

// Draw a given sprite
void blit(sprite_t sprite) {
  SDL_Rect dest = {
    .x = sprite.pos.x,
    .y = sprite.pos.y,
    .w = DISPLAY_SPRITE_SIZE,
    .h = DISPLAY_SPRITE_SIZE
  };
  SDL_RenderCopy(renderer, sprite.texture, NULL, &dest);
}

// Draw a sprite from the texture atlas
void blit_from_atlas(vector_2d_t atlas_pos, vector_2d_t draw_pos)
{
  SDL_Rect src = {
    .x = atlas_pos.x,
    .y = atlas_pos.y,
    .w = ASSET_SPRITE_SIZE,
    .h = ASSET_SPRITE_SIZE
  };
  SDL_Rect dest = {
    .x = draw_pos.x,
    .y = draw_pos.y,
    .w = DISPLAY_SPRITE_SIZE,
    .h = DISPLAY_SPRITE_SIZE
  };
  SDL_RenderCopy(renderer, texture_atlas, &src, &dest);
}

static void init_sdl_window()
{
  window = SDL_CreateWindow(
      "Bomberman",
      SDL_WINDOWPOS_CENTERED,
      SDL_WINDOWPOS_CENTERED,
      WINDOW_WIDTH,
      WINDOW_HEIGHT,
      0);
  if (!window)
  {
    printf("Error while trying to craete a new SDL Window! SDL_Error: %s \n", SDL_GetError());
    exit(EXIT_FAILURE);
  }
};

static void init_sdl_renderer()
{
  renderer = SDL_CreateRenderer(window, -1, 0);
  if (!renderer)
  {
    printf("Error while trying to create a new SDL Renderer! SDL_Error: %s \n", SDL_GetError());
    exit(EXIT_FAILURE);
  }
}

void setup_sdl()
{
  /** Init the SDL Packages */
  if (SDL_Init(SDL_INIT_EVERYTHING) != 0)
  {
    printf("Error while trying to setup Sdl! SDL_Error: %s\n", SDL_GetError());
    exit(EXIT_FAILURE);
  }
  if (!IMG_Init(IMG_INIT_PNG | IMG_INIT_JPG)) {
    printf("Error while trying to init SDL_Image! SDL_Error: %s\n", SDL_GetError());
    exit(EXIT_FAILURE);
  }

  init_sdl_window();
  init_sdl_renderer();
  // Init the texture atlas
  texture_atlas = load_texture("assets/texture_atlas.png");
};