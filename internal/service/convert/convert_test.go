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

	logs_ "github.com/jhekau/favicon/internal/core/logs"
	logs_mock_ "github.com/jhekau/favicon/internal/core/logs/mock"
	image_test_data_ "github.com/jhekau/favicon/internal/core/test_data/image"
	types_ "github.com/jhekau/favicon/internal/core/types"
	mock_convert_ "github.com/jhekau/favicon/internal/mocks/intr/service/convert"
	convert_ "github.com/jhekau/favicon/internal/service/convert"
	// checks_ "github.com/jhekau/favicon/internal/service/convert/checks"
	// converters_ "github.com/jhekau/favicon/internal/service/convert/converters"
	domain_ "github.com/jhekau/favicon/pkg/domain"

	// domain_ "github.com/jhekau/favicon/pkg/domain"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)



// image test data
type storage struct{
	l *logs_.Logger
	img interface{ 
		Base64Reader(l *logs_.Logger) (io.Reader, string, error)
	}
	ifExist bool
	ifExistError error
	storageKey domain_.StorageKey
}
func (s *storage) IsExists() ( bool, error ) {
	return s.ifExist, s.ifExistError
}
func (s *storage) Key() domain_.StorageKey {
	return s.storageKey
}
func (s *storage) Reader() (io.ReadCloser , error) {
	r, err := image_test_data_.GetFileReader(s.img, s.l)
	return io.NopCloser(r), err
}
func (s *storage) Writer() (io.WriteCloser, error) {
	var w io.WriteCloser
	return w, nil
}




// Unit *******

//
func TestUnit_Convert_JPGxICO( t *testing.T ) {

	// Data
	logger := &logs_.Logger{
		Typ: &logs_mock_.LoggerErrorf{},
	}

	original := &storage{
		l: logger,
		img: image_test_data_.PNG_16x16,
		storageKey: `TestConvertUnit/1.jpg`,
	}
	original_svg := false

	save := &storage{
		l: logger,
		storageKey: `testThumb/1.ico`,
	}

	size := 16
	typ := types_.ICO()

	errNil := (error)(nil)

	// Mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_converters := mock_convert_.NewMockConverterT(ctrl)
	mock_converters.EXPECT().Do(original, save, size, typ).Return(true, errNil)

	mock_check_preview := mock_convert_.NewMockCheckPreview(ctrl)
	mock_check_preview.EXPECT().Check(typ, size)

	mock_check_source := mock_convert_.NewMockCheckSource(ctrl)
	mock_check_source.EXPECT().Check(original, original_svg, size)
	
	// Test
	err := (&convert_.Converter{
		L: logger,
		Converters: []convert_.ConverterT{
			mock_converters,
		},
		CheckPreview: mock_check_preview,
		CheckSource: mock_check_source,

	}).Do( 
		original, save, original_svg, typ, size,
	)

	require.Equal(t, err, errNil, 
		fmt.Sprintf(`error: return '%v', want: '%v'`, err, errNil ))
}

//
func TestUnit_Convert_CoverterError( t *testing.T ) {

	// Data
	logger := &logs_.Logger{
		Typ: &logs_mock_.LoggerErrorf{},
	}

	original := &storage{
		l: logger,
		img: image_test_data_.PNG_16x16,
		storageKey: `TestConvertUnit/1.jpg`,
	}
	original_svg := false

	save := &storage{
		l: logger,
		storageKey: `testThumb/1.ico`,
	}
	
	size := 16
	typ := types_.ICO()

	errNil := (error)(nil)
	errConverter := errors.New(`error converter`)
	errReturn := logger.Typ.Error(convert_.LogFP, convert_.LogF03, errConverter)

	// Mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_converters := mock_convert_.NewMockConverterT(ctrl)
	mock_converters.EXPECT().Do(original, save, size, typ).Return(true, errConverter)

	mock_check_preview := mock_convert_.NewMockCheckPreview(ctrl)
	mock_check_preview.EXPECT().Check(typ, size).Return(errNil)

	mock_check_source := mock_convert_.NewMockCheckSource(ctrl)
	mock_check_source.EXPECT().Check(original, original_svg, size)
	
	// Test
	err := (&convert_.Converter{
		L: logger,
		Converters: []convert_.ConverterT{
			mock_converters,
		},
		CheckPreview: mock_check_preview,
		CheckSource: mock_check_source,

	}).Do( 
		original, save, original_svg, typ, size,
	)

	require.Equal(t, err, errReturn, 
		fmt.Sprintf(`error: return '%v', want: '%v'`, err, errReturn ))
}

