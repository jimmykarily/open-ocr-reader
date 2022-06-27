// Package ocr is reponsible for extracting text from an image
package ocr

import (
	"github.com/jimmykarily/open-ocr-reader/internal/img"
	"github.com/otiai10/gosseract/v2"
	"github.com/pkg/errors"
)

type OCR interface {
	Parse(img *img.Image) (string, error)
}

type TesseractOCR struct{}

func NewTesseractOCR() TesseractOCR {
	return TesseractOCR{}
}

func (t TesseractOCR) Parse(img *img.Image) (string, error) {
	//l, _ := gosseract.GetAvailableLanguages()
	//fmt.Printf("l = %+v\n", l)

	imgPath, err := img.StoreTmp()
	if err != nil {
		return "", errors.Wrap(err, "storing the image to a temp file")
	}

	client := gosseract.NewClient()
	//client.Languages = []string{"ell", "eng"}
	client.Languages = []string{"ell"}
	defer client.Close()
	client.SetImage(imgPath)
	text, err := client.Text()
	if err != nil {
		return "", errors.Wrap(err, "detecting text")
	}

	return text, nil
}
