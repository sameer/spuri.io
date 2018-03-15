package main

import (
	"net/http"
	"sync/atomic"
	"time"
	"strings"
	"net/url"
)

type handlerCoordinator struct {
	handlers []handlerInterface
}

func (c *handlerCoordinator) start(mux *http.ServeMux, notifyPeriod time.Duration) {
	// Do initialization
	for _, handler := range c.handlers {
		if hFS, ok := handler.(*handlerWithFinalState); ok {
			hFS.init()
		} else if hS, ok := handler.(*handlerWithUpdatableState); ok {
			hS.init()
			hS.startUpdater()
		}
	}
	// Install handlers
	for _, handlerI := range c.handlers {
		if handlerI.shouldStripPrefix() {
			if h, ok := handlerI.(*handlerWithFinalState); ok {
				h.handler = stripPrefix(h.path(), h.handler)
			} else if h, ok := handlerI.(*handlerWithUpdatableState); ok {
				h.handler = stripPrefix(h.path(), h.handler)
			} else if h, ok := handlerI.(*handlerWithoutState); ok {
				h.handler = http.StripPrefix(h.path(), h.handler)
			}
		}
		mux.Handle(handlerI.path(), handlerI)
	}
}

type state interface{}
type stateInitFunc func() state
type stateUpdateFunc func(s state) state

type handlerWithStateFunc func(w http.ResponseWriter, r *http.Request, s state)

type handlerWithFinalState struct {
	handlerGenericAttributes
	handler     handlerWithStateFunc
	state       atomic.Value
	initializer stateInitFunc
}

func (h *handlerWithFinalState) init() {
	s := h.initializer()
	h.state.Store(s)
}

func (h *handlerWithFinalState) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.handler(w, r, h.state.Load())
}

type handlerWithUpdatableState struct {
	handlerWithFinalState
	updater      stateUpdateFunc
	updatePeriod time.Duration
	updatedLast  time.Time
}

func (hS *handlerWithUpdatableState) startUpdater() {
	go func() {
		t := time.NewTicker(hS.updatePeriod)
		for {
			if now := time.Now(); hS.updatedLast.Add(hS.updatePeriod).Before(now) {
				hS.state.Store(hS.updater(hS.state.Load()))
				hS.updatedLast = now
			}
			<-t.C
		}
	}()
}

var _ handlerInterface = &handlerWithUpdatableState{}

type handlerWithoutState struct {
	handlerGenericAttributes
	handler http.Handler
}

func (h *handlerWithoutState) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.handler.ServeHTTP(w, r)
}

var _ handlerInterface = &handlerWithoutState{}

type handlerInterface interface {
	http.Handler
	path() string
	shouldStripPrefix() bool
}

type handlerGenericAttributes struct {
	pathStr     string
	stripPrefix bool
}

func (h handlerGenericAttributes) path() string {
	return h.pathStr
}

func (h handlerGenericAttributes) shouldStripPrefix() bool {
	return h.stripPrefix
}

func stripPrefix(prefix string, h handlerWithStateFunc) handlerWithStateFunc {
	if prefix == "" {
		return h
	}
	return handlerWithStateFunc(func(w http.ResponseWriter, r *http.Request, s state) {
		if p := strings.TrimPrefix(r.URL.Path, prefix); len(p) < len(r.URL.Path) {
			r2 := new(http.Request)
			*r2 = *r
			r2.URL = new(url.URL)
			*r2.URL = *r.URL
			r2.URL.Path = p
			h(w, r2, s)
		} else {
			http.NotFound(w, r)
		}
	})
}
