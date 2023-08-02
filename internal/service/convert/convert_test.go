package convert_test

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 26 July 2023
 */
import (
	"errors"
	"testing"

	logger_ "github.com/jhekau/favicon/internal/core/logger"
	logger_mock_ "github.com/jhekau/favicon/internal/core/logger/mock"
	types_ "github.com/jhekau/favicon/internal/core/types"
	mock_convert_ "github.com/jhekau/favicon/internal/mocks/intr/service/convert"
	mock_checks_ "github.com/jhekau/favicon/internal/mocks/intr/service/convert/checks"
	convert_ "github.com/jhekau/favicon/internal/service/convert"
	checks_ "github.com/jhekau/favicon/internal/service/convert/checks"
	domain_ "github.com/jhekau/favicon/pkg/domain"
	"go.uber.org/mock/gomock"
)

// конвертер, который непосредственно занимается конвертацией
type converter struct {
	err error
}
func (c converter) Proc(_, _ types_.FilePath, _ int, _ types_.FileType) error {
	return c.err
}


// проверяют тип, в который хотим конвертировать и запускают конвертер
type convertType struct{
	typ types_.FileType
	err error 
	converterExec converter
}
func (c convertType) Do(_, _ types_.FilePath, _ int, typ types_.FileType) (complete bool, err error) {
	return c.typ == typ, c.err
}

// проверка валидности запрашиваемой превьюхи
type checkPreview struct {
	err error
}
func (c checkPreview) Check(_ types_.FileType, _ int) error {
	return c.err
}

// проверка валидности исходника для нарезания превьюхи
type checkSource struct {
	err error
}
func (c checkSource) Check(_ types_.FilePath, _ bool, _ int) error {
	return c.err
}



// Unit
func TestConvertUnit( t *testing.T ) {

	for _, d := range []struct{
		// test data
		source 			types_.FilePath
		save 			types_.FilePath
		originalSVG    	bool
		typThumb		types_.FileType
		size_px 		int
		converters 		[]convert_.ConverterT
		check_preview 	checkPreview
		check_source 	checkSource
		// result
		complite_err 	error
	}{
		{ 	// -----------------------------------------------------------
			`TestConvertUnit/1.jpg`, `1.ico`, false, types_.ICO(), 16,
			[]convert_.ConverterT{
				&convertType{types_.ICO(), nil, converter{nil}},
				&convertType{types_.SVG(), nil, converter{nil}},
			},
			checkPreview{nil},
			checkSource{nil},
			nil,
		},
		{ 	// нулевой размер для готовой превьюхи -----------------------
			`TestConvertUnit/2.jpg`, `2.ico`, false, types_.ICO(), 0,
			[]convert_.ConverterT{
				&convertType{types_.ICO(), nil, converter{nil}},
				&convertType{types_.SVG(), nil, converter{nil}},
			},
			checkPreview{nil},
			checkSource{nil},
			errors.New(`error`),
		},
		{ 	// -----------------------------------------------------------
			`TestConvertUnit/2.svg`, `2.svg`, false, types_.SVG(), 16,
			[]convert_.ConverterT{
				&convertType{types_.ICO(), nil, converter{nil}},
				&convertType{types_.SVG(), nil, converter{nil}},
			},
			checkPreview{nil},
			checkSource{nil},
			nil,
		},
		{ 	// исходных файлов нет, для нарезания превьюхи ---------------
			``, `2.svg`, false, types_.SVG(), 16,
			[]convert_.ConverterT{
				&convertType{types_.ICO(), nil, converter{nil}},
				&convertType{types_.SVG(), nil, converter{nil}},
			},
			checkPreview{nil},
			checkSource{nil},
			errors.New(`error`),
		},
		{ 	// ошибка из проверки параметров нарезаемой превьюхи ---------
			`TestConvertUnit/3.jpg`, `3.ico`, false, types_.ICO(), 16,
			[]convert_.ConverterT{
				&convertType{types_.ICO(), nil, converter{nil}},
				&convertType{types_.SVG(), nil, converter{nil}},
			},
			checkPreview{errors.New(`error`)},
			checkSource{nil},
			errors.New(`error`),
		},
		{ 	// ошибка при проверке оригинального файла, с которого нарезается превьюха 
			`TestConvertUnit/3.jpg`, `3.ico`, false, types_.ICO(), 16,
			[]convert_.ConverterT{
				&convertType{types_.ICO(), nil, converter{nil}},
				&convertType{types_.SVG(), nil, converter{nil}},
			},
			checkPreview{nil},
			checkSource{errors.New(`error`)},
			errors.New(`error`),
		},
		{ 	// ошибка декоратора, который проверяет тип нарезаемой превьюхи и запускает свой конвертер
			`TestConvertUnit/3.jpg`, `3.ico`, false, types_.ICO(), 16,
			[]convert_.ConverterT{
				&convertType{types_.ICO(), errors.New(`error`), converter{nil}},
				&convertType{types_.SVG(), nil, converter{nil}},
			},
			checkPreview{nil},
			checkSource{nil},
			errors.New(`error`),
		},
		{ 	// отсутствуют конвертеры ------------------------------------
			`TestConvertUnit/3.jpg`, `3.ico`, false, types_.ICO(), 16,
			nil,
			checkPreview{nil},
			checkSource{nil},
			errors.New(`error`),
		},
	}{
		err := (&convert_.Converter{
			&logger_.Logger{
				Typ: &logger_mock_.LoggerErrorf{},
			},
			d.converters,
			d.check_preview,
			d.check_source,
		}).Do( 
			d.source, d.save,
			d.originalSVG,
			d.typThumb,
			d.size_px,
		)

		if (err == nil && d.complite_err != nil) || (err != nil && d.complite_err == nil) {
			t.Fatalf(`TestConvertUnit - error: '%v' data: %#v`, err, d)
		}
	}
}



