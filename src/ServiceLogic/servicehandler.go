//Package ServiceLogic handles the rest requests and forwards them to the cache handler
package ServiceLogic

import (
	"bytes"
	"io"
	"log"
	"net/http"

	"time"

	"../CacheLogic"
	"github.com/gorilla/mux"
)

// Servicehandler that implements the ServiceHandlerInterface
type ServiceHandler struct {
	*serviceHandler
}

type serviceHandler struct {
	paramName           string
	defaultTTLInMinutes int
	cacheHandler        CacheLogic.CacheHandlerInterface
}

func (s *serviceHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	//Read key value from address params
	vars := mux.Vars(r)
	key := (vars[s.paramName])

	//Timing
	now := time.Now()
	defer func() {
		log.Printf("Get\tkey:%q\ttime:%v", key, time.Since(now))
	}()

	if !(s.cacheHandler).KeyExists(key) {
		http.Error(w, "Key not found", http.StatusNotFound)
	} else {
		value := (s.cacheHandler).Get(key)
		io.WriteString(w, value)
	}
}

func (s *serviceHandler) HandlePost(w http.ResponseWriter, r *http.Request) {

	// Read key
	vars := mux.Vars(r)
	key := (vars[s.paramName])

	//Instrumentation
	now := time.Now()
	defer func() {
		log.Printf("Post\tkey:%q\ttime:%v", key, time.Since(now))
	}()

	// Read the value
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	value := buf.String()
	if len(value) > 100000 {
		http.Error(w, "Input too long", http.StatusBadRequest)
		return
	}
	(s.cacheHandler).Insert(key, &value, s.defaultTTLInMinutes)
}

// New returns a new service hanbdler. excepts a cache handler interface, parametername to exract the key from the url and default expiration time in minutes
func New(cacheHandler CacheLogic.CacheHandlerInterface, paramName string, defaultTTLInMinutes int) *ServiceHandler {
	s := &serviceHandler{
		cacheHandler:        cacheHandler,
		paramName:           paramName,
		defaultTTLInMinutes: defaultTTLInMinutes,
	}

	S := &ServiceHandler{s}

	return S
}
