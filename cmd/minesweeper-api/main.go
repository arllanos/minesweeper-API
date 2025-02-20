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
	gameHandler := handler.NewGameHandler(gameService)
	httpRouter := router.NewChiRouter()

	// register routes
	httpRouter.POST("/users", gameHandler.CreateUser)
	httpRouter.PUT("/games", gameHandler.CreateGame)
	httpRouter.POST("/games/{gamename}/{username}/click", gameHandler.ClickCell)
	httpRouter.GET("/games/{gamename}/{username}/board", gameHandler.GetBoard)

	// start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	if err := httpRouter.SERVE(port); err != nil {
		log.Fatalf("Failed to start server %v", err)
	}
}
