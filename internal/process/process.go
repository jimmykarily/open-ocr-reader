// Package process is responsible for taking a raw image (as capture by
// the "capture" package and process it to make it suitable for OCR.
package process

import (
	"fmt"
	goimage "image"
	"image/color"
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
	storeDebug(&cvImg, "1-original")

	convertToGrayscale(&cvImg)
	storeDebug(&cvImg, "2-grayscale")

	deskew(&cvImg)

	_ = gocv.Threshold(cvImg, &cvImg, 140, 255, gocv.ThresholdBinary)
	storeDebug(&cvImg, "12-black-and-white")

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

	// gocv.GaussianBlur(tmpImg, &tmpImg, goimage.Point{}, 1, 1, gocv.BorderDefault)
	// storeDebug(&tmpImg, "3-after-gaussionblur")

	_ = gocv.Threshold(tmpImg, &tmpImg, 0, 255, gocv.ThresholdBinaryInv+gocv.ThresholdOtsu)
	storeDebug(&tmpImg, "4-after-threshold")

	kernel := gocv.GetStructuringElement(gocv.MorphEllipse, goimage.Point{5, 5})
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

	straightRectPoints := gocv.NewPointsVectorFromPoints([][]goimage.Point{{
		{rect.Center.X - rect.Width/2, rect.Center.Y - rect.Height/2},
		{rect.Center.X + rect.Width/2, rect.Center.Y - rect.Height/2},
		{rect.Center.X + rect.Width/2, rect.Center.Y + rect.Height/2},
		{rect.Center.X - rect.Width/2, rect.Center.Y + rect.Height/2},
	}})
	straightCopy := i.Clone()
	defer straightCopy.Close()
	gocv.DrawContours(&straightCopy, straightRectPoints, -1, color.RGBA{255, 255, 255, 255}, 3)
	storeDebug(&straightCopy, "10-deskewed-min-rectangle")

	fmt.Printf("straightRectPoints.ToPoints() = %+v\n", straightRectPoints.ToPoints())
	// Now let's crop the rectangle
	straightRect := goimage.Rectangle{
		goimage.Point{
			max(rect.Center.X-rect.Width/2.0, 0), // x0
			max(rect.Center.Y-rect.Height/2, 0),  // y0
		},
		goimage.Point{
			min(rect.Center.X+rect.Width/2, i.Cols()),  // x1
			min(rect.Center.Y+rect.Height/2, i.Rows()), // y1
		},
	}

	fmt.Printf("straightRect = %+v\n", straightRect)
	fmt.Printf("x = %+v\n", straightRect.Min.X)
	fmt.Printf("y = %+v\n", straightRect.Min.Y)
	fmt.Printf("width = %+v\n", straightRect.Size().X)
	fmt.Printf("height = %+v\n", straightRect.Size().Y)
	fmt.Printf("x + width = %+v\n", straightRect.Min.X+straightRect.Size().X)
	fmt.Printf("y + height = %+v\n", straightRect.Min.Y+straightRect.Size().Y)
	fmt.Printf("i.Cols() = %+v\n", i.Cols())
	fmt.Printf("i.Rows() = %+v\n", i.Rows())

	croppedMat := i.Region(straightRect)
	if !croppedMat.IsContinuous() {
		croppedMat = croppedMat.Clone()
	}
	storeDebug(&croppedMat, "11-cropped")
	*i = croppedMat
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
