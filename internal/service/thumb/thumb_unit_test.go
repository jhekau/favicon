package thumb_test

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 10 August 2023
 */
import (
	"testing"

	logger_ "github.com/jhekau/favicon/internal/core/logger"
	logger_mock_ "github.com/jhekau/favicon/internal/core/logger/mock"
	thumb_ "github.com/jhekau/favicon/internal/service/thumb"
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

func Test_NewThumb( t *testing ) {

	logger := &logger_.Logger{
		Typ: &logger_mock_.LoggerErrorf{},
	}


	thumb_.NewThumb(
		logger,
		s thumb_.Storage, c thumb_.Converter)

}
