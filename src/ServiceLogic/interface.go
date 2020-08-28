package ServiceLogic

import "net/http"

type ServiceHandlerInterface interface {
	HandleGet(http.ResponseWriter, *http.Request)
	HandlePost(http.ResponseWriter, *http.Request)
}
