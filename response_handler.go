package asyncexecutor

type ResponseHandler interface {
	Handle(*ResponseObject)
}

type CleaningHandler struct {
}

func (*CleaningHandler) Handle(response *ResponseObject) {

}
