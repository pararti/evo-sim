package world

import (
	"log"
	"math"
	"math/rand/v2"
)

type TerrainType uint8

const (
	Water TerrainType = iota
	Sand
	Grass
)

// TerrainGrid holds the static map data.
// Resolution: 1 cell represents Scale x Scale world units.
type TerrainGrid struct {
	Width, Height int       // Dimensions in grid cells
	Scale         float64   // World units per cell (e.g., 20.0)
	Cells         []TerrainType
}

func NewTerrainGrid(worldW, worldH, scale float64) *TerrainGrid {
	w := int(math.Ceil(worldW / scale))
	h := int(math.Ceil(worldH / scale))

	t := &TerrainGrid{
		Width:  w,
		Height: h,
		Scale:  scale,
		Cells:  make([]TerrainType, w*h),
	}

	t.Generate()
	return t
}

func (t *TerrainGrid) Generate() {
	// Simple Perlin-like noise using overlapping sine waves
	seed := rand.Float64() * 100
	
	// Increased frequencies to fit more features into small grid (40x30)
	freq1 := 0.25 
	freq2 := 0.8
	
	counts := map[TerrainType]int{Water: 0, Sand: 0, Grass: 0}
	total := t.Width * t.Height

	for y := 0; y < t.Height; y++ {
		for x := 0; x < t.Width; x++ {
			nx := float64(x)
			ny := float64(y)
			
			// Combine waves
			val := math.Sin(nx*freq1 + seed) * math.Cos(ny*freq1 + seed)
			val += (math.Sin(nx*freq2 - seed) + math.Cos(ny*freq2 + seed*0.5)) * 0.2
			
			idx := y*t.Width + x
			if val < -0.2 {
				t.Cells[idx] = Water
				counts[Water]++
			} else if val < 0.2 {
				t.Cells[idx] = Sand
				counts[Sand]++
			} else {
				t.Cells[idx] = Grass
				counts[Grass]++
			}
		}
	}
	
	log.Printf("Terrain Generated: Water %.1f%%, Sand %.1f%%, Grass %.1f%%",
		float64(counts[Water])/float64(total)*100,
		float64(counts[Sand])/float64(total)*100,
		float64(counts[Grass])/float64(total)*100,
	)
}

func (t *TerrainGrid) GetType(worldX, worldY float64) TerrainType {
	gx := int(worldX / t.Scale)
	gy := int(worldY / t.Scale)

	// Clamp
	if gx < 0 { gx = 0 }
	if gx >= t.Width { gx = t.Width - 1 }
	if gy < 0 { gy = 0 }
	if gy >= t.Height { gy = t.Height - 1 }

	return t.Cells[gy*t.Width + gx]
}

func (t *TerrainGrid) GetMovementPenalty(worldX, worldY float64) (speedFactor, energyCostFactor float64) {
	tt := t.GetType(worldX, worldY)
	switch tt {
	case Water:
		return 0.3, 3.0 // Very slow, very tiring
	case Sand:
		return 0.6, 1.5 // Slow, tiring
	case Grass:
		return 1.0, 1.0 // Normal
	default:
		return 1.0, 1.0
	}
}
