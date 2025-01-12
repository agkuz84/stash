package tmdb

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v2"
)

type config struct {
	APIKey string `yaml:"api_key"`
}

func loadConfig(t *testing.T) string {
	// Get the project root directory
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	// Navigate up to the root of the project
	rootDir := filepath.Join(pwd, "..", "..", "..")
	configPath := filepath.Join(rootDir, "scraper", "tmdb", "tmdb.yml")

	// Read the config file
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	// Parse the YAML
	var cfg config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("Failed to parse config: %v", err)
	}

	if cfg.APIKey == "" {
		t.Fatal("No API key found in config file")
	}

	return cfg.APIKey
}

func TestTMDBScraper(t *testing.T) {
	apiKey := loadConfig(t)
	scraper := NewScraper(apiKey)

	// Test movie scraping
	t.Run("Test Movie Search", func(t *testing.T) {
		ctx := context.Background()
		result, err := scraper.ScrapeMovie(ctx, "The Matrix")
		if err != nil {
			t.Fatalf("Failed to scrape movie: %v", err)
		}

		if result == nil {
			t.Fatal("No result returned")
		}

		if result.Name == nil {
			t.Fatal("Movie name is nil")
		}

		if *result.Name != "The Matrix" {
			t.Errorf("Expected 'The Matrix', got '%s'", *result.Name)
		}

		// Check if we got a synopsis
		if result.Synopsis == nil || *result.Synopsis == "" {
			t.Error("No synopsis returned")
		}

		// Check if we got a release date
		if result.Date == nil || *result.Date == "" {
			t.Error("No release date returned")
		}
	})

	// Test TV show scraping
	t.Run("Test TV Show Search", func(t *testing.T) {
		ctx := context.Background()
		result, err := scraper.ScrapeTV(ctx, "Breaking Bad")
		if err != nil {
			t.Fatalf("Failed to scrape TV show: %v", err)
		}

		if result == nil {
			t.Fatal("No result returned")
		}

		if result.Name == nil {
			t.Fatal("Show name is nil")
		}

		if *result.Name != "Breaking Bad" {
			t.Errorf("Expected 'Breaking Bad', got '%s'", *result.Name)
		}

		// Check if we got a synopsis
		if result.Synopsis == nil || *result.Synopsis == "" {
			t.Error("No synopsis returned")
		}
	})
}
