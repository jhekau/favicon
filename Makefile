# solves the problem on the windows platform:
# process_begin: CreateProcess(NULL, Write-Output asdf/asdf, ...) failed.
# make (e=2): ═х єфрхЄё  эрщЄш єърчрээ√щ Їрщы.
ifeq ($(OS),Windows_NT)
SHELL := powershell.exe
.SHELLFLAGS := -NoProfile -Command
endif

cover:
	go test -short -count=1 -race -coverprofile="coverage.out" ./...
	go tool cover -html="coverage.out"
	rm coverage.out
