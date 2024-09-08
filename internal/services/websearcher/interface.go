package websearcher

import "context"

// Service defines the interface for the web searcher service.
type Service interface {
	// Search searches the given query and returns a few result URLs.
	// Parameters:
	// - ctx: The context for the request.
	// - query: The search query string.
	// - limit: The maximum number of URLs to return.
	// Returns: A slice of strings containing the result URLs and an error if any occurred.
	Search(ctx context.Context, query string, limit uint8) (urls []string, err error)
}
