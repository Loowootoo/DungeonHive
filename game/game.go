package game

import (
	. "DungeonHive/config"
	"DungeonHive/renderer"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	RawMap   [GridSize][GridSize]int
	TileMap  [GridSize][GridSize]int
	Rooms    []Room
	Ticks    int
	Money    int
	Renderer *renderer.Renderer
}

func NewGame() *Game {
	g := &Game{}
	rooms := GenerateMapAndRooms(DefaultConfig(), g.RawMap)
	BuildTileMap(g.RawMap, g.TileMap)
	g.Renderer = renderer.NewRenderer() // 初始化渲染器
	g.Renderer.DrawMapLayer(g.TileMap)
	g.Rooms = rooms // 初始化 Room 的 Grid
	g.Money = 0     // 初始化金钱
	return g
}

func (g *Game) Update() error {
	g.Ticks++
	g.Renderer.Camera.Update()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.Renderer.DrawLayer(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}
