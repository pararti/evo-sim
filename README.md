# 🧬 EvoSim — Neural Network Evolution Simulation

An artificial life simulation where creatures with neural network brains evolve, find food, reproduce, and adapt to their environment through natural selection.

![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go)
![License](https://img.shields.io/badge/license-MIT-blue)

## 🎥 Demo

The web interface shows the simulation in real-time: creatures (circles) move toward the nearest food (squares), consume energy, and reproduce when they accumulate enough reserves.

## 🚀 Quick Start

### Local Run

```bash
# Install dependencies
go mod download

# Run the simulation
go run cmd/app/main.go
```

Open http://localhost:8080 in your browser.

### Docker

```bash
# Run with Docker Compose
docker-compose up -d

# Simulation will be available on port 8080
# Caddy reverse proxy — on port 8089
```

## 🧠 Architecture

### Creature Neural Network

Each creature is controlled by a **Feed-Forward Neural Network (FFNN)** with three layers:

| Layer | Description |
|-------|-------------|
| **Input** | 3 neurons: vector to nearest food (dx, dy) + energy level |
| **Hidden** | 4 neurons with Tanh activation function |
| **Output** | 2 neurons: velocity X and Y |

Behavior is completely emergent — no hardcoded logic like "if hungry, go to food".

### Evolution

- **Reproduction**: When energy > 150, the creature splits in half
- **Inheritance**: Offspring receives a copy of the parent's weights
- **Mutation**: 10% of weights are changed with random noise
- **Selection**: Unlucky creatures die of starvation

## 📁 Project Structure

```
.
├── cmd/app/           # Entry point
├── internal/
│   ├── brain/         # Neural network (feed-forward, mutations)
│   ├── entity/        # Creatures and food
│   ├── world/         # Game engine and physics
│   ├── server/        # HTTP + WebSocket server
│   ├── storage/       # SQLite for state persistence
│   └── config/        # Configuration
├── web/               # Frontend (HTML5 Canvas + WebSocket)
└── docker-compose.yml
```

## ⚙️ Configuration

Via environment variables (`.env`):

```env
WORLD_WIDTH=800         # World width
WORLD_HEIGHT=600        # World height
INITIAL_POP=20          # Initial creature count
FOOD_COUNT=50           # Initial food count
FOOD_ENERGY=50.0        # Energy gained from food
MOVE_COST=0.1           # Energy cost per movement
SPEED_FACTOR=2.0        # Movement speed multiplier
MUTATION_RATE=0.1       # Probability of weight mutation
MUTATION_STRENGTH=0.2   # Strength of mutation noise
REPRODUCE_THRESHOLD=150 # Energy required to reproduce
HTTP_PORT=8080          # Server port
DB_PATH=./database.db   # SQLite path
```

## 🛠️ Tech Stack

- **Backend**: Go 1.25+, Gorilla WebSocket
- **Frontend**: Vanilla JS, HTML5 Canvas
- **Database**: SQLite
- **Infrastructure**: Docker, Caddy

## 📝 API

### WebSocket

```
ws://localhost:8080/ws
```

The server sends world state 60 times per second:

```json
{
  "creatures": [
    {"id": 1, "x": 100, "y": 200, "energy": 80}
  ],
  "food": [
    {"id": 1, "x": 150, "y": 250}
  ]
}
```

## 💾 State Snapshots

Simulation state is automatically saved to SQLite every 15 minutes.

## 📄 License

MIT

---

**Made with ❤️ in Go**
