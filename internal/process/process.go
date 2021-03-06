// Package process is responsible for taking a raw image (as capture by
// the "capture" package and process it to make it suitable for OCR.
package process

import (
	"fmt"
	goimage "image"
	"image/color"
	"math"
	"os"
	"path/filepath"
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
// Heavily inspired by these:
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
	storeDebug(&cvImg, "1-original")

	convertToGrayscale(&cvImg)
	storeDebug(&cvImg, "2-grayscale")

	deskew(&cvImg)

	_ = gocv.Threshold(cvImg, &cvImg, 127, 255, gocv.ThresholdBinary+gocv.ThresholdOtsu)
	storeDebug(&cvImg, "12-black-and-white")

	// tesseract likes borders:
	// https://tesseract-ocr.github.io/tessdoc/ImproveQuality#dilation-and-erosion
	gocv.CopyMakeBorder(cvImg, &cvImg, 10, 10, 10, 10, gocv.BorderConstant, color.RGBA{100, 100, 100, 255})
	storeDebug(&cvImg, "13-withborder")

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
	defer tmpImg.Close()

	// TODO: Blurring doesn't seem to improve things. Maybe it would work with
	// different Thresholding below.
	// gocv.GaussianBlur(tmpImg, &tmpImg, goimage.Point{}, 1, 1, gocv.BorderDefault)
	// storeDebug(&tmpImg, "3-after-gaussionblur")

	_ = gocv.Threshold(tmpImg, &tmpImg, 100, 255, gocv.ThresholdBinaryInv) //+gocv.ThresholdOtsu)
	storeDebug(&tmpImg, "4-after-threshold")

	// TODO: this is a hack. We decide what the kernel size is based on the resolution
	// of the image. This means, we assume what the approximate size of each character is
	// a certain percentage of the total page.
	kernelWidth := (i.Cols() / 150) / 2
	kernelHeight := i.Rows() / 150
	kernel := gocv.GetStructuringElement(gocv.MorphRect, goimage.Point{kernelWidth, kernelHeight})
	gocv.DilateWithParams(tmpImg, &tmpImg, kernel, goimage.Point{}, 5, gocv.BorderDefault, color.RGBA{})
	storeDebug(&tmpImg, "5-after-dilate")

	points := gocv.FindContours(tmpImg, gocv.RetrievalList, gocv.ChainApproxSimple)
	contours := ContoursBySize{}
	for j := 1; j < points.Size(); j++ {
		contours = append(contours, Contour{OriginalIdx: j, Contour: points.At(j)})
	}
	sort.Sort(contours)

	contouredImage := i.Clone()
	defer contouredImage.Close()
	colors := []color.RGBA{
		{255, 0, 0, 255},
		{0, 255, 0, 255},
		{0, 0, 255, 255},
	}
	for j := range contours {
		gocv.DrawContours(&contouredImage, points, j, colors[j%len(colors)], 1)
	}
	maxContour := contours[len(contours)-1]
	gocv.DrawContours(&contouredImage, points, maxContour.OriginalIdx, color.RGBA{100, 80, 20, 255}, 2)
	storeDebug(&contouredImage, "6-contoured")

	rect := gocv.MinAreaRect(maxContour.Contour)

	// Debug
	originalCopy := i.Clone()
	defer originalCopy.Close()
	rectV := gocv.NewPointsVectorFromPoints([][]goimage.Point{rect.Points})
	gocv.DrawContours(&originalCopy, rectV, -1, color.RGBA{0, 255, 0, 255}, 3)
	storeDebug(&originalCopy, "7-min-rectangle")

	skewAngle := calculateSkewAngle(rect.Angle)
	rotateImg(i, rect.Center, skewAngle)
	storeDebug(i, "9-deskew")

	// Construct the straight rectangle that contains our text (in the, now deskewed, image)
	var straightWidth, straightHeight int
	if math.Abs(rect.Angle) < 45 {
		straightWidth = rect.Width
		straightHeight = rect.Height
	} else {
		straightWidth = rect.Height
		straightHeight = rect.Width
	}

	straightRectPoints := gocv.NewPointsVectorFromPoints([][]goimage.Point{{
		{rect.Center.X - straightWidth/2, rect.Center.Y - straightHeight/2},
		{rect.Center.X + straightWidth/2, rect.Center.Y - straightHeight/2},
		{rect.Center.X + straightWidth/2, rect.Center.Y + straightHeight/2},
		{rect.Center.X - straightWidth/2, rect.Center.Y + straightHeight/2},
	}})
	// Draw the straight rectangle for debugging
	straightCopy := i.Clone()
	defer straightCopy.Close()
	gocv.DrawContours(&straightCopy, straightRectPoints, -1, color.RGBA{255, 255, 255, 255}, 3)
	storeDebug(&straightCopy, "10-deskewed-min-rectangle")

	// Now let's crop the rectangle
	straightRect := goimage.Rectangle{
		goimage.Point{
			max(rect.Center.X-straightWidth/2, 0),  // x0
			max(rect.Center.Y-straightHeight/2, 0), // y0
		},
		goimage.Point{
			min(rect.Center.X+straightWidth/2, i.Cols()),  // x1
			min(rect.Center.Y+straightHeight/2, i.Rows()), // y1
		},
	}

	croppedMat := i.Region(straightRect)
	// https://answers.opencv.org/question/22742/create-a-memory-continuous-cvmat-any-api-could-do-that/
	if !croppedMat.IsContinuous() {
		croppedMat = croppedMat.Clone()
	}
	storeDebug(&croppedMat, "11-cropped")
	*i = croppedMat
}

// calculateSkewAngle take the angle of the min area rectagle and return the
// angle to rotate the image in order to deskew the document.
// WarpPerspective does both in one step but it's a pain get the orientation right.
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
func rotateImg(i *gocv.Mat, center goimage.Point, angle float64) {
	size := i.Size()
	width := size[1]
	height := size[0]
	rMatrix := gocv.GetRotationMatrix2D(center, angle, 1.0)
	gocv.WarpAffineWithParams(*i, i, rMatrix, goimage.Point{X: width, Y: height}, gocv.InterpolationCubic, gocv.BorderReplicate, color.RGBA{0, 0, 0, 0})
}

// storeDebug writes an image to the filesystem if OOR_DEBUG env var is set
func storeDebug(i *gocv.Mat, filename string) {
	if os.Getenv("OOR_DEBUG") != "" {
		outPath := filepath.Join("tmp", filename+".jpg")
		if ok := gocv.IMWrite(outPath, *i); !ok {
			panic(fmt.Sprintf("Failed to write image: %s\n", outPath))
		}
	}
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
