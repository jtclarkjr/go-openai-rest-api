package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

// ImageRequest struct to parse incoming image generation requests
type ImageRequest struct {
	Prompt string `json:"prompt"`
}

// ImageResponse struct to format the response payload
type ImageResponse struct {
	ID            int    `json:"id"`
	RevisedPrompt string `json:"prompt"`
	URL           string `json:"url"`
}

func imageController(w http.ResponseWriter, r *http.Request) {
	apiKey := os.Getenv("API_KEY")
	var req ImageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	body, err := json.Marshal(map[string]interface{}{
		"model":      "dall-e-3",
		"prompt":     req.Prompt,
		"num_images": 1,
		"size":       "1024x1024",
	})

	if err != nil {
		http.Error(w, "Failed to encode request body", http.StatusInternalServerError)
		return
	}

	httpReq, err := http.NewRequest("POST", openAIImageURL, bytes.NewBuffer(body))
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
		http.Error(w, "Failed to generate image", resp.StatusCode)
		return
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		http.Error(w, "Failed to decode response from OpenAI", http.StatusInternalServerError)
		return
	}

	// Extract the necessary fields and create the custom response
	data := result["data"].([]interface{})[0].(map[string]interface{})
	customResponse := ImageResponse{
		ID:            int(result["created"].(float64)),
		RevisedPrompt: data["revised_prompt"].(string),
		URL:           data["url"].(string),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(customResponse); err != nil {
		http.Error(w, "Failed to encode response body", http.StatusInternalServerError)
		return
	}
}
