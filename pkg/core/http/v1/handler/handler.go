package handler

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 12 August 2023
 */
import (
	"io"
	"net/http"
	"time"

	"github.com/jhekau/favicon/pkg/core/models/logger"
)

const (
	logP = `github.com/jhekau/favicon/pkg/core/http/v1/handler/handler.go`
	logH1 = `H1: get file`
)

type Content interface {
	File(urlPath string) (content io.ReadSeekCloser, modtime time.Time, name string, exists bool, err error)
}

type Hanler struct {
	L logger.Logger
}

func (h *Hanler) Handle(w http.ResponseWriter, req *http.Request, c ...Content) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		for _, c := range c {
			content, modtime, name, exists, err := c.File(r.URL.Path)
			if err != nil {
				h.L.Error(logP, logH1, err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			if exists {
				http.ServeContent(w, r, name, modtime, content)
				return
			}
		}
		w.WriteHeader(http.StatusNotFound)
	}
}
