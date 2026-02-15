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

	response := struct {
		Terrain   *world.TerrainGrid `json:"terrain"`
		StartTime int64              `json:"startTime"` // Unix timestamp in milliseconds
	}{
		Terrain:   s.World.Terrain,
		StartTime: s.World.StartTime.UnixMilli(),
	}

	json.NewEncoder(w).Encode(response)
}
