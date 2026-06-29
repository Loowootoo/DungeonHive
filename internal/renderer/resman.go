package renderer

import (
	"fmt"
	"image"
	_ "image/png"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

type ResourceManager struct {
	TileImage *ebiten.Image
}

func loadResources(fn string) *ebiten.Image {
	f, err := os.Open(fn)
	if err != nil {
		fmt.Printf("Error opening %s: %s\n", fn, err)
		return nil
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil
	}
	return ebiten.NewImageFromImage(img)
}

func NewResourceManager() *ResourceManager {
	return &ResourceManager{
		TileImage: loadResources("assets/tileset.png"),
	}
}
