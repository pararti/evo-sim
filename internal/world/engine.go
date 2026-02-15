package world

import (
	"math"
	"math/rand/v2"
	"sync"
	"time"

	"evo-sim/internal/config"
	"evo-sim/internal/entity"
)

type World struct {
	Cfg       *config.Config
	Creatures []*entity.Creature
	Food      []entity.Food
	Grid      *Grid
	Terrain   *TerrainGrid
	SpeciesManager *SpeciesManager
	Mu        sync.RWMutex
	StartTime time.Time
}

func NewWorld(cfg *config.Config) *World {
	w := &World{
		Cfg:            cfg,
		Grid:           NewGrid(cfg.WorldWidth, cfg.WorldHeight, 40.0),
		Terrain:        NewTerrainGrid(cfg.WorldWidth, cfg.WorldHeight, 20.0),
		SpeciesManager: NewSpeciesManager(cfg.SpeciationThreshold),
		StartTime:      time.Now(),
	}

	w.spawnRandomCreatures(cfg.InitialPop)
	for i := 0; i < cfg.FoodCount; i++ {
		w.spawnFood()
	}

	return w
}

func (w *World) spawnRandomCreatures(count int) {
	for i := 0; i < count; i++ {
		// Try to spawn on land
		for attempt := 0; attempt < 5; attempt++ {
			x := rand.Float64() * w.Cfg.WorldWidth
			y := rand.Float64() * w.Cfg.WorldHeight
			if w.Terrain.GetType(x, y) != Water {
				c := entity.NewCreature(
					rand.IntN(10000000),
					x, y,
					w.Cfg.InputSize,
					w.Cfg.HiddenSize,
					w.Cfg.OutputSize,
				)
				c.SpeciesID = w.SpeciesManager.Classify(c.Genome)
				w.Creatures = append(w.Creatures, c)
				break
			}
		}
	}
}

func (w *World) spawnFood() {
	var x, y float64
	
	// Try to find a good spot (Grass preferred)
	for i := 0; i < 10; i++ {
		// 70% chance to spawn in the "Oasis" (center 40% of the world)
		if rand.Float64() < 0.7 {
			marginW := w.Cfg.WorldWidth * 0.3
			marginH := w.Cfg.WorldHeight * 0.3
			x = marginW + rand.Float64()*(w.Cfg.WorldWidth*0.4)
			y = marginH + rand.Float64()*(w.Cfg.WorldHeight*0.4)
		} else {
			x = rand.Float64() * w.Cfg.WorldWidth
			y = rand.Float64() * w.Cfg.WorldHeight
		}

		// Food grows best on Grass, okay on Sand, never on Water
		tile := w.Terrain.GetType(x, y)
		if tile == Grass {
			break // Good spot
		}
		if tile == Sand && rand.Float64() < 0.2 {
			break // Rare cactus?
		}
		// If Water, retry
	}

	w.Food = append(w.Food, entity.Food{
		ID: rand.IntN(10000000),
		X:  x,
		Y:  y,
	})
}

