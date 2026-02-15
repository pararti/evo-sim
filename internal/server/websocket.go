package server

import (
	"encoding/binary"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Разрешаем CORS для локальной разработки (в проде можно ужесточить)
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	log.Println("New client connected via WebSockets")

	ticker := time.NewTicker(33 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		s.World.Mu.RLock()

		creaturesCount := len(s.World.Creatures)
		foodCount := len(s.World.Food)

		// (2 bytes header) + (N * 18 bytes) + (2 bytes header) + (M * 8 bytes)
		packetSize := 2 + (creaturesCount * 18) + 2 + (foodCount * 8)
		buf := make([]byte, packetSize)
		offset := 0

		// === CREATURE SECTION ===
		binary.LittleEndian.PutUint16(buf[offset:], uint16(creaturesCount))
		offset += 2

		for _, c := range s.World.Creatures {
			// ID
			binary.LittleEndian.PutUint16(buf[offset:], uint16(c.ID))
			offset += 2
			// X, Y
			binary.LittleEndian.PutUint32(buf[offset:], math.Float32bits(float32(c.X)))
			offset += 4
			binary.LittleEndian.PutUint32(buf[offset:], math.Float32bits(float32(c.Y)))
			offset += 4
			// IsCarnivore (1 byte)
			if c.IsCarnivore {
				buf[offset] = 1
			} else {
				buf[offset] = 0
			}
			offset += 1
			// Size (4 bytes)
			binary.LittleEndian.PutUint32(buf[offset:], math.Float32bits(float32(c.Size)))
			offset += 4
			// Color (3 bytes)
			buf[offset] = uint8(c.Genome.ColorR * 255)
			offset++
			buf[offset] = uint8(c.Genome.ColorG * 255)
			offset++
			buf[offset] = uint8(c.Genome.ColorB * 255)
			offset++
		}

		// === FOOD SECTION ===
		binary.LittleEndian.PutUint16(buf[offset:], uint16(foodCount))
		offset += 2

		for _, f := range s.World.Food {
			// X
			binary.LittleEndian.PutUint32(buf[offset:], math.Float32bits(float32(f.X)))
			offset += 4
			// Y
			binary.LittleEndian.PutUint32(buf[offset:], math.Float32bits(float32(f.Y)))
			offset += 4
		}

		s.World.Mu.RUnlock()

		if err := conn.WriteMessage(websocket.BinaryMessage, buf); err != nil {
			log.Printf("Client disconnected: %v", err)
			break
		}
	}
}
