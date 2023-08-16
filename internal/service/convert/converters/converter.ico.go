package converters

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 24 July 2023
 * конвертер для создания ICO превьюх
 */
import (
	logger_ "github.com/jhekau/favicon/pkg/core/models/logger"
	types_ "github.com/jhekau/favicon/pkg/core/types"
	converter_ "github.com/jhekau/favicon/pkg/core/models/converter"
	storage_ "github.com/jhekau/favicon/pkg/core/models/storage"
	err_ "github.com/jhekau/favicon/internal/core/err"
)

//
const (
	logFI   = `internal/service/convert/converters/converter.ico.go`
	logFI01 = `F01: convert ico`
)


type ConverterICO struct{
	L logger_.Logger
	// пакет/утилита, который выполняет непосредственную конвертацию изображения
	ConverterExec converter_.ConverterExec
}

// интерфейс для конвертора ICO
func (t *ConverterICO) Do(source, save storage_.StorageOBJ, size_px int, typ types_.FileType/*, conv Converter*/) (complete bool, err error) {
	if typ != types_.ICO() {
		return false, nil
	}
    if err := t.ConverterExec.Proc(source, save, size_px, typ); err != nil {
		return false, err_.Err( t.L, logFI, logFI01, err )
	}
	return true, nil
}
