#include "sdl_display.h"
#include "constants.h"
#include <stdlib.h>
#include <SDL2/SDL.h>
#include <SDL2/SDL_image.h>
#include "util/sdl_helper.h"
#include "globals.h"
#include "types.h"

//! A very bad try to process the input
// But it should be okay here because im only using the events to
// quit the game! But i really have to change if i want to add a real
// user to control a character in the game
char quitted_game()
{
  SDL_Event event;
  char quitted_game = 0;

  // Polling events until we have handeld each input
  while (SDL_PollEvent(&event))
  {
    switch (event.type)
    {
    case SDL_QUIT:
      quitted_game = 1;
      break;
    case SDL_KEYDOWN:
      if (event.key.keysym.sym == SDLK_ESCAPE)
      {
        quitted_game = 1;
        break;
      }
    default:
      break;
    }
  }

  return quitted_game;
}

void prepare_scene()
{
  SDL_SetRenderDrawColor(renderer, 96, 128, 255, 255);
  SDL_RenderClear(renderer);
}

void present_scene()
{
  SDL_RenderPresent(renderer);
}

/**
 * This function takes a position and checks for the sourounding
 * blocks if the current position is the same block
 * depending on the outcome this function will return the equivalent
 * block variant.
 * WARNING! This function does not work on the AIR TYPE!
 */
static block_variant_type_t get_block_variant(block_t **map, cell_pos_t pos)
{
  block_t block = map[pos.y][pos.x];
  // The sourounding blocks
  block_t top = pos.y - 1 >= 0 ? map[pos.y - 1][pos.x] : AIR;
  block_t right = pos.x + 1 < MAP_WIDTH ? map[pos.y][pos.x + 1] : AIR;
  block_t bottom = pos.y + 1 < MAP_HEIGHT ? map[pos.y + 1][pos.x] : AIR;
  block_t left = pos.x - 1 >= 0 ? map[pos.y][pos.x - 1] : AIR;
  if (top == block && right == block && bottom == block && left == block)
    return CENTER_VARIANT;
  // ALONE_VARIANT
  if (top != block && right != block && bottom != block && left != block)
  {
    return ALONE_VARIANT;
  }
  // LEFT_RIGHT_VARIANT
  if (top != block && right == block && bottom != block && left == block)
  {
    return LEFT_RIGHT_VARIANT;
  }
  // TOP_BOTTOM_VARIANT
  if (top == block && right != block && bottom == block && left != block)
  {
    return TOP_BOTTOM_VARIANT;
  }
  // BOTTOM_RIGHT_VARIANT
  if (top != block && right == block && bottom == block && left != block)
  {
    return BOTTOM_RIGHT_VARIANT;
  }
  // TOP_RIGHT_VARIANT
  if (top == block && right == block && bottom != block && left != block)
  {
    return TOP_RIGHT_VARIANT;
  }
  // BOTTOM_LEFT_VARIANT
  if (top != block && right != block && bottom == block && left == block)
  {
    return BOTTOM_LEFT_VARIANT;
  }
  // TOP_LEFT_VARIANT
  if (top == block && right != block && bottom != block && left == block)
  {
    return TOP_LEFT_VARIANT;
  }
  // TOP_END_VARIANT
  if (top == block && right != block && bottom != block && left != block)
    return TOP_END_VARIANT;
  // RIGHT_END_VARIANT
  if (top != block && right == block && bottom != block && left != block)
    return RIGHT_END_VARIANT;
  // BOTTOM_END_VARIANT
  if (top != block && right != block && bottom == block && left != block)
    return BOTTOM_END_VARIANT;
  // LEFT_END_VARIANT
  if (top != block && right != block && bottom != block && left == block)
    return LEFT_END_VARIANT;
  // TOP_BOT_LEFT_VARIANT
  if (top == block && right != block && bottom == block && left == block)
    return TOP_BOT_LEFT_VARIANT;
  // TOP_RIGHT_LEFT_VARIANT
  if (top == block && right == block && bottom != block && left == block)
    return TOP_RIGHT_LEFT_VARIANT;
  // TOP_RIGHT_BOT_VARIANT
  if (top == block && right == block && bottom == block && left != block)
    return TOP_RIGHT_BOT_VARIANT;
  // RIGHT_BOT_LEFT_VARIANT
  if (top != block && right == block && bottom == block && left == block)
    return RIGHT_BOT_LEFT_VARIANT;

  printf("Could not make a decision in get_block_variant");
  exit(EXIT_FAILURE);
}

