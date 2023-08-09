package converters_test

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 26 July 2023
 */
import (
	"bytes"
	"errors"
	"io"
	"testing"

	logger_ "github.com/jhekau/favicon/internal/core/logger"
	logger_mock_ "github.com/jhekau/favicon/internal/core/logger/mock"
	image_test_data_ "github.com/jhekau/favicon/internal/core/test_data/image"
	types_ "github.com/jhekau/favicon/internal/core/types"
	mock_converters_ "github.com/jhekau/favicon/internal/mocks/intr/service/convert/converters"
	converters_ "github.com/jhekau/favicon/internal/service/convert/converters"
	domain_ "github.com/jhekau/favicon/pkg/domain"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// image test data

type obj struct {
	bytes.Buffer
}
func (o *obj) Close() error {
	return nil
}

// image test data
type storage struct{
	l *logger_.Logger
	obj *obj
	key domain_.StorageKey
}
func (s *storage) Reader() (io.ReadCloser , error) {
	return io.NopCloser(bytes.NewBuffer(s.obj.Bytes())), nil
}
func (s *storage) Writer() (io.WriteCloser, error) {
	return s.obj, nil
}
func (s *storage) Key() domain_.StorageKey {
	return s.key
}


func readTestImage(img image_test_data_.Imgb64, logger *logger_.Logger) ( *storage, error ) {
	obj := &storage{l: logger, obj: &obj{}}
	r, err := image_test_data_.GetFileReader(image_test_data_.PNG_32x32, logger)
	if err != nil {
		return nil, err
	}
	io.Copy(obj.obj, r)
	return obj, nil
}


func TestConverterICOUnit( t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	png16 := image_test_data_.PNG_16x16
	errNil := (error)(nil)
	logger := &logger_.Logger{
		Typ: &logger_mock_.LoggerErrorf{},
	}
	
	for _, d := range []struct{
		source image_test_data_.Imgb64
		save_key domain_.StorageKey
		size_px int
		typ types_.FileType
		converter_error error
		result_complite bool
		result_error error
	}{
		{ png16, `1_preview.png`, 16, types_.PNG(), nil, false, nil },
		{ png16, `1_preview.png`, 16, types_.ICO(), nil, true, nil },
		{ png16, `1_preview.png`, 16, types_.SVG(), nil, false, nil },
		{ png16, `1_preview.png`, 16, types_.ICO(), errors.New(`error`), false, errors.New(`error`) },
	}{

		orig, err := readTestImage(d.source, logger)
		require.Equal(t, err, errNil)

		save := &storage{ key: d.save_key }

		convExec := mock_converters_.NewMockConverterExec(ctrl)
		convExec.EXPECT().Proc(orig, save, d.size_px, d.typ).Return(d.converter_error).AnyTimes()

		res, err := (&converters_.ConverterICO{
			L: &logger_.Logger{
				Typ: &logger_mock_.LoggerErrorf{},
			},
			ConverterExec: convExec,
		}).Do(orig, save, d.size_px, d.typ)

		if ( err != nil && d.result_error == nil ) || ( err == nil && d.result_error != nil ) {
			t.Fatalf(`TestConverterICOUnit:error - result: '%v', err: '%v', testdata: '%#v'`, res, err, d)
		}
		if res != d.result_complite {
			t.Fatalf(`TestConverterICOUnit:result - result: '%v', err: '%v', testdata: '%#v'`, res, err, d)
		}
	}
}


func TestConverterPNGUnit( t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	png16 := image_test_data_.PNG_16x16
	errNil := (error)(nil)
	logger := &logger_.Logger{
		Typ: &logger_mock_.LoggerErrorf{},
	}

	for _, d := range []struct{
		source image_test_data_.Imgb64
		save_key domain_.StorageKey
		size_px int
		typ types_.FileType
		converter_error error
		result_complite bool
		result_error error
	}{
		{ png16, `1_preview.png`, 16, types_.PNG(), nil, true, nil },
		{ png16, `1_preview.png`, 16, types_.ICO(), nil, false, nil },
		{ png16, `1_preview.png`, 16, types_.SVG(), nil, false, nil },
		{ png16, `1_preview.png`, 16, types_.PNG(), errors.New(`error`), false, errors.New(`error`) },
	}{

		orig, err := readTestImage(d.source, logger)
		require.Equal(t, err, errNil)

		save := &storage{ key: d.save_key }

		convExec := mock_converters_.NewMockConverterExec(ctrl)
		convExec.EXPECT().Proc(orig, save, d.size_px, d.typ).Return(d.converter_error).AnyTimes()

		res, err := (&converters_.ConverterPNG{
			L: &logger_.Logger{
				Typ: &logger_mock_.LoggerErrorf{},
			},
			ConverterExec: convExec,
		}).Do(orig, save, d.size_px, d.typ)
		
		if ( err != nil && d.result_error == nil ) || ( err == nil && d.result_error != nil ) {
			t.Fatalf(`TestConverterICOUnit:error - result: '%v', err: '%v', testdata: '%#v'`, res, err, d)
		}
		if res != d.result_complite {
			t.Fatalf(`TestConverterICOUnit:result - result: '%v', err: '%v', testdata: '%#v'`, res, err, d)
		}
	}
}