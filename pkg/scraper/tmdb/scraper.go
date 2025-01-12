package tmdb

import (
	"context"
	"fmt"
	"strconv"
	"strings"

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

// Movie scraping
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

	return s.convertToScrapedMovie(details), nil
}

func (s *Scraper) ScrapeMovieURL(ctx context.Context, url string) (*models.ScrapedMovie, error) {
	movieID, err := extractIDFromURL(url)
	if err != nil {
		return nil, err
	}

	details, err := s.client.GetMovieDetails(movieID)
	if err != nil {
		return nil, fmt.Errorf("failed to get movie details: %v", err)
	}

	return s.convertToScrapedMovie(details), nil
}

// TV Show scraping
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

	return s.convertToScrapedTV(details), nil
}

func (s *Scraper) ScrapeTVURL(ctx context.Context, url string) (*models.ScrapedMovie, error) {
	tvID, err := extractIDFromURL(url)
	if err != nil {
		return nil, err
	}

	details, err := s.client.GetTVShowDetails(tvID)
	if err != nil {
		return nil, fmt.Errorf("failed to get TV show details: %v", err)
	}

	return s.convertToScrapedTV(details), nil
}

// Helper functions
func (s *Scraper) convertToScrapedMovie(details *MovieDetails) *models.ScrapedMovie {
	scrapedMovie := &models.ScrapedMovie{
		Name:     &details.Title,
		Date:     &details.ReleaseDate,
		Synopsis: &details.Overview,
	}

	if details.PosterPath != "" {
		posterURL := s.client.GetImageURL(details.PosterPath)
		scrapedMovie.FrontImage = &posterURL
	}

	// Add genres as tags
	var tags []string
	for _, genre := range details.Genres {
		tags = append(tags, genre.Name)
	}
	if len(tags) > 0 {
		scrapedMovie.Tags = &models.ScrapedTags{Values: tags}
	}

	// Add cast and crew
	var performers []models.ScrapedPerformer
	for _, cast := range details.Credits.Cast {
		name := cast.Name
		role := cast.Character
		performers = append(performers, models.ScrapedPerformer{
			Name:  &name,
			Notes: &role,
		})
	}
	if len(performers) > 0 {
		scrapedMovie.Performers = &performers
	}

	// Add director
	for _, crew := range details.Credits.Crew {
		if crew.Job == "Director" {
			director := crew.Name
			scrapedMovie.Director = &director
			break
		}
	}

	return scrapedMovie
}

func (s *Scraper) convertToScrapedTV(details *TVShowDetails) *models.ScrapedMovie {
	scrapedTV := &models.ScrapedMovie{
		Name:     &details.Name,
		Date:     &details.FirstAirDate,
		Synopsis: &details.Overview,
	}

	if details.PosterPath != "" {
		posterURL := s.client.GetImageURL(details.PosterPath)
		scrapedTV.FrontImage = &posterURL
	}

	// Add genres as tags
	var tags []string
	for _, genre := range details.Genres {
		tags = append(tags, genre.Name)
	}
	if len(tags) > 0 {
		scrapedTV.Tags = &models.ScrapedTags{Values: tags}
	}

	// Add cast and crew
	var performers []models.ScrapedPerformer
	for _, cast := range details.Credits.Cast {
		name := cast.Name
		role := cast.Character
		performers = append(performers, models.ScrapedPerformer{
			Name:  &name,
			Notes: &role,
		})
	}
	if len(performers) > 0 {
		scrapedTV.Performers = &performers
	}

	return scrapedTV
}

func extractIDFromURL(url string) (int, error) {
	parts := strings.Split(url, "/")
	if len(parts) == 0 {
		return 0, fmt.Errorf("invalid URL format")
	}

	// Get the last part of the URL and convert to int
	idStr := parts[len(parts)-1]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("invalid ID in URL")
	}

	return id, nil
}

func (s *Scraper) SupportedScrapes() []models.ScrapeType {
	return []models.ScrapeType{
		models.ScrapeTypeMovie,
		models.ScrapeTypeScene, // For TV episodes
	}
}
