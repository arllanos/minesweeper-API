package services

import (
	"errors"
	"math/rand"
	"time"
	"unicode"

	"github.com/arllanos/minesweeper-API/internal/domain"
)

// create local random number genrator
var randg = rand.New(rand.NewSource(time.Now().UnixNano()))

func generateBoard(game *domain.Game) {
	// initialize board
	game.Board = make([][]byte, game.Rows)
	for i := range game.Board {
		game.Board[i] = make([]byte, game.Cols)
		for j := 0; j < game.Cols; j++ {
			game.Board[i][j] = 'E'
		}
	}

	// plant mines randomly
	i := 0
	for i < game.Mines {
		x := randg.Intn(game.Rows)
		y := randg.Intn(game.Cols)
		if game.Board[x][y] != 'M' {
			game.Board[x][y] = 'M'
			i++
		}
	}
}

func clickCell(game *domain.Game, i int, j int) error {
	// NW, N, NE, SE, S, SW, W, E direction vectors
	dirVector := [8][2]int{{-1, -1}, {-1, 0}, {-1, 1}, {1, 1}, {1, 0}, {1, -1}, {0, -1}, {0, 1}}
	ASCII0 := 48

	var solve func(board [][]byte, r int, c int)
	solve = func(board [][]byte, r int, c int) {
		// check neighboring cells and compute mineCount
		mineCount := 0
		for i := 0; i < 8; i++ {
			x, y := r+dirVector[i][0], c+dirVector[i][1]
			if x >= 0 && x < len(board) && y >= 0 && y < len(board[0]) && (board[x][y] == 'M' || board[x][y] == 'm') {
				mineCount++
			}
		}
		if mineCount > 0 {
			// reveal cell with neighbor mine count
			board[r][c] = byte(mineCount + ASCII0)
			return
		}

		// reveal cell (no adjacent mines)
		board[r][c] = 'B'

		// recursively solve adjacent
		for i := 0; i < 8; i++ {
			x, y := r+dirVector[i][0], c+dirVector[i][1]
			if x >= 0 && x < len(board) && y >= 0 && y < len(board[0]) && (board[x][y] == 'E' || board[x][y] == 'e') {
				solve(board, x, y)
			}
		}
	}

	if !(i >= 0 && i < game.Rows && j >= 0 && j < game.Cols) {
		return errors.New("clicked cell out of bounds")
	}

	// return if it is a flagged cell
	if game.Board[i][j] == 'm' || game.Board[i][j] == 'e' {
		return errors.New("clicked cell is flagged")
	}

	// increment click count if it is a valid click
	if game.Board[i][j] == 'M' || game.Board[i][j] == 'E' {
		game.Clicks++
	}

	if game.Board[i][j] == 'M' {
		game.Board[i][j] = 'X'
		game.Status = "over"
		return nil
	}

	solve(game.Board, i, j)

	return nil
}

func flagCell(game *domain.Game, i int, j int) error {

	if !(i >= 0 && i < game.Rows && j >= 0 && j < game.Cols) {
		return errors.New("flagged cell out of bounds")
	}

	// only vealed cells M and E can be flagged / unflagged
	value := rune(game.Board[i][j])
	if value == 'M' || value == 'E' {
		game.Board[i][j] = byte(unicode.ToLower(value))
	}
	if value == 'm' || value == 'e' {
		game.Board[i][j] = byte(unicode.ToUpper(value))
	}

	return nil
}

func weHaveWinner(game *domain.Game) bool {
	// we have a winner if:	no 'E', no 'e' and no 'X'

	for i, _ := range game.Board {
		for j, _ := range game.Board[0] {
			if game.Board[i][j] == 'E' || game.Board[i][j] == 'e' || game.Board[i][j] == 'X' {
				return false
			}
		}
	}

	return true
}
