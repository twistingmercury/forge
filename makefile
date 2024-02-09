default: build

BIN_DIR=./bin/
BIN_NAME=forge
BUILD_DATE=$(shell date +"%Y-%m-%d")
BUILD_VER=0.0.2-$(shell date +"%Y%m%d%H%M%S")
GIT_COMMIT=$(shell git rev-parse --short HEAD)

build:
	go build \
	-ldflags "-X 'main.date=$(BUILD_DATE)' -X 'main.version=$(BUILD_VER)' -X 'main.commit=$(GIT_COMMIT)'" \
	-o $(BIN_DIR)/$(BIN_NAME) .

test: build
	go test -v ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out