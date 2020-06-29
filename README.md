# async-executor

[![Build Status](https://travis-ci.com/DACUS1995/async-executor.svg?token=EXzDMzSxfmwgYPg9Ttfx&branch=master)](https://travis-ci.com/DACUS1995/async-executor)

Small module that allows to easily schedule single asynchronous jobs or create asynchronous tasks that sequentially execute the grouped jobs.

---

By simply adding a new job to the executor, that job is inserted in a queue that is accessible by every created worker.

Every worker has access to two job queues, one shared and one private. The private queue is used by jobs tasks to ensure the sequential execution of the provided functions.

When waiting for the result there are two options: call `Await()` on that job or keep searching for job's id in the executor GlobalResponseQueue.

#### Installation
After making sure that Go is installed on your device.
You can use the following command in your terminal:

	go get github.com/DACUS1995/async-executor


#### Import package
Add following line in your `*.go` file:
```go
import "github.com/DACUS1005/async-executor"
```
If you are unhappy to use long `asyncexecutor`, you can do something like this:
```go
import (
	async "github.com/DACUS1005/async-executor"
)
```

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

executor.Stop()
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
executor.Stop()
```

* You can add a custom job response handler that can be used to save, log, send, etc. the responses by implementing the ResponseHandler interface:

```go
type CustomHandler struct{}

func (*CustomHandler) Handle(response *ResponseObject) {
	fmt.Printf("Finished job: %v", response.ID)
}

executor.SetResponseHandler(&CustomHandler{})
```
