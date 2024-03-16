package queue

import (
	"fmt"
	"os"
	"strconv"

	"github.com/codersidprogrammer/gonotif/pkg/queue/process"
	"github.com/gojek/work"
	"github.com/gomodule/redigo/redis"
)

type Queue struct {
	nameSpace string
	enq       *work.Enqueuer
}

var RedisPool = &redis.Pool{
	MaxActive: 5,
	MaxIdle:   5,
	Wait:      true,
	Dial: func() (redis.Conn, error) {
		port, err := strconv.Atoi(os.Getenv("REDIS_PORT"))
		if err != nil {
			return nil, err
		}
		host := os.Getenv("REDIS_HOST")

		url := fmt.Sprintf("%s:%d", host, port)
		c, err := redis.Dial("tcp", url)
		if err != nil {
			return nil, err
		}
		return c, err
	},
}

func NewQueue(namespace string) *Queue {
	return &Queue{
		nameSpace: namespace,
		enq:       work.NewEnqueuer(namespace, RedisPool),
	}
}

func (q *Queue) Register(job string, data map[string]interface{}) (string, error) {
	j, err := q.enq.Enqueue(job, data)
	if err != nil {
		return "", err
	}

	return j.ID, nil
}

func (q *Queue) RegisterWithSecondDelay(job string, seconds int, data map[string]interface{}) (string, error) {
	j, err := q.enq.EnqueueIn(job, int64(seconds), data)
	if err != nil {
		return "", err
	}

	return j.ID, nil
}

func QueueHandler() {
	pp := process.NewProcessPool()

	// Ensuring Pool handler active before registering
	go pp.PoolsHandler()

	wp := work.NewWorkerPool(process.SendNotifcation{}, 10, "development_test", RedisPool)
	wp2 := work.NewWorkerPool(process.SendMail{}, 10, "development_test", RedisPool)

	// Registering 1 process
	sn := process.NewSendNotifcationProcess("send_notifcation", wp)
	pp.RegisterProcess(sn)

	// Registering 2 process
	sn2 := process.NewSendMailProcess("send_email", wp2)
	pp.RegisterProcess(sn2)

}