static void draw_explosion(block_t **map, cell_pos_t pos, vector_2d_t draw_pos, int animation_value)
{
  block_variant_type_t variant = get_block_variant(map, pos);
  vector_2d_t atlas_pos;
  switch (variant)
  {
  case CENTER_VARIANT:
    atlas_pos = (vector_2d_t){
        .x = ((int)(animation_value / 3)) * ASSET_SPRITE_SIZE,
        .y = 7 * ASSET_SPRITE_SIZE};
    break;
  case LEFT_RIGHT_VARIANT:
    atlas_pos = (vector_2d_t){
        .x = (animation_value / 3) * ASSET_SPRITE_SIZE,
        .y = 9 * ASSET_SPRITE_SIZE};
    break;
  case TOP_BOTTOM_VARIANT:
    atlas_pos = (vector_2d_t){
        .x = (animation_value / 3) * ASSET_SPRITE_SIZE,
        .y = 8 * ASSET_SPRITE_SIZE};
    break;
  case BOTTOM_RIGHT_VARIANT:
    atlas_pos = (vector_2d_t){
        .x = ((animation_value / 3) * ASSET_SPRITE_SIZE) + (3 * ASSET_SPRITE_SIZE),
        .y = 10 * ASSET_SPRITE_SIZE};
    break;
  case TOP_RIGHT_VARIANT:
    atlas_pos = (vector_2d_t){
        .x = ((animation_value / 3) * ASSET_SPRITE_SIZE) + (3 * ASSET_SPRITE_SIZE),
        .y = 9 * ASSET_SPRITE_SIZE};
    break;
  case BOTTOM_LEFT_VARIANT:
    atlas_pos = (vector_2d_t){
        .x = ((animation_value / 3) * ASSET_SPRITE_SIZE) + (3 * ASSET_SPRITE_SIZE),
        .y = 7 * ASSET_SPRITE_SIZE};
    break;
  case TOP_LEFT_VARIANT:
    atlas_pos = (vector_2d_t){
        .x = ((animation_value / 3) * ASSET_SPRITE_SIZE) + (3 * ASSET_SPRITE_SIZE),
        .y = 8 * ASSET_SPRITE_SIZE};
    break;
  case TOP_END_VARIANT:
    atlas_pos = (vector_2d_t){
        .x = (animation_value / 3) * ASSET_SPRITE_SIZE,
        .y = 12 * ASSET_SPRITE_SIZE};
    break;
  case RIGHT_END_VARIANT:
    atlas_pos = (vector_2d_t){
        .x = (animation_value / 3) * ASSET_SPRITE_SIZE,
        .y = 13 * ASSET_SPRITE_SIZE};
    break;
  case BOTTOM_END_VARIANT:
    atlas_pos = (vector_2d_t){
        .x = (animation_value / 3) * ASSET_SPRITE_SIZE,
        .y = 10 * ASSET_SPRITE_SIZE};
    break;
  case LEFT_END_VARIANT:
    atlas_pos = (vector_2d_t){
        .x = (animation_value / 3) * ASSET_SPRITE_SIZE,
        .y = 11 * ASSET_SPRITE_SIZE};
    break;
  case TOP_BOT_LEFT_VARIANT:
    atlas_pos = (vector_2d_t){
        .x = ((animation_value / 3) * ASSET_SPRITE_SIZE) + (3 * ASSET_SPRITE_SIZE),
        .y = 11 * ASSET_SPRITE_SIZE};
    break;
  case TOP_RIGHT_LEFT_VARIANT:
    atlas_pos = (vector_2d_t){
        .x = ((animation_value / 3) * ASSET_SPRITE_SIZE) + (3 * ASSET_SPRITE_SIZE),
        .y = 12 * ASSET_SPRITE_SIZE};
    break;
  case TOP_RIGHT_BOT_VARIANT:
    atlas_pos = (vector_2d_t){
        .x = ((animation_value / 3) * ASSET_SPRITE_SIZE) + (3 * ASSET_SPRITE_SIZE),
        .y = 13 * ASSET_SPRITE_SIZE};
    break;
  case RIGHT_BOT_LEFT_VARIANT:
    atlas_pos = (vector_2d_t){
        .x = ((animation_value / 3) * ASSET_SPRITE_SIZE) + (3 * ASSET_SPRITE_SIZE),
        .y = 14 * ASSET_SPRITE_SIZE};
    break;
  case ALONE_VARIANT:
    break;
  }
  blit_from_atlas(atlas_pos, draw_pos);
}

