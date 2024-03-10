# Definiere Variablen für Quelldateien und Objektdateien
SRC_DIR := ./src
SRC_FILES := $(wildcard $(SRC_DIR)/*.c)
OBJ_FILES := $(patsubst $(SRC_DIR)/%.c, $(SRC_DIR)/%.o, $(SRC_FILES))

CC := gcc
CFLAGS := -Wall -Wextra -g

game: $(OBJ_FILES)
	$(CC) $(OBJ_FILES) -o game

%.o: %.c
	$(CC) -c $(CFLAGS) $< -o $@

clear:
	rm -f $(OBJ_FILES) game

build: game
run: game
	./game
