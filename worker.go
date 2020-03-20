package asyncexecutor

import "sync"

type Worker struct {
	id                  int
	globalJobQueue      chan *Job
	globalResponseQueue chan *ResponseObject
	jobQueue            chan *Job
	responseQueue       chan *ResponseObject
}

func NewWorker(id int, globalJobQueue chan *Job, globalResponseQueue chan *ResponseObject) *Worker {
	return &Worker{
		id,
		globalJobQueue,
		globalResponseQueue,
		make(chan *Job, 10),
		make(chan *ResponseObject, 10)}
}

func (worker *Worker) startWorker(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case job := <-worker.globalJobQueue:
			worker.globalResponseQueue <- job.call()
		case job := <-worker.jobQueue:
			worker.responseQueue <- job.call()
		}
	}
}
