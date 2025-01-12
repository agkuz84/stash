package tmdb

import (
	"context"
	"fmt"
	"strconv"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

type Scraper struct {
	client *Client
}

func NewScraper(apiKey string) *Scraper {
	return &Scraper{
		client: NewClient(apiKey),
	}
}

// ScrapeMovie tries to find movie information based on the name
func (s *Scraper) ScrapeMovie(ctx context.Context, name string) (*models.ScrapedMovie, error) {
	logger.Infof("Searching TMDB for movie: %s", name)

	results, err := s.client.SearchMovie(name)
	if err != nil {
		return nil, fmt.Errorf("TMDB search failed: %v", err)
	}

	if len(results.Results) == 0 {
		return nil, fmt.Errorf("no movies found matching: %s", name)
	}

	movie := results.Results[0]
	details, err := s.client.GetMovieDetails(movie.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get movie details: %v", err)
	}

	// Convert to ScrapedMovie
	scraped := &models.ScrapedMovie{}

	// Set basic information
	scraped.Name = &movie.Title
	scraped.Date = &movie.ReleaseDate
	scraped.Synopsis = &movie.Overview

	// Handle Rating
	if movie.VoteAverage > 0 {
		rating := strconv.FormatFloat(movie.VoteAverage, 'f', 2, 64)
		scraped.Rating = &rating
	}

	// Handle image
	if movie.PosterPath != "" {
		posterURL := s.client.GetImageURL(movie.PosterPath)
		scraped.FrontImage = &posterURL
	}

	// Handle tags (genres)
	var tags []*models.ScrapedTag
	for _, genre := range details.Genres {
		tags = append(tags, &models.ScrapedTag{
			Name: genre.Name,
		})
	}
	if len(tags) > 0 {
		scraped.Tags = tags
	}

	// Handle duration if available
	if details.Runtime > 0 {
		duration := strconv.Itoa(details.Runtime)
		scraped.Duration = &duration
	}

	// Handle director
	for _, crew := range details.Credits.Crew {
		if crew.Job == "Director" {
			scraped.Director = &crew.Name
			break
		}
	}

	// Handle URLs
	if movie.ID > 0 {
		movieURL := fmt.Sprintf("https://www.themoviedb.org/movie/%d", movie.ID)
		scraped.URLs = []string{movieURL}
	}

	return scraped, nil
}

// ScrapeTV converts a TV show to the movie format that Stash expects
func (s *Scraper) ScrapeTV(ctx context.Context, name string) (*models.ScrapedMovie, error) {
	logger.Infof("Searching TMDB for TV show: %s", name)

	results, err := s.client.SearchTV(name)
	if err != nil {
		return nil, fmt.Errorf("TMDB search failed: %v", err)
	}

	if len(results.Results) == 0 {
		return nil, fmt.Errorf("no TV shows found matching: %s", name)
	}

	show := results.Results[0]
	details, err := s.client.GetTVShowDetails(show.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get TV show details: %v", err)
	}

	// Convert to ScrapedMovie
	scraped := &models.ScrapedMovie{}

	// Set basic information
	scraped.Name = &show.Name
	scraped.Date = &show.FirstAirDate
	scraped.Synopsis = &show.Overview

	// Handle Rating
	if show.VoteAverage > 0 {
		rating := strconv.FormatFloat(show.VoteAverage, 'f', 2, 64)
		scraped.Rating = &rating
	}

	// Handle image
	if show.PosterPath != "" {
		posterURL := s.client.GetImageURL(show.PosterPath)
		scraped.FrontImage = &posterURL
	}

	// Handle tags (genres)
	var tags []*models.ScrapedTag
	for _, genre := range details.Genres {
		tags = append(tags, &models.ScrapedTag{
			Name: genre.Name,
		})
	}
	if len(tags) > 0 {
		scraped.Tags = tags
	}

	// Handle URLs
	if show.ID > 0 {
		showURL := fmt.Sprintf("https://www.themoviedb.org/tv/%d", show.ID)
		scraped.URLs = []string{showURL}
	}

	return scraped, nil
}
