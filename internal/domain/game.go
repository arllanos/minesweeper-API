package domain

import "time"

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
}

type User struct {
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"createdAt"`
}

type ClickData struct {
	Row  int    `json:"row"`
	Col  int    `json:"col"`
	Kind string `json:"kind"`
}
