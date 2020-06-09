package repository

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/arllanos/minesweeper-API/internal/logs"
	"github.com/arllanos/minesweeper-API/types"
	"github.com/gomodule/redigo/redis"
	"go.uber.org/zap"
)

type repo struct {
	redis.Conn
}

// creates a new redis repo with a connection
func NewRedisRepo() types.Repo {
	return &repo{
		getConnection(),
	}
}

func (r *repo) SaveGame(game *types.Game) error {
	key := game.Name

	// if status is new do not serialize the board
	if game.Status != "new" {
		k := key + types.BoardSuffix
		r.saveBoard(k, game.Board)
	}

	jData, err := json.Marshal(game)
	if err != nil {
		logs.Log().Error("Unable to marshal game data", zap.Error(err))
		return err
	}
	_, err = r.Do("SET", key, jData)

	return err
}

func (r *repo) SaveUser(key string, user *types.User) error {
	jData, err := json.Marshal(user)
	if err != nil {
		logs.Log().Error("Unable to marshal data", zap.Error(err))
		return err
	}
	_, err = r.Do("SET", key, jData)

	return err
}

func (r *repo) GetGame(key string) (*types.Game, error) {
	data, err := redis.String(r.Do("GET", key))
	logs.Log().Infof("Game -> %s", data)

	var game types.Game
	err = json.Unmarshal([]byte(data), &game)
	if err != nil {
		return nil, err
	}

	if game.Status != "new" {
		// unmarshal the board 2d slice properly from redis
		k := key + types.BoardSuffix
		rData, err := r.readBoard(k)
		if err != nil {
			return nil, err
		}

		logs.Log().Debug("Game Board:")
		for index, element := range rData {
			logs.Log().Debugf("%d => %s", index, string(element))
		}
		game.Board = rData
	}

	return &game, nil
}

func (r *repo) GetUser(key string) (*types.User, error) {
	data, err := redis.String(r.Do("GET", key))
	logs.Log().Infof("User -> %s", data)

	var user types.User
	err = json.Unmarshal([]byte(data), &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *repo) Exists(key string) bool {
	data, err := redis.Int(r.Do("EXISTS", key))

	if err != nil {
		return false
	}

	return data > 0
}

func (r *repo) Delete(key string) error {
	_, err := redis.Int(r.Do("DEL", key))

	return err
}

func (r *repo) readBoard(key string) ([][]byte, error) {
	// redis to string
	sData, err := redis.String(r.Do("GET", key))
	logs.Log().Infof("Board -> %s", sData)
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

func (r *repo) saveBoard(key string, data [][]byte) error {
	// board 2D slice to json
	boardData, err := json.Marshal(data)
	if err != nil {
		logs.Log().Error("Unable to marshal board data", zap.Error(err))
	}

	// json encoded to redis
	_, err = redis.String(r.Do("SET", key, boardData))

	return err
}

func getConnection() redis.Conn {
	localURL := os.Getenv("REDIS_URL")
	redisURL := fmt.Sprintf("redis://%s", localURL)

	logs.Log().Infof("Connecting to Redis %s ...", redisURL)

	c, err := redis.DialURL(redisURL)
	if err != nil {
		logs.Log().Fatal("Error connecting to Redis")
		panic(err)
	}
	return c
}
