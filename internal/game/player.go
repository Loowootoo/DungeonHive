package game

import (
	. "DungeonHive/pkg/config"
	"DungeonHive/internal/renderer"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	moveFrames = 8
	playerSize = 12
)

type Player struct {
	X, Y             float64
	GridX, GridY     int
	StartX, StartY   float64
	TargetX, TargetY float64
	MoveTimer        int
	Dir              DirType
	sprite           *ebiten.Image
}

// NewPlayer 在指定的 tile 中心建立玩家。
func NewPlayer(tileX, tileY int) *Player {
	px := float64(tileX*TileSize) + float64(TileSize)/2
	py := float64(tileY*TileSize) + float64(TileSize)/2

	p := &Player{
		X:      px,
		Y:      py,
		GridX:  tileX,
		GridY:  tileY,
		StartX: px,
		StartY: py,
		TargetX: px,
		TargetY: py,
		Dir:    DirUP,
	}

	// 用金色方塊作為玩家 sprite（可替換為 assets/ow*.png）
	p.sprite = ebiten.NewImage(playerSize, playerSize)
	p.sprite.Fill(color.RGBA{240, 220, 80, 255})

	return p
}

func (p *Player) Update(rawMap *[GridSize][GridSize]int) {
	if p.MoveTimer > 0 {
		p.MoveTimer++
		if p.MoveTimer >= moveFrames {
			p.X = p.TargetX
			p.Y = p.TargetY
			p.GridX = int(p.TargetX) / TileSize
			p.GridY = int(p.TargetY) / TileSize
			p.MoveTimer = 0
		} else {
			t := float64(p.MoveTimer) / float64(moveFrames)
			p.X = p.StartX + (p.TargetX-p.StartX)*t
			p.Y = p.StartY + (p.TargetY-p.StartY)*t
		}
		return
	}

	dx, dy := 0, 0
	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		dy = -1
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		dy = 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		dx = -1
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		dx = 1
	}

	if dx == 0 && dy == 0 {
		return
	}

	switch {
	case dy < 0 && dx == 0:
		p.Dir = DirUP
	case dy > 0 && dx == 0:
		p.Dir = DirDOWN
	case dx < 0 && dy == 0:
		p.Dir = DirLEFT
	case dx > 0 && dy == 0:
		p.Dir = DirRIGHT
	case dy < 0 && dx < 0:
		p.Dir = DirUPLEFT
	case dy < 0 && dx > 0:
		p.Dir = DirUPRIGHT
	case dy > 0 && dx < 0:
		p.Dir = DirDOWNLEFT
	case dy > 0 && dx > 0:
		p.Dir = DirDOWNRIGHT
	}

	targetGX := p.GridX + dx
	targetGY := p.GridY + dy

	if targetGX < 0 || targetGX >= GridSize || targetGY < 0 || targetGY >= GridSize {
		return
	}
	if rawMap[targetGY][targetGX] == 1 {
		return
	}

	p.StartX = p.X
	p.StartY = p.Y
	p.TargetX = float64(targetGX*TileSize) + float64(TileSize)/2
	p.TargetY = float64(targetGY*TileSize) + float64(TileSize)/2
	p.MoveTimer = 1
}

// Draw 以鏡頭轉換繪製玩家 sprite。
func (p *Player) Draw(screen *ebiten.Image, cam *renderer.Camera) {
	renderer.DrawPlayer(screen, p.sprite, p.X, p.Y, cam)
}