static void draw_wall(block_t **map, cell_pos_t pos, vector_2d_t draw_pos)
{
  block_variant_type_t variant = get_block_variant(map, pos);
  vector_2d_t atlas_pos = texture_atlas_positions.wall;
  switch (variant)
  {
  case ALONE_VARIANT:
    atlas_pos.x = ASSET_SPRITE_SIZE * 6;
    break;
  case LEFT_RIGHT_VARIANT:
    atlas_pos.x = 0;
    break;
  case TOP_BOTTOM_VARIANT:
    atlas_pos.x = ASSET_SPRITE_SIZE * 1;
    break;
  case BOTTOM_RIGHT_VARIANT:
    atlas_pos.x = ASSET_SPRITE_SIZE * 2;
    break;
  case TOP_RIGHT_VARIANT:
    atlas_pos.x = ASSET_SPRITE_SIZE * 4;
    break;
  case BOTTOM_LEFT_VARIANT:
    atlas_pos.x = ASSET_SPRITE_SIZE * 3;
    break;
  case TOP_LEFT_VARIANT:
    atlas_pos.x = ASSET_SPRITE_SIZE * 5;
    break;
  default:
    // Lets just draw the center variant for now!
    atlas_pos.x = ASSET_SPRITE_SIZE * 6;
    break;
  }
  blit_from_atlas(atlas_pos, draw_pos);
}

static int get_player_animation_movement_offset(player_action_t pl_act)
{
  switch (pl_act)
  {
  case MOVE_UP:
    return (9 * ASSET_SPRITE_SIZE);
  case MOVE_RIGHT:
    return (6 * ASSET_SPRITE_SIZE);
  case MOVE_LEFT:
    return (3 * ASSET_SPRITE_SIZE);
  case NONE:
    return ASSET_SPRITE_SIZE;
  case MOVE_DOWN:
  case PLANT_BOMB:
    return 0;
  default:
    return 0;
  }
}

static vector_2d_t get_player_draw_pos(player_t player, player_action_t pl_act, int anim_value)
{
  vector_2d_t draw_pos = {
      .x = player.cell_pos.x * DISPLAY_SPRITE_SIZE,
      .y = player.cell_pos.y * DISPLAY_SPRITE_SIZE};
  int move_distance = (int)DISPLAY_SPRITE_SIZE / 9;
  int movement = move_distance * (anim_value + 1);
  switch (pl_act)
  {
  case MOVE_UP:
    draw_pos.y -= movement;
    break;
  case MOVE_RIGHT:
    draw_pos.x += movement;
    break;
  case MOVE_LEFT:
    draw_pos.x -= movement;
    break;
  case MOVE_DOWN:
    draw_pos.y += movement;
    break;
  case PLANT_BOMB:
  case NONE:
    break;
  }
  return draw_pos;
}

