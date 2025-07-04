package game

import (
    "math/rand"
    "time"
)

func NewGameState() *GameState {
    // Example: 13 rows x 15 cols, fixed walls, random blocks
    rows, cols := 13, 15
    level := make([][]int, rows)
    for r := 0; r < rows; r++ {
        level[r] = make([]int, cols)
        for c := 0; c < cols; c++ {
            if r == 0 || c == 0 || r == rows-1 || c == cols-1 || (r%2 == 0 && c%2 == 0) {
                level[r][c] = 1 // wall
            } else {
                level[r][c] = 0 // empty
            }
        }
    }
    // Place random blocks (2), but not in corners (safe zones)
    rand.Seed(time.Now().UnixNano())
    for r := 1; r < rows-1; r++ {
        for c := 1; c < cols-1; c++ {
            if level[r][c] == 0 && !isCorner(r, c, rows, cols) && rand.Float64() < 0.3 {
                level[r][c] = 2 // block
            }
        }
    }
    return &GameState{
        Phase:   "waiting",
        Timer:   0,
        Level:   level,
        Players: make(map[string]*Player),
        Bombs:   []Bomb{},
        PowerUps: []PowerUp{},
    }
}

func isCorner(r, c, rows, cols int) bool {
    corners := [][2]int{{1, 1}, {1, cols-2}, {rows-2, 1}, {rows-2, cols-2}}
    for _, corner := range corners {
        if corner[0] == r && corner[1] == c {
            return true
        }
    }
    return false
}
