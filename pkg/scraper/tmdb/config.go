package tmdb

import (
	"github.com/stashapp/stash/pkg/match"
	"github.com/stashapp/stash/pkg/models"
)

// Config holds the scraper configuration
type Config struct {
	ApiKey string `json:"api_key"`
}

// The factory function that creates our scraper
func init() {
	registerScraper()
}

func registerScraper() {
	models.RegisterScraperFactory("tmdb", "TMDB", createScraper)
}

func createScraper(config map[string]interface{}) models.Scraper {
	// Extract API key from config
	apiKey, _ := config["api_key"].(string)
	return NewScraper(apiKey)
}

// GetScraperOptions returns the options for configuring the scraper
func GetScraperOptions() map[string]models.ScraperOption {
	return map[string]models.ScraperOption{
		"api_key": {
			Type:     models.OptionTypeString,
			Required: true,
			Help:     "TMDB API Key",
		},
	}
}

// QueryURLBase returns the base URL for TMDB
func QueryURLBase() string {
	return "https://www.themoviedb.org"
}

// These functions help with URL matching
func (s *Scraper) ValidateURL(url string) bool {
	return match.URL(url, QueryURLBase())
}

func (s *Scraper) NormaliseURL(url string) string {
	return match.NormaliseURL(url)
}
