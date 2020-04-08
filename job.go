package asyncexecutor

import (
	"errors"
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

func NewJob(id int, function callableType, paramObject *parameterObject) (*Job, error) {
	if function == nil {
		return nil, errors.New("Callable cannot be nil")
	}

	funcV := reflect.ValueOf(function)
	if funcV.Kind() != reflect.Func {
		return nil, errors.New("Callable argument must be a function")
	}

	return &Job{
		id,
		paramObject,
		new(ResponseObject),
		make(chan *ResponseObject, 1),
		function,
	}, nil
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