//
func TestUnit_Convert_CoverterMulti( t *testing.T ) {

	// Data
	logger := &logs_.Logger{
		Typ: &logs_mock_.LoggerErrorf{},
	}

	original := &storage{
		l: logger,
		img: image_test_data_.PNG_16x16,
		storageKey: `TestConvertUnit/1.jpg`,
	}
	original_svg := false

	save := &storage{
		l: logger,
		storageKey: `testThumb/1.ico`,
	}

	size := 16
	typ := types_.ICO()


	errNil := (error)(nil)

	// Mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_converters_v1 := mock_convert_.NewMockConverterT(ctrl)
	mock_converters_v1.EXPECT().Do(original, save, size, typ).Return(true, errNil)

	mock_converters_v2 := mock_convert_.NewMockConverterT(ctrl)
	mock_converters_v2.EXPECT().Do(original, save, size, typ).Return(false, errNil).AnyTimes()

	mock_check_preview := mock_convert_.NewMockCheckPreview(ctrl)
	mock_check_preview.EXPECT().Check(typ, size).Return(errNil)

	mock_check_source := mock_convert_.NewMockCheckSource(ctrl)
	mock_check_source.EXPECT().Check(original, original_svg, size)
	
	// Test
	err := (&convert_.Converter{
		L: logger,
		Converters: []convert_.ConverterT{
			mock_converters_v1,
			mock_converters_v2, 
		},
		CheckPreview: mock_check_preview,
		CheckSource: mock_check_source,

	}).Do( 
		original, save, original_svg, typ, size,
	)

	require.Equal(t, err, errNil, 
		fmt.Sprintf(`error: return '%v', want: '%v'`, err, errNil ))
}

//
func TestUnit_Convert_NoConverters( t *testing.T ) {

	// Data
	logger := &logs_.Logger{
		Typ: &logs_mock_.LoggerErrorf{},
	}

	original := &storage{
		l: logger,
		img: image_test_data_.PNG_16x16,
		storageKey: `TestConvertUnit/1.jpg`,
	}
	original_svg := false

	save := &storage{
		l: logger,
		storageKey: `testThumb/1.ico`,
	}

	size := 16
	typ := types_.ICO()

	errNil := (error)(nil)
	errReturn := logger.Typ.Error(convert_.LogFP, convert_.LogF06)

	// Mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_check_preview := mock_convert_.NewMockCheckPreview(ctrl)
	mock_check_preview.EXPECT().Check(typ, size).Return(errNil)

	mock_check_source := mock_convert_.NewMockCheckSource(ctrl)
	mock_check_source.EXPECT().Check(original, original_svg, size)
	
	// Test
	err := (&convert_.Converter{
		L: logger,
		Converters: []convert_.ConverterT{},
		CheckPreview: mock_check_preview,
		CheckSource: mock_check_source,

	}).Do( 
		original, save, original_svg, typ, size,
	)

	require.Equal(t, err, errReturn, 
		fmt.Sprintf(`error: return '%v', want: '%v'`, err, errReturn ))
}

// 
func TestUnit_Convert_SizePX0( t *testing.T ) {

	// Data
	logger := &logs_.Logger{
		Typ: &logs_mock_.LoggerErrorf{},
	}

	original := &storage{
		l: logger,
		img: image_test_data_.PNG_16x16,
		storageKey: `TestConvertUnit/2.jpg`,
	}
	original_svg := false

	save := &storage{
		l: logger,
		storageKey: `testThumb/2.ico`,
	}

	size := 0
	typ := types_.ICO()


	errReturn := logger.Typ.Error(convert_.LogFP, convert_.LogF05)

	// Mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_converters := mock_convert_.NewMockConverterT(ctrl)
	mock_check_preview := mock_convert_.NewMockCheckPreview(ctrl)
	mock_check_source := mock_convert_.NewMockCheckSource(ctrl)
	
	// Test
	err := (&convert_.Converter{
		L: logger,
		Converters: []convert_.ConverterT{
			mock_converters,
		},
		CheckPreview: mock_check_preview,
		CheckSource: mock_check_source,

	}).Do( 
		original, save, original_svg, typ, size,
	)

	require.EqualError(t, err, errReturn.Error(), 
		fmt.Sprintf(`error: return '%v', want: '%v'`, err, errReturn ))
}


