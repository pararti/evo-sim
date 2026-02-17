package world

import "math"

// PheromoneGrid tracks pheromone concentrations across the world.
// Same resolution as TerrainGrid (scale=20).
type PheromoneGrid struct {
	Width, Height int
	Scale         float64
	Cells         []float64
}

func NewPheromoneGrid(worldW, worldH, scale float64) *PheromoneGrid {
	w := int(math.Ceil(worldW / scale))
	h := int(math.Ceil(worldH / scale))

	return &PheromoneGrid{
		Width:  w,
		Height: h,
		Scale:  scale,
		Cells:  make([]float64, w*h),
	}
}

func (p *PheromoneGrid) cellIndex(worldX, worldY float64) int {
	gx := int(worldX / p.Scale)
	gy := int(worldY / p.Scale)

	if gx < 0 {
		gx = 0
	}
	if gx >= p.Width {
		gx = p.Width - 1
	}
	if gy < 0 {
		gy = 0
	}
	if gy >= p.Height {
		gy = p.Height - 1
	}

	return gy*p.Width + gx
}

// Deposit adds pheromone at the given world position.
func (p *PheromoneGrid) Deposit(worldX, worldY, amount float64) {
	idx := p.cellIndex(worldX, worldY)
	p.Cells[idx] += amount
	if p.Cells[idx] > 10.0 {
		p.Cells[idx] = 10.0 // Cap to prevent unbounded accumulation
	}
}

// Get returns the pheromone concentration at the given world position.
func (p *PheromoneGrid) Get(worldX, worldY float64) float64 {
	return p.Cells[p.cellIndex(worldX, worldY)]
}

// Decay multiplies all cells by the decay factor.
func (p *PheromoneGrid) Decay(factor float64) {
	for i := range p.Cells {
		p.Cells[i] *= factor
		if p.Cells[i] < 0.001 {
			p.Cells[i] = 0 // Zero out negligible values
		}
	}
}
