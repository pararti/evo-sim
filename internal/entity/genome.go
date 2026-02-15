package entity

import (
	"math"
	"math/rand/v2"
)

// Genome represents the genetic blueprint of a creature.
// Mutations happen here, not on the creature directly.
type Genome struct {
	// Physical Traits
	SizeGene         float64 // Affects mass and energy storage. Range: [0.5, 3.0]
	SpeedGene        float64 // Affects max speed and energy cost of movement. Range: [0.5, 2.0]
	SenseGene        float64 // Affects view radius. Range: [50.0, 300.0]
	DietGene         float64 // < 0.4: Herbivore, > 0.6: Carnivore, Middle: Omnivore (inefficient). Range: [0.0, 1.0]
	MetabolismGene   float64 // Affects energy usage rate and speed efficiency. Range: [0.5, 2.0]
	FertilityGene    float64 // Affects reproduction threshold. Range: [0.5, 0.9]
	ConstitutionGene float64 // Affects density and max health/energy. Range: [0.5, 1.5]

	// Visual Traits (for debugging/visualization)
	ColorR float64
	ColorG float64
	ColorB float64
}

// NewRandomGenome creates a genome with random traits.
func NewRandomGenome() Genome {
	return Genome{
		SizeGene:         0.5 + rand.Float64()*1.0, // Default slightly small
		SpeedGene:        1.0 + (rand.Float64()-0.5)*0.5,
		SenseGene:        100.0 + (rand.Float64()-0.5)*50.0,
		DietGene:         rand.Float64(), // Random diet strategy
		MetabolismGene:   1.0 + (rand.Float64()-0.5)*0.5,
		FertilityGene:    0.5 + rand.Float64()*0.4,
		ConstitutionGene: 1.0 + (rand.Float64()-0.5)*0.5,
		ColorR:           rand.Float64(),
		ColorG:           rand.Float64(),
		ColorB:           rand.Float64(),
	}
}

// Mutate returns a mutated copy of the genome.
// rate: probability of a gene mutating (0.0 - 1.0)
// strength: magnitude of change
func (g Genome) Mutate(rate, strength float64) Genome {
	ng := g // Copy struct

	// Helper to mutate a float gene within bounds
	mutateFloat := func(val *float64, min, max float64) {
		if rand.Float64() < rate {
			*val += rand.NormFloat64() * strength
			if *val < min {
				*val = min
			}
			if *val > max {
				*val = max
			}
		}
	}

	mutateFloat(&ng.SizeGene, 0.4, 4.0)
	mutateFloat(&ng.SpeedGene, 0.2, 3.0)
	mutateFloat(&ng.SenseGene, 30.0, 500.0)
	mutateFloat(&ng.DietGene, 0.0, 1.0)
	mutateFloat(&ng.MetabolismGene, 0.5, 2.5)
	mutateFloat(&ng.FertilityGene, 0.3, 0.95)
	mutateFloat(&ng.ConstitutionGene, 0.4, 2.0)

	// Color mutation is purely cosmetic but helps visualize lineage
	mutateFloat(&ng.ColorR, 0.0, 1.0)
	mutateFloat(&ng.ColorG, 0.0, 1.0)
	mutateFloat(&ng.ColorB, 0.0, 1.0)

	return ng
}

// Crossover creates a child genome by picking each gene from a random parent (uniform crossover).
func (g Genome) Crossover(other Genome) Genome {
	pick := func(a, b float64) float64 {
		if rand.Float64() < 0.5 {
			return a
		}
		return b
	}
	return Genome{
		SizeGene:         pick(g.SizeGene, other.SizeGene),
		SpeedGene:        pick(g.SpeedGene, other.SpeedGene),
		SenseGene:        pick(g.SenseGene, other.SenseGene),
		DietGene:         pick(g.DietGene, other.DietGene),
		MetabolismGene:   pick(g.MetabolismGene, other.MetabolismGene),
		FertilityGene:    pick(g.FertilityGene, other.FertilityGene),
		ConstitutionGene: pick(g.ConstitutionGene, other.ConstitutionGene),
		ColorR:           pick(g.ColorR, other.ColorR),
		ColorG:           pick(g.ColorG, other.ColorG),
		ColorB:           pick(g.ColorB, other.ColorB),
	}
}

// CalculateStats derives physical stats from genes (Phenotype) with Epistasis.
func (g Genome) CalculateStats() (mass, speed, viewRadius, bmr, maxEnergy, reproductionThreshold float64, isCarnivore bool) {
	// Epistasis: Traits depend on multiple genes

	// Mass depends on Size (volume) and Constitution (density)
	mass = (g.SizeGene * g.SizeGene) * g.ConstitutionGene

	// Speed depends on SpeedGene (muscle structure), Metabolism (energy output), and inversely on Mass
	// High metabolism boosts speed, but high mass slows it down significantly.
	// We use sqrt(mass) to simulate surface area/drag/inertia scaling.
	speed = (g.SpeedGene * g.MetabolismGene) / math.Sqrt(mass)

	viewRadius = g.SenseGene

	// BMR (Basal Metabolic Rate): The energy cost of existing per tick.
	// - Mass costs energy to maintain (scaled by metabolism).
	// - High Speed potential (muscle maintenance) costs energy.
	// - Brain/Senses cost energy.
	bmr = (mass * 0.05 * g.MetabolismGene) + (g.SpeedGene * 0.02) + (g.SenseGene * 0.0001)

	// Max Energy Storage depends on Mass and Constitution (fat reserves)
	maxEnergy = mass * 100.0 * g.ConstitutionGene + 50.0

	// Reproduction Threshold: When is the creature ready to reproduce?
	// Controlled by FertilityGene (strategy: r-selection vs K-selection)
	reproductionThreshold = maxEnergy * g.FertilityGene

	isCarnivore = g.DietGene > 0.6

	return
}

// Distance calculates the genetic distance between two genomes.
// Used for speciation and reproductive isolation.
func (g Genome) Distance(other Genome) float64 {
	dSize := g.SizeGene - other.SizeGene
	dSpeed := g.SpeedGene - other.SpeedGene
	dSense := (g.SenseGene - other.SenseGene) * 0.01 // Normalize sense (range 50-300)
	dDiet := g.DietGene - other.DietGene
	dMeta := g.MetabolismGene - other.MetabolismGene
	dFert := g.FertilityGene - other.FertilityGene
	dConst := g.ConstitutionGene - other.ConstitutionGene

	sumSq := dSize*dSize +
		dSpeed*dSpeed +
		dSense*dSense +
		dDiet*dDiet +
		dMeta*dMeta +
		dFert*dFert +
		dConst*dConst

	return math.Sqrt(sumSq)
}
