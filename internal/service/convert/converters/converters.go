package converters

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 24 July 2023
 * конвертеры для различных типов файлов
 */
import (
	err_ "github.com/jhekau/favicon/internal/core/err"
	types_ "github.com/jhekau/favicon/internal/core/types"
)

// 
const (
	logF01 = `F01: convert ico`
	logF02 = `F02: convert png`
)
func errF(i... interface{}) error {
	return err_.Err(err_.TypeError, `/internal/service/convert.go`, i...)
} 


// функция, которая непосредственно конвертирует изображение.
// Можно использовать внешнюю библиотеку, или внешний бинарник
// type Converter interface {
// 	Proc(source, save types_.FilePath, size_px int, typ types_.FileType) error
// }

type ConverterICO struct{
	// пакет/утилита, который выполняет непосредственную конвертацию изображения
	ConverterExec interface {
		Proc(source, save types_.FilePath, size_px int, typ types_.FileType) error
	}
}

// интерфейс для конвертора ICO
func (t *ConverterICO) Do(source, save types_.FilePath, size_px int, typ types_.FileType/*, conv Converter*/) (complete bool, err error) {
	if typ != types_.ICO() {
		return false, nil
	}
    if err := t.ConverterExec.Proc(source, save, size_px, typ); err != nil {
		return false, errF(logF01, err)
	}
	return true, nil
}



type ConverterPNG struct{
	ConverterExec interface {
		Proc(source, save types_.FilePath, size_px int, typ types_.FileType) error
	}
}

// интерфейс для конвертора PNG
func (t *ConverterPNG) Do(source, save types_.FilePath, size_px int, typ types_.FileType/*, conv Converter*/) (complete bool, err error) {
	if typ != types_.PNG() {
		return false, nil
	}
    if err := t.ConverterExec.Proc(source, save, size_px, typ); err != nil {
		return false, errF(logF02, err)
	}
	return true, nil
}
