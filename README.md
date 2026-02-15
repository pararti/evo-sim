# 🧬 EvoKakSim — Advanced Artificial Life Simulation

[GitHub Repository](https://github.com/pararti/evo-sim)

EvoKakSim is a high-performance evolutionary simulation written in **Go**. It models an ecosystem where creatures with unique **Genomes** and **Neural Networks** evolve to survive in a dynamic environment with distinct **Biomes**.

Unlike simple genetic algorithms, this project focuses on **Thermodynamics** and **Physical Constraints**: nothing is free, every advantage (speed, size, brain power) has an energy cost.

![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go)
![License](https://img.shields.io/badge/license-MIT-blue)

## Key Features

### 🧬 Advanced Genetics & Epistasis
Creatures possess a complex **Genome** where traits are not 1:1 with genes but interact via **Epistasis**:
- **Constitution Gene**: Affects bone density and health. High constitution increases max energy but significantly increases mass (slowing movement).
- **Metabolism Gene**: Determines how fast energy is converted to work. High metabolism boosts speed but burns calories rapidly (Red Queen hypothesis).
- **Fertility Gene**: Controls reproductive strategy (r/K selection). High fertility allows earlier reproduction but produces weaker offspring.
- **Physical Traits**: Size, Speed, Sense, Diet, and Color (for lineage visualization).

### 🔬 Speciation & Phylogeny
The simulation tracks evolutionary divergence in real-time.
- **Genetic Distance**: Species are defined by clustering genomes based on Euclidean distance in high-dimensional gene space.
- **Reproductive Isolation**: Creatures can only reproduce with genetically similar mates, leading to distinct species branches.
- **Visual Cladogram**: Lineages evolve distinct color patterns, making speciation visible on the map.

### 📉 Thermodynamics & Gradient Aging
Energy is the fundamental currency.
- **Basal Metabolic Rate (BMR)**: `BMR = f(Mass, BrainComplexity, SpeedPotential)`.
- **Gradient Aging**: Instead of a sudden death at `MaxAge`, creatures experience **Senescence**. Energy efficiency drops quadratically with age ($Cost \propto Age^2$), forcing older creatures to eat more or die.
- **Physics**: Movement cost is strictly $W = F \cdot d$. Moving through water or sand incurs heavy penalties.

### 🧠 Neural Network & Crossover
Each creature is controlled by a Feed-Forward Neural Network that evolves over time.
- **Inputs**: Vector to nearest food/enemy, terrain type underfoot, internal energy levels, and pheromones/smell.
- **Outputs**: Velocity vector (X, Y) driving movement.
- **Sexual Reproduction**: Offspring inherit a mix of brain weights and biases from two parents (Crossover) plus random mutations.
- **Asexual Reproduction**: Cloning with mutation for rapid colonization.

### Procedural Terrain & Biomes
The world is generated using **Perlin Noise** and divided into biomes:
1.  **Water**: High movement cost (3x), very slow speed. Safe from non-amphibious predators.
2.  **Sand**: Medium movement penalty.
3.  **Grass**: Normal speed. Food grows here abundantly.


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
