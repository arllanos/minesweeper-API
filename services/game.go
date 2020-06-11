package services

import (
	"encoding/json"
	"errors"
	"log"
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
	minRows      = 2
	minCols      = 2
)

type GameService interface {
	CreateGame(game *types.Game) (*types.Game, error)
	CreateUser(user *types.User) (*types.User, error)
	Exists(key string) bool
	Click(gameName string, userName string, data *types.ClickData) (*types.Game, error)
	Board(gameName string, userName string) ([]uint8, error)
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

	if !repo.Exists(game.Username) {
		return nil, errors.New("Username does not exits")
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

	// if no game name assign a short ID
	if game.Name == "" {
		game.Name = ksuid.New().String()
	}
	game.Board = nil
	game.CreatedAt = time.Now()

	// start the game with an initialized board
	game.Status = "ready"
	generateBoard(game)
	_, err := repo.SaveGame(game)

	if err != nil {
		return nil, errors.New("Error saving game")
	}

	return game, err
}

func (*service) CreateUser(user *types.User) (*types.User, error) {
	if user.Username == "" {
		return nil, errors.New("Username not provided")
	}

	if repo.Exists(user.Username) {
		return nil, errors.New("User already exists")
	}

	user.CreatedAt = time.Now()
	return repo.SaveUser(user)
}

func (*service) Exists(key string) bool {
	return repo.Exists(key)
}

func (*service) Click(gameName string, userName string, click *types.ClickData) (*types.Game, error) {
	if !repo.Exists(gameName) || !repo.Exists(userName) {
		return nil, errors.New("Game or user do not exists")
	}

	game, err := repo.GetGame(gameName)
	if err != nil {
		return nil, err
	}

	log.Printf("Click type [%s] request at (%d, %d) for game [%s] with status [%s]", click.Kind, click.Row, click.Col, game.Name, game.Status)

	if click.Kind != "click" && click.Kind != "flag" {
		return nil, errors.New("Click kind should be either 'click' or 'flag'")
	}

	if game.Status == "ready" {
		// first click: set in progress and set start time
		game.Status = "in_progress"
		game.StartedAt = time.Now()
	}

	if game.Status == "over" {
		return nil, errors.New("Game is over and does not accept clicks")
	}

	if game.Status == "won" {
		return nil, errors.New("Game is finished and does not accept clicks. You won!!!")
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

	game.TimeSpent = time.Now().Sub(game.StartedAt)

	if weHaveWinner(game) {
		game.Status = "won"
	}

	if _, err := repo.SaveGame(game); err != nil {
		return nil, err
	}

	return game, nil
}

func (*service) Board(gameName string, userName string) ([]uint8, error) {
	if !repo.Exists(gameName) || !repo.Exists(userName) {
		return nil, errors.New("Game or user do not exists")
	}

	game, err := repo.GetGame(gameName)
	if err != nil {
		return nil, err
	}

	if game.Board == nil {
		return nil, errors.New("This game has no board.")
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
