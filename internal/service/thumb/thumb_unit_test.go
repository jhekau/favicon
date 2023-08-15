package thumb_test

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 10 August 2023
 */
import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/url"
	"sync"
	"testing"

	// logger_ "github.com/jhekau/favicon/pkg/core/models/logger"
	logs_mock_ "github.com/jhekau/favicon/internal/core/logs/mock"

	typ_ "github.com/jhekau/favicon/internal/core/types"
	mock_converter_ "github.com/jhekau/favicon/internal/mocks/pkg/models/converter"
	mock_storage_ "github.com/jhekau/favicon/internal/mocks/pkg/models/storage"
	types_ "github.com/jhekau/favicon/pkg/core/types"
	storage_ "github.com/jhekau/favicon/pkg/core/models/storage"

	mock_thumb_ "github.com/jhekau/favicon/internal/mocks/intr/service/thumb"
	thumb_ "github.com/jhekau/favicon/internal/service/thumb"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)



func Test_NewThumb( t *testing.T ) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := &logs_mock_.LoggerErrorf{}

	keyThumb := `123`

	storage := mock_storage_.NewMockStorage(ctrl)
	storage.EXPECT().NewObject(keyThumb)

	conv := mock_converter_.NewMockConverter(ctrl)

	//
	thumb, _ := thumb_.NewThumb(keyThumb, thumb_.ICO, logger, storage, conv)

	require.IsType(t, thumb, &thumb_.Thumb{}, fmt.Sprintf(`%T`, thumb))
}

func Test_NewThumbError( t *testing.T ) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := &logs_mock_.LoggerErrorf{}

	keyThumb := `123`

	instanceErr := errors.New(`error new object`)
	storage :=mock_storage_.NewMockStorage(ctrl)
	storage.EXPECT().NewObject(keyThumb).Return(nil, instanceErr)

	conv := mock_converter_.NewMockConverter(ctrl)
	conv.EXPECT().Do(nil, nil, false, nil, 0).AnyTimes()

	_, err := thumb_.NewThumb(keyThumb, thumb_.ICO, logger, storage, conv)

	require.Equal(t, err, logger.Error(thumb_.LogTP, thumb_.LogT02, instanceErr))
}

func Test_Size( t *testing.T ) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := &logs_mock_.LoggerErrorf{}

	keyThumb := `123`

	storage :=mock_storage_.NewMockStorage(ctrl)
	storage.EXPECT().NewObject(keyThumb).AnyTimes()

	conv := mock_converter_.NewMockConverter(ctrl)
	conv.EXPECT().Do(nil, nil, false, nil, 0).AnyTimes()

	cache := mock_thumb_.NewMockcache(ctrl)
	cache.EXPECT().Delete(nil).AnyTimes()
	cache.EXPECT().Load(nil).AnyTimes()
	cache.EXPECT().Range( gomock.Any() )
	cache.EXPECT().Store(nil, nil).AnyTimes()

	size := 16
	thumb, _ := thumb_.NewThumb(keyThumb, thumb_.ICO, logger, storage, conv)

	thumb.TestCacheSwap(cache)
	thumb.SetSize(size)

	expect := thumb.GetSize()
	require.Equal(t, expect, size)
}

func Test_SetTagRel( t *testing.T ) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := &logs_mock_.LoggerErrorf{}

	keyThumb := `123`

	storage :=mock_storage_.NewMockStorage(ctrl)
	storage.EXPECT().NewObject(keyThumb).AnyTimes()

	conv := mock_converter_.NewMockConverter(ctrl)
	conv.EXPECT().Do(nil, nil, false, nil, 0).AnyTimes()

	cache := mock_thumb_.NewMockcache(ctrl)
	cache.EXPECT().Delete(nil).AnyTimes()
	cache.EXPECT().Load(gomock.Any()).Return(nil, false)
	cache.EXPECT().Range( gomock.Any() )
	cache.EXPECT().Store(gomock.Any(), gomock.Any()).AnyTimes()

	thumb, _ := thumb_.NewThumb(keyThumb, thumb_.TestTypEmpty, logger, storage, conv)
	thumb.TestCacheSwap(cache)

	tagRel := `apple-touch-icon`
	expect := `<link rel="apple-touch-icon" >`

	thumb.SetAttrRel(tagRel)

	tag := thumb.GetTAG()
	require.Equal(t, expect, tag)
}

