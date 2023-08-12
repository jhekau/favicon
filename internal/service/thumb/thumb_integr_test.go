package thumb_test

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 10 August 2023
 */
import (
	"io"
	"testing"

	logs_ "github.com/jhekau/favicon/internal/core/logs"
	logs_mock_ "github.com/jhekau/favicon/internal/core/logs/mock"
	mock_convert_ "github.com/jhekau/favicon/internal/mocks/intr/service/convert"
	mock_checks_ "github.com/jhekau/favicon/internal/mocks/intr/service/convert/checks"
	converter_ "github.com/jhekau/favicon/pkg/models/converter"
	storage_ "github.com/jhekau/favicon/pkg/models/storage"
	"github.com/stretchr/testify/require"

	mock_thumb_ "github.com/jhekau/favicon/internal/mocks/intr/service/thumb"
	mock_converter_ "github.com/jhekau/favicon/internal/mocks/pkg/models/converter"
	mock_storage_ "github.com/jhekau/favicon/internal/mocks/pkg/models/storage"
	convert_ "github.com/jhekau/favicon/internal/service/convert"
	checks_ "github.com/jhekau/favicon/internal/service/convert/checks"
	converters_ "github.com/jhekau/favicon/internal/service/convert/converters"
	thumb_ "github.com/jhekau/favicon/internal/service/thumb"
	types_ "github.com/jhekau/favicon/pkg/core/types"
	"go.uber.org/mock/gomock"
)

// conv Converter

/*



// Integration
func (t *Thumb) OriginalFileSet( filepath string ) {
	file := (&files_.Files{L: t.l}).NewObject(types_.FilePath(filepath))
	t.original = &original{
		obj: file,
	}
}
func (t *Thumb) OriginalFileSetSVG( filepath string ) {
	file := (&files_.Files{L: t.l}).NewObject(types_.FilePath(filepath))
	t.original = &original{
		typSVG: true,
		obj: file,
	}
}


*/

/*
go test ./internal/service/thumb/ -v -short -count=1 -race -coverprofile="coverage.out" -coverpkg='./internal/service/thumb,./internal/service/convert' -run="Test_Inegration_ConverterOnly" ;`
go tool cover -html="coverage.out" ;`
rm coverage.out
*/
func Test_Inegration_ConverterOnly( t *testing.T ) {

	// Data
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := &logs_.Logger{
		Typ: &logs_mock_.LoggerErrorf{},
	}

	cache := mock_thumb_.NewMockcache(ctrl)
	cache.EXPECT().Delete(nil).AnyTimes()
	cache.EXPECT().Load(gomock.Any()).Return(nil, false).AnyTimes()
	cache.EXPECT().Range( gomock.Any() ).AnyTimes()
	cache.EXPECT().Store(gomock.Any(), gomock.Any()).AnyTimes()

	thumbKey := `123`
	thumbSize := 16
	thumbTyp := types_.ICO()

	originalOBJ := mock_storage_.NewMockStorageOBJ(ctrl)

	thumbOBJ := mock_storage_.NewMockStorageOBJ(ctrl)
	thumbOBJ.EXPECT().IsExists().Return(false, nil)
	thumbOBJ.EXPECT().Reader().Return((io.ReadCloser)(nil), (error)(nil))

	storage := mock_storage_.NewMockStorage(ctrl)
	storage.EXPECT().NewObject(thumbKey).Return(thumbOBJ, nil)

	//
	converterTyp := mock_converter_.NewMockConverterTyp(ctrl)
	converterTyp.EXPECT().Do(originalOBJ, thumbOBJ, thumbSize, thumbTyp).Return(true, (error)(nil))	

	checkPreview := mock_convert_.NewMockCheckPreview(ctrl)
	checkPreview.EXPECT().Check(thumbTyp, thumbSize).Return((error)(nil))

	checkSource  := mock_convert_.NewMockCheckSource(ctrl)
	checkSource.EXPECT().Check(originalOBJ, false, thumbSize).Return((error)(nil))

	conv := &convert_.Converter{
		L            : logger,
		Converters   : []converter_.ConverterTyp{ 
			converterTyp,
		},
		CheckPreview : checkPreview,
		CheckSource  : checkSource,
	}

	//
	thumb, err := thumb_.NewThumb(thumbKey, thumb_.Typ(thumbTyp), logger, storage, conv)
	require.Equal(t, err, (error)(nil))
	require.NotNil(t, thumb)

	thumb.TestCacheSwap(cache).OriginalCustomSet(originalOBJ).SetSize(thumbSize)

	_, err = thumb.Read()
	require.Equal(t, err, (error)(nil))
}



