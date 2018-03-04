package main

import (
	"bytes"
	"fmt"
	"github.com/disintegration/imaging"
	"image"
	"image/png"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
	"unicode"
	"strconv"
	"log"
)

type c0dartContext struct {
	*staticContext
	Images      sync.Map
	ImageSlice  []c0dartImage
	NextUpdate  time.Time
	UpdateMutex sync.RWMutex
	logger      *log.Logger
}

type c0dartImage struct {
	Filename string
	Href     string
	Src      string
	Title    string
	Desc     string
	Data     []byte
}

const (
	resizerPath     = "resizer/"
	resizeFactor    = 2
	galleryWidth    = 1920 / resizeFactor // Reduction of transfer bandwidth
	galleryHeight   = 1080 / resizeFactor
	c0dartCacheTime = time.Duration(1 * time.Hour)
)

func (ctx *c0dartContext) serveResizer(w http.ResponseWriter, r *http.Request) {
	var fileName string
	var width, height int

	if _, err := fmt.Sscanf(r.URL.Path, "resizer/%q/%d/%d", &fileName, &width, &height); err != nil {
		http.NotFound(w, r)
	} else if width != galleryWidth || height != galleryHeight {
		http.NotFound(w, r)
	} else if value, ok := ctx.Images.Load(fileName); !ok {
		ctx.logger.Println("didn't find", value)
		http.NotFound(w, r)
	} else if value.(*c0dartImage).Data == nil {
		ctx.logger.Println("request for", value, "but was not resized")
		http.NotFound(w, r)
	} else {
		imgInfo := value.(*c0dartImage)
		w.Header().Add("Content-Length", strconv.Itoa(len(imgInfo.Data)))
		w.Header().Add("Cache-Control", fmt.Sprintf("private, max-age=%d", int(c0dartCacheTime.Seconds())))
		w.Write(imgInfo.Data)
	}
}

func (ctx *c0dartContext) doResize(imgInfo *c0dartImage) error {
	imgFileName := imgInfo.Filename
	defer ctx.Images.Store(imgFileName, imgInfo)
	imgFilePath := c0dartDir + imgFileName
	if _, err := os.Stat(imgFilePath); os.IsNotExist(err) {
		ctx.logger.Println("couldn't find file", imgFileName, ":", err)
		return err
	} else if os.IsPermission(err) {
		ctx.logger.Println("don't have permissions for file", imgFileName, ":", err)
	} else if _, ok := ctx.Images.Load(imgFileName); ok {
		return nil
	}
	imgFile, err := os.Open(imgFilePath)
	if err != nil {
		ctx.logger.Println("failed to open image", err)
		return err
	}
	decodedImgFile, _, err := image.Decode(imgFile)
	if err != nil {
		ctx.logger.Println("failed to decode image", err)
		return err
	}
	var bufout bytes.Buffer
	resizedImage := imaging.Resize(decodedImgFile, galleryWidth, galleryHeight, imaging.Lanczos)
	png.Encode(&bufout, resizedImage)
	imgInfo.Data = bufout.Bytes()
	return nil
}

var c0dartHandler = func() http.HandlerFunc {
	c0dartContext := &c0dartContext{logger: log.New(os.Stdout, "C0dart ", log.LstdFlags)}
	return func(w http.ResponseWriter, r *http.Request) {
		c0dartContext.refresh()
		c0dartContext.UpdateMutex.RLock()
		defer c0dartContext.UpdateMutex.RUnlock()
		if r.URL.Path == "" { // Gallery
			renderTemplate(w, "c0dart_gallery", c0dartContext)
		} else if strings.Count(r.URL.Path, "/") == 3 && strings.HasPrefix(r.URL.Path, resizerPath) { // Resizer
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
	if time.Now().After(ctx.NextUpdate) {
		if images, err := ioutil.ReadDir(c0dartDir); err == nil {
			if ctx.ImageSlice == nil {
				ctx.ImageSlice = make([]c0dartImage, 0, len(images))
			}
			resizeWaiter := sync.WaitGroup{}
			for i := range rand.New(rand.NewSource(time.Now().UnixNano())).Perm(len(images)) {
				if _, ok := ctx.Images.Load(images[i].Name()); !ok {
					imageFileName := images[i].Name()
					ctx.ImageSlice = append(ctx.ImageSlice, c0dartImage{
						Filename: imageFileName,
						Href:     staticHandlerPath + "c0dart/" + imageFileName,
						Src:      fmt.Sprintf(c0dartHandlerPath+resizerPath+"\""+imageFileName+"\"/%d/%d", galleryWidth, galleryHeight),
						Title:    fileNameToTitle(imageFileName),
						Desc:     "", // TODO: make these fields real
						Data:     nil,
					})
					resizeWaiter.Add(1)
					go func(imgInfo *c0dartImage) {
						defer resizeWaiter.Done()
						if err := ctx.doResize(imgInfo); err != nil {
							ctx.logger.Println("error in resizing", err)
						}
					}(&ctx.ImageSlice[len(ctx.ImageSlice)-1])
					ctx.logger.Println("resizing", images[i].Name())
				}
			}
			resizeWaiter.Wait()
		} else {
			ctx.logger.Println("Error reading c0dart directory", err)
		}
		ctx.staticContext, ctx.NextUpdate = staticCtx.Load().(*staticContext), time.Now().Add(c0dartCacheTime)
		ctx.logger.Println("refreshed")
	}
}
