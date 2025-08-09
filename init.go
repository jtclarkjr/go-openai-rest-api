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

type MediaFile struct {
	ID       string `json:"id"`
	S3URL    string `json:"s3_url"`
	Filename string `json:"filename"`
}

type TTSInput struct {
	Input string `json:"input"`
}

type TTSRequest struct {
	Model string `json:"model"`
	Voice string `json:"voice"`
	Input string `json:"input"`
}

type WhisperResponse struct {
	Text string `json:"text"`
}

type ChatRequest struct {
	Prompt string `json:"prompt"`
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