/*
go test ./internal/service/thumb/ -v -short -count=1 -race -coverprofile="coverage.out" -coverpkg='./internal/service/thumb,./internal/service/convert,./internal/service/convert/checks' -run="Test_Inegration_ConverterCheckPreview" ;`
go tool cover -html="coverage.out" ;`
rm coverage.out
*/
func Test_Inegration_ConverterCheckPreview( t *testing.T ) {

	// Data
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := &logs_.Logger{
		Typ: &logs_mock_.LoggerErrorf{},
	}

	cache := mock_thumb_.NewMockcache(ctrl)
	cache.EXPECT().Delete(nil).AnyTimes()
	cache.EXPECT().Load(gomock.Any()).Return(nil, false).AnyTimes()
	cache.EXPECT().Range( gomock.Any() ).AnyTimes()
	cache.EXPECT().Store(gomock.Any(), gomock.Any()).AnyTimes()

	thumbKey := `123`
	thumbSize := 16
	thumbTyp := types_.ICO()

	originalOBJ := mock_storage_.NewMockStorageOBJ(ctrl)

	thumbOBJ := mock_storage_.NewMockStorageOBJ(ctrl)
	thumbOBJ.EXPECT().IsExists().Return(false, nil).AnyTimes()
	thumbOBJ.EXPECT().Reader().Return((io.ReadCloser)(nil), (error)(nil))

	storage := mock_storage_.NewMockStorage(ctrl)
	storage.EXPECT().NewObject(thumbKey).Return(thumbOBJ, nil)

	//
	converterTyp := mock_converter_.NewMockConverterTyp(ctrl)
	converterTyp.EXPECT().Do(originalOBJ, thumbOBJ, thumbSize, thumbTyp).Return(true, (error)(nil))	

	checkPreview := checks_.Preview{L: logger}

	checkSource  := mock_convert_.NewMockCheckSource(ctrl)
	checkSource.EXPECT().Check(originalOBJ, false, thumbSize).Return((error)(nil))

	conv := &convert_.Converter{
		L            : logger,
		Converters   : []converter_.ConverterTyp{ 
			converterTyp,
		},
		CheckPreview : checkPreview,
		CheckSource  : checkSource,
	}

	//
	thumb, err := thumb_.NewThumb(thumbKey, thumb_.Typ(thumbTyp), logger, storage, conv)
	require.Equal(t, err, (error)(nil))
	require.NotNil(t, thumb)

	thumb.TestCacheSwap(cache).OriginalCustomSet(originalOBJ).SetSize(thumbSize)

	_, err = thumb.Read()
	require.Equal(t, err, (error)(nil))

	_, err = thumb.SetSize(0).Read()
	require.NotEqual(t, err, (error)(nil))

	_, err = thumb.SetSize(1).Read()
	require.NotEqual(t, err, (error)(nil))
}



/*
go test ./internal/service/thumb/ -v -short -count=1 -race -coverprofile="coverage.out" -coverpkg='./internal/service/thumb,./internal/service/convert,./internal/service/convert/checks' -run="Test_Inegration_ConverterCheckSource" ;`
go tool cover -html="coverage.out" ;`
rm coverage.out
*/
func Test_Inegration_ConverterCheckSource( t *testing.T ) {

	// Data
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := &logs_.Logger{
		Typ: &logs_mock_.LoggerErrorf{},
	}

	cache := mock_thumb_.NewMockcache(ctrl)
	cache.EXPECT().Delete(nil).AnyTimes()
	cache.EXPECT().Load(gomock.Any()).Return(nil, false).AnyTimes()
	cache.EXPECT().Range( gomock.Any() ).AnyTimes()
	cache.EXPECT().Store(gomock.Any(), gomock.Any()).AnyTimes()

	thumbKey := `123`
	thumbSize := 16
	thumbTyp := types_.ICO()

	origKey := storage_.StorageKey(`432`)

	originalOBJ := mock_storage_.NewMockStorageOBJ(ctrl)
	originalOBJ.EXPECT().Key().Return(origKey).AnyTimes()
	originalOBJ.EXPECT().IsExists().Return(true, nil)

	thumbOBJ := mock_storage_.NewMockStorageOBJ(ctrl)
	thumbOBJ.EXPECT().IsExists().Return(false, nil)
	thumbOBJ.EXPECT().Reader().Return((io.ReadCloser)(nil), (error)(nil))

	storage := mock_storage_.NewMockStorage(ctrl)
	storage.EXPECT().NewObject(thumbKey).Return(thumbOBJ, nil)

	//
	converterTyp := mock_converter_.NewMockConverterTyp(ctrl)
	converterTyp.EXPECT().Do(originalOBJ, thumbOBJ, thumbSize, thumbTyp).Return(true, (error)(nil))	

	checkPreview := mock_convert_.NewMockCheckPreview(ctrl)
	checkPreview.EXPECT().Check(thumbTyp, thumbSize).Return((error)(nil))

	checkSource_MockCache := mock_checks_.NewMockCache(ctrl)
	checkSource_MockCache.EXPECT().Status(origKey, false, thumbSize).Return(false, nil)
	checkSource_MockCache.EXPECT().SetErr(origKey, false, thumbSize, (error)(nil))

	checkSource_MockResolution := mock_checks_.NewMockResolution(ctrl)
	checkSource_MockResolution.EXPECT().Get(originalOBJ).Return(thumbSize, thumbSize, (error)(nil))

	checkSource  := checks_.Source{
		L : logger,
		Cache : checkSource_MockCache,
		Resolution : checkSource_MockResolution,
	}

	conv := &convert_.Converter{
		L            : logger,
		Converters   : []converter_.ConverterTyp{ 
			converterTyp,
		},
		CheckPreview : checkPreview,
		CheckSource  : &checkSource,
	}

	//
	thumb, err := thumb_.NewThumb(thumbKey, thumb_.Typ(thumbTyp), logger, storage, conv)
	require.Equal(t, err, (error)(nil))
	require.NotNil(t, thumb)

	thumb.TestCacheSwap(cache).OriginalCustomSet(originalOBJ).SetSize(thumbSize)

	_, err = thumb.Read()
	require.Equal(t, err, (error)(nil))
}


