package main

import (
	"github.com/N3moAhead/bombahead/website/internal/cfg"
	"github.com/N3moAhead/bombahead/website/internal/db"
	"github.com/N3moAhead/bombahead/website/internal/router"
)

func main() {
	cfg := cfg.Load()
	db.Init(cfg)
	router.Start(cfg)
}
