package main

import (
	"bytes"
	"fmt"
	"github.com/esimov/caire"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
	"math/rand"
	"runtime"
)

type c0dartContext struct {
	*globalContext
	Images        []c0dartImage
	ResizerImages sync.Map
	NextUpdate    time.Time
	UpdateMutex sync.Mutex
}

type c0dartImage struct {
	Filename string
	Href     string
	Src      string
	Title    string
	Desc     string
}

const (
	resizerPath     = "resizer/"
	resizeFactor = 2
	galleryWidth    = 1920 / resizeFactor // Reduction of transfer bandwidth
	galleryHeight   = 1080 / resizeFactor
	c0dartCacheTime = time.Duration(1 * time.Hour)
)

func (ctx *c0dartContext) serveResizer(w http.ResponseWriter, r *http.Request) {
	var fileName string
	var width, height int

	_, err := fmt.Sscanf(r.URL.Path, "resizer/%q/%d/%d", &fileName, &width, &height)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	for _, c0dartImage := range ctx.Images {
		if c0dartImage.Filename == fileName && width == galleryWidth && height == galleryHeight {
			resizedImageReader, ok := ctx.ResizerImages.Load(fileName)
			if !ok {
				fmt.Println("Found unresized image", c0dartImage.Filename, "in", ctx.ResizerImages)
				break
			} else if resizedImageReader == nil {
				fmt.Println("Found errored resize image", c0dartImage.Filename)
				break
			}
			w.Header().Add("Content-Length", strconv.Itoa(len(resizedImageReader.([]byte))))
			w.Header().Add("Cache-Control", fmt.Sprintf("private, max-age=%d", int(c0dartCacheTime.Seconds())))
			w.Write(resizedImageReader.([]byte))
			return
		}
	}
	http.NotFound(w, r)
}

func (ctx *c0dartContext) doResize(fileName string) error {
	if _, err := os.Stat(c0dartDir + fileName); os.IsNotExist(err) {
		fmt.Println("C0dart resizer didn't find file", fileName, "with error", err)
		return err
	} else if os.IsPermission(err) {
		fmt.Println("C0dart resizer doesn't have permissions for file", fileName, "with error", err)
	} else if _, ok := ctx.ResizerImages.Load(fileName); ok {
		return nil
	}
	file, err := os.Open(c0dartDir + fileName)
	if err != nil {
		fmt.Println("Failed to open image: ", err)
		return err
	}
	var imgBuf bytes.Buffer
	processor := &caire.Processor{
		BlurRadius:     1,
		SobelThreshold: 10,
		NewWidth:       0,
		NewHeight:      0,
		Percentage:     false,
		Debug:          false,
	}
	err = processor.Process(file, &imgBuf)
	if err != nil {
		fmt.Println("Failed to resize image")
		return err
	}
	ctx.ResizerImages.Store(fileName, imgBuf.Bytes())
	return nil
}

var c0dartHandler = func() http.HandlerFunc {
	c0dartContext := &c0dartContext{}
	return func(w http.ResponseWriter, r *http.Request) {
		globalSetHeaders(w, r)
		c0dartContext.refresh()
		if r.URL.Path == "" { // Gallery
			renderTemplate(w, "c0dart_gallery", c0dartContext)
			return
		} else if strings.HasPrefix(r.URL.Path, resizerPath) && strings.Count(r.URL.Path, "/") == 3 { // Resizer
			c0dartContext.serveResizer(w, r)
		} else {
			http.NotFound(w, r)
		}
	}
}()

func fileNameToTitle(fileName string) string {
	var outstring = ""
	for _, r := range fileName {
		if unicode.IsUpper(r) {
			outstring += " " + string(r)
		} else if r == '.' {
			break
		} else {
			outstring += string(r)
		}
	}
	return outstring
}

func (ctx *c0dartContext) refresh() {
	ctx.UpdateMutex.Lock()
	defer ctx.UpdateMutex.Unlock()
	ctx.globalContext.refresh()
	if time.Now().After(ctx.NextUpdate) {
		var c0dartImages []c0dartImage = nil
		if images, err := ioutil.ReadDir(c0dartDir); err == nil {
			c0dartImages = make([]c0dartImage, len(images))
			resizeWaiter := sync.WaitGroup{}

			for i, j := range rand.New(rand.NewSource(time.Now().UnixNano())).Perm(len(images)) {
				resizeWaiter.Add(1)
				imageName := images[i].Name()
				c0dartImages[j] = c0dartImage{
					Filename: imageName,
					Href:     staticHandlerPath + "c0dart/" + imageName ,
					Src:      fmt.Sprintf(c0dartHandlerPath+resizerPath+"\""+imageName+"\"/%d/%d", galleryWidth, galleryHeight),
					Title:    fileNameToTitle(imageName),
					Desc:     "", // TODO: make these fields real
				}
				go func(name string) {
					defer resizeWaiter.Done()
					ctx.doResize(name)
				}(imageName)
			}
			resizeWaiter.Wait()
			runtime.GC() // Force gc to collect image processing garbage
		} else {
			fmt.Printf("Error reading c0dart directory: %v\n", err)
		}
		ctx.globalContext, ctx.Images, ctx.NextUpdate = globalCtx, c0dartImages, time.Now().Add(c0dartCacheTime)
	}
}
