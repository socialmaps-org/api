package moderation

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"codeberg.org/socialmaps/api/internal/database"
	"codeberg.org/socialmaps/api/internal/mistral"
	"codeberg.org/socialmaps/api/internal/model"
	"github.com/stretchr/testify/require"
)

func TestConsume(t *testing.T) {
	// Arrange
	mistralSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/v1/chat/completions", r.URL.Path)

		auth := r.Header.Get("Authorization")
		require.Equal(t, "Bearer my-bearer-token", auth)

		var ccReq mistral.ChatCompletionRequest
		err := json.NewDecoder(r.Body).Decode(&ccReq)
		require.NoError(t, err)

		require.Equal(t, "mistral-large-2512", string(ccReq.Model))
		require.Len(t, ccReq.Messages, 2)
		require.Equal(t, "system", string(ccReq.Messages[0].Role))
		require.Equal(t, "user", string(ccReq.Messages[1].Role))
		require.Equal(t, "great little cafe!", ccReq.Messages[1].Content)
		require.Equal(t, "json_object", string(ccReq.ResponseFormat.Type))
		require.Equal(t, true, ccReq.SafePrompt)

		w.Header().Set("Content-Type", "application/json")
		w.Write(
			[]byte(`{
				"id": "bdb27a8e1fa84f5d9128f2bb4b29f42e",
				"created": 1765137528,
				"model": "mistral-large-2512",
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
			}`),
		)
	}))
	t.Cleanup(mistralSrv.Close)

	ctx := t.Context()

	db := database.Open(":memory:")
	plc := model.CreatePlace(ctx, db, "Izz Cafe", 51.8952597, -8.4715779, "node", 7095470096)
	usr := model.UpsertUser(ctx, db, "1", "Steve")
	rvw := model.CreateReview(ctx, db, plc.ID, usr.ID, true, "great little cafe!")

	mod := &MistralLarge2512v1{
		Client: mistral.Client{
			BaseURL:     mistralSrv.URL,
			SecretToken: "my-bearer-token",
		},
	}
	ch := make(chan *model.Review, 1)
	ch <- rvw

	// Act
	consume(ctx, db, mod, ch)

	// Assert
	dec := model.LoadLatestDecisionOfReview(ctx, db, rvw.ID)
	require.NotNil(t, dec)
	require.Equal(t, rvw.ID, dec.ReviewID)
	require.Equal(t, true, dec.Approved)
	require.Equal(
		t,
		"The review is positive and does not contain any harmful, unethical, prejudiced, or negative content. It is safe for the platform.",
		dec.Details,
	)
	require.Equal(t, "mistral-large-2512-v1", dec.Moderator)
}
