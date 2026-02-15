package entity

import (
	"math"
	"testing"

	"evo-sim/internal/brain"
)

func TestCreature_ReproduceSexual(t *testing.T) {
	// Create two parent creatures with known positions and energy
	p1 := &Creature{
		ID:                    1,
		X:                     100.0,
		Y:                     200.0,
		Energy:                300.0,
		Size:                  1.0,
		Mass:                  1.0,
		Speed:                 1.0,
		ViewRadius:            100.0,
		BMR:                   0.1,
		MaxEnergy:             1000.0,
		ReproductionThreshold: 500.0,
		IsCarnivore:           false,
		Genome: Genome{
			SizeGene: 1.0, SpeedGene: 1.0, SenseGene: 100.0, DietGene: 0.2,
			MetabolismGene: 1.0, FertilityGene: 0.5, ConstitutionGene: 1.0,
			ColorR: 0.0, ColorG: 1.0, ColorB: 0.0,
		},
		Brain: brain.NewNetwork(10, 6, 2),
	}
	p2 := &Creature{
		ID:                    2,
		X:                     120.0,
		Y:                     220.0,
		Energy:                240.0,
		Size:                  1.5,
		Mass:                  2.25,
		Speed:                 0.8,
		ViewRadius:            120.0,
		BMR:                   0.15,
		MaxEnergy:             2000.0,
		ReproductionThreshold: 1000.0,
		IsCarnivore:           false,
		Genome: Genome{
			SizeGene: 1.5, SpeedGene: 0.8, SenseGene: 120.0, DietGene: 0.3,
			MetabolismGene: 1.0, FertilityGene: 0.5, ConstitutionGene: 1.0,
			ColorR: 0.0, ColorG: 0.5, ColorB: 0.5,
		},
		Brain: brain.NewNetwork(10, 6, 2),
	}

	p1EnergyBefore := p1.Energy
	p2EnergyBefore := p2.Energy

	child := p1.ReproduceSexual(p2, 0.1, 0.2)

	// Each parent loses 1/3 of their energy
	expectedP1Loss := p1EnergyBefore / 3
	expectedP2Loss := p2EnergyBefore / 3

	if math.Abs(p1.Energy-(p1EnergyBefore-expectedP1Loss)) > 0.001 {
		t.Errorf("Parent1 energy: got %f, want %f", p1.Energy, p1EnergyBefore-expectedP1Loss)
	}
	if math.Abs(p2.Energy-(p2EnergyBefore-expectedP2Loss)) > 0.001 {
		t.Errorf("Parent2 energy: got %f, want %f", p2.Energy, p2EnergyBefore-expectedP2Loss)
	}

	// Child energy = sum of both contributions
	expectedChildEnergy := expectedP1Loss + expectedP2Loss
	if math.Abs(child.Energy-expectedChildEnergy) > 0.001 {
		t.Errorf("Child energy: got %f, want %f", child.Energy, expectedChildEnergy)
	}

	// Child position is midpoint of parents
	expectedX := (p1.X + p2.X) / 2 // Note: p1.X was not changed by ReproduceSexual
	expectedY := (p1.Y + p2.Y) / 2
	// p1.X=100, p2.X=120 -> midpoint was calculated before energy deduction
	// Actually the positions don't change, only energy does
	if math.Abs(child.X-110.0) > 0.001 || math.Abs(child.Y-210.0) > 0.001 {
		t.Errorf("Child position: got (%f, %f), want (%f, %f)", child.X, child.Y, expectedX, expectedY)
	}

	// Child has valid genome and brain
	if child.Brain == nil {
		t.Error("Child brain is nil")
	}
	if child.Age != 0 {
		t.Errorf("Child age: got %d, want 0", child.Age)
	}
	
	// Check that child has Mass and other new fields populated
	if child.Mass <= 0 {
		t.Errorf("Child mass should be positive, got %f", child.Mass)
	}
	if child.MaxEnergy <= 0 {
		t.Errorf("Child MaxEnergy should be positive, got %f", child.MaxEnergy)
	}
}
