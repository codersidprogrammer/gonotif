package process

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/gojek/work"
)

type SendNotifcation struct {
	jobName string
	pool    *work.WorkerPool
}

func NewSendNotifcationProcess(jobName string, pool *work.WorkerPool) Processor {
	return &SendNotifcation{
		jobName: jobName,
		pool:    pool,
	}
}

func (s *SendNotifcation) log(job *work.Job, next work.NextMiddlewareFunc) error {
	log.Info("Starting job: ", job.Name)
	return next()
}

func (s *SendNotifcation) execute(job *work.Job) error {
	topic := job.ArgString("topic")
	message := job.Args["message"]

	if err := job.ArgError(); err != nil {
		return err
	}

	log.Infof("Sending message: %s to topic: %v", message, topic)

	return nil
}

// GetWorkerPool implements Processor.
func (s *SendNotifcation) GetWorkerPool() *work.WorkerPool {
	return s.pool
}

func (s *SendNotifcation) Do() {
	s.pool.Middleware(s.log)
	s.pool.Job(s.jobName, s.execute)
}

func (s *SendNotifcation) ProcessName() string {
	return s.jobName
}
