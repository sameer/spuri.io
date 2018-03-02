package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"crypto/subtle"
	"os"
	"sync/atomic"
)

var studioStatisticsHandler http.HandlerFunc = func() func(w http.ResponseWriter, r *http.Request) {
	const apiKeyHeader = "x-api-key"
	var atomicStatsImage atomic.Value
	atomicStatsImage.Store([]byte{})
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			if apiKey := r.Header.Get(apiKeyHeader); apiKey != "" && r.Header.Get("content-type") == "image/png" && subtle.ConstantTimeCompare([]byte(apiKey), []byte(os.Getenv("x_api_key"))) == 1 {
				img, err := ioutil.ReadAll(r.Body)
				if err != nil {
					atomicStatsImage.Store(img)
				}
			}
		} else if r.Method == http.MethodGet {
			img := atomicStatsImage.Load().([]byte)
			w.Header().Add("Content-Length", strconv.Itoa(len(img)))
			w.Header().Add("Cache-Control", fmt.Sprintf("private, max-age=%d", 60))
			w.Write(img)
		}
	}
}()
