# DungeonHive — AGENTS.md

> **語言設定：請全程使用繁體中文（臺灣）回應。**

**Ebitengine v2** roguelike-dungeon generator.  
Build: `go run ./cmd/dungeonhive` (or `go run .` from repo root). Go 1.26, single module `DungeonHive`.

---

## Architecture (16 Go files, 6 packages)

| Directory | Role | Entrypoint |
|-----------|------|------------|
| `cmd/dungeonhive/main.go` | Ebitengine game loop, 960×960 window | `game.NewGame()` → `ebiten.RunGame` |
| `internal/game/` | Game state, player, Update/Draw orchestration | `Game` struct, `Update()` / `Draw()` / `Layout()` |
| `internal/dungeon/` | Map generation, autotiling, room/corridor placement | `dungeon.GenerateMapAndRooms()`, `dungeon.BuildTileMap()` |
| `internal/renderer/` | Camera (WASD+Q/E), resource loading, tile-map blit | `Renderer`, `Camera`, `ResourceManager` |
| `internal/pathfinding/` | A* pathfinding (standard + no-GC fixed-array) | `Search()` and `AStarFixedContext.FindPathFixedNoGC()` |
| `pkg/config/` | Constants: `GridSize=60`, `TileSize=16`, `DirType`, `Point` | `constant.go` |
| `assets/` | `tileset.png` (13×5 grid, 16px tiles) + 10 overworld sprites | Loaded via `assets/tileset.png` relative path |

**Global state** — `internal/game/game.go` declares `rmap` and `tmap` as package-level `[60][60]int` arrays; `Game` holds pointers to them. No other global mutable state.

---

## Key codebase facts

### Dot-import convention
- `internal/dungeon/`, `internal/game/`, `internal/pathfinding/` use `import . "DungeonHive/pkg/config"` — all config identifiers (GridSize, Point, DirType, TileSize, etc.) in scope without prefix.
- `internal/renderer/` uses `import cf "DungeonHive/pkg/config"` — prefixed with `cf.` instead.

### Map generation flow
```
dungeon.DefaultConfig() → GenerateMapAndRooms() → PlaceRooms() → ConnectRooms() (Prim MST + A*)
                                                  → BuildTileMap() (autotile wall-ID mapping)
```
- Raw map: `0` = floor, `1` = wall.
- Tile map: `0` = skip drawing, nonzero = tile index into `tileset.png` (column-major, 13 cols).
- `autotile.go` handles 8-direction neighbour-aware wall tile selection.

### A* — two implementations
1. **`internal/pathfinding/astar.go`** — standard `Search(start, end, CostFunc)`. Uses `map[Point]int` / `container/heap` — GC-allocating. Used by corridor carving.
2. **`internal/pathfinding/astarNoGC.go`** — `AStarFixedContext.FindPathFixedNoGC()` — zero-alloc, fixed `[60][60]` arrays, custom manual heap. Includes a commented-out `main()` demo (lines ~262). Not currently wired into game logic.

### Camera controls (runtime — `internal/renderer/camera.go`)
- **Auto-follows player** — camera X/Y set to player X/Y each frame in `Game.Update()`
- **E** — zoom in (×1.05, centered on cursor); **Q** — zoom out (×0.95)
- Default zoom: 2.0×
- Clamp: zoom ∈ [0.2, 5.0]

### Player
- `internal/game/player.go` — `Player` struct with smooth grid-based 8-direction movement.
- **State machine**: `MoveTimer` → 0 = idle, 1–8 = interpolating between tiles (~133ms/tile). Input locked during movement.
- **8-direction input**: WASD + arrows, diagonals from combinations (W+A → up-left, etc.).
- **Collision**: checks target grid tile against `rawMap[gy][gx] == 1` only — no bounding-box.
- Spawns at first room's center tile via `dungeon.Room.CenterX()`/`CenterY()`.
- Sprite: 12×12 gold-filled rect (placeholder — swappable for `assets/ow*.png`).
- Drawing: `renderer.DrawPlayer()` applies camera transform.

---

## Commands

| Command | What |
|---------|------|
| `go run ./cmd/dungeonhive` | Run the game |
| `go build -o dungeonhive.exe ./cmd/dungeonhive` | Build binary |
| `go vet ./...` | Static analysis |

No tests, no CI, no formatter config, no task runner, no `.gitignore` yet.

---

## Things an agent will get wrong without this

- **Run from repo root** — `assets/tileset.png` is loaded via relative `os.Open("assets/tileset.png")`. Don't `cd cmd/dungeonhive/` and run; always run from project root.
- **Dot-import pattern** — `internal/dungeon/`, `internal/game/`, `internal/pathfinding/` use `import . "DungeonHive/pkg/config"`. `internal/renderer/` uses `cf.` prefix. Keep consistency when adding files.
- **Fixed-size arrays** — all maps are `[60][60]int`, never `[][]int`. Passing by pointer (`*[GridSize][GridSize]int`) is the convention.
- **`game` and `player` are the same package** — `player.go` lives in `package game` (`internal/game/`). Don't split into a separate `player` package.
- **`astarNoGC.go` comment block** — lines ~262–351 contain a commented-out `main()` demo. Don't mistake it for live code.
- **Player movement locks input** — during interpolation (`MoveTimer > 0`), `Player.Update()` ignores new input until the tile is reached. Diagonals move to an adjacent diagonal tile, not two separate orthogonal tiles.
- **Standard Go Layout** — project follows `cmd/` (entry), `internal/` (private domain packages), `pkg/` (shared config). New domain packages belong under `internal/`.
