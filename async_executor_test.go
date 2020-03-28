package asyncexecutor

import (
	"testing"
)

func TestSimpleJobExecutor(t *testing.T) {
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

	if resObj := <-executor.GlobalResponseQueue; resObj.Responses[0] != expected {
		t.Errorf("Expected: [%v] | Returned: [%v]", expected, resObj.Responses[0])
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
		if resObj := <-executor.GlobalResponseQueue; resObj.Responses[0] != expected {
			t.Errorf("Expected: [%v] | Returned: [%v]", expected, resObj.Responses[0])
		}
	}

	executor.StopExecutor()
}

func TestCompletionMultipleJobsExecutor(t *testing.T) {
	expected := "Done"
	queueSize := 10
	numWorkers := 5
	numJobs := 10
	jobs := make(map[int]bool)

	executor := NewExecutor(queueSize)
	executor.StartExecutor(numWorkers)

	testFunction := func(str string) (string, error) {
		return str, nil
	}

	for i := 0; i < numJobs; i++ {
		newJob := executor.CreateJob(
			testFunction,
			[]interface{}{"Done"},
		)
		jobs[newJob.id] = false
	}

	for i := 0; i < numJobs; i++ {
		if resObj := <-executor.GlobalResponseQueue; resObj.Responses[0] != expected {
			jobs[resObj.id] = true
		}
	}

	unfinishedJobsCounter := 0
	for _, v := range jobs {
		if v != false {
			unfinishedJobsCounter++
		}
	}

	if unfinishedJobsCounter != 0 {
		t.Errorf("[%v] unfinished jobs are stuck in the executor.", unfinishedJobsCounter)
	}

	executor.StopExecutor()
}

func TestSimpleTaskExecutor(t *testing.T) {
	expected := "Done"
	queueSize := 10
	numWorkers := 1
	testFunction := func(str string) (string, error) {
		return str, nil
	}

	executor := NewExecutor(queueSize)
	executor.StartExecutor(numWorkers)

	taskList := []*Job{}

	for i := 0; i < queueSize; i++ {
		taskList = append(taskList, executor.CreateTaskJob(
			testFunction,
			[]interface{}{expected},
		))
	}

	lastJob := executor.CreateTask(taskList)

	if resObj := lastJob.Await(); resObj.Responses[0] != expected {
		t.Errorf("Expected: [%v] | Returned: [%v]", expected, resObj.Responses[0])
	}

	executor.StopExecutor()
}

func Benchmark(b *testing.B) {
	expected := "Done"
	queueSize := 10000
	numWorkers := 5

	executor := NewExecutor(queueSize)
	executor.StartExecutor(numWorkers)

	testFunction := func(str string) (string, error) {
		return str, nil
	}

	for i := 0; i < b.N; i++ {
		executor.CreateJob(
			testFunction,
			[]interface{}{expected},
		)
	}

	for i := 0; i < b.N; i++ {
		if resObj := <-executor.GlobalResponseQueue; resObj.Responses[0] != expected {
		}
	}

	executor.StopExecutor()
}
