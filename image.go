package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/openai/openai-go/v2"
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
	var req ImageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Minute)
	defer cancel()

	imgResp, err := openAIClient.Images.Generate(ctx, openai.ImageGenerateParams{
		Model:  openai.ImageModelDallE3,
		Prompt: req.Prompt,
		Size:   openai.ImageGenerateParamsSize("1024x1024"),
		N:      openai.Int(1),
	})
	if err != nil {
		http.Error(w, "Failed to generate image: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if len(imgResp.Data) == 0 {
		http.Error(w, "No image returned", http.StatusInternalServerError)
		return
	}
	d := imgResp.Data[0]
	customResponse := ImageResponse{
		ID:            int(imgResp.Created),
		RevisedPrompt: d.RevisedPrompt,
		URL:           d.URL,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(customResponse); err != nil {
		http.Error(w, "Failed to encode response body", http.StatusInternalServerError)
		return
	}
}
