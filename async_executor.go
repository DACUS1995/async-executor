package asyncexecutor

import (
	"sync"
)

type executor struct {
	wg                  sync.WaitGroup
	globalJobQueue      chan *Job
	globalResponseQueue chan *ResponseObject
	workerCounter       int
	jobCounter          int
}

type Executor interface {
	StartExecutor(int)
	AddGlobalJob(Job)
	StopExecutor()
	waitAllWorkers(chan Job)
}

func NewExecutor(queueSize int) *executor {
	return &executor{
		sync.WaitGroup{},
		make(chan *Job, queueSize),
		make(chan *ResponseObject, queueSize),
		0, 0,
	}
}

func (exec *executor) StartExecutor(numWorkers int) {
	for i := 0; i < numWorkers; i++ {
		worker := NewWorker(exec.workerCounter, exec.globalJobQueue, exec.globalResponseQueue)
		exec.wg.Add(1)
		exec.workerCounter++
		go worker.startWorker(&exec.wg)
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
	exec.waitAllWorkers(exec.globalJobQueue)
	return
}

func (exec *executor) waitAllWorkers(globalJobQueue chan *Job) {
	close(globalJobQueue)
	exec.wg.Wait()
}
