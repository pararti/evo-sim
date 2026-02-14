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
		IsCarnivore: rand.Float64() < 0.2, // 20% chance to be born a carnivore
		Size:        0.8 + rand.Float64()*0.4, // Random size [0.8, 1.2]
		Age:         0,
		Brain:       net,
	}
}

func (c *Creature) Update(foodX, foodY, enemyX, enemyY, targetIsCarnivore, speedCof, energyLoss, worldW, worldH float64) {
	// Normalized Inputs [-1, 1]
	// 1-2: relative food pos
	// 3-4: relative creature pos
	// 5: current energy (normalized to ~200 max)
	// 6: target role
	// 7-10: distances to walls (Left, Right, Top, Bottom)
	
	input := []float64{
		(foodX - c.X) / worldW,
		(foodY - c.Y) / worldH,
		(enemyX - c.X) / worldW,
		(enemyY - c.Y) / worldH,
		c.Energy / 100.0,
		targetIsCarnivore,
		c.X / worldW,             // Dist to left
		(worldW - c.X) / worldW,  // Dist to right
		c.Y / worldH,             // Dist to top
		(worldH - c.Y) / worldH,  // Dist to bottom
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
	if rand.Float64() < 0.1 { // 10% chance to mutate role
		child.IsCarnivore = !child.IsCarnivore
	}
	if rand.Float64() < mutationRate {
		// Change size slightly
		child.Size += (rand.Float64() - 0.5) * 0.1
		if child.Size < 0.5 {
			child.Size = 0.5
		}
		if child.Size > 2.0 {
			child.Size = 2.0
		}
	}

	c.Energy /= 2
	return child
}
