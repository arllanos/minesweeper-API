package main

import (
	"os"

	"github.com/arllanos/minesweeper-API/controller"
	router "github.com/arllanos/minesweeper-API/http"
	"github.com/arllanos/minesweeper-API/repository"
	"github.com/arllanos/minesweeper-API/services"
)

const defaultPort = "8080"

var (
	gameRepository repository.GameRepository = repository.NewRedisRepository()
	gameService    services.GameService      = services.NewGameService(gameRepository)
	httpRouter     router.Router             = router.NewChiRouter()
	gameController controller.GameController = controller.NewGameController(gameService)
)

func main() {
	httpRouter.POST("/users", gameController.CreateUser)
	httpRouter.POST("/games", gameController.CreateGame)
	httpRouter.POST("/games/{gamename}/{username}", gameController.StartGame)
	httpRouter.POST("/games/{gamename}/{username}/click", gameController.ClickCell)
	httpRouter.GET("/games/{gamename}/{username}/board", gameController.GetBoard)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	httpRouter.SERVE(port)
}
