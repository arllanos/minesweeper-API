package services

import (
	"errors"

	"github.com/arllanos/minesweeper-API/types"
)

const (
	defaultRows  = 6
	defaultCols  = 6
	defaultMines = 12
	maxRows      = 30
	maxCols      = 30
)

func Create(game *types.Game) error {
	if game.Name == "" {
		return errors.New("Game name not provided")
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
	if game.Rows > maxRows {
		game.Rows = maxRows
	}
	if game.Cols > maxCols {
		game.Cols = maxCols
	}
	if game.Mines > (game.Cols * game.Rows) {
		game.Mines = (game.Cols * game.Rows)
	}
	game.Status = "new"

	return nil
}

func Start(game *types.Game, name string) (*types.Game, error) {

	generateBoard(game)

	game.Status = "started"

	var err error
	err = nil
	return game, err
}

func Click(game *types.Game, name string, click *types.ClickData) (*types.Game, error) {

	if click.Kind == "click" {
		if err := clickCell(game, click.Row, click.Col); err != nil {
			return nil, err
		}
	} else if click.Kind == "flag" {
		if err := flagCell(game, click.Row, click.Col); err != nil {
			return nil, err
		}
	}

	return game, nil
}
