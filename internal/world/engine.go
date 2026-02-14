package world

import (
	"math"
	"math/rand/v2"
	"sync"

	"evo-sim/internal/config"
	"evo-sim/internal/entity"
)

type World struct {
	Cfg       *config.Config
	Creatures []*entity.Creature
	Food      []entity.Food
	Grid      *Grid
	Mu        sync.RWMutex
}

func NewWorld(cfg *config.Config) *World {
	w := &World{
		Cfg:  cfg,
		Grid: NewGrid(cfg.WorldWidth, cfg.WorldHeight, 40.0),
	}

	w.spawnRandomCreatures(cfg.InitialPop)
	for i := 0; i < cfg.FoodCount; i++ {
		w.spawnFood()
	}

	return w
}

func (w *World) spawnRandomCreatures(count int) {
	for i := 0; i < count; i++ {
		w.Creatures = append(w.Creatures, entity.NewCreature(
			rand.IntN(10000000),
			rand.Float64()*w.Cfg.WorldWidth,
			rand.Float64()*w.Cfg.WorldHeight,
			w.Cfg.InputSize,
			w.Cfg.HiddenSize,
			w.Cfg.OutputSize,
		))
	}
}

func (w *World) spawnFood() {
	w.Food = append(w.Food, entity.Food{
		ID: rand.IntN(10000000),
		X:  rand.Float64() * w.Cfg.WorldWidth,
		Y:  rand.Float64() * w.Cfg.WorldHeight,
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
		c.Update(foodX, foodY, targetX, targetY, roleVal, w.Cfg.SpeedFactor, w.Cfg.MoveCost, w.Cfg.WorldWidth, w.Cfg.WorldHeight)

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
				}
			}
		}

		// Life cycle
		if c.Energy > w.Cfg.ReproduceThreshold {
			child := c.Reproduce(w.Cfg.MutationRate, w.Cfg.MutationStrength)
			child.ID = rand.IntN(10000000)
			newChildren = append(newChildren, child)
		}

		if c.Energy <= 0 || c.Age > 10000 {
			deadCreatures[c.ID] = true
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

	cells, _ := w.Grid.GetNeighbors(c.X, c.Y, 300.0)
	for _, cell := range cells {
		for _, other := range cell {
			if other.ID == c.ID || dead[other.ID] {
				continue
			}
			dist := math.Hypot(other.X-c.X, other.Y-c.Y)
			if dist < minDist {
				minDist, nx, ny, targetID, isCarnivore = dist, other.X, other.Y, other.ID, other.IsCarnivore
			}
		}
	}
	return nx, ny, minDist, targetID, isCarnivore
}

func (w *World) findNearestFood(c *entity.Creature, eaten map[int]bool) (float64, float64, float64, int) {
	minDist := math.MaxFloat64
	var nx, ny float64
	var fid = -1

	// 1. Spatial
	_, foodCells := w.Grid.GetNeighbors(c.X, c.Y, 300.0)
	for _, cell := range foodCells {
		for _, f := range cell {
			if eaten[f.ID] { continue }
			dist := math.Hypot(f.X-c.X, f.Y-c.Y)
			if dist < minDist {
				minDist, nx, ny, fid = dist, f.X, f.Y, f.ID
			}
		}
	}
	// 2. Global Fallback
	if fid == -1 {
		for _, f := range w.Food {
			if eaten[f.ID] { continue }
			dist := math.Hypot(f.X-c.X, f.Y-c.Y)
			if dist < minDist {
				minDist, nx, ny, fid = dist, f.X, f.Y, f.ID
			}
		}
	}
	return nx, ny, minDist, fid
}
