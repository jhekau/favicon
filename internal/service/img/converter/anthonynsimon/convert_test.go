package converter_exec_anthonynsimon_test

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 26 August 2023
 */
import (
	"bytes"
	"fmt"
	"io"
	"image"
	"testing"

	logger_ "github.com/jhekau/favicon/internal/core/logger"
	logger_mock_ "github.com/jhekau/favicon/internal/core/logger/mock"
	image_test_data_ "github.com/jhekau/favicon/internal/core/test_data/image"
	types_ "github.com/jhekau/favicon/pkg/core/types"
	converter_exec_anthonynsimon_ "github.com/jhekau/favicon/internal/service/img/converter/anthonynsimon"
	resolution_ "github.com/jhekau/favicon/internal/service/img/resolution"
	"github.com/stretchr/testify/require"
	"github.com/pressly/goico"
	storage_ "github.com/jhekau/favicon/pkg/models/storage"
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
	l *logger_.Logger
	obj *obj
}
func (s *storage) Reader() (io.ReadCloser , error) {
	return io.NopCloser(bytes.NewBuffer(s.obj.Bytes())), nil
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


//
func TestConvert_CreatePNG( t *testing.T ) {

	// Data
	logger := &logger_.Logger{
		Typ: &logger_mock_.LoggerErrorf{},
	}

	original := &storage{l: logger, obj: &obj{} }
	original_size := 32

	save := &storage{l: logger, obj: &obj{}}
	thumb_size := 16

	errNil := (error)(nil)

	r, err := image_test_data_.GetFileReader(image_test_data_.PNG_32x32, logger)
	require.Equal(t, err, errNil)
	io.Copy(original.obj, r)
	res := resolution_.Resolution{L: logger}
		
	// check original size
	w, h, err := res.Get(original)
	require.Equal(t, err, errNil)
	require.Equal(t, w, original_size)
	require.Equal(t, h, original_size)

	// check convert
	err = (&converter_exec_anthonynsimon_.Exec{logger}).Proc(original, save, thumb_size, types_.PNG())
	require.Equal(t, err, errNil, fmt.Sprintf(`err: %v, %v`, err, original.obj.Bytes()))

	// check save thumb size
	w, h, err = res.Get(save)
	require.Equal(t, err, errNil, fmt.Sprintf(`%v`, save.obj.Bytes()))
	require.Equal(t, w, thumb_size)
	require.Equal(t, h, thumb_size)

}

func TestConvert_CreateICO( t *testing.T ) {

	// Data
	logger := &logger_.Logger{
		Typ: &logger_mock_.LoggerErrorf{},
	}

	original := &storage{l: logger, obj: &obj{} }
	original_size := 32

	save := &storage{l: logger, obj: &obj{}}
	thumb_size := 16

	errNil := (error)(nil)

	r, err := image_test_data_.GetFileReader(image_test_data_.PNG_32x32, logger)
	require.Equal(t, err, errNil)
	io.Copy(original.obj, r)
	res := resolution_.Resolution{L: logger}
		
	// check original size
	w, h, err := res.Get(original)
	require.Equal(t, err, errNil)
	require.Equal(t, w, original_size)
	require.Equal(t, h, original_size)

	// check convert
	err = (&converter_exec_anthonynsimon_.Exec{logger}).Proc(original, save, thumb_size, types_.ICO())
	require.Equal(t, err, errNil, fmt.Sprintf(`err: %v, %v`, err, original.obj.Bytes()))

	// check save thumb size
	w, h, err = res.Get(save)
	require.Equal(t, err, errNil, fmt.Sprintf(`%v`, save.obj.Bytes()))
	require.Equal(t, w, thumb_size)
	require.Equal(t, h, thumb_size)
}