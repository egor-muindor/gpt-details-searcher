package aiprocessor

import (
	"context"
	"encoding/json"

	"github.com/sashabaranov/go-openai"
)

// service represents the AI processing service.
type service struct {
	client *openai.Client // OpenAI client for making API requests.
	model  string         // Model identifier for the OpenAI model to use.
}

// NewService creates a new instance of the AI processing service.
// Parameters:
// - client: An instance of the OpenAI client.
// - model: The model identifier to use for the OpenAI API.
// Returns: A new instance of the Service interface.
func NewService(client *openai.Client, model string) Service {
	return &service{client: client, model: model}
}

// FindWeightsInText finds weights in the given text based on the search string.
// Parameters:
// - ctx: The context for the request.
// - text: The text to search within.
// - search: The search string to look for in the text.
// Returns: A WeightsInTextResult containing the result of the search and an error if any occurred.
func (s *service) FindWeightsInText(ctx context.Context, text string, search string) (
	result WeightsInTextResult,
	err error,
) {
	response, err := s.client.CreateChatCompletion(
		ctx, openai.ChatCompletionRequest{
			Model: s.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    "system",
					Content: prompt,
				}, {
					Role:    "user",
					Content: search + "\n" + text,
				},
			},
		},
	)
	if err != nil {
		return WeightsInTextResult{}, err
	}
	for _, choice := range response.Choices {
		err = json.Unmarshal([]byte(choice.Message.Content), &result)
		if err != nil {
			return WeightsInTextResult{}, err
		}

		return result, nil
	}
	return WeightsInTextResult{
		Gram:  0,
		Found: false,
		Sure:  false,
	}, nil
}
