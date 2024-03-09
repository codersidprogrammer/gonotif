package platform

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/codersidprogrammer/gonotif/pkg/utils"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2/log"
)

var RedisConnection *redis.Client
var lock = &sync.Mutex{}

func DoConnectRedis() (*redis.Client, error) {
	port, err := strconv.Atoi(os.Getenv("REDIS_PORT"))
	utils.ExitIfErr(err, "Failed to get redis port")

	host := os.Getenv("REDIS_HOST")

	url := fmt.Sprintf("%s:%d", host, port)
	options := &redis.Options{
		Addr: url,
	}

	if RedisConnection == nil {
		lock.Lock()
		defer lock.Unlock()

		if RedisConnection == nil {
			log.Info("Creating new redis connection")
			RedisConnection = redis.NewClient(options)
		} else {
			log.Info("Using existing redis connection")
		}
	}

	log.Info("Redis conneection established")
	return RedisConnection, nil

}
