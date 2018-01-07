package main

import (
	"net/http"
	"time"
	"fmt"
	"sort"
	"path/filepath"
	"os"
	"io/ioutil"
)

const (
	blogCacheTime = time.Duration(5 * time.Minute)
	blogIndexMax  = 10
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
	pages      map[time.Time]BlogPage
	dates      []time.Time
	NextUpdate time.Time
}

func blogHandler() http.HandlerFunc {
	blogContext := &BlogContext{

	}
	blogServeIndex := func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "blog_index", blogContext)
	}
	blogServePage := func(w http.ResponseWriter, r *http.Request) {
		var dateStr string
		var page int
		fmt.Sscanf(r.URL.Path, "/%q/%d", &dateStr, &page)

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
		pages := make(map[time.Time]BlogPage)
		filepath.Walk(blogDir, func(path string, info os.FileInfo, err error) error {
			if err == nil && !info.IsDir() {
				bytes, err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}
				// TODO: get all fields from file somehow
				pages[info.ModTime()] = BlogPage{
					Title:    info.Name(),
					Posted:   info.ModTime(),
					Modified: info.ModTime(),
					Author:   "",
					Content:  string(bytes),
				}
			}
			return err
		})
		dates := make([]time.Time, 0, len(pages))
		for k, _ := range pages {
			dates = append(dates, k)
		}

		sort.Slice(dates, func(i, j int) bool { return dates[i].After(dates[i]) })
		Index := make([]BlogPage, 0, blogIndexMax)
		for _, date := range dates {
			Index = append(Index, pages[date])
			if len(Index) == blogIndexMax {
				break
			}
			if len(Index) == blogIndexMax {
				break
			}
		}

		this.dates = dates
		this.Index = Index
		this.NextUpdate = time.Now().Add(blogCacheTime)
	}
	this.GlobalContext.Refresh()
}
