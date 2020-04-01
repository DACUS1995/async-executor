package asyncexecutor

import (
	"fmt"
)

type ResponseHandler interface {
	Handle(*ResponseObject)
}

type CleaningHandler struct {
}

func (*CleaningHandler) Handle(response *ResponseObject) {
	fmt.Printf("Finished job: %v", response.ID)
}
