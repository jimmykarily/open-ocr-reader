// Package tts is reponsible for reading out loud the given text
package tts

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

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

	ttsIP := os.Getenv("TTS_IP")
	if ttsIP == "" {
		ttsIP = "127.0.0.1"
	}

	ttsPort := os.Getenv("TTS_PORT")
	if ttsPort == "" {
		ttsPort = "5002"
	}

	ttsVoice := os.Getenv("TTS_VOICE")
	if ttsVoice == "" {
		ttsVoice = "en-us/harvard-glow_tts"
	}

	data := url.Values{
		"voice":            {ttsVoice},
		"vocoder":          {"hifi_gan/universal_large"},
		"denoiserStrength": {"0.005"},
		"noiseScale":       {"0.667"},
		"lengthScale":      {"1"},
		"ssml":             {"on"},
		"text":             {text},
	}
	resp, err := http.Get("http://" + ttsIP + ":" + ttsPort + "/api/tts?" + data.Encode())
	if err != nil {
		logger.Log(err.Error())
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Log(err.Error())
		return err
	}
	_, err = f.Write(body)
	if err != nil {
		logger.Log(err.Error())
		return err
	}
	return nil
}
