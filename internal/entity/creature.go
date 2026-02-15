package entity

import (
	"math"

	"evo-sim/internal/brain"
)

type Creature struct {
	ID          int
	X, Y        float64
	Energy      float64
	
	// Phenotype (derived from Genome)
	Size        float64
	Speed       float64
	ViewRadius  float64
	IsCarnivore bool
	BMR         float64 // Basal Metabolic Rate (energy cost per tick)

	Age         int

	// Genotype
	Genome Genome
	Brain  *brain.Network
}

func NewCreature(id int, x, y float64, inputSize, hiddenSize, outputSize int) *Creature {
	net := brain.NewNetwork(inputSize, hiddenSize, outputSize)
	genome := NewRandomGenome()
	
	// Calculate Phenotype from Genotype
	mass, speed, view, bmr, isCarn := genome.CalculateStats()

	return &Creature{
		ID:          id,
		X:           x,
		Y:           y,
		Energy:      100.0 + (mass * 10), // Larger creatures start with more energy reserves
		
		Size:        genome.SizeGene, // Using SizeGene directly as visual size for now
		Speed:       speed,
		ViewRadius:  view,
		BMR:         bmr,
		IsCarnivore: isCarn,

		Age:    0,
		Genome: genome,
		Brain:  net,
	}
}

func (c *Creature) Update(foodX, foodY, enemyX, enemyY, targetIsCarnivore, terrainSpeedFactor, terrainEnergyFactor, worldW, worldH, maxAge float64) {
	// Inputs normalized relative to ViewRadius where possible
	// 1-2: relative food pos
	// 3-4: relative creature pos
	// 5: current energy 
	// 6: target role
	// 7-10: distances to walls
	
	// Normalize distance by ViewRadius to allow evolution of sensing range
	// If things are outside ViewRadius, inputs should be 0 (simulated by caller or here)
	// For now, we keep global normalization but scaled 
	
	input := []float64{
		(foodX - c.X) / c.ViewRadius,
		(foodY - c.Y) / c.ViewRadius,
		(enemyX - c.X) / c.ViewRadius,
		(enemyY - c.Y) / c.ViewRadius,
		c.Energy / 200.0,
		targetIsCarnivore,
		c.X / worldW,             
		(worldW - c.X) / worldW,  
		c.Y / worldH,             
		(worldH - c.Y) / worldH,  
	}

	output := c.Brain.FeedForward(input)

	// Movement
	// Output is [-1, 1]. Speed is max speed.
	// Apply terrain penalty
	currentMaxSpeed := c.Speed * terrainSpeedFactor
	
	dx := output[0] * currentMaxSpeed
	dy := output[1] * currentMaxSpeed
	
	c.X += dx
	c.Y += dy
	
	// Energy Calculation (Thermodynamics)
	// 1. Basal Metabolic Rate (Living cost)
	// 2. Movement Cost (Work = Force * Distance). F = ma. Heavier creatures spend more energy moving.
	// 3. Terrain Resistance (Mud/Water makes it harder)
	// 4. Aging Cost (Gradient Aging). As creatures age, they become less efficient.
	movementDist := math.Sqrt(dx*dx + dy*dy)
	movementCost := movementDist * (c.Size * c.Size) * 0.1 * terrainEnergyFactor // Mass ~ Size^2

	// Calculate aging factor: 1 + (Age/MaxAge)^2
	// This means young creatures pay ~1x BMR, but old ones pay significantly more.
	// At MaxAge, they pay 2x BMR. Past MaxAge, it skyrockets.
	ageRatio := float64(c.Age) / maxAge
	agingFactor := 1.0 + (ageRatio * ageRatio)

	c.Energy -= (c.BMR * agingFactor) + movementCost
	
	c.Age++
}

func (c *Creature) ReproduceAsexual(mutationRate, mutationStrength float64) *Creature {
	childBrain := c.Brain.Clone()
	childBrain.Mutate(mutationRate, mutationStrength)

	// Mutate Genome
	childGenome := c.Genome.Mutate(mutationRate, mutationStrength)
	
	// Calculate new Phenotype
	_, speed, view, bmr, isCarn := childGenome.CalculateStats()

	child := &Creature{
		ID:          0, // To be assigned by world
		X:           c.X,
		Y:           c.Y,
		Energy:      c.Energy / 2, // Parent gives half energy
		
		Size:        childGenome.SizeGene,
		Speed:       speed,
		ViewRadius:  view,
		BMR:         bmr,
		IsCarnivore: isCarn,
		
		Age:         0,
		Genome:      childGenome,
		Brain:       childBrain,
	}

	c.Energy /= 2
	return child
}

func (c *Creature) ReproduceSexual(mate *Creature, mutationRate, mutationStrength float64) *Creature {
	// Crossover genomes + mutate
	childGenome := c.Genome.Crossover(mate.Genome)
	childGenome = childGenome.Mutate(mutationRate, mutationStrength)

	// Crossover brains + mutate
	childBrain := c.Brain.Crossover(mate.Brain)
	childBrain.Mutate(mutationRate, mutationStrength)

	_, speed, view, bmr, isCarn := childGenome.CalculateStats()

	// Each parent gives 1/3 of energy
	energyFromP1 := c.Energy / 3
	energyFromP2 := mate.Energy / 3
	c.Energy -= energyFromP1
	mate.Energy -= energyFromP2

	return &Creature{
		ID:          0,
		X:           (c.X + mate.X) / 2,
		Y:           (c.Y + mate.Y) / 2,
		Energy:      energyFromP1 + energyFromP2,
		Size:        childGenome.SizeGene,
		Speed:       speed,
		ViewRadius:  view,
		BMR:         bmr,
		IsCarnivore: isCarn,
		Age:         0,
		Genome:      childGenome,
		Brain:       childBrain,
	}
}
