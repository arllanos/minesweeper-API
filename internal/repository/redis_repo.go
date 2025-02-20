package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/arllanos/minesweeper-API/internal/domain"
	"github.com/arllanos/minesweeper-API/internal/services"
	"github.com/gomodule/redigo/redis"
)

const BoardSuffix = "-Board"

var (
	ErrDeleteBoard   = errors.New("error deleting game board")
	ErrMarshalData   = errors.New("unable to marshal data")
	ErrUnmarshalData = errors.New("unable to unmarshal data")
	ErrGameNotFound  = errors.New("game not found")
)

type redisRepo struct {
	pool *redis.Pool
}

func NewRedisRepository() services.GameRepository {
	return &redisRepo{
		pool: newRedisPool(),
	}
}

func (r *redisRepo) getConn() redis.Conn {
	return r.pool.Get()
}

func (r *redisRepo) SaveGame(game *domain.Game) (*domain.Game, error) {
	conn := r.getConn()
	defer conn.Close()

	k := game.Name + BoardSuffix
	if err := r.Delete(k); err != nil {
		return nil, ErrDeleteBoard
	}
	if err := r.saveBoard(k, game.Board); err != nil {
		return nil, err
	}

	jData, err := json.Marshal(game)
	if err != nil {
		log.Printf("Error: Unable to marshal game data: %q", err)
		return nil, ErrMarshalData
	}

	_, err = conn.Do("SET", game.Name, jData)
	return game, err
}

func (r *redisRepo) GetUser(key string) (*domain.User, error) {
	conn := r.getConn()
	defer conn.Close()

	data, err := redis.String(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}

	var user domain.User
	if err := json.Unmarshal([]byte(data), &user); err != nil {
		return nil, ErrUnmarshalData
	}

	return &user, nil
}

func (r *redisRepo) SaveUser(user *domain.User) (*domain.User, error) {
	conn := r.getConn()
	defer conn.Close()

	jData, err := json.Marshal(user)
	if err != nil {
		log.Printf("Error: Unable to marshal data: %q", err)
		return nil, ErrMarshalData
	}

	_, err = conn.Do("SET", user.Username, jData)
	return user, err
}

func (r *redisRepo) GetGame(key string) (*domain.Game, error) {
	conn := r.getConn()
	defer conn.Close()

	if !r.Exists(key) {
		return nil, ErrGameNotFound
	}

	data, err := redis.String(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}

	var game domain.Game
	if err := json.Unmarshal([]byte(data), &game); err != nil {
		return nil, ErrUnmarshalData
	}

	k := key + BoardSuffix
	board, err := r.readBoard(k)
	if err != nil {
		return nil, err
	}

	game.Board = board
	return &game, nil
}

func (r *redisRepo) Exists(key string) bool {
	conn := r.getConn()
	defer conn.Close()

	data, err := redis.Int(conn.Do("EXISTS", key))
	if err != nil {
		return false
	}

	return data > 0
}

func (r *redisRepo) Delete(key string) error {
	conn := r.getConn()
	defer conn.Close()

	_, err := redis.Int(conn.Do("DEL", key))
	return err
}

func (r *redisRepo) readBoard(key string) ([][]byte, error) {
	conn := r.getConn()
	defer conn.Close()

	sData, err := redis.String(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}

	var slcData [][]byte
	if err := json.Unmarshal([]byte(sData), &slcData); err != nil {
		return nil, ErrUnmarshalData
	}

	return slcData, nil
}

func (r *redisRepo) saveBoard(key string, data [][]byte) error {
	conn := r.getConn()
	defer conn.Close()

	boardData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error: Unable to marshal board data: %q", err)
		return ErrMarshalData
	}

	_, err = conn.Do("SET", key, boardData)
	return err
}

func newRedisPool() *redis.Pool {
	redisURL := os.Getenv("REDIS_URL")
	return &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.DialURL(fmt.Sprintf("redis://%s", redisURL))
		},
	}
}
