package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/arllanos/minesweeper-API/errors"
	"github.com/arllanos/minesweeper-API/services"
	"github.com/arllanos/minesweeper-API/types"
)

type controller struct{}

var (
	gameService services.GameService
)

type GameController interface {
	CreateUser(response http.ResponseWriter, request *http.Request)
	CreateGame(response http.ResponseWriter, request *http.Request)
	StartGame(response http.ResponseWriter, request *http.Request)
	ClickCell(response http.ResponseWriter, request *http.Request)
	GetBoard(response http.ResponseWriter, request *http.Request)
}

func NewGameController(service services.GameService) GameController {
	gameService = service
	return &controller{}
}

func (*controller) CreateUser(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var user types.User
	err := json.NewDecoder(request.Body).Decode(&user)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(errors.ServiceError{Message: "Error unmarshalling data"})
		return
	}

	if gameService.Exists(user.Username) {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(errors.ServiceError{Message: "User already exists"})
		return
	}

	user.CreatedAt = time.Now()

	result, err1 := gameService.CreateUser(&user)
	if err1 != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(errors.ServiceError{Message: "Error saving user"})
		return
	}
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(result)
}

func (*controller) CreateGame(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var game types.Game
	err := json.NewDecoder(request.Body).Decode(&game)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(errors.ServiceError{Message: "Error unmarshalling data"})
		return
	}

	if !gameService.Exists(game.Username) {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(errors.ServiceError{Message: "Cannot create game for non existent username"})
		return
	}

	result, err1 := gameService.CreateGame(&game)
	if err1 != nil {
		errMsg := err1.Error()
		fmt.Println(errMsg)
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(errors.ServiceError{Message: "Error saving game"})
		return
	}
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(result)
}

func (*controller) StartGame(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	// TODO: delegate parsing of route variables to the http router
	URLpath := strings.Split(request.URL.Path, "/")
	username := URLpath[len(URLpath)-1]
	gamename := URLpath[len(URLpath)-2]

	if !gameService.Exists(gamename) {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(errors.ServiceError{Message: "Cannot start game for non existent game ID"})
		return
	}

	if !gameService.Exists(username) {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(errors.ServiceError{Message: "Cannot start game for non existent username"})
		return
	}

	result, err1 := gameService.Start(gamename)
	if err1 != nil {
		errMsg := err1.Error()
		fmt.Println(errMsg)
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(errors.ServiceError{Message: "Error starting game"})
		return
	}
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(result)
}

func (*controller) ClickCell(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var click types.ClickData
	err := json.NewDecoder(request.Body).Decode(&click)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(errors.ServiceError{Message: "Error unmarshalling data"})
		return
	}

	// TODO: delegate parsing of route variables to the http router
	URLpath := strings.Split(request.URL.Path, "/")
	username := URLpath[len(URLpath)-2]
	gamename := URLpath[len(URLpath)-3]

	fmt.Println("click cell route variables", username, gamename, URLpath)

	if !gameService.Exists(gamename) {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(errors.ServiceError{Message: "Cannot start game for non existent game ID"})
		return
	}

	if !gameService.Exists(username) {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(errors.ServiceError{Message: "Cannot start game for non existent username"})
		return
	}

	result, err1 := gameService.Click(gamename, &click)
	if err1 != nil {
		errMsg := err1.Error()
		fmt.Println(errMsg)
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(errors.ServiceError{Message: "Error starting game"})
		return
	}
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(result)
}

func (*controller) GetBoard(response http.ResponseWriter, request *http.Request) {
	//TODO: GetBoard code goes here...
	_, _ = response.Write([]byte("Printing board coming soon"))

}
