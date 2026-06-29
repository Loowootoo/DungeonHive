package config

// for a*pathfinding
const (
	GridSize = 60
	MaxCells = GridSize * GridSize
	Inf      = 1 << 30
)

const (
	TileSize     = 16
	ScreenWidth  = GridSize * TileSize
	ScreenHeight = GridSize * TileSize
)
const (
	TileSetCol = 13
	TileSetRow = 5
)

// common types
type Point struct {
	X, Y int
}

type TileType int
type DirType int

const (
	DirUP    DirType = 0
	DirDOWN  DirType = 1
	DirLEFT  DirType = 2
	DirRIGHT DirType = 3
)
