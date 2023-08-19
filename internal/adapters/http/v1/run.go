package adapter_http_v1

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 27 July 2023
 */
import (
	"flag"
	"net/http"
	"path/filepath"

	"github.com/jhekau/gdown"

	config_ "github.com/jhekau/favicon/internal/adapters/http/v1/config"
	err_ "github.com/jhekau/favicon/internal/core/err"
	logs_ "github.com/jhekau/favicon/internal/core/logs"
	"github.com/jhekau/favicon/internal/storage/files"
	handler_ "github.com/jhekau/favicon/pkg/core/http/v1/handler"
	thumbs_ "github.com/jhekau/favicon/pkg/service/thumbs"
)

const (
	appPort = ":80" // default
	workOpt = `httpv1`
)

const (
	logP = `github.com/jhekau/favicon/internal/adapters/http/v1/run.go`

	logR1 = `R1: read yaml config`
	logR2 = `R2: create defaults icons`
	logR3 = `R3: incorrect filename original image`
	logR4 = `R4: incorrect filename original image svg`
	logR5 = `R5: empty original image`
	logR6 = `R6: set object storage image`
	logR7 = `R7: set object storage image`
)

func Run() {

	log := &logs_.Logger{}

	yamlFile := flag.String("conf", ``, `config yaml file`)
	img := flag.String("img", ``, `original image`)
	svg := flag.String("svg", ``, `original image`)
	flag.Parse()

	//
	conf, err := config_.Parse(*yamlFile)
	if err != nil {
		panic(err_.Err(log, logP, logR1, err))
	}
	if conf.Port == `` {
		conf.Port = appPort
	}

	var imgF, svgF string
	if *img != `` {
		imgF, err = filepath.Abs(*img)
		if err != nil {
			panic(err_.Err(log, logP, logR3, err))
		}
	}
	if *svg != `` {
		imgF, err = filepath.Abs(*svg)
		if err != nil {
			panic(err_.Err(log, logP, logR4, err))
		}
	}
	if imgF == `` && svgF == `` {
		panic(err_.Err(log, logP, logR5, err))
	}

	//
	icons, err := thumbs_.NewThumbs_DefaultsIcons()
	if err != nil {
		panic(err_.Err(log, logP, logR2, err))
	}
	icons.LoggerSet(log)
	if imgF != `` {
		img, err := files.Files{L:log}.NewObject(imgF)
		if err != nil {
			panic(err_.Err(log, logP, logR6, err))
		}
		icons.SetOriginal(img)
	}
	if svgF != `` {
		img, err := files.Files{L:log}.NewObject(svgF)
		if err != nil {
			panic(err_.Err(log, logP, logR7, err))
		}
		icons.SetOriginal(img)
	}
	// manifest

	handler := (&handler_.Handler{L: log}).Handle(icons)
	
	// graceful shutdown
	server, s := gdown.HTTPNewServerWithHandler(handler)
	server.Addr = conf.Port 
	s.Logger(log)

	//
    if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        panic(err)
    }
}

