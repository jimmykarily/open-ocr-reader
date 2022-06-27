// Package tts is reponsible for reading out loud the given text
package tts

type TTS interface {
	Speak(text string) error
}

type DefaultTTS struct{}

func NewDefaultTTS() DefaultTTS {
	return DefaultTTS{}
}

func (t DefaultTTS) Speak(text string) error {
	return nil
}
