VERSION ?= $(shell git describe --tags)
CGO_ENABLED ?= 1 # https://github.com/otiai10/gosseract/issues/138#issuecomment-488233324

export LDFLAGS += -X github.com/epinio/epinio/internal/version.Version=$(VERSION)

build:  build-amd64

build-amd64:
	GOARCH="amd64" GOOS="linux" CGO_ENABLED=$(CGO_ENABLED) go build $(BUILD_ARGS) -ldflags '$(LDFLAGS)' -o dist/oor-linux-amd64

test:
	ginkgo ./...
