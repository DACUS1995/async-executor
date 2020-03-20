package asyncexecutor

import (
	"reflect"
)

type parameterObject struct {
	parameters []interface{}
}

type ResponseObject struct {
	id        int
	responses []interface{}
}

type CallableJob interface {
	call() *ResponseObject
}

type callableType interface{}

type Job struct {
	id       int
	paramObj *parameterObject
	response *ResponseObject
	callable callableType
}

func NewJob(id int, function callableType, paramObject *parameterObject) *Job {
	return &Job{
		id,
		paramObject,
		new(ResponseObject),
		function,
	}
}

func (job *Job) call() *ResponseObject {
	reflectedParams := reflect.TypeOf(job.callable)
	reflectedFunc := reflect.ValueOf(job.callable)

	params := make([]reflect.Value, reflectedParams.NumIn())
	for i := 0; i < reflectedParams.NumIn(); i++ {
		params[i] = reflect.ValueOf(job.paramObj.parameters[i])
	}

	returnValues := reflectedFunc.Call(params)
	responses := []interface{}{}

	for _, value := range returnValues {
		responses = append(responses, value.Interface())
	}

	return &ResponseObject{
		job.id,
		responses,
	}
}
