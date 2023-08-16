package defaults

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 23 March 2023
 */
import (
	thumb_ "github.com/jhekau/favicon/internal/service/thumb"
	converter_ "github.com/jhekau/favicon/pkg/core/models/converter"
	logger_ "github.com/jhekau/favicon/pkg/core/models/logger"
	storage_ "github.com/jhekau/favicon/pkg/core/models/storage"
	err_ "github.com/jhekau/favicon/internal/core/err"
)

const (
	logP = `github.com/jhekau/favicon/internal/service/defaults/defaults.go`
	logD1 = `D1: create new thumb`
)

type attrSize struct{
	empty bool
	customVal string
}

func Defaults(l logger_.Logger, s storage_.Storage, c converter_.Converter) ([]*thumb_.Thumb, error) {

	t := make([]*thumb_.Thumb, 0)
	for _, d := range []struct{
		key string
		typ thumb_.Typ
		urlPath string
		size int
		attrRel string
		attrSize attrSize
		manifest bool
	}{
		// Нам нужны sizes="any" для <link> на файл .ico, 
		// чтобы исправить ошибку Chrome, из-за которой он выбирает файл ICO вместо SVG.
		// <link rel="icon" href="/favicon.ico" sizes="any"><!-- 32×32 -->
		{
			urlPath: `/favicon.ico`, 
			typ: thumb_.ICO, 
			size: 32, 
			attrRel: `icon`, 
			attrSize: attrSize{ customVal: `any`},
		},

		// SVG: <link rel="icon" href="/icon.svg" type="image/svg+xml">
		{
			urlPath: `/icon.svg`, 
			typ: thumb_.SVG, 
			attrRel: `icon`, 
			attrSize: attrSize{ empty: true },
		},

		// APPLE

		// <link rel="apple-touch-icon" href="/touch-icon-iphone.png"> <!-- 180x180 -->
		{
			urlPath: `/touch-icon-iphone.png`, 
			size: 180, 
			typ: thumb_.PNG, 
			attrRel: `apple-touch-icon`, 
			attrSize: attrSize{ empty: true },
		},
		
		// <link rel="apple-touch-icon" sizes="152x152" href="/touch-icon-ipad.png">
		{ 
			urlPath: `/touch-icon-ipad.png`, 
			size: 152, 
			typ: thumb_.PNG, 
			attrRel: `apple-touch-icon`, 
		},

		// <link rel="apple-touch-icon" sizes="180x180" href="/touch-icon-iphone-retina.png">
		{
			urlPath: `/touch-icon-iphone-retina.png`, 
			size: 180, 
			typ: thumb_.PNG, 
			attrRel: `apple-touch-icon`, 
		},

		// <link rel="apple-touch-icon" sizes="167x167" href="/touch-icon-ipad-retina.png">
		{
			urlPath: `/touch-icon-ipad-retina.png`, 
			size: 167, 
			typ: thumb_.PNG, 
			attrRel: `apple-touch-icon`, 
		},

		// For all browsers

		// <link rel="icon" type="image/png" sizes="32x32" href="/favicon-32x32.png">
		{
			urlPath: `/favicon-32x32.png`, 
			size: 32, 
			typ: thumb_.PNG, 
			attrRel: `icon`, 
		},

		// <link rel="icon" type="image/png" sizes="16x16" href="/favicon-16x16.png">
		{
			urlPath: `/favicon-16x16.png`, 
			size: 16, 
			typ: thumb_.PNG, 
			attrRel: `icon`, 
		},

		// For Google and Chrome

		// <link rel="icon" type="image/png" sizes="48x48" href="/favicon-48x48.png">
		{
			urlPath: `/favicon-48x48.png`, 
			size: 48, 
			typ: thumb_.PNG, 
			attrRel: `icon`, 
		},

		// <link rel="icon" type="image/png" sizes="192x192" href="/favicon-192x192.png">
		{
			urlPath: `/favicon-192x192.png`, 
			size: 192, 
			typ: thumb_.PNG, 
			attrRel: `icon`, 
			manifest: true,
		},

		// <link rel="icon" type="image/png" sizes="512x512" href="/favicon-512x512.png">
		{
			urlPath: `/favicon-512x512.png`, 
			size: 512, 
			typ: thumb_.PNG, 
			attrRel: `icon`, 
			manifest: true,
		},
	}{
		tb, err := thumb_.NewThumb(d.urlPath, d.typ, l, s, c)
		if err != nil {
			return nil, err_.Err(l, logP, logD1, err)
		}
		tb.SetSize(d.size)
		tb.SetAttrRel(d.attrRel)
		tb.SetURLPath(d.urlPath)

		if v := d.attrSize.customVal; v != `` {
			tb.SetAttrSize_Custom(v)
		} else if d.attrSize.empty {
			tb.SetAttrSize_Empty()
		}
		
		if d.manifest {
			tb.SetManifestUsed()
		}

		t = append(t, tb)
	}
	return t, nil
}
