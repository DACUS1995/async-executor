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

	if resObj := <-executor.GlobalResponseQueue; resObj.responses[0] != expected {
		t.Errorf("Expected: [%v] | Returned: [%v]", expected, resObj.responses[0])
	}

	executor.StopExecutor()
}

func TestMultipleJobsExecutor(t *testing.T) {
	expected := "Done"
	queueSize := 10
	numWorkers := 5
	numJobs := 10

	executor := NewExecutor(queueSize)
	executor.StartExecutor(numWorkers)

	testFunction := func(str string) (string, error) {
		return str, nil
	}

	for i := 0; i < numJobs; i++ {
		executor.CreateJob(
			testFunction,
			[]interface{}{"Done"},
		)
	}

	for i := 0; i < numJobs; i++ {
		if resObj := <-executor.GlobalResponseQueue; resObj.responses[0] != expected {
			t.Errorf("Expected: [%v] | Returned: [%v]", expected, resObj.responses[0])
		}
	}

	executor.StopExecutor()
}

func Benchmark(b *testing.B) {
	expected := "Done"
	queueSize := 10
	numWorkers := 5

	executor := NewExecutor(queueSize)
	executor.StartExecutor(numWorkers)

	testFunction := func(str string) (string, error) {
		return str, nil
	}

	for i := 0; i < b.N; i++ {
		executor.CreateJob(
			testFunction,
			[]interface{}{"Done"},
		)
	}

	for i := 0; i < b.N; i++ {
		if resObj := <-executor.GlobalResponseQueue; resObj.responses[0] != expected {
		}
	}
}
