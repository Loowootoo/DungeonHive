// Package player 管理玩家角色在地圖上的移動與繪製。
package game

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Player 代表玩家角色，位置以地圖格座標表示。
type Player struct {
	TileX, TileY int // 地圖格座標
	size         int // 玩家繪製大小（等同 tilesize）
}

// New 建立一個位於 (startX, startY) 格的新玩家。
func New(startX, startY, tileSize int) *Player {
	return &Player{
		TileX: startX,
		TileY: startY,
		size:  tileSize,
	}
}

// HandleInput 處理 WASD / 方向鍵移動。
// rawMap[y][x] == 0 才算可通行，牆壁(1)與超出邊界都無法通過。
func (p *Player) HandleInput(rawMap [][]int) {
	dx, dy := 0, 0
	switch {
	case rl.IsKeyPressed(rl.KeyW) || rl.IsKeyPressed(rl.KeyUp):
		dy = -1
	case rl.IsKeyPressed(rl.KeyS) || rl.IsKeyPressed(rl.KeyDown):
		dy = 1
	case rl.IsKeyPressed(rl.KeyA) || rl.IsKeyPressed(rl.KeyLeft):
		dx = -1
	case rl.IsKeyPressed(rl.KeyD) || rl.IsKeyPressed(rl.KeyRight):
		dx = 1
	}
	if dx == 0 && dy == 0 {
		return
	}

	newX, newY := p.TileX+dx, p.TileY+dy

	// 邊界檢查
	mapH := len(rawMap)
	if mapH == 0 {
		return
	}
	mapW := len(rawMap[0])
	if newX < 0 || newX >= mapW || newY < 0 || newY >= mapH {
		return
	}

	// 碰撞檢測：只有地板(0)才能走
	if rawMap[newY][newX] != 0 {
		return
	}

	p.TileX = newX
	p.TileY = newY
}

// WorldX 回傳玩家在畫面上的 X 座標（像素）。
func (p *Player) WorldX() float32 {
	return float32(p.TileX * p.size)
}

// WorldY 回傳玩家在畫面上的 Y 座標（像素）。
func (p *Player) WorldY() float32 {
	return float32(p.TileY * p.size)
}

// Draw 在目前 2D 攝影機範圍內繪製玩家。
func (p *Player) Draw() {
	cx := p.WorldX() + float32(p.size)/2
	cy := p.WorldY() + float32(p.size)/2
	radius := float32(p.size) / 2.8

	// 發光外圍
	rl.DrawCircle(int32(cx), int32(cy), radius+2, rl.Fade(rl.SkyBlue, 0.3))
	// 主體
	rl.DrawCircle(int32(cx), int32(cy), radius, rl.SkyBlue)
	// 外框
	rl.DrawCircleLines(int32(cx), int32(cy), radius, rl.NewColor(20, 60, 120, 255))

	// 眼睛（兩個小白點）
	eyeOff := float32(p.size) / 6
	eyeR := float32(p.size) / 12
	rl.DrawCircle(int32(cx-eyeOff), int32(cy-eyeOff/2), eyeR, rl.White)
	rl.DrawCircle(int32(cx+eyeOff), int32(cy-eyeOff/2), eyeR, rl.White)
}
