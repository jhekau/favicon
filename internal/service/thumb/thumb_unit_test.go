package thumb_test

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 10 August 2023
 */
import (
	"fmt"
	"testing"

	logger_ "github.com/jhekau/favicon/internal/core/logger"
	logger_mock_ "github.com/jhekau/favicon/internal/core/logger/mock"
	mock_thumb_ "github.com/jhekau/favicon/internal/mocks/intr/service/thumb"
	thumb_ "github.com/jhekau/favicon/internal/service/thumb"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

/*

type Thumb struct {

	s sync.RWMutex
	l *logger_.Logger

	original *original
	thumb StorageOBJ

	storage Storage
	conv Converter

	size_px uint16
	size_attr_value attr_size
	comment string // <!-- comment -->
	url_href types_.URLHref // domain{/name_url}, first -> `/`
	url_href_clear types_.URLHref
	tag_rel string
	manifest bool
	mimetype types_.FileType
	typ types_.FileType
	cache cache

}

func (*Thumb).GetHREF() types_.URLHref
func (*Thumb).GetHREFClear() types_.URLHref
func (*Thumb).GetOriginalKey() string
func (*Thumb).GetSize() uint16
func (*Thumb).GetTAG() string
func (*Thumb).GetType() types_.FileType
func (*Thumb).OriginalCustomSet(obj StorageOBJ)
func (*Thumb).OriginalCustomSetSVG(obj StorageOBJ)
func (*Thumb).OriginalFileSet(filepath string)
func (*Thumb).OriginalFileSetSVG(filepath string)
func (*Thumb).Read() (io.ReadCloser, error)
func (*Thumb).SetHREF(src string) *Thumb
func (*Thumb).SetHTMLComment(comment string) *Thumb
func (*Thumb).SetManifestUsed() *Thumb
func (*Thumb).SetSize(px uint16) *Thumb
func (*Thumb).SetSizeAttrCustom(val string) *Thumb
func (*Thumb).SetSizeAttrDefault() *Thumb
func (*Thumb).SetSizeAttrEmpty() *Thumb
func (*Thumb).SetTagRel(tagRel string) *Thumb
func (*Thumb).SetType(typ types_.FileType) *Thumb
func (*Thumb).StatusManifest() bool




func (*Thumb).get_href() types_.URLHref
func (*Thumb).get_href_clear() types_.URLHref
func (*Thumb).get_original_key() string
func (*Thumb).get_size() uint16
func (*Thumb).get_tag() string
func (*Thumb).get_type() types_.FileType
func (*Thumb).original_get(filepath string) *original
func (*Thumb).read() (io.ReadCloser, error)
func (*Thumb).set_href(src string) *Thumb
func (*Thumb).set_html_comment(comment string) *Thumb
func (*Thumb).set_manifest_used() *Thumb
func (*Thumb).set_size(px uint16) *Thumb
func (*Thumb).set_size_attr_custom(val string) *Thumb
func (*Thumb).set_size_attr_default() *Thumb
func (*Thumb).set_size_attr_empty() *Thumb
func (*Thumb).set_tag_rel(tagRel string) *Thumb
func (*Thumb).set_type(typ types_.FileType) *Thumb
func (*Thumb).status_manifest() bool
func (*Thumb).thumb_create() error

*/

func Test_NewThumb( t *testing.T ) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := &logger_.Logger{
		Typ: &logger_mock_.LoggerErrorf{},
	}

	storage := mock_thumb_.NewMockStorage(ctrl)
	storage.EXPECT().NewObject().AnyTimes()

	conv := mock_thumb_.NewMockConverter(ctrl)
	conv.EXPECT().Do(nil, nil, false, nil, 0).AnyTimes()

	thumb := thumb_.NewThumb(logger, storage, conv)

	require.IsType(t, thumb, &thumb_.Thumb{}, fmt.Sprintf(`%T`, thumb))
}

func Test_Size( t *testing.T ) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := &logger_.Logger{
		Typ: &logger_mock_.LoggerErrorf{},
	}

	storage := mock_thumb_.NewMockStorage(ctrl)
	storage.EXPECT().NewObject().AnyTimes()

	conv := mock_thumb_.NewMockConverter(ctrl)
	conv.EXPECT().Do(nil, nil, false, nil, 0).AnyTimes()

	cache := mock_thumb_.NewMockcache(ctrl)
	cache.EXPECT().Delete(nil).AnyTimes()
	cache.EXPECT().Load(nil).AnyTimes()
	cache.EXPECT().Range( gomock.Any() )
	cache.EXPECT().Store(nil, nil).AnyTimes()

	size := 16
	thumb := thumb_.NewThumb(logger, storage, conv)

	thumb.TestCacheSwap(cache)
	thumb.SetSize(size)

	res := thumb.GetSize()

	require.Equal(t, res, size)
}

func Test_SetTagRel( t *testing.T ) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := &logger_.Logger{
		Typ: &logger_mock_.LoggerErrorf{},
	}

	storage := mock_thumb_.NewMockStorage(ctrl)
	storage.EXPECT().NewObject().AnyTimes()

	conv := mock_thumb_.NewMockConverter(ctrl)
	conv.EXPECT().Do(nil, nil, false, nil, 0).AnyTimes()

	cache := mock_thumb_.NewMockcache(ctrl)
	cache.EXPECT().Delete(nil).AnyTimes()
	cache.EXPECT().Load(gomock.Any()).Return(nil, false)
	cache.EXPECT().Range( gomock.Any() )
	cache.EXPECT().Store(gomock.Any(), gomock.Any()).AnyTimes()

	thumb := thumb_.NewThumb(logger, storage, conv)
	thumb.TestCacheSwap(cache)

	tagRel := `apple-touch-icon`
	expect := `<link rel="apple-touch-icon" >`

	thumb.SetTagRel(tagRel)

	res := thumb.GetTAG()

	require.Equal(t, expect, res)
}


/*
func (t *Thumb) SetManifestUsed() *Thumb {
	return t.set_manifest_used()
}

func (t *Thumb) SetHTMLComment(comment string) *Thumb {
	return t.set_html_comment(comment)
}

func (t *Thumb) SetType(typ types_.FileType) *Thumb {
	return t.set_type(typ)
}

func (t *Thumb) GetType() types_.FileType {
	return t.get_type()
}

func (t *Thumb) SetHREF(src string) *Thumb {
	return t.set_href(src)
}

func (t *Thumb) GetHREF() types_.URLHref {
	return t.get_href()
}

func (t *Thumb) GetHREFClear() types_.URLHref {
	return t.get_href_clear()
}

func (t *Thumb) StatusManifest() bool { // ( string, bool /*true - used )
	return t.status_manifest()
}

func (t *Thumb) GetTAG() string {
	return t.tagGet()
}

func (t *Thumb) SetSizeAttrEmpty() *Thumb {
	return t.set_size_attr_empty()
}

func (t *Thumb) SetSizeAttrDefault() *Thumb {
	return t.set_size_attr_default()
}

func (t *Thumb) SetSizeAttrCustom(val string) *Thumb {
	return t.set_size_attr_custom(val)
}

func (t *Thumb) GetOriginalKey() string{
	return t.get_original_key()
}

func (t *Thumb) Read() (io.ReadCloser, error) {
	return t.read()
}



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
func (t *Thumb) OriginalCustomSet( obj StorageOBJ ) {
	t.original = &original{
		obj: obj,
	}
}
func (t *Thumb) OriginalCustomSetSVG( obj StorageOBJ ) {
	t.original = &original{
		typSVG: true,
		obj: obj,
	}
}

*/