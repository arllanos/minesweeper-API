package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/arllanos/minesweeper-API/types"
	"github.com/gomodule/redigo/redis"
)

const BoardSuffix = "-Board"

type redisRepo struct {
	redis.Conn
}

// creates a new redis repo with a connection
func NewRedisRepository() GameRepository {
	return &redisRepo{
		getConnection(),
	}
}

func (r *redisRepo) SaveGame(game *types.Game) (*types.Game, error) {
	// reset the auxiliar board (delete & save)
	k := game.Name + BoardSuffix
	if err := r.Delete(k); err != nil {
		return nil, errors.New("Error deleting game board")
	}
	r.saveBoard(k, game.Board)

	jData, err := json.Marshal(game)
	if err != nil {
		log.Printf("Error: Unable to marshal game data: %q", err)
		return nil, err
	}
	_, err = r.Do("SET", game.Name, jData)

	return game, err
}

func (r *redisRepo) GetUser(key string) (*types.User, error) {
	data, err := redis.String(r.Do("GET", key))

	var user types.User
	err = json.Unmarshal([]byte(data), &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *redisRepo) SaveUser(user *types.User) (*types.User, error) {
	key := user.Username
	jData, err := json.Marshal(user)
	if err != nil {
		log.Printf("Error: Unable to marshal data: %q", err)
		return nil, err
	}
	_, err = r.Do("SET", key, jData)

	return user, err
}

func (r *redisRepo) GetGame(key string) (*types.Game, error) {
	data, err := redis.String(r.Do("GET", key))

	var game types.Game
	err = json.Unmarshal([]byte(data), &game)
	if err != nil {
		return nil, err
	}

	// unmarshal the board 2d slice properly from redis
	k := key + BoardSuffix
	rData, err := r.readBoard(k)
	if err != nil {
		return nil, err
	}

	log.Printf("Game Board:")
	for index, element := range rData {
		log.Printf("%d => %s", index, string(element))
	}
	game.Board = rData

	return &game, nil
}

func (r *redisRepo) Exists(key string) bool {
	data, err := redis.Int(r.Do("EXISTS", key))

	if err != nil {
		return false
	}

	return data > 0
}

func (r *redisRepo) Delete(key string) error {
	_, err := redis.Int(r.Do("DEL", key))

	return err
}

func (r *redisRepo) readBoard(key string) ([][]byte, error) {
	// redis to string
	sData, err := redis.String(r.Do("GET", key))
	if err != nil {
		return nil, err
	}

	// string to 2D slice (unmarshal)
	var slcData [][]byte
	err = json.Unmarshal([]byte(sData), &slcData)
	if err != nil {
		return nil, err
	}

	return slcData, nil
}

func (r *redisRepo) saveBoard(key string, data [][]byte) error {
	// board 2D slice to json
	boardData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error: Unable to marshal board data: %q", err)
	}

	// json encoded to redis
	_, err = redis.String(r.Do("SET", key, boardData))

	return err
}

func getConnection() redis.Conn {
	localURL := os.Getenv("REDIS_URL")
	redisURL := fmt.Sprintf("redis://%s", localURL)

	log.Printf("Connecting to Redis %q ...", redisURL)

	c, err := redis.DialURL(redisURL)
	if err != nil {
		log.Fatal(err)
	}
	return c
}
