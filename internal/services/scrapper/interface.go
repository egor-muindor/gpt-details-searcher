package scrapper

// Service defines the interface for the scrapper service.
type Service interface {
	// ScrapeURL scrapes the given URLs and returns the text content without HTML tags.
	// Parameters:
	// - urls: A slice of strings containing the URLs to scrape.
	// Returns: A map where the keys are the URLs and the values are the text content without HTML tags, and an error if
	// any occurred.
	ScrapeURL(urls []string) (text map[string]string, err error)
}
