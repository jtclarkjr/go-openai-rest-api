package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/openai/openai-go/v2"
)

func voiceChatFromAudioController(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("audio")
	if err != nil {
		http.Error(w, "Failed to get audio file from request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	tempFile, err := os.CreateTemp("", "audio-*.mp3")
	if err != nil {
		http.Error(w, "Failed to create temporary file", http.StatusInternalServerError)
		return
	}
	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	_, err = io.Copy(tempFile, file)
	if err != nil {
		http.Error(w, "Failed to save audio file", http.StatusInternalServerError)
		return
	}

	text, err := transcribeAudio(tempFile.Name())
	if err != nil {
		http.Error(w, "Failed to transcribe audio", http.StatusInternalServerError)
		return
	}

	// Create a new request with the transcribed text
	req := ChatRequest{Prompt: text}
	jsonReq, err := json.Marshal(req)
	if err != nil {
		http.Error(w, "Failed to encode request body", http.StatusInternalServerError)
		return
	}

	r.Body = io.NopCloser(bytes.NewBuffer(jsonReq))
	textVoiceChatController(w, r)
}

func transcribeAudioController(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("audio")
	if err != nil {
		http.Error(w, "No file found", http.StatusBadRequest)
		return
	}
	defer file.Close()

	log.Printf("Uploaded file: %s\n", handler.Filename)

	tempFile, err := os.CreateTemp("", "audio-*.mp3")
	if err != nil {
		http.Error(w, "Could not create temp file", http.StatusInternalServerError)
		return
	}
	defer os.Remove(tempFile.Name())

	if _, err := io.Copy(tempFile, file); err != nil {
		http.Error(w, "Could not save temp file", http.StatusInternalServerError)
		return
	}

	transcribedText, err := transcribeAudio(tempFile.Name())
	if err != nil {
		http.Error(w, fmt.Sprintf("Error transcribing audio: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"text": transcribedText})
}

func transcribeAudio(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// openai.File helper to set filename & type
	resp, err := openAIClient.Audio.Transcriptions.New(ctx, openai.AudioTranscriptionNewParams{
		Model: "whisper-1",
		File:  openai.File(f, filepath.Base(filePath), "audio/mpeg"),
	})
	if err != nil {
		return "", err
	}
	return resp.Text, nil
}
