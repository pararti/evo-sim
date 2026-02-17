package entity

type Food struct {
	ID         int
	X, Y       float64
	Energy     float64 // >0 for carrion (variable energy); 0 for regular plant food
	DecayTicks int     // >0 for carrion (countdown each tick); 0 for regular food
}
