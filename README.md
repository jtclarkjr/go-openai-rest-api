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
curl -X POST http://localhost:3000/chat/voice \
     -H "Content-Type: multipart/form-data" \
     -F "audio=@speech.mp3"
```

Text input to assistant voice response. use ChatGPT and Text to speech

`/chat/text_voice`

```
curl -X POST http://localhost:3000/chat/text_audio \
     -H "Content-Type: application/json" \
     -d '{
           "prompt": "What is Boyer Moore algorithm?"
         }'
```

### Text to speech

Text to speech. The text is converted in to speech in a audio file.

```
curl -X POST http://localhost:3000/tts \
     -H "Content-Type: application/json" \
     -d '{
           "input": "Today is a wonderful day to build something people love!"
         }'
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
