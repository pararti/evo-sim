package entity

import (
	"math"
	"math/rand/v2"
)

// Genome represents a diploid genetic blueprint.
// Each trait gene has two alleles. Phenotype is derived via dominance rules.
type Genome struct {
	// Diploid trait alleles (two per gene)
	SizeAllele1, SizeAllele2                 float64 // Dominant: max(a1,a2)
	SpeedAllele1, SpeedAllele2               float64 // Dominant: max(a1,a2)
	SenseAllele1, SenseAllele2               float64 // Additive: avg(a1,a2)
	DietAllele1, DietAllele2                 float64 // Additive: avg(a1,a2)
	MetabolismAllele1, MetabolismAllele2     float64 // Additive: avg(a1,a2)
	FertilityAllele1, FertilityAllele2       float64 // Additive: avg(a1,a2)
	ConstitutionAllele1, ConstitutionAllele2 float64 // Additive: avg(a1,a2)
	HiddenAllele1, HiddenAllele2             float64 // Additive: avg(a1,a2)

	// Visual Traits (haploid — simple pick)
	ColorR float64
	ColorG float64
	ColorB float64
}

// Phenotypic accessors — expressed gene values from diploid alleles.

func (g Genome) ExpressedSize() float64         { return math.Max(g.SizeAllele1, g.SizeAllele2) }
func (g Genome) ExpressedSpeed() float64        { return math.Max(g.SpeedAllele1, g.SpeedAllele2) }
func (g Genome) ExpressedSense() float64        { return (g.SenseAllele1 + g.SenseAllele2) / 2 }
func (g Genome) ExpressedDiet() float64         { return (g.DietAllele1 + g.DietAllele2) / 2 }
func (g Genome) ExpressedMetabolism() float64   { return (g.MetabolismAllele1 + g.MetabolismAllele2) / 2 }
func (g Genome) ExpressedFertility() float64    { return (g.FertilityAllele1 + g.FertilityAllele2) / 2 }
func (g Genome) ExpressedConstitution() float64 { return (g.ConstitutionAllele1 + g.ConstitutionAllele2) / 2 }
func (g Genome) ExpressedHidden() float64       { return (g.HiddenAllele1 + g.HiddenAllele2) / 2 }

// NewRandomGenome creates a genome with random diploid traits.
func NewRandomGenome() Genome {
	randSize := func() float64 { return 0.5 + rand.Float64()*1.0 }
	randSpeed := func() float64 { return 1.0 + (rand.Float64()-0.5)*0.5 }
	randSense := func() float64 { return 100.0 + (rand.Float64()-0.5)*50.0 }
	randDiet := func() float64 { return rand.Float64() }
	randMeta := func() float64 { return 1.0 + (rand.Float64()-0.5)*0.5 }
	randFert := func() float64 { return 0.5 + rand.Float64()*0.4 }
	randConst := func() float64 { return 1.0 + (rand.Float64()-0.5)*0.5 }
	randHidden := func() float64 { return 4.0 + rand.Float64()*4.0 }

	return Genome{
		SizeAllele1: randSize(), SizeAllele2: randSize(),
		SpeedAllele1: randSpeed(), SpeedAllele2: randSpeed(),
		SenseAllele1: randSense(), SenseAllele2: randSense(),
		DietAllele1: randDiet(), DietAllele2: randDiet(),
		MetabolismAllele1: randMeta(), MetabolismAllele2: randMeta(),
		FertilityAllele1: randFert(), FertilityAllele2: randFert(),
		ConstitutionAllele1: randConst(), ConstitutionAllele2: randConst(),
		HiddenAllele1: randHidden(), HiddenAllele2: randHidden(),
		ColorR: rand.Float64(),
		ColorG: rand.Float64(),
		ColorB: rand.Float64(),
	}
}

// Mutate returns a mutated copy of the genome.
// Each allele mutates independently.
func (g Genome) Mutate(rate, strength float64) Genome {
	ng := g

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

	// Each allele mutates independently
	mutateFloat(&ng.SizeAllele1, 0.4, 4.0)
	mutateFloat(&ng.SizeAllele2, 0.4, 4.0)
	mutateFloat(&ng.SpeedAllele1, 0.2, 3.0)
	mutateFloat(&ng.SpeedAllele2, 0.2, 3.0)
	mutateFloat(&ng.SenseAllele1, 30.0, 500.0)
	mutateFloat(&ng.SenseAllele2, 30.0, 500.0)
	mutateFloat(&ng.DietAllele1, 0.0, 1.0)
	mutateFloat(&ng.DietAllele2, 0.0, 1.0)
	mutateFloat(&ng.MetabolismAllele1, 0.5, 2.5)
	mutateFloat(&ng.MetabolismAllele2, 0.5, 2.5)
	mutateFloat(&ng.FertilityAllele1, 0.3, 0.95)
	mutateFloat(&ng.FertilityAllele2, 0.3, 0.95)
	mutateFloat(&ng.ConstitutionAllele1, 0.4, 2.0)
	mutateFloat(&ng.ConstitutionAllele2, 0.4, 2.0)

	// Brain size alleles: ±1 step, clamped [3, 12]
	mutateHidden := func(val *float64) {
		if rand.Float64() < rate {
			*val += rand.NormFloat64() * 1.0
			if *val < 3.0 {
				*val = 3.0
			}
			if *val > 12.0 {
				*val = 12.0
			}
		}
	}
	mutateHidden(&ng.HiddenAllele1)
	mutateHidden(&ng.HiddenAllele2)

	// Color mutation (haploid)
	mutateFloat(&ng.ColorR, 0.0, 1.0)
	mutateFloat(&ng.ColorG, 0.0, 1.0)
	mutateFloat(&ng.ColorB, 0.0, 1.0)

	return ng
}

