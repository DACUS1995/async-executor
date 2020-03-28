package asyncexecutor

import (
	"testing"
)

func TestSimpleJob(t *testing.T) {
	expected := "Done"
	newJob := NewJob(
		/*id*/ 0,
		func(str string) (string, error) {
			return str, nil
		},
		/*paramObject*/ &parameterObject{[]interface{}{expected}},
	)

	newJob.call()
	response := newJob.Await()

	if response.Responses[0] != expected {
		t.Errorf("Expected: [%v] | Returned: [%v]", expected, response.Responses[0])
	}
}
