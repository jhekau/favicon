package adapter_http_v1

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 27 July 2023
 */
import (
	"github.com/jhekau/gdown"

	logger_default_ "github.com/jhekau/favicon/internal/core/logs/default"
	handler_ "github.com/jhekau/favicon/pkg/core/http/v1/handler"
	thumbs_ "github.com/jhekau/favicon/pkg/service/thumbs"
)

func Run() {

	// gracefull shotdown

	logger := &logger_default_.Logger{}
	
	icons, err := thumbs_.NewThumbs_DefaultsIcons()
	if err != nil {
		// error
	}
	icons.LoggerSet(logger)
	// manifest

	handler := (&handler_.Handler{}).Handle(icons)
	server, s := gdown.HTTPNewServerWithHandler(handler)
	s.Logger(logger)
	
    if err := server.ListenAndServe(); err != nil {
        panic(err)
    }

}
