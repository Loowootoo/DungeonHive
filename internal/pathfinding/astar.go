package pathfinding

import (
	. "DungeonHive/pkg/config"
	"container/heap"
)

// item 是優先佇列內部使用的節點。
type item struct {
	point    Point
	priority int // fScore = gScore + heuristic
	index    int
}

// priorityQueue 實作 heap.Interface，以 fScore 最小優先。
type priorityQueue []*item

func (pq priorityQueue) Len() int           { return len(pq) }
func (pq priorityQueue) Less(i, j int) bool { return pq[i].priority < pq[j].priority }
func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}
func (pq *priorityQueue) Push(x interface{}) {
	n := len(*pq)
	it := x.(*item)
	it.index = n
	*pq = append(*pq, it)
}
func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	it := old[n-1]
	old[n-1] = nil
	it.index = -1
	*pq = old[0 : n-1]
	return it
}

func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// heuristic 回傳曼哈頓距離作為啟發函式。
func heuristic(a, b Point) int {
	return absInt(a.X-b.X) + absInt(a.Y-b.Y)
}

// CostFunc 是一個回呼，讓呼叫方決定移動到 next 的代價。
// 回傳 -1 表示此格完全不可通行（直接跳過）。
type CostFunc func(next Point) int

// Search 執行 A* 搜尋，從 start 到 end。
// costFn 由呼叫方提供，決定每個候選格子的移動代價。
// 回傳從 start 到 end 的路徑（含兩端點），若找不到則回傳 nil。
func Search(start, end Point, costFn CostFunc) []Point {
	pq := make(priorityQueue, 0)
	heap.Init(&pq)

	gScore := make(map[Point]int)
	cameFrom := make(map[Point]Point)

	gScore[start] = 0
	heap.Push(&pq, &item{
		point:    start,
		priority: heuristic(start, end),
	})

	dirs := []Point{{X: 0, Y: -1}, {X: 1, Y: 0}, {X: 0, Y: 1}, {X: -1, Y: 0}}

	for pq.Len() > 0 {
		curr := heap.Pop(&pq).(*item).point

		if curr == end {
			return reconstructPath(cameFrom, start, end)
		}

		for _, d := range dirs {
			next := Point{X: curr.X + d.X, Y: curr.Y + d.Y}

			cost := costFn(next)
			if cost < 0 {
				continue // 不可通行
			}

			tentative := gScore[curr] + cost
			if val, exists := gScore[next]; !exists || tentative < val {
				gScore[next] = tentative
				cameFrom[next] = curr
				heap.Push(&pq, &item{
					point:    next,
					priority: tentative + heuristic(next, end),
				})
			}
		}
	}

	return nil // 無路可走
}

// reconstructPath 從 cameFrom 表回溯路徑。
func reconstructPath(cameFrom map[Point]Point, start, end Point) []Point {
	path := []Point{}
	curr := end
	for {
		path = append(path, curr)
		next, exists := cameFrom[curr]
		if !exists || curr == start {
			break
		}
		curr = next
	}
	// 反轉，讓路徑從 start 到 end
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path
}
