package services

import "github.com/arllanos/minesweeper-API/internal/domain"

type GameRepository interface {
	SaveGame(game *domain.Game) (*domain.Game, error)
	SaveUser(game *domain.User) (*domain.User, error)
	GetGame(key string) (*domain.Game, error)
	GetUser(key string) (*domain.User, error)
	Exists(key string) bool
	Delete(key string) error
}
