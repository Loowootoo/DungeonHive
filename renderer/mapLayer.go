package renderer

import (
	cf "DungeonHive/config"

	"github.com/hajimehoshi/ebiten/v2"
)

func (r *Renderer) DrawMapLayer(m [cf.GridSize][cf.GridSize]int) {
	op := &ebiten.DrawImageOptions{}
	for y := 0; y < cf.GridSize; y++ {
		for x := 0; x < cf.GridSize; x++ {
			tilesID := m[y][x]
			tx := tilesID % cf.TileSetCol
			ty := tilesID / cf.TileSetCol
			dX, dY := float64(x*cf.TileSize), float64(y*cf.TileSize)
			op.GeoM.Reset()
			op.GeoM.Translate(dX, dY)
			r.MapLayer.DrawImage(r.ResMan.TileImage.SubImage(GetSubImageRect(tx, ty)).(*ebiten.Image), op)
		}
	}
}
