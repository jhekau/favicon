package converter_exec_anthonynsimon_test

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 26 August 2023
 */
import (
	"bytes"
	"fmt"
	"image"
	"io"
	"testing"
	"time"

	image_test_data_ "github.com/jhekau/favicon/internal/test/test_data/image"
	mock_logger_ "github.com/jhekau/favicon/internal/test/mocks/interfaces/logger"
	converter_exec_anthonynsimon_ "github.com/jhekau/favicon/internal/pkg/img/convert.exec/anthonynsimon"
	resolution_ "github.com/jhekau/favicon/internal/pkg/img/resolution"
	logger_ "github.com/jhekau/favicon/interfaces/logger"
	storage_ "github.com/jhekau/favicon/interfaces/storage"
	types_ "github.com/jhekau/favicon/domain/types"
	"github.com/pressly/goico"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)


func init() {
	image.RegisterFormat("ico", "\x00\x00\x01\x00", ico.Decode, ico.DecodeConfig)
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
func (s *storage) IsExists() (bool, error){
	return false, nil
}
func (s *storage) Key() storage_.StorageKey {
	return ``
}
func (s *storage) ModTime() time.Time {
	return time.Time{}
}


//
func TestConvert_CreatePNG( t *testing.T ) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logs := mock_logger_.NewMockLogger(ctrl)
	logs.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
	 

	// Data
	original := &storage{l: logs, obj: &obj{} }
	original_size := 32

	save := &storage{l: logs, obj: &obj{}}
	thumb_size := 16

	errNil := (error)(nil)

	r, err := image_test_data_.GetFileReader(image_test_data_.PNG_32x32, logs)
	require.Equal(t, err, errNil)
	io.Copy(original.obj, r)
	res := resolution_.Resolution{L: logs}
		
	// check original size
	w, h, err := res.Get(original)
	require.Equal(t, err, errNil)
	require.Equal(t, w, original_size)
	require.Equal(t, h, original_size)

	// check convert
	err = (&converter_exec_anthonynsimon_.Exec{logs}).Proc(original, save, thumb_size, types_.PNG())
	require.Equal(t, err, errNil, fmt.Sprintf(`err: %v, %v`, err, original.obj.Bytes()))

	// check save thumb size
	w, h, err = res.Get(save)
	require.Equal(t, err, errNil, fmt.Sprintf(`%v`, save.obj.Bytes()))
	require.Equal(t, w, thumb_size)
	require.Equal(t, h, thumb_size)

}

func TestConvert_CreateICO( t *testing.T ) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logs := mock_logger_.NewMockLogger(ctrl)
	logs.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
	 

	// Data
	original := &storage{l: logs, obj: &obj{} }
	original_size := 32

	save := &storage{l: logs, obj: &obj{}}
	thumb_size := 16

	errNil := (error)(nil)

	r, err := image_test_data_.GetFileReader(image_test_data_.PNG_32x32, logs)
	require.Equal(t, err, errNil)
	io.Copy(original.obj, r)
	res := resolution_.Resolution{L: logs}
		
	// check original size
	w, h, err := res.Get(original)
	require.Equal(t, err, errNil)
	require.Equal(t, w, original_size)
	require.Equal(t, h, original_size)

	// check convert
	err = (&converter_exec_anthonynsimon_.Exec{logs}).Proc(original, save, thumb_size, types_.ICO())
	require.Equal(t, err, errNil, fmt.Sprintf(`err: %v, %v`, err, original.obj.Bytes()))

	// check save thumb size
	w, h, err = res.Get(save)
	require.Equal(t, err, errNil, fmt.Sprintf(`%v`, save.obj.Bytes()))
	require.Equal(t, w, thumb_size)
	require.Equal(t, h, thumb_size)
}