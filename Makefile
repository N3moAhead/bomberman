# Definiere Variablen für Quelldateien und Objektdateien
SRC_DIR := ./src
SRC_FILES := $(wildcard $(SRC_DIR)/*.c)
OBJ_FILES := $(patsubst $(SRC_DIR)/%.c, $(SRC_DIR)/%.o, $(SRC_FILES))

# Verwendung des gcc-Compilers als Standard
CC := gcc
# Hinzufügen von -g für Debugging-Informationen, -Wall und -Wextra für Warnungen
CFLAGS := -Wall -Wextra -g

# Bestimme das Betriebssystem
ifeq ($(OS),Windows_NT)
	# Wenn Windows verwendet wird
	EXECUTABLE := game.exe
	# Ersetze 'rm' durch 'del' für Windows
	RM := del /Q
	# Ersetze '/' durch '\' in den Dateipfaden
	FIXPATH = $(subst /,\,$1)
else
	# Wenn ein Unix-basiertes System verwendet wird (Linux, macOS usw.)
	EXECUTABLE := game
	# Verwende 'rm' zum Löschen von Dateien
	RM := rm -f
	# Das Dateisystem verwendet bereits Schrägstriche, daher kein Bedarf für die Fixpath-Funktion
	FIXPATH = $1
endif

# Phony-Ziele (Ziele, die nicht tatsächlich Dateien darstellen)
.PHONY: all clear

# Standardziel ist das Spiel
all: $(EXECUTABLE)

# Regel zum Erstellen der ausführbaren Datei
$(EXECUTABLE): $(OBJ_FILES)
	$(CC) $(OBJ_FILES) -o $(call FIXPATH,$(EXECUTABLE))

# Regel zum Erstellen von Objektdateien
$(SRC_DIR)/%.o: $(SRC_DIR)/%.c
	$(CC) -c $(CFLAGS) $< -o $@

# Regel zum Löschen von Objektdateien und der ausführbaren Datei
clear:
	$(RM) $(call FIXPATH,$(OBJ_FILES)) $(call FIXPATH,$(EXECUTABLE))

# Build- und Ausführungsziele
build: $(EXECUTABLE)
run: $(EXECUTABLE)
	$(call FIXPATH,./$(EXECUTABLE))
