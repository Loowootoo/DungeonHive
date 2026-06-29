package game

import (
	. "DungeonHive/config"
	"DungeonHive/renderer"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	playerSpeed = 2.0  // 移動速度（像素/ frame）
	playerSize  = 12   // 玩家 sprite 尺寸（比 tile 16px 略小，方便鑽通道）
)

type Player struct {
	X, Y   float64 // 世界座標（像素），sprite 中心位置
	Dir    DirType
	sprite *ebiten.Image
}

// NewPlayer 在指定的 tile 中心建立玩家。
func NewPlayer(tileX, tileY int) *Player {
	px := float64(tileX*TileSize) + float64(TileSize)/2
	py := float64(tileY*TileSize) + float64(TileSize)/2

	p := &Player{
		X:   px,
		Y:   py,
		Dir: DirUP,
	}

	// 用金色方塊作為玩家 sprite（可替換為 assets/ow*.png）
	p.sprite = ebiten.NewImage(playerSize, playerSize)
	p.sprite.Fill(color.RGBA{240, 220, 80, 255})

	return p
}

// Update 處理 WASD / 方向鍵輸入與牆壁碰撞。
func (p *Player) Update(rawMap *[GridSize][GridSize]int) {
	dx, dy := 0.0, 0.0

	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		dy = -playerSpeed
		p.Dir = DirUP
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		dy = playerSpeed
		p.Dir = DirDOWN
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		dx = -playerSpeed
		p.Dir = DirLEFT
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		dx = playerSpeed
		p.Dir = DirRIGHT
	}

	// 分別處理 X/Y 軸碰撞，達成貼牆滑行（wall sliding）
	if dx != 0 {
		newX := p.X + dx
		if !p.collidesAt(newX, p.Y, rawMap) {
			p.X = newX
		}
	}
	if dy != 0 {
		newY := p.Y + dy
		if !p.collidesAt(p.X, newY, rawMap) {
			p.Y = newY
		}
	}
}

// collidesAt 檢查玩家 bounding box 位於 (x, y) 時是否會碰撞。
func (p *Player) collidesAt(x, y float64, rawMap *[GridSize][GridSize]int) bool {
	half := float64(playerSize) / 2

	// 檢查四個角落
	corners := [][2]float64{
		{x - half, y - half}, // 左上
		{x + half, y - half}, // 右上
		{x + half, y + half}, // 右下
		{x - half, y + half}, // 左下
	}

	for _, c := range corners {
		tx := int(c[0]) / TileSize
		ty := int(c[1]) / TileSize
		if tx < 0 || tx >= GridSize || ty < 0 || ty >= GridSize {
			return true // 超出地圖邊界視為碰撞
		}
		if rawMap[ty][tx] == 1 {
			return true
		}
	}
	return false
}

// Draw 以鏡頭轉換繪製玩家 sprite。
func (p *Player) Draw(screen *ebiten.Image, cam *renderer.Camera) {
	renderer.DrawPlayer(screen, p.sprite, p.X, p.Y, cam)
}
