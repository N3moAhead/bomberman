```
    ____                  __
   / __ )____  ____ ___  / /_  ___  _________ ___  ____ _____
  / __  / __ \/ __ `__ \/ __ \/ _ \/ ___/ __ `__ \/ __ `/ __ \
 / /_/ / /_/ / / / / / / /_/ /  __/ /  / / / / / / /_/ / / / /
/_____/\____/_/ /_/ /_/_.___/\___/_/  /_/ /_/ /_/\__,_/_/ /_/
```

This repository contains everything needed to play the classic game Bomberman on the terminal. It can easily be extended by developing bots for the four possible players. If you want, you can fork the repository or simply push the written bot here. The documentation below explains how to write a bot for the game.

# Installation

Clone the repository:

```bash
# HTTPS
git clone https://github.com/N3moAhead/bomberman.git
# SSH
git clone git@github.com:N3moAhead/bomberman.git
```

I am using make and gcc, so please install them if you haven't already.
After that you can just run the following commands:

**Build the binary**

```bash
make
```

**Build and run**

```bash
make run
```

**Clear all build files**

```bash
make clear
```

# Docs

The principle is simple. In each round, four player functions are called up one after the other to control one of the four bots. Each function receives the current playing field, the position of all players and the current round number. Each function must return one of six actions performed by the respective bot.

The playing field is a two-dimensional array. It is filled with different types of fields, such as walls, bombs or explosions.

The player positions are passed as a struct, which contains a further struct for each player. Each player struct contains the position of the current player and the current lives.

Attention! The playing field array does not contain any players. These are passed separately to prevent them from covering fields. For example, a player could stand on a bomb and hide it from others.

Finally, the player function is passed an int containing the current turn number. This simply counts up each round by 1.

The individual actions are briefly explained below, but if you want to get started right away, you can simply take a look at the file `types.h` and if you can make sense of it, you can get started right away.

## Player Actions

Each player function can return one of six possible actions in each round to control one of the bots. There are four actions to move the bot. One action to place a bomb on the current field and one function to simply chill for a round. Actually quite simple :)

- MOVE_UP: Moves the player a field up
- MOVE_DOWN: Moves the player a field down
- MOVE_LEFT: Moves the plyer a field to the left
- MOVE_RIGHT: Moves the player field to the right
- PLANT_BOMB: Places a bomb on the current field of the player
- NONE: Do nothing

## The Map

The playing field is a two-dimensional array that represents the current playing field. Each value in the array is a block_t enum. Each block can be one of the following:

- BOMB1: A bomb that has been freshly placed
- BOMB2: A bomb thats one round old
- BOMB3: A bomb that will explode in the next round
- WALL: A wall a player is unable to go through
- EXPLOSION: An exploding field, that harms players
- AIR: Nothing

The playing field is updated after each round. This mainly affects the bombs. A newly placed bomb is called BOMB1. This counts up one each round until BOMB3. In the next round, the bomb (BOMB3) explodes in a cross shape. An explosion lasts one round. The explosion fields then become air fields.

The dimensions of the map are defined in the file `src/constants.h` as `MAP_HEIGHT` and `MAP_WIDTH`. If you would like to iterate through the whole map the code could look like this:

```c
#include "constants.h"

void some_function(block_t **map) {
  for (int row = 0; row < MAP_HEIGHT; row++) {
    for (int col = 0; col < MAP_WIDTH; col++) {
      block_t current_block = map[row][col];
      // Do something...
    }
  }
}
```

## Example of checking if a field is a wall

```c
if (map[3][3] == WALL) {
  printf("Is a wall");
} else {
  printf("This field must be something else");
}
```

## Example of checking if a player stands on a specific field

```c
if (players->player2.cell_pos.x == 3 && players->player2.cell_pos.y == 4) {
  printf("Player2 is standing on the field 3/4");
} else {
  printf("Player2 ist standing somewhere else");
}
```

## Example player function

We write a bot for player 1. The bot should always repeat the following sequence: First go one step down, then one step to the right and then place a bomb.

The code for this could look like this:

```c
// src/player1.c
#include "player1.h"
#include "types.h"

// This function gets called on each tick
player_action_t get_player_1_action(block_t **map, players_t *players, int game_round) {
  player_action_t actions[3] = {MOVE_DOWN, MOVE_RIGHT, PLANT_BOMB};
  return actions[game_round % 3];
}
```

First, we define an array with our sequence of actions that we want to execute one after the other. Then we simply return the next action in the sequence for each call. To do this, we select the action that corresponds to game round modulo 3. For example, in game round 6, this is MOVE_DOWN, as 6 % 3 = 0.

# Rules

Not that many in general do what you like hahaa

- It is not permitted to adapt code outside the player functions for your own benefit.
- You may use malloc, but after each call of the player function no more memory may be allocated and must have been released.
- Have fun :)
