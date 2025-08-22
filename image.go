package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/openai/openai-go/v2"
)

// imageGenerate generates an image using OpenAI's DALL-E
// @Summary Generate an image from text prompt
// @Description Generate an image using OpenAI's DALL-E model based on a text prompt
// @Tags images
// @Accept json
// @Produce json
// @Param request body ImageRequest true "Image generation request"
// @Success 200 {object} ImageResponse "Successfully generated image"
// @Failure 400 {object} ErrorResponse "Invalid request body"
// @Failure 500 {object} ErrorResponse "Failed to generate image"
// @Router /image [post]
func imageGenerate(w http.ResponseWriter, r *http.Request) {
	var req ImageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Minute)
	defer cancel()

	imgResp, err := openAIClient.Images.Generate(ctx, openai.ImageGenerateParams{
		Model:  openai.ImageModelGPTImage1,
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
