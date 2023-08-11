package converters

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 24 July 2023
 * конвертеры для различных типов файлов
 */
import (
	logger_ "github.com/jhekau/favicon/internal/core/logger"
	types_ "github.com/jhekau/favicon/pkg/core/types"
	converter_ "github.com/jhekau/favicon/pkg/models/converter"
	storage_ "github.com/jhekau/favicon/pkg/models/storage"
)

//
const (
	logF   = `/internal/service/convert.go`
	logF01 = `F01: convert ico`
	logF02 = `F02: convert png`
)


// type StorageOBJ interface{
// 	Reader() (io.ReadCloser, error)
// 	Writer() (io.WriteCloser, error)
// }

// type ConverterExec interface {
// 	Proc(source, save StorageOBJ, size_px int, typ types_.FileType) error
// }

type ConverterICO struct{
	L *logger_.Logger
	// пакет/утилита, который выполняет непосредственную конвертацию изображения
	ConverterExec converter_.ConverterExec
}

// интерфейс для конвертора ICO
func (t *ConverterICO) Do(source, save storage_.StorageOBJ, size_px int, typ types_.FileType/*, conv Converter*/) (complete bool, err error) {
	if typ != types_.ICO() {
		return false, nil
	}
    if err := t.ConverterExec.Proc(source, save, size_px, typ); err != nil {
		return false, t.L.Typ.Error( logF, logF01, err )
	}
	return true, nil
}



type ConverterPNG struct{
	L *logger_.Logger
	// пакет/утилита, который выполняет непосредственную конвертацию изображения
	ConverterExec interface {
		Proc(source, save storage_.StorageOBJ, size_px int, typ types_.FileType) error
	}
}

// интерфейс для конвертора PNG
func (t *ConverterPNG) Do(source, save storage_.StorageOBJ, size_px int, typ types_.FileType/*, conv Converter*/) (complete bool, err error) {
	if typ != types_.PNG() {
		return false, nil
	}
    if err := t.ConverterExec.Proc(source, save, size_px, typ); err != nil {
		return false, t.L.Typ.Error( logF, logF02, err )
	}
	return true, nil
}
