package convert_test

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 26 July 2023
 */
import (
	"errors"
	"fmt"
	"io"
	"testing"
	"time"

	err_ "github.com/jhekau/favicon/internal/pkg/err"
	image_test_data_ "github.com/jhekau/favicon/internal/test/test_data/image"
	mock_convert_ "github.com/jhekau/favicon/internal/test/mocks/intr/service/convert"
	mock_converter_ "github.com/jhekau/favicon/internal/test/mocks/pkg/core/models/converter"
	mock_logger_ "github.com/jhekau/favicon/internal/test/mocks/pkg/core/models/logger"
	convert_ "github.com/jhekau/favicon/internal/pkg/img/convert"
	converter_ "github.com/jhekau/favicon/pkg/core/models/converter"
	logger_ "github.com/jhekau/favicon/pkg/core/models/logger"
	storage_ "github.com/jhekau/favicon/pkg/core/models/storage"
	types_ "github.com/jhekau/favicon/pkg/core/types"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// image test data
type storage struct {
	l   logger_.Logger
	img interface {
		Base64Reader(l logger_.Logger) (io.Reader, string, error)
	}
	ifExist      bool
	ifExistError error
	storageKey   storage_.StorageKey
}

func (s *storage) IsExists() (bool, error) {
	return s.ifExist, s.ifExistError
}
func (s *storage) Key() storage_.StorageKey {
	return s.storageKey
}
type reader struct {
	io.ReadCloser
}
func (r *reader) Seek(offset int64, whence int) (int64, error){
	return 0,nil
}

func (s *storage) Reader() (io.ReadSeekCloser, error) {
	r, err := image_test_data_.GetFileReader(s.img, s.l)
	return &reader{
		io.NopCloser(r),
	}, err
}
func (s *storage) Writer() (io.WriteCloser, error) {
	var w io.WriteCloser
	return w, nil
}
func (s *storage) ModTime() time.Time {
	return time.Time{}
}

// Unit *******

func TestUnit_Convert_JPGxICO(t *testing.T) {
	
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logs := mock_logger_.NewMockLogger(ctrl)
	logs.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()

	// Data
	original := &storage{
		l:          logs,
		img:        image_test_data_.PNG_16x16,
		storageKey: `TestConvertUnit/1.jpg`,
	}
	original_svg := false

	save := &storage{
		l:          logs,
		storageKey: `testThumb/1.ico`,
	}

	size := 16
	typ := types_.ICO()

	errNil := (error)(nil)

	// Mock

	mock_converters := mock_converter_.NewMockConverterTyp(ctrl)
	mock_converters.EXPECT().Do(original, save, size, typ).Return(true, errNil)

	mock_check_preview := mock_convert_.NewMockCheckPreview(ctrl)
	mock_check_preview.EXPECT().Check(typ, size)

	mock_check_source := mock_convert_.NewMockCheckSource(ctrl)
	mock_check_source.EXPECT().Check(original, original_svg, size)

	// Test
	err := (&convert_.Converter{
		L: logs,
		Converters: []converter_.ConverterTyp{
			mock_converters,
		},
		CheckPreview: mock_check_preview,
		CheckSource:  mock_check_source,
	}).Do(
		original, save, original_svg, typ, size,
	)

	require.Equal(t, err, errNil,
		fmt.Sprintf(`error: return '%v', want: '%v'`, err, errNil))
}

func TestUnit_Convert_CoverterError(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logs := mock_logger_.NewMockLogger(ctrl)
	logs.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
	 
	// Data
	original := &storage{
		l:          logs,
		img:        image_test_data_.PNG_16x16,
		storageKey: `TestConvertUnit/1.jpg`,
	}
	original_svg := false

	save := &storage{
		l:          logs,
		storageKey: `testThumb/1.ico`,
	}

	size := 16
	typ := types_.ICO()

	errNil := (error)(nil)
	errConverter := errors.New(`error converter`)
	errReturn := err_.Err(logs, convert_.LogFP, convert_.LogF03, errConverter)

	// Mock
	mock_converters := mock_converter_.NewMockConverterTyp(ctrl)
	mock_converters.EXPECT().Do(original, save, size, typ).Return(true, errConverter)

	mock_check_preview := mock_convert_.NewMockCheckPreview(ctrl)
	mock_check_preview.EXPECT().Check(typ, size).Return(errNil)

	mock_check_source := mock_convert_.NewMockCheckSource(ctrl)
	mock_check_source.EXPECT().Check(original, original_svg, size)

	// Test
	err := (&convert_.Converter{
		L: logs,
		Converters: []converter_.ConverterTyp{
			mock_converters,
		},
		CheckPreview: mock_check_preview,
		CheckSource:  mock_check_source,
	}).Do(
		original, save, original_svg, typ, size,
	)

	require.Equal(t, err, errReturn,
		fmt.Sprintf(`error: return '%v', want: '%v'`, err, errReturn))
}

func TestUnit_Convert_CoverterMulti(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logs := mock_logger_.NewMockLogger(ctrl)
	logs.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
	 

	// Data
	original := &storage{
		l:          logs,
		img:        image_test_data_.PNG_16x16,
		storageKey: `TestConvertUnit/1.jpg`,
	}
	original_svg := false

	save := &storage{
		l:          logs,
		storageKey: `testThumb/1.ico`,
	}

	size := 16
	typ := types_.ICO()

	errNil := (error)(nil)

	// Mock
	mock_converters_v1 := mock_converter_.NewMockConverterTyp(ctrl)
	mock_converters_v1.EXPECT().Do(original, save, size, typ).Return(true, errNil)

	mock_converters_v2 := mock_converter_.NewMockConverterTyp(ctrl)
	mock_converters_v2.EXPECT().Do(original, save, size, typ).Return(false, errNil).AnyTimes()

	mock_check_preview := mock_convert_.NewMockCheckPreview(ctrl)
	mock_check_preview.EXPECT().Check(typ, size).Return(errNil)

	mock_check_source := mock_convert_.NewMockCheckSource(ctrl)
	mock_check_source.EXPECT().Check(original, original_svg, size)

	// Test
	err := (&convert_.Converter{
		L: logs,
		Converters: []converter_.ConverterTyp{
			mock_converters_v1,
			mock_converters_v2,
		},
		CheckPreview: mock_check_preview,
		CheckSource:  mock_check_source,
	}).Do(
		original, save, original_svg, typ, size,
	)

	require.Equal(t, err, errNil,
		fmt.Sprintf(`error: return '%v', want: '%v'`, err, errNil))
}

func TestUnit_Convert_NoConverters(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logs := mock_logger_.NewMockLogger(ctrl)
	logs.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
	 

	// Data
	original := &storage{
		l:          logs,
		img:        image_test_data_.PNG_16x16,
		storageKey: `TestConvertUnit/1.jpg`,
	}
	original_svg := false

	save := &storage{
		l:          logs,
		storageKey: `testThumb/1.ico`,
	}

	size := 16
	typ := types_.ICO()

	errNil := (error)(nil)
	errReturn := err_.Err(logs, convert_.LogFP, convert_.LogF06)

	// Mock
	mock_check_preview := mock_convert_.NewMockCheckPreview(ctrl)
	mock_check_preview.EXPECT().Check(typ, size).Return(errNil)

	mock_check_source := mock_convert_.NewMockCheckSource(ctrl)
	mock_check_source.EXPECT().Check(original, original_svg, size)

	// Test
	err := (&convert_.Converter{
		L:            logs,
		Converters:   []converter_.ConverterTyp{},
		CheckPreview: mock_check_preview,
		CheckSource:  mock_check_source,
	}).Do(
		original, save, original_svg, typ, size,
	)

	require.Equal(t, err, errReturn,
		fmt.Sprintf(`error: return '%v', want: '%v'`, err, errReturn))
}

func TestUnit_Convert_SizePX0(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logs := mock_logger_.NewMockLogger(ctrl)
	logs.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
	 

	// Data
	original := &storage{
		l:          logs,
		img:        image_test_data_.PNG_16x16,
		storageKey: `TestConvertUnit/2.jpg`,
	}
	original_svg := false

	save := &storage{
		l:          logs,
		storageKey: `testThumb/2.ico`,
	}

	size := 0
	typ := types_.ICO()

	errReturn := err_.Err(logs, convert_.LogFP, convert_.LogF05)

	// Mock
	mock_converters := mock_converter_.NewMockConverterTyp(ctrl)
	mock_check_preview := mock_convert_.NewMockCheckPreview(ctrl)
	mock_check_source := mock_convert_.NewMockCheckSource(ctrl)

	// Test
	err := (&convert_.Converter{
		L: logs,
		Converters: []converter_.ConverterTyp{
			mock_converters,
		},
		CheckPreview: mock_check_preview,
		CheckSource:  mock_check_source,
	}).Do(
		original, save, original_svg, typ, size,
	)

	require.EqualError(t, err, errReturn.Error(),
		fmt.Sprintf(`error: return '%v', want: '%v'`, err, errReturn))
}

func TestUnit_Convert_PreviewCheckError(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logs := mock_logger_.NewMockLogger(ctrl)
	logs.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
	 

	// Data
	original := &storage{
		l:          logs,
		img:        image_test_data_.PNG_16x16,
		storageKey: `TestConvertUnit/1.jpg`,
	}
	original_svg := false

	save := &storage{
		l:          logs,
		storageKey: `testThumb/1.ico`,
	}

	size := 16
	typ := types_.ICO()

	errCheck := errors.New(`error check`)
	errReturn := err_.Err(logs, convert_.LogFP, convert_.LogF02, errCheck)

	// Mock
	mock_converters := mock_converter_.NewMockConverterTyp(ctrl)

	mock_check_preview := mock_convert_.NewMockCheckPreview(ctrl)
	mock_check_preview.EXPECT().Check(typ, size).Return(errCheck)

	mock_check_source := mock_convert_.NewMockCheckSource(ctrl)

	// Test
	err := (&convert_.Converter{
		L: logs,
		Converters: []converter_.ConverterTyp{
			mock_converters,
		},
		CheckPreview: mock_check_preview,
		CheckSource:  mock_check_source,
	}).Do(
		original, save, original_svg, typ, size,
	)

	require.EqualError(t, err, errReturn.Error(),
		fmt.Sprintf(`error: return '%v', want: '%v'`, err, errReturn))
}

func TestUnit_Convert_OriginalCheckError(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logs := mock_logger_.NewMockLogger(ctrl)
	logs.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()

	// Data
	original := &storage{
		l:          logs,
		img:        image_test_data_.PNG_16x16,
		storageKey: `TestConvertUnit/3.jpg`,
	}
	original_svg := false

	save := &storage{
		l:          logs,
		storageKey: `testThumb/3.ico`,
	}

	size := 16
	typ := types_.ICO()

	errNil := (error)(nil)
	errCheck := errors.New(`error check`)
	errReturn := err_.Err(logs, convert_.LogFP, convert_.LogF04, errCheck)

	// Mock
	mock_converters := mock_converter_.NewMockConverterTyp(ctrl)

	mock_check_preview := mock_convert_.NewMockCheckPreview(ctrl)
	mock_check_preview.EXPECT().Check(typ, size).Return(errNil)

	mock_check_source := mock_convert_.NewMockCheckSource(ctrl)
	mock_check_source.EXPECT().Check(original, original_svg, size).Return(errCheck)

	// Test
	err := (&convert_.Converter{
		L: logs,
		Converters: []converter_.ConverterTyp{
			mock_converters,
		},
		CheckPreview: mock_check_preview,
		CheckSource:  mock_check_source,
	}).Do(
		original, save, original_svg, typ, size,
	)

	require.EqualError(t, err, errReturn.Error(),
		fmt.Sprintf(`error: return '%v', want: '%v'`, err, errReturn))
}
