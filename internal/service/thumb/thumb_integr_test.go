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
	mock_thumb_ "github.com/jhekau/favicon/internal/mocks/intr/service/thumb"
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

	storage := mock_thumb_.NewMockStorage(ctrl)
	storage.EXPECT().NewObject(keyThumb)

	//
	converterTyp := mock_convert_.NewMockConverterT(ctrl)
	checkPreview := mock_convert_.NewMockCheckPreview(ctrl)
	checkSource  := mock_convert_.NewMockCheckSource(ctrl)

	// conv := &convert_.Converter{
	// 	L            : logger,
	// 	Converters   : []convert_.ConverterT{ 
	// 		converterTyp,
	// 	},
	// 	CheckPreview : checkPreview,
	// 	CheckSource  : checkSource,
	// }

	var convT thumb_.Converter
	convT = &convert_.Converter{
		L            : logger,
		Converters   : []convert_.ConverterT{ 
			converterTyp,
		},
		CheckPreview : checkPreview,
		CheckSource  : checkSource,
	}

	//
	thumb, err := thumb_.NewThumb(keyThumb, thumb_.ICO, logger, storage, convT)

}
