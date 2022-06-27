// Package process is responsible for taking a raw image (as capture by
// the "capture" package and process it to make it suitable for OCR.
package process

import "github.com/jimmykarily/open-ocr-reader/internal/image"

func Process(img image.Image) (image.Image, error) {
	return img, nil
}
