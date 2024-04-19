#include "sdl_helper.h"
#include "../globals.h"
#include "../constants.h"
#include "../types.h"
#include <stdlib.h>
#include <SDL2/SDL.h>
#include <SDL2/SDL_ttf.h>
#include <SDL2/SDL_image.h>

SDL_Texture *load_texture(char *filename)
{
  SDL_Texture *texture;
  printf("Loading Image: %s\n", filename);
  texture = IMG_LoadTexture(renderer, filename);
  if (!texture)
  {
    printf("Error while trying to load Image %s. SDL_Error: %s", filename, SDL_GetError());
    exit(EXIT_FAILURE);
  }
  return texture;
}

// Draw a given sprite
void blit(sprite_t sprite)
{
  SDL_Rect dest = {
      .x = sprite.pos.x,
      .y = sprite.pos.y,
      .w = DISPLAY_SPRITE_SIZE,
      .h = DISPLAY_SPRITE_SIZE};
  SDL_RenderCopy(renderer, sprite.texture, NULL, &dest);
}

// Draw a sprite from the texture atlas
void blit_from_atlas(vector_2d_t atlas_pos, vector_2d_t draw_pos)
{
  SDL_Rect src = {
      .x = atlas_pos.x,
      .y = atlas_pos.y,
      .w = ASSET_SPRITE_SIZE,
      .h = ASSET_SPRITE_SIZE};
  SDL_Rect dest = {
      .x = draw_pos.x,
      .y = draw_pos.y,
      .w = DISPLAY_SPRITE_SIZE,
      .h = DISPLAY_SPRITE_SIZE};
  SDL_RenderCopy(renderer, texture_atlas, &src, &dest);
}

void display_text(char *text, vector_2d_t draw_pos)
{
  SDL_Color color = {255, 255, 255, 255};
  SDL_Surface *surface = TTF_RenderText_Solid(font, text, color);
  SDL_Texture *texture = SDL_CreateTextureFromSurface(renderer, surface);
  int draw_offset = (int)((DISPLAY_SPRITE_SIZE - surface->w) / 2);
  SDL_Rect dest = {draw_pos.x + draw_offset, draw_pos.y, surface->w, surface->h};
  SDL_RenderCopy(renderer, texture, NULL, &dest);
}

void blit_custom_from_atlas(vector_2d_t atlas_pos, vector_2d_t atlas_size, vector_2d_t draw_pos, vector_2d_t draw_size)
{
  SDL_Rect src = {
      .x = atlas_pos.x,
      .y = atlas_pos.y,
      .w = atlas_size.x,
      .h = atlas_size.y};
  SDL_Rect dest = {
      .x = draw_pos.x,
      .y = draw_pos.y,
      .w = draw_size.x,
      .h = draw_size.y};
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

static void init_font()
{
  if (TTF_Init() != 0)
  {
    printf("Error while trying to setup TTF! SDL_Error: %s\n", SDL_GetError());
    exit(EXIT_FAILURE);
  }
  font = TTF_OpenFont("assets/font.ttf", (int)(DISPLAY_SPRITE_SIZE / 3));
  if (!font)
  {
    printf("Error while trying to load the font! %s\n", SDL_GetError());
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
  if (!IMG_Init(IMG_INIT_PNG | IMG_INIT_JPG))
  {
    printf("Error while trying to init SDL_Image! SDL_Error: %s\n", SDL_GetError());
    exit(EXIT_FAILURE);
  }

  init_sdl_window();
  init_sdl_renderer();
  init_font();
  // Init the texture atlas
  texture_atlas = load_texture("assets/texture_atlas.png");
};