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
	"time"

	mock_logger_ "github.com/jhekau/favicon/internal/test/mocks/interfaces/logger"
	image_test_data_ "github.com/jhekau/favicon/internal/test/test_data/image"
	mock_converter_ "github.com/jhekau/favicon/internal/test/mocks/interfaces/converter"
	converters_ "github.com/jhekau/favicon/internal/pkg/img/convert/converters"
	types_ "github.com/jhekau/favicon/domain/types"
	logger_ "github.com/jhekau/favicon/interfaces/logger"
	storage_ "github.com/jhekau/favicon/interfaces/storage"
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
type storage struct {
	l   logger_.Logger
	obj *obj
	key storage_.StorageKey
}

type reader struct {
	io.ReadCloser
}
func (r *reader) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

func (s *storage) Reader() (io.ReadSeekCloser, error) {
	return &reader{
		io.NopCloser(bytes.NewBuffer(s.obj.Bytes())),
	}, nil
}
func (s *storage) Writer() (io.WriteCloser, error) {
	return s.obj, nil
}
func (s *storage) Key() storage_.StorageKey {
	return s.key
}
func (s *storage) IsExists() (bool, error) {
	return true, nil
}
func (s *storage) ModTime() time.Time {
	return time.Time{}
}

func readTestImage(img image_test_data_.Imgb64, logger logger_.Logger) (*storage, error) {
	obj := &storage{l: logger, obj: &obj{}}
	r, err := image_test_data_.GetFileReader(image_test_data_.PNG_32x32, logger)
	if err != nil {
		return nil, err
	}
	io.Copy(obj.obj, r)
	return obj, nil
}

func TestConverterICOUnit(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logs := mock_logger_.NewMockLogger(ctrl)
	logs.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
	 

	png16 := image_test_data_.PNG_16x16
	errNil := (error)(nil)

	for _, d := range []struct {
		source          image_test_data_.Imgb64
		save_key        storage_.StorageKey
		size_px         int
		typ             types_.FileType
		converter_error error
		result_complite bool
		result_error    error
	}{
		{png16, `1_preview.png`, 16, types_.PNG(), nil, false, nil},
		{png16, `1_preview.png`, 16, types_.ICO(), nil, true, nil},
		{png16, `1_preview.png`, 16, types_.SVG(), nil, false, nil},
		{png16, `1_preview.png`, 16, types_.ICO(), errors.New(`error`), false, errors.New(`error`)},
	} {

		orig, err := readTestImage(d.source, logs)
		require.Equal(t, err, errNil)

		save := &storage{key: d.save_key}

		convExec := mock_converter_.NewMockConverterExec(ctrl)
		convExec.EXPECT().Proc(orig, save, d.size_px, d.typ).Return(d.converter_error).AnyTimes()

		res, err := (&converters_.ConverterICO{
			L: logs,
			ConverterExec: convExec,
		}).Do(orig, save, d.size_px, d.typ)

		if (err != nil && d.result_error == nil) || (err == nil && d.result_error != nil) {
			t.Fatalf(`TestConverterICOUnit:error - result: '%v', err: '%v', testdata: '%#v'`, res, err, d)
		}
		if res != d.result_complite {
			t.Fatalf(`TestConverterICOUnit:result - result: '%v', err: '%v', testdata: '%#v'`, res, err, d)
		}
	}
}

func TestConverterPNGUnit(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	png16 := image_test_data_.PNG_16x16
	errNil := (error)(nil)
	
	logs := mock_logger_.NewMockLogger(ctrl)
	logs.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
	 

	for _, d := range []struct {
		source          image_test_data_.Imgb64
		save_key        storage_.StorageKey
		size_px         int
		typ             types_.FileType
		converter_error error
		result_complite bool
		result_error    error
	}{
		{png16, `1_preview.png`, 16, types_.PNG(), nil, true, nil},
		{png16, `1_preview.png`, 16, types_.ICO(), nil, false, nil},
		{png16, `1_preview.png`, 16, types_.SVG(), nil, false, nil},
		{png16, `1_preview.png`, 16, types_.PNG(), errors.New(`error`), false, errors.New(`error`)},
	} {

		orig, err := readTestImage(d.source, logs)
		require.Equal(t, err, errNil)

		save := &storage{key: d.save_key}

		convExec := mock_converter_.NewMockConverterExec(ctrl)
		convExec.EXPECT().Proc(orig, save, d.size_px, d.typ).Return(d.converter_error).AnyTimes()

		res, err := (&converters_.ConverterPNG{
			L: logs,
			ConverterExec: convExec,
		}).Do(orig, save, d.size_px, d.typ)

		if (err != nil && d.result_error == nil) || (err == nil && d.result_error != nil) {
			t.Fatalf(`TestConverterICOUnit:error - result: '%v', err: '%v', testdata: '%#v'`, res, err, d)
		}
		if res != d.result_complite {
			t.Fatalf(`TestConverterICOUnit:result - result: '%v', err: '%v', testdata: '%#v'`, res, err, d)
		}
	}
}
