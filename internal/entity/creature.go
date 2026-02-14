package entity

import (
	"math/rand/v2"

	"evo-sim/internal/brain"
)

type Creature struct {
	ID          int
	X, Y        float64
	Energy      float64
	IsCarnivore bool
	Size        float64
	Age         int

	Brain *brain.Network
}

func NewCreature(id int, x, y float64, inputSize, hiddenSize, outputSize int) *Creature {
	net := brain.NewNetwork(inputSize, hiddenSize, outputSize)

	return &Creature{
		ID:          id,
		X:           x,
		Y:           y,
		Energy:      100.0,
		IsCarnivore: false, // Default to herbivore
		Size:        1.0,   // Base size
		Age:         0,
		Brain:       net,
	}
}

func (c *Creature) Update(foodX, foodY, enemyX, enemyY, targetIsCarnivore, speedCof, energyLoss float64) {
	// New expanded input: 
	// 1-2: relative food pos
	// 3-4: relative creature/threat pos
	// 5: current energy
	// 6: is the target a carnivore? (1.0 or -1.0)
	input := []float64{
		foodX - c.X,
		foodY - c.Y,
		enemyX - c.X,
		enemyY - c.Y,
		c.Energy,
		targetIsCarnivore,
	}

	output := c.Brain.FeedForward(input)

	// Metabolic cost depends on size and speed
	actualSpeed := speedCof / (c.Size * 0.5) // Larger creatures are slightly slower
	c.X += output[0] * actualSpeed
	c.Y += output[1] * actualSpeed
	
	// Energy loss: base + movement + size penalty
	movement := (output[0]*output[0] + output[1]*output[1]) * 0.05
	c.Energy -= energyLoss + movement + (c.Size * 0.02)
	c.Age++
}

func (c *Creature) Reproduce(mutationRate, mutationStrength float64) *Creature {
	childBrain := c.Brain.Clone()
	childBrain.Mutate(mutationRate, mutationStrength)

	child := &Creature{
		ID:          0,
		X:           c.X,
		Y:           c.Y,
		Energy:      c.Energy / 2,
		IsCarnivore: c.IsCarnivore,
		Size:        c.Size,
		Age:         0,
		Brain:       childBrain,
	}

	// Mutation of biological traits
	if rand.Float64() < mutationRate {
		// Toggle carnivore role
		child.IsCarnivore = !child.IsCarnivore
	}
	if rand.Float64() < mutationRate {
		// Change size slightly
		child.Size += (rand.Float64() - 0.5) * 0.2
		if child.Size < 0.5 {
			child.Size = 0.5
		}
	}

	c.Energy /= 2
	return child
}
