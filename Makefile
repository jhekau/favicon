ifndef $(GOPATH)
    GOPATH=$(shell go env GOPATH)
    export GOPATH
endif

# *****
# process_begin: CreateProcess(NULL, Write-Output asdf/asdf, ...) failed.
# make (e=2): ═х єфрхЄё  эрщЄш єърчрээ√щ Їрщы.
ifeq ($(OS),Windows_NT)
SHELL := powershell.exe
.SHELLFLAGS := -NoProfile -Command
endif


##
.PHONY: clean build run cover mockgen

APP_NAME?=autofav
APP_PATH?=cmd/http

clean:
	rm -f ${APP_NAME}

build: clean
	go build -o ${APP_NAME} ./cmd/http

run: build
	./${APP_NAME} -conf="http.yaml" -img="img.jpg" -svg="image.svg" 

cover:
	@echo "Generating coverprofile from test..."
	go test -short -count=1 -race -coverprofile="coverage.out" ./...
	go tool cover -html="coverage.out"
	rm coverage.out

mockgen: 
	scripts/mockgen.sh



