package moderation

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.socialmaps.org/api/internal/database"
	"golang.socialmaps.org/api/internal/mistral"
	"golang.socialmaps.org/api/internal/model"
	"golang.socialmaps.org/api/internal/must"
	"golang.socialmaps.org/api/internal/mytime"
)

func TestConsume(t *testing.T) {
	// Arrange
	ctx := t.Context()
	now := mytime.Now()

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
		must.Get(w.Write(
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
		))
	}))
	t.Cleanup(mistralSrv.Close)

	db := database.OpenInTest(t)
	qs := model.New(db)
	plc := must.Get(qs.CreatePlace(ctx, "Woo", 7.4192941, 43.7330475, model.OSMTypeNode, 12802966710, mytime.Now()))
	usr := must.Get(qs.CreateUser(ctx, now, 1, "Steve"))
	rvw := must.Get(qs.CreateReview(ctx, plc.ID, usr.ID, 4, new("great little cafe!"), mytime.Now(), mytime.Now()))

	mod := &MistralLarge2512v1{
		Client: mistral.Client{
			BaseURL:     mistralSrv.URL,
			SecretToken: "my-bearer-token",
		},
	}
	ch := make(chan model.Review, 1)
	ch <- rvw

	// Act
	consume(ctx, qs, mod, ch)

	// Assert
	dec := must.Get(qs.LoadLatestDecisionOfReview(ctx, rvw.ID))
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
