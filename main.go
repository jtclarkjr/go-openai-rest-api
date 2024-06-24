package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

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

	// serves from server
	r.Get("/files/*", func(w http.ResponseWriter, r *http.Request) {
		fs := http.FileServer(http.Dir("/data"))
		http.StripPrefix("/files/", fs).ServeHTTP(w, r)
	})

	// serves from storage bucket
	r.Get("/audio/*", func(w http.ResponseWriter, r *http.Request) {
		path := chi.URLParam(r, "*")
		s3URL := fmt.Sprintf("%s/audio/%s", bucketUrl, path)
		http.Redirect(w, r, s3URL, http.StatusFound)
	})

	log.Println("Server starting on port 8080...")
	http.ListenAndServe(":8080", r)
}
