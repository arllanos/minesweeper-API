package services

import (
	"errors"
	"time"

	"github.com/arllanos/minesweeper-API/repository"
	"github.com/arllanos/minesweeper-API/types"
)

const (
	defaultRows  = 10
	defaultCols  = 10
	defaultMines = 15
	maxRows      = 30
	maxCols      = 30
	minRows      = 5
	minCols      = 5
)

type GameService struct {
	GameRepo types.GameServiceRepo
}

func NewGameService(db types.Repo) types.GameService {
	return &GameService{
		GameRepo: repository.NewGameRepository(db),
	}
}

func (s *GameService) Create(game *types.Game) error {
	if game.Name == "" {
		return errors.New("Game name not provided")
	}

	// defaults
	if game.Rows == 0 {
		game.Rows = defaultRows
	}
	if game.Rows == 0 {
		game.Rows = defaultRows
	}
	if game.Cols == 0 {
		game.Cols = defaultCols
	}
	if game.Mines == 0 {
		game.Mines = defaultMines
	}

	// maximum
	if game.Rows > maxRows {
		game.Rows = maxRows
	}
	if game.Cols > maxCols {
		game.Cols = maxCols
	}

	// minimum
	if game.Rows < minRows {
		game.Rows = minRows
	}
	if game.Cols < minCols {
		game.Cols = minCols
	}

	if game.Mines > (game.Cols * game.Rows) {
		game.Mines = (game.Cols * game.Rows)
	}

	game.Status = "new"
	game.Board = nil
	game.CreatedAt = time.Now()

	err := s.GameRepo.CreateGame(game)

	return err
}

func (s *GameService) Start(name string) (*types.Game, error) {
	game, err := s.GameRepo.GetGame(name)
	if err != nil {
		return nil, err
	}

	game.Status = "started"
	game.StartedAt = time.Now()

	generateBoard(game)

	err = s.GameRepo.UpdateGame(game)

	return game, err
}

func (s *GameService) Click(name string, click *types.ClickData) (*types.Game, error) {
	game, err := s.GameRepo.GetGame(name)
	if err != nil {
		return nil, err
	}

	if click.Kind == "click" {
		if err := clickCell(game, click.Row, click.Col); err != nil {
			return nil, err
		}
	} else if click.Kind == "flag" {
		if err := flagCell(game, click.Row, click.Col); err != nil {
			return nil, err
		}
	}

	if err := s.GameRepo.UpdateGame(game); err != nil {
		return nil, err
	}

	return game, nil
}
