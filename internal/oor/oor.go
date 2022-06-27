// Package oor is the package that collects all the steps to go from a photograph
// to audio
package oor

import (
	"fmt"

	"github.com/jimmykarily/open-ocr-reader/internal/img"
	"github.com/jimmykarily/open-ocr-reader/internal/logger"
	"github.com/jimmykarily/open-ocr-reader/internal/ocr"
	"github.com/jimmykarily/open-ocr-reader/internal/process"
	"github.com/jimmykarily/open-ocr-reader/internal/tts"

	"github.com/pkg/errors"
)

type ParserDeps struct {
	Processor process.Processor
	OCR       ocr.OCR
	TTS       tts.TTS
}

// Parse takes all the steps needed to go from a photo of a book page to audio
func Parse(imgPath string, deps ParserDeps) error {
	logger := logger.New()

	textImg, err := img.New(imgPath)
	if err != nil {
		return errors.Wrap(err, "reading image file")
	}

	// TODO: It's easy to capture an image with external tools and pass it
	// to this program. Do we want to really deal with v4l and such?
	// E.g.
	// ffmpeg -f video4linux2 -s 640x480 -i /dev/video1 -ss 0:0:2 -frames 1 /tmp/out.jpg
	//
	// logger.Log("Taking a photo...")
	// img, err := capture.TakePhoto()
	// if err != nil {
	// 	return errors.Wrap(err, "taking an image")
	// }

	logger.Log("Processing the photo...")
	processedImg, err := deps.Processor.Process(textImg)
	if err != nil {
		return errors.Wrap(err, "processing the image")
	}

	// TODO: Detect blocks of text?

	// TODO: Make OCR an interface
	logger.Log("Running OCR on the photo...")
	text, err := deps.OCR.Parse(processedImg)
	if err != nil {
		return errors.Wrap(err, "running OCR on the image")
	}

	fmt.Printf("text = %+v\n", text)

	// TODO: Split in 2 steps? One to generate audio and one to play it?
	// Maybe the tts package can "stream" the audio, as in "play before the whole
	// text is parsed"?
	logger.Log("Running text to speech on the photo...")
	err = deps.TTS.Speak(text)
	if err != nil {
		return errors.Wrap(err, "running text to speech on the text")
	}

	return nil
}
