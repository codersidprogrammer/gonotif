package process

import (
	"context"

	"github.com/codersidprogrammer/gonotif/platform/transport"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gojek/courier-go"
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
	payload := job.Args["payload"]

	if err := job.ArgError(); err != nil {
		return err
	}

	if err := transport.MqttClient.Publish(context.Background(), topic, payload, courier.QOSTwo); err != nil {
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
