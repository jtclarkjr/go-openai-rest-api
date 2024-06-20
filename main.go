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

	// Chat routes
	r.Route("/chat", func(r chi.Router) {

		// Text to text chat
		r.Post("/text", completionsController)

		// Audio to audio chat
		r.Post("/audio", voiceChatFromAudioController)

		// Text to chat assistant audio
		r.Post("/text_audio", textVoiceChatController)
	})

	// Text to speech
	r.Post("/tts", ttsController)

	// Speech to text
	r.Post("/stt", transcribeAudioController)

	// Text to image
	r.Post("/image", imageController)

	log.Println("Server starting on port 8080...")
	http.ListenAndServe(":8080", r)
}
