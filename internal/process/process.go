// Package process is responsible for taking a raw image (as capture by
// the "capture" package and process it to make it suitable for OCR.
package process

import (
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

// Process prepares a photo of a book page for OCR
// These are the steps
// - Make the image black and white
// - Find the biggest block of text
// - Find a containing rectangle of that block of text and deskew the image
//   based on that rectangle (align the text vertically)
// - Crop image to that rectangle
// Heavily inspried by this:
// https://github.com/JPLeoRX/opencv-text-deskew/blob/master/python-service/services/deskew_service.py
// https://becominghuman.ai/how-to-automatically-deskew-straighten-a-text-image-using-opencv-a0c30aed83df
// https://github.com/milosgajdos/gocv-playground/blob/master/04_Geometric_Transformations/README.md#perspective-transformation
func (p DefaultProcessor) Process(image *img.Image) (*img.Image, error) {

	imgPath, err := image.StoreTmp()
	if err != nil {
		return nil, errors.Wrap(err, "storing the image to a temp file")
	}
	defer os.Remove(imgPath)

	cvImg := gocv.IMRead(imgPath, gocv.IMReadColor)

	convertToGrayscale(&cvImg)
	deskew(&cvImg)
	//showImg(cvImg)

	result, err := cvImg.ToImage()
	if err != nil {
		return nil, errors.Wrap(err, "converting Mat to image")
	}
	image.Object = result

	return image, nil
}

func convertToGrayscale(i *gocv.Mat) {
	gocv.CvtColor(*i, i, gocv.ColorBGRToGray)
}

func deskew(i *gocv.Mat) {
	tmpImg := i.Clone()

	gocv.GaussianBlur(tmpImg, &tmpImg, goimage.Point{}, 1, 1, gocv.BorderDefault)
	_ = gocv.Threshold(tmpImg, &tmpImg, 0, 255, gocv.ThresholdBinaryInv+gocv.ThresholdOtsu)

	kernel := gocv.GetStructuringElement(gocv.MorphRect, goimage.Point{5, 10})
	gocv.DilateWithParams(tmpImg, &tmpImg, kernel, goimage.Point{}, 5, gocv.BorderDefault, color.RGBA{})

	points := gocv.FindContours(tmpImg, gocv.RetrievalList, gocv.ChainApproxSimple)
	contours := ContoursBySize{}
	for j := 1; j < points.Size(); j++ {
		contours = append(contours, Contour{OriginalIdx: j, Contour: points.At(j)})
	}
	sort.Sort(contours)

	// TODO: Debug lines, remove
	// colors := []color.RGBA{
	// 	{255, 0, 0, 255},
	// 	{0, 255, 0, 255},
	// 	{0, 0, 255, 255},
	// }
	// for j := range contours {
	// 	gocv.DrawContours(&tmpImg, points, j, colors[j%len(colors)], 1)
	// }
	maxContour := contours[len(contours)-1]
	//gocv.DrawContours(&tmpImg, points, maxContour.OriginalIdx, color.RGBA{100, 80, 20, 255}, 2)

	rect := gocv.MinAreaRect(maxContour.Contour)
	//rectV := gocv.NewPointsVectorFromPoints([][]goimage.Point{rect.Points})
	//gocv.DrawContours(i, rectV, -1, color.RGBA{0, 255, 0, 255}, 3)
	// showImg(*i)

	// https://github.com/milosgajdos/gocv-playground/blob/master/04_Geometric_Transformations/README.md#perspective-transformation
	newPoints := []goimage.Point{
		{0, rect.Height},
		{0, 0},
		{rect.Width, 0},
		{rect.Width, rect.Height},
	}

	transform := gocv.GetPerspectiveTransform(
		gocv.NewPointVectorFromPoints(rect.Points),
		gocv.NewPointVectorFromPoints(newPoints),
	)
	gocv.WarpPerspective(*i, i, transform, goimage.Point{rect.Width, rect.Height})

	// TODO: Not needed? We do it in one step above
	// TODO: If we want to detect more than one blocks of text, then we need
	// to first deskew the original image without cropping.
	//
	// skewAngle := calculateSkewAngle(rect.Angle)
	// rotateImg(i, skewAngle)
}

// calculateSkewAngle take the angle of the min area rectagle and return the
// angle to rotate the image in order to deskew the document.
// Currently WarpPerspective does both in one step.
func calculateSkewAngle(angle float64) float64 {
	if angle < -45 {
		return 90 + angle
	}
	if angle > 45 {
		return -1 * (90 - angle)
	}

	return angle
}

// Rotate the image around its center
func rotateImg(i *gocv.Mat, angle float64) {
	size := i.Size()
	width := size[1]
	height := size[0]
	rMatrix := gocv.GetRotationMatrix2D(goimage.Point{X: width / 2.0, Y: height / 2.0}, angle, 1.0)
	gocv.WarpAffineWithParams(*i, i, rMatrix, goimage.Point{X: width, Y: height}, gocv.InterpolationCubic, gocv.BorderReplicate, color.RGBA{0, 0, 0, 0})
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
