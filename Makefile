TAG_VERSION := $(shell git describe --tags --always)

test:
	go test  -v -coverpkg=./... -race -covermode=atomic -coverprofile=coverage.txt ./... -run . -timeout=2m

build:
	go build -ldflags='-s -w -X main.TAG_version=${TAG_VERSION}' -o archive ./cli/main.go