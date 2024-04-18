// Map Settings
#define MAP_HEIGHT 11 // should be odd
#define MAP_WIDTH 11 // should be odd

// Player Settings
#define MAX_LIVES 3

/** SDL SETTINGS **/
// Window size zoom
#define WINDOW_ZOOM 2
// The size of the source sprites
#define ASSET_SPRITE_SIZE 32
// The size of the sprites to print to the screen
#define DISPLAY_SPRITE_SIZE (ASSET_SPRITE_SIZE * WINDOW_ZOOM)
#define WINDOW_WIDTH (MAP_WIDTH * DISPLAY_SPRITE_SIZE)
#define WINDOW_HEIGHT (MAP_HEIGHT * DISPLAY_SPRITE_SIZE)