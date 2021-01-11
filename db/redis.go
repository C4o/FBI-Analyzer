package db

import (
	"FBI-Analyzer/logger"
	"os"
	"time"

	"github.com/go-redis/redis"
)

var (
	RedSess *redis.Client
)

type Redis struct {
	RedisAddr string
	RedisPass string
	RedisDB   int
}

func (r *Redis) Conn() {

	RedSess = redis.NewClient(&redis.Options{
		DB:         r.RedisDB,
		Addr:       r.RedisAddr,
		Password:   r.RedisPass,
		PoolSize:   10,
		MaxConnAge: 600 * time.Second,
	})
}

func (r *Redis) Health() {

	var err error
	s5 := time.NewTicker(5 * time.Second)

	defer func() {
		RedSess.Close()
		s5.Stop()
	}()

	for {
		select {
		case <-s5.C:
			_, err = RedSess.Ping().Result()
			if err != nil {
				logger.Print(logger.ERROR, "error in ping redis : %v", err)
				os.Exit(0)
			}
		}
	}
}
