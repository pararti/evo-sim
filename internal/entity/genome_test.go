package entity

import (
	"math"
	"testing"
)

func TestGenome_CalculateStats(t *testing.T) {
	// 1. Small herbivore
	g1 := Genome{
		SizeAllele1: 0.5, SizeAllele2: 0.5,
		SpeedAllele1: 1.0, SpeedAllele2: 1.0,
		SenseAllele1: 100.0, SenseAllele2: 100.0,
		DietAllele1: 0.2, DietAllele2: 0.2,
		MetabolismAllele1: 1.0, MetabolismAllele2: 1.0,
		FertilityAllele1: 0.7, FertilityAllele2: 0.7,
		ConstitutionAllele1: 1.0, ConstitutionAllele2: 1.0,
		HiddenAllele1: 6.0, HiddenAllele2: 6.0,
	}
	mass1, speed1, _, bmr1, _, _, isCarn1, _ := g1.CalculateStats(0.005)

	if isCarn1 {
		t.Errorf("Expected herbivore, got carnivore")
	}

	// 2. Large carnivore
	g2 := Genome{
		SizeAllele1: 2.0, SizeAllele2: 2.0,
		SpeedAllele1: 1.0, SpeedAllele2: 1.0,
		SenseAllele1: 100.0, SenseAllele2: 100.0,
		DietAllele1: 0.8, DietAllele2: 0.8,
		MetabolismAllele1: 1.0, MetabolismAllele2: 1.0,
		FertilityAllele1: 0.7, FertilityAllele2: 0.7,
		ConstitutionAllele1: 1.0, ConstitutionAllele2: 1.0,
		HiddenAllele1: 6.0, HiddenAllele2: 6.0,
	}
	mass2, speed2, _, bmr2, _, _, isCarn2, _ := g2.CalculateStats(0.005)

	if !isCarn2 {
		t.Errorf("Expected carnivore, got herbivore")
	}

	if mass2 <= mass1 {
		t.Errorf("Expected larger creature to have more mass")
	}

	if speed2 >= speed1 {
		t.Errorf("Expected larger creature to be slower with same muscles. S1: %f, S2: %f", speed1, speed2)
	}

	if bmr2 <= bmr1 {
		t.Errorf("Expected larger creature to have higher BMR")
	}
}

func TestGenome_Epistasis(t *testing.T) {
	gBase := Genome{
		SizeAllele1: 1.0, SizeAllele2: 1.0,
		SpeedAllele1: 1.0, SpeedAllele2: 1.0,
		SenseAllele1: 100.0, SenseAllele2: 100.0,
		DietAllele1: 0.5, DietAllele2: 0.5,
		MetabolismAllele1: 1.0, MetabolismAllele2: 1.0,
		FertilityAllele1: 0.7, FertilityAllele2: 0.7,
		ConstitutionAllele1: 1.0, ConstitutionAllele2: 1.0,
		HiddenAllele1: 6.0, HiddenAllele2: 6.0,
	}

	massBase, _, _, _, maxEnergyBase, _, _, _ := gBase.CalculateStats(0.005)

	gDense := gBase
	gDense.ConstitutionAllele1 = 1.5
	gDense.ConstitutionAllele2 = 1.5

	massDense, _, _, _, maxEnergyDense, _, _, _ := gDense.CalculateStats(0.005)

	if massDense <= massBase {
		t.Errorf("High constitution should increase mass")
	}
	if maxEnergyDense <= maxEnergyBase {
		t.Errorf("High constitution should increase max energy storage")
	}

	// Test that Metabolism affects Speed and BMR
	gFastMeta := gBase
	gFastMeta.MetabolismAllele1 = 1.5
	gFastMeta.MetabolismAllele2 = 1.5

	_, speedFast, _, bmrFast, _, _, _, _ := gFastMeta.CalculateStats(0.005)
	_, speedBase, _, bmrBase, _, _, _, _ := gBase.CalculateStats(0.005)

	if speedFast <= speedBase {
		t.Errorf("High metabolism should increase speed (assuming constant mass)")
	}
	if bmrFast <= bmrBase {
		t.Errorf("High metabolism should increase BMR")
	}
}

