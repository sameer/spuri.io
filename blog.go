package main

import (
	"net/http"
	"time"
	"sort"
)

const (
	blogCacheTime = time.Duration(5 * time.Minute)
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
	Index      []BlogPage
	Page       *BlogPage
	pages      []BlogPage
	NextUpdate time.Time
}

func blogHandler() http.HandlerFunc {
	blogContext := &BlogContext{}
	blogServeIndex := func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "blog_index", blogContext)
	}
	blogServePage := func(w http.ResponseWriter, r *http.Request) {

	}
	return func(w http.ResponseWriter, r *http.Request) {
		globalSetHeaders(w, r)
		blogContext.Refresh()
		if len(r.URL.Path) == 0 { // Request for index
			blogServeIndex(w, r)
		} else { // Req for page, need to do handling of this
			blogServePage(w, r)
		}
	}
}

func (this *BlogContext) Refresh() {
	if time.Now().After(this.NextUpdate) {
		// TODO: refresh articles from dir
		sort.Slice(this.pages, func(i, j int) bool { return this.pages[i].Modified.Before(this.pages[i].Modified) })
		*this = BlogContext{
			GlobalContext: globalContext,
			Index:         this.pages,
			Page:          nil,
			pages:         this.pages,
			NextUpdate:    time.Now().Add(blogCacheTime),
		}
	}
	this.GlobalContext.Refresh()
}
