ifndef $(GOPATH)
    GOPATH=$(shell go env GOPATH)
    export GOPATH
endif


# *****
# solves the problem on the windows platform
# решает проблему на винде
# ------------------------------------------
# process_begin: CreateProcess(NULL, Write-Output asdf/asdf, ...) failed.
# make (e=2): ═х єфрхЄё  эрщЄш єърчрээ√щ Їрщы.

ifeq ($(OS),Windows_NT)
SHELL := powershell.exe
.SHELLFLAGS := -NoProfile -Command
endif




.PHONY: cover
cover:
	@echo "Generating coverprofile from test..."
	go test -short -count=1 -race -coverprofile="coverage.out" ./...
	go tool cover -html="coverage.out"
	rm coverage.out



# *****
# rename the destination for subfolders - "internal", otherwise it will be impossible to import packages
# ренейм для вложенных папок назначения - "internal", такая сложнота нужна, в противном случае импортировать моки будет нереал
# <!> go install go.uber.org/mock/mockgen@latest

MOCKS_FILEPATH = pkg/models/converter/converter.exe.go $\
pkg/models/converter/converter.type.go $\
pkg/models/converter/converter.go $\
pkg/models/storage/storage.go $\
pkg/models/storage/storage.key.go $\
pkg/models/storage/storage.obj.go $\
internal/service/thumb/thumb.go $\
internal/service/convert/convert.go $\
internal/service/convert/convert/checks/source.go

MOCKGEN ?= $(GOPATH)/bin/mockgen
ifeq ($(OS),Windows_NT)
MOCKGEN := $(GOPATH)"\bin\mockgen.exe"
endif
MOCKS_DESTINATION ?= internal/mocks

.PHONY: mockgen
mockgen: $(MOCKS_FILEPATH)
	@echo "Generating mocks..."
	@rm -rf $(MOCKS_DESTINATION)
	@for file in $^; 									   \
		do			 									   \
			echo $$file | sed "s/internal/intr/g" ; \
			$(MOCKGEN) -source=$$file -destination=$(MOCKS_DESTINATION)/`echo $$file | sed "s/internal/intr/g"`; \
		done;



