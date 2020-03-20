package asyncexecutor

import (
	"testing"
)

func TestSimpleExecutor(t *testing.T) {
	expected := "Done"
	queueSize := 10
	numWorkers := 1

	executor := NewExecutor(queueSize)
	executor.StartExecutor(numWorkers)
	executor.CreateJob(
		func(str string) (string, error) {
			return str, nil
		},
		[]interface{}{"Done"},
	)

	if resObj := <-executor.globalResponseQueue; resObj.responses[0] != expected {
		t.Errorf("Expected: [%v] | Returned: [%v]", expected, resObj.responses[0])
	}
}