func Test_HTMLComment( t *testing.T ) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := &logs_mock_.LoggerErrorf{}

	keyThumb := `123`

	storage :=mock_storage_.NewMockStorage(ctrl)
	storage.EXPECT().NewObject(keyThumb).AnyTimes()

	conv := mock_converter_.NewMockConverter(ctrl)
	conv.EXPECT().Do(nil, nil, false, nil, 0).AnyTimes()

	cache := mock_thumb_.NewMockcache(ctrl)
	cache.EXPECT().Delete(nil).AnyTimes()
	cache.EXPECT().Load(gomock.Any()).Return(nil, false)
	cache.EXPECT().Range( gomock.Any() ).AnyTimes()
	cache.EXPECT().Store(gomock.Any(), gomock.Any()).AnyTimes()

	thumb, _ := thumb_.NewThumb(keyThumb, thumb_.TestTypEmpty, logger, storage, conv)
	thumb.TestCacheSwap(cache)

	tagRel := `apple-touch-icon`
	htmlComment := `hello`
	expect := `<link rel="`+tagRel+`" ><!-- `+htmlComment+` -->`

	thumb.
		SetAttrRel(tagRel).
		SetHTMLComment(htmlComment)

	tag := thumb.GetTAG()
	require.Equal(t, expect, tag)
}

func Test_HTMLCommentEmptyTag( t *testing.T ) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := &logs_mock_.LoggerErrorf{}

	keyThumb := `123`

	storage :=mock_storage_.NewMockStorage(ctrl)
	storage.EXPECT().NewObject(keyThumb).AnyTimes()

	conv := mock_converter_.NewMockConverter(ctrl)
	conv.EXPECT().Do(nil, nil, false, nil, 0).AnyTimes()

	cache := mock_thumb_.NewMockcache(ctrl)
	cache.EXPECT().Delete(nil).AnyTimes()
	cache.EXPECT().Load(gomock.Any()).Return(nil, false)
	cache.EXPECT().Range( gomock.Any() ).AnyTimes()
	cache.EXPECT().Store(gomock.Any(), gomock.Any()).AnyTimes()

	thumb, _ := thumb_.NewThumb(keyThumb, thumb_.TestTypEmpty, logger, storage, conv)
	thumb.TestCacheSwap(cache)

	htmlComment := `hello`
	expect := ``

	thumb.SetHTMLComment(htmlComment)

	tag := thumb.GetTAG()
	require.Equal(t, expect, tag)
}

func Test_ManifestUsed( t *testing.T ) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := &logs_mock_.LoggerErrorf{}

	keyThumb := `123`

	storage :=mock_storage_.NewMockStorage(ctrl)
	storage.EXPECT().NewObject(keyThumb).AnyTimes()

	conv := mock_converter_.NewMockConverter(ctrl)
	conv.EXPECT().Do(nil, nil, false, nil, 0).AnyTimes()

	cache := mock_thumb_.NewMockcache(ctrl)
	cache.EXPECT().Delete(nil).AnyTimes()
	cache.EXPECT().Load(gomock.Any()).Return(nil, false).AnyTimes()
	cache.EXPECT().Range( gomock.Any() )
	cache.EXPECT().Store(gomock.Any(), gomock.Any()).AnyTimes()

	thumb, _ := thumb_.NewThumb(keyThumb, thumb_.ICO, logger, storage, conv)
	thumb.TestCacheSwap(cache)

	thumb.SetManifestUsed()
	expect := thumb.StatusManifest()

	require.Equal(t, expect, true)
}

func Test_TypeThumb( t *testing.T ) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := &logs_mock_.LoggerErrorf{}

	keyThumb := `123`

	storage :=mock_storage_.NewMockStorage(ctrl)
	storage.EXPECT().NewObject(keyThumb).AnyTimes()

	conv := mock_converter_.NewMockConverter(ctrl)
	conv.EXPECT().Do(nil, nil, false, nil, 0).AnyTimes()

	cache := mock_thumb_.NewMockcache(ctrl)
	cache.EXPECT().Delete(nil).AnyTimes()
	cache.EXPECT().Load(gomock.Any()).Return(nil, false)
	cache.EXPECT().Store(gomock.Any(), gomock.Any()).AnyTimes()

	typ := types_.PNG()

	thumb, _ := thumb_.NewThumb(keyThumb, thumb_.Typ(typ), logger, storage, conv)
	thumb.TestCacheSwap(cache)

	expect := `<link type="`+string(typ)+`" >`

	tag := thumb.GetTAG()
	require.Equal(t, expect, tag)

	mimetype := thumb.GetType()
	require.Equal(t, mimetype, typ)

}

