// Package process is responsible for taking a raw image (as capture by
// the "capture" package and process it to make it suitable for OCR.
package process

import (
	"fmt"
	goimage "image"
	"image/color"
	"os"
	"sort"

	"github.com/jimmykarily/open-ocr-reader/internal/img"
	"github.com/pkg/errors"
	"gocv.io/x/gocv"
)

type Processor interface {
	Process(*img.Image) (*img.Image, error)
}

type DefaultProcessor struct{}

type Contour struct {
	OriginalIdx int
	Contour     gocv.PointVector
}
type ContoursBySize []Contour

func (s ContoursBySize) Len() int { return len(s) }
func (s ContoursBySize) Less(i, j int) bool {
	return gocv.ContourArea(s[i].Contour) < gocv.ContourArea(s[j].Contour)
}
func (s ContoursBySize) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// NewDefaultProcessor returns a DefaultProcessor
func NewDefaultProcessor() DefaultProcessor {
	return DefaultProcessor{}
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
	// imgObject := image.Object
	// imgObject = imaging.Grayscale(imgObject)
	// imgObject = imaging.AdjustContrast(imgObject, 20)
	// imgObject = imaging.Sharpen(imgObject, 2)

	imgPath, err := image.StoreTmp()
	if err != nil {
		return nil, errors.Wrap(err, "storing the image to a temp file")
	}
	defer os.Remove(imgPath)

	cvImgOrigin := gocv.IMRead(imgPath, gocv.IMReadColor)
	//showImg(cvImg)

	cvImg := cvImgOrigin.Clone()

	gocv.CvtColor(cvImg, &cvImg, gocv.ColorBGRToGray)
	showImg(cvImg)

	gocv.GaussianBlur(cvImg, &cvImg, goimage.Point{}, 1, 1, gocv.BorderDefault)
	//showImg(cvImg)

	_ = gocv.Threshold(cvImg, &cvImg, 0, 255, gocv.ThresholdBinaryInv+gocv.ThresholdOtsu)

	kernel := gocv.GetStructuringElement(gocv.MorphRect, goimage.Point{5, 10})

	gocv.DilateWithParams(cvImg, &cvImg, kernel, goimage.Point{}, 5, gocv.BorderDefault, color.RGBA{})
	showImg(cvImg)

	// gocv.ErodeWithParams(cvImg, &cvImg, kernel, goimage.Point{}, 5, int(gocv.BorderDefault))
	// showImg(cvImg)

	points := gocv.FindContours(cvImg, gocv.RetrievalList, gocv.ChainApproxSimple)
	contours := ContoursBySize{}
	for i := 1; i < points.Size(); i++ {
		contours = append(contours, Contour{OriginalIdx: i, Contour: points.At(i)})
	}
	sort.Sort(contours)

	colors := []color.RGBA{
		{255, 0, 0, 255},
		{0, 255, 0, 255},
		{0, 0, 255, 255},
	}
	totalColors := len(colors)
	for i, c := range contours {
		gocv.DrawContours(&cvImgOrigin, points, i, colors[i%totalColors], 4)
		fmt.Printf("gocv.ContourArea(c) = %+v\n", gocv.ContourArea(c.Contour))
	}

	gocv.DrawContours(&cvImgOrigin, points, contours[len(contours)-1].OriginalIdx, color.RGBA{100, 80, 20, 255}, 10)
	showImg(cvImgOrigin)

	return image, nil
}

func showImg(i gocv.Mat) {
	w := gocv.NewWindow("Debug")

	for {
		w.IMShow(i)
		if w.WaitKey(1) >= 0 {
			break
		}
	}
	w.Close()
}
