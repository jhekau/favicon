# solves the problem on the windows platform:
# process_begin: CreateProcess(NULL, Write-Output asdf/asdf, ...) failed.
# make (e=2): ═х єфрхЄё  эрщЄш єърчрээ√щ Їрщы.

ifndef $(GOPATH)
    GOPATH=$(shell go env GOPATH)
    export GOPATH
endif

MOCKGEN ?= $(GOPATH)/bin/mockgen

MOCKS_DESTINATION ?= internal/mocks
ifeq ($(OS),Windows_NT)

SHELL := powershell.exe
.SHELLFLAGS := -NoProfile -Command

MOCKGEN := $(GOPATH)"\bin\mockgen.exe"
endif


cover:
	@echo "Generating coverprofile from test..."
	go test -short -count=1 -race -coverprofile="coverage.out" ./...
	go tool cover -html="coverage.out"
	rm coverage.out

.PHONY: mockgen
mockgen: internal/service/convert/convert.go internal/service/convert/checks/source.go
	@echo "Generating mocks..."
	@rm -rf $(MOCKS_DESTINATION)
	@for file in $^; 															\
		do			 															\
			IFS=', '; read -r -a array <<< "$string"; \
			for element in $$array; 		\
			do 							\
				echo $$element'*'; 			\
			done; 						\
			$(MOCKGEN) -source=$$file -destination=$(MOCKS_DESTINATION)/$$file; \
		done;
