package world

import (
	"math"
	"evo-sim/internal/entity"
)

type Grid struct {
	cellSize float64
	cols, rows int
	creatureCells [][][]*entity.Creature
	foodCells     [][][]entity.Food
}

func NewGrid(width, height, cellSize float64) *Grid {
	cols := int(math.Ceil(width / cellSize))
	rows := int(math.Ceil(height / cellSize))

	g := &Grid{
		cellSize: cellSize,
		cols:     cols,
		rows:     rows,
	}
	g.Clear()
	return g
}

func (g *Grid) Clear() {
	g.creatureCells = make([][][]*entity.Creature, g.cols)
	for i := range g.creatureCells {
		g.creatureCells[i] = make([][]*entity.Creature, g.rows)
	}

	g.foodCells = make([][][]entity.Food, g.cols)
	for i := range g.foodCells {
		g.foodCells[i] = make([][]entity.Food, g.rows)
	}
}

func (g *Grid) InsertCreature(c *entity.Creature) {
	col := int(c.X / g.cellSize)
	row := int(c.Y / g.cellSize)

	if col >= 0 && col < g.cols && row >= 0 && row < g.rows {
		g.creatureCells[col][row] = append(g.creatureCells[col][row], c)
	}
}

func (g *Grid) InsertFood(f entity.Food) {
	col := int(f.X / g.cellSize)
	row := int(f.Y / g.cellSize)

	if col >= 0 && col < g.cols && row >= 0 && row < g.rows {
		g.foodCells[col][row] = append(g.foodCells[col][row], f)
	}
}

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

	var cResults [][]*entity.Creature
	var fResults [][]entity.Food

	for c := colStart; c <= colEnd; c++ {
		for r := rowStart; r <= rowEnd; r++ {
			if len(g.creatureCells[c][r]) > 0 {
				cResults = append(cResults, g.creatureCells[c][r])
			}
			if len(g.foodCells[c][r]) > 0 {
				fResults = append(fResults, g.foodCells[c][r])
			}
		}
	}

	return cResults, fResults
}
