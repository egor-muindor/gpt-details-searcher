package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/gocolly/colly/v2"
	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
	"golang.org/x/time/rate"

	"github.com/egor-muindor/gpt-details-searcher/internal/services/aiprocessor"
	"github.com/egor-muindor/gpt-details-searcher/internal/services/scrapper"
	"github.com/egor-muindor/gpt-details-searcher/internal/services/websearcher"
)

const (
	braveTokenEnv     = "BRAVE_TOKEN"
	openAIEndpointEnv = "OPENAI_ENDPOINT"
	openAITokenEnv    = "OPENAI_TOKEN"
	aiModelEnv        = "AI_MODEL"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load(".env")
	if err != nil {
		panic("can't load environment")
	}

	// Retrieve BRAVE_TOKEN from environment variables
	braveToken, ok := os.LookupEnv(braveTokenEnv)
	if !ok {
		panic("can't load BRAVE_TOKEN from env")
	}

	// Retrieve OPENAI_ENDPOINT from environment variables
	openAIEndpoint, ok := os.LookupEnv(openAIEndpointEnv)
	if !ok {
		panic("can't load OPENAI_ENDPOINT from env")
	}

	// Retrieve OPENAI_TOKEN from environment variables
	openAIToken, ok := os.LookupEnv(openAITokenEnv)
	if !ok {
		panic("can't load OPENAI_TOKEN from env")
	}

	// Retrieve AI_MODEL from environment variables
	aiModel, ok := os.LookupEnv(aiModelEnv)
	if !ok {
		panic("can't load AI_MODEL from env")
	}

	// Initialize synchronization mutex and rate limiter
	mu := &sync.Mutex{}
	limiter := rate.NewLimiter(1, 1) // 1 request per second

	// Create a new instance of the web searcher service
	webSearchService, err := websearcher.NewService(braveToken, mu, limiter)
	if err != nil {
		panic(err)
	}

	// Create a new instance of the scrapper service
	collector := colly.NewCollector(
		colly.Async(true),
		colly.AllowURLRevisit(),
		colly.MaxDepth(1),
	)
	scrapperService := scrapper.NewService(collector)

	// Initialize OpenAI client
	openaiConfig := openai.DefaultConfig(openAIToken)
	openaiConfig.BaseURL = openAIEndpoint
	openaiClient := openai.NewClientWithConfig(openaiConfig)

	// Initialize AI processor service
	textProcessor := aiprocessor.NewService(openaiClient, aiModel)

	// =============== test area below VVV =================

	SearchQuery := "weight of nike air force 1 in box"

	ctx := context.Background()
	result, err := webSearchService.Search(ctx, SearchQuery, 5)
	if err != nil {
		panic(err)
	}

	// urls for testing
	// 	return []string{
	//		"https://www.xtremeloaded.com/14025/nike-air-force-one-here-is-everything-you-need-to-know",
	//		"https://www.reddit.com/r/FashionReps/comments/cahymu/can_you_tell_me_the_weight_of_af1_lows/?rdt=61342",
	//		"https://runrepeat.com/nike-air-force-1-07",
	//		"https://wearablyweird.com/how-much-do-air-force-ones-weigh-are-they-heavy-to-wear/",
	//	}, nil

	scrapedText, err := scrapperService.ScrapeURL(result)
	if err != nil {
		panic(err)
	}

	for url, text := range scrapedText {
		result, err := textProcessor.FindWeightsInText(ctx, text, SearchQuery)
		if err != nil {
			fmt.Println("Error on processing text:", err, "text:", result)
			panic(err)
		}
		if result.Found {
			fmt.Println("Found weight:", result.Gram, "grams")
			j, _ := json.Marshal(result)
			fmt.Printf("json: %v\n", string(j))
			break
		} else {
			fmt.Println("Weight not found in url: " + url)
		}
	}

}
