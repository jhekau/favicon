package resolution_test

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 28 July 2023
 */
import (
	"errors"
	"io"
	"testing"

	image_test_data_ "github.com/jhekau/favicon/internal/core/test_data/image"
	resolution_ "github.com/jhekau/favicon/internal/service/img/resolution"
)

type storageReader struct{
	r io.Reader
}
func (s storageReader) Read() io.Reader {
	return s.r
}

func TestGetResolution(t *testing.T){
	
	for _, d := range []struct{
		img interface{ Base64Reader() (io.Reader, string, error) }
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
		reader, err := image_test_data_.GetFileReader(d.img)
		if err != nil {
			t.Fatalf(`TestGetResolution:image_test_data_.GetFileReader - error: '%v' data: %#v`, err, d)
		}

		w, h, err := (&resolution_.Resolution{
			storageReader{reader},
		}).Get()
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