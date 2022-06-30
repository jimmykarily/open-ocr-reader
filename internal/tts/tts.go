// Package tts is reponsible for reading out loud the given text
package tts

import (
	"io"
	"os"
	"os/exec"

	"github.com/jimmykarily/open-ocr-reader/internal/logger"
)

type TTS interface {
	Speak(text string) error
}

type DefaultTTS struct{}

func NewDefaultTTS() DefaultTTS {
	return DefaultTTS{}
}

func (t DefaultTTS) Speak(text string) error {
	logger := logger.New()
	f, err1 := os.Create("output.wav")
	if err1 != nil {
		return err1
	}
	defer f.Close()
	larynxCmd := exec.Command("larynx", "--raw-stream", text)
	audioStream, _ := larynxCmd.StdoutPipe()
	larynxCmd.Start()
	bytesRead := make([]byte, 1000)
	n, err := audioStream.Read(bytesRead)
	for ; err == nil; n, err = audioStream.Read(bytesRead) {
		_, err = f.Write(bytesRead[:n])
		if err != nil {
			return err
		}
	}
	if err != io.EOF {
		return err
	}
	// The generated stream is a stream raw 16-bit 22050Hz mono PCM audio to play it use cat output.wav | aplay -r 22050 -c 1 -f S16_LE
	logger.Log("End of the audio steam\n")
	larynxCmd.Wait()
	return nil
}
