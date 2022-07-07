VERSION ?= $(shell git describe --tags)
CGO_ENABLED ?= 1 # https://github.com/otiai10/gosseract/issues/138#issuecomment-488233324
USER_ID ?= $(shell id -u ${USER})
USER_GID ?= $(shell id -g ${USER})

export LDFLAGS += -X github.com/epinio/epinio/internal/version.Version=$(VERSION)

build:  build-amd64

# Build inside docker as the current user
container-build:
	docker run -u $(USER_ID):$(USER_GID) --env GOCACHE=/oor/tmp/gocache --env GOPATH=/oor/tmp/gomodcache -v "${PWD}":/oor --workdir /oor open-ocr-reader-ubuntu make

build-amd64:
	GOARCH="amd64" GOOS="linux" CGO_ENABLED=$(CGO_ENABLED) go build $(BUILD_ARGS) -ldflags '$(LDFLAGS)' -o dist/oor-linux-amd64

build-container-image:
	docker build --build-arg VERSION=$(VERSION) . -t open-ocr-reader:$(VERSION)

build-container-image-ubuntu:
	docker build --build-arg VERSION=$(VERSION) . -t open-ocr-reader-ubuntu -f Dockerfile-ubuntu

test:
	ginkgo ./...