// Crossover creates a child genome via diploid meiosis.
// Each parent donates one random allele per gene.
func (g Genome) Crossover(other Genome) Genome {
	// Pick one allele from each parent per gene
	pickOne := func(a1, a2 float64) float64 {
		if rand.Float64() < 0.5 {
			return a1
		}
		return a2
	}
	pick := func(a, b float64) float64 {
		if rand.Float64() < 0.5 {
			return a
		}
		return b
	}

	return Genome{
		// Child allele1 = one from parent1, allele2 = one from parent2
		SizeAllele1: pickOne(g.SizeAllele1, g.SizeAllele2),
		SizeAllele2: pickOne(other.SizeAllele1, other.SizeAllele2),

		SpeedAllele1: pickOne(g.SpeedAllele1, g.SpeedAllele2),
		SpeedAllele2: pickOne(other.SpeedAllele1, other.SpeedAllele2),

		SenseAllele1: pickOne(g.SenseAllele1, g.SenseAllele2),
		SenseAllele2: pickOne(other.SenseAllele1, other.SenseAllele2),

		DietAllele1: pickOne(g.DietAllele1, g.DietAllele2),
		DietAllele2: pickOne(other.DietAllele1, other.DietAllele2),

		MetabolismAllele1: pickOne(g.MetabolismAllele1, g.MetabolismAllele2),
		MetabolismAllele2: pickOne(other.MetabolismAllele1, other.MetabolismAllele2),

		FertilityAllele1: pickOne(g.FertilityAllele1, g.FertilityAllele2),
		FertilityAllele2: pickOne(other.FertilityAllele1, other.FertilityAllele2),

		ConstitutionAllele1: pickOne(g.ConstitutionAllele1, g.ConstitutionAllele2),
		ConstitutionAllele2: pickOne(other.ConstitutionAllele1, other.ConstitutionAllele2),

		HiddenAllele1: pickOne(g.HiddenAllele1, g.HiddenAllele2),
		HiddenAllele2: pickOne(other.HiddenAllele1, other.HiddenAllele2),

		// Colors remain haploid (simple pick)
		ColorR: pick(g.ColorR, other.ColorR),
		ColorG: pick(g.ColorG, other.ColorG),
		ColorB: pick(g.ColorB, other.ColorB),
	}
}

// CalculateStats derives physical stats from expressed (phenotypic) genes with Epistasis.
func (g Genome) CalculateStats(brainCostPerNeuron float64) (mass, speed, viewRadius, bmr, maxEnergy, reproductionThreshold float64, isCarnivore bool, hiddenSize int) {
	// Express diploid alleles to phenotype
	sizeGene := g.ExpressedSize()
	speedGene := g.ExpressedSpeed()
	senseGene := g.ExpressedSense()
	dietGene := g.ExpressedDiet()
	metabolismGene := g.ExpressedMetabolism()
	fertilityGene := g.ExpressedFertility()
	constitutionGene := g.ExpressedConstitution()
	hiddenGene := g.ExpressedHidden()

	// Mass depends on Size (volume) and Constitution (density)
	mass = (sizeGene * sizeGene) * constitutionGene

	// Speed depends on SpeedGene, Metabolism, and inversely on Mass
	speed = (speedGene * metabolismGene) / math.Sqrt(mass)

	viewRadius = senseGene

	// Hidden layer size from HiddenGene
	hiddenSize = int(math.Round(hiddenGene))
	if hiddenSize < 3 {
		hiddenSize = 3
	}
	if hiddenSize > 12 {
		hiddenSize = 12
	}

	// BMR
	bmr = (mass * 0.05 * metabolismGene) + (speedGene * 0.02) + (senseGene * 0.0001) + (float64(hiddenSize) * brainCostPerNeuron)

	// Max Energy Storage
	maxEnergy = mass * 100.0 * constitutionGene + 50.0

	// Reproduction Threshold
	reproductionThreshold = maxEnergy * fertilityGene

	isCarnivore = dietGene > 0.6

	return
}

// Distance calculates the phenotypic distance between two genomes.
// Uses expressed (phenotypic) values for comparison.
func (g Genome) Distance(other Genome) float64 {
	dSize := g.ExpressedSize() - other.ExpressedSize()
	dSpeed := g.ExpressedSpeed() - other.ExpressedSpeed()
	dSense := (g.ExpressedSense() - other.ExpressedSense()) * 0.01
	dDiet := g.ExpressedDiet() - other.ExpressedDiet()
	dMeta := g.ExpressedMetabolism() - other.ExpressedMetabolism()
	dFert := g.ExpressedFertility() - other.ExpressedFertility()
	dConst := g.ExpressedConstitution() - other.ExpressedConstitution()
	dHidden := (g.ExpressedHidden() - other.ExpressedHidden()) * 0.1

	sumSq := dSize*dSize +
		dSpeed*dSpeed +
		dSense*dSense +
		dDiet*dDiet +
		dMeta*dMeta +
		dFert*dFert +
		dConst*dConst +
		dHidden*dHidden

	return math.Sqrt(sumSq)
}
