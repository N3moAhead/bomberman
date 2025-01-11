#ifndef TYPES_H
#define TYPES_H
#include <SDL2/SDL_image.h>

typedef enum player_action {
  // Move the player a field up
  MOVE_UP,
  // Move the player a field down
  MOVE_DOWN,
  // Move the player a field to the left
  MOVE_LEFT,
  // Move the player a field to the right
  MOVE_RIGHT,
  // Place a new bomb at the current position of the player
  PLANT_BOMB,
  // Do nothing this round, just chilling a bit
  NONE,
} player_action_t;

typedef enum block {
  /**
   * Bombs explode after 10 stages.
   * Each bomb gets updated after each tick and explode
   * after the tenth stage.
   */
  BOMB1, 
  BOMB2,
  BOMB3,
  BOMB4,
  BOMB5,
  BOMB6,
  BOMB7,
  BOMB8,
  BOMB9,
  BOMB10,
  WALL,
  EXPLOSION, 
  AIR,
  BOX,
} block_t;

typedef enum block_variant_type {
  ALONE_VARIANT,
  CENTER_VARIANT,
  LEFT_RIGHT_VARIANT,
  TOP_BOTTOM_VARIANT,
  BOTTOM_RIGHT_VARIANT,
  TOP_RIGHT_VARIANT,
  BOTTOM_LEFT_VARIANT,
  TOP_LEFT_VARIANT,
  TOP_END_VARIANT,
  RIGHT_END_VARIANT,
  BOTTOM_END_VARIANT,
  LEFT_END_VARIANT,
  TOP_BOT_LEFT_VARIANT,
  TOP_RIGHT_LEFT_VARIANT,
  TOP_RIGHT_BOT_VARIANT,
  RIGHT_BOT_LEFT_VARIANT,
} block_variant_type_t;

typedef struct bot_description {
  char author_name[50];
} bot_description_t;

typedef struct vector_2d {
  int x;
  int y;
} vector_2d_t;
/** Macro to create vector 2d */
#define VECTOR_2D(row,col) ((vector_2d_t){.y = row, .x = col})

/** Used for the cell position inside of the field grid */
typedef vector_2d_t cell_pos_t;
/** Macro to create cell positions */
#define CELL_POS(row,col) ((cell_pos_t){.y = row, .x = col})

// Used to draw sprites outside of the texture atlas
typedef struct sprite {
  vector_2d_t pos;
  SDL_Texture *texture;
} sprite_t;

// Used for coordinates inside of the texture atlas
typedef struct atlas_sprite {
  vector_2d_t atlas_pos;
  vector_2d_t draw_pos;
} atlas_sprite_t;

typedef struct atlas_positions {
  vector_2d_t player1;
  vector_2d_t player2;
  vector_2d_t player3;
  vector_2d_t player4;
  vector_2d_t bomb; 
  vector_2d_t wall;
  vector_2d_t explosion; 
  vector_2d_t air;
  vector_2d_t box;
} atlas_positions_t;

typedef struct player {
  int id;
  cell_pos_t cell_pos;
  bot_description_t bot_description;
  int lives;
} player_t;

typedef struct players {
  player_t player1;
  player_t player2;
  player_t player3; 
  player_t player4; 
} players_t;

typedef struct marker_node {
  vector_2d_t pos;
  SDL_Color color;
  char* text;
  struct marker_node* next;
} marker_node_t;

/** SDL COLORS */
#define SDL_Red ((SDL_Color){255, 0, 0, 255})
#define SDL_Green ((SDL_Color){0, 255, 0, 255})
#define SDL_Blue ((SDL_Color){0, 0, 255, 255})

#endif