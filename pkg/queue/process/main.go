package process

import (
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2/log"
	"github.com/gojek/work"
)

// For closing signal handlers
var sigint = make(chan os.Signal, 1)

// All processor should implement this interface
type Processor interface {

	// This function was main function
	// to implement process that should be queued
	Do()

	// Returning processor name
	ProcessName() string

	// Accessing worker pool
	GetWorkerPool() *work.WorkerPool
}

type ProcessPool struct {
	process         map[*work.WorkerPool]Processor
	registerProcess chan workpoolProcess
}

type workpoolProcess struct {
	worker  *work.WorkerPool
	process Processor
}

func NewProcessPool() *ProcessPool {
	return &ProcessPool{
		process:         make(map[*work.WorkerPool]Processor),
		registerProcess: make(chan workpoolProcess),
	}
}

func (p *ProcessPool) RegisterProcess(ps Processor) {
	log.Info("Registering pool processors ", ps.ProcessName())
	p.registerProcess <- workpoolProcess{
		worker:  ps.GetWorkerPool(),
		process: ps,
	}
}

func (p *ProcessPool) GetProcess(pool *work.WorkerPool) Processor {
	return p.process[pool]
}

func (p *ProcessPool) PoolsHandler() {
	go func() {
		signal.Notify(sigint, os.Interrupt)
	}()

	for {
		select {
		case newPool := <-p.registerProcess:
			p.process[newPool.worker] = newPool.process

			for mp, po := range p.process {
				if !mp.Started() {
					po.Do()
					mp.Start()
					log.Infof("Process %s started", po.ProcessName())
				}
			}

		case <-sigint:
			for pool := range p.process {
				if pool.Started() {
					log.Infof("Stopping process %v \n", pool)
					pool.Stop()
				}
			}
			return
		}
	}
}
