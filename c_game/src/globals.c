#include <SDL2/SDL.h>
#include <SDL2/SDL_ttf.h>
#include "types.h"
#include "util/sdl_helper.h"
#include "constants.h"

SDL_Window *window = NULL;
SDL_Renderer *renderer = NULL;
SDL_Texture *texture_atlas = NULL;
TTF_Font *font = NULL;
const atlas_positions_t texture_atlas_positions = {
  .player1 = {
    .x = 0,
    .y = 96
  }, 
  .player2 = {
    .x = 0,
    .y = 128
  }, 
  .player3 = {
    .x = 0,
    .y = 160
  }, 
  .player4 = {
    .x = 0,
    .y = 192
  }, 
  .bomb = {
    .x = 0,
    .y = 64
  } , 
  .wall = {
    .x = 0,
    .y = 32
  }, 
  .explosion = {
    .x = 0,
    .y = 224
  } , 
  .air = {
    .x = 64,
    .y = 0
  },
  .box = {
    .y = 0,
    .x = 32,
  }
};

marker_node_t *ll_marker = NULL;

// Function to create a new node
marker_node_t *create_node(SDL_Color color, vector_2d_t pos, char *text)
{
  marker_node_t *new_node = (marker_node_t *)malloc(sizeof(marker_node_t));
  if (new_node == NULL)
  {
    fprintf(stderr, "Memory allocation failed in function create_node in globals.c\n");
    exit(EXIT_FAILURE);
  }
  new_node->color = color;
  new_node->text = text;
  new_node->pos = pos;
  new_node->next = NULL;
  return new_node;
}

// Function to print the elements of the list
void print_markers()
{
  marker_node_t *temp = ll_marker;
  while (temp != NULL)
  {
    printf("%s; ", temp->text);
    printf("R: %d, G: %d, B: %d, A: %d; ", temp->color.r, temp->color.g, temp->color.b, temp->color.a);
    printf("Row: %d, Col: %d\n", temp->pos.y, temp->pos.x);
    temp = temp->next;
  }
  printf("\n");
}

void display_markers()
{
  marker_node_t *temp = ll_marker;
  while (temp != NULL)
  {
    display_marker(temp->color, temp->pos, temp->text);
    temp = temp->next;
  }
}

// Function to delete the entire list
void delete_markers()
{
  marker_node_t *current = ll_marker;
  marker_node_t *next;
  while (current != NULL)
  {
    next = current->next;
    free(current);
    current = next;
  }
  ll_marker = NULL;
}

// Function to insert a new node at the beginning of the list
void add_marker(SDL_Color color, cell_pos_t pos, char *text)
{
  marker_node_t *newNode = create_node(
    color,
    VECTOR_2D(pos.y * DISPLAY_SPRITE_SIZE, pos.x * DISPLAY_SPRITE_SIZE),
    text
  );
  newNode->next = ll_marker;
  ll_marker = newNode;
}