package world

import (
	"math"
	"math/rand"
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

	for i := len(w.Creatures) - 1; i >= 0; i-- {
		c := w.Creatures[i]

		fx, fy, minDist := w.findNearestFood(c)

		c.Update(fx, fy, w.Cfg.SpeedFactor, w.Cfg.MoveCost)

		//NO out from world!
		if c.X < 0 {
			c.X = 0
		}
		if c.X > w.Cfg.WorldWidth {
			c.X = w.Cfg.WorldWidth
		}
		if c.Y < 0 {
			c.Y = 0
		}
		if c.Y > w.Cfg.WorldHeight {
			c.Y = w.Cfg.WorldHeight
		}

		//ate ?
		if minDist < w.Cfg.EatRadius {
			c.Energy += w.Cfg.FoodEnergy
			w.eatFood(fx, fy) // delete food
		}

		//die?
		if c.Energy <= 0 {
			w.Creatures[i] = w.Creatures[len(w.Creatures)-1]
			w.Creatures = w.Creatures[:len(w.Creatures)-1]
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
