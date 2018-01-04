package main

import (
	"net/http"
	"time"
	"sort"
)

type BlogPage struct {
	Title    string
	Posted   time.Time
	Modified time.Time
	Author   string
	Content  string
}

type BlogContext struct {
	*GlobalContext
	Index []BlogPage
	Page  *BlogPage
}

var pages []BlogPage = nil

func blogHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path[len(blogHandlerPath):]
	if len(url) == 0 { // Request for index
		blogServeIndex(w, r)
	} else { // Something else
		blogServePage(w, r, url)
	}
}

func blogServeIndex(w http.ResponseWriter, r *http.Request) {
	sort.Slice(pages, func(i, j int) bool { return pages[i].Modified.Before(pages[i].Modified) })
	ctx := BlogContext{GlobalContext: globalContext, Index: pages, Page: nil}
	renderTemplate(w, "blog_index", ctx)
}

func blogServePage(w http.ResponseWriter, r *http.Request, path string) {

}
