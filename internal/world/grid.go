package world

import (
	"math"
	"evo-sim/internal/entity"
)

type Cell struct {
	Creatures []*entity.Creature
	Food      []entity.Food
}

type Grid struct {
	cellSize float64
	cols     int
	rows     int
	cells    []Cell
}

func NewGrid(width, height, cellSize float64) *Grid {
	cols := int(math.Ceil(width / cellSize))
	rows := int(math.Ceil(height / cellSize))

	g := &Grid{
		cellSize: cellSize,
		cols:     cols,
		rows:     rows,
		cells:    make([]Cell, cols*rows),
	}
	
	// Pre-allocate inner slices to avoid initial resizing churn
	for i := range g.cells {
		g.cells[i].Creatures = make([]*entity.Creature, 0, 16)
		g.cells[i].Food = make([]entity.Food, 0, 16)
	}
	
	return g
}

func (g *Grid) Clear() {
	// Reset slices without deallocating memory
	for i := range g.cells {
		g.cells[i].Creatures = g.cells[i].Creatures[:0]
		g.cells[i].Food = g.cells[i].Food[:0]
	}
}

func (g *Grid) InsertCreature(c *entity.Creature) {
	col := int(c.X / g.cellSize)
	row := int(c.Y / g.cellSize)

	if col >= 0 && col < g.cols && row >= 0 && row < g.rows {
		index := row*g.cols + col
		g.cells[index].Creatures = append(g.cells[index].Creatures, c)
	}
}

func (g *Grid) InsertFood(f entity.Food) {
	col := int(f.X / g.cellSize)
	row := int(f.Y / g.cellSize)

	if col >= 0 && col < g.cols && row >= 0 && row < g.rows {
		index := row*g.cols + col
		g.cells[index].Food = append(g.cells[index].Food, f)
	}
}

// GetNeighbors returns flattened lists of entities in the radius
func (g *Grid) GetNeighbors(x, y float64, radius float64) ([][]*entity.Creature, [][]entity.Food) {
	colStart := int((x - radius) / g.cellSize)
	colEnd := int((x + radius) / g.cellSize)
	rowStart := int((y - radius) / g.cellSize)
	rowEnd := int((y + radius) / g.cellSize)

	// Clamp boundaries
	if colStart < 0 { colStart = 0 }
	if colEnd >= g.cols { colEnd = g.cols - 1 }
	if rowStart < 0 { rowStart = 0 }
	if rowEnd >= g.rows { rowEnd = g.rows - 1 }

	// Optimization: We could pre-allocate these result buffers too if needed, 
	// but for now, let's keep the signature. 
	// The return type is slightly awkward for 1D grid ([][]...), 
	// but changing it requires refactoring Engine too much. 
	// We will return a slice of slices to maintain API compatibility with Engine.
	
	var cResults [][]*entity.Creature
	var fResults [][]entity.Food

	for r := rowStart; r <= rowEnd; r++ {
		for c := colStart; c <= colEnd; c++ {
			index := r*g.cols + c
			cell := &g.cells[index]
			
			if len(cell.Creatures) > 0 {
				cResults = append(cResults, cell.Creatures)
			}
			if len(cell.Food) > 0 {
				fResults = append(fResults, cell.Food)
			}
		}
	}

	return cResults, fResults
}
