package services

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/arllanos/minesweeper-API/internal/domain"
	"github.com/arllanos/minesweeper-API/internal/repository"
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
	CreateGame(game *domain.Game) (*domain.Game, error)
	CreateUser(user *domain.User) (*domain.User, error)
	Exists(key string) bool
	Click(gameName string, userName string, data *domain.ClickData) (*domain.Game, error)
	Board(gameName string, userName string) ([]uint8, error)
}

type service struct {
	repo repository.GameRepository
}

func NewGameService(db repository.GameRepository) GameService {
	return &service{repo: db}
}

func (s *service) CreateGame(game *domain.Game) (*domain.Game, error) {

	if !s.repo.Exists(game.Username) {
		return nil, errors.New("user_not_found")
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
	_, err := s.repo.SaveGame(game)

	if err != nil {
		return nil, errors.New("error saving game")
	}

	return game, err
}

func (s *service) CreateUser(user *domain.User) (*domain.User, error) {
	if s.repo.Exists(user.Username) {
		return nil, errors.New("user_already_exist")
	}

	user.CreatedAt = time.Now()
	return s.repo.SaveUser(user)
}

func (s *service) Exists(key string) bool {
	return s.repo.Exists(key)
}

func (s *service) Click(gameName string, userName string, click *domain.ClickData) (*domain.Game, error) {
	if !s.repo.Exists(gameName) {
		return nil, errors.New("game_not_found")
	}
	if !s.repo.Exists(userName) {
		return nil, errors.New("user_not_found")
	}

	game, err := s.repo.GetGame(gameName)
	if err != nil {
		return nil, err
	}

	log.Printf("Click type [%s] request at (%d, %d) for game [%s] with status [%s]", click.Kind, click.Row, click.Col, game.Name, game.Status)

	if click.Kind != "click" && click.Kind != "flag" {
		return nil, errors.New("bad_click_kind")
	}

	if game.Status == "ready" {
		// first click: set in progress and set start time
		game.Status = "in_progress"
		game.StartedAt = time.Now()
	}

	if game.Status == "over" {
		return nil, errors.New("game_over")
	}

	if game.Status == "won" {
		return nil, errors.New("game_won")
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

	game.TimeSpent = time.Since(game.StartedAt)

	if weHaveWinner(game) {
		game.Status = "won"
	}

	if _, err := s.repo.SaveGame(game); err != nil {
		return nil, err
	}

	return game, nil
}

func (s *service) Board(gameName string, userName string) ([]uint8, error) {
	if !s.repo.Exists(gameName) {
		return nil, errors.New("game_not_found")
	}
	if !s.repo.Exists(userName) {
		return nil, errors.New("user_not_found")
	}

	game, err := s.repo.GetGame(gameName)
	if err != nil {
		return nil, err
	}

	if game.Board == nil {
		return nil, errors.New("this game has no board")
	}

	boardToJSON := func(data [][]byte) ([]uint8, error) {
		tmp := make([][]string, len(data), len(data[0]))
		for i := range data {
			row := data[i][:]
			var fRow []string
			for _, v := range row {
				fRow = append(fRow, string(v))
			}
			tmp[i] = fRow
		}
		tmpJSON, err := json.Marshal(tmp)
		if err != nil {
			return nil, errors.New("cannot encode to json")
		}
		return tmpJSON, nil
	}

	jBoard, err1 := boardToJSON(game.Board)
	if err1 != nil {
		return nil, err1
	}
	return jBoard, nil
}
