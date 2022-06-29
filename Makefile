VERSION ?= $(shell git describe --tags)
CGO_ENABLED ?= 1 # https://github.com/otiai10/gosseract/issues/138#issuecomment-488233324

export LDFLAGS += -X github.com/epinio/epinio/internal/version.Version=$(VERSION)

build:  build-amd64

build-amd64:
	docker build --build-arg VERSION=$(VERSION) . -t open-ocr-reader:$(VERSION)
test:
	ginkgo ./...
