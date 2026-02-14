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
	Mu        sync.RWMutex
}

func NewWorld(cfg *config.Config) *World {
	w := &World{
		Cfg: cfg,
	}

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

	var newChildren []*entity.Creature

	for i := len(w.Creatures) - 1; i >= 0; i-- {
		c := w.Creatures[i]

		// 1. Find targets
		foodX, foodY, foodDist := w.findNearestFood(c)
		targetX, targetY, targetDist, isTargetCarnivore := w.findNearestCreature(c, i)

		// 2. Update Brain and Physics
		roleVal := -1.0
		if isTargetCarnivore {
			roleVal = 1.0
		}
		w.updateCreature(c, foodX, foodY, targetX, targetY, roleVal)

		// 3. Interactions
		// Eating Food (only if not strictly carnivore, or let everyone eat food for now but carnivores get less?)
		// Let's say: Herbivores eat food, Carnivores eat creatures.
		if !c.IsCarnivore && foodDist < w.Cfg.EatRadius*c.Size {
			c.Energy += w.Cfg.FoodEnergy
			w.eatFood(foodX, foodY)
		}

		// Predation logic
		if c.IsCarnivore && targetDist < w.Cfg.EatRadius*c.Size {
			target := w.getCreatureAt(targetX, targetY)
			if target != nil && target.Size < c.Size*1.2 { // Can only eat if not much bigger
				c.Energy += target.Energy * 0.8 // Gain most of prey's energy
				w.removeCreature(target.ID)
			}
		}

		// 4. Life Cycle
		// Reproduction
		if c.Energy > w.Cfg.ReproduceThreshold {
			child := c.Reproduce(w.Cfg.MutationRate, w.Cfg.MutationStrength)
			child.ID = rand.IntN(1000000) // Simple ID generation
			newChildren = append(newChildren, child)
		}

		// Death (Energy or Age)
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

func (w *World) findNearestCreature(c *entity.Creature, currentIndex int) (float64, float64, float64, bool) {
	minDist := math.MaxFloat64
	var nx, ny float64
	var isCarnivore bool

	for i, other := range w.Creatures {
		if i == currentIndex {
			continue
		}
		dist := math.Hypot(other.X-c.X, other.Y-c.Y)
		if dist < minDist {
			minDist = dist
			nx, ny = other.X, other.Y
			isCarnivore = other.IsCarnivore
		}
	}
	return nx, ny, minDist, isCarnivore
}

func (w *World) getCreatureAt(x, y float64) *entity.Creature {
	for _, c := range w.Creatures {
		if c.X == x && c.Y == y {
			return c
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

	for _, f := range w.Food {
		dist := math.Hypot(f.X-c.X, f.Y-c.Y)
		if dist < minDist {
			minDist = dist
			nearestX = f.X
			nearestY = f.Y
		}
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
