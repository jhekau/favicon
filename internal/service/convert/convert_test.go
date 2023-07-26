package convert_test

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 26 July 2023
 */
import (
	"errors"
	"testing"

	types_ "github.com/jhekau/favicon/internal/core/types"
	convert_ "github.com/jhekau/favicon/internal/service/convert"
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
}
func (c convertType) Do(_, _ types_.FilePath, _ int, typ types_.FileType, _ convert_.Converter) (complete bool, err error) {
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
func (c checkSource) Check(_ types_.FilePath, _ types_.FileType, _ int) error {
	return c.err
}



// Unit
func TestConvertUnit( t *testing.T ) {

	for _, d := range []struct{
		// 
		source 			types_.FilePath
		source_svg 		types_.FilePath
		save 			types_.FilePath
		typ 			types_.FileType
		size_px 		int
		converters 		[]convert_.Converters
		check_preview 	checkPreview
		check_source 	checkSource
		// result
		complite 		bool
		complite_err 	error
	}{
		{ 	// -----------------------------------------------------------
			`TestConvertUnit/1.jpg`, ``, `1.ico`, types_.ICO(), 16,
			[]convert_.Converters{
				{converter{nil}, convertType{types_.ICO(), nil}},
				{converter{nil}, convertType{types_.SVG(), nil}},
			},
			checkPreview{nil},
			checkSource{nil},
			true, nil,
		},
		{ 	// нулевой размер для готовой превьюхи -----------------------
			`TestConvertUnit/2.jpg`, ``, `2.ico`, types_.ICO(), 0,
			[]convert_.Converters{
				{converter{nil}, convertType{types_.ICO(), nil}},
				{converter{nil}, convertType{types_.SVG(), nil}},
			},
			checkPreview{nil},
			checkSource{nil},
			false, nil,
		},
		{ 	// -----------------------------------------------------------
			``, `TestConvertUnit/2.svg`, `2.svg`, types_.SVG(), 16,
			[]convert_.Converters{
				{converter{nil}, convertType{types_.ICO(), nil}},
				{converter{nil}, convertType{types_.SVG(), nil}},
			},
			checkPreview{nil},
			checkSource{nil},
			true, nil,
		},
		{ 	// исходных файлов нет, для нарезания превьюхи ---------------
			``, ``, `2.svg`, types_.SVG(), 16,
			[]convert_.Converters{
				{converter{nil}, convertType{types_.ICO(), nil}},
				{converter{nil}, convertType{types_.SVG(), nil}},
			},
			checkPreview{nil},
			checkSource{nil},
			false, errors.New(`error`),
		},
		{ 	// ошибка из проверки параметров нарезаемой превьюхи ---------
			`TestConvertUnit/3.jpg`, ``, `3.ico`, types_.ICO(), 16,
			[]convert_.Converters{
				{converter{nil}, convertType{types_.ICO(), nil}},
				{converter{nil}, convertType{types_.SVG(), nil}},
			},
			checkPreview{errors.New(`error`)},
			checkSource{nil},
			false, errors.New(`error`),
		},
		{ 	// ошибка при проверке оригинального файла, с которого нарезается превьюха 
			`TestConvertUnit/3.jpg`, ``, `3.ico`, types_.ICO(), 16,
			[]convert_.Converters{
				{converter{nil}, convertType{types_.ICO(), nil}},
				{converter{nil}, convertType{types_.SVG(), nil}},
			},
			checkPreview{nil},
			checkSource{errors.New(`error`)},
			false, errors.New(`error`),
		},
		{ 	// ошибка декоратора, который проверяет тип нарезаемой превьюхи и запускает свой конвертер
			`TestConvertUnit/3.jpg`, ``, `3.ico`, types_.ICO(), 16,
			[]convert_.Converters{
				{converter{nil}, convertType{types_.ICO(), errors.New(`error`)}},
				{converter{nil}, convertType{types_.SVG(), nil}},
			},
			checkPreview{nil},
			checkSource{nil},
			false, errors.New(`error`),
		},
		{ 	// отсутствуют конвертеры ------------------------------------
			`TestConvertUnit/3.jpg`, ``, `3.ico`, types_.ICO(), 16,
			nil,
			checkPreview{nil},
			checkSource{nil},
			false, nil,
		},
	}{
		complite, err := convert_.Convert( 
			d.source, d.source_svg, d.save,
			d.typ,
			d.size_px,
			d.converters,
			d.check_preview,
			d.check_source,
		)

		if (err == nil && d.complite_err != nil) || (err != nil && d.complite_err == nil) {
			t.Fatalf(`TestConvertUnit - error: complite: %v err: '%v' data: %#v`, d.complite, err, d)
		}
		
		if complite != d.complite {
			t.Fatalf(`TestConvertUnit - complite: complite: %v err: '%v' data: %#v`, d.complite, err, d)
		}
	}
}



// Integration