static char will_player_get_dmg(block_t **map, player_t player)
{
  return (map[player.cell_pos.y][player.cell_pos.x] == EXPLOSION);
}

static void draw_lives(player_t player, vector_2d_t player_draw_pos, char will_get_dmg, int anim_val)
{
  vector_2d_t asset_heart_size = {
      .x = ASSET_SPRITE_SIZE / 2,
      .y = ASSET_SPRITE_SIZE / 2};
  vector_2d_t draw_heart_size = {
      .x = (int)((asset_heart_size.x * WINDOW_ZOOM) / 2),
      .y = (int)((asset_heart_size.y * WINDOW_ZOOM) / 2)};
  int draw_total_size = MAX_LIVES * draw_heart_size.x;
  // The offset of the lives to center them above the player
  int draw_offset = (int)((DISPLAY_SPRITE_SIZE - draw_total_size) / 2);
  vector_2d_t hearts = {
      .x = 3 * ASSET_SPRITE_SIZE,
      .y = 0};
  vector_2d_t draw_pos = {
      .y = player_draw_pos.y - draw_heart_size.x,
      .x = player_draw_pos.x + draw_offset};
  for (int i = 0; i < MAX_LIVES; i++)
  {
    if (i < player.lives && !(i + 1 == player.lives && will_get_dmg && anim_val % 3 == 0))
    {
      // Draw Heart
      blit_custom_from_atlas(hearts, asset_heart_size, draw_pos, draw_heart_size);
    }
    else
    {
      // Draw Empty Heart
      blit_custom_from_atlas((vector_2d_t){.x = hearts.x + asset_heart_size.x, .y = hearts.y}, asset_heart_size, draw_pos, draw_heart_size);
    }
    draw_pos.x += draw_heart_size.x;
  }
}

static void draw_players(block_t **map, players_t players, int animation_value, player_action_t pl1_act, player_action_t pl2_act, player_action_t pl3_act, player_action_t pl4_act)
{
  int animation_offset = ((int)(animation_value / 3)) * ASSET_SPRITE_SIZE;
  if (players.player1.lives > 0)
  {
    vector_2d_t player1_draw_pos = get_player_draw_pos(players.player1, pl1_act, animation_value);
    char will_get_dmg = will_player_get_dmg(map, players.player1);
    char will_die = (players.player1.lives == 1 && will_get_dmg);
    draw_lives(players.player1, player1_draw_pos, will_get_dmg, animation_value);
    blit_from_atlas((vector_2d_t){
                        .y = texture_atlas_positions.player1.y,
                        .x = (pl1_act == NONE && !will_die ? 0 : animation_offset) + (will_die ? 12 * ASSET_SPRITE_SIZE : get_player_animation_movement_offset(pl1_act))},
                    player1_draw_pos);
    display_text(players.player1.bot_description.author_name, player1_draw_pos);
  }
  if (players.player2.lives > 0)
  {
    vector_2d_t player2_draw_pos = get_player_draw_pos(players.player2, pl2_act, animation_value);
    char will_get_dmg = will_player_get_dmg(map, players.player2);
    char will_die = (players.player2.lives == 1 && will_get_dmg);
    draw_lives(players.player2, player2_draw_pos, will_get_dmg, animation_value);
    blit_from_atlas((vector_2d_t){
                        .y = texture_atlas_positions.player2.y,
                        .x = (pl2_act == NONE && !will_die ? 0 : animation_offset) + (will_die ? 12 * ASSET_SPRITE_SIZE : get_player_animation_movement_offset(pl2_act))},
                    player2_draw_pos);
    display_text(players.player2.bot_description.author_name, player2_draw_pos);
  }
  if (players.player3.lives > 0)
  {
    vector_2d_t player3_draw_pos = get_player_draw_pos(players.player3, pl3_act, animation_value);
    char will_get_dmg = will_player_get_dmg(map, players.player3);
    char will_die = (players.player3.lives == 1 && will_get_dmg);
    draw_lives(players.player3, player3_draw_pos, will_get_dmg, animation_value);
    blit_from_atlas((vector_2d_t){
                        .y = texture_atlas_positions.player3.y,
                        .x = (pl3_act == NONE && !will_die ? 0 : animation_offset) + (will_die ? 12 * ASSET_SPRITE_SIZE : get_player_animation_movement_offset(pl3_act))},
                    player3_draw_pos);
    display_text(players.player3.bot_description.author_name, player3_draw_pos);
  }
  if (players.player4.lives > 0)
  {
    vector_2d_t player4_draw_pos = get_player_draw_pos(players.player4, pl4_act, animation_value);
    char will_get_dmg = will_player_get_dmg(map, players.player4);
    char will_die = (players.player4.lives == 1 && will_get_dmg);
    draw_lives(players.player4, player4_draw_pos, will_get_dmg, animation_value);
    blit_from_atlas((vector_2d_t){
                        .y = texture_atlas_positions.player4.y,
                        .x = (pl4_act == NONE && !will_die ? 0 : animation_offset) + (will_die ? 12 * ASSET_SPRITE_SIZE : get_player_animation_movement_offset(pl4_act))},
                    get_player_draw_pos(players.player4, pl4_act, animation_value));
    display_text(players.player4.bot_description.author_name, player4_draw_pos);
  }
}

