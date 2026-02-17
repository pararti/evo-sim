package entity

import (
	"math"
	"testing"

	"evo-sim/internal/brain"
)

func TestCreature_ReproduceSexual(t *testing.T) {
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
			SizeAllele1: 1.0, SizeAllele2: 1.0,
			SpeedAllele1: 1.0, SpeedAllele2: 1.0,
			SenseAllele1: 100.0, SenseAllele2: 100.0,
			DietAllele1: 0.2, DietAllele2: 0.2,
			MetabolismAllele1: 1.0, MetabolismAllele2: 1.0,
			FertilityAllele1: 0.5, FertilityAllele2: 0.5,
			ConstitutionAllele1: 1.0, ConstitutionAllele2: 1.0,
			HiddenAllele1: 6.0, HiddenAllele2: 6.0,
			ColorR: 0.0, ColorG: 1.0, ColorB: 0.0,
		},
		Brain: brain.NewNetwork(11, 6, 2),
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
			SizeAllele1: 1.5, SizeAllele2: 1.5,
			SpeedAllele1: 0.8, SpeedAllele2: 0.8,
			SenseAllele1: 120.0, SenseAllele2: 120.0,
			DietAllele1: 0.3, DietAllele2: 0.3,
			MetabolismAllele1: 1.0, MetabolismAllele2: 1.0,
			FertilityAllele1: 0.5, FertilityAllele2: 0.5,
			ConstitutionAllele1: 1.0, ConstitutionAllele2: 1.0,
			HiddenAllele1: 6.0, HiddenAllele2: 6.0,
			ColorR: 0.0, ColorG: 0.5, ColorB: 0.5,
		},
		Brain: brain.NewNetwork(11, 6, 2),
	}

	p1EnergyBefore := p1.Energy
	p2EnergyBefore := p2.Energy

	child := p1.ReproduceSexual(p2, 0.1, 0.2, 0.15, 0.2, 0.005)

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
	expectedX := (p1.X + p2.X) / 2
	expectedY := (p1.Y + p2.Y) / 2
	if math.Abs(child.X-110.0) > 0.001 || math.Abs(child.Y-210.0) > 0.001 {
		t.Errorf("Child position: got (%f, %f), want (%f, %f)", child.X, child.Y, expectedX, expectedY)
	}

	if child.Brain == nil {
		t.Error("Child brain is nil")
	}
	if child.Age != 0 {
		t.Errorf("Child age: got %d, want 0", child.Age)
	}

	if child.Mass <= 0 {
		t.Errorf("Child mass should be positive, got %f", child.Mass)
	}
	if child.MaxEnergy <= 0 {
		t.Errorf("Child MaxEnergy should be positive, got %f", child.MaxEnergy)
	}
}
