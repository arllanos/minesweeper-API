package main

import (
	"log"
	"os"

	"github.com/arllanos/minesweeper-API/internal/api/handler"
	"github.com/arllanos/minesweeper-API/internal/api/router"
	"github.com/arllanos/minesweeper-API/internal/repository"
	"github.com/arllanos/minesweeper-API/internal/services"
)

const defaultPort = "8080"

func main() {
	// initialize dependencies
	gameRepository := repository.NewRedisRepository()
	gameService := services.NewGameService(gameRepository)
	gameController := handler.NewGameController(gameService)
	httpRouter := router.NewChiRouter()

	// register routes
	httpRouter.POST("/users", gameController.CreateUser)
	httpRouter.PUT("/games", gameController.CreateGame)
	httpRouter.POST("/games/{gamename}/{username}/click", gameController.ClickCell)
	httpRouter.GET("/games/{gamename}/{username}/board", gameController.GetBoard)

	// start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	if err := httpRouter.SERVE(port); err != nil {
		log.Fatalf("Failed to start server %v", err)
	}
}
