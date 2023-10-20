package promptCollection

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"my_chatbot/models"
	"net/http"

	openai "github.com/sashabaranov/go-openai"
)

type Openai struct {
	apiURL     string
	authBearer string
}

func (o *Openai) SetAPIURL(url string) {
	o.apiURL = url
}

func (o *Openai) SetAuthBearer(bearer string) {
	o.authBearer = bearer
}

// TODO: Initialise these fields within Openai struct and remove these constants
const (
	apiURL     = "https://api.openai.com/v1/chat/completions"
	authBearer = "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"
)

func (o *Openai) CallOpenAIRest(prompt string) string {
	data := models.RequestPayload{
		Model:     "gpt-3.5-turbo",
		Prompt:    prompt,
		MaxTokens: 150,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return "Error occurred"
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+authBearer)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "Error occurred while calling OpenAI Completion API"
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var response models.ResponsePayload
	json.Unmarshal(body, &response)

	if len(response.Choices) > 0 {
		return response.Choices[0].Text
	}
	return "No reply generated"
}

func (o *Openai) CallOpenAIUsingPackage(prompt string) string {
	client := openai.NewClient(authBearer)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)
	if err != nil {
		return "Error occurred while calling OpenAI Completion"
	}

	if len(resp.Choices) > 0 {
		return resp.Choices[0].Message.Content
	}
	return "No reply generated"
}
