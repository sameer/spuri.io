package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"encoding/base64"
	"crypto/sha256"
)

type staticContext struct {
	NavItems    []NavItem
	CssFileHash string
}

type NavItem struct {
	Name    string
	Link    string
	NewPage bool
}

var staticCtx staticContext

func init() {
	sha_512 := sha256.New()
	if cssFileContent, err := ioutil.ReadFile(cssFilePath); err == nil {
		sha_512.Write(cssFileContent)
	} else { // If we failed here, in all likelihood, the CSS handler will fail too...
		sha_512.Write([]byte(cssDefault))
	}
	staticCtx = staticContext{
		NavItems: []NavItem{
			{"c0dart", "/c0dart/", false},
			{"Blog", "/blog/", false},
			{"Github", "https://github.com/sameer", false},
			{"About", "/about", false},
		},
		CssFileHash: fmt.Sprintf("sha256-%s", base64.StdEncoding.EncodeToString(sha_512.Sum(nil))),
	}
}

func globalSetHeaders(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Via", "0.2 spuri.io")
	w.Header().Add("Upgrade-Insecure-Requests", "1")
	w.Header().Add("X-Powered-By", runtime.Version())
	w.Header().Add("X-XSS-Protection", "1; mode=block")
	w.Header().Add("Content-Language", "en-US")

	if dnt := r.Header.Get("DNT"); dnt != "" { // Do Not Track header, no cookies here except for Cloudflare.
		w.Header().Add("Tk", "N")
	}
}
