# go-openai-rest-api

Go OpenAI REST API

## ChatGPT

`GET /chat/text`

`GET /chat/text?stream="true"`

Request body

```
{
  prompt: "string"
}
```

Response is same as [OpenAI Response](https://platform.openai.com/docs/api-reference/making-requests)

## VoiceGPT

Audio to Audio. uses ChatGPT, Text to speech, and Whisper

`/chat/voice`

```
curl -X POST http://localhost:8080/chat/voice \
     -H "Content-Type: multipart/form-data" \
     -F "audio=@speech.mp3"
```

Text input to assistant voice response. use ChatGPT and Text to speech

`/chat/text_voice`

```
curl -X POST http://localhost:8080/chat/text_audio \
     -H "Content-Type: application/json" \
     -d '{
           "prompt": "What is the meaning of the word Anagram?"
         }'
```

### Text to speech

Text to speech. The text is converted in to speech in a audio file.

`/tts`

```
curl -X POST http://localhost:8080/tts \
     -H "Content-Type: application/json" \
     -d '{
           "input": "Today is a wonderful day to build something people love!"
         }'
```

### Speech to text

Speech to text by transcribing an audio file

```
curl -X POST http://localhost:8080/stt \
     -H "Content-Type: multipart/form-data" \
     -F "audio=@./speech.mp3"
```

### GET output file

Dowloads output file that is saved from voice reponses

```
curl http://your-app-name.fly.dev/files/output.wav -o output.wav
```

## Dalle

`GET /image`

Request body

```
{
  prompt: "string"
}
```

Response

```
{
  id: int
  prompt: string
  url: string
}
```
