package renderer

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// DrawPlayer 以鏡頭轉換將 sprite 繪製到世界座標 (worldX, worldY) 中心。
func DrawPlayer(screen *ebiten.Image, sprite *ebiten.Image, worldX, worldY float64, cam *Camera) {
	op := &ebiten.DrawImageOptions{}
	sw := float64(screen.Bounds().Dx())
	sh := float64(screen.Bounds().Dy())

	// sprite 中心對齊世界座標 → 鏡頭位移 → 縮放 → 畫面置中
	op.GeoM.Translate(-float64(sprite.Bounds().Dx())/2, -float64(sprite.Bounds().Dy())/2)
	op.GeoM.Translate(worldX, worldY)
	op.GeoM.Translate(-cam.X, -cam.Y)
	op.GeoM.Scale(cam.Zoom, cam.Zoom)
	op.GeoM.Translate(sw/2, sh/2)

	screen.DrawImage(sprite, op)
}
