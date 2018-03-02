package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync/atomic"
	"strconv"
	"sync"
)

var cssHandler = func() http.HandlerFunc {
	// Keep the file in-memory because it's only several KB & lowers load time.
	var atomicCssFile atomic.Value
	atomicCssFile.Store([]byte{})
	var loadMutex sync.Mutex
	const defaultCSS = "* { font-family: sans-serif; }"

	return func(w http.ResponseWriter, r *http.Request) {
		loadMutex.Lock()
		cssFile := atomicCssFile.Load().([]byte)
		if len(cssFile) == 0 {
			var err error
			cssFile, err = ioutil.ReadFile(cssFilePath)
			if err != nil {
				cssFile = []byte(defaultCSS)
				fmt.Println(err)
			}
			atomicCssFile.Store(cssFile)
		}
		loadMutex.Unlock()
		w.Header().Set("Content-Type", "text/css")
		w.Header().Set("Content-Length", strconv.Itoa(len(cssFile)))
		w.Write(cssFile)
	}
}()
