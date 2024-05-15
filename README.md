# go-openai

Go OpenAI ChatGPT REST API

## ChatGPT

`GET /completions`

`GET /completions/stream`

Request body

```
{
  prompt: "string"
}
```

Response is same as [OpenAI Response](https://platform.openai.com/docs/api-reference/making-requests)

## Dalle

`GET /image`

Request body

```
{
  promp: "string"
}
```

Response

```
{
  id: number
  prompt: string
  url: string
}
```