func (w *World) Update() {
	w.Mu.Lock()
	defer w.Mu.Unlock()

	// 1. Rebuild grid
	w.Grid.Clear()
	for _, c := range w.Creatures {
		w.Grid.InsertCreature(c)
	}
	for _, f := range w.Food {
		w.Grid.InsertFood(f)
	}

	var newChildren []*entity.Creature
	deadCreatures := make(map[int]bool)
	eatenFood := make(map[int]bool)
	matedThisTick := make(map[int]bool)

	// 2. Main Simulation Loop
	for _, c := range w.Creatures {
		if deadCreatures[c.ID] {
			continue
		}

		// Find targets
		foodX, foodY, foodDist, foodID := w.findNearestFood(c, eatenFood)
		targetX, targetY, targetDist, targetID, isTargetCarnivore := w.findNearestCreature(c, deadCreatures)

		// Update Brain
		roleVal := -1.0
		if isTargetCarnivore { roleVal = 1.0 }

		// Get Terrain Physics
		speedFactor, energyCostFactor := w.Terrain.GetMovementPenalty(c.X, c.Y)
		
		c.Update(foodX, foodY, targetX, targetY, roleVal, speedFactor, energyCostFactor, w.Cfg.WorldWidth, w.Cfg.WorldHeight, w.Cfg.MaxAge)

		// Boundaries
		if c.X < 0 { c.X = 0 } else if c.X > w.Cfg.WorldWidth { c.X = w.Cfg.WorldWidth }
		if c.Y < 0 { c.Y = 0 } else if c.Y > w.Cfg.WorldHeight { c.Y = w.Cfg.WorldHeight }

		// Interactions
		if !c.IsCarnivore && foodID != -1 && foodDist < w.Cfg.EatRadius*c.Size {
			if !eatenFood[foodID] {
				c.Energy += w.Cfg.FoodEnergy
				eatenFood[foodID] = true
			}
		}

		if c.IsCarnivore && targetID != -1 && targetDist < w.Cfg.EatRadius*c.Size {
			if !deadCreatures[targetID] {
				target := w.getCreatureByID(targetID)
				if target != nil && target.Size < c.Size*1.2 {
					c.Energy += target.Energy * 0.8
					deadCreatures[targetID] = true
					w.SpeciesManager.RemoveCreature(target.SpeciesID)
				}
			}
		}

		// Reproduction
		if c.Energy > c.ReproductionThreshold && !matedThisTick[c.ID] {
			mate := w.findMate(c, deadCreatures, matedThisTick)
			var child *entity.Creature
			if mate != nil {
				child = c.ReproduceSexual(mate, w.Cfg.MutationRate, w.Cfg.MutationStrength)
				matedThisTick[mate.ID] = true
			} else if c.Energy > c.ReproductionThreshold*w.Cfg.AsexualThresholdMult {
				child = c.ReproduceAsexual(w.Cfg.MutationRate, w.Cfg.MutationStrength)
			}
			if child != nil {
				child.ID = rand.IntN(10000000)
				child.SpeciesID = w.SpeciesManager.Classify(child.Genome)
				newChildren = append(newChildren, child)
				matedThisTick[c.ID] = true
			}
		}

		if c.Energy <= 0 {
			deadCreatures[c.ID] = true
			w.SpeciesManager.RemoveCreature(c.SpeciesID)
		}
	}

	// 3. Cleanup & Finalize
	// Remove dead creatures
	newCreatureList := make([]*entity.Creature, 0, len(w.Creatures))
	for _, c := range w.Creatures {
		if !deadCreatures[c.ID] {
			newCreatureList = append(newCreatureList, c)
		}
	}
	w.Creatures = append(newCreatureList, newChildren...)

	// Remove eaten food and respawn
	newFoodList := make([]entity.Food, 0, len(w.Food))
	for _, f := range w.Food {
		if !eatenFood[f.ID] {
			newFoodList = append(newFoodList, f)
		}
	}
	w.Food = newFoodList
	for len(w.Food) < w.Cfg.FoodCount {
		w.spawnFood()
	}

	// 4. Rescue population
	if len(w.Creatures) < 10 {
		w.spawnRandomCreatures(5)
	}
}

func (w *World) getCreatureByID(id int) *entity.Creature {
	for _, c := range w.Creatures {
		if c.ID == id {
			return c
		}
	}
	return nil
}

func (w *World) findNearestCreature(c *entity.Creature, dead map[int]bool) (float64, float64, float64, int, bool) {
	minDist := math.MaxFloat64
	var nx, ny float64
	var targetID = -1
	var isCarnivore bool

	// Use creature's view radius, but clamp it to reasonable grid lookup limits if needed
	w.Grid.ForEachNeighbor(c.X, c.Y, c.ViewRadius, func(other *entity.Creature) {
		if other.ID == c.ID || dead[other.ID] {
			return
		}
		dist := math.Hypot(other.X-c.X, other.Y-c.Y)
		// Only see things within ViewRadius
		if dist < c.ViewRadius && dist < minDist {
			minDist, nx, ny, targetID, isCarnivore = dist, other.X, other.Y, other.ID, other.IsCarnivore
		}
	}, nil)

	return nx, ny, minDist, targetID, isCarnivore
}

func (w *World) findNearestFood(c *entity.Creature, eaten map[int]bool) (float64, float64, float64, int) {
	minDist := math.MaxFloat64
	var nx, ny float64
	var fid = -1

	// 1. Spatial
	w.Grid.ForEachNeighbor(c.X, c.Y, c.ViewRadius, nil, func(f entity.Food) {
		if eaten[f.ID] {
			return
		}
		dist := math.Hypot(f.X-c.X, f.Y-c.Y)
		if dist < c.ViewRadius && dist < minDist {
			minDist, nx, ny, fid = dist, f.X, f.Y, f.ID
		}
	})

	// 2. Global Fallback (only if no spatial result found? Or remove global fallback to enforce blindness?)
	// To make evolution real, if they can't see it, they can't find it.
	// But to prevent total extinction of "blind" early gens, maybe keep a small "smell" range?
	// For now, let's remove the global fallback to strictly enforce ViewRadius.
	// This makes "SenseGene" actually valuable.
	
	return nx, ny, minDist, fid
}

func (w *World) findMate(c *entity.Creature, dead, mated map[int]bool) *entity.Creature {
	var best *entity.Creature
	bestDist := math.MaxFloat64

	w.Grid.ForEachNeighbor(c.X, c.Y, c.ViewRadius, func(other *entity.Creature) {
		if other.ID == c.ID || dead[other.ID] || mated[other.ID] {
			return
		}
		if other.Energy <= other.ReproductionThreshold {
			return
		}
		if other.IsCarnivore != c.IsCarnivore {
			return
		}
		if c.Genome.Distance(other.Genome) > w.Cfg.MatingDistanceThreshold {
			return
		}
		dist := math.Hypot(other.X-c.X, other.Y-c.Y)
		if dist < c.ViewRadius && dist < w.Cfg.EatRadius*c.Size && dist < bestDist {
			bestDist = dist
			best = other
		}
	}, nil)

	return best
}