// Integration

type CheckSourceCacheDisable struct{}
func (c CheckSourceCacheDisable) Status(_ domain_.StorageKey, _ bool, _ int) (bool, error) {
	return false, nil
}
func (c CheckSourceCacheDisable) SetErr(_ domain_.StorageKey, _ bool, _ int, err error) error {
	return err
}

type resolution struct{
	f func() (w int, h int, err error)
}
func (r resolution) Get() (w int, h int, err error){
	return r.f()
}

func TestConvertIntegration( t *testing.T ) {

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mock_check_source := mock_convert_.NewMockCheckSource(ctl)
	// mock_checks_.MockStorageOBJ

	for _, d := range []struct{
		// test data
		source 			types_.FilePath
		save 			types_.FilePath
		originalSVG    	bool
		typThumb		types_.FileType
		size_px 		int
		converters 		[]convert_.ConverterT
		check_preview 	interface { Check(typ types_.FileType, size_px int) error }
		check_source 	interface { Check(_ mock_checks_.MockStorageOBJ, _ bool, _ int) error }
		// result
		complite_err 	error
	}{
		{ 	// + checkPreview ---------------------------------------------
			`TestConvertUnit/1.jpg`, `1.ico`, false, types_.ICO(), 16,
			[]convert_.ConverterT{
				&convertType{types_.ICO(), nil, converter{nil}},
				&convertType{types_.SVG(), nil, converter{nil}},
			},
			checks_.Preview{
				L: &logger_.Logger{
					Typ: &logger_mock_.LoggerErrorf{},
				},
			},
			mock_check_source,
			nil,
		},
		{ 	// + checkPreview, ошибка, превьюха размером меньше, чем нужно -
			`TestConvertUnit/1.jpg`, `1.ico`, false, types_.ICO(), 1,
			[]convert_.ConverterT{
				&convertType{types_.ICO(), nil, converter{nil}},
				&convertType{types_.SVG(), nil, converter{nil}},
			},
			checks_.Preview{
				L: &logger_.Logger{
					Typ: &logger_mock_.LoggerErrorf{},
				},
			},
			mock_check_source,
			errors.New(`error`),
		},
		{ 	// + checkSource -----------------------------------------------
			`TestConvertUnit/2.jpg`, `2.ico`, false, types_.ICO(), 16,
			[]convert_.ConverterT{
				&convertType{types_.ICO(), nil, converter{nil}},
				&convertType{types_.SVG(), nil, converter{nil}},
			},
			checkPreview{nil},
			&checks_.Source{ 
				L: &logger_.Logger{
					Typ: &logger_mock_.LoggerErrorf{},
				},
				Cache: CheckSourceCacheDisable{},
				// Resolution: resolution{f: func() (w int, h int, err error){ return 16, 16, nil } },
			},
			nil,
		},
	}{
		err := (&convert_.Converter{
			&logger_.Logger{
				Typ: &logger_mock_.LoggerErrorf{},
			},
			d.converters,
			d.check_preview,
			d.check_source,
		}).Do( 
			d.source, d.save,
			d.originalSVG,
			d.typThumb,
			d.size_px,
		)

		if (err == nil && d.complite_err != nil) || (err != nil && d.complite_err == nil) {
			t.Fatalf(`TestConvertIntegration - error: '%v' data: %#v`, err, d)
		}
	}
}
