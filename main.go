package main

import (
	"github.com/jimmykarily/open-ocr-reader/internal/logger"
	"github.com/jimmykarily/open-ocr-reader/internal/oor"
	"github.com/jimmykarily/open-ocr-reader/internal/process"
	"github.com/jimmykarily/open-ocr-reader/internal/version"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           "oor",
	Short:         "Open OCR reader",
	Long:          `open-ocr-reader is tool that can read out loud paper books`,
	Version:       version.Version,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		logger := logger.New()
		logger.Logf("args = %+v\n", args)

		processor := process.NewDefaultProcessor()

		if err := oor.Parse(args[0], processor); err != nil {
			logger.Error(err.Error())
		}
	},
}

func main() {
	rootCmd.Execute()
}
