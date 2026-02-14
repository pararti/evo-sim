# GEMINI.md - EvoSim Project Context

## Project Overview
**EvoSim** is an artificial life simulation where autonomous creatures controlled by Feed-Forward Neural Networks (FFNN) evolve and adapt to their environment. The project demonstrates emergent behavior through natural selection, reproduction, and mutation.

### Core Technologies
- **Backend**: Go (1.25+), Gorilla WebSocket, SQLite, `godotenv`.
- **Frontend**: Vanilla JavaScript, HTML5 Canvas, WebSockets (Binary Protocol).
- **Infrastructure**: Docker, Docker Compose, Caddy.

### Architecture
- **`cmd/app/`**: Application entry point. Handles initialization and the main simulation loop (60 FPS).
- **`internal/world/`**: The core simulation engine managing entities, physics, and world logic.
- **`internal/brain/`**: Neural network logic, including forward propagation and mutations.
- **`internal/entity/`**: Definitions for `Creature` and `Food`.
- **`internal/server/`**: HTTP server for frontend assets and a high-performance binary WebSocket for real-time state broadcasting (approx. 30 FPS).
- **`internal/storage/`**: SQLite implementation for periodic simulation state snapshots.
- **`internal/config/`**: Configuration management via environment variables.

---

## Building and Running

### Prerequisites
- Go 1.25 or higher
- Docker & Docker Compose (optional, for containerized deployment)

### Local Development
1. **Install Dependencies**:
   ```bash
   go mod download
   ```
2. **Configure Environment**:
   Create a `.env` file based on the defaults in `internal/config/config.go`.
3. **Run Application**:
   ```bash
   go run cmd/app/main.go
   ```
4. **Access UI**:
   Open `http://localhost:8080` in your browser.

### Docker Deployment
```bash
docker-compose up -d
```
The simulation will be available on port `8080`, and Caddy reverse proxy on port `8089`.

### Testing
- **TODO**: No unit tests were found in the initial analysis. Add tests using standard `go test ./...`.

---

## Development Conventions

### Coding Style
- Follow standard Go idioms and `gofmt`.
- Logic is strictly separated into packages within `internal/`.
- Use `sync.RWMutex` for thread-safe access to the global world state.

### Communication Protocol
- The WebSocket uses a **custom binary protocol** for performance:
  - **Header**: Uint16 (Creature Count)
  - **Creatures**: Array of { ID (Uint16), X (Float32), Y (Float32) }
  - **Header**: Uint16 (Food Count)
  - **Food**: Array of { X (Float32), Y (Float32) }
- All binary data is encoded in **Little Endian**.

### Configuration
- All simulation parameters (speed, mutation rate, world size) are configurable via environment variables.
- Default values are provided in `internal/config/config.go`.

### State Persistence
- Simulation state is snapshotted to SQLite every 15 minutes.
- Database path is configurable via `DB_PATH`.
