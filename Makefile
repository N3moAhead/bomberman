output: ./src/main.o
	gcc ./src/main.o -o game

main.o: ./src/main.c
	gcc -c ./src/main.c

clear:
	rm ./src/*.o game