package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/jtclarkjr/router-go"
	"github.com/jtclarkjr/router-go/middleware"
)

func main() {
	r := router.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RateLimiter)
	r.Use(middleware.Throttle(100))
	r.Use(middleware.EnvVarChecker(
		"OPENAI_API_KEY",
		"BUCKET_NAME",
		"AWS_ENDPOINT_URL_S3",
		"AWS_REGION",
		"AWS_ACCESS_KEY_ID",
		"AWS_SECRET_ACCESS_KEY",
	))

	// Chat routes
	r.Route("/chat", func(r *router.Router) {

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

	r.Get("/audio/*", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/audio/")
		s3URL := fmt.Sprintf("%s/audio/%s", bucketUrl, path)
		http.Redirect(w, r, s3URL, http.StatusFound)
	})

	log.Println("Server starting on port 8080...")
	http.ListenAndServe(":8080", r)
}
