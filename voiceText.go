package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// WhisperResponse struct to parse Whisper API responses
type WhisperResponse struct {
	Text string `json:"text"`
}

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

func transcribeAudio(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		return "", err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return "", err
	}
	writer.WriteField("model", "whisper-1")
	writer.Close()

	req, err := http.NewRequest("POST", openAIWhisperURL, body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to transcribe audio: %s", string(bodyBytes))
	}

	var whisperResp WhisperResponse
	if err := json.NewDecoder(resp.Body).Decode(&whisperResp); err != nil {
		return "", err
	}

	return whisperResp.Text, nil
}
