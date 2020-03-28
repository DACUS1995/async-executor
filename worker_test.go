package asyncexecutor

import (
	"sync"
	"testing"
)

func TestSimpleWorker(t *testing.T) {
	expected := "Done"
	newJob := NewJob(
		/*id*/ 0,
		func(str string) (string, error) {
			return str, nil
		},
		/*paramObject*/ &parameterObject{[]interface{}{expected}},
	)

	worker := NewWorker(
		0,
		make(chan *Job),
		make(chan *ResponseObject),
	)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go worker.start(wg)

	worker.jobQueue <- newJob
	response := newJob.Await()

	if response.Responses[0] != expected {
		t.Errorf("Expected: [%v] | Returned: [%v]", expected, response.Responses[0])
	}

	worker.stop()

	wg.Wait()
}
