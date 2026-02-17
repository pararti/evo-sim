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
	Cfg            *config.Config
	Creatures      []*entity.Creature
	Food           []entity.Food
	Grid           *Grid
	Terrain        *TerrainGrid
	Pheromone      *PheromoneGrid
	SpeciesManager *SpeciesManager
	Mu             sync.RWMutex
	StartTime      time.Time

	// Control Logic
	FoodSpawnAccumulator float64
}

func NewWorld(cfg *config.Config) *World {
	w := &World{
		Cfg:                  cfg,
		Grid:                 NewGrid(cfg.WorldWidth, cfg.WorldHeight, 40.0),
		Terrain:              NewTerrainGrid(cfg.WorldWidth, cfg.WorldHeight, 20.0),
		Pheromone:            NewPheromoneGrid(cfg.WorldWidth, cfg.WorldHeight, 20.0),
		SpeciesManager:       NewSpeciesManager(cfg.SpeciationThreshold),
		StartTime:            time.Now(),
		FoodSpawnAccumulator: 0.0,
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
					w.Cfg.OutputSize,
					w.Cfg.BrainCostPerNeuron,
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
	var newCarrion []entity.Food
	deadCreatures := make(map[int]bool)
	eatenFood := make(map[int]bool)
	matedThisTick := make(map[int]bool)
	maturityAge := int(w.Cfg.MaxAge * w.Cfg.MaturityAgeFraction)

	// 2. Main Simulation Loop
	for _, c := range w.Creatures {
		if deadCreatures[c.ID] {
			continue
		}

		// Find targets
		foodX, foodY, foodDist, foodID, foodEnergy := w.findNearestFood(c, eatenFood)
		targetX, targetY, targetDist, targetID, targetDiet := w.findNearestCreature(c, deadCreatures)

		// Continuous diet signal: maps DietGene [0,1] → [-1,1]
		roleVal := targetDiet*2.0 - 1.0

		// Get Terrain Physics
		speedFactor, energyCostFactor := w.Terrain.GetMovementPenalty(c.X, c.Y)

		// Calculate Crowding Stress
		neighbors := 0
		w.Grid.ForEachNeighbor(c.X, c.Y, w.Cfg.CrowdingDistance, func(other *entity.Creature) {
			if other.ID != c.ID && !deadCreatures[other.ID] {
				neighbors++
			}
		}, nil)
		stressFactor := 1.0 + float64(neighbors)*w.Cfg.CrowdingMultiplier

		pheromoneVal := w.Pheromone.Get(c.X, c.Y)
		c.Update(foodX, foodY, targetX, targetY, roleVal, speedFactor, energyCostFactor, w.Cfg.WorldWidth, w.Cfg.WorldHeight, w.Cfg.MaxAge, stressFactor, pheromoneVal)

		// Deposit pheromone trail
		w.Pheromone.Deposit(c.X, c.Y, w.Cfg.PheromoneDeposit)

		// Boundaries
		if c.X < 0 {
			c.X = 0
		} else if c.X > w.Cfg.WorldWidth {
			c.X = w.Cfg.WorldWidth
		}
		if c.Y < 0 {
			c.Y = 0
		} else if c.Y > w.Cfg.WorldHeight {
			c.Y = w.Cfg.WorldHeight
		}

		// Interactions — continuous diet spectrum
		// Food eating: any creature can eat, efficiency depends on DietGene
		if foodID != -1 && foodDist < w.Cfg.EatRadius*c.Size {
			if !eatenFood[foodID] {
				if foodEnergy > 0 {
					// Carrion (dead creature remains): carnivores benefit more
					c.Energy += foodEnergy * c.Genome.ExpressedDiet()
				} else {
					// Plant: herbivores benefit more
					c.Energy += w.Cfg.FoodEnergy * (1.0 - c.Genome.ExpressedDiet())
				}
				eatenFood[foodID] = true
			}
		}

		// Hunting: efficiency scales with DietGene
		if targetID != -1 && targetDist < w.Cfg.EatRadius*c.Size && c.Genome.ExpressedDiet() > 0 {
			if !deadCreatures[targetID] {
				target := w.getCreatureByID(targetID)
				if target != nil && target.Size < c.Size*1.2 {
					c.Energy += target.Energy * c.Genome.ExpressedDiet() * 0.8
					deadCreatures[targetID] = true
					w.SpeciesManager.RemoveCreature(target.SpeciesID)
					// Spawn carrion from kill remains (30% of body mass remains)
					newCarrion = append(newCarrion, entity.Food{
						ID:         rand.IntN(10000000),
						X:          target.X,
						Y:          target.Y,
						Energy:     target.Mass * w.Cfg.CarrionEnergyMult * 0.3,
						DecayTicks: w.Cfg.CarrionLifespan,
					})
				}
			}
		}

		// Reproduction — requires maturity age
		if c.Energy > c.ReproductionThreshold && !matedThisTick[c.ID] && c.Age >= maturityAge {
			mate := w.findMate(c, deadCreatures, matedThisTick)
			var child *entity.Creature
			if mate != nil {
				child = c.ReproduceSexual(mate, w.Cfg.MutationRate, w.Cfg.MutationStrength, w.Cfg.InbreedingThreshold, w.Cfg.InbreedingPenalty, w.Cfg.BrainCostPerNeuron)
				matedThisTick[mate.ID] = true
			} else if c.Energy > c.ReproductionThreshold*w.Cfg.AsexualThresholdMult {
				child = c.ReproduceAsexual(w.Cfg.MutationRate, w.Cfg.MutationStrength, w.Cfg.BrainCostPerNeuron)
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
			// Spawn carrion from natural death
			newCarrion = append(newCarrion, entity.Food{
				ID:         rand.IntN(10000000),
				X:          c.X,
				Y:          c.Y,
				Energy:     c.Mass * w.Cfg.CarrionEnergyMult,
				DecayTicks: w.Cfg.CarrionLifespan,
			})
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

	// Remove eaten food and decay carrion
	newFoodList := make([]entity.Food, 0, len(w.Food))
	for _, f := range w.Food {
		if eatenFood[f.ID] {
			continue
		}
		if f.DecayTicks > 0 {
			f.DecayTicks--
			if f.DecayTicks <= 0 {
				continue // Carrion fully decayed
			}
		}
		newFoodList = append(newFoodList, f)
	}
	w.Food = append(newFoodList, newCarrion...)

	// Dynamic Food Spawning (Entropy control)
	w.FoodSpawnAccumulator += w.Cfg.FoodSpawnChance
	for w.FoodSpawnAccumulator >= 1.0 {
		w.spawnFood()
		w.FoodSpawnAccumulator -= 1.0
	}

	// 4. Decay pheromones
	w.Pheromone.Decay(w.Cfg.PheromoneDecay)

	// 5. Rescue population
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

func (w *World) findNearestCreature(c *entity.Creature, dead map[int]bool) (float64, float64, float64, int, float64) {
	minDist := math.MaxFloat64
	var nx, ny float64
	var targetID = -1
	var targetDiet float64

	w.Grid.ForEachNeighbor(c.X, c.Y, c.ViewRadius, func(other *entity.Creature) {
		if other.ID == c.ID || dead[other.ID] {
			return
		}
		dist := math.Hypot(other.X-c.X, other.Y-c.Y)
		if dist < c.ViewRadius && dist < minDist {
			minDist, nx, ny, targetID, targetDiet = dist, other.X, other.Y, other.ID, other.Genome.ExpressedDiet()
		}
	}, nil)

	return nx, ny, minDist, targetID, targetDiet
}

func (w *World) findNearestFood(c *entity.Creature, eaten map[int]bool) (float64, float64, float64, int, float64) {
	minDist := math.MaxFloat64
	var nx, ny float64
	var fid = -1
	var fEnergy float64

	w.Grid.ForEachNeighbor(c.X, c.Y, c.ViewRadius, nil, func(f entity.Food) {
		if eaten[f.ID] {
			return
		}
		dist := math.Hypot(f.X-c.X, f.Y-c.Y)
		if dist < c.ViewRadius && dist < minDist {
			minDist, nx, ny, fid, fEnergy = dist, f.X, f.Y, f.ID, f.Energy
		}
	})

	return nx, ny, minDist, fid, fEnergy
}

func (w *World) findMate(c *entity.Creature, dead, mated map[int]bool) *entity.Creature {
	var best *entity.Creature
	bestScore := -1.0

	w.Grid.ForEachNeighbor(c.X, c.Y, c.ViewRadius, func(other *entity.Creature) {
		if other.ID == c.ID || dead[other.ID] || mated[other.ID] {
			return
		}
		if other.Energy <= other.ReproductionThreshold {
			return
		}
		geneticDist := c.Genome.Distance(other.Genome)
		if geneticDist > w.Cfg.MatingDistanceThreshold {
			return
		}
		dist := math.Hypot(other.X-c.X, other.Y-c.Y)
		if dist < c.ViewRadius && dist < w.Cfg.EatRadius*c.Size {
			// Sexual selection: score based on fitness indicators
			energyRatio := other.Energy / other.MaxEnergy
			geneticCompat := 1.0 - (geneticDist / w.Cfg.MatingDistanceThreshold)
			score := energyRatio*0.5 + (other.Size/4.0)*0.2 + geneticCompat*0.3
			if score > bestScore {
				bestScore = score
				best = other
			}
		}
	}, nil)

	return best
}
