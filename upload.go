package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
)

// uploadFileToS3 uploads a file to S3 and returns the file metadata
func uploadFileToS3(w http.ResponseWriter, audioFilePath string, s3Key string) {
	audioFile, err := os.Open(audioFilePath)
	if err != nil {
		http.Error(w, "Failed to open audio file", http.StatusInternalServerError)
		return
	}
	defer audioFile.Close()

	fileID := uuid.New().String()

	_, err = s3Service.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(s3Key),
		Body:   audioFile,
	})
	if err != nil {
		http.Error(w, "Failed to upload file to S3", http.StatusInternalServerError)
		return
	}

	s3URL := fmt.Sprintf("%s/%s", bucketUrl, s3Key)
	mediaFile := MediaFile{
		ID:       fileID,
		S3URL:    s3URL,
		Filename: s3Key,
	}

	// Save metadata locally
	filesStore.Store(fileID, mediaFile)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(mediaFile)
}

// writeAudioDataToFile writes audio data from the response body to a file
func writeAudioDataToFile(w http.ResponseWriter, body io.Reader, audioFilePath string) error {
	out, err := os.Create(audioFilePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	return nil
}
