package game

import (
	. "DungeonHive/config"
	"DungeonHive/renderer"

	"github.com/hajimehoshi/ebiten/v2"
)

var rmap = [GridSize][GridSize]int{}
var tmap = [GridSize][GridSize]int{}

type Game struct {
	RawMap   *[GridSize][GridSize]int
	TileMap  *[GridSize][GridSize]int
	Rooms    []Room
	Ticks    int
	Money    int
	Renderer *renderer.Renderer
	Player   *Player
}

func NewGame() *Game {
	g := &Game{
		RawMap:  &rmap,
		TileMap: &tmap,
	}
	rooms := GenerateMapAndRooms(DefaultConfig(), g.RawMap)
	BuildTileMap(g.RawMap, g.TileMap)
	g.Renderer = renderer.NewRenderer() // 初始化渲染器
	g.Renderer.DrawMapLayer(g.TileMap)
	g.Rooms = rooms // 初始化 Room 的 Grid
	g.Money = 0     // 初始化金钱

	if len(rooms) > 0 {
		g.Player = NewPlayer(rooms[0].CenterX(), rooms[0].CenterY())
	}

	return g
}

func (g *Game) Update() error {
	g.Ticks++
	g.Renderer.Camera.Update()
	if g.Player != nil {
		g.Player.Update(g.RawMap)
		// 鏡頭跟隨玩家
		g.Renderer.Camera.X = g.Player.X
		g.Renderer.Camera.Y = g.Player.Y
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.Renderer.DrawLayer(screen)
	if g.Player != nil {
		g.Player.Draw(screen, g.Renderer.Camera)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}
