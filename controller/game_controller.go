package controller

import (
	"encoding/json"
	"net/http"
	"strings"

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

	if user.Username == "" {
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(errors.ServiceError{Message: "User name not provided"})
		return
	}

	result, err1 := gameService.CreateUser(&user)
	if err1 != nil {
		if err1.Error() == "user_already_exist" {
			response.WriteHeader(http.StatusConflict)
			json.NewEncoder(response).Encode(errors.ServiceError{Message: err1.Error()})
			return
		}
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(errors.ServiceError{Message: err1.Error()})
		return
	}
	response.WriteHeader(http.StatusCreated)
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

	result, err1 := gameService.CreateGame(&game)
	if err1 != nil {

		if err1.Error() == "user_not_found" {
			response.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(response).Encode(errors.ServiceError{Message: "Username not exists"})
			return
		}
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(errors.ServiceError{Message: err1.Error()})
		return
	}

	response.WriteHeader(http.StatusCreated)
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
	userName := URLpath[len(URLpath)-2]
	gameName := URLpath[len(URLpath)-3]

	result, err1 := gameService.Click(gameName, userName, &click)
	if err1 != nil {
		if err1.Error() == "bad_click_kind" || err1.Error() == "game_over" || err1.Error() == "game_won" {
			response.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(response).Encode(errors.ServiceError{Message: err1.Error()})
			return
		}
		response.WriteHeader(http.StatusNotFound)
		json.NewEncoder(response).Encode(errors.ServiceError{Message: err1.Error()})
		return
	}
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(result)
}

func (*controller) GetBoard(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	// TODO: delegate parsing of route variables to the http router
	URLpath := strings.Split(request.URL.Path, "/")
	userName := URLpath[len(URLpath)-2]
	gameName := URLpath[len(URLpath)-3]

	board, err := gameService.Board(gameName, userName)
	if err != nil {
		if err.Error() == "user_not_found" || err.Error() == "game_not_found" {
			response.WriteHeader(http.StatusNotFound)
			json.NewEncoder(response).Encode(errors.ServiceError{Message: "Username or game not exists"})
			return
		}
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(errors.ServiceError{Message: err.Error()})
		return
	}

	response.WriteHeader(http.StatusOK)
	_, _ = response.Write([]byte(board))

}
