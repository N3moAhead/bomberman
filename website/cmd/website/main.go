package main

import (
	"context"
	"os"

	"github.com/N3moAhead/bomberman/website/internal/templates/home"
)

func main() {
	h := home.Home("Lukas")
	h.Render(context.Background(), os.Stdout)
}
