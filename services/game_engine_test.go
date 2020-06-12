package services

import (
	"testing"

	"github.com/arllanos/minesweeper-API/types"
	"github.com/stretchr/testify/assert"
)

func TestClickCellOutOfBounds(t *testing.T) {
	var game types.Game
	game.Name = "TestGame1"
	game.Username = "TestPlayer1"
	game.Rows = 4
	game.Cols = 4
	generateBoard(&game)

	err := clickCell(&game, 4, 4)

	assert.NotNil(t, err)
}

func TestClickCellWithinBounds(t *testing.T) {
	var game types.Game
	game.Name = "TestGame1"
	game.Username = "TestPlayer1"
	game.Rows = 4
	game.Cols = 4
	generateBoard(&game)

	err := clickCell(&game, 3, 3)

	assert.Nil(t, err)
}
