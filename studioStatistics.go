package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

var studioStatisticsHandler = func() func(w http.ResponseWriter, r *http.Request) {
	var studioStatisticsImage []byte
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			if r.Header.Get("x-api-key") != "" && r.Header.Get("content-type") == "image/png" {
				img, err := ioutil.ReadAll(r.Body)
				if err != nil {
					studioStatisticsImage = img
				}
			}
		} else if r.Method == http.MethodGet {
			globalSetHeaders(w, r)
			w.Header().Add("Content-Length", strconv.Itoa(len(studioStatisticsImage)))
			w.Header().Add("Cache-Control", fmt.Sprintf("private, max-age=%d", 60))
			w.Write(studioStatisticsImage)
		}
	}
}()
