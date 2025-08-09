package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/openai/openai-go/v2"
)

// ChatGPT voice assistant
func textVoiceChatController(w http.ResponseWriter, r *http.Request) {
	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get the completion from OpenAI using SDK
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Minute)
	defer cancel()
	completion, err := openAIClient.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{openai.UserMessage(req.Prompt)},
		Model:    openai.ChatModelGPT5,
	})
	if err != nil {
		http.Error(w, "Failed to get response from OpenAI: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if len(completion.Choices) == 0 {
		http.Error(w, "No completions found", http.StatusInternalServerError)
		return
	}
	completionText := completion.Choices[0].Message.Content

	// TTS via SDK
	speechResp, err := openAIClient.Audio.Speech.New(ctx, openai.AudioSpeechNewParams{
		Model: openai.SpeechModel("gpt-4o-mini-tts"),
		Input: completionText,
		Voice: openai.AudioSpeechNewParamsVoice("shimmer"),
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
	audioFilePath := "./data/output.wav"
	writeAudioDataToFile(w, speechResp.Body, audioFilePath)

	// Upload the file to S3
	uploadFileToS3(w, audioFilePath, "output-voice.wav")
}
