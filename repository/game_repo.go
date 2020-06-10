package repository

import (
	"github.com/arllanos/minesweeper-API/types"
)

type GameRepository interface {
	SaveGame(game *types.Game) (*types.Game, error)
	SaveUser(game *types.User) (*types.User, error)
	GetGame(key string) (*types.Game, error)
	GetUser(key string) (*types.User, error)
	Exists(key string) bool
	Delete(key string) error
}
