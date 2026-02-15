# 🧬 EvoKakSim — Advanced Artificial Life Simulation

EvoKakSim is a high-performance evolutionary simulation written in **Go**. It models an ecosystem where creatures with unique **Genomes** and **Neural Networks** evolve to survive in a dynamic environment with distinct **Biomes**.

Unlike simple genetic algorithms, this project focuses on **Thermodynamics** and **Physical Constraints**: nothing is free, every advantage (speed, size, brain power) has an energy cost.

![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go)
![License](https://img.shields.io/badge/license-MIT-blue)

## Key Features

### Genotype vs Phenotype
Creatures are not born random; they are built from a **Genome** that mutates over generations.
- **Size Gene**: Determines mass and health. Larger creatures are stronger but need more food (Square-Cube Law).
- **Speed Gene**: Determines muscle density. Fast creatures can catch prey but burn energy rapidly.
- **Sense Gene**: Determines visual range. Better vision helps find food but brains consume more energy.
- **Diet Gene**: Determines placement on the Carnivore-Herbivore spectrum.

### Thermodynamics & BMR
The simulation enforces a strict energy budget via **Basal Metabolic Rate (BMR)**.
- **Living Cost**: `BMR = f(Mass, BrainSize, SpeedPotential)`.
- **Movement Cost**: Moving through water or sand requires more work ($W = F \cdot d$).
- **Evolutionary Pressure**: Inefficient creatures (e.g., huge body with small mouth) starve and die out.

### Procedural Terrain & Biomes
The world is generated using **Perlin Noise** and divided into biomes:
1.  **Water**: High movement cost (3x), very slow speed. Safe from non-amphibious predators.
2.  **Sand**: Medium movement penalty.
3.  **Grass**: Normal speed. Food grows here abundantly.

### Neural Network Brain
Each creature makes decisions using a mutable Feed-Forward Neural Network.
- **Inputs**: Vector to food/enemy, terrain type underfoot, internal energy, smell.
- **Outputs**: Velocity vector (X, Y).
- **Neuro-Evolution**: Brain weights mutate along with physical genes.

## Quick Start

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
docker compose up -d --build
```

## Interface

The web interface is optimized for both **Desktop** and **Mobile**:
- **Real-time Visualization**: HTML5 Canvas rendering at 60 FPS using OffscreenCanvas for performance.
- **Responsive HUD**: Adapts layout for small screens.
- **Live Stats**: FPS, Population count, Food abundance.

## Architecture

### Backend (Go)
- **Engine**: Custom physics engine with Spatial Partitioning (Grid) to support thousands of entities on low-end hardware (VPS optimized).
- **Concurrency**: Parallelized update loops and thread-safe data access.
- **Networking**: Binary WebSocket protocol for minimal latency and bandwidth.

### Frontend (Vanilla JS)
- **Rendering**: Optimized 2D Context with off-screen buffering for static terrain.
- **Protocol**: Binary parsing (`ArrayBuffer` / `DataView`) of world state.

## Configuration

Tune the simulation via environment variables or `config.go`:

| Variable | Description |
|----------|-------------|
| `WORLD_WIDTH` | Map width in pixels (e.g., 800) |
| `WORLD_HEIGHT` | Map height in pixels (e.g., 600) |
| `INITIAL_POP` | Starting creature count |
| `MUTATION_RATE` | DNA mutation probability |
| `FOOD_COUNT` | Max food on map |

## License

MIT
