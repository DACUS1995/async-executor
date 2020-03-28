package asyncexecutor

import (
	"fmt"
	"testing"
	"time"
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

	executor.Stop()
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

	executor.Stop()
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

	executor.Stop()
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

	executor.Stop()
}

func TestWaitAndStopExecutor(t *testing.T) {
	expected := "Done"
	queueSize := 10
	numWorkers := 1
	testFunction := func(str string) (string, error) {
		return str, nil
	}

	executor := NewExecutor(queueSize)
	executor.StartExecutor(numWorkers)

	for i := 0; i < queueSize; i++ {
		executor.CreateJob(
			testFunction,
			[]interface{}{expected},
		)
	}

	time.Sleep(60 * time.Millisecond)
	numOfJobsWaited := executor.WaitAndStop()

	if numOfJobsWaited != queueSize-1 {
		t.Errorf("Expected: [%v] | Returned: [%v]", queueSize, numOfJobsWaited)
	}
}

func Benchmark(b *testing.B) {
	fmt.Printf("Benchmark with [%v] elements.", b.N)
	expected := "Done"
	queueSize := 1000
	numWorkers := 5

	executor := NewExecutor(queueSize)
	executor.StartExecutor(numWorkers)

	testFunction := func(str string) (string, error) {
		return str, nil
	}

	var lastJob *Job
	for i := 0; i < b.N; i++ {
		lastJob = executor.CreateJob(
			testFunction,
			[]interface{}{expected},
		)
	}

	lastJob.Await()

	executor.Stop()
}
