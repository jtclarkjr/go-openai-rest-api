package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

// OpenAI API details
const (
	openAIURL        = "https://api.openai.com/v1/chat/completions"
	openAIImageURL   = "https://api.openai.com/v1/images/generations"
	openAITTSURL     = "https://api.openai.com/v1/audio/speech"
	openAIWhisperURL = "https://api.openai.com/v1/audio/transcriptions"
	apiKey           = ""
)

func main() {
	r := chi.NewRouter()

	// Text to text chat
	r.Post("/chat/text", completionsController)

	// Audio to audio chat
	r.Post("/chat/audio", voiceChatFromAudioController)

	// Text to chat assistant audio
	r.Post("/chat/text_audio", textVoiceChatController)

	// Text to speech
	r.Post("/tts", ttsController)

	// Text to image
	r.Post("/image", imageController)

	log.Println("Server starting on port 3000...")
	http.ListenAndServe(":3000", r)
}
