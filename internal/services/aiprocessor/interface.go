package aiprocessor

import "context"

// WeightsInTextResult represents the result of finding weights in text.
// Fields:
// - Gram: The weight in grams found in the text.
// - Found: A boolean indicating if the weight was found.
// - Sure: A boolean indicating if the weight found is accurate.
type WeightsInTextResult struct {
	Gram  int64 `json:"weight"`
	Found bool  `json:"found"`
	Sure  bool  `json:"sure"`
}

// Service defines the interface for the AI processing service.
type Service interface {
	// FindWeightsInText finds the weight in grams of a given item in a text.
	// Parameters:
	// - ctx: The context for the request.
	// - text: The text to search within.
	// - search: The search string to look for in the text.
	// Returns: A WeightsInTextResult containing the result of the search and an error if any occurred.
	FindWeightsInText(ctx context.Context, text string, search string) (result WeightsInTextResult, err error)
}

// prompt is the instruction set for the AI model to find weight information in the text.
const prompt = `You need a found about weight  information in my text
RULE 1
return answer in json {"weight": int64 IN GRAMS, "found": bool true if you found information in my info, "sure": you sure about weight}
RULE 2
if the weight is not true - return sure as false
RULE 3
answer only as plain text json without markdown markup`
