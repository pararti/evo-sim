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
		Grid: NewGrid(cfg.WorldWidth, cfg.WorldHeight, 40.0), // 40px cell size
	}
    // ... rest of NewWorld

	for i := 0; i < cfg.InitialPop; i++ {
		w.Creatures = append(w.Creatures, entity.NewCreature(
			i,
			rand.Float64()*cfg.WorldWidth,
			rand.Float64()*cfg.WorldHeight,
			cfg.InputSize,
			cfg.HiddenSize,
			cfg.OutputSize,
		))
	}

	for i := 0; i < cfg.FoodCount; i++ {
		w.spawnFood()
	}

	return w
}

func (w *World) spawnFood() {
	w.Food = append(w.Food, entity.Food{
		X: rand.Float64() * w.Cfg.WorldWidth,
		Y: rand.Float64() * w.Cfg.WorldHeight,
	})
}

func (w *World) Update() {
	w.Mu.Lock()
	defer w.Mu.Unlock()

	// Rebuild grid for spatial optimization
	w.Grid.Clear()
	for _, c := range w.Creatures {
		w.Grid.InsertCreature(c)
	}
	for _, f := range w.Food {
		w.Grid.InsertFood(f)
	}

	var newChildren []*entity.Creature

	for i := len(w.Creatures) - 1; i >= 0; i-- {
		c := w.Creatures[i]

		// 1. Find targets using Grid (look within 200px radius for performance/relevance)
		foodX, foodY, foodDist := w.findNearestFood(c)
		targetX, targetY, targetDist, isTargetCarnivore := w.findNearestCreature(c)

		// 2. Update Brain and Physics
		roleVal := -1.0
		if isTargetCarnivore {
			roleVal = 1.0
		}
		w.updateCreature(c, foodX, foodY, targetX, targetY, roleVal)

		// 3. Interactions
		if !c.IsCarnivore && foodDist < w.Cfg.EatRadius*c.Size {
			c.Energy += w.Cfg.FoodEnergy
			w.eatFood(foodX, foodY)
		}

		if c.IsCarnivore && targetDist < w.Cfg.EatRadius*c.Size {
			target := w.getCreatureAt(targetX, targetY)
			if target != nil && target.ID != c.ID && target.Size < c.Size*1.2 {
				c.Energy += target.Energy * 0.8
				w.removeCreature(target.ID)
			}
		}

		// 4. Life Cycle
		if c.Energy > w.Cfg.ReproduceThreshold {
			child := c.Reproduce(w.Cfg.MutationRate, w.Cfg.MutationStrength)
			child.ID = rand.IntN(1000000)
			newChildren = append(newChildren, child)
		}

		if c.Energy <= 0 || c.Age > 10000 {
			w.Creatures = append(w.Creatures[:i], w.Creatures[i+1:]...)
		}
	}

	w.Creatures = append(w.Creatures, newChildren...)
}

func (w *World) updateCreature(c *entity.Creature, fx, fy, tx, ty, roleVal float64) {
	c.Update(fx, fy, tx, ty, roleVal, w.Cfg.SpeedFactor, w.Cfg.MoveCost)

	// World boundaries
	if c.X < 0 { c.X = 0 }
	if c.X > w.Cfg.WorldWidth { c.X = w.Cfg.WorldWidth }
	if c.Y < 0 { c.Y = 0 }
	if c.Y > w.Cfg.WorldHeight { c.Y = w.Cfg.WorldHeight }
}

func (w *World) findNearestCreature(c *entity.Creature) (float64, float64, float64, bool) {
	minDist := math.MaxFloat64
	var nx, ny float64
	var isCarnivore bool

	// Query nearby cells (radius 150.0 is enough for most interactions)
	cells, _ := w.Grid.GetNeighbors(c.X, c.Y, 150.0)

	for _, cell := range cells {
		for _, other := range cell {
			if other.ID == c.ID {
				continue
			}
			dist := math.Hypot(other.X-c.X, other.Y-c.Y)
			if dist < minDist {
				minDist = dist
				nx, ny = other.X, other.Y
				isCarnivore = other.IsCarnivore
			}
		}
	}
	return nx, ny, minDist, isCarnivore
}

func (w *World) getCreatureAt(x, y float64) *entity.Creature {
	// Precise lookup still needs care, but with Grid we can narrow it down
	cells, _ := w.Grid.GetNeighbors(x, y, 5.0)
	for _, cell := range cells {
		for _, c := range cell {
			if c.X == x && c.Y == y {
				return c
			}
		}
	}
	return nil
}

func (w *World) removeCreature(id int) {
	for i, c := range w.Creatures {
		if c.ID == id {
			w.Creatures = append(w.Creatures[:i], w.Creatures[i+1:]...)
			return
		}
	}
}

func (w *World) findNearestFood(c *entity.Creature) (float64, float64, float64) {
	minDist := math.MaxFloat64
	var nearestX, nearestY float64

	_, foodCells := w.Grid.GetNeighbors(c.X, c.Y, 200.0)

	for _, cell := range foodCells {
		for _, f := range cell {
			dist := math.Hypot(f.X-c.X, f.Y-c.Y)
			if dist < minDist {
				minDist = dist
				nearestX = f.X
				nearestY = f.Y
			}
		}
	}
	
	// If no food nearby, return far away point
	if minDist == math.MaxFloat64 {
		return -1000, -1000, minDist
	}

	return nearestX, nearestY, minDist
}

func (w *World) eatFood(x, y float64) {
	for i, f := range w.Food {
		if f.X == x && f.Y == y {
			//Delete food
			w.Food[i] = w.Food[len(w.Food)-1]
			w.Food = w.Food[:len(w.Food)-1]

			w.spawnFood()
			return
		}
	}
}
