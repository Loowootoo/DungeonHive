package pathfinding

import (
	. "DungeonHive/pkg/config"
)

// Grid 使用固定陣列，避免 [][]int 動態配置
type Grid struct {
	W, H  int
	Cells [GridSize][GridSize]uint8 // 0 = walkable, 1 = blocked
}

func (g *Grid) InitSquare(size int) {
	if size <= 0 || size > GridSize {
		panic("grid size out of range")
	}
	g.W, g.H = size, size
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			g.Cells[y][x] = 0
		}
	}
}

func (g *Grid) InBounds(p Point) bool {
	return p.X >= 0 && p.X < g.W && p.Y >= 0 && p.Y < g.H
}

func (g *Grid) Walkable(p Point) bool {
	return g.InBounds(p) && g.Cells[p.Y][p.X] == 0
}

// Path 也用固定陣列，避免 []Point 配置
type Path struct {
	Points [MaxCells]Point
	Len    int
}

func (p *Path) Reset() {
	p.Len = 0
}

type heapNode struct {
	P Point
	F int
}

// AStarFixedContext 可重複使用，避免每次搜尋都配置記憶體
type AStarFixedContext struct {
	heap     [MaxCells]heapNode
	heapSize int

	gScore  [GridSize][GridSize]int
	parentX [GridSize][GridSize]int16
	parentY [GridSize][GridSize]int16
	openIdx [GridSize][GridSize]int
	closed  [GridSize][GridSize]bool
}

func (ctx *AStarFixedContext) heuristic(a, b Point) int {
	return abs(a.X-b.X) + abs(a.Y-b.Y)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (ctx *AStarFixedContext) reset(w, h int) {
	ctx.heapSize = 0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			ctx.gScore[y][x] = Inf
			ctx.parentX[y][x] = -1
			ctx.parentY[y][x] = -1
			ctx.openIdx[y][x] = -1
			ctx.closed[y][x] = false
		}
	}
}

func (ctx *AStarFixedContext) swapHeap(i, j int) {
	ctx.heap[i], ctx.heap[j] = ctx.heap[j], ctx.heap[i]
	pi := ctx.heap[i].P
	pj := ctx.heap[j].P
	ctx.openIdx[pi.Y][pi.X] = i
	ctx.openIdx[pj.Y][pj.X] = j
}

func (ctx *AStarFixedContext) siftUp(i int) {
	for i > 0 {
		parent := (i - 1) / 2
		if ctx.heap[i].F >= ctx.heap[parent].F {
			break
		}
		ctx.swapHeap(i, parent)
		i = parent
	}
}

func (ctx *AStarFixedContext) siftDown(i int) {
	for {
		left := 2*i + 1
		right := 2*i + 2
		smallest := i

		if left < ctx.heapSize && ctx.heap[left].F < ctx.heap[smallest].F {
			smallest = left
		}
		if right < ctx.heapSize && ctx.heap[right].F < ctx.heap[smallest].F {
			smallest = right
		}
		if smallest == i {
			return
		}

		ctx.swapHeap(i, smallest)
		i = smallest
	}
}

func (ctx *AStarFixedContext) pushOpen(p Point, f int) {
	if ctx.heapSize >= MaxCells {
		panic("fixed heap overflow")
	}
	idx := ctx.heapSize
	ctx.heap[idx] = heapNode{P: p, F: f}
	ctx.openIdx[p.Y][p.X] = idx
	ctx.heapSize++
	ctx.siftUp(idx)
}

func (ctx *AStarFixedContext) decreaseOrPushOpen(p Point, f int) {
	idx := ctx.openIdx[p.Y][p.X]
	if idx == -1 {
		ctx.pushOpen(p, f)
		return
	}
	if f >= ctx.heap[idx].F {
		return
	}
	ctx.heap[idx].F = f
	ctx.siftUp(idx)
}

func (ctx *AStarFixedContext) popOpen() (heapNode, bool) {
	if ctx.heapSize == 0 {
		return heapNode{}, false
	}

	top := ctx.heap[0]
	ctx.openIdx[top.P.Y][top.P.X] = -1
	ctx.heapSize--

	if ctx.heapSize > 0 {
		ctx.heap[0] = ctx.heap[ctx.heapSize]
		p := ctx.heap[0].P
		ctx.openIdx[p.Y][p.X] = 0
		ctx.siftDown(0)
	}

	return top, true
}

func (ctx *AStarFixedContext) buildPath(end Point, out *Path) {
	out.Len = 0
	cur := end

	for {
		if out.Len >= MaxCells {
			panic("path overflow")
		}
		out.Points[out.Len] = cur
		out.Len++

		px := ctx.parentX[cur.Y][cur.X]
		py := ctx.parentY[cur.Y][cur.X]

		if px == -1 && py == -1 {
			break
		}

		cur = Point{X: int(px), Y: int(py)}
	}

	// reverse in place
	for i, j := 0, out.Len-1; i < j; i, j = i+1, j-1 {
		out.Points[i], out.Points[j] = out.Points[j], out.Points[i]
	}
}

