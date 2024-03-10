#ifndef TYPES_H
#define TYPES_H

typedef enum player_action {
  MOVE_UP,
  MOVE_DOWN,
  MOVE_LEFT,
  MOVE_RIGHT,
  PLANT_BOMB
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
  AIR
} block_t;

typedef struct {
  int x;
  int y;
} player;

#endif