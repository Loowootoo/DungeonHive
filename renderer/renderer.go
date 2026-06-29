package renderer

import (
	cf "DungeonHive/config"

	"github.com/hajimehoshi/ebiten/v2"
)

type Renderer struct {
	ResMan   *ResourceManager
	Camera   *Camera
	MapLayer *ebiten.Image
}

func NewRenderer() *Renderer {
	renderer := &Renderer{}
	renderer.ResMan = NewResourceManager()
	renderer.Camera = NewCamera(float64(cf.ScreenWidth/2), float64(cf.ScreenHeight/2), 2.0)
	renderer.MapLayer = ebiten.NewImage(cf.ScreenWidth, cf.ScreenHeight)

	return renderer
}

func (r *Renderer) DrawLayer(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	screenW := float64(screen.Bounds().Dx())
	screenH := float64(screen.Bounds().Dy())
	op.GeoM.Reset()
	op.GeoM.Translate(-r.Camera.X, -r.Camera.Y)
	op.GeoM.Scale(r.Camera.Zoom, r.Camera.Zoom)
	op.GeoM.Translate(screenW/2, screenH/2)
	screen.DrawImage(r.MapLayer, op)
}
