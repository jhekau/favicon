package defaults

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 23 March 2023
 */
import (
	thumb_ "github.com/jhekau/favicon/internal/service/thumb"
	"github.com/jhekau/favicon/internal/core/types"
)

func Defaults() []*thumb_.Thumb {
	return []*thumb_.Thumb{
		// Нам нужны sizes="any" для <link> на файл .ico, 
		// чтобы исправить ошибку Chrome, из-за которой он выбирает файл ICO вместо SVG.
		// <link rel="icon" href="/favicon.ico" sizes="any"><!-- 32×32 -->
		(&thumb_.Thumb{}).SetTagRel(`icon`).URLPathSet(`/favicon.ico`).SetSize(32).SetType(types.ICO()).SetSizeAttrCustom(`any`),

		// SVG
		// <link rel="icon" href="/icon.svg" type="image/svg+xml">
		(&thumb_.Thumb{}).SetTagRel(`icon`).URLPathSet(`/icon.svg`).SetType(types.SVG()).SetSizeAttrEmpty(),

		// APPLE
		// <link rel="apple-touch-icon" href="/touch-icon-iphone.png"> <!-- 180x180 -->
		// <link rel="apple-touch-icon" sizes="152x152" href="/touch-icon-ipad.png">
		// <link rel="apple-touch-icon" sizes="180x180" href="/touch-icon-iphone-retina.png">
		// <link rel="apple-touch-icon" sizes="167x167" href="/touch-icon-ipad-retina.png">
		(&thumb_.Thumb{}).SetTagRel(`apple-touch-icon`).URLPathSet(`/touch-icon-iphone.png`).SetSize(180).SetType(types.PNG()).SetSizeAttrEmpty(),
		(&thumb_.Thumb{}).SetTagRel(`apple-touch-icon`).URLPathSet(`/touch-icon-ipad.png`).SetSize(152).SetType(types.PNG()),
		(&thumb_.Thumb{}).SetTagRel(`apple-touch-icon`).URLPathSet(`/touch-icon-iphone-retina.png`).SetSize(180).SetType(types.PNG()),
		(&thumb_.Thumb{}).SetTagRel(`apple-touch-icon`).URLPathSet(`/touch-icon-ipad-retina.png`).SetSize(167).SetType(types.PNG()),

		// For all browsers
		// <link rel="icon" type="image/png" sizes="32x32" href="/favicon-32x32.png">
		// <link rel="icon" type="image/png" sizes="16x16" href="/favicon-16x16.png">
		(&thumb_.Thumb{}).SetTagRel(`icon`).URLPathSet(`/favicon-32x32.png`).SetSize(32).SetType(types.PNG()),
		(&thumb_.Thumb{}).SetTagRel(`icon`).URLPathSet(`/favicon-16x16.png`).SetSize(13).SetType(types.PNG()),

		// For Google and Chrome
		// <link rel="icon" type="image/png" sizes="48x48" href="/favicon-48x48.png">
		// <link rel="icon" type="image/png" sizes="192x192" href="/favicon-192x192.png">
		// <link rel="icon" type="image/png" sizes="512x512" href="/favicon-512x512.png">
		(&thumb_.Thumb{}).SetTagRel(`icon`).URLPathSet(`/favicon-48x48.png`).SetSize(48).SetType(types.PNG()),
		(&thumb_.Thumb{}).SetTagRel(`icon`).URLPathSet(`/favicon-192x192.png`).SetSize(192).SetType(types.PNG()).SetManifestUsed(),
		(&thumb_.Thumb{}).SetTagRel(`icon`).URLPathSet(`/favicon-512x512.png`).SetSize(512).SetType(types.PNG()).SetManifestUsed(),
	}
}
