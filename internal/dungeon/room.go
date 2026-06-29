package dungeon

import (
	. "DungeonHive/pkg/config"
	"math/rand"
)

// Room 代表地圖上一個矩形房間。
type Room struct {
	X, Y int // 左上角座標
	W, H int // 寬、高
}

// CenterX 回傳房間水平中心格座標。
func (r Room) CenterX() int { return r.X + r.W/2 }

// CenterY 回傳房間垂直中心格座標。
func (r Room) CenterY() int { return r.Y + r.H/2 }

// Contains 判斷格子 (x, y) 是否在房間內。
func (r Room) Contains(x, y int) bool {
	return x >= r.X && x < r.X+r.W && y >= r.Y && y < r.Y+r.H
}

// overlaps 判斷兩個房間（加上各自的 padding）是否重疊。
func overlaps(a, b Room, padding int) bool {
	return a.X-padding < b.X+b.W+padding &&
		a.X+a.W+padding > b.X-padding &&
		a.Y-padding < b.Y+b.H+padding &&
		a.Y+a.H+padding > b.Y-padding
}

// PlaceRooms 在 mapW×mapH 的地圖上隨機放置房間，
// 並將房間所在格子設為 0（地板），其餘保持 1（牆壁）。
// 回傳成功放置的 []Room 與初始化後的地圖。
func PlaceRooms(rng *rand.Rand, minRooms, maxRooms int, roomMin, roomMax int, padding int, gmap *[GridSize][GridSize]int) []Room {
	// 初始化全牆地圖
	makeWallMap(gmap)
	target := minRooms + rng.Intn(maxRooms-minRooms+1)
	rooms := make([]Room, 0, target)
	attempts := 0

	for len(rooms) < target && attempts < 2000 {
		w := roomMin + rng.Intn(roomMax-roomMin+1)
		h := roomMin + rng.Intn(roomMax-roomMin+1)
		x := 4 + rng.Intn(GridSize-w-8)
		y := 4 + rng.Intn(GridSize-h-8)

		candidate := Room{x, y, w, h}
		if hasOverlap(candidate, rooms, padding) {
			attempts++
			continue
		}

		rooms = append(rooms, candidate)
		carveRoom(gmap, candidate)
		attempts = 0
	}

	// 至少保留一個房間，防止空地圖
	if len(rooms) == 0 {
		fallback := Room{5, 5, 6, 6}
		rooms = append(rooms, fallback)
		carveRoom(gmap, fallback)
	}

	return rooms
}

// ---- 私有輔助函式 ----

func makeWallMap(gmap *[GridSize][GridSize]int) {
	for y := 0; y < GridSize; y++ {
		for x := 0; x < GridSize; x++ {
			gmap[y][x] = 1
		}
	}
}

func hasOverlap(candidate Room, rooms []Room, padding int) bool {
	for _, r := range rooms {
		if overlaps(candidate, r, padding) {
			return true
		}
	}
	return false
}

func carveRoom(m *[GridSize][GridSize]int, r Room) {
	for dy := 0; dy < r.H; dy++ {
		for dx := 0; dx < r.W; dx++ {
			m[r.Y+dy][r.X+dx] = 0
		}
	}
}
