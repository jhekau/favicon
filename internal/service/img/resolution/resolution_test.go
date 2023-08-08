package resolution_test

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 28 July 2023
 */
import (
	"errors"
	"io"
	"testing"

	logger_ "github.com/jhekau/favicon/internal/core/logger"
	logger_mock_ "github.com/jhekau/favicon/internal/core/logger/mock"
	image_test_data_ "github.com/jhekau/favicon/internal/core/test_data/image"
	resolution_ "github.com/jhekau/favicon/internal/service/img/resolution"
)

type storageReader struct{
	r io.Reader
}
func (s storageReader) Read() io.Reader {
	return s.r
}

type storage struct{
	l *logger_.Logger
	img interface{ 
		Base64Reader(l *logger_.Logger) (io.Reader, string, error)
	}
}
func (s *storage) Read() (io.ReadCloser , error) {
	r, err := image_test_data_.GetFileReader(s.img, s.l)
	return io.NopCloser(r), err
}


func TestGetResolution(t *testing.T){
	
	logger := &logger_.Logger{
		Typ: &logger_mock_.LoggerErrorf{},
	}

	for _, d := range []struct{
		img interface{ Base64Reader(l *logger_.Logger) (io.Reader, string, error) }
		w, h int
		err error
	}{
		{image_test_data_.PNG_1_1, 1, 1, nil},
		{image_test_data_.JPG_1_1, 1, 1, nil},
		{image_test_data_.PNG_16x16, 16, 16, nil},
		{image_test_data_.JPG_16_16, 16, 16, nil},
		{image_test_data_.JPG_10001_10001, 10001, 10001, nil},
		{image_test_data_.SVG, 0, 0, errors.New(`image: unknown format`)},
	}{
		w, h, err := (&resolution_.Resolution{logger}).Get(
			&storage{
				logger,
				d.img,
			},
		)
		if (err == nil && d.err != nil) || (err != nil && d.err == nil) {
			t.Fatalf(`TestGetResolution - error: '%v' data: %#v`, err, d)
		}
		if w != d.w {
			t.Fatalf(`TestGetResolution - width: '%v' data: %#v`, err, d)
		}
		if h != d.h {
			t.Fatalf(`TestGetResolution - height: '%v' data: %#v`, err, d)
		}
	}

}