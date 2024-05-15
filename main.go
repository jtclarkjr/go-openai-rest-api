package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

// OpenAI API details
const (
	openAIURL      = "https://api.openai.com/v1/chat/completions"
	openAIImageURL = "https://api.openai.com/v1/images/generations"
	apiKey         = "" // use os.Getenv("OPENAI_KEY")

)

func main() {
	r := chi.NewRouter()
	r.Post("/completions", completionsController)
	r.Post("/image", imageController)

	log.Println("Server starting on port 3000...")
	http.ListenAndServe(":3000", r)
}
