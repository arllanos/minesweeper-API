package main

import (
	"os"

	"github.com/arllanos/minesweeper-API/internal/api"
	"github.com/arllanos/minesweeper-API/internal/api/router"
	"github.com/arllanos/minesweeper-API/internal/repository"
	"github.com/arllanos/minesweeper-API/internal/services"
)

const defaultPort = "8080"

var (
	gameRepository repository.GameRepository = repository.NewRedisRepository()
	gameService    services.GameService      = services.NewGameService(gameRepository)
	httpRouter     router.Router             = router.NewChiRouter()
	gameController api.GameController        = api.NewGameController(gameService)
)

func main() {
	httpRouter.POST("/users", gameController.CreateUser)
	httpRouter.PUT("/games", gameController.CreateGame)
	httpRouter.POST("/games/{gamename}/{username}/click", gameController.ClickCell)
	httpRouter.GET("/games/{gamename}/{username}/board", gameController.GetBoard)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	httpRouter.SERVE(port)
}