static void draw_map(block_t **map, int animation_value)
{
  {
    for (int row = 0; row < MAP_HEIGHT; row++)
    {
      for (int col = 0; col < MAP_WIDTH; col++)
      {
        vector_2d_t draw_pos = {
            .y = row * DISPLAY_SPRITE_SIZE,
            .x = col * DISPLAY_SPRITE_SIZE,
        };
        blit_from_atlas(texture_atlas_positions.air, draw_pos);
        switch (map[row][col])
        {
        case BOX:
          blit_from_atlas(texture_atlas_positions.box, draw_pos);
          break;
        case AIR:
          break;
        case BOMB1:
        case BOMB2:
        case BOMB3:
        case BOMB4:
        case BOMB5:
        case BOMB6:
        case BOMB7:
        case BOMB8:
        case BOMB9:
        case BOMB10:
          blit_from_atlas((vector_2d_t){.y = texture_atlas_positions.bomb.y, .x = ((int)(animation_value / 3)) * ASSET_SPRITE_SIZE}, draw_pos);
          break;
        case WALL:
          draw_wall(map, (cell_pos_t){.y = row, .x = col}, draw_pos);
          break;
        case EXPLOSION:
          draw_explosion(map, (cell_pos_t){.y = row, .x = col}, draw_pos, animation_value);
          break;
        }
      }
    }
  }
}

// Ive copied this function from this tutorial :/
// https://www.parallelrealities.co.uk/tutorials/shooter/shooter5.php
static void cap_frame_rate(long *then, float *remainder)
{
  long wait, frame_time;
  // 33 cause 1000 / 30 ~ 33 To achieve 30 fps
  wait = 33 + *remainder;
  *remainder -= (int)*remainder;
  frame_time = SDL_GetTicks() - *then;
  wait -= frame_time;
  if (wait < 1)
  {
    wait = 1;
  }
  SDL_Delay(wait);
  *remainder += 0.667;
  *then = SDL_GetTicks();
}

void display_map(block_t **map, players_t players, player_action_t pl1_act, player_action_t pl2_act, player_action_t pl3_act, player_action_t pl4_act)
{
  long then = SDL_GetTicks();
  float remainder = 0;
  for (int i = 0; i < 9; i++)
  {
    cap_frame_rate(&then, &remainder);
    prepare_scene();
    draw_map(map, i);
    draw_players(map, players, i, pl1_act, pl2_act, pl3_act, pl4_act);
    display_markers();
    present_scene();
  }
}
