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
	blogAuthor    = "Sameer Puri"
)

type BlogPage struct {
	Title    string
	Checksum uint32
	Modified time.Time
	Author   string
	Content  template.HTML
}

type blogContext struct {
	staticContext
	Index     [blogIndexMax]*BlogPage
	Page      *BlogPage
	pages     map[uint32]*BlogPage
	checksums []uint32
}

func (ctx *blogContext) serveIndex(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "blog_index", ctx)
}

func (ctx *blogContext) servePage(w http.ResponseWriter, r *http.Request) {
	var crcUint uint32
	if _, err := fmt.Sscanf(r.URL.Path, "%d", &crcUint); err != nil {
		http.NotFound(w, r)
	} else if page, ok := ctx.pages[uint32(crcUint)]; !ok {
		http.NotFound(w, r)
	} else {
		ctx.Page = page
		renderTemplate(w, "blog_page", ctx)
		ctx.Page = nil
	}
}

var blogHandler = handlerWithUpdatableState{
	handlerWithFinalState: handlerWithFinalState{
		handlerGenericAttributes: handlerGenericAttributes{
			pathStr:     blogHandlerPath,
			stripPrefix: true,
		},
		handler: func(w http.ResponseWriter, r *http.Request, s state) {
			ctx := s.(blogContext)
			if len(r.URL.Path) == 0 { // Request for index
				ctx.serveIndex(w, r)
			} else { // Req for page, need to do handling of this
				ctx.servePage(w, r)
			}
		},
		initializer: func() state {
			return blogContext{}
		},
	},
	updater: func(_ state) state {
		ctx := blogContext{}
		pages := make(map[uint32]*BlogPage)
		filepath.Walk(blogDir, func(path string, info os.FileInfo, err error) error {
			if err == nil && !info.IsDir() && strings.HasSuffix(info.Name(), "md") {
				bytes, err := ioutil.ReadFile(path)
				fmt.Println("adding", path)
				if err != nil {
					return err
				}
				checksum := crc32.ChecksumIEEE([]byte(info.Name()))
				pages[checksum] = &BlogPage{
					Title:    strings.TrimSuffix(info.Name(), filepath.Ext(info.Name())),
					Checksum: checksum,
					Modified: info.ModTime(),
					Author:   blogAuthor,
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
		var Index [blogIndexMax]*BlogPage
		for i, checksum := range checksums {
			Index[i] = pages[checksum]
			if i+1 == blogIndexMax {
				break
			}
		}

		ctx.staticContext = staticCtx
		ctx.pages = pages
		ctx.checksums = checksums
		ctx.Index = Index
		return ctx
	},
	updatePeriod: blogCacheTime,
}
