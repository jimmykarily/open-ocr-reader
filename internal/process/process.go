// Package process is responsible for taking a raw image (as capture by
// the "capture" package and process it to make it suitable for OCR.
package process

import (
	"github.com/disintegration/imaging"
	"github.com/jimmykarily/open-ocr-reader/internal/img"
)

type DefaultProcessor struct{}

// NewDefaultProcessor returns a DefaultProcessor
func NewDefaultProcessor() DefaultProcessor {
	return DefaultProcessor{}
}

type Processor interface {
	Process(*img.Image) (*img.Image, error)
}

// Process should do at least these:
// - Make the image black and white
// - Deskew (align the text vertically)
// - Remove noise
// - Improve contrast
// - Detect block of text (find bounding box)
// - Crop image to the bounding box (opencv)
// - Apply 4 point image transformation (opencv)
func (p DefaultProcessor) Process(image *img.Image) (*img.Image, error) {
	imgObject := image.Object

	imgObject = imaging.Grayscale(imgObject)
	imgObject = imaging.AdjustContrast(imgObject, 20)
	imgObject = imaging.Sharpen(imgObject, 2)

	return image, nil
}
