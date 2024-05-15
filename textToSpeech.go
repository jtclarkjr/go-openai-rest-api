package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
)

type TTSInput struct {
	Input string `json:"input"`
}

type TTSRequest struct {
	Model string `json:"model"`
	Voice string `json:"voice"`
	Input string `json:"input"`
}

func ttsController(w http.ResponseWriter, r *http.Request) {
	var input TTSInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req := TTSRequest{
		Model: "tts-1",
		Voice: "shimmer",
		Input: input.Input,
	}

	// Prepare the request to OpenAI API
	body, err := json.Marshal(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	client := &http.Client{}
	openAIReq, err := http.NewRequest("POST", openAITTSURL, bytes.NewBuffer(body))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	openAIReq.Header.Set("Content-Type", "application/json")
	openAIReq.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := client.Do(openAIReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		http.Error(w, string(bodyBytes), resp.StatusCode)
		return
	}

	// Write audio data directly to file
	audioFilePath := "./speech.mp3"
	out, err := os.Create(audioFilePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send response with file path
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"file": audioFilePath})
}
