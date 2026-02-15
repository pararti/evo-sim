package server

import (
	"encoding/json"
	"net/http"

	"evo-sim/internal/world"
)

type Server struct {
	World *world.World
}

func NewServer(w *world.World) *Server {
	return &Server{World: w}
}

func (s *Server) Start(port string) error {
	http.Handle("/", http.FileServer(http.Dir("./web")))
	http.HandleFunc("/ws", s.handleWebSocket)
	http.HandleFunc("/api/map", s.handleMap)

	return http.ListenAndServe(":"+port, nil)
}

func (s *Server) handleMap(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// CORS for dev (if needed, otherwise can remove)
	w.Header().Set("Access-Control-Allow-Origin", "*")

	s.World.Mu.RLock()
	defer s.World.Mu.RUnlock()

	json.NewEncoder(w).Encode(s.World.Terrain)
}
