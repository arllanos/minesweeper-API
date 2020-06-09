package main

import (
	"github.com/arllanos/minesweeper-API/internal/logs"
	"github.com/arllanos/minesweeper-API/repository"
	"github.com/arllanos/minesweeper-API/services"
	"github.com/arllanos/minesweeper-API/types"
)

func main() {
	logs.InitLogger()

	type Services struct {
		gameService services.GameService
	}
	db := repository.NewRedisRepo()
	services := &Services{
		gameService: services.NewGameService(db),
	}

	game := &types.Game{
		Name:     "TestGame",
		Username: "ariel",
		Rows:     10,
		Cols:     10,
		Mines:    7,
	}
	services.gameService.Create(game)
	services.gameService.Start("TestGame")

	clickData := &types.ClickData{
		Row:  1,
		Col:  1,
		Kind: "click",
	}

	services.gameService.Click("TestGame", clickData)
}
