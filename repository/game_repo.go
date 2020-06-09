package repository

import (
	"errors"

	"github.com/arllanos/minesweeper-API/types"
)

type GameRepository struct {
	repo types.Repo
}

func NewGameRepository(repo types.Repo) *GameRepository {
	return &GameRepository{repo: repo}
}

func (r *GameRepository) GetGame(name string) (*types.Game, error) {
	game, err := r.repo.GetGame(name)
	if err != nil || game == nil {
		return nil, errors.New("Game not found")
	}

	return game, nil
}

func (r *GameRepository) CreateGame(game *types.Game) error {
	if r.repo.Exists(game.Name) && game.Status == "in_progress" {
		return errors.New("Game already exists")
	} else {
		k := game.Name + types.BoardSuffix
		if err := r.repo.Delete(k); err != nil {
			return errors.New("Error deleting game board")
		}
	}

	if err := r.repo.SaveGame(game); err != nil {
		return errors.New("Error saving Game")
	}

	return nil
}

func (r *GameRepository) UpdateGame(game *types.Game) error {
	if _, err := r.repo.GetGame(game.Name); err != nil {
		return errors.New("Game not found")
	}

	return r.repo.SaveGame(game)
}

func (r *GameRepository) GetUser(username string) (*types.User, error) {
	user, err := r.repo.GetUser(username)
	if err != nil {
		return nil, errors.New("User not found")
	}

	return user, nil
}

func (r *GameRepository) CreateUser(user *types.User) error {
	if r.repo.Exists(user.Username) {
		return errors.New("User already exists")
	}

	if err := r.repo.SaveUser(user.Username, user); err != nil {
		return errors.New("Error saving user")
	}

	return nil
}
