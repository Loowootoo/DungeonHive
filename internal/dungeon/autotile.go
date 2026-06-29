package dungeon

import (
	. "DungeonHive/pkg/config"
)

// neighbours 儲存一個格子八個方向的牆壁狀態。
type neighbours struct {
	N, E, S, W     bool
	NE, SE, SW, NW bool
}

// fallbackTileID 是當 mask 不符合任何特殊規則時使用的對照表（16 種 4-bit 組合）。
var fallbackTileID = [16]int{
	0, 3, 43, 17,
	15, 15, 26, 18,
	43, 13, 43, 43,
	16, 43, 43, 44,
}

// getTile 安全地讀取地圖格子值，超出邊界視為牆壁（1）。
func getTile(m *[GridSize][GridSize]int, x, y int) int {
	if y < 0 || y >= len(m) || x < 0 || x >= len(m[0]) {
		return 1
	}
	return m[y][x]
}

// wallTileID 根據八方向鄰居計算該格子應使用的 tile ID。
// 若該格子是地板（0）則回傳 0（不繪製）。
func wallTileID(m *[GridSize][GridSize]int, x, y int) int {
	if getTile(m, x, y) == 0 {
		return 0
	}
	nb := neighbours{
		N:  getTile(m, x, y-1) == 1,
		E:  getTile(m, x+1, y) == 1,
		S:  getTile(m, x, y+1) == 1,
		W:  getTile(m, x-1, y) == 1,
		NE: getTile(m, x+1, y-1) == 1,
		SE: getTile(m, x+1, y+1) == 1,
		SW: getTile(m, x-1, y+1) == 1,
		NW: getTile(m, x-1, y-1) == 1,
	}

	mask := 0
	if nb.N {
		mask |= 1
	}
	if nb.E {
		mask |= 2
	}
	if nb.S {
		mask |= 4
	}
	if nb.W {
		mask |= 8
	}

	// 內角（凹角）優先判斷
	if nb.N && nb.W && !nb.NW && mask != 13 && mask != 11 {
		return 44
	}
	if nb.N && nb.E && !nb.NE && mask != 11 && mask != 7 {
		return 41
	}
	if nb.S && nb.W && !nb.SW && mask != 14 && mask != 13 {
		return 5
	}
	if nb.S && nb.E && !nb.SE && mask != 14 {
		return 2
	}

	// 直線牆
	switch mask {
	case 3:
		return 43
	case 6:
		return 16
	case 9:
		return 43
	case 12:
		return 17
	}

	// 孤立牆與雜項
	if mask == 5 {
		return 4
	}
	if mask == 10 {
		return 43
	}
	if mask == 7 {
		return 18
	}
	if mask == 11 {
		return 43
	}
	if mask == 13 {
		return 15
	}
	if mask == 14 {
		return 43
	}
	if mask == 15 {
		return 22
	}

	return fallbackTileID[mask]
}

// BuildTileMap 將原始地圖轉換成 tile ID 地圖。
// 回傳的二維陣列與輸入等大，每格存放對應的 tile ID。
func BuildTileMap(rawMap *[GridSize][GridSize]int, tileMap *[GridSize][GridSize]int) {
	for y := 0; y < GridSize; y++ {
		for x := 0; x < GridSize; x++ {
			tileMap[y][x] = wallTileID(rawMap, x, y)
		}
	}
}
