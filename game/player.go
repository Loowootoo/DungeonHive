// Package player 管理玩家角色在地圖上的移動與繪製。
package game

import (
	cf "DungeonHive/config"
)

type Player struct {
	X, Y int
	Dir  cf.DirType
}

func NewPlayer(x, y int) *Player {
	return &Player{
		X:   x,
		Y:   y,
		Dir: cf.DirUP,
	}
}
