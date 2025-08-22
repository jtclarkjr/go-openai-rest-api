package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/openai/openai-go/v2"
)

// ttsController converts text to speech
// @Summary Convert text to speech
// @Description Convert text to speech using OpenAI's TTS model and save the audio file
// @Tags audio
// @Accept json
// @Produce application/octet-stream
// @Param request body TTSInput true "Text-to-speech request"
// @Success 200 {object} MediaFile "Successfully generated speech audio"
// @Failure 400 {object} ErrorResponse "Invalid request body"
// @Failure 500 {object} ErrorResponse "TTS generation failed"
// @Router /tts [post]
func ttsController(w http.ResponseWriter, r *http.Request) {

	var input TTSInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Minute)
	defer cancel()

	// Determine voice (default shimmer)
	voice := input.Voice
	if voice == "" {
		voice = "shimmer"
	}

	// Using SDK speech endpoint (returns http.Response)
	speechResp, err := openAIClient.Audio.Speech.New(ctx, openai.AudioSpeechNewParams{
		Model: openai.SpeechModel("gpt-4o-mini-tts"), // TTS capable model
		Input: input.Input,
		Voice: openai.AudioSpeechNewParamsVoice(voice),
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

	// Write audio data directly to timestamped file
	timestamp := time.Now().Unix()
	fileName := fmt.Sprintf("output-voice-%d.wav", timestamp)
	audioFilePath := fmt.Sprintf("./data/%s", fileName)
	writeAudioDataToFile(w, speechResp.Body, audioFilePath)

	// Upload the file to S3
	uploadFileToS3(w, audioFilePath, fileName)
}