func Test_Href( t *testing.T ) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := &logs_mock_.LoggerErrorf{}

	keyThumb := `123`

	storage :=mock_storage_.NewMockStorage(ctrl)
	storage.EXPECT().NewObject(keyThumb).AnyTimes()

	conv := mock_converter_.NewMockConverter(ctrl)
	conv.EXPECT().Do(nil, nil, false, nil, 0).AnyTimes()

	cache := mock_thumb_.NewMockcache(ctrl)
	cache.EXPECT().Delete(nil).AnyTimes()
	cache.EXPECT().Load(gomock.Any()).Return(nil, false)
	cache.EXPECT().Range( gomock.Any() )
	cache.EXPECT().Store(gomock.Any(), gomock.Any()).AnyTimes()

	thumb,  _ := thumb_.NewThumb(keyThumb, thumb_.TestTypEmpty, logger, storage, conv)
	thumb.TestCacheSwap(cache)

	href := `/path/thumbs/image.png`
	expect := `<link href="`+string(href)+`" >`

	thumb.SetURLPath(href)

	tag := thumb.GetTAG()
	require.Equal(t, expect, tag)

	h := thumb.GetURLPath()
	require.Equal(t, href, string(h))
}

func Test_SizeAttr( t *testing.T ) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := &logs_mock_.LoggerErrorf{}

	keyThumb := `123`

	storage :=mock_storage_.NewMockStorage(ctrl)
	storage.EXPECT().NewObject(keyThumb).AnyTimes()

	conv := mock_converter_.NewMockConverter(ctrl)
	conv.EXPECT().Do(nil, nil, false, nil, 0).AnyTimes()

	cache := mock_thumb_.NewMockcache(ctrl)
	cache.EXPECT().Delete(nil).AnyTimes()
	cache.EXPECT().Load(gomock.Any()).Return(nil, false).AnyTimes()
	cache.EXPECT().Range( gomock.Any() ).AnyTimes()
	cache.EXPECT().Store(gomock.Any(), gomock.Any()).AnyTimes()

	size := 16

	//
	thumb, _ := thumb_.NewThumb(keyThumb, thumb_.TestTypEmpty, logger, storage, conv)
	thumb.
		TestCacheSwap(cache).
		SetSize(size).
		SetAttrRel(`apple-touch-icon`)

	thumb.SetAttrSize_Empty()
	tag := thumb.GetTAG()
	require.Equal(t, tag, `<link rel="apple-touch-icon" >`)

	//
	thumb, _ = thumb_.NewThumb(keyThumb, thumb_.TestTypEmpty, logger, storage, conv)
	thumb.
		TestCacheSwap(cache).
		SetSize(size)

	thumb.SetAttrSize_Default()
	tag = thumb.GetTAG()
	require.Equal(t, tag, fmt.Sprintf(`<link sizes="%vx%v" >`, size, size), 
		fmt.Sprintf(`size: '%v'`, thumb.GetSize()))

	//
	attrCustom := `1000xYYYY`
	thumb, _ = thumb_.NewThumb(keyThumb, thumb_.TestTypEmpty, logger, storage, conv)
	thumb.
		TestCacheSwap(cache).
		SetSize(size)

	thumb.SetAttrSize_Custom(attrCustom)
	tag = thumb.GetTAG()
	require.Equal(t, tag, fmt.Sprintf(`<link sizes="%s" >`, attrCustom))

}


func Test_OriginalCustomSet( t *testing.T ) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := &logs_mock_.LoggerErrorf{}

	keyThumb := `123`

	storage :=mock_storage_.NewMockStorage(ctrl)
	storage.EXPECT().NewObject(keyThumb).AnyTimes()

	conv := mock_converter_.NewMockConverter(ctrl)
	conv.EXPECT().Do(nil, nil, false, nil, 0).AnyTimes()

	obj :=mock_storage_.NewMockStorageOBJ(ctrl)
	cache := mock_thumb_.NewMockcache(ctrl)

	//
	thumb, _ := thumb_.NewThumb(keyThumb, thumb_.ICO, logger, storage, conv)
	thumb.
		TestCacheSwap(cache). 
		SetOriginal(obj)

	objExpect := thumb.GetOriginalStorageObj()
	require.Equal(t, obj, objExpect)

	//
	thumb, _ = thumb_.NewThumb(keyThumb, thumb_.ICO, logger, storage, conv)
	thumb.
		TestCacheSwap(cache). 
		SetOriginal(obj)

	objExpect = thumb.GetOriginalStorageObj()
	require.Equal(t, obj, objExpect)
}






