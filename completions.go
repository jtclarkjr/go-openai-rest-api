package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

// ChatRequest struct to parse incoming requests
type ChatRequest struct {
	Prompt string `json:"prompt"`
}

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

	body, err := json.Marshal(map[string]interface{}{
		"model":    "gpt-4o",
		"messages": []map[string]string{{"role": "user", "content": req.Prompt}},
		"stream":   queryResult,
	})

	log.Println(stream)
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

	// Set header for chunked transfer
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	reader := bufio.NewReader(resp.Body)
	buffer := make([]byte, 1024)
	for {
		n, err := reader.Read(buffer)
		if n > 0 {
			_, writeErr := w.Write(buffer[:n])
			if writeErr != nil {
				log.Println("Error writing to response:", writeErr)
				return
			}
			w.(http.Flusher).Flush()
		}
		if err != nil {
			break
		}
	}
}
