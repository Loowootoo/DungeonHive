package renderer

import (
	cf "DungeonHive/pkg/config"
	"image"
)

func GetSubImageRect(x, y int) image.Rectangle {
	return image.Rect(x*cf.TileSize, y*cf.TileSize, (x+1)*cf.TileSize, (y+1)*cf.TileSize)
}
