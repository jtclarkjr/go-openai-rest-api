package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
)

// ChatGPT voice assistant
func textVoiceChatController(w http.ResponseWriter, r *http.Request) {
	apiKey := os.Getenv("API_KEY")
	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get the completion from OpenAI
	body, err := json.Marshal(map[string]interface{}{
		"model":    "gpt-4o",
		"messages": []map[string]string{{"role": "user", "content": req.Prompt}},
	})
	if err != nil {
		http.Error(w, "Failed to encode request body", http.StatusInternalServerError)
		return
	}

	httpReq, err := http.NewRequest("POST", openAIURL, bytes.NewBuffer(body))
	if err != nil {
		http.Error(w, "Failed to create request to OpenAI", http.StatusInternalServerError)
		return
	}
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		http.Error(w, "Failed to get response from OpenAI", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		http.Error(w, string(bodyBytes), resp.StatusCode)
		return
	}

	var completionResponse struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&completionResponse); err != nil {
		http.Error(w, "Failed to decode OpenAI response", http.StatusInternalServerError)
		return
	}

	if len(completionResponse.Choices) == 0 {
		http.Error(w, "No completions found", http.StatusInternalServerError)
		return
	}

	completionText := completionResponse.Choices[0].Message.Content

	// Use the completion text as input for TTS
	ttsReq := TTSRequest{
		Model: "tts-1",
		Voice: "shimmer",
		Input: completionText,
	}

	ttsBody, err := json.Marshal(ttsReq)
	if err != nil {
		http.Error(w, "Failed to encode TTS request body", http.StatusInternalServerError)
		return
	}

	ttsReqHTTP, err := http.NewRequest("POST", openAITTSURL, bytes.NewBuffer(ttsBody))
	if err != nil {
		http.Error(w, "Failed to create request to OpenAI TTS", http.StatusInternalServerError)
		return
	}
	ttsReqHTTP.Header.Set("Content-Type", "application/json")
	ttsReqHTTP.Header.Set("Authorization", "Bearer "+apiKey)

	ttsResp, err := client.Do(ttsReqHTTP)
	if err != nil {
		http.Error(w, "Failed to get response from OpenAI TTS", http.StatusInternalServerError)
		return
	}
	defer ttsResp.Body.Close()

	if ttsResp.StatusCode != http.StatusOK {
		ttsBodyBytes, _ := io.ReadAll(ttsResp.Body)
		http.Error(w, string(ttsBodyBytes), ttsResp.StatusCode)
		return
	}

	// Write audio data directly to file
	audioFilePath := "./data/output.wav"
	out, err := os.Create(audioFilePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, ttsResp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send response with file path
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"file": audioFilePath})
}
