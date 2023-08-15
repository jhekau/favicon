# Favicon — automatic creation of popular icons for popular platform

doc 

*use cases*

*pkg use*

...

### Quick start
```
go run ...
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
    github.com/jhekau/favicon/pkg/core/models/logger
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
    github.com/jhekau/favicon/pkg/core/models/storage
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
    github.com/jhekau/favicon/pkg/core/models/converter
)

# check implementation
var _ storage.Converter = (YourConverter)(nil)

t := thumbs.NewThumbs()
t.ConvertSet( YourConverter )
```


#### Chapters

- v0.0.1: set architecture project, added dependency inversion, unit test, integration test;
- v0.0.0: PoC version;
