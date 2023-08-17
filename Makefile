ifndef $(GOPATH)
    GOPATH=$(shell go env GOPATH)
    export GOPATH
endif

APP_NAME?=afav


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



clean_httpv1:
	rm -f ${APP_NAME}_httpv1

build_httpv1: clean_httpv1
	go build -o ${APP_NAME}_httpv1

run_httpv1: build_httpv1
	./${APP_NAME}_httpv1 -adapter="httpv1" -conf="conf_httpv1.yaml" -img="image.png" -svg="image.svg" 


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

MOCKS_FILEPATH = pkg/core/models/converter/converter.exe.go $\
pkg/core/models/converter/converter.type.go $\
pkg/core/models/converter/converter.go $\
pkg/core/models/storage/storage.go $\
pkg/core/models/storage/storage.key.go $\
pkg/core/models/storage/storage.obj.go $\
pkg/core/models/logger/logger.go $\
internal/service/thumb/thumb.go $\
internal/service/convert/convert.go $\
internal/service/convert/checks/source.go

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



