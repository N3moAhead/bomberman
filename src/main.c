#include <stdlib.h>
#include <stdio.h>
#include "constants.h"
#include "types.h"
#include "display.h"
#include "map.h"
#include "players.h"
#include "player1.h"
#include "player2.h"
#include "player3.h"
#include "player4.h"

int main() {
  /**
   * THE MAP
   * A 2D Array composed by the block_t type.
   * Warning: You wont find the positions of players in this object.
   * If you want to get the positions of the players you should look in
   * the players object.
  */
  block_t **map = init_map();
  /** 
   * A copy of the current map that will be given to the user functions
   * It will be resetted after every user function has been called
   * because some evil player might try to modify it XD 
  */
  block_t **map_copy = init_map();
  /**
   * THE PLAYERS
   * A struct that holds the position of all 4 players.
  */
  players_t *players = init_players();
  /**
   * A copy of the current players object that will be given to the user functions
   * It will be resetted after every user function has been called
   * because some evil player might try to modify it XD
  */
  players_t *players_copy = init_players();

  char game_is_running = 1;
  int index = 0;
  while (game_is_running) {
    // GETTING THE PLAYER INPUT
    // player 1
    copy_map(map_copy, map);
    copy_players(players_copy, players);
    player_action_t player1_action = get_player_1_action(map_copy, players_copy);

    // player 2
    copy_map(map_copy, map);
    copy_players(players_copy, players);
    player_action_t player2_action = get_player_2_action(map_copy, players_copy);

    // player 3
    copy_map(map_copy, map);   
    copy_players(players_copy, players);
    player_action_t player3_action = get_player_3_action(map_copy, players_copy);

    // player 4
    copy_map(map_copy, map);
    copy_players(players_copy, players);
    player_action_t player4_action = get_player_4_action(map_copy, players_copy);

    // UPDATE THE MAP
    update_map(map);
    //TODO CHECKING THE PLAYER INPUT
    //TODO ADD THE PLAYER INPUT TO THE MAP
    //TODO CHECK PlAYER HEALTH 
    // CLEAR THE DISPLAY
    clear_display();
    // DISPLAYING THE MAP
    /**
     * Im reusing the player map here to save a bit of memory
     * if I figure out that it is a bad idea i will change it later on
     */
    copy_map(map_copy, map);
    add_players(map_copy, players);
    display(map_copy);

    // only run the loop 10 time for now
    if (index++ > 10) {
      game_is_running = 0;
    }
  }
  return 0;
}