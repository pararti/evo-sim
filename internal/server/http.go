package server

import (
	"log"
	"net/http"

	"evo-sim/internal/world"
)

type Server struct {
	World *world.World
}

func NewServer(w *world.World) *Server {
	return &Server{World: w}
}

func (s *Server) Start(port string) {
	fs := http.FileServer(http.Dir("./web"))
	http.Handle("/", fs)

	http.HandleFunc("/ws", s.handleWebSocket)

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
