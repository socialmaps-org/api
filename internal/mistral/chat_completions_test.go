package mistral

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChatCompletionRequest(t *testing.T) {
	// Arrange
	req := ChatCompletionRequest{
		Model: MistralLarge2512,
		Messages: []Message{
			{
				Role:    System,
				Content: "system-prompt",
			},
			{
				Role:    User,
				Content: "user-prompt",
			},
		},
		ResponseFormat: ResponseFormat{
			Type: JSONObject,
		},
		SafePrompt: true,
	}

	// Act
	b, err := json.MarshalIndent(req, "", "\t")

	// Assert
	require.NoError(t, err)
	require.Equal(
		t,
		`{
	"model": "mistral-large-2512",
	"messages": [
		{
			"role": "system",
			"content": "system-prompt"
		},
		{
			"role": "user",
			"content": "user-prompt"
		}
	],
	"response_format": {
		"type": "json_object"
	},
	"safe_prompt": true
}`,
		string(b),
	)
}

func TestChatCompletionResponse(t *testing.T) {
	// Arrange
	const doc = `{
	"id": "bdb27a8e1fa84f5d9128f2bb4b29f42e",
	"created": 1765137528,
	"model": "mistral-small-2506",
	"usage": {
		"prompt_tokens": 753,
		"total_tokens": 796,
		"completion_tokens": 43
	},
	"object": "chat.completion",
	"choices": [
		{
			"index": 0,
			"finish_reason": "stop",
			"message": {
				"role": "assistant",
				"tool_calls": null,
				"content": "{\n    \"approved\": true,\n    \"details\": \"The review is positive and does not contain any harmful, unethical, prejudiced, or negative content. It is safe for the platform.\"\n}"
			}
		}
	]
}`

	// Act
	var res ChatCompletionResponse
	err := json.Unmarshal([]byte(doc), &res)

	// Assert
	require.NoError(t, err)
	require.Equal(t, "bdb27a8e1fa84f5d9128f2bb4b29f42e", res.ID)
	require.Equal(t, int64(1765137528), res.Created)
	require.Equal(t, 753, res.Usage.PromptTokens)
	require.Equal(t, 796, res.Usage.TotalTokens)
	require.Equal(t, 43, res.Usage.CompletionTokens)
	require.Equal(t, "chat.completion", res.Object)
	require.Len(t, res.Choices, 1)
	require.Equal(t, 0, res.Choices[0].Index)
	require.Equal(t, "stop", res.Choices[0].FinishReason)
	require.Equal(t, Assistant, res.Choices[0].Message.Role)
	require.Equal(
		t,
		"{\n    \"approved\": true,\n    \"details\": \"The review is positive and does not contain any harmful, unethical, prejudiced, or negative content. It is safe for the platform.\"\n}",
		res.Choices[0].Message.Content,
	)
}
