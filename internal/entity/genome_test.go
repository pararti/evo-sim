package entity

import (
	"math"
	"testing"
)

func TestGenome_CalculateStats(t *testing.T) {
	// 1. Small creature
	g1 := Genome{
		SizeGene:  0.5,
		SpeedGene: 1.0,
		SenseGene: 100.0,
		DietGene:  0.2, // Herbivore
	}
	mass1, speed1, _, bmr1, isCarn1 := g1.CalculateStats()

	if isCarn1 {
		t.Errorf("Expected herbivore, got carnivore")
	}

	// 2. Large creature (Size 2.0 vs 0.5 -> Mass 4x -> 16x?)
	// Formula: mass = size * size. 0.5^2 = 0.25. 2.0^2 = 4.0. Mass is 16x.
	g2 := Genome{
		SizeGene:  2.0,
		SpeedGene: 1.0,
		SenseGene: 100.0,
		DietGene:  0.8, // Carnivore
	}
	mass2, speed2, _, bmr2, isCarn2 := g2.CalculateStats()

	if !isCarn2 {
		t.Errorf("Expected carnivore, got herbivore")
	}

	if mass2 <= mass1 {
		t.Errorf("Expected larger creature to have more mass")
	}

	// Speed = SpeedGene / Sqrt(Size).
	// g1: 1.0 / sqrt(0.5) = 1.41
	// g2: 1.0 / sqrt(2.0) = 0.707
	if speed2 >= speed1 {
		t.Errorf("Expected larger creature to be slower with same muscles")
	}

	// BMR check: Larger mass should have higher BMR
	if bmr2 <= bmr1 {
		t.Errorf("Expected larger creature to have higher BMR")
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
			math.Abs(mutated.DietGene-g.DietGene) > 0.001 {
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
		ColorR: 0.0, ColorG: 0.0, ColorB: 0.0,
	}
	g2 := Genome{
		SizeGene: 3.5, SpeedGene: 2.8, SenseGene: 450.0, DietGene: 0.9,
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
