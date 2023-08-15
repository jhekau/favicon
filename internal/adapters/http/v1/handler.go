package adapter_http_v1

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 27 July 2023
 */
import (
	"flag"
	"net/http"

	"github.com/jhekau/gdown"

	config_ "github.com/jhekau/favicon/internal/adapters/http/v1/config"
	logger_default_ "github.com/jhekau/favicon/internal/core/logs/default"
	handler_ "github.com/jhekau/favicon/pkg/core/http/v1/handler"
	thumbs_ "github.com/jhekau/favicon/pkg/service/thumbs"
)

const (
	appPort = ":80" // default
	workOpt = `httpv1`
)

const (
	logP = `github.com/jhekau/favicon/internal/adapters/http/v1/handler.go`

	logR1 = `R1: read yaml config`
	logR2 = `R2: create defaults icons`
	// logR3 = `R3: `
)

func Run() bool {

	log := &logger_default_.Logger{}

	workOption := flag.String(`adapter`, ``, `worker option`)
	yamlFile := flag.String("conf", ``, `config yaml file`)
	flag.Parse()

	//
	if *workOption != workOpt {
		return false
	}

	//
	conf, err := config_.Parse(*yamlFile)
	if err != nil {
		panic(log.Error(logP, logR1, err))
	}
	if conf.Port == `` {
		conf.Port = appPort
	}

	//
	icons, err := thumbs_.NewThumbs_DefaultsIcons()
	if err != nil {
		panic(log.Error(logP, logR1, err))
	}
	icons.LoggerSet(log)
	// manifest

	handler := (&handler_.Handler{}).Handle(icons)

	// graceful shutdown
	server, s := gdown.HTTPNewServerWithHandler(handler)
	server.Addr = conf.Port 
	s.Logger(log)

	//
    if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        panic(err)
    }

	return true
}

