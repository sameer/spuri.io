package main

import (
	"net/http"
	"runtime"
	"time"
	"sync/atomic"
)

const (
	globalContextCacheTime = time.Duration(time.Hour * 24 * 360) // Never refresh, it's essentially static
)

type staticContext struct {
	NavItems   []NavItem
	NextUpdate time.Time
}

type NavItem struct {
	Name    string
	Link    string
	NewPage bool
}

var staticCtx atomic.Value

func (this *staticContext) init() {
	*this = staticContext{
		NavItems: []NavItem{
			{"c0dart", "/c0dart/", false},
			{"Blog", "/blog/", false},
			{"Github", "https://github.com/sameer", false},
			{"About", "/about", false},
		},
		NextUpdate: time.Now().Add(globalContextCacheTime),
	}
}

func globalSetHeaders(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Via", "0.1 spuri.io")
	w.Header().Add("Upgrade-Insecure-Requests", "1")
	w.Header().Add("X-Powered-By", runtime.Version())
	w.Header().Add("X-XSS-Protection", "1; mode=block")
	w.Header().Add("Content-Language", "en-US")

	if dnt := r.Header.Get("DNT"); dnt != "" { // Do Not Track header, no cookies here except for Cloudflare.
		w.Header().Add("Tk", "N")
	}
}