func Test_Read( t *testing.T ) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := &logs_mock_.LoggerErrorf{}

	keyThumb := `123`
	instanceData := []byte(`1234`)
	instanceReader := io.NopCloser(bytes.NewBuffer(instanceData))

	storageObj :=mock_storage_.NewMockStorageOBJ(ctrl)
	storageObj.EXPECT().IsExists().Return(true, (error)(nil))
	storageObj.EXPECT().Reader().Return(instanceReader, nil)
	
	storage :=mock_storage_.NewMockStorage(ctrl)
	storage.EXPECT().NewObject(keyThumb).Return(storageObj, (error)(nil))

	conv := mock_converter_.NewMockConverter(ctrl)
	conv.EXPECT().Do(nil, nil, false, nil, 0).AnyTimes()

	cache := mock_thumb_.NewMockcache(ctrl)

	//
	thumb, _ := thumb_.NewThumb(keyThumb, thumb_.ICO, logger, storage, conv)
	thumb.TestCacheSwap(cache)

	expectReader, err := thumb.Read()
	require.Equal(t, err, (error)(nil))

	expectData, err := io.ReadAll(expectReader)
	require.Equal(t, err, (error)(nil))

	require.Equal(t, expectData, instanceData)
}

func Test_Read_ExistError( t *testing.T ) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := &logs_mock_.LoggerErrorf{}

	keyThumb := `123`

	storageObj :=mock_storage_.NewMockStorageOBJ(ctrl)
	storageObj.EXPECT().IsExists().Return(true, errors.New(`error exist`))
	
	storage :=mock_storage_.NewMockStorage(ctrl)
	storage.EXPECT().NewObject(keyThumb).Return(storageObj, (error)(nil))

	conv := mock_converter_.NewMockConverter(ctrl)
	conv.EXPECT().Do(nil, nil, false, nil, 0).AnyTimes()

	cache := mock_thumb_.NewMockcache(ctrl)

	//
	thumb, _ := thumb_.NewThumb(keyThumb, thumb_.ICO, logger, storage, conv)
	thumb.TestCacheSwap(cache)

	_, err := thumb.Read()
	require.Equal(t, err, logger.Error(thumb_.LogTP, thumb_.LogT10, errors.New(`error exist`)))

}

func Test_Read_CreateOriginalIsNil( t *testing.T ) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := &logs_mock_.LoggerErrorf{}

	keyThumb := `123`

	originalObj :=mock_storage_.NewMockStorageOBJ(ctrl)

	thumbObj :=mock_storage_.NewMockStorageOBJ(ctrl)
	thumbObj.EXPECT().IsExists().Return(false, (error)(nil))
	
	storage :=mock_storage_.NewMockStorage(ctrl)
	storage.EXPECT().NewObject(keyThumb).Return(thumbObj, (error)(nil))

	cache := mock_thumb_.NewMockcache(ctrl)


	conv := mock_converter_.NewMockConverter(ctrl)
	conv.EXPECT().Do(originalObj, thumbObj, false, types_.ICO(), 16).Return((error)(nil)).AnyTimes()

	thumb, _ := thumb_.NewThumb(keyThumb, thumb_.ICO, logger, storage, conv)
	thumb.TestCacheSwap(cache)

	_, err := thumb.Read()
	require.Equal(t, err, logger.Error(thumb_.LogTP, thumb_.LogT11, logger.Error(thumb_.LogTP, thumb_.LogT03)))

}


func Test_Read_Create( t *testing.T ) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := &logs_mock_.LoggerErrorf{}

	keyThumb := `123`
	size := 16

	originalObj :=mock_storage_.NewMockStorageOBJ(ctrl)

	thumbObj :=mock_storage_.NewMockStorageOBJ(ctrl)
	thumbObj.EXPECT().IsExists().Return(false, (error)(nil)).AnyTimes()
	thumbObj.EXPECT().Reader().Return(nil, (error)(nil))
	
	storage :=mock_storage_.NewMockStorage(ctrl)
	storage.EXPECT().NewObject(keyThumb).Return(thumbObj, (error)(nil))

	cache := mock_thumb_.NewMockcache(ctrl)
	cache.EXPECT().Range(gomock.Any())

	//
	conv := mock_converter_.NewMockConverter(ctrl)
	conv.EXPECT().Do(originalObj, thumbObj, false, types_.ICO(), size).Return((error)(nil))

	thumb, _ := thumb_.NewThumb(keyThumb, thumb_.ICO, logger, storage, conv)
	thumb.TestCacheSwap(cache).SetOriginal(originalObj).SetSize(16)

	_, err := thumb.Read()
	require.Equal(t, err, (error)(nil))

}

