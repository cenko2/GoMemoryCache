package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"./CacheLogic"
	"./ServiceLogic"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
)

const (
	defaultAddr = ":http"
)

func main() {

	//  to run locally
	//  docker pull redis
	//  docker run -d -p 6379:6379 redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     getRedisHost(),
		Password: "",
		DB:       0,
	})

	redisClient := CacheLogic.RedisCache{Rdb: rdb}

	x := ServiceLogic.ServiceHandler{ParamName: "key", CacheHandler: redisClient, DefaultTTLInMinutes: 30}

	r := mux.NewRouter()
	r.HandleFunc("/cache/{key}", x.HandleGet).Methods("GET")
	r.HandleFunc("/cache/{key}", x.HandlePost).Methods("POST")
	r.HandleFunc("/status", handleHearBeat).Methods("GET")
	http.Handle("/", r)

	srv := &http.Server{
		Handler:      r,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("Instance initialized routing complete")
	log.Fatal(srv.ListenAndServe())
}

func handleHearBeat(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Heartbeat recevied")
}

func getRedisHost() string {
	envRedis := os.Getenv("REDISHOST")
	if envRedis != "" {
		return envRedis
	} else {
		return "localhost:6379"
	}
}
