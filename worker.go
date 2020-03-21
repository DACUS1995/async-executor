package asyncexecutor

import (
	"fmt"
	"sync"
)

type Worker struct {
	id                  int
	globalJobQueue      chan *Job
	GlobalResponseQueue chan *ResponseObject
	jobQueue            chan *Job
	responseQueue       chan *ResponseObject
}

func NewWorker(workerID int, globalJobQueue chan *Job, GlobalResponseQueue chan *ResponseObject) *Worker {
	return &Worker{
		workerID,
		globalJobQueue,
		GlobalResponseQueue,
		make(chan *Job, 10),
		make(chan *ResponseObject, 10)}
}

func (worker *Worker) start(wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("Started worker [id = %v] \n", worker.id)

	for {
		select {
		case job, ok := <-worker.globalJobQueue:
			if !ok {
				return
			}
			worker.GlobalResponseQueue <- job.call()
		case job, ok := <-worker.jobQueue:
			if !ok {
				return
			}
			worker.responseQueue <- job.call()
		}
	}
}

func (worker *Worker) stop() {
	fmt.Printf("Stopping worker [id = %v] \n", worker.id)
	close(worker.jobQueue)
}
