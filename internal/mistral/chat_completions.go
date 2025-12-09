package mistral

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Model string

const (
	MistralLarge2512 Model = "mistral-large-2512"
)

type Role string

const (
	System    Role = "system"
	User      Role = "user"
	Assistant Role = "assistant"
)

type ResponseFormatType string

const (
	JSONObject ResponseFormatType = "json_object"
)

type Message struct {
	Role    Role   `json:"role"`
	Content string `json:"content"`
}

type ResponseFormat struct {
	Type ResponseFormatType `json:"type"`
}

type ChatCompletionRequest struct {
	Model          Model          `json:"model"`
	Messages       []Message      `json:"messages"`
	ResponseFormat ResponseFormat `json:"response_format"`
	SafePrompt     bool           `json:"safe_prompt"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	TotalTokens      int `json:"total_tokens"`
	CompletionTokens int `json:"completion_tokens"`
}

type Choice struct {
	Index        int     `json:"index"`
	FinishReason string  `json:"finish_reason"`
	Message      Message `json:"message"`
}

type ChatCompletionResponse struct {
	ID      string   `json:"id"`
	Created int64    `json:"created"`
	Model   Model    `json:"model"`
	Usage   Usage    `json:"usage"`
	Object  string   `json:"object"`
	Choices []Choice `json:"choices"`
}

type Client struct {
	BaseURL     string
	SecretToken string
}

func NewClient(secretToken string) *Client {
	return &Client{
		BaseURL:     "https://api.mistral.ai",
		SecretToken: secretToken,
	}
}

func (cl *Client) ChatComplete(req *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest(http.MethodPost, cl.BaseURL+"/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Authorization", "Bearer "+cl.SecretToken)
	res, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}

	var ccRes ChatCompletionResponse
	err = json.NewDecoder(res.Body).Decode(&ccRes)
	if err != nil {
		return nil, err
	}

	return &ccRes, nil
}
