// Package img implements an Image class that can be share between the other
// packages.
package img

import (
	"image"
	"io/ioutil"
	"os"

	"image/png"

	"github.com/pkg/errors"
)

type Image struct {
	Object image.Image
}

func New(imgPath string) (*Image, error) {
	f, err := os.Open(imgPath)
	if err != nil {
		return nil, errors.Wrap(err, "reading image file")
	}
	defer f.Close()

	imgObject, _, err := image.Decode(f)
	if err != nil {
		return nil, errors.Wrap(err, "decoding the image")
	}

	return &Image{Object: imgObject}, nil
}

// StoreTmp can be used to store the image to a temporary location.
// Some libraries only work with file paths, not Image objects.
// It's the callers responsibility to delete the temporary file.
// Returns the path to the temporary file or an error if something goes wrong.
func (image Image) StoreTmp() (string, error) {
	file, err := ioutil.TempFile("", "oor")
	if err != nil {
		return "", errors.Wrap(err, "creating a temporary file")
	}
	defer file.Close()

	// Encode to `PNG` with `DefaultCompression` level then save to file
	err = png.Encode(file, image.Object)
	if err != nil {
		return "", errors.Wrap(err, "encoding the image as png")
	}

	return file.Name(), nil
}
