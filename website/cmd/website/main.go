package main

import (
	"github.com/N3moAhead/bomberman/website/internal/cfg"
	"github.com/N3moAhead/bomberman/website/internal/db"
	"github.com/N3moAhead/bomberman/website/internal/router"
)

func main() {
	cfg := cfg.Load()
	db.Init(cfg)
	router.Start(cfg)
}
