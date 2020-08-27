package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang/groupcache"
	"golang.org/x/crypto/bcrypt"
)

type cache struct {
	group *groupcache.Group
}

func (ac cache) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	key := r.FormValue("key")
	if key == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	now := time.Now()
	defer func() {
		log.Printf("cacheHandler: group[%s]\tkey[%q]\ttime[%v]", ac.group.Name(), key, time.Since(now))
	}()
	var respBody []byte
	if err := ac.group.Get(r.Context(), key, groupcache.AllocatingByteSliceSink(&respBody)); err != nil {
		log.Printf("Get/cache.Get error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(respBody))
}

// Get is am arbitrary getter function. Bcrypt is nice here because, it:
//	1) takes a long time
//	2) uses a random seed so non-cache results for the same key are obvious
func (ac cache) Get(ctx context.Context, key string, dst groupcache.Sink) error {
	now := time.Now()
	defer func() {
		log.Printf("distributedStringCacheKey/key:%q\ttime:%v", key, time.Since(now))
	}()
	out, err := bcrypt.GenerateFromPassword([]byte(key), 14)
	if err != nil {
		return err
	}
	if err := dst.SetBytes(out); err != nil {
		return err
	}
	return nil
}
