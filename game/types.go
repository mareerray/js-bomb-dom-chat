package game

type Player struct {
    ID     string `json:"id"`
    Name   string `json:"name"`
    Color  string `json:"color"`
    Row    int    `json:"row"`
    Col    int    `json:"col"`
    Lives  int    `json:"lives"`
    Bombs  int    `json:"bombs"`
    Flame  int    `json:"flame"`
    Speed  int    `json:"speed"`
    Alive  bool   `json:"alive"`
}

type Bomb struct {
    Row    int    `json:"row"`
    Col    int    `json:"col"`
    Owner  string `json:"owner"`
    Timer  int    `json:"timer"`
    Flame  int    `json:"flame"`
}

type PowerUp struct {
    Row  int    `json:"row"`
    Col  int    `json:"col"`
    Type string `json:"type"` // "bomb", "flame", "speed"
}

type GameState struct {
    Phase    string              `json:"phase"`
    Timer    int                 `json:"timer"`
    Level    [][]int             `json:"level"`
    Players  map[string]*Player  `json:"players"`
    Bombs    []Bomb              `json:"bombs"`
    PowerUps []PowerUp           `json:"powerups"`
}
