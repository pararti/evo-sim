package storage

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"evo-sim/internal/entity"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	DB *sql.DB
}

type WorldSnapshot struct {
	Timestamp int64              `json:"timestamp"`
	Stats     map[string]int     `json:"stats"` // Например: кол-во живых
	Creatures []*entity.Creature `json:"creatures"`
	Food      []entity.Food      `json:"food"`
}

func NewStorage(dbPath string) *Storage {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal("Failed to open DB:", err)
	}

	query := `
	CREATE TABLE IF NOT EXISTS snapshots (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		data JSON
	);`

	if _, err := db.Exec(query); err != nil {
		log.Fatal("Failed to create table:", err)
	}

	return &Storage{DB: db}
}

func (s *Storage) SaveSnapshot(creatures []*entity.Creature, food []entity.Food) {
	snapshot := WorldSnapshot{
		Timestamp: time.Now().Unix(),
		Stats: map[string]int{
			"creatures_count": len(creatures),
			"food_count":      len(food),
		},
		Creatures: creatures,
		Food:      food,
	}

	data, err := json.Marshal(snapshot)
	if err != nil {
		log.Println("Error marshalling snapshot:", err)
		return
	}

	_, err = s.DB.Exec("INSERT INTO snapshots (data) VALUES (?)", data)
	if err != nil {
		log.Println("Error saving to DB:", err)
	} else {
		log.Printf("Snapshot saved. Size: %d bytes", len(data))
	}
}
