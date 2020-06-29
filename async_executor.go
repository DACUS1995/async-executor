package asyncexecutor

import (
	"sync"

	"github.com/pkg/errors"
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
	Stop()
	WaitAndStop()
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

func (exec *executor) CreateJob(function callableType, parameters []interface{}) (*Job, error) {
	newJobID := exec.jobCounter
	newJob, err := NewJob(
		newJobID,
		function,
		&parameterObject{parameters},
	)

	if err != nil {
		return nil, errors.Wrap(err, "Job creation failed.")
	}

	exec.addGlobalJob(newJob)
	exec.jobCounter++

	return newJob, nil
}

func (exec *executor) CreateTaskJob(function callableType, parameters []interface{}) (*Job, error) {
	newJobID := exec.jobCounter
	newJob, err := NewJob(
		newJobID,
		function,
		&parameterObject{parameters},
	)

	if err != nil {
		return nil, errors.Wrap(err, "Job creation failed.")
	}

	exec.jobCounter++
	return newJob, nil
}

func (exec *executor) CreateTask(taskJobList []*Job) *Job {
	sizeOfTask := len(taskJobList)

	// Iterate over the started workers to search for the most free worker
	// or the first one that has enough space for all the jobs from the task.
	bestWorker := exec.workers[0]
	for _, worker := range exec.workers {
		currWorkerAvailablity := cap(worker.responseQueue) - len(worker.responseQueue)
		if currWorkerAvailablity >= sizeOfTask {
			bestWorker = worker
			break
		}

		if currWorkerAvailablity > cap(bestWorker.responseQueue)-len(bestWorker.responseQueue) {
			bestWorker = worker
		}
	}

	for _, job := range taskJobList {
		bestWorker.jobQueue <- job
	}

	return taskJobList[len(taskJobList)-1]
}

func (exec *executor) SetResponseHandler(handler ResponseHandler) {
	for _, worker := range exec.workers {
		worker.SetResponseHandler(handler)
	}
}

func (exec *executor) SetStatusResponseHandler(enabled bool) {
	for _, worker := range exec.workers {
		worker.responseHandlerEnabled = enabled
	}
}

func (exec *executor) addGlobalJob(job *Job) {
	exec.globalJobQueue <- job
}

func (exec *executor) Stop() {
	exec.waitAllWorkers()
	return
}

func (exec *executor) WaitAndStop() int {
	jobCounter := 0

	for {
		select {
		case <-exec.GlobalResponseQueue:
			jobCounter++
			// Maybe save the responses somewhere
		default:
			exec.Stop()
			return jobCounter
		}
	}
}

func (exec *executor) waitAllWorkers() {
	close(exec.globalJobQueue)

	for _, worker := range exec.workers {
		worker.stop()
	}
	exec.wg.Wait()
}
