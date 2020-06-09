package services

import (
	"errors"
	"math/rand"
	"time"
	"unicode"

	"github.com/arllanos/minesweeper-API/types"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func generateBoard(game *types.Game) {
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
		x := rand.Intn(game.Rows)
		y := rand.Intn(game.Cols)
		if game.Board[x][y] != 'M' {
			game.Board[x][y] = 'M'
			i++
		}
	}
}

func clickCell(game *types.Game, i int, j int) error {
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

	if !(i >= 0 || i < len(game.Board) || j >= 0 || j < len(game.Board[0])) {
		return errors.New("Cell out of bounds")
	}

	// click on a flagged cell -> do nothing
	if game.Board[i][j] == 'm' || game.Board[i][j] == 'e' {
		return nil
	}

	if game.Board[i][j] == 'M' {
		game.Board[i][j] = 'X'
		game.Status = "over"
		return nil
	}

	solve(game.Board, i, j)
	return nil
}

func flagCell(game *types.Game, i int, j int) error {

	if !(i >= 0 || i < len(game.Board) || j >= 0 || j < len(game.Board[0])) {
		return errors.New("Cell out of bounds")
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
