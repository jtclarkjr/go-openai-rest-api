package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	openai "github.com/openai/openai-go/v2"
)

func completionsController(w http.ResponseWriter, r *http.Request) {
	var req ChatRequest
	stream := r.URL.Query().Get("stream")
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// ?stream="true" param; no param return without stream
	queryResult, err := strconv.ParseBool(stream)
	if err != nil {
		queryResult = false
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Minute)
	defer cancel()

	if !queryResult { // non-streaming simple completion
		completion, err := openAIClient.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.UserMessage(req.Prompt),
			},
			Model: openai.ChatModelGPT5,
		})
		if err != nil {
			http.Error(w, "OpenAI error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"content": completion.Choices[0].Message.Content,
		})
		return
	}

	streamResp := openAIClient.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{openai.UserMessage(req.Prompt)},
		Model:    openai.ChatModelGPT5,
	})
	defer streamResp.Close()

	w.Header().Set("Content-Type", "application/json")
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	acc := openai.ChatCompletionAccumulator{}
	for streamResp.Next() {
		chunk := streamResp.Current()
		acc.AddChunk(chunk)
		if len(chunk.Choices) > 0 {
			delta := chunk.Choices[0].Delta.Content
			if delta != "" {
				// Send SSE-like JSON lines (not official SSE)
				payload, _ := json.Marshal(map[string]string{"delta": delta})
				w.Write(payload)
				w.Write([]byte("\n"))
				flusher.Flush()
			}
		}
	}
	if streamResp.Err() != nil {
		log.Println("stream error:", streamResp.Err())
	}
}
