package main

import (
	"github.com/jimmykarily/open-ocr-reader/internal/logger"
	"github.com/jimmykarily/open-ocr-reader/internal/ocr"
	"github.com/jimmykarily/open-ocr-reader/internal/oor"
	"github.com/jimmykarily/open-ocr-reader/internal/process"
	"github.com/jimmykarily/open-ocr-reader/internal/tts"
	"github.com/jimmykarily/open-ocr-reader/internal/version"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           "oor",
	Short:         "Open OCR reader",
	Long:          `open-ocr-reader is tool that can read out loud paper books`,
	Version:       version.Version,
	SilenceErrors: true,
}

var parseCmd = &cobra.Command{
	Use:           "parse <image-file>",
	Short:         "parse an image file of text and produce audio in the command line",
	Long:          `This command can be as a cli to produce audio from an image file`,
	SilenceErrors: true,
	Args:          cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		logger := logger.New()
		//logger.Logf("args = %+v\n", args)

		parserDeps := oor.ParserDeps{
			Processor: process.NewDefaultProcessor(),
			OCR:       ocr.NewTesseractOCR(),
			TTS:       tts.NewDefaultTTS(),
		}

		if err := oor.Parse(args[0], parserDeps); err != nil {
			logger.Error(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(parseCmd)
}

func main() {
	rootCmd.Execute()
}
