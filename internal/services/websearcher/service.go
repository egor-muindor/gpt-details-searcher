package websearcher

import (
	"context"
	"fmt"
	"sync"

	"dev.freespoke.com/brave-search"
	"golang.org/x/time/rate"
)

// service represents the web searcher service.
type service struct {
	braveClient brave.Brave   // Client for interacting with the Brave search API.
	mu          *sync.Mutex   // Mutex for synchronizing access to shared resources.
	limiter     *rate.Limiter // Rate limiter to control the rate of API requests.
}

// NewService creates a new instance of the web searcher service.
// Parameters:
// - braveToken: The token for authenticating with the Brave search API.
// - mu: A mutex for synchronizing access to shared resources.
// - limiter: A rate limiter to control the rate of API requests.
// Returns: A new instance of the Service interface and an error if any occurred.
func NewService(braveToken string, mu *sync.Mutex, limiter *rate.Limiter) (Service, error) {
	client, err := brave.New(braveToken)
	if err != nil {
		return nil, err
	}

	return &service{
		braveClient: client,
		mu:          mu,
		limiter:     limiter,
	}, nil
}

// Search searches the given query and returns a few result URLs.
// Parameters:
// - ctx: The context for the request.
// - query: The search query string.
// - limit: The maximum number of URLs to return (capped at 20).
// Returns: A slice of strings containing the result URLs and an error if any occurred.
func (s *service) Search(ctx context.Context, query string, limit uint8) (urls []string, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err = s.limiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit exceeded: %w", err)
	}
	if limit > 20 {
		limit = 20
	}

	result, err := s.braveClient.WebSearch(ctx, query, brave.WithCount(int(limit)))
	if err != nil {
		return nil, fmt.Errorf("failed to search: %w", err)
	}

	if result == nil {
		return nil, fmt.Errorf("no results found")
	}

	urls = make([]string, 0, len(result.Web.Results))
	for _, r := range result.Web.Results {
		urls = append(urls, r.URL)
	}

	return urls, nil
}
