package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/openai/openai-go/v2"
)

func ttsController(w http.ResponseWriter, r *http.Request) {

	var input TTSInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Minute)
	defer cancel()

	// Using SDK speech endpoint (returns http.Response)
	speechResp, err := openAIClient.Audio.Speech.New(ctx, openai.AudioSpeechNewParams{
		Model: openai.SpeechModel("gpt-4o-mini-tts"), // TTS capable model
		Input: input.Input,
		Voice: openai.AudioSpeechNewParamsVoice("shimmer"),
		// default format (mp3 or wav depending on API defaults)
	})
	if err != nil {
		http.Error(w, "TTS error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer speechResp.Body.Close()
	if speechResp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(speechResp.Body)
		http.Error(w, string(b), speechResp.StatusCode)
		return
	}

	// Write audio data directly to file
	audioFilePath := "./data/output.wav"
	writeAudioDataToFile(w, speechResp.Body, audioFilePath)

	// Upload the file to S3
	uploadFileToS3(w, audioFilePath, "output-voice.wav")
}
