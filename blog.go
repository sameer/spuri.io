package main

import (
	"fmt"
	"github.com/sameer/unsafe-markdown"
	"hash/crc32"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	blogCacheTime = time.Duration(5 * time.Minute)
	blogIndexMax  = 10
	author        = "Sameer Puri"
)

type BlogPage struct {
	Title    string
	Checksum uint32
	Modified time.Time
	Author   string
	Content  template.HTML
}

type BlogContext struct {
	*GlobalContext
	Index      []*BlogPage
	Page       *BlogPage
	pages      map[uint32]*BlogPage
	checksums  []uint32
	NextUpdate time.Time
}

func blogHandler() http.HandlerFunc {
	blogContext := &BlogContext{}
	blogServeIndex := func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "blog_index", blogContext)
	}
	blogServePage := func(w http.ResponseWriter, r *http.Request) {
		var crcUint uint32
		_, err := fmt.Sscanf(r.URL.Path, "%d", &crcUint)
		if err != nil {
			http.NotFound(w, r)
			return
		} else if page, ok := blogContext.pages[uint32(crcUint)]; !ok {
			http.NotFound(w, r)
			return
		} else {
			blogContext.Page = page
			renderTemplate(w, "blog_page", blogContext)
		}
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
		pages := make(map[uint32]*BlogPage)
		filepath.Walk(blogDir, func(path string, info os.FileInfo, err error) error {
			if err == nil && !info.IsDir() {
				bytes, err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}
				checksum := crc32.ChecksumIEEE([]byte(info.Name()))
				pages[checksum] = &BlogPage{
					Title:    strings.TrimSuffix(info.Name(), filepath.Ext(info.Name())),
					Checksum: checksum,
					Modified: info.ModTime(),
					Author:   author,
					Content:  template.HTML(markdown.MarkdownToHtmlString(string(bytes))),
				}
			}
			return err
		})
		checksums := make([]uint32, 0, len(pages))
		for k := range pages {
			checksums = append(checksums, k)
		}

		sort.Slice(checksums, func(i, j int) bool { return pages[checksums[i]].Modified.After(pages[checksums[j]].Modified) })
		Index := make([]*BlogPage, 0, blogIndexMax)
		for _, checksum := range checksums {
			Index = append(Index, pages[checksum])
			if len(Index) == blogIndexMax {
				break
			}
			if len(Index) == blogIndexMax {
				break
			}
		}

		this.GlobalContext = globalContext
		this.pages = pages
		this.checksums = checksums
		this.Index = Index
		this.NextUpdate = time.Now().Add(blogCacheTime)
	}
	this.GlobalContext.Refresh()
}
