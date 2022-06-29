FROM opensuse/tumbleweed AS build
ENV CGO_ENABLED 1
ENV GOOS "linux"
ENV GOARCH "amd64"
ARG VERSION
ARG BUILD_ARGS
RUN zypper update
RUN zypper install -y opencv-devel tesseract-ocr-traineddata-greek tesseract-ocr-devel leptonica-devel git go gcc-c++
ADD . / open-ocr-reader/
WORKDIR /open-ocr-reader/
RUN go build ${BUILD_ARGS} -ldflags "-X github.com/epinio/epinio/internal/version.Version=${VERSION}" -o dist/oor-linux-amd64