// FindPathFixedNoGC 是熱路徑版本：
// - 不分配 slice
// - 不用 interface
// - 不用 sync.Pool
// - 不用 container/heap
func (ctx *AStarFixedContext) FindPathFixedNoGC(g *Grid, start, end Point, out *Path) bool {
	out.Reset()

	if g.W <= 0 || g.H <= 0 || g.W > GridSize || g.H > GridSize {
		return false
	}
	if !g.Walkable(start) || !g.Walkable(end) {
		return false
	}

	ctx.reset(g.W, g.H)

	ctx.gScore[start.Y][start.X] = 0
	ctx.pushOpen(start, heuristic(start, end))

	dirs := [4]Point{
		{X: 0, Y: -1},
		{X: 0, Y: 1},
		{X: -1, Y: 0},
		{X: 1, Y: 0},
	}

	for ctx.heapSize > 0 {
		node, ok := ctx.popOpen()
		if !ok {
			break
		}

		cur := node.P
		if ctx.closed[cur.Y][cur.X] {
			continue
		}
		ctx.closed[cur.Y][cur.X] = true

		if cur == end {
			ctx.buildPath(end, out)
			return true
		}

		curG := ctx.gScore[cur.Y][cur.X]

		for _, d := range dirs {
			nb := Point{X: cur.X + d.X, Y: cur.Y + d.Y}

			if !g.Walkable(nb) || ctx.closed[nb.Y][nb.X] {
				continue
			}

			newG := curG + 1
			if newG < ctx.gScore[nb.Y][nb.X] {
				ctx.gScore[nb.Y][nb.X] = newG
				ctx.parentX[nb.Y][nb.X] = int16(cur.X)
				ctx.parentY[nb.Y][nb.X] = int16(cur.Y)

				newF := newG + heuristic(nb, end)
				ctx.decreaseOrPushOpen(nb, newF)
			}
		}
	}

	return false
}

/*
// ====================== 固定陣列地圖產生器 ======================

// 用簡單 PRNG，避免 math/rand 帶來額外配置干擾
type XorShift64 struct {
	state uint64
}

func NewXorShift64(seed uint64) XorShift64 {
	if seed == 0 {
		seed = 0x9e3779b97f4a7c15
	}
	return XorShift64{state: seed}
}

func (r *XorShift64) Next() uint64 {
	x := r.state
	x ^= x << 13
	x ^= x >> 7
	x ^= x << 17
	r.state = x
	return x
}

func GenerateMazeFixed(g *Grid, size int, seed uint64) {
	g.InitSquare(size)

	rng := NewXorShift64(seed)
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			if rng.Next()%100 < 25 {
				g.Cells[y][x] = 1
			} else {
				g.Cells[y][x] = 0
			}
		}
	}

	g.Cells[0][0] = 0
	g.Cells[size-1][size-1] = 0
}

// 反覆換 seed，直到產生可通行地圖
func GenerateSolvableMazeFixed(g *Grid, size int, startSeed uint64, ctx *AStarFixedContext, out *Path) uint64 {
	start := Point{0, 0}
	end := Point{size - 1, size - 1}

	for seed := startSeed; ; seed++ {
		GenerateMazeFixed(g, size, seed)
		if ctx.FindPathFixedNoGC(g, start, end, out) {
			return seed
		}
	}
}

func main() {
	// 示範：整個程式先關自動 GC
	// 實務上比較建議只在 benchmark/test 時關
	oldGC := debug.SetGCPercent(-1)
	defer debug.SetGCPercent(oldGC)

	const size = 50

	var grid Grid
	var ctx AStarFixedContext
	var path Path

	seed := GenerateSolvableMazeFixed(&grid, size, 1, &ctx, &path)

	start := Point{0, 0}
	end := Point{size - 1, size - 1}

	ok := ctx.FindPathFixedNoGC(&grid, start, end, &path)
	if !ok {
		fmt.Println("找不到路")
		return
	}

	fmt.Println("fixed-only / no-GC 版本")
	fmt.Printf("grid size   : %dx%d\n", grid.W, grid.H)
	fmt.Printf("used seed   : %d\n", seed)
	fmt.Printf("path length : %d\n", path.Len)
	fmt.Printf("start       : %+v\n", start)
	fmt.Printf("end         : %+v\n", end)
	for i:=0;i<path.Len;i++ {
		fmt.Printf("%d = [%d,%d] -> ", i, path.Points[i].X, path.Points[i].Y)
	}

}
*/
