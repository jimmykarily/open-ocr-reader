// Package oor is the package that collects all the steps to go from a photograph
// to audio
package oor

import (
	"github.com/jimmykarily/open-ocr-reader/internal/capture"
	"github.com/jimmykarily/open-ocr-reader/internal/logger"
	"github.com/jimmykarily/open-ocr-reader/internal/ocr"
	"github.com/jimmykarily/open-ocr-reader/internal/process"
	"github.com/jimmykarily/open-ocr-reader/internal/tts"

	"github.com/pkg/errors"
)

// Parse takes all the steps needed to go from a photograph to audio
func Parse() error {
	logger := logger.New()

	logger.Log("Taking a photo...")
	img, err := capture.TakePhoto()
	if err != nil {
		return errors.Wrap(err, "taking an image")
	}

	logger.Log("Processing the photo...")
	improvedImg, err := process.Process(img)
	if err != nil {
		return errors.Wrap(err, "processing the image")
	}

	// TODO: Detect blocks of text?

	logger.Log("Running OCR on the photo...")
	text, err := ocr.Parse(improvedImg)
	if err != nil {
		return errors.Wrap(err, "running OCR on the image")
	}

	// TODO: Split in 2 steps? One to generate audio and one to play it?
	// Maybe the tts package can "stream" the audio, as in "play before the whole
	// text is parsed"?
	logger.Log("Running text to speech on the photo...")
	err = tts.Speak(text)
	if err != nil {
		return errors.Wrap(err, "running text to speech on the text")
	}

	return nil
}
