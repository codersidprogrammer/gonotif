package transport

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/codersidprogrammer/gonotif/pkg/utils"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2/log"
)

var RedisClient *redis.Client

type redisTransport struct {
	*Transport
	name string
}

func NewRedisTransport(name string) TransportService {
	return &redisTransport{
		name: name,
		Transport: &Transport{
			lock: &sync.Mutex{},
		},
	}
}

// DoConnect implements TransportService.
func (t *redisTransport) DoConnect() error {
	port, err := strconv.Atoi(os.Getenv("REDIS_PORT"))
	utils.ExitIfErr(err, "Failed to get redis port")

	host := os.Getenv("REDIS_HOST")

	url := fmt.Sprintf("%s:%d", host, port)
	options := &redis.Options{
		Addr: url,
	}

	if RedisClient == nil {
		t.Transport.lock.Lock()
		defer t.Transport.lock.Unlock()

		if RedisClient == nil {
			log.Infof("[REDIS] Creating New %s Transport Connection", t.name)
			RedisClient = redis.NewClient(options)
		} else {
			log.Info("Using existing redis connection")
		}
	}

	log.Infof("[REDIS] %s connection established", t.name)
	return nil
}

func (t *redisTransport) Close() error {
	log.Info("Closing transport connection: ", t.name)
	return RedisClient.Close()
}
