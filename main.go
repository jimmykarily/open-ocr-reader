package main

import (
	"net"
	"net/http"
	"os"

	"github.com/jimmykarily/open-ocr-reader/controllers"
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

var serverCmd = &cobra.Command{
	Use:           "server",
	Short:         "start the web server",
	Long:          `This command start a web server that hosts the web interface of this application. It starts on a random port or "PORT" environment variable value, if set`,
	SilenceErrors: true,
	Args:          cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		logger := logger.New()
		logger.Log("Starting the server")

		mux := http.NewServeMux()
		fileServer := http.FileServer(http.Dir("./static/"))
		mux.Handle("/static/", http.StripPrefix("/static", fileServer))

		mux.HandleFunc("/", controllers.Home)

		// https://gist.github.com/xcsrz/538e291d12be6ee9a8c7
		var port string
		if port = os.Getenv("PORT"); port == "" {
			port = "0"
		}
		listener, err := net.Listen("tcp", "0.0.0.0:"+port)
		if err != nil {
			logger.Errorf("starting the server %s", err.Error())
		}
		logger.Logf("listening on %s", listener.Addr().String())
		http.Serve(listener, mux)
	},
}

func init() {
	rootCmd.AddCommand(parseCmd)
	rootCmd.AddCommand(serverCmd)
}

func main() {
	rootCmd.Execute()
}