func TestGenome_Mutate(t *testing.T) {
	g := NewRandomGenome()

	mutated := g
	changed := false

	for i := 0; i < 100; i++ {
		mutated = mutated.Mutate(0.5, 0.1)

		if math.Abs(mutated.SizeAllele1-g.SizeAllele1) > 0.001 ||
			math.Abs(mutated.SpeedAllele1-g.SpeedAllele1) > 0.001 ||
			math.Abs(mutated.DietAllele1-g.DietAllele1) > 0.001 ||
			math.Abs(mutated.MetabolismAllele1-g.MetabolismAllele1) > 0.001 {
			changed = true
			break
		}
	}

	if !changed {
		t.Errorf("Mutation failed to change genome after 100 attempts")
	}
}

func TestGenome_Crossover(t *testing.T) {
	// Two maximally different genomes (homozygous for distinct values)
	g1 := Genome{
		SizeAllele1: 0.5, SizeAllele2: 0.5,
		SpeedAllele1: 0.3, SpeedAllele2: 0.3,
		SenseAllele1: 50.0, SenseAllele2: 50.0,
		DietAllele1: 0.1, DietAllele2: 0.1,
		MetabolismAllele1: 0.5, MetabolismAllele2: 0.5,
		FertilityAllele1: 0.5, FertilityAllele2: 0.5,
		ConstitutionAllele1: 0.5, ConstitutionAllele2: 0.5,
		HiddenAllele1: 4.0, HiddenAllele2: 4.0,
		ColorR: 0.0, ColorG: 0.0, ColorB: 0.0,
	}
	g2 := Genome{
		SizeAllele1: 3.5, SizeAllele2: 3.5,
		SpeedAllele1: 2.8, SpeedAllele2: 2.8,
		SenseAllele1: 450.0, SenseAllele2: 450.0,
		DietAllele1: 0.9, DietAllele2: 0.9,
		MetabolismAllele1: 2.0, MetabolismAllele2: 2.0,
		FertilityAllele1: 0.9, FertilityAllele2: 0.9,
		ConstitutionAllele1: 1.5, ConstitutionAllele2: 1.5,
		HiddenAllele1: 10.0, HiddenAllele2: 10.0,
		ColorR: 1.0, ColorG: 1.0, ColorB: 1.0,
	}

	// Run crossover many times to verify mixing
	sawG1Size := false
	sawG2Size := false
	for i := 0; i < 100; i++ {
		child := g1.Crossover(g2)

		// Child allele1 comes from g1, allele2 from g2 (both homozygous)
		if child.SizeAllele1 != g1.SizeAllele1 && child.SizeAllele1 != g1.SizeAllele2 {
			t.Fatalf("Child SizeAllele1 %f is not from parent1", child.SizeAllele1)
		}
		if child.SizeAllele2 != g2.SizeAllele1 && child.SizeAllele2 != g2.SizeAllele2 {
			t.Fatalf("Child SizeAllele2 %f is not from parent2", child.SizeAllele2)
		}

		if child.SizeAllele1 == g1.SizeAllele1 {
			sawG1Size = true
		}
		if child.SizeAllele2 == g2.SizeAllele1 {
			sawG2Size = true
		}
	}

	if !sawG1Size || !sawG2Size {
		t.Errorf("Crossover not working: sawG1=%v sawG2=%v", sawG1Size, sawG2Size)
	}
}

func TestGenome_Diploid_Dominance(t *testing.T) {
	// Size uses dominant (max) expression
	g := Genome{
		SizeAllele1: 1.0, SizeAllele2: 2.0,
		SpeedAllele1: 0.5, SpeedAllele2: 1.5,
		SenseAllele1: 100.0, SenseAllele2: 200.0,
		DietAllele1: 0.3, DietAllele2: 0.7,
		MetabolismAllele1: 1.0, MetabolismAllele2: 1.0,
		FertilityAllele1: 0.7, FertilityAllele2: 0.7,
		ConstitutionAllele1: 1.0, ConstitutionAllele2: 1.0,
		HiddenAllele1: 6.0, HiddenAllele2: 6.0,
	}

	// Dominant traits: max
	if g.ExpressedSize() != 2.0 {
		t.Errorf("Size should be dominant (max): got %f, want 2.0", g.ExpressedSize())
	}
	if g.ExpressedSpeed() != 1.5 {
		t.Errorf("Speed should be dominant (max): got %f, want 1.5", g.ExpressedSpeed())
	}

	// Additive traits: avg
	if g.ExpressedSense() != 150.0 {
		t.Errorf("Sense should be additive (avg): got %f, want 150.0", g.ExpressedSense())
	}

	// Diet is dominant (max) — carnivory dominates
	if g.ExpressedDiet() != 0.7 {
		t.Errorf("Diet should be dominant (max): got %f, want 0.7", g.ExpressedDiet())
	}
}
