SRC_DIRS := ./src ./src/bots ./src/util
SRC_FILES := $(foreach dir,$(SRC_DIRS),$(wildcard $(dir)/*.c))
OBJ_FILES := $(patsubst %.c,%.o,$(SRC_FILES))

CC := gcc
CFLAGS := -Wall -Wextra -lSDL2 -lSDL2_image -g
LDFLAGS := -lSDL2 -lSDL2_image

ifeq ($(OS),Windows_NT)
	EXECUTABLE := game.exe
	RM := del /Q
	FIXPATH = $(subst /,\,$1)
else
	EXECUTABLE := game
	RM := rm -f
	FIXPATH = $1
endif

.PHONY: all clear

all: $(EXECUTABLE)

$(EXECUTABLE): $(OBJ_FILES)
	$(CC) $(OBJ_FILES) -o $(call FIXPATH,$(EXECUTABLE)) $(LDFLAGS)

%.o: %.c
	$(CC) -c $(CFLAGS) $< -o $@

clear:
	$(RM) $(call FIXPATH,$(OBJ_FILES)) $(call FIXPATH,$(EXECUTABLE))

build: $(EXECUTABLE)
run: $(EXECUTABLE)
	$(call FIXPATH,./$(EXECUTABLE))
