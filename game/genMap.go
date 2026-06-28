package game

import (
	. "DungeonHive/config"
	"math/rand"
	"time"
)

// Config 包含地圖生成的所有參數。
type Config struct {
	MinRooms, MaxRooms int
	RoomMin, RoomMax   int
	RoomPadding        int // 房間之間保留的最小牆壁格數
}

// DefaultConfig 回傳預設配置。
func DefaultConfig() Config {
	return Config{
		MinRooms:    6,
		MaxRooms:    12,
		RoomMin:     5,
		RoomMax:     9,
		RoomPadding: 2,
	}
}

// Generate 產生一張完整的隨機地牢原始地圖（0=地板，1=牆壁）。
// 同時回傳成功放置的房間列表，供 player 重生點等用途。
func GenerateMapAndRooms(cf Config, gmap [GridSize][GridSize]int) []Room {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	rooms := PlaceRooms(
		rng,
		cf.MinRooms, cf.MaxRooms,
		cf.RoomMin, cf.RoomMax,
		cf.RoomPadding,
		gmap,
	)
	ConnectRooms(gmap, rooms)
	return rooms
}