func TestUnit_Convert_PreviewCheckError( t *testing.T ) {

	// Data
	logger := &logs_.Logger{
		Typ: &logs_mock_.LoggerErrorf{},
	}

	original := &storage{
		l: logger,
		img: image_test_data_.PNG_16x16,
		storageKey: `TestConvertUnit/1.jpg`,
	}
	original_svg := false

	save := &storage{
		l: logger,
		storageKey: `testThumb/1.ico`,
	}

	size := 16
	typ := types_.ICO()

	// errNil := (error)(nil)
	errCheck := errors.New(`error check`)
	errReturn := logger.Typ.Error(convert_.LogFP, convert_.LogF02, errCheck)

	// Mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_converters := mock_convert_.NewMockConverterT(ctrl)

	mock_check_preview := mock_convert_.NewMockCheckPreview(ctrl)
	mock_check_preview.EXPECT().Check(typ, size).Return(errCheck)

	mock_check_source := mock_convert_.NewMockCheckSource(ctrl)
	
	// Test
	err := (&convert_.Converter{
		L: logger,
		Converters: []convert_.ConverterT{
			mock_converters,
		},
		CheckPreview: mock_check_preview,
		CheckSource: mock_check_source,

	}).Do( 
		original, save, original_svg, typ, size,
	)

	require.EqualError(t, err, errReturn.Error(), 
		fmt.Sprintf(`error: return '%v', want: '%v'`, err, errReturn ))
}


//
func TestUnit_Convert_OriginalCheckError( t *testing.T ) {

	// Data
	logger := &logs_.Logger{
		Typ: &logs_mock_.LoggerErrorf{},
	}

	original := &storage{
		l: logger,
		img: image_test_data_.PNG_16x16,
		storageKey: `TestConvertUnit/3.jpg`,
	}
	original_svg := false

	save := &storage{
		l: logger,
		storageKey: `testThumb/3.ico`,
	}

	size := 16
	typ := types_.ICO()

	errNil := (error)(nil)
	errCheck := errors.New(`error check`)
	errReturn := logger.Typ.Error(convert_.LogFP, convert_.LogF04, errCheck)

	// Mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock_converters := mock_convert_.NewMockConverterT(ctrl)

	mock_check_preview := mock_convert_.NewMockCheckPreview(ctrl)
	mock_check_preview.EXPECT().Check(typ, size).Return(errNil)

	mock_check_source := mock_convert_.NewMockCheckSource(ctrl)
	mock_check_source.EXPECT().Check(original, original_svg, size).Return(errCheck)
	
	// Test
	err := (&convert_.Converter{
		L: logger,
		Converters: []convert_.ConverterT{
			mock_converters,
		},
		CheckPreview: mock_check_preview,
		CheckSource: mock_check_source,

	}).Do( 
		original, save, original_svg, typ, size,
	)

	require.EqualError(t, err, errReturn.Error(), 
		fmt.Sprintf(`error: return '%v', want: '%v'`, err, errReturn ))
}













/// *** Integation ****

// CheckPreview
// CheckSource
// ConvertersPNG
// ConvertersISO
// 
// Original_16, Preview_16, ICO, ConverterICO [PASS]
// Original_1, Previce_16, ICO, ConvertICO [FALSE]
// Original_16, Previce_16, PNG, ConverterPNG [PASS]
// Original_0, Previre_16. ICO, nil [FALSE]
// Original_16, Preview_16, ICO, nil [FALSE]

// !!!! Переделать


// Original_16, Preview_16, ICO, ConverterICO [PASS]
func TestIntegrationConvert( t *testing.T ) {

	// Data 
	logger := &logs_.Logger{
		Typ: &logs_mock_.LoggerErrorf{},
	}
	
	_ = logger



	/*
	// Test
	err := (&convert_.Converter{
		L: logger,
		Converters: []convert_.ConverterT{
			&converters_.ConverterICO{
				L: logger,
			}, 
			&converters_.ConverterPNG{
				L: logger,
			},
		},
		CheckPreview: checks_.Preview{
			L: logger,
		},
		CheckSource: checks_.Source{
			L: logger,
		},

	}).Do( 
		original, save, original_svg, types_.ICO(), size,
	)
	*/


	/*
	// Data
	original := types_.FilePath(`TestConvertUnit/1.jpg`)
	original_svg := false
	save := types_.FilePath(`1.ico`)
	size := 16

	logger := &logs_.Logger{
		Typ: &logs_mock_.LoggerErrorf{},
	}

	// Mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()


	// Test
	err := (&convert_.Converter{
		L: logger,
		Converters: []convert_.ConverterT{
			&convertType{types_.ICO(), nil, converter{nil}},
			&convertType{types_.SVG(), nil, converter{nil}},
		},
		CheckPreview: checks_.Preview{
			L: logger,
		},
		CheckSource: mock_check_source,

	}).Do( 
		original, save, original_svg, types_.ICO(), size,
	)


	require.Equal(t, err, nil, fmt.Sprintf(`error: '%v'`, err ))
	*/
}

