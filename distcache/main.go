package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/golang/groupcache"
	"github.com/pomerium/autocache"
)

const (
	defaultAddr = ":http"
)

func main() {
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = defaultAddr
	}

	existing := []string{}
	if nodes := os.Getenv("NODES"); nodes != "" {
		existing = strings.Split(nodes, ",")
	}

	ac, err := autocache.New(&autocache.Options{})
	if err != nil {
		log.Fatal(err)
	}
	if _, err := ac.Join(existing); err != nil {
		log.Fatal(err)
	}

	var distributedCache cache
	distributedCache.group = groupcache.NewGroup("distributedStringCache", 1<<20, distributedCache)

	mux := http.NewServeMux()
	mux.Handle("/get/", distributedCache)

	mux.Handle("/_groupcache/", ac)
	log.Fatal(http.ListenAndServe(addr, mux))

}
