#ifndef GLOBALS_H
#define GLOBALS_H
#include <SDL2/SDL.h>
#include <SDL2/SDL_ttf.h>
#include "types.h"

extern SDL_Window *window;
extern SDL_Renderer *renderer;
extern SDL_Texture *texture_atlas;
extern TTF_Font *font;
extern atlas_positions_t texture_atlas_positions;

#endif
