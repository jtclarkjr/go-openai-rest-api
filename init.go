package main

import (
	"log"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
)

// Bucket
var (
	filesStore   = sync.Map{}
	s3Session    *session.Session
	s3Service    *s3.S3
	bucketName   = os.Getenv("BUCKET_NAME")
	bucketUrl    = os.Getenv("AWS_ENDPOINT_URL_S3")
	bucketRegion = os.Getenv("AWS_REGION")
	bucketId     = os.Getenv("AWS_ACCESS_KEY_ID")
	bucketKey    = os.Getenv("AWS_SECRET_ACCESS_KEY")
)

var openAIClient openai.Client

// MediaFile represents a media file stored in S3
type MediaFile struct {
	ID       string `json:"id" example:"12345" description:"Unique identifier for the media file"`
	S3URL    string `json:"s3_url" example:"https://s3.amazonaws.com/bucket/file.wav" description:"S3 URL for the media file"`
	Filename string `json:"filename" example:"output.wav" description:"Original filename"`
}

// TTSInput represents input for text-to-speech conversion
type TTSInput struct {
	Input string `json:"input" binding:"required" example:"Hello world" description:"Text to convert to speech"`
	Voice string `json:"voice,omitempty" example:"alloy" description:"Voice to use for synthesis (alloy, ash, ballad, coral, echo, sage, shimmer, verse)"`
}

// TTSRequest represents a text-to-speech request
type TTSRequest struct {
	Model string `json:"model" example:"tts-1" description:"TTS model to use"`
	Voice string `json:"voice" example:"alloy" description:"Voice for synthesis"`
	Input string `json:"input" example:"Hello world" description:"Text to synthesize"`
}

// WhisperResponse represents the response from speech-to-text transcription
type WhisperResponse struct {
	Text string `json:"text" example:"Hello, how are you?" description:"Transcribed text from audio"`
}

// ChatRequest represents a chat completion request
type ChatRequest struct {
	Prompt string `json:"prompt" binding:"required" example:"What is the weather today?" description:"The prompt for the chat completion"`
	Voice  string `json:"voice,omitempty" example:"alloy" description:"Voice for audio response (if applicable)"`
}

// ChatResponse represents a chat completion response
type ChatResponse struct {
	Content string `json:"content" example:"The weather today is sunny and warm." description:"The AI-generated response"`
}

// StreamDelta represents a streaming response chunk
type StreamDelta struct {
	Delta string `json:"delta" example:"Hello" description:"Partial response content"`
}

// ImageRequest represents an image generation request
type ImageRequest struct {
	Prompt string `json:"prompt" binding:"required" example:"A beautiful sunset over mountains" description:"Description of the image to generate"`
}

// ImageResponse represents an image generation response
type ImageResponse struct {
	ID            int    `json:"id" example:"1723223344" description:"Timestamp ID for the generated image"`
	RevisedPrompt string `json:"prompt" example:"A beautiful sunset over mountains with vibrant colors" description:"The possibly revised prompt used for generation"`
	URL           string `json:"url" example:"https://oaidalleapiprodscus.blob.core.windows.net/private/..." description:"URL to download the generated image"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Message string `json:"message" example:"Invalid request body" description:"Error message"`
}

func init() {
	var err error

	s3Session, err = session.NewSession(&aws.Config{
		Region:   aws.String(bucketRegion),
		Endpoint: aws.String(bucketUrl),
		Credentials: credentials.NewStaticCredentials(
			bucketId,
			bucketKey,
			""),
		S3ForcePathStyle: aws.Bool(true),
	})
	if err != nil {
		log.Fatalf("failed to create S3 session: %v", err)
	}

	s3Service = s3.New(s3Session)

	// Initialize OpenAI client
	openai_key := os.Getenv("OPENAI_API_KEY")
	openAIClient = openai.NewClient(option.WithAPIKey(openai_key))
}
