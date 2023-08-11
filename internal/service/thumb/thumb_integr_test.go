package thumb_test

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 10 August 2023
 */
import (
	"testing"

	logger_ "github.com/jhekau/favicon/internal/core/logger"
	logger_mock_ "github.com/jhekau/favicon/internal/core/logger/mock"
	mock_convert_ "github.com/jhekau/favicon/internal/mocks/intr/service/convert"
	converter_ "github.com/jhekau/favicon/pkg/models/converter"
	"github.com/stretchr/testify/require"

	// mock_thumb_ "github.com/jhekau/favicon/internal/mocks/intr/service/thumb"
	mock_converter_ "github.com/jhekau/favicon/internal/mocks/pkg/models/converter"
	mock_storage_ "github.com/jhekau/favicon/internal/mocks/pkg/models/storage"
	convert_ "github.com/jhekau/favicon/internal/service/convert"
	thumb_ "github.com/jhekau/favicon/internal/service/thumb"
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



func Test_Inegration_Converter( t *testing.T ) {

	// Data
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := &logger_.Logger{
		Typ: &logger_mock_.LoggerErrorf{},
	}

	keyThumb := `123`

	storage := mock_storage_.NewMockStorage(ctrl)
	storage.EXPECT().NewObject(keyThumb)

	//
	converterTyp := mock_converter_.NewMockConverterTyp(ctrl)
	checkPreview := mock_convert_.NewMockCheckPreview(ctrl)
	checkSource  := mock_convert_.NewMockCheckSource(ctrl)

	conv := &convert_.Converter{
		L            : logger,
		Converters   : []converter_.ConverterTyp{ 
			converterTyp,
		},
		CheckPreview : checkPreview,
		CheckSource  : checkSource,
	}

	//
	_, err := thumb_.NewThumb(keyThumb, thumb_.ICO, logger, storage, conv)
	require.Equal(t, err, (error)(nil))

}
