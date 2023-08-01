package converters_test

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 26 July 2023
 */
import (
	"errors"
	"testing"

	types_ "github.com/jhekau/favicon/internal/core/types"
	converters_ "github.com/jhekau/favicon/internal/service/convert/converters"
)

type converter struct{
	err error
}

func (c converter) Proc(source, save types_.FilePath, size_px int, typ types_.FileType) error {
	return c.err
}

func TestConverterICOUnit( t *testing.T) {

	for _, d := range []struct{
		source, save types_.FilePath
		size_px int
		typ types_.FileType
		conv converter
		result_complite bool
		result_error error
	}{
		{ `1.jpg`, `1_preview.png`, 16, types_.PNG(), converter{ nil }, false, nil },
		{ `2.jpg`, `1_preview.png`, 16, types_.ICO(), converter{ nil }, true, nil },
		{ `3.jpg`, `1_preview.png`, 16, types_.SVG(), converter{ nil }, false, nil },
		{ `3.jpg`, `1_preview.png`, 16, types_.ICO(), converter{ errors.New(`error`) }, false, errors.New(`error`) },
	}{
		res, err := (&converters_.ConverterICO{
			ConverterExec: d.conv,
		}).Do(d.source, d.save, d.size_px, d.typ)

		if ( err != nil && d.result_error == nil ) || ( err == nil && d.result_error != nil ) {
			t.Fatalf(`TestConverterICOUnit:error - result: '%v', err: '%v', testdata: '%#v'`, res, err, d)
		}
		if res != d.result_complite {
			t.Fatalf(`TestConverterICOUnit:result - result: '%v', err: '%v', testdata: '%#v'`, res, err, d)
		}
	}
}


func TestConverterPNGUnit( t *testing.T) {

	for _, d := range []struct{
		source, save types_.FilePath
		size_px int
		typ types_.FileType
		conv converter
		result_complite bool
		result_error error
	}{
		{ `1.jpg`, `1_preview.png`, 16, types_.PNG(), converter{ nil }, true, nil },
		{ `2.jpg`, `1_preview.png`, 16, types_.ICO(), converter{ nil }, false, nil },
		{ `3.jpg`, `1_preview.png`, 16, types_.SVG(), converter{ nil }, false, nil },
		{ `3.jpg`, `1_preview.png`, 16, types_.PNG(), converter{ errors.New(`error`) }, false, errors.New(`error`) },
	}{
		res, err := (&converters_.ConverterPNG{
			ConverterExec: d.conv,
		}).Do(d.source, d.save, d.size_px, d.typ)
		
		if ( err != nil && d.result_error == nil ) || ( err == nil && d.result_error != nil ) {
			t.Fatalf(`TestConverterICOUnit:error - result: '%v', err: '%v', testdata: '%#v'`, res, err, d)
		}
		if res != d.result_complite {
			t.Fatalf(`TestConverterICOUnit:result - result: '%v', err: '%v', testdata: '%#v'`, res, err, d)
		}
	}
}