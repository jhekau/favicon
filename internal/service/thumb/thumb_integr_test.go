package thumb_test

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 10 August 2023
 */
import (
	"io"
	"testing"

	logger_ "github.com/jhekau/favicon/internal/core/logger"
	logger_mock_ "github.com/jhekau/favicon/internal/core/logger/mock"
	mock_convert_ "github.com/jhekau/favicon/internal/mocks/intr/service/convert"
	converter_ "github.com/jhekau/favicon/pkg/models/converter"
	"github.com/stretchr/testify/require"

	mock_thumb_ "github.com/jhekau/favicon/internal/mocks/intr/service/thumb"
	mock_converter_ "github.com/jhekau/favicon/internal/mocks/pkg/models/converter"
	mock_storage_ "github.com/jhekau/favicon/internal/mocks/pkg/models/storage"
	convert_ "github.com/jhekau/favicon/internal/service/convert"
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
rm coverage.out
go test ./internal/service/thumb/ -v -short -count=1 -race -coverprofile="coverage.out" -coverpkg=./... -run="Test_Inegration_Converter"
go tool cover -html="coverage.out"
*/
func Test_Inegration_Converter( t *testing.T ) {

	// Data
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := &logger_.Logger{
		Typ: &logger_mock_.LoggerErrorf{},
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
