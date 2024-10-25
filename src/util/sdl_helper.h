#ifndef SDL_HELPER_H
#define SDL_HELPER_H
#include <SDL2/SDL_image.h>
#include "../types.h"

SDL_Texture *load_texture(char *filename);
void blit(sprite_t sprite);
void setup_sdl();
void blit_from_atlas(vector_2d_t atlas_pos, vector_2d_t draw_pos);
void blit_custom_from_atlas(vector_2d_t atlas_pos, vector_2d_t atlas_size, vector_2d_t draw_pos, vector_2d_t draw_size);
void display_text(char *text, vector_2d_t draw_pos);
void display_marker(SDL_Color rect_color, vector_2d_t draw_pos, char *text);

#endif