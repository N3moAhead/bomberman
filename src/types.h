#ifndef TYPES_H
#define TYPES_H

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
  PLAYER1,
  PLAYER2,
  PLAYER3,
  PLAYER4,
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
} block_t;

typedef struct bot_description {
  char author_name[50];
} bot_description_t;

/** Used for the cell position inside of the field grid */
typedef struct cell_pos {
  int x;
  int y;
} cell_pos_t;

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

#endif