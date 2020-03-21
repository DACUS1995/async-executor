package asyncexecutor

import (
	"sync"
)

type executor struct {
	wg                  *sync.WaitGroup
	globalJobQueue      chan *Job
	GlobalResponseQueue chan *ResponseObject
	workerCounter       int
	jobCounter          int
	workers             []*Worker
}

type Executor interface {
	StartExecutor(int)
	AddGlobalJob(Job)
	StopExecutor()
	waitAllWorkers(chan Job)
}

func NewExecutor(queueSize int) *executor {
	return &executor{
		&sync.WaitGroup{},
		make(chan *Job, queueSize),
		make(chan *ResponseObject, queueSize),
		0, 0,
		[]*Worker{},
	}
}

func (exec *executor) StartExecutor(numWorkers int) {
	for i := 0; i < numWorkers; i++ {
		worker := NewWorker(exec.workerCounter, exec.globalJobQueue, exec.GlobalResponseQueue)
		exec.workers = append(exec.workers, worker)
		exec.wg.Add(1)
		exec.workerCounter++
		go worker.start(exec.wg)
	}

	return
}

func (exec *executor) CreateJob(function callableType, parameters []interface{}) {
	exec.addGlobalJob(NewJob(
		exec.jobCounter,
		function,
		&parameterObject{parameters},
	))
	exec.jobCounter++
}

func (exec *executor) addGlobalJob(job *Job) {
	exec.globalJobQueue <- job
}

func (exec *executor) StopExecutor() {
	exec.waitAllWorkers()
	return
}

func (exec *executor) waitAllWorkers() {
	close(exec.globalJobQueue)

	for _, worker := range exec.workers {
		worker.stop()
	}
	exec.wg.Wait()
}
