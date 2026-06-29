# DungeonHive — AGENTS.md

> **語言設定：請全程使用繁體中文（臺灣）回應。**

**Ebitengine v2** roguelike-dungeon generator.  
Build: `go run .` from repo root. Go 1.26, single module `DungeonHive`.

---

## Architecture (15 Go files, 4 packages)

| Directory | Role | Entrypoint |
|-----------|------|------------|
| `main.go` | Ebitengine game loop, 960×960 window | `game.NewGame()` → `ebiten.RunGame` |
| `game/` | Map generation, autotiling, player, game state | `Game` struct, `Update()` / `Draw()` / `Layout()` |
| `renderer/` | Camera (WASD+Q/E), resource loading, tile-map blit | `Renderer`, `Camera`, `ResourceManager` |
| `config/` | Constants only: `GridSize=60`, `TileSize=16`, `Point` type | `constant.go` |
| `utils/` | A* pathfinding (standard + no-GC fixed-array variant) | `Search()` and `AStarFixedContext.FindPathFixedNoGC()` |
| `assets/` | `tileset.png` (13×5 grid, 16px tiles) + 10 overworld sprites | Loaded via `assets/tileset.png` relative path |

**Global state** — `game/game.go` declares `rmap` and `tmap` as package-level `[60][60]int` arrays; `Game` holds pointers to them. No other global mutable state.

---

## Key codebase facts

### Dot-import convention
- `game/` and `utils/` use `import . "DungeonHive/config"` — all config identifiers (GridSize, Point, DirType, TileSize, etc.) are in scope without prefix.
- `renderer/` uses `import cf "DungeonHive/config"` — prefixed with `cf.` instead.

### Map generation flow
```
DefaultConfig() → GenerateMapAndRooms() → PlaceRooms() → ConnectRooms() (Prim MST + A*)
                                        → BuildTileMap() (autotile wall-ID mapping)
```
- Raw map: `0` = floor, `1` = wall.
- Tile map: `0` = skip drawing, nonzero = tile index into `tileset.png` (column-major, 13 cols).
- `autotile.go` handles 8-direction neighbour-aware wall tile selection.

### A* — two implementations
1. **`utils/astar.go`** — standard `Search(start, end, CostFunc)`. Uses `map[Point]int` / `container/heap` — GC-allocating. Used by corridor carving.
2. **`utils/astarNoGC.go`** — `AStarFixedContext.FindPathFixedNoGC()` — zero-alloc, fixed `[60][60]` arrays, custom manual heap. Includes a commented-out `main()` demo. Not currently wired into game logic.

### Camera controls (runtime — `renderer/camera.go`)
- **Auto-follows player** — camera X/Y set to player X/Y each frame in `Game.Update()`
- **E** — zoom in (×1.05, centered on cursor); **Q** — zoom out (×0.95)
- Default zoom: 2.0×
- Clamp: zoom ∈ [0.2, 5.0]

### Player
- `game/player.go` — `Player` struct with `float64 X, Y` for smooth pixel movement. WASD + arrow keys, speed=2 px/frame.
- Spawns at first room's center tile (pixel coordinates: `CenterX()*16+8`, `CenterY()*16+8`).
- Corner-based wall collision with separate X/Y checks (wall sliding). Sprite: 12×12 gold rect.
- Wired into `Game.Update()` and `Game.Draw()` with camera transform.
- Drawing delegated to `renderer/drawPlayer.go` — `DrawPlayer()` function applies camera transform.

---

## Commands

| Command | What |
|---------|------|
| `go run .` | Run the game |
| `go build -o dungeonhive.exe .` | Build binary |
| `go vet ./...` | Static analysis (no custom linter config exists) |

No tests, no CI, no formatter config, no task runner, no `.gitignore` yet.

---

## Things an agent will get wrong without this

- **Run from repo root** — `assets/tileset.png` is loaded via relative `os.Open("assets/tileset.png")`. Don't `cd game/` and run.
- **Dot-import pattern** — adding new code to `game/` or `utils/` should match the `import . "DungeonHive/config"` style for consistency.
- **Fixed-size arrays** — all maps are `[60][60]int`, never `[][]int`. Passing by pointer (`*[GridSize][GridSize]int`) is the convention.
- **`game` and `player` are the same package** — `player.go` is in `package game`. Don't split into a separate `player` package.
- **No Ebitengine `Update()` input handling yet** — the game tick only increments `Ticks` and calls `Camera.Update()`. Input/handle logic should go in `Game.Update()`.
- **`astarNoGC.go` comment block** — lines 262–351 contain a commented-out `main()` demo with its own package-level code. Don't mistake it for live code.
