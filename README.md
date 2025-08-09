# go-openai-rest-api

Go OpenAI REST API

## Chat (GPT-5 via SDK)

`POST /chat/text`

`POST /chat/text?stream=true` (newline-delimited JSON chunks: {"delta":"..."})

Request body

```
{
  prompt: "string"
}
```

Non-stream response:

```
{ "content": "..." }
```

Streaming: each line is JSON with a delta field until completion.

## VoiceGPT

Audio to Audio. uses ChatGPT, Text to speech, and Whisper

`/chat/voice` (audio -> text -> GPT-5 -> speech)

```
curl -X POST http://localhost:8080/chat/voice \
     -H "Content-Type: multipart/form-data" \
     -F "audio=@speech.mp3"
```

Text input to assistant voice response. use ChatGPT and Text to speech

`/chat/text_voice` (text -> GPT-5 -> speech)

```
curl -X POST http://localhost:8080/chat/text_audio \
     -H "Content-Type: application/json" \
     -d '{
           "prompt": "What is the meaning of the word Anagram?"
         }'
```

### Text to speech (SDK Audio.Speech)

`POST /tts`

```
curl -X POST http://localhost:8080/tts \
     -H "Content-Type: application/json" \
     -d '{
           "input": "Today is a wonderful day to build something people love!"
         }'
```

### Speech to text (SDK Audio.Transcriptions)

Transcribes an audio file using whisper-1.

```
curl -X POST http://localhost:8080/stt \
     -H "Content-Type: multipart/form-data" \
     -F "audio=@./speech.mp3"
```

### GET output file

Dowloads output file that is saved from voice reponses

```
curl http://localhost:8080/files/output.wav -o output.wav
```

## Images (DALL·E 3 via SDK)

`POST /image`

Request body

```
{
  prompt: "string"
}
```

Response

```
{
  "id": 1723223344,
  "prompt": "...possibly revised...",
  "url": "https://..."
}
```

## Environment

Set `OPENAI_API_KEY` for all endpoints.

## Notes

This service now uses the official OpenAI Go SDK v2 for chat, images, speech synthesis, and transcription.

## Docker Compose

Run the API with Docker Compose (uses `docker-compose.yml`):

1. Create an `.envrc` (or `.env`) file in the project root (loaded via `env_file`). Example:

```
API_KEY=sk-your-openai-key
BUCKET_NAME=your-bucket
AWS_ENDPOINT_URL_S3=https://s3.amazonaws.com
AWS_REGION=us-east-1
AWS_ACCESS_KEY_ID=your-access-key
AWS_SECRET_ACCESS_KEY=your-secret-key
```

2. Start the service:

```
docker compose up --build
```

3. Access the API at `http://localhost:8080`.
4. Generated audio files appear in `./data` (mounted into the container at `/data`).

To rebuild after code changes:

```
docker compose build app && docker compose up -d
```
