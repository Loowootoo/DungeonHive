package game

import (
	. "DungeonHive/config"
	"DungeonHive/utils"
)

// ConnectRooms 用 Prim MST 決定房間連接順序，
// 再對每對房間呼叫 A* 挖出 1 格寬通道。
func ConnectRooms(m *[GridSize][GridSize]int, rooms []Room) {
	if len(rooms) == 0 {
		return
	}
	mapH := GridSize
	mapW := GridSize

	connected := make([]bool, len(rooms))
	connected[0] = true
	for {
		src, dst := primNextEdge(rooms, connected)
		if src < 0 {
			break
		}
		connected[dst] = true

		start := Point{X: rooms[src].CenterX(), Y: rooms[src].CenterY()}
		end := Point{X: rooms[dst].CenterX(), Y: rooms[dst].CenterY()}

		carveCorridor(m, mapW, mapH, start, end, rooms[src], rooms[dst])
	}
}

// primNextEdge 回傳 Prim MST 下一條最短邊的 (src, dst) 房間索引。
// 若所有房間已連接則回傳 (-1, -1)。
func primNextEdge(rooms []Room, connected []bool) (int, int) {
	bestDist := int(^uint(0) >> 1) // MaxInt
	bestSrc, bestDst := -1, -1

	for i, ri := range rooms {
		if !connected[i] {
			continue
		}
		for j, rj := range rooms {
			if connected[j] {
				continue
			}
			dx := ri.CenterX() - rj.CenterX()
			dy := ri.CenterY() - rj.CenterY()
			if d := dx*dx + dy*dy; d < bestDist {
				bestDist = d
				bestSrc, bestDst = i, j
			}
		}
	}
	return bestSrc, bestDst
}

// carveCorridor 以 A* 找路並逐格挖開，形成 1 格寬通道。
func carveCorridor(m *[GridSize][GridSize]int, mapW, mapH int, start, end Point, r1, r2 Room) {
	costFn := makeCostFn(m, mapW, mapH, r1, r2)
	path := utils.Search(start, end, costFn)
	if path == nil {
		return
	}
	for _, p := range path {
		m[p.Y][p.X] = 0
	}
}

// makeCostFn 產生 A* 的代價函式。
// 規則：
//   - 超出邊界          → -1（不可通行）
//   - 3x3 範圍內碰到其他房間的地板 → 代價 100（高代價繞道）
//   - 走既有通道        → 代價 1
//   - 挖新牆壁          → 代價 2
func makeCostFn(m *[GridSize][GridSize]int, mapW, mapH int, r1, r2 Room) utils.CostFunc {
	return func(next Point) int {
		// 邊界：留 2 格緩衝避免緊貼地圖邊緣
		if next.X < 2 || next.X >= mapW-2 || next.Y < 2 || next.Y >= mapH-2 {
			return -1
		}

		// 掃描 next 的 3x3 鄰域，偵測是否侵犯其他房間
		if neighbourForbidden(m, mapW, mapH, next, r1, r2) {
			return 100
		}

		if m[next.Y][next.X] == 1 {
			return 2 // 挖牆
		}
		return 1 // 走通道
	}
}

// neighbourForbidden 掃描 (next) 周圍 3x3 格，
// 若發現不屬於 r1/r2 的地板則回傳 true。
// 這確保通道不會緊貼其他房間，autotile 才能正確顯示牆壁分隔。
func neighbourForbidden(m *[GridSize][GridSize]int, mapW, mapH int, next Point, r1, r2 Room) bool {
	for ty := next.Y - 1; ty <= next.Y+1; ty++ {
		for tx := next.X - 1; tx <= next.X+1; tx++ {
			if ty < 0 || ty >= mapH || tx < 0 || tx >= mapW {
				continue
			}
			if r1.Contains(tx, ty) || r2.Contains(tx, ty) {
				continue
			}
			if m[ty][tx] == 0 {
				return true
			}
		}
	}
	return false
}
