# Favicon — automatic creation of popular icons for popular platform
# Demo <!>

### Introduction
The program allows you to add a favicon to the project. It is enough to put the original images (svg, img) in the root of the site, launch the application. It will automatically cut previews for popular formats and raise the necessary service.

By default, the service generates a manifest at domain/manifest.webmanifest (!in operation).

You can override the behavior, for example, add your own sets, change the storage system, add your own logger implementation, and more.

### Default icons
```
<link rel="icon" href="/favicon.ico" sizes="any"><!-- 32×32 -->
<link rel="icon" href="/icon.svg" type="image/svg+xml">

// APPLE
<link rel="apple-touch-icon" href="/touch-icon-iphone.png"> <!-- 180x180 -->
<link rel="apple-touch-icon" sizes="152x152" href="/touch-icon-ipad.png">
<link rel="apple-touch-icon" sizes="180x180" href="/touch-icon-iphone-retina.png">
<link rel="apple-touch-icon" sizes="167x167" href="/touch-icon-ipad-retina.png">

// For all browsers
<link rel="icon" type="image/png" sizes="32x32" href="/favicon-32x32.png">
<link rel="icon" type="image/png" sizes="16x16" href="/favicon-16x16.png">

// For Google and Chrome
<link rel="icon" type="image/png" sizes="48x48" href="/favicon-48x48.png">
<link rel="icon" type="image/png" sizes="192x192" href="/favicon-192x192.png">
<link rel="icon" type="image/png" sizes="512x512" href="/favicon-512x512.png">
```

### Quick start
```
go run ...
```

You can get image and manifest tags to add to your project:
```Go
add docs
icons := ...
```

### Изменения директории хранения нарезанных иконок стандартной системы хранения
Новые иконки хранятся по умолчанию в директории icons рабочего бинарника, чтобы изменить место хранения:
```Go
import (
    github.com/jhekau/favicon/pkg/service/storage_default
)
storagedefault.SetFolderIcons(YourFolder)
```

### Использование альтернативной реализации логгера:
```Go
import (
    github.com/jhekau/favicon/pkg/thumbs
    github.com/jhekau/favicon/interfaces/logger
)

# check implementation
var _ logger.Logger = (YourLogger)(nil)

t := thumbs.NewThumbs()
t.LoggerSet( YourLogger )
```

### Использование альтернативной системы хранения:
```Go
import (
    github.com/jhekau/favicon/pkg/thumbs
    github.com/jhekau/favicon/interfaces/storage
)

# check implementation
var _ storage.Storage = (YourStorage)(nil)

t := thumbs.NewThumbs()
t.StorageSet( YourStorage )
```

### Использование альтернативного конвертера для создания превьюх:
```Go
import (
    github.com/jhekau/favicon/pkg/thumbs
    github.com/jhekau/favicon/interfaces/converter
)

# check implementation
var _ storage.Converter = (YourConverter)(nil)

t := thumbs.NewThumbs()
t.ConvertSet( YourConverter )
```


#### Chapters
- v0.0.2: working version, partially tested;
- v0.0.1: set architecture project, added dependency inversion, unit test, integration test;
- v0.0.0: PoC version;

#### TODO`s
- [x] test unit;
- [x] test integration (thumb-convert-conv.exec) (part);
- [x] debug;
- [x] clean arch;
- [x] scripts;
- [ ] Thumb - manifest, test;
- [ ] Docker - manifest, deploy;
- [ ] other testing;
- [ ] doc local;
- [ ] doc other;
