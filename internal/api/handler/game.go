package handler

import (
	"encoding/json"
	"net/http"

	"github.com/arllanos/minesweeper-API/internal/domain"
	"github.com/arllanos/minesweeper-API/internal/errors"
	"github.com/arllanos/minesweeper-API/internal/services"
)

type handler struct {
	gameService services.GameService
}

type GameHandler interface {
	CreateUser(response http.ResponseWriter, request *http.Request)
	CreateGame(response http.ResponseWriter, request *http.Request)
	ClickCell(response http.ResponseWriter, request *http.Request)
	GetBoard(response http.ResponseWriter, request *http.Request)
}

func NewGameHandler(service services.GameService) GameHandler {
	return &handler{
		gameService: service,
	}
}

func (h *handler) CreateUser(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var user domain.User
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

	result, err1 := h.gameService.CreateUser(&user)
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

func (h *handler) CreateGame(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var game domain.Game
	err := json.NewDecoder(request.Body).Decode(&game)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(errors.ServiceError{Message: "Error unmarshalling data"})
		return
	}

	result, err1 := h.gameService.CreateGame(&game)
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

func (h *handler) ClickCell(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var click domain.ClickData
	err := json.NewDecoder(request.Body).Decode(&click)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(errors.ServiceError{Message: "Error unmarshalling data"})
		return
	}

	// extract path variables from request context (router-specific logic handled externally)
	gameName := request.Context().Value("gameName").(string)
	userName := request.Context().Value("userName").(string)

	result, err1 := h.gameService.Click(gameName, userName, &click)
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

func (h *handler) GetBoard(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	// extract path variables from request context (router-specific logic handled externally)
	gameName := request.Context().Value("gameName").(string)
	userName := request.Context().Value("userName").(string)

	board, err := h.gameService.Board(gameName, userName)
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