/*
go test ./internal/service/thumb/ -v -short -count=1 -race -coverprofile="coverage.out" -coverpkg='./internal/service/thumb,./internal/service/convert,./internal/service/convert/converters' -run="Test_Inegration_ConverterConverters" ;`
go tool cover -html="coverage.out" ;`
rm coverage.out
*/
func Test_Inegration_ConverterConverters( t *testing.T ) {

	// Data
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := &logs_.Logger{
		Typ: &logs_mock_.LoggerErrorf{},
	}

	cache := mock_thumb_.NewMockcache(ctrl)
	cache.EXPECT().Delete(nil).AnyTimes()
	cache.EXPECT().Load(gomock.Any()).Return(nil, false).AnyTimes()
	cache.EXPECT().Range( gomock.Any() ).AnyTimes()
	cache.EXPECT().Store(gomock.Any(), gomock.Any()).AnyTimes()

	thumbKey := `123`
	thumbSize := 16
	thumbTyp := types_.ICO()

	origKey := storage_.StorageKey(`432`)

	originalOBJ := mock_storage_.NewMockStorageOBJ(ctrl)
	originalOBJ.EXPECT().Key().Return(origKey).AnyTimes()
	originalOBJ.EXPECT().IsExists().Return(true, nil)

	thumbOBJ := mock_storage_.NewMockStorageOBJ(ctrl)
	thumbOBJ.EXPECT().IsExists().Return(false, nil)
	thumbOBJ.EXPECT().Reader().Return((io.ReadCloser)(nil), (error)(nil))

	storage := mock_storage_.NewMockStorage(ctrl)
	storage.EXPECT().NewObject(thumbKey).Return(thumbOBJ, nil)

	converterExec := mock_converter_.NewMockConverterExec(ctrl)
	converterExec.EXPECT().Proc(originalOBJ, thumbOBJ, thumbSize, thumbTyp).Return((error)(nil))

	//
	checkPreview := mock_convert_.NewMockCheckPreview(ctrl)
	checkPreview.EXPECT().Check(thumbTyp, thumbSize).Return((error)(nil))

	checkSource_MockCache := mock_checks_.NewMockCache(ctrl)
	checkSource_MockCache.EXPECT().Status(origKey, false, thumbSize).Return(false, nil)
	checkSource_MockCache.EXPECT().SetErr(origKey, false, thumbSize, (error)(nil))

	checkSource_MockResolution := mock_checks_.NewMockResolution(ctrl)
	checkSource_MockResolution.EXPECT().Get(originalOBJ).Return(thumbSize, thumbSize, (error)(nil))

	checkSource  := checks_.Source{
		L : logger,
		Cache : checkSource_MockCache,
		Resolution : checkSource_MockResolution,
	}

	conv := &convert_.Converter{
		L            : logger,
		Converters   : []converter_.ConverterTyp{ 
			&converters_.ConverterPNG{L: logger, ConverterExec: converterExec},
			&converters_.ConverterICO{L: logger, ConverterExec: converterExec},
		},
		CheckPreview : checkPreview,
		CheckSource  : &checkSource,
	}

	//
	thumb, err := thumb_.NewThumb(thumbKey, thumb_.Typ(thumbTyp), logger, storage, conv)
	require.Equal(t, err, (error)(nil))
	require.NotNil(t, thumb)

	thumb.TestCacheSwap(cache).OriginalCustomSet(originalOBJ).SetSize(thumbSize)

	_, err = thumb.Read()
	require.Equal(t, err, (error)(nil))
}