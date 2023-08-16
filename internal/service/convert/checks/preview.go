package checks

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 24 July 2023
 * проверка исходного изображения
 */
import (
	"fmt"

	config_ "github.com/jhekau/favicon/internal/config"
	err_ "github.com/jhekau/favicon/internal/core/err"
	logger_ "github.com/jhekau/favicon/pkg/core/models/logger"
	types_ "github.com/jhekau/favicon/pkg/core/types"
)

const (
	logPP  = `/internal/service/convert/checks/preview.go`
	logP01 = `P01: incorrect resolution source file`
	// logP02 = `P02: `
	// logP03 = `P03: `
	// logP04 = `P04: `
)


type Preview struct{
	L logger_.Logger
}

func (p Preview) Check(typ types_.FileType, size_px int) error {

	if typ == types_.SVG() {
		return nil
	}

	if 	size_px < config_.ImagePreviewResolutionMin || 
		size_px > config_.ImagePreviewResolutionMax {
		return err_.Err(
			p.L,
			logPP,
			fmt.Sprintf(
				`Minimum Resolution: %d, Maximum Resolution %d, Current Value: %d`,
				config_.ImagePreviewResolutionMin,
				config_.ImagePreviewResolutionMax,
				size_px,
			),
			logP01,
		)
	}
	return nil
}
