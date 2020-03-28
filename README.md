# async-executor

[![Build Status](https://travis-ci.com/DACUS1995/async-executor.svg?token=EXzDMzSxfmwgYPg9Ttfx&branch=master)](https://travis-ci.com/DACUS1995/async-executor)

Small package that allows to easily schedule single asynchronous jobs or create asynchronous tasks that sequentially execute the grouped jobs.

---

By simply adding a new job to the executor, that job is inserted in a queue that is accessible by every created worker.

Every worker has access to two job queues, one shared and one private. The private queue is used by jobs tasks to ensure the sequential execution of the provided functions.

When waiting for the result there are two options: call `Await()` on that job or search for job's id in the executor GlobalResponseQueue.

## Basic instructions

* Creating an executor and adding a job:
```go
queueSize := 10
numWorkers := 5

executor := NewExecutor(queueSize)
executor.StartExecutor(numWorkers)

job := executor.CreateJob(
	func(str string) (string, error) {
		return str, nil
	},
	[]interface{}{"Done"}, // the function parameters
)

returnValue := job.Await()
fmt.Println(returnValue.Responses[0])

executor.StopExecutor()
```
* Creating an executor and adding a task:

```go
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
	lastJob.Await()
	executor.StopExecutor()
```