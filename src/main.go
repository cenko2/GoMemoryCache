package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"./CacheLogic"
	"./ServiceLogic"
	"github.com/gorilla/mux"
)

func main() {

	memoryCache := CacheLogic.NewMemoryCache(time.Duration(10) * time.Second)

	x := ServiceLogic.New(memoryCache, "key", 30)

	r := mux.NewRouter()
	r.HandleFunc("/cache/{key}", x.HandleGet).Methods("GET")
	r.HandleFunc("/cache/{key}", x.HandlePost).Methods("POST")
	r.HandleFunc("/status", handleHeartBeat).Methods("GET")
	http.Handle("/", r)

	srv := &http.Server{
		Handler:      r,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("Instance initialized routing complete")
	log.Fatal(srv.ListenAndServe())
}

func handleHeartBeat(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Heartbeat recevied")
}
