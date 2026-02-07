package main

import (
	"log"
	"time"

	"evo-sim/internal/config"
	"evo-sim/internal/server"
	"evo-sim/internal/storage"
	"evo-sim/internal/world"
)

func main() {
	cfg := config.Load()
	log.Println("Config loaded. World size:", cfg.WorldWidth, "x", cfg.WorldHeight)

	store := storage.NewStorage(cfg.DBPath)

	w := world.NewWorld(cfg)

	srv := server.NewServer(w)
	go srv.Start(cfg.HTTPPort)

	go func() {
		ticker := time.NewTicker(15 * time.Minute)
		for range ticker.C {
			w.Mu.RLock()
			store.SaveSnapshot(w.Creatures, w.Food)
			w.Mu.RUnlock()
		}
	}()

	log.Println("Simulation started...")

	ticker := time.NewTicker(time.Second / 60)
	for range ticker.C {
		w.Update()
	}
}
