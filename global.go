package main

import (
	"net/http"
	"runtime"
)

type GlobalContext struct {
	NavItems NavItems
}

type NavItem struct {
	Name    string
	Link    string
	NewPage bool
}

type NavItems []NavItem

var globalContext *GlobalContext = nil

func initGlobalContext() {
	globalContext = &GlobalContext{NavItems: NavItems{
		NavItem{"c0dart", "/c0dart/", false,},
		NavItem{"Blog", "/blog/", false,},
		NavItem{"Github", "https://github.com/sameer", false,},
		NavItem{"About", "/about", false},
	}}
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
