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
	"sync"
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
	*globalContext
	Index      []*BlogPage
	Page       *BlogPage
	pages      map[uint32]*BlogPage
	checksums  []uint32
	NextUpdate time.Time
	UpdateMutex sync.Mutex
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

var blogHandler = func() http.HandlerFunc {
	ctx := &blogContext{}
	return func(w http.ResponseWriter, r *http.Request) {
		globalSetHeaders(w, r)
		ctx.refresh()
		if len(r.URL.Path) == 0 { // Request for index
			ctx.serveIndex(w, r)
		} else { // Req for page, need to do handling of this
			ctx.servePage(w, r)
		}
	}
}()

func (ctx *blogContext) refresh() {
	ctx.UpdateMutex.Lock()
	defer ctx.UpdateMutex.Unlock()
	ctx.globalContext.refresh()
	if time.Now().After(ctx.NextUpdate) {
		pages := make(map[uint32]*BlogPage)
		filepath.Walk(blogDir, func(path string, info os.FileInfo, err error) error {
			if err == nil && !info.IsDir() && strings.HasSuffix(info.Name(), "md") {
				bytes, err := ioutil.ReadFile(path)
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

		ctx.globalContext = globalCtx
		ctx.pages = pages
		ctx.checksums = checksums
		ctx.Index = Index
		ctx.NextUpdate = time.Now().Add(blogCacheTime)
	}
}
