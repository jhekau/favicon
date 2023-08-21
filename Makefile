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





APP_NAME?=autofav
APP_PATH?=cmd/http/v1


clean:
	rm -f ${APP_NAME}

build_http_v1: clean
	go build -o ${APP_NAME} ./cmd/http/v1

# example
# run_http_v1: build_httpv1
# 	./${APP_NAME} -conf="conf.yaml" -img="image.png" -svg="image.svg" 


.PHONY: cover
cover:
	@echo "Generating coverprofile from test..."
	go test -short -count=1 -race -coverprofile="coverage.out" ./...
	go tool cover -html="coverage.out"
	rm coverage.out

.PHONY: mockgen
mockgen: 
	scripts/mockgen.sh



