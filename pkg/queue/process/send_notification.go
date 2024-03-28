package process

import (
	"github.com/codersidprogrammer/gonotif/app/notification/service"
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
	log.Infof("Starting job %s:%s", job.Name, job.ID)
	return next()
}

func (s *SendNotifcation) execute(job *work.Job) error {
	topic := job.ArgString("topic")
	username := job.ArgString("username")
	payload := job.Args["payload"]

	if err := job.ArgError(); err != nil {
		return err
	}

	log.Debugf("Executing topic: %s", topic)
	ps := service.NewPushService()

	if err := ps.Send(username, topic, payload); err != nil {
		log.Error("Error publishing payload, error: ", err)
		return err
	}

	log.Infof("Job completed successfully for %s:%s", job.Name, job.ID)
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
