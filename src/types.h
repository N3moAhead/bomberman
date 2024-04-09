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
   * Bombs explode after 3 stages.
   * Each bomb gets updated after each tick and explode
   * after the third stage.
   */
  BOMB1, 
  BOMB2,
  BOMB3,
  WALL,
  EXPLOSION, 
  AIR,
} block_t;

typedef struct player {
  int id;
  int x;
  int y;
  int lives;
} player_t;

typedef struct players {
  player_t player1;
  player_t player2;
  player_t player3; 
  player_t player4; 
} players_t;

#endif