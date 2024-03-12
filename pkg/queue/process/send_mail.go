package process

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/gojek/work"
)

type SendMail struct {
	jobName string
	pool    *work.WorkerPool
}

func NewSendMailProcess(jobName string, pool *work.WorkerPool) Processor {
	return &SendMail{
		jobName: jobName,
		pool:    pool,
	}
}

func (s *SendMail) log(job *work.Job, next work.NextMiddlewareFunc) error {
	log.Info("Starting job: ", job.Name)
	return next()
}

func (s *SendMail) execute(job *work.Job) error {
	topic := job.ArgString("topic")
	message := job.Args["message"]

	if err := job.ArgError(); err != nil {
		return err
	}

	log.Info("Sending message: %s to topic: %v", message, topic)

	return nil
}

// GetWorkerPool implements Processor.
func (s *SendMail) GetWorkerPool() *work.WorkerPool {
	return s.pool
}

func (s *SendMail) Do() {
	log.Info("Running Do Job ", s.jobName)
	s.pool.Middleware(s.log)
	s.pool.Job(s.jobName, s.execute)
	log.Infof("Running Job %v", s.execute)
}

func (s *SendMail) ProcessName() string {
	return s.jobName
}