func Test_Read_CreateConverterError( t *testing.T ) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := &logs_mock_.LoggerErrorf{}

	keyThumb := `123`
	size := 16

	originalObj :=mock_storage_.NewMockStorageOBJ(ctrl)

	thumbObj :=mock_storage_.NewMockStorageOBJ(ctrl)
	thumbObj.EXPECT().IsExists().Return(false, (error)(nil)).AnyTimes()
	
	storage :=mock_storage_.NewMockStorage(ctrl)
	storage.EXPECT().NewObject(keyThumb).Return(thumbObj, (error)(nil))

	cache := mock_thumb_.NewMockcache(ctrl)
	cache.EXPECT().Range(gomock.Any())

	//
	convErr := errors.New(`error converter`)
	conv := mock_converter_.NewMockConverter(ctrl)
	conv.EXPECT().Do(originalObj, thumbObj, false, types_.ICO(), size).Return( convErr )

	thumb, _ := thumb_.NewThumb(keyThumb, thumb_.ICO, logger, storage, conv)
	thumb.TestCacheSwap(cache).SetOriginal(originalObj).SetSize(16)

	_, err := thumb.Read()
	require.Equal(t, err, logger.Error(thumb_.LogTP, thumb_.LogT11, logger.Error(thumb_.LogTP, thumb_.LogT01, convErr)))

}

func Test_OriginalKeyGet( t *testing.T ) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := &logs_mock_.LoggerErrorf{}

	keyThumb := `123`

	instanceKey := storage_.StorageKey(`325`)
	originalObj :=mock_storage_.NewMockStorageOBJ(ctrl)
	originalObj.EXPECT().Key().Return(instanceKey)
	
	thumbObj :=mock_storage_.NewMockStorageOBJ(ctrl)
	
	storage :=mock_storage_.NewMockStorage(ctrl)
	storage.EXPECT().NewObject(keyThumb).Return(thumbObj, (error)(nil))

	cache := mock_thumb_.NewMockcache(ctrl)
	conv := mock_converter_.NewMockConverter(ctrl)

	thumb, _ := thumb_.NewThumb(keyThumb, thumb_.ICO, logger, storage, conv)
	thumb.TestCacheSwap(cache)

	expectKey := thumb.GetOriginalKey()
	require.Equal(t, expectKey, `` )

	//
	thumb.SetOriginal(originalObj)

	expectKey = thumb.GetOriginalKey()
	require.Equal(t, expectKey, string(instanceKey) )
}

func Test_Cache( t *testing.T ) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := &logs_mock_.LoggerErrorf{}

	keyThumb := `123`

	thumbObj :=mock_storage_.NewMockStorageOBJ(ctrl)

	storage :=mock_storage_.NewMockStorage(ctrl)
	storage.EXPECT().NewObject(keyThumb).Return(thumbObj, (error)(nil))

	conv := mock_converter_.NewMockConverter(ctrl)
	

	instanceKey := `tag`
	instanceTag := `<link test />`
	var cache sync.Map
	
	thumb, _ := thumb_.NewThumb(keyThumb, thumb_.TestTypEmpty, logger, storage, conv)
	thumb.TestCacheSwap(&cache)

	//
	cache.Store(instanceKey, instanceTag)
	expectTag := thumb.GetTAG()
	require.Equal(t, expectTag, instanceTag)

	//
	thumb.SetSize(0)
	expectTag = thumb.GetTAG()
	require.Equal(t, expectTag, ``)
}

func Test_URLPath_Get( t *testing.T ) {

	href := `https://domain.org/stories/icon.png`
	u, err := url.ParseRequestURI(href)
	require.Equal(t, err, (error)(nil))

	instanceThumb := thumb_.Thumb{}
	m := map[typ_.URLPath]*thumb_.Thumb{
		`/stories/icon.png`: &instanceThumb,
	}

	expectThumb, exist := thumb_.URLPath_Get(u.Path, m)
	require.True(t, exist, u.Path)
	require.Equal(t, *expectThumb, instanceThumb)
}

