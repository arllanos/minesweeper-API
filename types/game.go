package types

import "time"

/* Values in Game.Board cells follows this logic/rules:
M		Veiled Mine.
E		Veiled Empty.
B		Revealed Blank w/o adjacent mines.
1-8		Digit representing mine count. These are set during the game rather than when game starts.
*/

type Game struct {
	Name      string        `json:"name"`
	Username  string        `json:"username"`
	Rows      int           `json:"rows"`
	Cols      int           `json:"cols"`
	Mines     int           `json:"mines"`
	Status    string        `json:"status"`
	Board     [][]byte      `json:"board"`
	Clicks    int           `json:"clicks"`
	CreatedAt time.Time     `json:"created_at,omitempty"`
	StartedAt time.Time     `json:"started_at"`
	TimeSpent time.Duration `json:"time_spent"`
	Points    float32       `json:"points,omitempty"`
}

type User struct {
	Username  string    `json:"username"`
	CreatedOn time.Time `json:"createdOn"`
}

type ClickData struct {
	Row  int    `json:"row"`
	Col  int    `json:"col"`
	Kind string `json:"kind"`
}
