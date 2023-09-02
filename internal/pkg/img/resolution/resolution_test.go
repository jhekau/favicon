package resolution_test

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 28 July 2023
 */
import (
	"bytes"
	"errors"
	"io"
	"testing"
	"time"

	image_test_data_ "github.com/jhekau/favicon/internal/test/test_data/image"
	mock_logger_ "github.com/jhekau/favicon/internal/test/mocks/interfaces/logger"
	resolution_ "github.com/jhekau/favicon/internal/pkg/img/resolution"
	logger_ "github.com/jhekau/favicon/interfaces/logger"
	storage_ "github.com/jhekau/favicon/interfaces/storage"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type storageReader struct{
	r io.Reader
}
func (s storageReader) Read() io.Reader {
	return s.r
}

type obj struct {
	bytes.Buffer
}
func (o *obj) Close() error {
	return nil
}

// image test data
type storage struct{
	l logger_.Logger
	obj *obj
}

type reader struct {
	io.ReadCloser
}
func (r *reader) Seek(offset int64, whence int) (int64, error){
	return 0,nil
}

func (s *storage) Reader() (io.ReadSeekCloser , error) {
	return &reader{
		io.NopCloser(bytes.NewBuffer(s.obj.Bytes())),
	}, nil
}
func (s *storage) Writer() (io.WriteCloser, error) {
	return s.obj, nil
}
func (s *storage) Key() storage_.StorageKey{
	return ``
}
func (s *storage) IsExists() ( bool, error ){
	return false, nil
}
func (s *storage) ModTime() time.Time {
	return time.Time{}
}


func TestGetResolution(t *testing.T){
	
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logs := mock_logger_.NewMockLogger(ctrl)
	logs.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
	 

	errNil := (error)(nil)

	for _, d := range []struct{
		img interface{ Base64Reader(l logger_.Logger) (io.Reader, string, error) }
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
		img := &storage{l: logs, obj: &obj{} }

		r, err := image_test_data_.GetFileReader(d.img, logs)
		require.Equal(t, err, errNil)
		io.Copy(img.obj, r)

		w, h, err := (&resolution_.Resolution{logs}).Get(img)
		
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