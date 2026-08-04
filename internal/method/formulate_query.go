package method

import (
	"context"
	_ "embed"
	"errors"
	"log/slog"

	"github.com/openai/openai-go/v3"
	"github.com/theory/sqljson/path"

	"golang.socialmaps.org/api/internal/resource"
)

type FormulateQuery struct {
	Common
	openaiClient openai.Client
}

type formulateQueryArgs struct {
	Query string `query:"query" doc:"A user query in natural language for querying **Place**s." example:"Historical tourist attractions"`
}

//go:embed system-prompt.md
var systemPrompt string

func (m *FormulateQuery) Execute(ctx context.Context, args *formulateQueryArgs) (*Response[resource.Query], error) {
	comp, err := m.openaiClient.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemPrompt),
			openai.UserMessage("Italian or French restaurants with wheelchair access"),
			openai.AssistantMessage(
				`$.amenity == "restaurant" && ($.cuisine == "italian" || $.cuisine == "french") && $.wheelchair == "yes"`,
			),
			openai.UserMessage(args.Query),
		},
		Model:               "gpt-oss-120b",
		MaxCompletionTokens: openai.Int(500),
	})
	if err != nil {
		return nil, err
	}

	slog.InfoContext(ctx, "CANONICAL-OPENAI-USAGE-LINE",
		"use_case", "formulate_query",
		"completion_tokens", comp.Usage.CompletionTokens,
	)

	choice := comp.Choices[0]
	if choice.FinishReason == string(openai.CompletionChoiceFinishReasonLength) {
		return nil, errors.New("thought too hard")
	}

	predicate := choice.Message.Content

	_, err = path.Parse(predicate)
	if err != nil {
		return nil, err
	}

	return &Response[resource.Query]{Body: resource.Query{Predicate: predicate}}, nil
}
