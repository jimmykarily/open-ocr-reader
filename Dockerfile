FROM opensuse/tumbleweed AS build
RUN zypper refresh
RUN zypper install -y opencv-devel tesseract-ocr-devel leptonica-devel git go gcc-c++ openblas-devel libgomp1 libatomic1 python38 python38-pip python38-devel
RUN pip3 install --upgrade pip
RUN pip3 install --upgrade wheel setuptools
RUN pip3 install -f 'https://synesthesiam.github.io/prebuilt-apps/' -f 'https://download.pytorch.org/whl/cpu/torch_stable.html' larynx==1.1.0

FROM build AS ocr-container
ENV CGO_ENABLED 1
ENV GOOS "linux"
ENV GOARCH "amd64"
ARG VERSION
ARG BUILD_ARGS
RUN zypper install -y tesseract-ocr-traineddata-english tesseract-ocr-traineddata-greek
ADD . / open-ocr-reader/
WORKDIR /open-ocr-reader/
RUN go build ${BUILD_ARGS} -ldflags "-X github.com/epinio/epinio/internal/version.Version=${VERSION}" -o dist/oor-linux-amd64
ENTRYPOINT ["dist/oor-linux-amd64", "server"]
