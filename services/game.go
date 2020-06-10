package services

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/arllanos/minesweeper-API/repository"
	"github.com/arllanos/minesweeper-API/types"
	"github.com/segmentio/ksuid"
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

type GameService interface {
	CreateGame(game *types.Game) (*types.Game, error)
	CreateUser(user *types.User) (*types.User, error)
	Exists(key string) bool
	Start(name string) (*types.Game, error)
	Click(name string, data *types.ClickData) (*types.Game, error)
	Board(name string) ([]uint8, error)
}

type service struct{}

var (
	repo repository.GameRepository
)

func NewGameService(db repository.GameRepository) GameService {
	repo = db
	return &service{}
}

func (*service) CreateGame(game *types.Game) (*types.Game, error) {
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

	// if no game name assign a short ID
	if game.Name == "" {
		game.Name = ksuid.New().String()
	}
	game.Status = "new"
	game.Board = nil
	game.CreatedAt = time.Now()

	return repo.SaveGame(game)
}

func (*service) CreateUser(user *types.User) (*types.User, error) {
	if user.Username == "" {
		return nil, errors.New("Username not provided")
	}
	return repo.SaveUser(user)
}

func (*service) Exists(key string) bool {
	return repo.Exists(key)
}

func (*service) Start(name string) (*types.Game, error) {
	game, err := repo.GetGame(name)
	if err != nil {
		return nil, err
	}

	game.Status = "started"
	game.StartedAt = time.Now()

	generateBoard(game)

	_, err = repo.SaveGame(game)

	return game, err
}

func (*service) Click(name string, click *types.ClickData) (*types.Game, error) {
	game, err := repo.GetGame(name)
	if err != nil {
		return nil, err
	}

	game.TimeSpent = game.StartedAt.Sub(time.Now())

	if click.Kind == "click" {
		if err := clickCell(game, click.Row, click.Col); err != nil {
			return nil, err
		}
	} else if click.Kind == "flag" {
		if err := flagCell(game, click.Row, click.Col); err != nil {
			return nil, err
		}
	}

	if _, err := repo.SaveGame(game); err != nil {
		return nil, err
	}

	return game, nil
}

func (*service) Board(name string) ([]uint8, error) {
	game, err := repo.GetGame(name)
	if err != nil {
		return nil, err
	}

	var boardToJSON func(data [][]byte) ([]uint8, error)
	boardToJSON = func(data [][]byte) ([]uint8, error) {
		tmp := make([][]string, len(data), len(data[0]))
		for i, _ := range data {
			row := data[i][:]
			var fRow []string
			for _, v := range row {
				fRow = append(fRow, string(v))
			}
			tmp[i] = fRow
		}
		tmpJson, err := json.Marshal(tmp)
		if err != nil {
			return nil, errors.New("Cannot encode to JSON")
		}
		return tmpJson, nil
	}

	jBoard, err1 := boardToJSON(game.Board)
	if err1 != nil {
		return nil, err1
	}
	return jBoard, nil
}
