package renderer

import (
	. "DungeonHive/config"

	"github.com/hajimehoshi/ebiten/v2"
)

// Camera 鏡頭系統
type Camera struct {
	X, Y float64 // 世界座標中心
	Zoom float64
}

func NewCamera(x, y, zoom float64) *Camera {
	return &Camera{
		X:    x,
		Y:    y,
		Zoom: zoom,
	}
}

func (c *Camera) Move(dx, dy float64) {
	c.X += dx / c.Zoom
	c.Y += dy / c.Zoom
}

func (c *Camera) ZoomAt(screenX, screenY, factor float64) {
	// 以鼠標位置為中心縮放
	halfW, halfH := float64(ScreenWidth)/2, float64(ScreenHeight)/2
	worldX := (screenX-halfW)/c.Zoom + c.X
	worldY := (screenY-halfH)/c.Zoom + c.Y

	c.Zoom *= factor
	if c.Zoom < 0.2 {
		c.Zoom = 0.2
	}
	if c.Zoom > 5.0 {
		c.Zoom = 5.0
	}

	c.X = worldX - (screenX-halfW)/c.Zoom
	c.Y = worldY - (screenY-halfH)/c.Zoom
}
func (c *Camera) Update() {
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		factor := 1.05
		mx, my := ebiten.CursorPosition()
		c.ZoomAt(float64(mx), float64(my), factor)
	}
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		factor := 0.95
		mx, my := ebiten.CursorPosition()
		c.ZoomAt(float64(mx), float64(my), factor)
	}
}
