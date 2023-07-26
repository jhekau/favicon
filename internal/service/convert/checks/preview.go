package checks

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 24 July 2023
 * проверка исходного изображения
 */
import (
	// err_ "github.com/jhekau/favicon/internal/core/err"
	types_ "github.com/jhekau/favicon/internal/core/types"
	config_ "github.com/jhekau/favicon/internal/config"
)
/*
const (
	logP01 = `P01: `
	logP02 = `P02: `
	logP03 = `P03: `
	logP04 = `P04: `
)
func errP(i... interface{}) error {
	return err_.Err(err_.TypeError, `/internal/service/convert/checks/preview.go`, i...)
}
*/

type Preview struct {}

func (p Preview) Check(typ types_.FileType, size_px int) (bool, error) {

	if typ == types_.SVG() {
		return true, nil
	}

	if 	size_px < config_.ImagePreviewResolutionMin || 
		size_px > config_.ImagePreviewResolutionMax {
		return false, nil
	}

	return true, nil
}
