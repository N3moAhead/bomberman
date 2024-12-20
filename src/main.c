#include <SDL2/SDL.h>
#include "constants.h"
#include "types.h"
#include "sdl_display.h"
#include "util/sdl_helper.h"
#include "map.h"
#include "players.h"
#include "player1.h"
#include "player2.h"
#include "player3.h"
#include "player4.h"
#include "globals.h"

int main()
{
  /**
   * SDL
   */
  setup_sdl();
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
  get_player1_bot_description(players->player1.bot_description.author_name);
  get_player2_bot_description(players->player2.bot_description.author_name);
  get_player3_bot_description(players->player3.bot_description.author_name);
  get_player4_bot_description(players->player4.bot_description.author_name);
  /**
   * A copy of the current players object that will be given to the user functions
   * It will be resetted after every user function has been called
   * because some evil player might try to modify it XD
   */
  players_t *players_copy = init_players();

  char game_is_running = 1;
  int game_round = 0;
  char action_valid = 0;
  while (game_is_running)
  {
    // CHECKING FOR EXIT
    if (quitted_game() == 1)
    {
      printf("The game has been stopped\n");
      game_is_running = 0;
      break;
    }
    // GETTING THE PLAYER INPUT
    // player 1
    copy_map(map_copy, map);
    copy_players(players_copy, players);
    player_action_t player1_action = players->player1.lives > 0 
      ? get_player_1_action(map_copy, players_copy, game_round)
      : NONE;

    action_valid = validate_action(map, &players->player1, player1_action);
    if (!action_valid)
      player1_action = NONE;

    // player 2
    copy_map(map_copy, map);
    copy_players(players_copy, players);
    player_action_t player2_action = players->player2.lives > 0
      ? get_player_2_action(map_copy, players_copy, game_round)
      : NONE;
    action_valid = validate_action(map, &players->player2, player2_action);
    if (!action_valid)
      player2_action = NONE;

    // player 3
    copy_map(map_copy, map);
    copy_players(players_copy, players);
    player_action_t player3_action = players->player3.lives > 0
      ? get_player_3_action(map_copy, players_copy, game_round) 
      : NONE;
    action_valid = validate_action(map, &players->player3, player3_action);
    if (!action_valid)
      player3_action = NONE;

    // player 4
    copy_map(map_copy, map);
    copy_players(players_copy, players);
    player_action_t player4_action = players->player4.lives > 0
      ? get_player_4_action(map_copy, players_copy, game_round)
      : NONE;
    action_valid = validate_action(map, &players->player4, player4_action);
    if (!action_valid)
      player4_action = NONE;

    // UPDATE THE MAP
    update_map(map);
    // DISPLAYING THE MAP
    copy_map(map_copy, map);
    display_map(map_copy, *players, player1_action, player2_action, player3_action, player4_action);
    // ADD THE PLAYER INPUT TO THE MAP
    apply_player_input(map, &players->player1, player1_action);
    apply_player_input(map, &players->player2, player2_action);
    apply_player_input(map, &players->player3, player3_action);
    apply_player_input(map, &players->player4, player4_action);
    // UPDATE THE PLAYER OBJECTS
    update_player(&players->player1, player1_action, map);
    update_player(&players->player2, player2_action, map);
    update_player(&players->player3, player3_action, map);
    update_player(&players->player4, player4_action, map);
    // CHECK IF ENOUGH PLAYERS ARE STILL ALIVE
    int alive_players = get_alive_player_count(players);
    // ending the game if only one player is left to play
    if (alive_players < 2)
    {
      game_is_running = 0;
    }
    delete_markers();
    game_round++;
  }
  return 0;
}
