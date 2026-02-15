package entity

import (
	"math"
	"testing"
)

func TestGenome_CalculateStats(t *testing.T) {
	// 1. Small creature
	g1 := Genome{
		SizeGene:         0.5,
		SpeedGene:        1.0,
		SenseGene:        100.0,
		DietGene:         0.2, // Herbivore
		MetabolismGene:   1.0,
		FertilityGene:    0.7,
		ConstitutionGene: 1.0,
	}
	mass1, speed1, _, bmr1, _, _, isCarn1 := g1.CalculateStats()

	if isCarn1 {
		t.Errorf("Expected herbivore, got carnivore")
	}

	// 2. Large creature
	g2 := Genome{
		SizeGene:         2.0,
		SpeedGene:        1.0,
		SenseGene:        100.0,
		DietGene:         0.8, // Carnivore
		MetabolismGene:   1.0,
		FertilityGene:    0.7,
		ConstitutionGene: 1.0,
	}
	mass2, speed2, _, bmr2, _, _, isCarn2 := g2.CalculateStats()

	if !isCarn2 {
		t.Errorf("Expected carnivore, got herbivore")
	}

	if mass2 <= mass1 {
		t.Errorf("Expected larger creature to have more mass")
	}

	// Speed = SpeedGene / Sqrt(Mass).
	// g1: Mass = 0.25 * 1.0 = 0.25. Sqrt = 0.5. Speed = 1.0 / 0.5 = 2.0
	// g2: Mass = 4.0 * 1.0 = 4.0. Sqrt = 2.0. Speed = 1.0 / 2.0 = 0.5
	if speed2 >= speed1 {
		t.Errorf("Expected larger creature to be slower with same muscles. S1: %f, S2: %f", speed1, speed2)
	}

	// BMR check: Larger mass should have higher BMR
	if bmr2 <= bmr1 {
		t.Errorf("Expected larger creature to have higher BMR")
	}
}

func TestGenome_Epistasis(t *testing.T) {
	// Test that Constitution affects Mass and MaxEnergy
	gBase := Genome{
		SizeGene:         1.0,
		SpeedGene:        1.0,
		SenseGene:        100.0,
		DietGene:         0.5,
		MetabolismGene:   1.0,
		FertilityGene:    0.7,
		ConstitutionGene: 1.0,
	}

	massBase, _, _, _, maxEnergyBase, _, _ := gBase.CalculateStats()

	gDense := gBase
	gDense.ConstitutionGene = 1.5

	massDense, _, _, _, maxEnergyDense, _, _ := gDense.CalculateStats()

	if massDense <= massBase {
		t.Errorf("High constitution should increase mass")
	}
	if maxEnergyDense <= maxEnergyBase {
		t.Errorf("High constitution should increase max energy storage")
	}

	// Test that Metabolism affects Speed and BMR
	gFastMeta := gBase
	gFastMeta.MetabolismGene = 1.5

	_, speedFast, _, bmrFast, _, _, _ := gFastMeta.CalculateStats()
	_, speedBase, _, bmrBase, _, _, _ := gBase.CalculateStats()

	if speedFast <= speedBase {
		t.Errorf("High metabolism should increase speed (assuming constant mass)")
	}
	if bmrFast <= bmrBase {
		t.Errorf("High metabolism should increase BMR")
	}
}

func TestGenome_Mutate(t *testing.T) {
	g := NewRandomGenome()

	// Mutate many times to ensure *some* change happens
	mutated := g
	changed := false

	for i := 0; i < 100; i++ {
		mutated = mutated.Mutate(0.5, 0.1) // 50% chance, 0.1 strength

		if math.Abs(mutated.SizeGene-g.SizeGene) > 0.001 ||
			math.Abs(mutated.SpeedGene-g.SpeedGene) > 0.001 ||
			math.Abs(mutated.DietGene-g.DietGene) > 0.001 ||
			math.Abs(mutated.MetabolismGene-g.MetabolismGene) > 0.001 {
			changed = true
			break
		}
	}

	if !changed {
		t.Errorf("Mutation failed to change genome after 100 attempts")
	}
}

func TestGenome_Crossover(t *testing.T) {
	// Two maximally different genomes
	g1 := Genome{
		SizeGene: 0.5, SpeedGene: 0.3, SenseGene: 50.0, DietGene: 0.1,
		MetabolismGene: 0.5, FertilityGene: 0.5, ConstitutionGene: 0.5,
		ColorR: 0.0, ColorG: 0.0, ColorB: 0.0,
	}
	g2 := Genome{
		SizeGene: 3.5, SpeedGene: 2.8, SenseGene: 450.0, DietGene: 0.9,
		MetabolismGene: 2.0, FertilityGene: 0.9, ConstitutionGene: 1.5,
		ColorR: 1.0, ColorG: 1.0, ColorB: 1.0,
	}

	// Run crossover many times to verify mixing
	sawG1Size := false
	sawG2Size := false
	for i := 0; i < 100; i++ {
		child := g1.Crossover(g2)

		// Child gene must be exactly one of the two parents
		if child.SizeGene != g1.SizeGene && child.SizeGene != g2.SizeGene {
			t.Fatalf("Child SizeGene %f is neither parent1 %f nor parent2 %f",
				child.SizeGene, g1.SizeGene, g2.SizeGene)
		}
		if child.SizeGene == g1.SizeGene {
			sawG1Size = true
		}
		if child.SizeGene == g2.SizeGene {
			sawG2Size = true
		}
	}

	if !sawG1Size || !sawG2Size {
		t.Errorf("Crossover not mixing: sawG1=%v sawG2=%v", sawG1Size, sawG2Size)
	}
}
