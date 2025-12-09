package moderation

import (
	_ "embed"
	"encoding/json"

	"codeberg.org/socialmaps/api/internal/mistral"
)

//go:embed data/system_prompt.md
var systemPrompt string

type MistralLarge2512v1 struct {
	Client mistral.Client
}

func NewMistralLarge2512v1(client *mistral.Client) *MistralLarge2512v1 {
	return &MistralLarge2512v1{
		Client: *client,
	}
}

type content struct {
	Approved bool   `json:"approved"`
	Details  string `json:"details"`
}

func (mod *MistralLarge2512v1) ID() string {
	return "mistral-large-2512-v1"
}

func (mod *MistralLarge2512v1) Moderate(review string) (*Decision, error) {
	res, err := mod.Client.ChatComplete(&mistral.ChatCompletionRequest{
		Model: mistral.MistralLarge2512,
		Messages: []mistral.Message{
			{
				Role:    mistral.System,
				Content: systemPrompt,
			},
			{
				Role:    mistral.User,
				Content: CleanUp(review),
			},
		},
		ResponseFormat: mistral.ResponseFormat{
			Type: mistral.JSONObject,
		},
		SafePrompt: true,
	})
	if err != nil {
		return nil, err
	}

	var content content
	err = json.Unmarshal([]byte(res.Choices[0].Message.Content), &content)
	if err != nil {
		return nil, err
	}

	return &Decision{Approved: content.Approved, Details: content.Details}, nil
}
