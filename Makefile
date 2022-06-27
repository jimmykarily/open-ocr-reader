VERSION ?= $(shell git describe --tags)
CGO_ENABLED ?= 0

export LDFLAGS += -X github.com/epinio/epinio/internal/version.Version=$(VERSION)

build:  build-amd64

build-amd64:
	GOARCH="amd64" GOOS="linux" CGO_ENABLED=$(CGO_ENABLED) go build $(BUILD_ARGS) -ldflags '$(LDFLAGS)' -o dist/oor-linux-amd64
