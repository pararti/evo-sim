package world

import (
	"evo-sim/internal/entity"
	"sync"
)

type Species struct {
	ID       int
	Centroid entity.Genome
	Count    int
}

type SpeciesManager struct {
	NextID    int
	Species   map[int]*Species
	Threshold float64
	Mu        sync.RWMutex
}

func NewSpeciesManager(threshold float64) *SpeciesManager {
	return &SpeciesManager{
		NextID:    1,
		Species:   make(map[int]*Species),
		Threshold: threshold,
	}
}

// Classify determines the species of a genome.
// Returns the SpeciesID.
// If the genome is close enough to an existing species, it joins it.
// Otherwise, a new species is created.
func (sm *SpeciesManager) Classify(g entity.Genome) int {
	sm.Mu.Lock()
	defer sm.Mu.Unlock()

	var bestSpecies *Species
	bestDist := sm.Threshold // Start with threshold as max allowed

	for _, s := range sm.Species {
		dist := g.Distance(s.Centroid)
		if dist < bestDist {
			bestDist = dist
			bestSpecies = s
		}
	}

	if bestSpecies != nil {
		bestSpecies.Count++
		return bestSpecies.ID
	}

	// New Species
	newID := sm.NextID
	sm.NextID++
	sm.Species[newID] = &Species{
		ID:       newID,
		Centroid: g, // The founder defines the species
		Count:    1,
	}
	return newID
}

// Register adds a creature to a known species (e.g. initial loading or forced assignment)
// If species doesn't exist, it creates it.
// Used primarily when we want to increment count for an existing ID.
func (sm *SpeciesManager) Register(speciesID int, g entity.Genome) {
	sm.Mu.Lock()
	defer sm.Mu.Unlock()

	if s, ok := sm.Species[speciesID]; ok {
		s.Count++
	} else {
		// Should not happen usually, but for safety
		sm.Species[speciesID] = &Species{
			ID:       speciesID,
			Centroid: g,
			Count:    1,
		}
		if speciesID >= sm.NextID {
			sm.NextID = speciesID + 1
		}
	}
}

func (sm *SpeciesManager) RemoveCreature(speciesID int) {
	sm.Mu.Lock()
	defer sm.Mu.Unlock()

	if s, ok := sm.Species[speciesID]; ok {
		s.Count--
		if s.Count <= 0 {
			delete(sm.Species, speciesID)
		}
	}
}

func (sm *SpeciesManager) GetSpeciesCount() int {
	sm.Mu.RLock()
	defer sm.Mu.RUnlock()
	return len(sm.Species)
}

func (sm *SpeciesManager) Clear() {
	sm.Mu.Lock()
	defer sm.Mu.Unlock()
	sm.Species = make(map[int]*Species)
	sm.NextID = 1
}
