package asyncexecutor

import (
	"fmt"
	"sync"
)

type Worker struct {
	id                     int
	globalJobQueue         chan *Job
	GlobalResponseQueue    chan *ResponseObject
	jobQueue               chan *Job
	responseQueue          chan *ResponseObject
	responseHandler        ResponseHandler
	responseHandlerEnabled bool
}

func NewWorker(workerID int, globalJobQueue chan *Job, GlobalResponseQueue chan *ResponseObject) *Worker {
	return &Worker{
		workerID,
		globalJobQueue,
		GlobalResponseQueue,
		make(chan *Job, 10),
		make(chan *ResponseObject, 10),
		&CleaningHandler{},
		true}
}

func (worker *Worker) AddResponseHandler(handler ResponseHandler) {
	worker.responseHandler = handler
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
			response := job.call()

			// Check if the global response queue is full and remove elements to prevent deadlock.
			// When adding jobs if there is no consumer to remove the result the CreateJob function will block and produce a deadlock.
			if worker.responseHandlerEnabled == true {
				if cap(worker.GlobalResponseQueue)-len(worker.GlobalResponseQueue) == 1 {
					worker.responseHandler.Handle(<-worker.GlobalResponseQueue)
				}
			}
			worker.GlobalResponseQueue <- response

		case job, ok := <-worker.jobQueue:
			if !ok {
				return
			}

			response := job.call()

			if worker.responseHandlerEnabled == true {
				if cap(worker.responseQueue)-len(worker.responseQueue) == 1 {
					worker.responseHandler.Handle(<-worker.responseQueue)
				}
			}
			worker.responseQueue <- response
		}
	}
}

func (worker *Worker) stop() {
	fmt.Printf("Stopping worker [id = %v] \n", worker.id)
	close(worker.jobQueue)
}
