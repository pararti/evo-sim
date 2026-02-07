package entity

import "evo-sim/internal/brain"

type Creature struct {
	ID     int
	X, Y   float64
	Energy float64

	Brain *brain.Network
}

func NewCreature(id int, x, y float64, inputSize, hiddenSize, outputSize int) *Creature {
	net := brain.NewNetwork(inputSize, hiddenSize, outputSize)

	return &Creature{
		ID:     id,
		X:      x,
		Y:      y,
		Energy: 100.0,
		Brain:  net,
	}
}

func (c *Creature) Update(foodX, foodY, speedCof, energyLoss float64) {
	input := []float64{
		foodX - c.X,
		foodY - c.Y,
		c.Energy,
	}

	output := c.Brain.FeedForward(input)

	c.X += output[0] * speedCof
	c.Y += output[1] * speedCof
	c.Energy -= energyLoss
}

func (c *Creature) Reproduce(mutationRate, mutationStrength float64) *Creature {
	// 1. create  child
	childBrain := c.Brain.Clone()

	// 2. random mutation
	childBrain.Mutate(mutationRate, mutationStrength)

	child := &Creature{
		// ID присвоим позже в World, или можно использовать UUID
		// Пока поставим 0
		ID:     0,
		X:      c.X,
		Y:      c.Y,
		Energy: c.Energy / 2, // Родитель отдает половину энергии ребенку
		Brain:  childBrain,
	}

	// Родитель теряет энергию
	c.Energy /= 2

	return child
}
