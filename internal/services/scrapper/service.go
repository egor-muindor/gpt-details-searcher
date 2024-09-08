package scrapper

import (
	"log/slog"
	"regexp"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-shiori/go-readability"
	"github.com/gocolly/colly/v2"
)

// service represents the scrapper service.
type service struct {
	collector *colly.Collector // Collector for scraping web pages.
	mu        *sync.Mutex      // Mutex for synchronizing access to shared resources.
}

// NewService creates a new instance of the scrapper service.
// Parameters:
// - collector: An instance of the colly.Collector to use for scraping.
// Returns: A new instance of the Service interface.
func NewService(collector *colly.Collector) Service {
	return &service{
		collector: collector,
		mu:        &sync.Mutex{},
	}
}

// ScrapeURL scrapes the given URLs and returns the text content without HTML tags.
// Parameters:
// - urls: A slice of strings containing the URLs to scrape.
// Returns: A map where the keys are the URLs and the values are the text content without HTML tags, and an error if any
// occurred.
func (s *service) ScrapeURL(urls []string) (text map[string]string, err error) {
	text = make(map[string]string, len(urls))
	for _, url := range urls {
		s.collector.OnHTML(
			"html", func(e *colly.HTMLElement) {
				s.mu.Lock()
				defer s.mu.Unlock()
				if e.DOM == nil {
					return
				}
				text[e.Request.URL.String()] = htmlToText(e.DOM)
			},
		)
		if err = s.collector.Visit(url); err != nil {
			slog.Error("Error on parsing", "error", err, "url", url)
		}
	}
	s.collector.Wait()
	return text, nil
}

// htmlToText converts the HTML content of a goquery.Selection to plain text.
// Parameters:
// - d: The goquery.Selection containing the HTML content.
// Returns: A string containing the plain text content.
func htmlToText(d *goquery.Selection) string {
	if d == nil {
		return ""
	}
	content := ""
	if readability.CheckDocument(d.Get(0)) {
		a, _ := readability.FromDocument(d.Get(0), nil)
		content = a.TextContent
	} else {
		d = d.Find("body")
		d.Find("style").Remove()
		d.Find("script").Remove()
		d.Find("meta").Remove()
		d.Find("link").Remove()
		d.Find("img").Remove()
		content = d.Text()
	}
	trimmed := strings.ReplaceAll(content, "\t", "")
	trimmed = strings.Join(strings.Fields(trimmed), " ")
	trimmed = strings.ReplaceAll(trimmed, "\n", "")
	trimmed = strings.ReplaceAll(trimmed, "\r", "")
	trimmed = strings.ReplaceAll(trimmed, "\t", "")
	trimmed = strings.ReplaceAll(trimmed, "  ", " ")
	re := regexp.MustCompile(`<[^>]*>`)
	trimmed = re.ReplaceAllString(trimmed, "")
	return trimmed
}
