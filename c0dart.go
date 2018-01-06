package main

import (
	"net/http"
	"io/ioutil"
	"fmt"
	"github.com/tobi/mogrify-go"
	"strings"
	"time"
	"os"
	"bytes"
	"bufio"
	"strconv"
	"unicode"
)

type C0dartContext struct {
	*GlobalContext
	Images     []C0dartImage
	NextUpdate time.Time
}

var c0dartContext *C0dartContext = nil

type C0dartImage struct {
	Filename string
	Href     string
	Title    string
	Desc     string
}

const (
	resizerPath     = "resizer/"
	galleryWidth    = 1920 / 2 // Reduction of transfer bandwidth
	galleryHeight   = 1080 / 2
	c0dartCacheTime = time.Duration(5 * time.Minute)
)

var resizerImages = make(map[string]bytes.Buffer)

func c0dartHandler(w http.ResponseWriter, r *http.Request) {
	globalSetHeaders(w, r)
	c0dartUpdateContext()
	if r.URL.Path == "" { // Gallery
		renderTemplate(w, "c0dart_gallery", c0dartContext)
		return
	} else if strings.HasPrefix(r.URL.Path, resizerPath) && strings.Count(r.URL.Path, "/") == 3 { // Resizer
		var fileName string
		var width, height int

		_, err := fmt.Sscanf(r.URL.Path, "resizer/%q/%d/%d", &fileName, &width, &height)
		if err != nil {
			fmt.Printf("%v: %v\n", err, r.URL.Path)
		}
		for _, c0dartImage := range c0dartContext.Images {
			if c0dartImage.Filename == fileName {
				if err == nil && width == galleryWidth && height == galleryHeight {
					resizedImageBuffer, ok := resizerImages[fileName]
					if !ok {
						file, err := os.Open(c0dartDir + fileName)
						if err != nil {
							fmt.Printf("C0dart resizer unable to open file %s with error %v\n", fileName, err)
							break
						}
						img, err := mogrify.DecodePng(file)
						if err != nil {
							fmt.Printf("C0dart resizer unable to decode file %s with error %v\n", fileName, err)
							break
						}
						resizedImage, err := img.NewResized(mogrify.Bounds{Width: galleryWidth, Height: galleryHeight})
						img.Destroy()
						if err != nil {
							fmt.Printf("C0dart resizer unable to resize file %s with error %v\n", fileName, err)
							break
						}
						mogrify.EncodePng(bufio.NewWriter(&resizedImageBuffer), resizedImage)
						resizerImages[fileName] = resizedImageBuffer
					}
					w.Header().Add("Content-Length", strconv.Itoa(resizedImageBuffer.Len()))
					w.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", c0dartImage.Filename))
					w.Header().Add("Cache-Control", fmt.Sprintf("private, max-age=%d", int(c0dartCacheTime.Seconds())))
					resizedImageBuffer.WriteTo(w)
					return
				}
				break
			}
		}
	}
	http.NotFound(w, r)
}

func c0dartUpdateContext() {
	if c0dartContext == nil || time.Now().After(c0dartContext.NextUpdate) {
		var c0dartImages []C0dartImage = nil
		if images, err := ioutil.ReadDir(c0dartDir); err == nil {
			c0dartImages = make([]C0dartImage, len(images))
			for i, img := range images {
				c0dartImages[i] = C0dartImage{
					Filename: img.Name(),
					Href:     fmt.Sprintf(c0dartHandlerPath+resizerPath+"\""+img.Name()+"\"/%d/%d", galleryWidth, galleryHeight),
					Title:    c0dartFileNameToTitle(img.Name()),
					Desc:     "", // TODO: make these real
				}
			}
		} else {
			fmt.Printf("Error reading c0dart directory: %v\n", err)
		}
		c0dartContext = &C0dartContext{
			GlobalContext: globalContext,
			Images:        c0dartImages,
			NextUpdate:    time.Now().Add(c0dartCacheTime),
		}
	}
}

func c0dartFileNameToTitle(fileName string) string {
	var outstring string = ""
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