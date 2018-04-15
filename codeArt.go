package main

import (
	"bytes"
	"fmt"
	"github.com/disintegration/imaging"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
)

type codeArtContext struct {
	staticContext
	Images     map[string]*codeArtImage
	ImageSlice []codeArtImage
	logger     *log.Logger
}

type codeArtImage struct {
	Filename string
	Href     string
	Src      string
	Title    string
	Desc     string
	Data     []byte
}

const (
	resizerPath      = "resizer/"
	resizeFactor     = 2
	galleryWidth     = 1920 / resizeFactor // Reduction of transfer bandwidth
	galleryHeight    = 1080 / resizeFactor
	codeArtCacheTime = time.Duration(1 * time.Hour)
)

var codeArtHandler = handlerWithUpdatableState{
	updatePeriod: codeArtCacheTime,
	handlerWithFinalState: handlerWithFinalState{
		handlerGenericAttributes: handlerGenericAttributes{codeArtHandlerPath, true},
		handler: func(w http.ResponseWriter, r *http.Request, s state) {
			if ctx := s.(codeArtContext); r.URL.Path == "" { // Gallery
				renderTemplate(w, "codeArt_gallery", ctx)
			} else if strings.Count(r.URL.Path, "/") == 3 && strings.HasPrefix(r.URL.Path, resizerPath) { // Resizer
				ctx.serveResized(w, r)
			} else {
				http.NotFound(w, r)
			}
		},
		initializer: func() state {
			return codeArtContext{logger: log.New(os.Stdout, "codeArt ", log.LstdFlags), Images: make(map[string]*codeArtImage)}
		},
	},
	updater: func(s state) state {
		ctx := codeArtContext{logger: s.(codeArtContext).logger, Images: make(map[string]*codeArtImage)}
		if images, err := ioutil.ReadDir(codeArtDir); err == nil {
			if ctx.ImageSlice == nil {
				ctx.ImageSlice = make([]codeArtImage, 0, len(images))
			}
			resizeWaiter := sync.WaitGroup{}
			for i := range rand.New(rand.NewSource(time.Now().UnixNano())).Perm(len(images)) {
				if _, ok := ctx.Images[images[i].Name()]; !ok {
					imageFileName := images[i].Name()
					ctx.ImageSlice = append(ctx.ImageSlice, codeArtImage{
						Filename: imageFileName,
						Href:     staticHandlerPath + "codeArt/" + imageFileName,
						Src:      fmt.Sprintf(codeArtHandlerPath+resizerPath+"\""+imageFileName+"\"/%d/%d", galleryWidth, galleryHeight),
						Title:    fileNameToTitle(imageFileName),
						Desc:     "", // TODO: make these fields real
						Data:     nil,
					})
					resizeWaiter.Add(1)
					go func(imgInfo *codeArtImage) {
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
			ctx.logger.Println("Error reading codeArt directory", err)
		}
		ctx.staticContext = staticCtx
		ctx.logger.Println("refreshed")
		return ctx
	},
}

func (ctx *codeArtContext) serveResized(w http.ResponseWriter, r *http.Request) {
	var fileName string
	var width, height int

	if _, err := fmt.Sscanf(r.URL.Path, "resizer/%q/%d/%d", &fileName, &width, &height); err != nil {
		http.NotFound(w, r)
	} else if width != galleryWidth || height != galleryHeight {
		http.NotFound(w, r)
	} else if value, ok := ctx.Images[fileName]; !ok {
		ctx.logger.Println("didn't find", value)
		http.NotFound(w, r)
	} else if value.Data == nil {
		ctx.logger.Println("request for", value, "but was not resized")
		http.NotFound(w, r)
	} else {
		w.Header().Add("Content-Length", strconv.Itoa(len(value.Data)))
		w.Header().Add("Cache-Control", fmt.Sprintf("private, max-age=%d", int(codeArtCacheTime.Seconds())))
		w.Write(value.Data)
	}
}

func (ctx *codeArtContext) doResize(imgInfo *codeArtImage) error {
	imgFileName := imgInfo.Filename
	defer func() { ctx.Images[imgFileName] = imgInfo }()
	imgFilePath := codeArtDir + imgFileName
	if _, err := os.Stat(imgFilePath); os.IsNotExist(err) {
		ctx.logger.Println("couldn't find file", imgFileName, ":", err)
		return err
	} else if os.IsPermission(err) {
		ctx.logger.Println("don't have permissions for file", imgFileName, ":", err)
	} else if _, ok := ctx.Images[imgFileName]; ok {
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
