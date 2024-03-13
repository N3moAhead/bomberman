SRC_DIR := ./src
SRC_FILES := $(wildcard $(SRC_DIR)/*.c)
OBJ_FILES := $(patsubst $(SRC_DIR)/%.c, $(SRC_DIR)/%.o, $(SRC_FILES))

CC := gcc
CFLAGS := -Wall -Wextra -g

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
	$(CC) $(OBJ_FILES) -o $(call FIXPATH,$(EXECUTABLE))

$(SRC_DIR)/%.o: $(SRC_DIR)/%.c
	$(CC) -c $(CFLAGS) $< -o $@

clear:
	$(RM) $(call FIXPATH,$(OBJ_FILES)) $(call FIXPATH,$(EXECUTABLE))

build: $(EXECUTABLE)
run: $(EXECUTABLE)
	$(call FIXPATH,./$(EXECUTABLE))
