package asyncexecutor

import (
	"reflect"
)

type parameterObject struct {
	parameters []interface{}
}

type ResponseObject struct {
	ID        int
	Responses []interface{}
}

type CallableJob interface {
	call() *ResponseObject
}

type callableType interface{}

type Job struct {
	id              int
	paramObj        *parameterObject
	response        *ResponseObject
	responseChannel chan *ResponseObject
	callable        callableType
}

func NewJob(id int, function callableType, paramObject *parameterObject) *Job {
	return &Job{
		id,
		paramObject,
		new(ResponseObject),
		make(chan *ResponseObject, 1),
		function,
	}
}

func (job *Job) Await() *ResponseObject {
	return <-job.responseChannel
}

func (job *Job) call() *ResponseObject {
	reflectedParams := reflect.TypeOf(job.callable)
	reflectedFunc := reflect.ValueOf(job.callable)

	params := make([]reflect.Value, reflectedParams.NumIn())
	for i := 0; i < reflectedParams.NumIn(); i++ {
		params[i] = reflect.ValueOf(job.paramObj.parameters[i])
	}

	returnValues := reflectedFunc.Call(params)
	Responses := []interface{}{}

	for _, value := range returnValues {
		Responses = append(Responses, value.Interface())
	}

	job.response = &ResponseObject{
		job.id,
		Responses,
	}

	job.responseChannel <- job.response

	return job.response
}
